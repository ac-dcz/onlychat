package sd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

const DefaultNetwork = "tcp"

type Service struct {
	Key     string     `json:"key"`
	Network string     `json:"network"`
	Addr    string     `json:"addr"`
	TTLOpt  *TTLOption `json:"ttlopt"`
}

func (s *Service) ID() string {
	return fmt.Sprintf("%s(%s@%s)", s.Key, s.Network, s.Addr)
}

func (s *Service) Encode() (string, error) {
	buff := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buff).Encode(s); err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (s *Service) Decode(data []byte) error {
	return json.Unmarshal(data, s)
}

type TTLOption struct {
	HeartBeat time.Duration `json:"heartbeat"`
	TTL       time.Duration `json:"ttl"`
}

const (
	HeartBeatTime = time.Second * 3
	TTLTime       = time.Second * 10
)

var DefaultTTLOption = &TTLOption{
	HeartBeat: HeartBeatTime,
	TTL:       TTLTime,
}

func NewTTLOption(heartbeat, ttl time.Duration) *TTLOption {
	opt := &TTLOption{}
	if heartbeat < HeartBeatTime {
		heartbeat = HeartBeatTime
	}
	if ttl < TTLTime {
		ttl = TTLTime
	}
	if heartbeat > ttl {
		ttl = heartbeat
	}
	opt.HeartBeat, opt.TTL = heartbeat, ttl
	return opt
}

type Registrar interface {
	Register(context.Context, *Service) error
	UnRegister(context.Context, *Service) error
}
