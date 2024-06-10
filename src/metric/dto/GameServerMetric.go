package dto

import (
	"fmt"
	proto "mc_reverse_proxy/src/control/controlProto"
	metricUtils "mc_reverse_proxy/src/metric/utils"
	"time"
)

type GameServerMetric struct {
	*proto.MetricData
}

func (m *GameServerMetric) GetGameServerMetric() string {
	if starttime == nil {
		t := time.Now()
		starttime = &t
	}
	filter := map[string]string{
		"service": m.ServerID,
	}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_server_system_tps", fmt.Sprintf("%.2f", m.Tps), filter)
	formatter.Add("mcproxy_server_system_cpu_percentage", fmt.Sprintf("%.2f", m.CpuUsage), filter)
	formatter.Add("mcproxy_server_system_heap_reserved", fmt.Sprintf("%d", m.MemoryMax), filter)
	formatter.Add("mcproxy_server_system_heap_used", fmt.Sprintf("%d", m.MemoryUsed), filter)

	formatter.Add("mcproxy_server_game_online", fmt.Sprintf("%d", m.OnlineCount), filter)
	return formatter.Get()
}
