package cluster

import (
	"context"
	"errors"

	"go-redis/client"

	pool "github.com/jolestar/go-commons-pool"
)

type connectionFactory struct {
	// 链接的地址。
	Peer string
}

// ActivateObject implements [pool.PooledObjectFactory].
func (c *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

// DestroyObject implements [pool.PooledObjectFactory].
func (c *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	customer, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("failed to close client")
	}

	customer.Close()

	return nil
}

// MakeObject implements [pool.PooledObjectFactory].
func (c *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	customer := client.MakeClient(c.Peer)
	customer.Start()

	return pool.NewPooledObject(customer), nil
}

// PassivateObject implements [pool.PooledObjectFactory].
func (c *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

// ValidateObject implements [pool.PooledObjectFactory].
func (c *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

var _ pool.PooledObjectFactory = (*connectionFactory)(nil)
