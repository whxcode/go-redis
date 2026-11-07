package cluster

import (
	"strings"

	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
)

// set k1 v1
func set(c *ClusterDatabase,
	client resp.Connection,
	db databaseImp.Database,
	args CmdLine,
	originMsg []byte,
) resp.Reply {
	peer := c.peerPicker.PickNode(string(args[1]))

	if peer == c.self {
		return db.Exec(client, args, nil)
	}

	localC, _ := c.getPeerClient(peer)

	l := strings.Split(peer, ":")

	return localC.Send(l[1], originMsg)
}

func init() {
	RegisterCommand("set", set, 3)
}
