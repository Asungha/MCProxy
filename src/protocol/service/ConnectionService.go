package network

import (
	"context"
	"errors"
	"fmt"
	"log"
	metricDTO "mc_reverse_proxy/src/metric/dto"
	utils "mc_reverse_proxy/src/utils"
	"net"
	"sync"
)

type ConnectionService struct {
	NetworkMetric metricDTO.NetworkMetric

	TargetHostname  string
	ClientAddress   string
	ClientConn      *net.Conn
	ServerConn      *net.Conn
	ClientData      chan []byte
	ServerData      chan []byte
	Ctx             context.Context
	Cancle          context.CancelCauseFunc
	StateChangeLock *sync.Mutex
	Listener        *net.Listener
}

func (c *ConnectionService) WaitClientConnection() error {
	if c.Listener == nil {
		return errors.New("no listener available")
	}
	clientConn, err := (*c.Listener).Accept()
	if err != nil {
		log.Printf("Failed to accept client connection: %v", err)
		return err
	}
	log.Printf("[Proxy] Initiated connection between %s and %s", clientConn.RemoteAddr().String(), "localhost:25565")
	// clientConn.SetDeadline(time.Now().Add(5 * time.Second))
	c.ClientConn = &clientConn
	c.ClientAddress = clientConn.RemoteAddr().String()
	return nil
}

func (c *ConnectionService) ConnectServer(host string) error {
	upstreamConn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("Failed to connect to upstream server: %v", err)
		return err
	}
	log.Printf("[Proxy] Initiated connection between %s and %s", "proxy", upstreamConn.RemoteAddr().String())
	// upstreamConn.SetDeadline(time.Now().Add(5 * time.Second))
	c.ServerConn = &upstreamConn
	return nil
}

func (c *ConnectionService) PreConditionCheck() error {
	if c.ClientConn == nil {
		return errors.New("client connection not established")
	}
	if c.ServerConn == nil {
		return errors.New("upstream connection not established")
	}
	return nil
}

func (c *ConnectionService) ListenClient() error {
	defer log.Println("[client listener controller] Thread exit")
	// defer c.WaitGroup.Done()
	errs := make(chan error)
	// done := make(chan bool)
	go func(errs chan error) {
		defer log.Println("[client listener] Thread exit")
		for {
			buf := make([]byte, 1024)
			// (*c.ClientConn).SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := (*c.ClientConn).Read(buf)
			// log.Printf("[client listener Debug] Reading %d bytes from client: %x", n, buf[:n])
			if err != nil {
				log.Printf("[client listener] Failed to read from client connection: %v", err)
				c.Cancle(err)
				buf = nil
				errs <- err
				return
			}
			if n == 0 {
				continue
			}
			c.NetworkMetric.ClientPacketRx += 1
			c.NetworkMetric.ClientDataRx += uint(n)
			data := buf[:n]
			// c.StateChangeLock.Lock()
			// log.Printf("Sending data")
			fragments, err := utils.SplitDataframe(data)
			if err != nil {
				log.Printf(err.Error())
				c.Cancle(err)
				buf = nil
				errs <- err
				return
			}
			for _, f := range fragments {
				c.ClientData <- f
			}
			// c.ClientData <- data
			// log.Printf("Sending datato %v done", c.ClientData)
			// c.StateChangeLock.Unlock()
			buf = nil
		}
	}(errs)
	for {
		select {
		case e := <-errs:
			(*c.ClientConn).Close()
			log.Printf("[client listener] Exit due to error: %v", e)
			c.Cancle(e)
			<-c.Ctx.Done()
			return nil
		case <-c.Ctx.Done():
			c.Cancle(nil)
			(*c.ClientConn).Close()
			<-errs
			log.Printf("[client listener] Exit due context canceled")
			return nil
		}
	}
}

func (c *ConnectionService) ListenServer() error {
	defer log.Println("[server listener controller] Thread exit")
	errs := make(chan error)
	// datas := make(chan []byte)
	go func(errs chan error) {
		defer log.Println("[server listener] Thread exit")
		for {
			buf := make([]byte, 12400)
			// (*c.ServerConn).SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := (*c.ServerConn).Read(buf)
			// log.Printf("[server listener Debug] Reading %d bytes from upstream: %x", n, buf[:n])
			if err != nil {
				log.Printf("[server listener] Failed to read from upstream connection: %v", err)
				c.Cancle(err)
				buf = nil
				errs <- err
				return
			}
			if n == 0 {
				continue
			}
			c.NetworkMetric.ServerPacketRx += 1
			c.NetworkMetric.ServerDataRx += uint(n)
			data := buf[:n]
			c.ServerData <- data
			// c.StateChangeLock.Unlock()
			buf = nil
		}
	}(errs)
	for {
		select {
		case e := <-errs:
			c.Cancle(e)
			(*c.ServerConn).Close()
			log.Printf("[server listener] Thread exit due to error: %v", e)
			<-c.Ctx.Done()
			return nil
		case <-c.Ctx.Done():
			c.Cancle(nil)
			(*c.ServerConn).Close()
			<-errs
			log.Printf("[server listener] Exit due context canceled")
			return nil
		}
	}
}

func (c *ConnectionService) WriteClient(input []byte) error {
	// log.Printf("Writing to client: %s", input.String())
	// data := input.Encode()
	// log.Printf("[client writter Debug] Writing data to client: %x", input)
	n, err := (*c.ClientConn).Write(input)
	if err != nil {
		// log.Printf("[] Failed to write to client connection: %v", err)
		c.Cancle(err)
		return err
	}
	c.NetworkMetric.ClientPacketTx += 1
	c.NetworkMetric.ClientDataTx += uint(n)
	return nil
}

func (c *ConnectionService) WriteServer(input []byte) error {
	// log.Printf("Writing to upstream: %s", input.String())
	// data := input.Encode()
	// log.Printf("[server writter Debug]  Writing to server: %x", input)
	n, err := (*c.ServerConn).Write(input)
	if err != nil {
		// log.Printf("Failed to write to upstream connection: %v", err)
		c.Cancle(err)
		return err
	}
	c.NetworkMetric.ServerPacketTx += 1
	c.NetworkMetric.ServerDataTx += uint(n)
	return nil
}

func (c *ConnectionService) CloseClientConn() error {
	if c.ClientConn != nil {
		(*c.ClientConn).Close()
	}
	return nil
}

func (c *ConnectionService) CloseServerConn() error {
	if c.ServerConn != nil {
		(*c.ServerConn).Close()
	}
	return nil
}

func (c *ConnectionService) CloseConn() error {
	// c.WaitGroup.Wait()
	clientErr := c.CloseClientConn()
	serverErr := c.CloseServerConn()
	if clientErr != nil || serverErr != nil {
		return errors.New(fmt.Sprintf("error closing connection: Server %v, Client %v", serverErr, clientErr))
	}
	return nil
}

func (c *ConnectionService) Destroy() {
	c.CloseConn()
}

func NewConnectionService(mutex *sync.Mutex, ctx context.Context, cancle context.CancelCauseFunc, listener *net.Listener) *ConnectionService {
	return &ConnectionService{
		StateChangeLock: mutex,
		Ctx:             ctx,
		Cancle:          cancle,
		Listener:        listener,
		NetworkMetric:   metricDTO.NetworkMetric{},
		ClientData:      make(chan []byte, 16),
		ServerData:      make(chan []byte, 16),
	}
}
