package dto

import (
	"fmt"
	metricUtils "mc_reverse_proxy/src/metric/utils"
	"strconv"
	"time"
)

var ProcessorCount *uint
var starttime *time.Time

type SystemMetric struct {
	StartTime          time.Time
	ProcessorCount     uint
	ThreadCount        uint
	ProxyCPUPercentage float64
	// SystemCPUPercentage float64
	CPUTime        float64
	HeapMemoryUsed uint
	HeapMemoryFree uint
}

func (m *SystemMetric) GetSystemMetric() string {
	if starttime == nil {
		t := time.Now()
		starttime = &t
	}
	filter := map[string]string{}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_sys_uptime", strconv.FormatInt(int64(time.Since(*starttime).Seconds()), 10), filter)
	formatter.Add("mcproxy_sys_processor_count", strconv.FormatInt(int64(m.ProcessorCount), 10), filter)
	formatter.Add("mcproxy_sys_threads_count", strconv.FormatInt(int64(m.ThreadCount), 10), filter)
	formatter.Add("mcproxy_sys_proxy_cpu_percentage", fmt.Sprint(m.ProxyCPUPercentage), filter)
	// formatter.Add("mcproxy_sys_cpu_percentage", fmt.Sprint(m.SystemCPUPercentage), filter)
	formatter.Add("mcproxy_sys_cpu_time", fmt.Sprint(m.CPUTime), filter)
	formatter.Add("mcproxy_sys_memory_heap_used", strconv.FormatInt(int64(m.HeapMemoryUsed), 10), filter)
	formatter.Add("mcproxy_sys_memory_heap_free", strconv.FormatInt(int64(m.HeapMemoryFree), 10), filter)
	return formatter.Get()
}

type NetworkMetric struct {
	ClientPacketTx uint
	ClientPacketRx uint
	ServerPacketTx uint
	ServerPacketRx uint

	ClientDataTx uint
	ClientDataRx uint
	ServerDataTx uint
	ServerDataRx uint
}

func (m *NetworkMetric) Sum(a NetworkMetric) {
	m.ClientDataRx += a.ClientDataRx
	m.ClientDataTx += a.ClientDataTx
	m.ClientPacketRx += a.ClientPacketRx
	m.ClientPacketTx += a.ClientPacketTx
	m.ServerDataRx += a.ServerDataRx
	m.ServerDataTx += a.ServerDataTx
	m.ServerPacketRx += a.ServerPacketRx
	m.ServerPacketTx += a.ServerPacketTx
}

func (m *NetworkMetric) GetNetworkMetric() string {
	filter := map[string]string{}
	formatter := metricUtils.PrometheusFormatter{}
	formatter.Add("mcproxy_network_client_packet_tx", strconv.FormatInt(int64(m.ClientPacketTx), 10), filter)
	formatter.Add("mcproxy_network_client_packet_rx", strconv.FormatInt(int64(m.ClientPacketRx), 10), filter)
	formatter.Add("mcproxy_network_server_packet_tx", strconv.FormatInt(int64(m.ServerPacketTx), 10), filter)
	formatter.Add("mcproxy_network_server_packet_rx", strconv.FormatInt(int64(m.ServerPacketRx), 10), filter)
	formatter.Add("mcproxy_network_client_data_tx", strconv.FormatInt(int64(m.ClientDataTx), 10), filter)
	formatter.Add("mcproxy_network_client_data_rx", strconv.FormatInt(int64(m.ClientPacketRx), 10), filter)
	formatter.Add("mcproxy_network_server_data_tx", strconv.FormatInt(int64(m.ServerDataTx), 10), filter)
	formatter.Add("mcproxy_network_server_data_rx", strconv.FormatInt(int64(m.ServerDataRx), 10), filter)
	return formatter.Get()
}
