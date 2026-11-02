// Package consistenhash 实现一致性hash
package consistenhash

import (
	"hash/crc32"
	"sort"
)

type HashFunc func(data []byte) uint32

type NodeMap struct {
	// 实现一致性hash的 hash 函数
	hashFunc HashFunc
	// 各个节点的 hash 值
	nodeHashes []int // 123 -- 4534  -- 132413 --> 开到前面 123 第0个节点
	// 节点的hash值与节点的映射关系
	nodehashMap map[int]string
}

func (n *NodeMap) IsEmpty() bool {
	return len(n.nodeHashes) == 0
}

func (n *NodeMap) Len() int {
	return len(n.nodeHashes)
}

func (n *NodeMap) hash(data string) int {
	h := n.hashFunc([]byte(data))
	/*
		字节3 字节2 字节1 字节0
		(最高)            (最低)
		var x uint32 = 0x12345678
		topByte := byte(x >> 24)   // = 0x12
	*/

	return int(h) // = 0x1234
}

// AddNode 添加节点
func (n *NodeMap) AddNode(nodeAddr ...string) {
	for _, addr := range nodeAddr {
		hash := n.hash(addr)
		n.nodeHashes = append(n.nodeHashes, hash)
		n.nodehashMap[hash] = addr
	}

	sort.Ints(n.nodeHashes)
}

// PickNode 根据 key 选择节点
func (n *NodeMap) PickNode(key string) string {
	top := n.hash(key)
	// 找到第一个大于等于 top 的节点
	idx := sort.Search(len(n.nodeHashes), func(i int) bool {
		return n.nodeHashes[i] >= top
	})

	if idx == len(n.nodeHashes) {
		idx = 0
	}

	// 获取节点地址
	return n.nodehashMap[n.nodeHashes[idx]]
}

func MakeNodeMap(hashFunc HashFunc) *NodeMap {
	m := &NodeMap{
		hashFunc:    hashFunc,
		nodeHashes:  []int{},
		nodehashMap: make(map[int]string),
	}

	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}

	return m
}
