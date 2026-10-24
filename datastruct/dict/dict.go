// Package dict 表示一个字典数据
package dict

type Consumer func(key string, val any) bool

// Dict represents a dictionary data structure.
type Dict interface {
	// 根据 key 获取字典中的值
	Get(key string) (val any, exists bool)
	// 获取字典中键值对的数量
	Len() int
	// 设置字典中键值对的数量
	// 成功 返回 1，失败返回 0
	Put(key string, val any) (result int)
	// key 不存在时才设置值
	// 成功 返回 1，失败返回 0
	PutIfAbsent(key string, val any) (result int)
	// key 存在时才设置值
	// 成功 返回 1，失败返回 0
	PutIfExists(key string, val any) (result int)
	// 移除字典中指定的键值对
	// 成功 返回 1，失败返回 0
	Remove(key string) (result int)
	// 遍历字典中的所有键值对
	ForEach(consumer Consumer)
	// 获取字典中所有的键
	Keys() []string
	// 随机获取字典中指定数量的键
	RandomKeys(limit int) []string
	// 随机获取字典中指定数量的不重复键
	RandomDistinctKeys(limit int) []string
	// 清空字典中的所有键值对
	Clear()
}
