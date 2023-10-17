package storage

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"sync"
)

type ConsulOption func(*ConsulOptions)

type ConsulOptions struct {
	Address string
}

func ConsulAddress(h string) ConsulOption {
	return func(options *ConsulOptions) {
		options.Address = h
	}
}

var (
	consulOnce     sync.Once
	consulInstance *Consul
)

type Consul struct {
	client *api.Client
}

func NewConsul(options ...ConsulOption) (*Consul, error) {
	var err error
	consulOnce.Do(func() {
		opt := ConsulOptions{}
		for _, setter := range options {
			setter(&opt)
		}

		conf := api.DefaultConfig()
		conf.Address = viper.GetString(opt.Address)
		client, e := api.NewClient(conf)
		if e != nil {
			err = e
			return
		}
		consulInstance = &Consul{client: client}
	})

	if err != nil {
		return nil, err
	}
	return consulInstance, nil
}

func (c *Consul) Put(key string, val any) error {

	byteData, _ := json.Marshal(val)
	kvPair := &api.KVPair{Key: key, Value: byteData}
	_, err := c.client.KV().Put(kvPair, nil)
	if err != nil {
		logger.Errorf("put pair to consul failed, %v", err)
		return err
	}
	return nil
}

func (c *Consul) Get(key string) ([]byte, error) {
	var err error
	kvPair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		logger.Errorf("get pair from consul failed, key: %s error: %v", key, err)
		return nil, err
	}
	return kvPair.Value, nil
}
