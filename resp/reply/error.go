package reply

// UnknowErrReply 表示 Redis 中返回未知错误的回复。
type UnknowErrReply struct{}

var _ ErrorReply = &UnknowErrReply{}

var unknowErrBytes = []byte("-Err unknow\r\n")

func (r *UnknowErrReply) Error() string {
	return "Err unknow"
}

func (r *UnknowErrReply) ToBytes() []byte {
	return unknowErrBytes
}

func MakeUnknowErrReply() *UnknowErrReply {
	return &UnknowErrReply{}
}

// ArgNumErrReply 表示 Redis 中返回参数数量错误的回复。
type ArgNumErrReply struct {
	Cmd string
}

var _ ErrorReply = &ArgNumErrReply{}

func (r *ArgNumErrReply) Error() string {
	return "Err wrong number of arguments for '" + r.Cmd + "' command"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

// An SyntaxErrReply 表示 Redis 中返回语法错误的回复。
type SyntaxErrReply struct{}

var _ ErrorReply = &SyntaxErrReply{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

// WrongTypeErrReply 表示 Redis 中返回类型错误的回复。
type WrongTypeErrReply struct{}

var _ ErrorReply = &WrongTypeErrReply{}

func (r *WrongTypeErrReply) Error() string {
	return "Err wrong type"
}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return []byte("-Err wrong type\r\n")
}

func MakeWrongTypeErrReply() *WrongTypeErrReply {
	return &WrongTypeErrReply{}
}

// ProtocolErrReply 表示 Redis 中返回协议错误的回复。
type ProtocolErrReply struct{}

var _ ErrorReply = &ProtocolErrReply{}

func (r *ProtocolErrReply) Error() string {
	return "Err protocol error"
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-Err protocol error\r\n")
}

func MakeProtocolErrReply() *ProtocolErrReply {
	return &ProtocolErrReply{}
}
