package service

import (
	"errors"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"

	controlService "mc_reverse_proxy/src/control/service"
	metricDTO "mc_reverse_proxy/src/metric/dto"
)

type Loggable interface {
	UUID() string
	Log() metricDTO.Log
}

type MetricService struct {
	metricChannel chan controlService.EventData

	PushChannel chan metricDTO.Log
	PushBuffer  []*metricDTO.Log
	lastMetric  metricDTO.Metric

	LogEntities            map[string]Loggable
	mutex                  *sync.Mutex
	LastProxyCPUTime       float64
	LastSystemTotalCPUTime float64
	PushedLog              int
	LogOverflowCount       int
	pushMutex              *sync.Mutex
}

func (c *MetricService) readPushedLog() {
	for {
		log := <-c.PushChannel
		c.pushMutex.Lock()
		if len(c.PushBuffer) >= 80960 {
			c.PushBuffer = c.PushBuffer[1:]
		}
		c.PushBuffer = append(c.PushBuffer, &log)
		c.pushMutex.Unlock()
	}
}

func (c *MetricService) readLogEvent() {
	for {
		log := <-c.metricChannel
		c.pushMutex.Lock()
		if len(c.PushBuffer) >= 80960 {
			c.PushBuffer = c.PushBuffer[1:]
		}
		c.PushBuffer = append(c.PushBuffer, &metricDTO.Log{
			GameServerMetric: &metricDTO.GameServerMetric{
				MetricData: log.MetricData,
			},
		})
		c.pushMutex.Unlock()
	}
}

func (c *MetricService) systemMetric() (metricDTO.SystemMetric, error) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	if metricDTO.ProcessorCount == nil {
		cpuNum := uint(runtime.NumCPU())
		metricDTO.ProcessorCount = &cpuNum
	}
	pid := int32(os.Getpid())
	proc, err := process.NewProcess(pid)
	if err != nil {
		log.Printf("[Metric] Error getting cpu profile: %v", err)
		return metricDTO.SystemMetric{}, err
	}
	cpuTimes, errProc := proc.Times()
	if errProc != nil {
		log.Printf("[Metric] Error getting cpu utilization: %v", errProc)
		return metricDTO.SystemMetric{}, err
	}
	overallCPUTimes, errCpu := cpu.Times(false)
	if errCpu != nil {
		log.Printf("[Metric] Error getting cpu utilization: %v", errCpu)
		return metricDTO.SystemMetric{}, err
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
		sm := metricDTO.NewSystemMetric()
		sm.ThreadCount = uint(runtime.NumGoroutine())
		sm.HeapMemoryUsed = uint(mem.HeapAlloc)
		sm.HeapMemoryFree = uint(mem.HeapIdle)
		sm.ProcessorCount = *metricDTO.ProcessorCount
		sm.ProxyCPUPercentage = percentageProxyCPU // CPU used by proxy
		sm.CPUTime = cpuTimes.Total()
		return sm, nil
	}
	return metricDTO.SystemMetric{}, errors.New("Failed to get system metric")
}

func (c *MetricService) Register(l Loggable) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	uuid := l.UUID()
	c.LogEntities[uuid] = l
	return uuid
}

func (c *MetricService) Unregister(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.LogEntities, key)
}

func (c *MetricService) Collect() (metricDTO.Metric, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	sys, err := c.systemMetric()
	if err != nil {
		return metricDTO.Metric{}, err
	}
	c.lastMetric.SystemMetric = sys
	for _, e := range c.LogEntities {
		log := e.Log()
		c.lastMetric.ProxyMetric.Sum(log.ProxyMetric)
		c.lastMetric.NetworkMetric.Sum(log.NetworkMetric)
		c.lastMetric.ErrorMetric.Sum(log.ErrorMetric)
		if log.GameServerMetric != nil {
			c.lastMetric.GameServerMetric[log.GameServerMetric.ServerID] = log.GameServerMetric
		}
	}
	c.pushMutex.Lock()
	for _, log := range c.PushBuffer {
		if log == nil {
			continue
		}
		c.lastMetric.ProxyMetric.Sum(log.ProxyMetric)
		c.lastMetric.NetworkMetric.Sum(log.NetworkMetric)
		c.lastMetric.ErrorMetric.Sum(log.ErrorMetric)
		if log.GameServerMetric != nil {
			c.lastMetric.GameServerMetric[log.GameServerMetric.ServerID] = log.GameServerMetric
		}
	}
	c.PushBuffer = []*metricDTO.Log{}
	c.pushMutex.Unlock()
	return c.lastMetric, nil
}

func (c *MetricService) PushLog(log metricDTO.Log) error {
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

func NewMetricService(eventService *controlService.EventService) *MetricService {
	cputime := 0.0
	overallCPUTimes, errCpu := cpu.Times(false)
	if errCpu == nil {
		cputime = overallCPUTimes[0].Total()
	}
	_, metricChannel := eventService.Subscribe("metric")
	mc := &MetricService{
		LogEntities:            make(map[string]Loggable),
		mutex:                  &sync.Mutex{},
		LastProxyCPUTime:       0.0,
		LastSystemTotalCPUTime: cputime,
		PushChannel:            make(chan metricDTO.Log, 16),
		PushBuffer:             make([]*metricDTO.Log, 80960),
		pushMutex:              &sync.Mutex{},
		metricChannel:          metricChannel,
		lastMetric:             metricDTO.Metric{GameServerMetric: make(map[string]*metricDTO.GameServerMetric)},
	}
	go mc.readPushedLog()
	go mc.readLogEvent()
	return mc
}
