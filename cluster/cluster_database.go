package cluster

import (
	"context"

	"go-redis/config"
	"go-redis/database"
	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/lib/consistenhash"

	pool "github.com/jolestar/go-commons-pool"
)

/**
* 三种执行模式
* - 本节点 比如 ping
* - 远程节点 比如 set、get key
* - 集群模式 比如 keys *、flushdb 需要在所有节点上执行
*
*
* */

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenhash.NodeMap
	peerConnection map[string]*pool.ObjectPool
	db             databaseImp.Database
}

// AfterClientClose implements [databaseImp.Database].
func (*ClusterDatabase) AfterClientClose(c resp.Connection) {
	panic("unimplemented")
}

// Close implements [databaseImp.Database].
func (c *ClusterDatabase) Close() {
	panic("unimplemented")
}

// Exec implements [databaseImp.Database].
func (c *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	panic("unimplemented")
}

func MakeClustrtDatabase() *ClusterDatabase {
	// Self:           "127.0.0.1:6378",
	// Peers:          []string{"127.0.0.1:6379", "127.0.0.1:6380"},
	//
	c := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database.MakeStandaloneDatabase(),
		peerPicker:     consistenhash.MakeNodeMap(nil),
		nodes:          config.Properties.Peers[:],
		peerConnection: make(map[string]*pool.ObjectPool),
	}

	c.peerPicker.AddNode(c.nodes...)

	for _, peer := range c.nodes {
		// 创建连接池
		c.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(context.Background(), &connectionFactory{Peer: peer})
	}

	return c
}

var _ databaseImp.Database = (*ClusterDatabase)(nil)
