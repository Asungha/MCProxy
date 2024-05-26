package metric

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
)

type PrometheusFormatter struct {
	data string
}

func (f *PrometheusFormatter) Add(metricName string, value string, filter map[string]string) *PrometheusFormatter {
	if len(filter) == 0 {
		f.data += fmt.Sprintf("\n%s %s", metricName, value)
	} else {
		buf := "{"
		for k, v := range filter {
			if buf != "{" {
				buf += ","
			}
			buf += fmt.Sprintf(`%s="%s"`, k, v)
		}
		buf += "}"
		f.data += fmt.Sprintf("\n%s%s %s", metricName, buf, value)
	}
	return f
}

func (f *PrometheusFormatter) Get() string {
	return f.data
}

var ProcessorCount *uint

type Loggable interface {
	UUID() string
	Log() Log
}

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
	formatter := PrometheusFormatter{}
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
	formatter := PrometheusFormatter{}
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

type ProxyMetric struct {
	PlayerGetStatus int
	PlayerLogin     int
	PlayerPlaying   int
	Ping            int
	Connected       int
}

func (m *ProxyMetric) GetProxyMetric() string {
	filter := map[string]string{}
	formatter := PrometheusFormatter{}
	formatter.Add("mcproxy_proxy_get_status", strconv.FormatInt(int64(m.PlayerGetStatus), 10), filter)
	formatter.Add("mcproxy_proxy_login", strconv.FormatInt(int64(m.PlayerLogin), 10), filter)
	formatter.Add("mcproxy_proxy_ping", strconv.FormatInt(int64(m.Ping), 10), filter)
	formatter.Add("mcproxy_proxy_playing", strconv.FormatInt(int64(m.PlayerPlaying), 10), filter)
	formatter.Add("mcproxy_proxy_connected", strconv.FormatInt(int64(m.Connected), 10), filter)
	return formatter.Get()
}

func (m *ProxyMetric) Sum(a ProxyMetric) {
	m.PlayerGetStatus += a.PlayerGetStatus
	m.PlayerLogin += a.PlayerLogin
	m.PlayerPlaying += a.PlayerPlaying
	m.Connected += a.Connected
}

type ErrorMetric struct {
	AcceptFailed            uint
	HandshakeFailed         uint
	PacketDeserializeFailed uint
	HostnameResolveFailed   uint
	ServerConnectFailed     uint

	LogOverflow uint
}

func (m *ErrorMetric) Sum(a ErrorMetric) {
	m.AcceptFailed += a.AcceptFailed
	m.HandshakeFailed += a.HandshakeFailed
	m.PacketDeserializeFailed += m.PacketDeserializeFailed
	m.HostnameResolveFailed += a.HostnameResolveFailed
	m.ServerConnectFailed += a.ServerConnectFailed
	m.LogOverflow += a.LogOverflow
}

func (m *ErrorMetric) GetErrorMetric() string {
	filter := map[string]string{}
	formatter := PrometheusFormatter{}
	formatter.Add("mcproxy_error_accept_failed", strconv.FormatInt(int64(m.AcceptFailed), 10), filter)
	formatter.Add("mcproxy_error_hanhshake_failed", strconv.FormatInt(int64(m.HandshakeFailed), 10), filter)
	formatter.Add("mcproxy_error_deserialization_failed", strconv.FormatInt(int64(m.PacketDeserializeFailed), 10), filter)
	formatter.Add("mcproxy_error_hostname_resolve_failed", strconv.FormatInt(int64(m.HostnameResolveFailed), 10), filter)
	formatter.Add("mcproxy_error_server_connect_failed", strconv.FormatInt(int64(m.ServerConnectFailed), 10), filter)
	formatter.Add("mcproxy_error_log_overflow", strconv.FormatInt(int64(m.LogOverflow), 10), filter)
	return formatter.Get()
}

type Metric struct {
	SystemMetric
	NetworkMetric
	ErrorMetric
	ProxyMetric
	playerMetric map[string]PlayerMetric
}

type PlayerMetric struct {
	LoggedOut  bool
	PlayerName string
	IP         string
	Port       string
	LogginTime time.Time
	Playtime   time.Duration
	*NetworkMetric
	ErrorMetric
}

func NewPlayerMetric(addr string, name string) *PlayerMetric {
	s := strings.Split(addr, ":")
	return &PlayerMetric{
		LogginTime: time.Now(),
		PlayerName: name,
		IP:         s[0],
		Port:       s[1],
	}
}

func (m *PlayerMetric) GetPlayerMetric() string {
	if m.PlayerName == "" || m.IP == "" {
		return ""
	}
	filter := map[string]string{"player": m.PlayerName, "ip": m.IP}
	formatter := PrometheusFormatter{}
	formatter.Add("mcproxy_player_online", "1", filter)
	formatter.Add("mcproxy_player_playtime", fmt.Sprint(time.Since(m.LogginTime).Seconds()), filter)

	formatter.Add("mcproxy_player_error_accept_failed", fmt.Sprint(m.AcceptFailed), filter)
	formatter.Add("mcproxy_player_error_hanhshake_failed", fmt.Sprint(m.HandshakeFailed), filter)
	formatter.Add("mcproxy_player_error_deserialization_failed", fmt.Sprint(m.PacketDeserializeFailed), filter)
	formatter.Add("mcproxy_player_error_hostname_resolve_failed", fmt.Sprint(m.HostnameResolveFailed), filter)
	formatter.Add("mcproxy_player_error_server_connect_failed", fmt.Sprint(m.ServerConnectFailed), filter)
	formatter.Add("mcproxy_player_error_server_connect_failed", fmt.Sprint(m.ServerConnectFailed), filter)

	if m.NetworkMetric != nil {
		formatter.Add("mcproxy_player_network_client_packet_tx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketTx), 10), filter)
		formatter.Add("mcproxy_player_network_client_packet_rx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_server_packet_tx", strconv.FormatInt(int64(m.NetworkMetric.ServerPacketTx), 10), filter)
		formatter.Add("mcproxy_player_network_server_packet_rx", strconv.FormatInt(int64(m.NetworkMetric.ServerPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_client_data_tx", strconv.FormatInt(int64(m.NetworkMetric.ClientDataTx), 10), filter)
		formatter.Add("mcproxy_player_network_client_data_rx", strconv.FormatInt(int64(m.NetworkMetric.ClientPacketRx), 10), filter)
		formatter.Add("mcproxy_player_network_server_data_tx", strconv.FormatInt(int64(m.NetworkMetric.ServerDataTx), 10), filter)
		formatter.Add("mcproxy_player_network_server_data_rx", strconv.FormatInt(int64(m.NetworkMetric.ServerDataRx), 10), filter)
	}

	return formatter.Get()
}

func (m *Metric) GetMetric() string {
	m.ProxyMetric.PlayerPlaying = len(m.playerMetric) - 1
	d := (m.GetNetworkMetric() + m.GetErrorMetric() + m.GetSystemMetric() + m.GetProxyMetric())
	for _, playerMetric := range m.playerMetric {
		d += playerMetric.GetPlayerMetric()
	}
	d = strings.ReplaceAll(d, "\t", "")
	d = strings.ReplaceAll(d, "\n\n", "\n")
	return d
}

type MetricCollecter struct {
	LogEntities            map[string]Loggable
	mutex                  *sync.Mutex
	LastProxyCPUTime       float64
	LastSystemTotalCPUTime float64

	PushChannel      chan Log
	PushBuffer       []*Log
	PushedLog        int
	LogOverflowCount int
	pushMutex        *sync.Mutex

	lastMetric Metric
}

func (c *MetricCollecter) readPushedLog() {
	for {
		log := <-c.PushChannel
		// fmt.Println("Got log")s
		c.pushMutex.Lock()
		if len(c.PushBuffer) >= 80960 {
			c.PushBuffer = c.PushBuffer[1:]
		}
		c.PushBuffer = append(c.PushBuffer, &log)
		c.pushMutex.Unlock()
	}
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
		// systemTotal := overallCPUTimes[0].Total() - cpuTimes.Total()
		// processCPUUsage := (cpuTimes.Total() - c.lastCPUTime/diffSystem) * 100
		// c.lastCPUTime = cpuTimes.Total()
		proxyCPUTime := cpuTimes.Total()
		systemCPUTime := overallCPUTimes[0].Total()
		diffProxyCPUTime := proxyCPUTime - c.LastProxyCPUTime
		diffSystemCPUTime := systemCPUTime - c.LastSystemTotalCPUTime

		percentageProxyCPU := (diffProxyCPUTime / diffSystemCPUTime) * 100
		// percentageSystemCPU := ((diffSystemCPUTime - diffProxyCPUTime) / diffSystemCPUTime) * 100

		c.LastProxyCPUTime = proxyCPUTime
		c.LastSystemTotalCPUTime = systemCPUTime
		return SystemMetric{
			ThreadCount:        uint(runtime.NumGoroutine()),
			HeapMemoryUsed:     uint(mem.HeapAlloc),
			HeapMemoryFree:     uint(mem.HeapIdle),
			ProcessorCount:     *ProcessorCount,
			ProxyCPUPercentage: percentageProxyCPU, // CPU used by proxy
			// SystemCPUPercentage: percentageSystemCPU, // CPU used by the rest of the host system
			CPUTime: cpuTimes.Total(),
		}, nil
	}
	return SystemMetric{}, errors.New("Failed to get system metric")
}

func (c *MetricCollecter) Register(l Loggable) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	uuid := l.UUID()
	c.LogEntities[uuid] = l
	return uuid
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
	c.lastMetric.SystemMetric = sys
	c.lastMetric.playerMetric = map[string]PlayerMetric{}
	// if
	for _, e := range c.LogEntities {
		log := e.Log()
		// if (ProxyMetric{}) != log.ProxyMetric {
		// 	metric.ProxyMetric.Sum(log.ProxyMetric)
		// }
		c.lastMetric.ProxyMetric.Sum(log.ProxyMetric)
		c.lastMetric.NetworkMetric.Sum(log.NetworkMetric)
		c.lastMetric.ErrorMetric.Sum(log.ErrorMetric)
		c.lastMetric.playerMetric[log.PlayerName+log.IP] = log.PlayerMetric
	}
	c.pushMutex.Lock()
	log.Printf("[Metric Collector] Reading %d logs", len(c.PushBuffer))
	for _, log := range c.PushBuffer {
		if log == nil {
			continue
		}
		// fmt.Println(log)
		// log := e.Log()
		// if (ProxyMetric{}) != log.ProxyMetric {
		// 	metric.ProxyMetric.Sum(log.ProxyMetric)
		// }
		c.lastMetric.ProxyMetric.Sum(log.ProxyMetric)
		c.lastMetric.NetworkMetric.Sum(log.NetworkMetric)
		c.lastMetric.ErrorMetric.Sum(log.ErrorMetric)
		c.lastMetric.playerMetric[log.PlayerName+log.IP] = log.PlayerMetric
	}
	c.PushBuffer = []*Log{}
	c.pushMutex.Unlock()
	return c.lastMetric, nil
}

func (c *MetricCollecter) PushLog(log Log) error {
	done := make(chan bool)
	go func() {
		c.PushChannel <- log
		done <- true
	}()
	select {
	case <-time.After(time.Millisecond * 100):
		c.pushMutex.Lock()
		defer c.pushMutex.Unlock()
		c.LogOverflowCount++
		return errors.New("Log buffer overflow")
	case <-done:
		c.pushMutex.Lock()
		c.PushedLog += 1
		c.pushMutex.Unlock()
	}
	return nil
}

func NewMetricCollector() *MetricCollecter {
	cputime := 0.0
	overallCPUTimes, errCpu := cpu.Times(false)
	if errCpu == nil {
		// log.Printf("[Metric] Error getting cpu utilization: %v", errCpu)
		cputime = overallCPUTimes[0].Total()
	}
	mc := &MetricCollecter{
		LogEntities:            make(map[string]Loggable),
		mutex:                  &sync.Mutex{},
		LastProxyCPUTime:       0.0,
		LastSystemTotalCPUTime: cputime,
		PushChannel:            make(chan Log, 16),
		PushBuffer:             make([]*Log, 80960),
		pushMutex:              &sync.Mutex{},
	}
	go mc.readPushedLog()
	return mc
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
