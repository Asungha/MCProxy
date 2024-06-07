package adaptor

import (
	"log"

	// metric "mc_reverse_proxy/src/metric"
	metricAdaptor "mc_reverse_proxy/src/metric/adaptor"
	metricDTO "mc_reverse_proxy/src/metric/dto"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	statemachine "mc_reverse_proxy/src/statemachine/service"
	"net"
	"sync"
)

type Iproxy interface {
	ImplProxy()
	Serve()
	// UseMetricExporter(string)
	metricService.Loggable
}

type MinecraftProxy struct {
	Listener *net.Listener

	routerLock sync.Mutex

	threadWaitGroup *sync.WaitGroup

	MetricCollector *metricService.MetricService

	ProxyMetric metricDTO.ProxyMetric
	ErrorMetric metricDTO.ErrorMetric

	metricExporter metricAdaptor.MetricAdaptor

	init bool

	Repository proxyService.ServerRepositoryService

	// logger *Logger.Logger
}

func (p *MinecraftProxy) ImplProxy() {}

func (p *MinecraftProxy) UUID() string {
	return "0"
}

func (p *MinecraftProxy) Log() metricDTO.Log {
	return metricDTO.Log{
		ProxyMetric: p.ProxyMetric,
		ErrorMetric: p.ErrorMetric,
	}
}
func (p *MinecraftProxy) Serve() {
	if p.metricExporter != nil && !p.init {
		go func(p *MinecraftProxy) {
			defer log.Println("[metricExporter] Thread exit")
			p.metricExporter.Serve()
		}(p)
		p.init = true
	}
	for {
		statemachine := statemachine.NewNetworkStatemachine(p.Listener, p.Repository, &p.ProxyMetric, p.MetricCollector)
		go statemachine.Run()
		<-statemachine.ClientConnected

		log.Printf("[Proxy] Connection between proxy and client established")
		// uuid := p.MetricCollector.Register(statemachine)
		go func(uuid string) {
			<-statemachine.Ctx.Done()
			log.Printf("[sm manager] cleanup")
			// p.MetricCollector.Unregister(uuid)
		}("")
	}
}

func NewProxy(listenAddr string, metricService *metricService.MetricService) (Iproxy, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("[Proxy] Accepting connection at %s", listenAddr)
	repo := proxyService.NewQLServerRepositoryService()
	err = repo.Load()
	if err != nil {
		return nil, err
	}
	proxy := &MinecraftProxy{Listener: &listener, threadWaitGroup: &sync.WaitGroup{}, MetricCollector: metricService, Repository: repo}
	proxy.MetricCollector.Register(proxy)
	return proxy, nil
}
