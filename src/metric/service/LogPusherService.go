package service

import (
	"fmt"
	metricDTO "mc_reverse_proxy/src/metric/dto"
)

type LogPusher struct {
	Collector MetricService
}

func (p *LogPusher) PushErrorMetric(log metricDTO.ErrorMetric) error {
	fmt.Println("ErrorMetric")
	return p.Collector.PushLog(metricDTO.Log{ErrorMetric: log})
}

func (p *LogPusher) PushProxyMetric(log metricDTO.ProxyMetric) error {
	fmt.Println("ProxyMetric")
	return p.Collector.PushLog(metricDTO.Log{ProxyMetric: log})
}

func (p *LogPusher) PushPlayerMetric(log metricDTO.PlayerMetric) error {
	fmt.Println("PlayerMetric")
	return p.Collector.PushLog(metricDTO.Log{PlayerMetric: log})
}

func (p *LogPusher) PushNetworkMetric(log metricDTO.NetworkMetric) error {
	fmt.Println("NetworkMetric")
	return p.Collector.PushLog(metricDTO.Log{NetworkMetric: log})
}
