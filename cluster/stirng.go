package cluster

import (
	"fmt"
	"strings"

	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
)

// set k1 v1
func set(c *ClusterDatabase, client resp.Connection, db databaseImp.Database, args CmdLine) resp.Reply {
	peer := c.peerPicker.PickNode(string(args[1]))

	fmt.Printf("设置[%v],[%v]\n", string(args[1]), peer)

	if peer == c.self {
		return db.Exec(client, args)
	}

	localC, _ := c.getPeerClient(peer)
	l := strings.Split(peer, ":")

	return localC.Send(l[1], args)
}

func init() {
	RegisterCommand("set", set, 3)
}
