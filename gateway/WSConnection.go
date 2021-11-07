package gateway

import (
	"github.com/gorilla/websocket"
	"live-chat/common"
	"sync"
	"time"
)

type WSConnection struct {
	mutex             sync.Mutex
	connId            uint64
	wsSocket          *websocket.Conn
	inChan            chan *common.WSMessage
	outChan           chan *common.WSMessage
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time
	rooms             map[string]bool
}

func InitWSConnection(connId uint64, wsSocket *websocket.Conn) (wsConnetion *WSConnection) {
	wsConnetion = &WSConnection{
		wsSocket:          wsSocket,
		connId:            connId,
		inChan:            make(chan *common.WSMessage, G_config.WsInChannelSize),
		outChan:           make(chan *common.WSMessage, G_config.WsOutChannelSize),
		closeChan:         make(chan byte),
		lastHeartbeatTime: time.Now(),
		rooms:             make(map[string]bool),
	}

	go wsConnetion.readLoop()
	go wsConnetion.writeLoop()

	return
}

// 发送消息到outChan,让协程写出去
func (wsConnection *WSConnection) SendMessage(message *common.WSMessage) (err error) {
	select {
	case wsConnection.outChan <- message:
		SendMessageTotal_INCR() // 统计发送成功的条数
	case <-wsConnection.closeChan:
		err = common.ERR_CONNECTION_LOSS
	default:
		err = common.ERR_SEND_MESSAGE_FULL
		SendMessageFail_INCR() //
	}
	return
}

// inchan读消息
func (wsConnection *WSConnection) ReadMessage() (message *common.WSMessage, err error) {
	select {
	case message = <-wsConnection.inChan:
	case <-wsConnection.closeChan:
		err = common.ERR_CONNECTION_LOSS
	}
	return
}

func (wsConnection *WSConnection) readLoop() {
	var (
		msgType int
		msgData []byte
		message *common.WSMessage
		err     error
	)
	for {
		if msgType, msgData, err = wsConnection.wsSocket.ReadMessage(); err != nil {
			goto ERR
		}

		message = common.BuildWSMessage(msgType, msgData)

		select {
		case wsConnection.inChan <- message: // inchan 可写
		case <-wsConnection.closeChan: //closeChan 可读
			goto CLOSED
		}

	}

ERR:
	wsConnection.Close()
CLOSED:
}

func (wsConnection *WSConnection) writeLoop() {
	var (
		message *common.WSMessage
		err     error
	)
	for {
		select {
		case message = <-wsConnection.outChan:
			if err = wsConnection.wsSocket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				goto ERR
			}
		case <-wsConnection.closeChan:
			goto CLOSED
		}
	}
ERR:
	wsConnection.Close()
CLOSED:
}

// 关闭连接
func (wsConnection *WSConnection) Close() {
	wsConnection.wsSocket.Close()

	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()
	if !wsConnection.isClosed {
		wsConnection.isClosed = true
		close(wsConnection.closeChan)
	}
}

// 检查心跳(不需要太频繁)
func (wsConnection *WSConnection) IsAlive() bool {
	var (
		now = time.Now()
	)
	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()

	//连接已关闭 或者 太久没有心跳
	if wsConnection.isClosed || now.Sub(wsConnection.lastHeartbeatTime) > time.Duration(G_config.WsHeartbeatInterval)*time.Second {
		return false
	}
	return true
}

// 更新心跳值
func (wsConnection *WSConnection) KeepAlive() {
	var (
		now = time.Now()
	)
	wsConnection.mutex.Lock()
	defer wsConnection.mutex.Unlock()
	wsConnection.lastHeartbeatTime = now
}
