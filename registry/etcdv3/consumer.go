package etcdv3

import (
	"fmt"
	"golang.org/x/net/context"
)

import (
	log "github.com/AlexStocks/log4go"
	jerrors "github.com/juju/errors"
)

import (
	"github.com/dubbo/dubbo-go/plugins"
	"github.com/dubbo/dubbo-go/registry"
)

// name: service@protocol
func (r *EtcdRegistry) GetService(conf registry.ServiceConfig) ([]registry.ServiceURL, error) {

	var (
		err         error
		dubboPath   string
		listener    *etcdEventListener
		serviceURL  registry.ServiceURL
		serviceConf registry.ServiceConfig
	)
	r.listenerLock.Lock()
	listener = r.listener
	r.listenerLock.Unlock()

	if listener != nil {
		listener.listenServiceEvent(conf)
	}



	dubboPath = fmt.Sprintf("/dubbo/%s/providers", conf.Service())
	//err = r.validateZookeeperClient()
	if err != nil {
		return nil, jerrors.Trace(err)
	}
	r.cltLock.Lock()
	nodes, err := r.client.kv.Get(context.Background(), dubboPath)
	r.cltLock.Unlock()
	if err != nil {
		log.Warn("getChildren(dubboPath{%s}) = error{%v}", dubboPath, err)
		return nil, jerrors.Trace(err)
	}

	var listenerServiceMap = make(map[string]registry.ServiceURL)
	for _, n := range nodes {

		serviceURL, err = plugins.DefaultServiceURL()(n)
		if err != nil {
			log.Error("NewDefaultServiceURL({%s}) = error{%v}", n, err)
			continue
		}
		if !serviceConf.ServiceEqual(serviceURL) {
			log.Warn("serviceURL{%s} is not compatible with ServiceConfig{%#v}", serviceURL, serviceConf)
			continue
		}

		_, ok := listenerServiceMap[serviceURL.Query().Get(serviceURL.Location())]
		if !ok {
			listenerServiceMap[serviceURL.Location()] = serviceURL
			continue
		}
	}

	var services []registry.ServiceURL
	for _, service := range listenerServiceMap {
		services = append(services, service)
	}

	return services, nil
}
