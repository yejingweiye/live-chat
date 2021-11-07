package gateway

import (
	"live-chat/common"
	"sync"
)

// 直播房间
type Room struct {
	rwMutex sync.RWMutex
	roomId  string
	id2Conn map[uint64]*WSConnection // 房间多个连接
}

func InitRoom(roomId string) (room *Room) {
	room = &Room{
		roomId:  roomId,
		id2Conn: make(map[uint64]*WSConnection),
	}
	return
}

//房间添加连接
func (room *Room) Join(wsConn *WSConnection) (err error) {
	var (
		existed bool
	)

	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed = room.id2Conn[wsConn.connId]; !existed {
		err = common.ERR_NOT_IN_ROOM
		return
	}
	room.id2Conn[wsConn.connId] = wsConn
	return
}

// 删除房间连接
func (room *Room) Leave(wsConn *WSConnection) (err error) {
	var (
		existed bool
	)
	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()
	if _, existed = room.id2Conn[wsConn.connId]; !existed {
		err = common.ERR_NOT_IN_ROOM
		return
	}
	delete(room.id2Conn, wsConn.connId)
	return
}

// 房间连接数
func (room *Room) Count() int {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()
	return len(room.id2Conn)
}

//
func (room *Room) Push(wsMsg *common.WSMessage) {
	var (
		wsConn *WSConnection
	)
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	for _, wsConn = range room.id2Conn {
		wsConn.SendMessage(wsMsg) // 房间所有连接发送消息
	}
}
