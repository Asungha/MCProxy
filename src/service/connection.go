package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
)

type Connection struct {
	TargetHostname string
	ClientAddress  string
	ClientConn     *net.Conn
	ServerConn     *net.Conn
	ClientData     chan []byte
	ServerData     chan []byte
	Ctx            context.Context
	Cancle         context.CancelCauseFunc

	ServerList      map[string]map[string]string
	StateChangeLock *sync.Mutex
	Listener        *net.Listener

	// ListenAddress string
}

func (c *Connection) WaitClientConnection() error {
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

func (c *Connection) ConnectServer(host string) error {
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

func (c *Connection) PreConditionCheck() error {
	if c.ClientConn == nil {
		return errors.New("client connection not established")
	}
	if c.ServerConn == nil {
		return errors.New("upstream connection not established")
	}
	return nil
}

func (c *Connection) ListenClient() error {
	// defer c.WaitGroup.Done()
	errs := make(chan error)
	// done := make(chan bool)
	go func(errs chan error) {
		defer close(c.ClientData)
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
				break
			}
			if n == 0 {
				continue
			}
			data := buf[:n]
			// c.StateChangeLock.Lock()
			// log.Printf("Sending data")
			c.ClientData <- data
			// log.Printf("Sending data done")
			// c.StateChangeLock.Unlock()
			buf = nil
		}
		log.Printf("[client listener] Stop listening client")
	}(errs)
	for {
		select {
		case e := <-errs:
			(*c.ClientConn).Close()
			log.Printf("[ListenClient] exit due to error: %v", e)
			c.Cancle(e)
			return nil
			// case <-c.ctx.Done():
			// 	(*c.ClientConn).Close()
			// 	log.Printf("[client listener] Exit due context canceled")
			// 	return nil
		}
	}
}

func (c *Connection) ListenServer() error {
	// defer close(c.ServerData)
	errs := make(chan error)
	// datas := make(chan []byte)
	go func(errs chan error) {
		for {
			buf := make([]byte, 1024)
			// (*c.ServerConn).SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := (*c.ServerConn).Read(buf)
			// log.Printf("[server listener Debug] Reading %d bytes from upstream: %x", n, buf[:n])
			if err != nil {
				log.Printf("[server listener] Failed to read from upstream connection: %v", err)
				c.Cancle(err)
				buf = nil
				errs <- err
				break
			}
			if n == 0 {
				continue
			}
			data := buf[:n]
			// c.StateChangeLock.Lock()
			c.ServerData <- data
			// c.StateChangeLock.Unlock()
			buf = nil
		}
		log.Printf("[server listener] Stop listening server")
	}(errs)
	for {
		select {
		case e := <-errs:
			c.Cancle(e)
			(*c.ServerConn).Close()
			log.Printf("[ListenServer] exit due to error: %v", e)
			return nil
			// case <-c.ctx.Done():
			// 	(*c.ServerConn).Close()
			// 	log.Printf("[server listener] Exit due context canceled")
			// 	return nil
		}
	}
}

func (c *Connection) WriteClient(input []byte) error {
	// log.Printf("Writing to client: %s", input.String())
	// data := input.Encode()
	log.Printf("[client writter Debug] Writing data to client: %x", input)
	_, err := (*c.ClientConn).Write(input)
	if err != nil {
		// log.Printf("[] Failed to write to client connection: %v", err)
		c.Cancle(err)
		return err
	}
	return nil
}

func (c *Connection) WriteServer(input []byte) error {
	// log.Printf("Writing to upstream: %s", input.String())
	// data := input.Encode()
	log.Printf("[server writter Debug]  Writing to server: %x", input)
	_, err := (*c.ServerConn).Write(input)
	if err != nil {
		// log.Printf("Failed to write to upstream connection: %v", err)
		c.Cancle(err)
		return err
	}
	return nil
}

func (c *Connection) CloseClientConn() error {
	if c.ClientConn != nil {
		(*c.ClientConn).Close()
	}
	return nil
}

func (c *Connection) CloseServerConn() error {
	if c.ServerConn != nil {
		(*c.ServerConn).Close()
	}
	return nil
}

func (c *Connection) CloseConn() error {
	// c.WaitGroup.Wait()
	clientErr := c.CloseClientConn()
	serverErr := c.CloseServerConn()
	if clientErr != nil || serverErr != nil {
		return errors.New(fmt.Sprintf("error closing connection: Server %v, Client %v", serverErr, clientErr))
	}
	return nil
}

func (c *Connection) Destroy() {
	c.CloseConn()
	runtime.GC()
}

func NewConnection(mutex *sync.Mutex, ctx context.Context, cancle context.CancelCauseFunc, listener *net.Listener) Connection {
	return Connection{StateChangeLock: mutex, Ctx: ctx, Cancle: cancle, Listener: listener}
}
