package gateway

import (
	"live-chat/common"
	"sync"
)

type Bucket struct {
	rwMutex sync.RWMutex
	index   int                      //桶索引 第几个桶
	id2Conn map[uint64]*WSConnection // 连接列表(key=连接唯一ID)
	rooms   map[string]*Room         // 房间列表
}

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2Conn: make(map[uint64]*WSConnection),
		rooms:   make(map[string]*Room),
	}
	return
}

// 桶添加连接
func (bucket *Bucket) AddConn(wsConn *WSConnection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	bucket.id2Conn[wsConn.connId] = wsConn
}

// 桶删除连接
func (bucket *Bucket) DelConn(wsConn *WSConnection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	delete(bucket.id2Conn, wsConn.connId)
}

// 初始化房间，在桶加入房间，房间再添加连接
func (bucket *Bucket) JoinRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		existed bool
		room    *Room
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	// 找到房间
	if room, existed = bucket.rooms[roomId]; !existed { // 不存在，创建房间
		room = InitRoom(roomId)
		bucket.rooms[roomId] = room
		RoomCount_INCR()
	}

	// 加入房间
	err = room.Join(wsConn)
	return
}

// 连接删 无连接桶房间再删
func (bucket *Bucket) LeaveRoom(roomId string, wsConn *WSConnection) (err error) {
	var (
		existed bool
		room    *Room
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	if room, existed = bucket.rooms[roomId]; !existed {
		err = common.ERR_NOT_IN_ROOM
		return
	}

	err = room.Leave(wsConn) // 删除连接

	// 房间连接为空，则删除
	if room.Count() == 0 {
		delete(bucket.rooms, roomId)
		RoomCount_DESC()
	}
	return
}

// 推送消息给Bucket内所有用户
func (bucket *Bucket) PushAll(wsMsg *common.WSMessage) {
	var (
		wsConn *WSConnection
	)
	bucket.rwMutex.RLock()
	defer bucket.rwMutex.RUnlock()

	for _, wsConn = range bucket.id2Conn {
		wsConn.SendMessage(wsMsg)
	}
}

// 推送给某房间的所有用户
func (bucket *Bucket) PushRoom(roomId string, wsMsg *common.WSMessage) {
	var (
		room    *Room
		existed bool
	)
	bucket.rwMutex.RLock()
	room, existed = bucket.rooms[roomId]
	bucket.rwMutex.RUnlock()

	// 房间不存在
	if !existed {
		return
	}

	// 向房间做推送
	room.Push(wsMsg)
}
