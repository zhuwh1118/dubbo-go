package etcdv3

import (
	log "github.com/AlexStocks/log4go"
	etcd "github.com/etcd-io/etcd/clientv3"
	jerrors "github.com/juju/errors"
	"github.com/samuel/go-zookeeper/zk"
	"golang.org/x/net/context"
	"path"
	"strings"
	"sync"
	"time"
)

type etcdClient struct {
	etcdAddrs  []string
	timeout    time.Duration
	name       string
	client     *etcd.Client
	kv         etcd.KV
	ctx        context.Context
	sync.Mutex // for conn
}

func newEtcdClient(name string, etcdAddrs []string, timeout time.Duration) (*etcdClient, error) {
	var (
		e   *etcdClient
		cfg etcd.Config
		err error
	)

	e = &etcdClient{
		etcdAddrs: etcdAddrs,
		name:      name,
		timeout:   timeout,
		ctx:       context.Background(),
	}

	cfg = etcd.Config{
		Endpoints:   e.etcdAddrs,
		DialTimeout: e.timeout,
	}

	e.client, err = etcd.New(cfg)
	if err != nil {
		return nil, jerrors.Annotatef(err, "etcd.New(etcdAddrs:%+v)", etcdAddrs)
	}

	e.kv = etcd.NewKV(e.client)
	return e, nil
}

func (e *etcdClient) Create(basePath string) error {
	var (
		err     error
		tmpPath string
	)

	log.Debug("etcdkv.Put(basePath{%s})", basePath)
	for _, str := range strings.Split(basePath, "/")[1:] {
		tmpPath = path.Join(tmpPath, "/", str)
		e.Lock()
		if e.kv != nil {
			_, err = e.kv.Put(e.ctx, tmpPath, "")
		}
		e.Unlock()
		if err != nil {
			if err == zk.ErrNodeExists {
				log.Error("zk.create(\"%s\") exists\n", tmpPath)
			} else {
				log.Error("zk.create(\"%s\") error(%v)\n", tmpPath, jerrors.ErrorStack(err))
				return jerrors.Annotatef(err, "zk.Create(path:%s)", basePath)
			}
		}
	}

	return nil
}

func (e *etcdClient) Delete(basePath string) error {
	return nil
}
