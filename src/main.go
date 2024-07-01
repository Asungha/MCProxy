package main

import (
	configService "mc_reverse_proxy/src/configuration/service"
	controlAdaptor "mc_reverse_proxy/src/control/adaptor"
	controlService "mc_reverse_proxy/src/control/service"
	metricAdaptor "mc_reverse_proxy/src/metric/adaptor"
	metricService "mc_reverse_proxy/src/metric/service"
	packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
	proxyAdaptor "mc_reverse_proxy/src/proxy/adaptor"
	proxyService "mc_reverse_proxy/src/proxy/service"
	utils "mc_reverse_proxy/src/utils"
	webui "mc_reverse_proxy/src/webui"
	"os"
)

func main() {
	os.Remove("./.cf0fd1a71be9dd15298c0c29bf1e6b13a4433b34") // for hot reload
	config, err := configService.NewConfigurationService("config.json")
	if err != nil {
		panic(err)
	}

	event := controlService.NewEventService(8)
	metricService := metricService.NewMetricService(event)

	p, err := proxyAdaptor.NewProxy(config.ServerAddress, metricService, config)
	if err != nil {
		panic(err.Error())
	}

	if config.GRPCAddress != "" {
		GRPCServer := controlAdaptor.NewGRPCControlCenter(config.GRPCAddress, event)
		go GRPCServer.Serve()
	}

	if config.PrometheusAddress != "" {
		metricExporter := &metricAdaptor.PrometheusAdaptor{MetricCollecter: metricService, ListenAddress: config.PrometheusAddress}
		go metricExporter.Serve()
	}

	if config.HTTPApiAddress != "" {
		httpBackend := webui.NewHTTPBackend(config.HTTPApiAddress, metricService, p.(*proxyAdaptor.MinecraftProxy).Repository, event)
		go httpBackend.Serve()
	}

	if config.WebuiAddress != "" {
		httpFrontend := webui.NewHTTPFrontend(config.WebuiAddress, config.HTTPHostname)
		go httpFrontend.Serve()
	}

	defer p.(*proxyAdaptor.MinecraftProxy).Repository.(proxyService.UpdatableRepositoryService).Destroy()

	if config.LoggerMongoDBName != "" {
		err := packetLoggerService.InitPacketLogger(config)
		if err != nil {
			// log.Printf("Packet logger error : %v", err)
			utils.FLogErr.PacketLogger("Packet logger error: %v", err)
		}
	}
	p.Serve()
}
