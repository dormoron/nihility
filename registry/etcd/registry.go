package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nothingZero/nihility/registry"
	client3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type Registry struct {
	c    *client3.Client
	sess *concurrency.Session
}

func InitRegister(c *client3.Client) (*Registry, error) {
	sess, err := concurrency.NewSession(c)
	if err != nil {
		return nil, err
	}
	return &Registry{
		c:    c,
		sess: sess,
	}, nil
}

func (r *Registry) Register(ctx context.Context, si registry.ServiceInstance) error {
	val, err := json.Marshal(si)
	if err != nil {
		return err
	}
	_, err = r.c.Put(ctx, r.instanceKey(si), string(val), client3.WithLease(r.sess.Lease()))
	return err
}

func (r *Registry) UnRegister(ctx context.Context, si registry.ServiceInstance) error {
	_, err := r.c.Delete(ctx, r.instanceKey(si))
	return err
}

func (r *Registry) ListServices(ctx context.Context, serviceName string) ([]registry.ServiceInstance, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Registry) Subscribe(serviceName string) (<-chan registry.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Registry) Close() error {
	err := r.sess.Close()
	return err
}

func (r *Registry) instanceKey(si registry.ServiceInstance) string {
	return fmt.Sprintf("/nihility/%s/%s", si.Name, si.Address)
}
