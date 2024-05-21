package logger

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

var ProcessorCount *uint

type Loggable interface {
	UUID() string
	Log() Log
}

type SystemMetric struct {
	ProcessorCount uint
	ThreadCount    uint
	CPUPercentage  float64
	CPUTime        float64
	HeapMemoryUsed uint
	HeapMemoryFree uint
}

func (m *SystemMetric) GetSystemMetric() string {
	return fmt.Sprintf(`
	mcproxy_sys_processor_count %d
	mcproxy_sys_threads_count %d
	mcproxy_sys_cpu_percentage %.2f
	mcproxy_sys_cpu_time %.2f
	mcproxy_sys_memory_heap_used %d
	mcproxy_sys_memory_heap_free %d

	`, m.ProcessorCount, m.ThreadCount, m.CPUPercentage, m.CPUTime, m.HeapMemoryUsed, m.HeapMemoryFree)
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
	return fmt.Sprintf(`
	mcproxy_network_client_packet_tx %d
	mcproxy_network_client_packet_rx %d
	mcproxy_network_server_packet_tx %d
	mcproxy_network_server_packet_rx %d
	mcproxy_network_client_data_tx %d
	mcproxy_network_client_data_rx %d
	mcproxy_network_server_data_tx %d
	mcproxy_network_server_data_rx %d

	`, m.ClientPacketTx, m.ClientPacketRx, m.ServerPacketTx, m.ServerPacketRx, m.ClientDataTx, m.ClientDataRx, m.ServerDataTx, m.ServerDataRx)
}

type ProxyMetric struct {
	PlayerGetStatus uint
	PlayerLogin     uint
	PlayerPlaying   uint
}

func (m *ProxyMetric) GetProxyMetric() string {
	return fmt.Sprintf(`
	mcproxy_proxy_get_status %d
	mcproxy_proxy_login %d
	mcproxy_proxy_playing %d

	`, m.PlayerGetStatus, m.PlayerLogin, m.PlayerPlaying)
}

type ErrorMetric struct {
	AcceptFailed            uint
	HandshakeFailed         uint
	PacketDeserializeFailed uint
	HostnameResolveFailed   uint
	ServerConnectFailed     uint
}

func (m *ErrorMetric) Sum(a ErrorMetric) {
	m.AcceptFailed += a.AcceptFailed
	m.HandshakeFailed += a.HandshakeFailed
	m.PacketDeserializeFailed += m.PacketDeserializeFailed
	m.HostnameResolveFailed += a.HostnameResolveFailed
	m.ServerConnectFailed += a.ServerConnectFailed
}

func (m *ErrorMetric) GetErrorMetric() string {
	return fmt.Sprintf(`
	mcproxy_error_accept_failed %d
	mcproxy_error_hanhshake_failed %d
	mcproxy_error_deserialization_failed %d
	mcproxy_error_hostname_resolve_failed %d
	mcproxy_error_server_connect_failed %d

	`, m.AcceptFailed, m.HandshakeFailed, m.PacketDeserializeFailed, m.HostnameResolveFailed, m.ServerConnectFailed)
}

type Metric struct {
	SystemMetric
	NetworkMetric
	ErrorMetric
	ProxyMetric
}

func (m *Metric) GetMetric() string {
	d := (m.GetNetworkMetric() + m.GetErrorMetric() + m.GetSystemMetric() + m.GetProxyMetric())
	d = strings.ReplaceAll(d, "\t", "")
	d = strings.ReplaceAll(d, "\n\n", "\n")
	return d[1:]
}

type MetricCollecter struct {
	LogEntities map[string]Loggable
	mutex       *sync.Mutex
}

func (c *MetricCollecter) systemMetric() (SystemMetric, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	if ProcessorCount == nil {
		cpuNum := uint(runtime.NumCPU())
		ProcessorCount = &cpuNum
	}
	pid := int32(os.Getpid())
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Printf("[Metric] Error getting cpu profile: %v", err)
		return SystemMetric{}, err
	}
	cpuTimes, errProc := proc.Times()
	if errProc != nil {
		log.Printf("[Metric] Error getting cpu utilization: %v", errProc)
		return SystemMetric{}, err
	}
	overallCPUTimes, errCpu := cpu.Times(false)
	if errCpu != nil {
		log.Printf("[Metric] Error getting cpu utilization: %v", errCpu)
		return SystemMetric{}, err
	}
	if len(overallCPUTimes) > 0 {
		totalCPUDelta := overallCPUTimes[0].Total() - cpuTimes.Total()
		processCPUUsage := (cpuTimes.Total() / totalCPUDelta) * 100
		return SystemMetric{
			ThreadCount:    uint(runtime.NumGoroutine()),
			HeapMemoryUsed: uint(mem.HeapAlloc),
			HeapMemoryFree: uint(mem.HeapIdle),
			ProcessorCount: *ProcessorCount,
			CPUPercentage:  processCPUUsage,
			CPUTime:        cpuTimes.Total(),
		}, nil
	}
	return SystemMetric{}, errors.New("Failed to get system metric")
}

func (c *MetricCollecter) Register(l Loggable) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.LogEntities[l.UUID()] = l
}

func (c *MetricCollecter) Unregister(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.LogEntities, key)
}

func (c *MetricCollecter) Collect() (Metric, error) {
	// logs := make([]Log, 1024)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sys, err := c.systemMetric()
	if err != nil {
		return Metric{}, err
	}
	metric := Metric{SystemMetric: sys}
	for _, e := range c.LogEntities {
		log := e.Log()
		if (ProxyMetric{}) != log.ProxyMetric {
			metric.ProxyMetric = log.ProxyMetric
		}
		metric.NetworkMetric.Sum(log.NetworkMetric)
		metric.ErrorMetric.Sum(log.ErrorMetric)
	}
	return metric, nil
}

func NewMetricCollector() *MetricCollecter {
	return &MetricCollecter{LogEntities: make(map[string]Loggable), mutex: &sync.Mutex{}}
}

type IMetricExporter interface {
	Serve()
}

type PrometheusExporter struct {
	MetricCollecter *MetricCollecter
	Port            uint
	IMetricExporter
}

func (e *PrometheusExporter) handler(w http.ResponseWriter, r *http.Request) {
	// Set the content type to plain text
	w.Header().Set("Content-Type", "text/plain")

	// Write the text response
	l, err := e.MetricCollecter.Collect()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		log.Printf("[Prometheus Exporter] Error: %v", err)
		return
	}
	// fmt.Println(l.GetMetric())
	w.Write([]byte(l.GetMetric()))
}

func (e *PrometheusExporter) Serve() {
	<-time.After(1 * time.Second)
	http.HandleFunc("/metrics", e.handler)
	log.Printf("[Prometheus Exporter] Starting server at :%d", e.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", e.Port), nil); err != nil {
		log.Fatalf("[Prometheus Exporter] Error starting server: %v", err)
	}
}
