package dict

import (
	"math/rand"
	"sync"

	"go-redis/tools"
)

type SyncDict struct {
	m   sync.Map
	len int
}

// Clear implements [Dict].
func (s *SyncDict) Clear() {
	s.m = sync.Map{}
	s.len = 0
}

// ForEach implements [Dict].
func (s *SyncDict) ForEach(consumer Consumer) {
	s.m.Range(func(key, value any) bool {
		return consumer(key.(string), value)
	})
}

// Get 返回指定 key 对应的值，如果不存在则返回 nil 和 false
func (s *SyncDict) Get(key string) (val any, exists bool) {
	return s.m.Load(key)
}

// Keys implements [Dict].
func (s *SyncDict) Keys() (result []string) {
	result = make([]string, s.Len())

	s.m.Range(func(key, _ any) bool {
		result = append(result, key.(string))
		return true
	})

	return result
}

// Len 返回字典中键值对的数量
func (s *SyncDict) Len() int {
	return s.len
}

// Put 插入或更新指定 key 对应的值，如果 key 已存在则返回 0，否则返回 1
func (s *SyncDict) Put(key string, val any) (result int) {
	_, existed := s.m.Load(key)
	s.m.Store(key, val)
	s.len++

	if existed {
		return 0
	}

	return 1
}

// PutIfAbsent 指定的key不存在时；插入 val
func (s *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, ok := s.m.Load(key)
	if ok {
		return 0
	}

	s.m.Store(key, val)
	s.len++
	return 1
}

// PutIfExists 指定的key存在时；插入 val
func (s *SyncDict) PutIfExists(key string, val any) (result int) {
	_, ok := s.m.Load(key)
	if ok {
		s.m.Store(key, val)
		s.len++
		return 1
	}

	return 0
}

// RandomDistinctKeys 随机获取字典中指定数量的不重复键
func (s *SyncDict) RandomDistinctKeys(limit int) []string {
	all := s.Keys()

	m := tools.SyncSet[string]{}

	if limit >= len(all) {
		return all
	}

	// 随机获取 limit 个键
	for i := 0; i < limit; i++ {

		// 生成一个随机索引
		index := rand.Intn(len(all))
		if m.Has(all[index]) {
			i-- // 如果已经存在，则重新生成索引
			continue
		}

		m.Store(all[index])
	}

	return m.Keys()
}

func (s *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, 0, limit)

	all := s.Keys()

	if limit >= len(all) {
		return all
	}

	// 随机获取 limit 个键
	for range limit {
		// 生成一个随机索引
		index := rand.Intn(len(all))
		result = append(result, all[index])
	}

	return result
}

func (s *SyncDict) Remove(key string) (result int) {
	_, ok := s.m.Load(key)
	s.m.Delete(key)

	if ok {
		s.len--
		return 1
	}

	return 0
}

var _ Dict = (*SyncDict)(nil)

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}
