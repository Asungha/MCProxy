package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	config "mc_reverse_proxy/src/configuration/service"
	metricDTO "mc_reverse_proxy/src/metric/dto"
	packetLoggerService "mc_reverse_proxy/src/packet-logger/service"
	utils "mc_reverse_proxy/src/utils"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

	configService *config.ConfigurationService

	StrictValidation bool
	SessionId        string
}

func (c *ConnectionService) WaitClientConnection() error {
	if c.Listener == nil {
		return errors.New("no listener available")
	}
	clientConn, err := (*c.Listener).Accept()
	if err != nil {
		utils.FLogDebug.Connection("Session %s, Failed to accept client connection: %v", c.SessionId, err)
		return err
	}
	utils.FLogDebug.Connection("Session %s, Initiated connection between client (%s) and the proxy", c.SessionId, clientConn.RemoteAddr().String())
	c.ClientConn = &clientConn
	c.ClientAddress = clientConn.RemoteAddr().String()
	return nil
}

func (c *ConnectionService) ConnectServer(host string) error {
	upstreamConn, err := net.Dial("tcp", host)
	if err != nil {
		utils.FLogDebug.Connection("Session %s, Failed to connect to upstream server %s: %v", c.SessionId, host, err)
		return err
	}
	utils.FLogDebug.Connection("Session %s, Initiated connection between the proxy and server (%s)", c.SessionId, (*c.ClientConn).RemoteAddr().String())
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
	defer utils.FLogDebug.Connection("Session %s, Client %s socket terminated", c.SessionId, (*c.ClientConn).RemoteAddr().String())
	errs := make(chan error)
	go func(errs chan error) {
		for {
			buf := make([]byte, 1024)
			n, err := (*c.ClientConn).Read(buf)
			if err != nil {
				// log.Printf("[client listener] Failed to read from client connection: %v", err)
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
			fragments, err := utils.SplitDataframe(data, c.StrictValidation)
			if err != nil {
				if len(fragments) != 0 {
					for _, f := range fragments {
						addr := (*c.ClientConn).RemoteAddr().String()
						addr = strings.Replace(addr, "[::1]", "localhost", 1)
						packetLoggerService.Send(packetLoggerService.PacketLog{
							Type:      f.Type,
							IP:        strings.Split(addr, ":")[0],
							Port:      strings.Split(addr, ":")[1],
							Timestamp: time.Now(),
							Data:      hex.EncodeToString(f.Data),
						})
					}
				}
				c.Cancle(err)
				buf = nil
				errs <- err
				return
			}
			for _, f := range fragments {
				c.ClientData <- f.Data
			}
			buf = nil
		}
	}(errs)
	for {
		select {
		case e := <-errs:
			(*c.ClientConn).Close()
			// log.Printf("[client listener] Exit due to error: %v", e)
			c.Cancle(e)
			<-c.Ctx.Done()
			return nil
		case <-c.Ctx.Done():
			c.Cancle(nil)
			(*c.ClientConn).Close()
			<-errs
			// log.Printf("[client listener] Exit due context canceled")
			return nil
		}
	}
}

func (c *ConnectionService) ListenServer() error {
	defer utils.FLogDebug.Connection("Session %s, Server %s socket terminated", c.SessionId, (*c.ServerConn).RemoteAddr().String())
	errs := make(chan error)
	go func(errs chan error) {
		for {
			buf := make([]byte, 12400)
			n, err := (*c.ServerConn).Read(buf)
			if err != nil {
				// log.Printf("[server listener] Failed to read from upstream connection: %v", err)
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
			// log.Printf("[server listener] Thread exit due to error: %v", e)
			<-c.Ctx.Done()
			return nil
		case <-c.Ctx.Done():
			c.Cancle(nil)
			(*c.ServerConn).Close()
			<-errs
			// log.Printf("[server listener] Exit due context canceled")
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

func NewConnectionService(configService *config.ConfigurationService, mutex *sync.Mutex, ctx context.Context, cancle context.CancelCauseFunc, listener *net.Listener) *ConnectionService {
	return &ConnectionService{
		configService:   configService,
		StateChangeLock: mutex,
		Ctx:             ctx,
		Cancle:          cancle,
		Listener:        listener,
		NetworkMetric:   metricDTO.NetworkMetric{},
		ClientData:      make(chan []byte, 16),
		ServerData:      make(chan []byte, 16),
		SessionId:       uuid.NewString(),
	}
}
