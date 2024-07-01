package adaptor

import (

	// metric "mc_reverse_proxy/src/metric"
	config "mc_reverse_proxy/src/configuration/service"
	metricAdaptor "mc_reverse_proxy/src/metric/adaptor"
	metricDTO "mc_reverse_proxy/src/metric/dto"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyService "mc_reverse_proxy/src/proxy/service"
	statemachine "mc_reverse_proxy/src/statemachine/service"
	"mc_reverse_proxy/src/utils"
	"net"
	"sync"
)

const serviceName = "Game Proxy"

type Iproxy interface {
	ImplProxy()
	Serve()
}

type MinecraftProxy struct {
	MetricCollector *metricService.MetricService
	ProxyMetric     metricDTO.ProxyMetric
	ErrorMetric     metricDTO.ErrorMetric
	metricExporter  metricAdaptor.MetricAdaptor

	Repository proxyService.ServerRepositoryService

	Listener        *net.Listener
	routerLock      sync.Mutex
	threadWaitGroup *sync.WaitGroup
	init            bool

	configService *config.ConfigurationService
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
			p.metricExporter.Serve()
		}(p)
		p.init = true
	}
	for {
		statemachine := statemachine.NewNetworkStatemachine(p.configService, p.Listener, p.Repository, &p.ProxyMetric, p.MetricCollector)
		go statemachine.Run()
		<-statemachine.ClientConnected
		go func(uuid string) {
			<-statemachine.Ctx.Done()
		}("")
	}
}

func NewProxy(listenAddr string, metricService *metricService.MetricService, configService *config.ConfigurationService) (Iproxy, error) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	utils.FLog.Proxy("Accepting connection at %s", listenAddr)
	repo := proxyService.NewQLServerRepositoryService()
	err = repo.Load()
	if err != nil {
		return nil, err
	}
	proxy := &MinecraftProxy{Listener: &listener, threadWaitGroup: &sync.WaitGroup{}, MetricCollector: metricService, Repository: repo, configService: configService}
	proxy.MetricCollector.Register(proxy)
	return proxy, nil
}
