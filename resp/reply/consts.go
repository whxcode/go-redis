package reply

// PongReply 表示 Redis 中返回 PONG 的回复。
type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}

// OkReply 表示 Redis 中返回 OK 的回复。
type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

func MakeOkReply() *OkReply {
	return &OkReply{}
}

// NullBulkReply 表示 Redis 中返回 nil 的回复。
type NullBulkReply struct{}

func (r *NullBulkReply) ToBytes() []byte {
	// 特殊回复 nil 回复；并不是空字符串
	return []byte("$-1\r\n")
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

// EmptyBulkReply 表示 Redis 中返回空字符串的回复。
type EmptyBulkReply struct{}

func (r *EmptyBulkReply) ToBytes() []byte {
	// 回复空字符串
	return []byte("$0\r\n\r\n")
}

func MakeEmptyBulkReply() *EmptyBulkReply {
	return &EmptyBulkReply{}
}

// EmptyMultiBulkReply 表示 Redis 中返回空数组的回复。
type EmptyMultiBulkReply struct{}

var emptyMultiBulkBytres = []byte("*0\r\n")

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytres
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

// NoReply 表示 Redis 中没有返回值的回复。
type NoReply struct{}

func (r *NoReply) ToBytes() []byte {
	return []byte("")
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
