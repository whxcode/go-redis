package reply

import (
	"bytes"
	"strconv"

	"go-redis/interface/resp"
)

// 导入接口所在的包

var (
	// 空-nil 回复
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
	argLen := len(r.Args)
	var buf bytes.Buffer

	buf.WriteString("*")
	buf.WriteString(strconv.Itoa(argLen))
	buf.WriteString(CRLF)

	for _, arg := range r.Args {
		if arg == nil {
			buf.WriteString("$-1")
			buf.WriteString(CRLF)
		} else {
			buf.WriteString("$")
			buf.WriteString(strconv.Itoa(len(arg)))
			buf.WriteString(CRLF)
			buf.WriteString(string(arg))
			buf.WriteString(CRLF)
		}
	}

	return buf.Bytes()
}

func MakeMutiBulkReply(args [][]byte) *MutiBulkReply {
	return &MutiBulkReply{
		Args: args,
	}
}

// StatusReply 表示 Redis 中的状态回复。
type StatusReply struct {
	status string
}

var _ resp.Reply = &StatusReply{}

func (r *StatusReply) ToBytes() []byte {
	return []byte("+" + r.status + CRLF)
}

func MakeStatusReply(status string) *StatusReply {
	return &StatusReply{
		status: status,
	}
}

// IntReply 表示 Redis 中的整数回复。
type IntReply struct {
	status int
}

var _ resp.Reply = &IntReply{}

func (r *IntReply) ToBytes() []byte {
	return []byte(":" + strconv.Itoa(r.status) + CRLF)
}

func MakeIntReply(status int) *IntReply {
	return &IntReply{
		status: status,
	}
}

// StandarErrReply 表示 Redis 中的标准错误回复。
type StandarErrReply struct {
	status string
}

var _ resp.Reply = &StandarErrReply{}

func (r *StandarErrReply) ToBytes() []byte {
	return []byte("-" + r.status + CRLF)
}

func MakeStandarErrReply(status string) *StandarErrReply {
	return &StandarErrReply{
		status: status,
	}
}

// IsErrReply 判断给定的 resp.Reply 是否是一个错误回复。
func IsErrReply(r resp.Reply) bool {
	return r.ToBytes()[0] == '-'
}
