package etcdv3

import (
	"context"
	"errors"
	"math/rand"
	"onlychat/common/sd"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Client interface {
	GetService(ctx context.Context, key string) (*sd.Service, error)

	Register(context.Context, *sd.Service) error
	UnRegister(context.Context, *sd.Service) error
	Close() error
}

type BalanceType int

const (
	Randomness BalanceType = iota
	RoundRobin
)

var (
	ErrInvaildKey  = errors.New("service key is invaild")
	ErrInvaildAddr = errors.New("service addr is invaild")

	ErrClosed = errors.New("client has been closed")
)

type Option func(opt *ClientOption)

type ClientOption struct {
	endpoints        []string
	dialTimeout      time.Duration
	keepAliveTime    time.Duration
	keepAliveTimeout time.Duration
	balance          BalanceType
}

func WithDialTimeout(d time.Duration) Option {
	return func(opt *ClientOption) {
		opt.dialTimeout = d
	}
}

func WithKeepAlive(t, timeout time.Duration) Option {
	return func(opt *ClientOption) {
		opt.keepAliveTime = t
		opt.keepAliveTimeout = timeout
	}
}

func WithLoadBalance(typ BalanceType) Option {
	return func(opt *ClientOption) {
		opt.balance = typ
	}
}

type client struct {
	cli    *clientv3.Client
	opt    *ClientOption
	closed bool
}

func NewClient(endpoints []string, opts ...Option) (Client, error) {
	opt := &ClientOption{
		endpoints: endpoints,
	}
	for _, o := range opts {
		o(opt)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            opt.endpoints,
		DialTimeout:          opt.dialTimeout,
		DialKeepAliveTime:    opt.keepAliveTime,
		DialKeepAliveTimeout: opt.keepAliveTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &client{
		cli: cli,
		opt: opt,
	}, nil
}

func (c *client) GetService(ctx context.Context, key string) (*sd.Service, error) {
	if resp, err := c.cli.KV.Get(ctx, key, clientv3.WithPrefix()); err != nil {
		return nil, err
	} else {
		var ret []*sd.Service
		for _, kv := range resp.Kvs {
			s := &sd.Service{}
			err := s.Decode(kv.Value)
			if err != nil {
				return nil, err
			}
			ret = append(ret, s)
		}
		return c.loadbalance(ret), nil
	}
}

func (c *client) loadbalance(srvcs []*sd.Service) *sd.Service {
	switch c.opt.balance {
	case Randomness:
		return srvcs[rand.Intn(len(srvcs))]
	case RoundRobin: //TODO: "Randomness"
		return srvcs[rand.Intn(len(srvcs))]
	}
	return nil
}

func (c *client) Register(ctx context.Context, srvc *sd.Service) error {
	if srvc.Key == "" {
		return ErrInvaildKey
	}
	if srvc.Addr == "" {
		return ErrInvaildAddr
	}
	if srvc.Network == "" {
		srvc.Network = sd.DefaultNetwork
	}
	if srvc.TTLOpt == nil {
		srvc.TTLOpt = sd.DefaultTTLOption
	}

	val, err := srvc.Encode()
	if err != nil {
		return err
	}

	leaseResp, err := c.cli.Lease.Grant(ctx, int64(srvc.TTLOpt.TTL.Seconds()))
	if err != nil {
		return err
	}
	_, err = c.cli.Put(ctx, srvc.ID(), val, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return err
	}

	hbch, err := c.cli.Lease.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case p := <-hbch:
				if p == nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *client) UnRegister(ctx context.Context, srvc *sd.Service) error {
	_, err := c.cli.Delete(ctx, srvc.ID(), clientv3.WithIgnoreLease())
	if err != nil {
		return err
	}
	return nil
}

func (c *client) Close() error {
	if c.closed {
		return ErrClosed
	}
	return c.cli.Close()
}
