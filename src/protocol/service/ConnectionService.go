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
		log.Printf("[Proxy] Failed to accept client connection: %v", err)
		return err
	}
	log.Printf("[Proxy] Initiated connection between %s and %s", clientConn.RemoteAddr().String(), "localhost:25565")
	c.ClientConn = &clientConn
	c.ClientAddress = clientConn.RemoteAddr().String()
	return nil
}

func (c *ConnectionService) ConnectServer(host string) error {
	upstreamConn, err := net.Dial("tcp", host)
	if err != nil {
		log.Printf("[Proxy] Failed to connect to upstream server: %v", err)
		return err
	}
	log.Printf("[Proxy] Initiated connection between %s and %s", "proxy", upstreamConn.RemoteAddr().String())
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
	errs := make(chan error)
	go func(errs chan error) {
		for {
			buf := make([]byte, 1024)
			n, err := (*c.ClientConn).Read(buf)
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
			fragments, err := utils.SplitDataframe(data)
			if err != nil {
				c.Cancle(err)
				buf = nil
				errs <- err
				return
			}
			for _, f := range fragments {
				c.ClientData <- f
			}
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
	errs := make(chan error)
	go func(errs chan error) {
		for {
			buf := make([]byte, 12400)
			n, err := (*c.ServerConn).Read(buf)
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
	n, err := (*c.ClientConn).Write(input)
	if err != nil {
		c.Cancle(err)
		return err
	}
	c.NetworkMetric.ClientPacketTx += 1
	c.NetworkMetric.ClientDataTx += uint(n)
	return nil
}

func (c *ConnectionService) WriteServer(input []byte) error {
	n, err := (*c.ServerConn).Write(input)
	if err != nil {
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
