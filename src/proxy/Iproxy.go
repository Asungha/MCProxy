package proxy

import (
	"encoding/json"
	"log"
	Logger "mc_reverse_proxy/src/logger"
	state "mc_reverse_proxy/src/state"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

type Iproxy interface {
	ImplProxy()

	Serve()
}

type Proxy struct {
	Listener *net.Listener

	routerLock sync.Mutex

	threadWaitGroup *sync.WaitGroup

	logger *Logger.Logger
}

func (p *Proxy) ImplProxy() {}

// func (p *Proxy) WaitForConnection() error {
// 	clientConn, err := (*GetListener(p.ListenAddress)).Accept()
// 	if err != nil {
// 		log.Printf("Failed to accept client connection: %v", err)
// 		return err
// 	}
// 	log.Printf("[Inbound] Initiating connection between %s and %s", clientConn.RemoteAddr().String(), "localhost:25565")
// 	// clientConn.SetDeadline(time.Now().Add(5 * time.Second))
// 	p.client = &clientConn
// 	return nil
// }

// func (p *Proxy) SetUpStreams(host string) error {
// 	upstreamConn, err := net.Dial("tcp", host)
// 	if err != nil {
// 		log.Printf("Failed to connect to upstream server: %v", err)
// 		return err
// 	}
// 	log.Printf("[Outbound] Initiated connection between %s and %s", "proxy", upstreamConn.RemoteAddr().String())
// 	// upstreamConn.SetDeadline(time.Now().Add(5 * time.Second))
// 	p.upstream = &upstreamConn
// 	return nil
// }

// func (p *Proxy) preConditionCheck() error {
// 	if p.client == nil {
// 		return errors.New("Client connection not established")
// 	}
// 	if p.upstream == nil {
// 		return errors.New("Upstream connection not established")
// 	}
// 	return nil
// }

// func (p *Proxy) listenClient(wg *sync.WaitGroup, ctx context.Context, cancleFunc context.CancelCauseFunc, output chan []byte) error {
// 	defer wg.Done()
// 	for {
// 		buf := make([]byte, 512)
// 		n, err := (*p.client).Read(buf)
// 		// log.Printf("Reading from client: %x", buf[:n])
// 		if err != nil {
// 			// log.Printf("Failed to read from client connection: %v", err)
// 			cancleFunc(err)
// 			buf = nil
// 			return err
// 		}
// 		data := buf[:n]
// 		// packet := pac.NewPacket(&pac.Raw{})
// 		// err = packet.Decode(&data, len(data))
// 		// if err != nil {
// 		// 	log.Printf("Failed to decode packet: %v", err)
// 		// 	cancleFunc(err)
// 		// 	return err
// 		// }
// 		// p.routerLock.Lock()
// 		output <- data
// 		// p.routerLock.Unlock()
// 		buf = nil
// 	}
// }

// func (p *Proxy) listenUpstream(wg *sync.WaitGroup, ctx context.Context, cancleFunc context.CancelCauseFunc, output chan []byte) error {
// 	defer wg.Done()
// 	for {
// 		buf := make([]byte, 512)
// 		n, err := (*p.upstream).Read(buf)
// 		if n == 0 {
// 			continue
// 		}
// 		log.Printf("Reading %d from upstream: %x", n, buf[:n])
// 		if err != nil {
// 			// log.Printf("Failed to read from upstream connection: %v", err)
// 			cancleFunc(err)
// 			buf = nil
// 			return err
// 		}
// 		data := buf[:n]
// 		// packet := pac.NewPacket(&pac.Raw{})
// 		// err = packet.Decode(&data, len(data))
// 		// if err != nil {
// 		// 	log.Printf("Failed to decode packet: %v", err)
// 		// 	cancleFunc(err)
// 		// 	return err
// 		// }
// 		// p.routerLock.Lock()
// 		output <- data
// 		// p.routerLock.Unlock()
// 		buf = nil
// 	}
// }

// func (p *Proxy) writeClient(ctx context.Context, cancleFunc context.CancelCauseFunc, input []byte) error {
// 	// log.Printf("Writing to client: %s", input.String())
// 	// data := input.Encode()
// 	// log.Printf("Writing hex to client: %x", input)
// 	_, err := (*p.client).Write(input)
// 	if err != nil {
// 		// log.Printf("Failed to write to client connection: %v", err)
// 		cancleFunc(err)
// 		return err
// 	}
// 	return nil
// }

// func (p *Proxy) writeUpstream(ctx context.Context, cancleFunc context.CancelCauseFunc, input []byte) error {
// 	// log.Printf("Writing to upstream: %s", input.String())
// 	// data := input.Encode()
// 	// log.Printf("Writing to upstream: %x", input)
// 	_, err := (*p.upstream).Write(input)
// 	if err != nil {
// 		// log.Printf("Failed to write to upstream connection: %v", err)
// 		cancleFunc(err)
// 		return err
// 	}
// 	return nil
// }

// func (p *Proxy) Close() error {
// 	if p.client != nil {
// 		(*p.client).Close()
// 		p.client = nil
// 	}
// 	if p.upstream != nil {
// 		(*p.upstream).Close()
// 		p.upstream = nil
// 	}
// 	return nil
// }

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

func (p *Proxy) Serve() {
	startTime := time.Now()
	defer func() {
		log.Printf("[Proxy] Cleanup session")
		runtime.GC()
		if r := recover(); r != nil {
			log.Printf("[Proxy] panic: ", r)
			return
		}
	}()

	statemachine := state.NewStateMachine(p.Listener, GetServerList(), p.logger)
	err := statemachine.Run() // Block until someone connected
	if err != nil {
		log.Printf("[Proxy] Connection accept failed: %v", err)
		// p.StateMachine.Destroy()
		return
	}

	log.Printf("[Proxy] Connection between proxy and client established")
	// p.threadWaitGroup.Add(1)
	go func(_sm *state.StateMachine) {
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
	}(statemachine)
}

func NewProxy(port string, logger *Logger.Logger) (Iproxy, error) {
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
	return &Proxy{Listener: &listener, threadWaitGroup: &sync.WaitGroup{}, logger: logger}, nil
}
