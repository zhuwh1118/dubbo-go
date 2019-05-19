package etcdv3

import (
	"sync"
	"time"
)

type EtcdRegistryConfig struct {
	Address    []string      `required:"true" yaml:"address"  json:"address,omitempty"`
	UserName   string        `yaml:"user_name" json:"user_name,omitempty"`
	Password   string        `yaml:"password" json:"password,omitempty"`
	TimeoutStr string        `yaml:"timeout" default:"5s" json:"timeout,omitempty"` // unit: second
	Timeout    time.Duration `yaml:"-"  json:"-"`
}

type EtcdRegistry struct {
	birth int64          // time of file birth, seconds since Epoch; 0 if unknown
	wg    sync.WaitGroup // wg+done for zk restart
	done  chan struct{}

	cltLock  sync.Mutex
	client   *etcdClient


	listenerLock sync.Mutex
	listener     *etcdEventListener

	//for provider
	zkPath map[string]int // key = protocol://ip:port/interface
}