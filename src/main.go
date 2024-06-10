package main

import (
	"encoding/json"
	"fmt"
	controlAdaptor "mc_reverse_proxy/src/control/adaptor"
	controlService "mc_reverse_proxy/src/control/service"
	metricAdaptor "mc_reverse_proxy/src/metric/adaptor"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyAdaptor "mc_reverse_proxy/src/proxy/adaptor"
	proxyService "mc_reverse_proxy/src/proxy/service"
	webui "mc_reverse_proxy/src/webui"
	"os"
)

func ReadConfig() map[string]string {
	config := map[string]string{}
	config_file, err := os.Open("config.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to open config file: %v", err))
	}
	defer config_file.Close()

	decoder := json.NewDecoder(config_file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode config file: %v", err))
	}
	return config
}

func ConfigGRPCApi(eventService *controlService.EventService, address string) {

}

func main() {
	os.Remove("./.cf0fd1a71be9dd15298c0c29bf1e6b13a4433b34") // for hot reload
	config := ReadConfig()

	event := controlService.NewEventService()
	metricService := metricService.NewMetricService(event)

	isBackendStarted := false

	p, err := proxyAdaptor.NewProxy(config["listen_address"], metricService)
	if err != nil {
		panic(err.Error())
	}

	if grpcIaAddr, ok := config["grpc_metric_address"]; ok {
		GRPCService := controlAdaptor.NewGRPCControlCenter(grpcIaAddr, event)
		go GRPCService.Serve()
	}

	if v, ok := config["prometheus_address"]; ok {
		metricExporter := &metricAdaptor.PrometheusAdaptor{MetricCollecter: metricService, ListenAddress: v}
		go metricExporter.Serve()
	}

	if ba, ok := config["http_api_address"]; ok {
		backend := webui.NewHTTPBackend(ba, metricService, p.(*proxyAdaptor.MinecraftProxy).Repository)
		go backend.Serve()
		isBackendStarted = true
	}

	if fa, ok := config["webui_address"]; ok {
		if !isBackendStarted {
			panic("Webui frontend required http api to work. Add config 'http_api_address' in the config.json with appropiated address.")
		}
		frontend := webui.NewHTTPFrontend(fa)
		go frontend.Serve()
	}

	defer p.(*proxyAdaptor.MinecraftProxy).Repository.(proxyService.UpdatableRepositoryService).Destroy()

	p.Serve()
}
