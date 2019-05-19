package etcdv3

import (
	"fmt"
	"github.com/dubbo/dubbo-go/registry"
	"golang.org/x/net/context"
	"sync"
)

type etcdEventListener struct {
	client         *etcdClient
	serviceMapLock sync.Mutex
	serviceMap     map[string]struct{}
	wg             sync.WaitGroup
	registry       *EtcdRegistry
}

func (l *etcdEventListener) listenServiceNodeEvent(path string) bool {
	l.wg.Add(1)
	defer l.wg.Done()
	for {
		keyEventCh, err := l.client.kv.Get(context.Background(), path)
		if err != nil {
			//log.Error("existW{key:%s} = error{%v}", path, err)
			return false
		}
        fmt.Println(keyEventCh)
		select {
		}
	}

	return false
}

func (l *etcdEventListener) handleetcdNodeEvent(Path string, children []string, conf registry.ServiceConfig) {

	newChildren, err := l.client.kv.Get(context.Background(),Path)
	if err != nil {
		return
	}
	var (

	)
	for _, n := range newChildren.Kvs {
		fmt.Println(n.Key)

	}


}

func (l *etcdEventListener) listenDirEvent(zkPath string, conf registry.ServiceConfig) {
	l.wg.Add(1)
	defer l.wg.Done()

	var (
		event     chan struct{}
	)
	event = make(chan struct{}, 4)
	defer close(event)
}

// this func is invoked by ZkConsumerRegistry::Registe/ZkConsumerRegistry::get/ZkConsumerRegistry::getListener
// registry.go:Listen -> listenServiceEvent -> listenDirEvent -> listenServiceNodeEvent
//                            |
//                            --------> listenServiceNodeEvent
func (l *etcdEventListener) listenServiceEvent(conf registry.ServiceConfig) {
	var (
		err        error
		etcdPath     string
	)

	etcdPath = fmt.Sprintf("/dubbo/%s/providers", conf.Service())

	l.serviceMapLock.Lock()
	_, ok := l.serviceMap[etcdPath]
	l.serviceMapLock.Unlock()
	if ok {
		//log.Warn("@zkPath %s has already been listened.", zkPath)
		return
	}

	l.serviceMapLock.Lock()
	l.serviceMap[etcdPath] = struct{}{}
	l.serviceMapLock.Unlock()

	//log.Info("listen dubbo provider path{%s} event and wait to get all provider zk nodes", zkPath)
	children, err := l.client.kv.Get(context.Background(), etcdPath)
	if err != nil {
		children = nil
	}
	//todo get children
	fmt.Println(children)

}

func (l *etcdEventListener) Next() (*registry.ServiceEvent, error) {
	for {
		select {
		}
	}
}

func (l *etcdEventListener) valid() bool {
	return true
}

func (l *etcdEventListener) Close() {
	l.registry.wg.Done()
	l.wg.Wait()
}

