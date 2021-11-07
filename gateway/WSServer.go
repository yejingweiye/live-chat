package gateway

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net"
	_ "net"
	"net/http"
	"strconv"
	_ "strconv"
	"sync/atomic"
	"time"
	_ "time"
)

// Websocket服务端,
type WSServer struct {
	server    *http.Server
	curConnId uint64
}

var (
	G_WSServer *WSServer
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许跨域请求
		},
	}
)

func handleConnect(res http.ResponseWriter, req *http.Request) {
	var (
		err      error
		wsSocket *websocket.Conn
		connId   uint64
		wsConn   *WSConnection
	)
	// 由http请求升级为websocket
	if wsSocket, err = wsUpgrader.Upgrade(res, req, nil); err != nil {
		return
	}

	connId = atomic.AddUint64(&G_WSServer.curConnId, 1)

	wsConn = InitWSConnection(connId, wsSocket) // 初始化收发监听

	fmt.Println("socket connId: ", connId)

	// 开始处理websocket消息
	wsConn.WSHandle()
}

func InitWSServer() (err error) {
	var (
		mux      *http.ServeMux
		server   *http.Server
		listener net.Listener
	)

	mux = http.NewServeMux()
	mux.HandleFunc("/connect", handleConnect)

	server = &http.Server{
		ReadTimeout:  time.Duration(G_config.WsReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.WsWriteTimeout) * time.Millisecond,
		Handler:      mux, //路由
	}

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.WsPort)); err != nil {
		return
	}

	G_WSServer = &WSServer{
		server:    server,
		curConnId: uint64(time.Now().Unix()),
	}

	// 启动服务
	go server.Serve(listener)
	return

}
