package etcdv3

import (
	"context"
	"onlychat/common/sd"
)

type Registrar struct {
	client Client
}

func NewRegistrar(client Client) *Registrar {
	return &Registrar{
		client: client,
	}
}

func (r *Registrar) Register(ctx context.Context, s *sd.Service) error {
	return r.client.Register(ctx, s)
}

func (r *Registrar) UnRegister(ctx context.Context, s *sd.Service) error {
	return r.client.UnRegister(ctx, s)
}
