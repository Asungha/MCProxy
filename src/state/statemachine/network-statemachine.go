package statemachine

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mc_reverse_proxy/src/metric"
	pac "mc_reverse_proxy/src/packet"
	service "mc_reverse_proxy/src/service"
	networkstate "mc_reverse_proxy/src/state/networkstate"
	state "mc_reverse_proxy/src/state/state"
	"net"
	"sync"

	"github.com/google/uuid"
	// statemachine "mc_reverse_proxy/src/statemachine"
)

type NetworkStateMachine struct {
	errorMetric     metric.ErrorMetric
	Conn            *service.Connection
	hostname        string
	Data            *pac.Packet[*pac.Handshake]
	StateChangeLock sync.Mutex
	ClientConnected chan bool
	AStateMachine
	metric.Loggable
}

func (sm *NetworkStateMachine) UUID() string {
	return uuid.NewString()
}

func (sm *NetworkStateMachine) Log() metric.Log {
	return metric.Log{ErrorMetric: sm.errorMetric, NetworkMetric: sm.Conn.NetworkMetric}
}

func (sm *NetworkStateMachine) Serve() error {
	return sm.Run()
}

func NewNetworkStatemachine(listener *net.Listener, serverlist map[string]map[string]string) *NetworkStateMachine {
	ctx, cancle := context.WithCancelCause(context.Background())
	sm := &NetworkStateMachine{
		StateChangeLock: sync.Mutex{},
		errorMetric:     metric.ErrorMetric{},
		ClientConnected: make(chan bool),
	}
	sm.Conn = service.NewConnection(&sm.StateChangeLock, ctx, cancle, listener)
	sm.Conn.ClientData = make(chan []byte)
	sm.Conn.ServerData = make(chan []byte)
	sm.Conn.ServerList = serverlist
	sm.Conn.Ctx = ctx
	sm.Conn.Cancle = cancle
	hostname := &sm.hostname
	Data := sm.Data
	errorMetric := &sm.errorMetric
	StateChangeLock := &sm.StateChangeLock
	ClientConnected := sm.ClientConnected

	// handler
	var initHandler state.Function = func(i state.IState) error {
		err := sm.Conn.WaitClientConnection()
		if err != nil {
			log.Printf(err.Error())
			return err
		}
		ClientConnected <- true
		log.Printf("Creating log")
		go sm.Conn.ListenClient()
		return nil
	}

	var handshakeHandler state.Function = func(i state.IState) error {
		log.Printf("[state machine: Handshake] Enter handshake state")
		log.Printf("%v", sm.Conn)
		log.Printf("Chan %v", sm.Conn.ClientData)
		select {
		case rawData := <-sm.Conn.ClientData:
			log.Printf("raw %v", rawData)
			hs := pac.Handshake{}
			data := pac.Packet[*pac.Handshake]{Data: &hs}
			err := data.Decode(&rawData, len(rawData))
			if err != nil {
				return err
			}
			*hostname = data.Data.Hostname
			Data = &data
		case <-sm.Conn.Ctx.Done():
			log.Printf("Done")
			return errors.New("Context Done")
		}
		log.Printf("[state machine: Handshake] handshake data recieved")

		if data, ok := sm.Conn.ServerList[*hostname]; ok {
			if target, ok := data["target"]; ok {
				err := sm.Conn.ConnectServer(target)
				if err != nil {
					errorMetric.ServerConnectFailed += 1
					return err
				}
				err = sm.Conn.PreConditionCheck()
				if err != nil {
					log.Printf("[handshake state] Precondition failed %v", err)
					return err
				}
				go sm.Conn.ListenServer()
				hs_packet, err := Data.Encode()
				if err != nil {
					errorMetric.PacketDeserializeFailed += 1
					log.Printf("[handshake state] Encode handshale failed %v", err)
					return err
				}
				StateChangeLock.Lock()
				err = sm.Conn.WriteServer(hs_packet)
				if err != nil {
					log.Printf("[handshake state] Handshake packet send failed %v", err)
					return err
				}
				if Data.Data.Tail != nil {
					err = sm.Conn.WriteServer(Data.Data.Tail)
					if err != nil {
						log.Printf("[handshake state] additional packet send failed %v", err)
						return err
					}
				}
				StateChangeLock.Unlock()
				return nil
			}
			return errors.New("[Handshake State] host config file malformed")
		}
		errorMetric.HostnameResolveFailed += 1
		return errors.New(fmt.Sprintf("[Handshake State] Host %s not found", *hostname))
	}

	var passthroughHandler state.Function = func(i state.IState) error {
		select {
		case cData := <-sm.Conn.ClientData:
			StateChangeLock.Lock()
			sm.Conn.WriteServer(cData)
			StateChangeLock.Unlock()
			return nil
		case sData := <-sm.Conn.ServerData:
			StateChangeLock.Lock()
			sm.Conn.WriteClient(sData)
			StateChangeLock.Unlock()
			return nil
		case <-sm.Conn.Ctx.Done():
			return errors.New("context Done")
		}
	}

	var rejectHandler state.Function = func(i state.IState) error {
		// sm.Cancle(errors.New("Rejected"))
		sm.Conn.CloseConn()
		log.Printf("Closed connection")
		return errors.New("Rejected")
	}

	// state
	initState := networkstate.InitState{}
	initState.Init(&initHandler)

	handshakeState := networkstate.HandshakeState{}
	handshakeState.Init(&handshakeHandler)
	// handshakeState.SetTimeout(1 * time.Second)

	passthroughState := networkstate.PassthroughState{}
	passthroughState.Init(&passthroughHandler)

	rejectState := networkstate.RejectState{}
	rejectState.Init(&rejectHandler)

	sm.RegisterState("init", &initState)
	sm.RegisterState("handshake", &handshakeState)
	sm.RegisterState("passthough", &passthroughState)
	sm.RegisterState("reject", &rejectState)

	var initTransistion func() bool = func() bool {
		log.Printf("init >>> hs")
		return true
	}
	sm.TransistionCondition(TransistionPair{Source: "init", Destination: "handshake"}, &initTransistion)

	var handshakeTransistion func() bool = func() bool {
		return true
	}
	sm.TransistionCondition(TransistionPair{Source: "handshake", Destination: "passthough"}, &handshakeTransistion)

	var HaltTransistion func() bool = func() bool {
		log.Println(sm.Conn.Ctx.Err())
		return sm.Conn.Ctx.Err() != nil
	}
	sm.TransistionCondition(TransistionPair{Source: "passthough", Destination: "reject"}, &HaltTransistion)

	var passthoughTransistion func() bool = func() bool {
		return true
	}
	sm.TransistionCondition(TransistionPair{Source: "passthough", Destination: "passthough"}, &passthoughTransistion)

	sm.Construct()
	sm.SetRoot("init")
	go sm.Run()
	return sm
}
