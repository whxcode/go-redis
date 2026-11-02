// Package consistenhash 实现一致性hash
package consistenhash

import (
	"testing"
)

func TestNodeMap(t *testing.T) {
	m := MakeNodeMap(nil)
	nodes := map[string]map[string]string{
		"001": {},
		"002": {},
	}

	m.AddNode("001")
	m.AddNode("002")

	t.Logf("len: %d", m.Len())
	t.Logf("nodes: %v", m.nodeHashes)

	nodes[m.PickNode("key1")]["key1"] = "value1"
	nodes[m.PickNode("key0")]["key0"] = "value0"
	nodes[m.PickNode("key2")]["key2"] = "value2"

	t.Logf("nodes: %v", nodes)
}
