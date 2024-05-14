package state

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
)

const (
	ACTION_UNDEFINED = iota
	ACTION_TRANSPARENT
	ACTION_REJECT
	ACTION_ACCEPT
	ACTION_CANCLE
)

type Data struct {
	bytes  *[]byte
	length int
}
type IState interface {
	ImplState()

	Enter(*StateMachine) error
	Action() error
	Exit() IState
}

type Connection struct {
	TargetHostname string
	ClientAddress  string
	ClientConn     *net.Conn
	ServerConn     *net.Conn
	ClientData     chan []byte
	ServerData     chan []byte
	ctx            context.Context
	cancle         context.CancelCauseFunc

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
		return errors.New("Client connection not established")
	}
	if c.ServerConn == nil {
		return errors.New("Upstream connection not established")
	}
	return nil
}

func (c *Connection) ListenClient() error {
	// defer c.WaitGroup.Done()
	errs := make(chan error)
	// done := make(chan bool)
	go func() {
		defer close(c.ClientData)
		for {
			buf := make([]byte, 1024)
			// (*c.ClientConn).SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := (*c.ClientConn).Read(buf)
			log.Printf("[client listener Debug] Reading %d bytes from client: %x", n, buf[:n])
			if err != nil {
				log.Printf("[client listener] Failed to read from client connection: %v", err)
				c.cancle(err)
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
	}()
	for {
		select {
		case e := <-errs:
			(*c.ClientConn).Close()
			log.Printf("[ListenClient] exit due to error: %v", e)
			c.cancle(e)
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
	go func() {
		for {
			buf := make([]byte, 1024)
			// (*c.ServerConn).SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := (*c.ServerConn).Read(buf)
			log.Printf("[server listener Debug] Reading %d bytes from upstream: %x", n, buf[:n])
			if err != nil {
				log.Printf("[server listener] Failed to read from upstream connection: %v", err)
				c.cancle(err)
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
	}()
	for {
		select {
		case e := <-errs:
			c.cancle(e)
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
		c.cancle(err)
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
		c.cancle(err)
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
		return errors.New(fmt.Sprintf("Error closing connection: Server %v, Client %v", serverErr, clientErr))
	}
	return nil
}

func (c *Connection) Destroy() {
	c.CloseConn()
	runtime.GC()
}

func NewConnection(mutex *sync.Mutex, ctx context.Context, cancle context.CancelCauseFunc, listener *net.Listener) Connection {
	return Connection{StateChangeLock: mutex, ctx: ctx, cancle: cancle, Listener: listener}
}

const (
	STATUS_UNKNOWN = iota
	STATUS_OK
	STATUS_EXIT
	STATUS_CANCLE
	STATUS_ERROR
)

type StateMachine struct {
	Conn            Connection
	currentState    IState
	states          map[string]IState
	StateChangeLock *sync.Mutex
	ctx             context.Context
	cancle          context.CancelCauseFunc
}

func (sm *StateMachine) setState(s IState) {
	sm.currentState = s
}

func (sm *StateMachine) Run() error {
	// sm.StateChangeLock.Lock()
	err := sm.currentState.Enter(sm)
	if err != nil {
		return err
	}
	err = sm.currentState.Action()
	if err != nil {
		return err
	}
	return nil
}

func (sm *StateMachine) Transition() int {
	sm.StateChangeLock.Lock()
	sm.currentState = sm.currentState.Exit()
	if sm.currentState == nil {
		// log.Printf("Exit")
		err := sm.Conn.CloseConn()
		if err != nil {
			log.Printf("[state machine] Error when transistioning state %v", err)
		}
		return STATUS_EXIT
	}
	sm.StateChangeLock.Unlock()
	err := sm.currentState.Enter(sm)
	if err != nil {
		return STATUS_ERROR
	}
	// log.Printf("Action")
	err = sm.currentState.Action()
	if err != nil {
		log.Printf("[state machine] Action failed: %v", err)
		if err.Error() == "Context Done" {
			return STATUS_EXIT
		}
		return STATUS_ERROR
	}
	// log.Printf("Action Done")
	return STATUS_OK
}

func (sm *StateMachine) Destroy() {
	sm.cancle(errors.New("force closed"))
	sm.Conn.cancle(errors.New("force closed"))
	sm.Conn.CloseConn()
	sm.Conn.Destroy()
	runtime.GC()
}

// func (sm *StateMachine) Action() error {
// 	err := sm.currentState.Action()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

type Event struct {
	Type string
	Data map[string]string
}

func NewStateMachine(listener *net.Listener, serverList map[string]map[string]string) *StateMachine {
	ctx, cancle := context.WithCancelCause(context.Background())
	sm := &StateMachine{
		currentState:    &InitState{},
		states:          make(map[string]IState),
		StateChangeLock: &sync.Mutex{},
		ctx:             ctx,
		cancle:          cancle,
	}
	sm.Conn = NewConnection(sm.StateChangeLock, sm.ctx, sm.cancle, listener)
	sm.Conn.ServerList = serverList
	sm.Conn.ClientData = make(chan []byte)
	sm.Conn.ServerData = make(chan []byte)
	return sm
}
