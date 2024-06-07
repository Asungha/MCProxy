package main

import (
	"encoding/json"
	"fmt"
	metricAdaptor "mc_reverse_proxy/src/metric/adaptor"
	metricService "mc_reverse_proxy/src/metric/service"
	proxyAdaptor "mc_reverse_proxy/src/proxy/adaptor"
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

func main() {
	config := ReadConfig()
	metricService := metricService.NewMetricService()
	p, err := proxyAdaptor.NewProxy(config["listen_address"], metricService)
	if err != nil {
		panic(err.Error())
	}

	if v, ok := config["prometheus_address"]; ok {
		metricExporter := &metricAdaptor.PrometheusAdaptor{MetricCollecter: metricService, ListenAddress: v}
		go metricExporter.Serve()
	}

	if v, ok := config["webui_address"]; ok {
		webui := webui.NewWebUI(metricService, p.(*proxyAdaptor.MinecraftProxy).Repository)
		go webui.Serve(v)
	}

	p.Serve()
}
