package resp

type Connection interface {
	// 给客户端回复消息
	Write(bytes []byte) error

	// 获取当前使用的 DB
	GetDBIndex() int

	// 切换当前的数据库
	SelectDB(db int)
}
