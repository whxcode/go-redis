package reply

import (
	"strconv"

	"go-redis/interface/resp"
)

// 导入接口所在的包

var (
	nullBulkReplyBytes = []byte("$-1")
	CRLF               = "\r\n"
)

type ErrorReply interface {
	Error() string
	ToBytes() []byte
}

// BulkReply 表示 Redis 中返回的批量回复。
type BulkReply struct {
	Arg []byte // eg "moody" => "$5\r\nmoody\r\n"
}

var _ resp.Reply = &BulkReply{}

func (r *BulkReply) ToBytes() []byte {
	if len(r.Arg) == 0 {
		return nullBulkReplyBytes
	}

	return []byte("$" + strconv.Itoa(len(r.Arg)) + CRLF + string(r.Arg) + CRLF)
}

func MakeBulkReply(arg []byte) *BulkReply {
	return &BulkReply{
		Arg: arg,
	}
}

// MutiBulkReply 表示 Redis 中返回的多批量回复。
type MutiBulkReply struct {
	Args [][]byte
}

var _ resp.Reply = &MutiBulkReply{}

func (r *MutiBulkReply) ToBytes() []byte {
	size := len(r.Args)

	if size == 0 {
		return nullBulkReplyBytes
	}

	res := make([]byte, 0, 100)
	res = append(res, "*"+strconv.Itoa(size)+CRLF...)

	for _, arg := range r.Args {
		res = append(res, "$"+CRLF+strconv.Itoa(len(arg))+CRLF...)
	}

	return res
}

func MakeMutiBulkReply(args [][]byte) *MutiBulkReply {
	return &MutiBulkReply{
		Args: args,
	}
}
