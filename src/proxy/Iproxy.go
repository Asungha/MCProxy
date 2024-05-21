package proxy

import (
	"encoding/json"
	"log"
	metric "mc_reverse_proxy/src/logger"
	state "mc_reverse_proxy/src/state"
	"net"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Iproxy interface {
	ImplProxy()

	Serve()
	GetMC() *metric.MetricCollecter
	UseMetricExporter(uint)
	// metric.Loggable
}

type Proxy struct {
	Listener *net.Listener

	routerLock sync.Mutex

	threadWaitGroup *sync.WaitGroup

	MetricCollector *metric.MetricCollecter

	ProxyMetric metric.ProxyMetric
	ErrorMetric metric.ErrorMetric

	metricExporter metric.IMetricExporter

	init bool

	// logger *Logger.Logger
}

func (p *Proxy) ImplProxy() {}

func (p *Proxy) UUID() string {
	return "0"
}

func (p *Proxy) Log() metric.Log {
	return metric.Log{
		ProxyMetric: p.ProxyMetric,
		ErrorMetric: p.ErrorMetric,
	}
}

func (p *Proxy) UseMetricExporter(port uint) {
	p.metricExporter = &metric.PrometheusExporter{MetricCollecter: p.MetricCollector, Port: port}
}

var serverlist map[string]map[string]string

func GetServerList() map[string]map[string]string {
	if serverlist == nil {
		host_file, err := os.Open("host.json")
		if err != nil {
			log.Fatalf("Failed to open host config file: %v", err)
		}
		defer host_file.Close()

		host := make(map[string]map[string]string)
		decoder := json.NewDecoder(host_file)
		err = decoder.Decode(&host)
		if err != nil {
			log.Fatalf("Failed to decode config file: %v", err)
		}

		backends := map[string]map[string]string{}
		for k, v := range host {
			backends[k] = make(map[string]string)
			backends[k]["target"] = v["ip"] + ":" + v["port"]
			backends[k]["hostname"] = v["hostname"]
		}
		serverlist = backends
	}
	return serverlist
}

func (p *Proxy) GetMC() *metric.MetricCollecter {
	return p.MetricCollector
}

func (p *Proxy) Serve() {
	defer func() {
		log.Printf("[Proxy] Cleanup session")
		runtime.GC()
		// if r := recover(); r != nil {
		// 	log.Printf("[Proxy] panic: ", r)
		// 	return
		// }
	}()

	if p.metricExporter != nil && !p.init {
		go func(p *Proxy) {
			p.metricExporter.Serve()
		}(p)
		p.init = true
	}

	statemachine := state.NewStateMachine(p.Listener, GetServerList(), &p.ErrorMetric, &p.ProxyMetric)
	err := statemachine.Run() // Block until someone connected
	if err != nil {
		log.Printf("[Proxy] Connection accept failed: %v", err)
		// p.StateMachine.Destroy()
		return
	}

	log.Printf("[Proxy] Connection between proxy and client established")
	p.MetricCollector.Register(statemachine)
	// p.threadWaitGroup.Add(1)
	go func(_sm state.IStateMachine, startTime time.Time) {
		// defer p.threadWaitGroup.Done()
		for {
			switch _sm.Transition() {
			case state.STATUS_OK:
				continue
			default:
				log.Printf("[Proxy] Connection Terminated after %v", time.Since(startTime))
				_sm.Destroy()
				return
			}
		}
	}(statemachine, time.Now())
}

func NewProxy(port string) (Iproxy, error) {
	config := map[string]string{}
	config_file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer config_file.Close()

	decoder := json.NewDecoder(config_file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	listenAddr := config["listen"] + ":" + config["port"]
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("[Proxy] Accepting connection at %s", listenAddr)
	proxy := &Proxy{Listener: &listener, threadWaitGroup: &sync.WaitGroup{}, MetricCollector: metric.NewMetricCollector()}
	proxy.MetricCollector.Register(proxy)
	if v, ok := config["prometheus_port"]; ok {
		port, err := strconv.ParseInt(v, 10, 16)
		if err != nil {
			proxy.UseMetricExporter(8080) // default
		} else {
			proxy.UseMetricExporter(uint(port))
		}

	}
	return proxy, nil
}
