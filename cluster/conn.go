package cluster

import (
	"context"
	"errors"

	"go-redis/client"
)

// 在链接池拿到一个客户端
func (c *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {
	pool, ok := c.peerConnection[peer]

	if !ok {
		return nil, errors.New("peer not found")
	}

	obj, err := pool.BorrowObject(context.Background())
	if err != nil {
		return nil, err
	}

	customer, ok := obj.(*client.Client)

	if !ok {
		return nil, errors.New("failed to cast to client")
	}

	return customer, nil
}

func (c *ClusterDatabase) returnPeerClient(peer string, customer *client.Client) error {
	pool, ok := c.peerConnection[peer]

	if !ok {
		return errors.New("peer not found")
	}

	return pool.ReturnObject(context.Background(), customer)
}
