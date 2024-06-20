package statemachine

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	metricDTO "mc_reverse_proxy/src/metric/dto"
	metricService "mc_reverse_proxy/src/metric/service"
	pac "mc_reverse_proxy/src/protocol/packet"
	networkService "mc_reverse_proxy/src/protocol/service"

	// statemachine "mc_reverse_proxy/src/proxy/service/statemachine"
	state "mc_reverse_proxy/src/statemachine/dto"

	proxyService "mc_reverse_proxy/src/proxy/service"

	// state "mc_reverse_proxy/src/state/state"
	"mc_reverse_proxy/src/utils"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	// statemachine "mc_reverse_proxy/src/statemachine"
)

type NetworkStateMachine struct {
	playerMetric metricDTO.PlayerMetric
	Conn         *networkService.ConnectionService
	serverRepo   proxyService.ServerRepositoryService
	Data         *pac.Handshake

	StateMachine

	hostname        string
	StateChangeLock sync.Mutex
	ClientConnected chan bool
	// metric.Loggable
}

// func (sm *NetworkStateMachine) Run() error {
// 	return sm.Run()
// }

func (sm *NetworkStateMachine) UUID() string {
	return uuid.NewString()
}

// func (sm *NetworkStateMachine) Log() metric.Log {
// 	return metric.Log{ErrorMetric: sm.playerMetric.ErrorMetric, NetworkMetric: sm.Conn.NetworkMetric, PlayerMetric: sm.playerMetric, ProxyMetric: sm.proxyMetric}
// }

func (sm *NetworkStateMachine) Serve() error {
	return sm.Run()
}

func NewNetworkStatemachine(listener *net.Listener, serverRepo proxyService.ServerRepositoryService, serverMetric *metricDTO.ProxyMetric, metricCollector *metricService.MetricService) *NetworkStateMachine {
	ctx, cancle := context.WithCancelCause(context.Background())
	sm := &NetworkStateMachine{
		StateChangeLock: sync.Mutex{},
		playerMetric:    metricDTO.PlayerMetric{ErrorMetric: metricDTO.ErrorMetric{}},
		// proxyMetric:     metric.ProxyMetric{},
		ClientConnected: make(chan bool, 1),
	}
	sm.Conn = networkService.NewConnectionService(&sm.StateChangeLock, ctx, cancle, listener)
	// sm.playerMetric.NetworkMetric = &sm.Conn.NetworkMetric
	sm.serverRepo = serverRepo
	sm.Ctx = ctx
	sm.Cancle = cancle

	logPusher := &metricService.LogPusher{Collector: *metricCollector}
	hostname := &sm.hostname
	Data := sm.Data
	playerMetric := &sm.playerMetric
	StateChangeLock := &sm.StateChangeLock
	ClientConnected := sm.ClientConnected
	LoggedIn := false
	isLoggedIn := &LoggedIn
	isHttp := false
	// loginpayload := &pac.Login{}
	var loginpayload *pac.Login
	// playername := &sm.playername

	sm.DeferFunc = func() {
		sm.Conn.CloseConn()
		logPusher.PushNetworkMetric(sm.Conn.NetworkMetric)
	}
	// handler
	var initHandler state.Function = func(i state.IState) error {
		err := sm.Conn.WaitClientConnection()
		if err != nil {
			playerMetric.AcceptFailed += 1
			logPusher.PushErrorMetric(metricDTO.ErrorMetric{AcceptFailed: 1})
			return err
		}
		logPusher.PushProxyMetric(metricDTO.ProxyMetric{Connected: 1})
		s := strings.Split(sm.Conn.ClientAddress, ":")
		sm.playerMetric.IP = s[0]
		sm.playerMetric.Port = s[1]
		ClientConnected <- true
		go sm.Conn.ListenClient()
		return nil
	}

	var handshakeHandler state.Function = func(i state.IState) error {
		select {
		case rawData := <-sm.Conn.ClientData:
			hs := pac.Handshake{}
			if utils.IsHTTPMethod(strings.Split(string(rawData), " ")[0]) {
				isHttp = true
				playerMetric.PacketDeserializeFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{PacketDeserializeFailed: 1})
				return nil
			}
			err := hs.Decode(rawData)
			if err != nil {
				playerMetric.PacketDeserializeFailed += 1
				if bytes.Equal(rawData[:2], []byte{0xfe, 0x01}) {
					logPusher.PushProxyMetric(metricDTO.ProxyMetric{PlayerGetStatus: 1})
				}
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{PacketDeserializeFailed: 1})
				return err
			}
			*hostname = hs.GetHostname()
			Data = &hs
		case <-sm.Conn.Ctx.Done():
			return errors.New("Context Done")
		}
		if len(Data.Tail) > 2 {
			l := pac.Login{}
			err := l.Decode(Data.Tail)
			if err != nil {
				playerMetric.PacketDeserializeFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{PacketDeserializeFailed: 1})
				return err
			}
			loginpayload = &l
		}
		if Data.NextState == 0x01 {
			logPusher.PushProxyMetric(metricDTO.ProxyMetric{PlayerGetStatus: 1})
		} else if Data.NextState == 0x02 {
			logPusher.PushProxyMetric(metricDTO.ProxyMetric{PlayerLogin: 1})
		}
		if target, err := sm.serverRepo.Resolve(*hostname); err == nil {
			err := sm.Conn.ConnectServer(target)
			if err != nil {
				playerMetric.ServerConnectFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{ServerConnectFailed: 1})
				return err
			}
			err = sm.Conn.PreConditionCheck()
			if err != nil {
				log.Printf("[handshake state] Precondition failed %v", err)
				playerMetric.ServerConnectFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{ServerConnectFailed: 1})
				return err
			}
			go sm.Conn.ListenServer()
			hs_packet, err := Data.Encode()
			if err != nil {
				playerMetric.PacketDeserializeFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{PacketDeserializeFailed: 1})
				log.Printf("[handshake state] Encode handshale failed %v", err)
				return err
			}
			StateChangeLock.Lock()
			err = sm.Conn.WriteServer(hs_packet)
			if err != nil {
				log.Printf("[handshake state] Handshake packet send failed %v", err)
				playerMetric.ServerConnectFailed += 1
				logPusher.PushErrorMetric(metricDTO.ErrorMetric{ServerConnectFailed: 1})
				return err
			}
			if Data.Tail != nil {
				err = sm.Conn.WriteServer(Data.Tail)
				if err != nil {
					log.Printf("[handshake state] additional packet send failed %v", err)
					playerMetric.ServerConnectFailed += 1
					logPusher.PushErrorMetric(metricDTO.ErrorMetric{ServerConnectFailed: 1})
					return err
				}
			}
			StateChangeLock.Unlock()
			return nil
		} else {
			playerMetric.HostnameResolveFailed += 1
			logPusher.PushErrorMetric(metricDTO.ErrorMetric{HostnameResolveFailed: 1})
			return errors.New(fmt.Sprintf("[Handshake State] Can't resolve host %s. Error %v", *hostname, err))
		}
	}

	var statusReqHandler state.Function = func(i state.IState) error {
		return nil
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
			if *isLoggedIn {
				logPusher.PushProxyMetric(metricDTO.ProxyMetric{PlayerPlaying: -1})
			}
			return errors.New("context Done")
		}
	}

	var rejectHandler state.Function = func(i state.IState) error {
		StateChangeLock.Lock()
		defer StateChangeLock.Unlock()
		if *isLoggedIn {
			logPusher.PushProxyMetric(metricDTO.ProxyMetric{PlayerPlaying: -1})
		}
		sm.Conn.CloseConn()
		return nil
	}

	var pingHandler state.Function = func(i state.IState) error {
		return nil
	}

	initState := state.NewState(initHandler)
	handshakeState := state.NewState(handshakeHandler)
	statusReqState := state.NewState(statusReqHandler)
	loginState := state.NewState(func(i state.IState) error {
		if loginpayload != nil {
			res, err := loginpayload.Encode()
			if err != nil {
				return err
			}
			go func() {
				sm.Conn.ClientData <- res
			}()
		}
		select {
		case cData := <-sm.Conn.ClientData:
			StateChangeLock.Lock()
			p := pac.Login{}
			err := p.Decode(cData)
			if err != nil {
				return err
			}
			log.Printf("[Loging state] Player %s logged in", p.Name)
			sm.playerMetric.PlayerName = p.Name
			sm.playerMetric.LogginTime = time.Now()
			sm.Conn.WriteServer(cData)
			StateChangeLock.Unlock()
			return nil
		case <-time.After(15 * time.Second):
			return errors.New("loging timeout")
		case <-sm.Conn.Ctx.Done():
			return errors.New("context Done")
		}
	})
	passthroughState := state.NewState(passthroughHandler)
	rejectState := state.NewState(rejectHandler)
	pingState := state.NewState(pingHandler)

	sm.RegisterState("init", initState)
	sm.RegisterState("handshake", handshakeState)
	sm.RegisterState("status", statusReqState)
	sm.RegisterState("ping", pingState)
	sm.RegisterState("login", loginState)
	sm.RegisterState("passthough", passthroughState)
	sm.RegisterState("reject", rejectState)

	sm.TransistionCondition(TransistionPair{Source: "init", Destination: "handshake"}, state.True)

	sm.TransistionFunction("handshake", func(i state.IState) (state.IState, error) {
		if isHttp || *hostname == "" || sm.Conn.PreConditionCheck() != nil {
			return nil, nil
		}
		if Data.NextState == 0x01 && Data.ID == 0 {
			// go collector.Collect()
			return statusReqState, nil
		}
		return loginState, nil
	})

	sm.TransistionCondition(TransistionPair{Source: "status", Destination: "ping"}, state.True)
	sm.TransistionCondition(TransistionPair{Source: "ping", Destination: "passthough"}, state.True)

	// sm.TransistionCondition(TransistionPair{Source: "ping", Destination: ""}, state.True)
	sm.TransistionCondition(TransistionPair{Source: "login", Destination: "passthough"}, state.True)
	sm.TransistionFunction("passthough", func(i state.IState) (state.IState, error) {
		if sm.Conn.Ctx.Err() != nil {
			return rejectState, nil
		} else {
			return passthroughState, nil
		}
	})

	// sm.TransistionCondition(TransistionPair{Source: "handshake", Destination: "reject"}, RejectHandshake)
	// sm.TransistionCondition(TransistionPair{Source: "handshake", Destination: ""}, func() bool {
	// 	return Data.Data.Type == 1
	// })
	// sm.TransistionCondition(TransistionPair{Source: "handshake", Destination: "passthough"}, state.True)

	// sm.TransistionCondition(TransistionPair{Source: "passthough", Destination: "reject"}, HaltTransistion)
	// sm.TransistionCondition(TransistionPair{Source: "passthough", Destination: "passthough"}, state.True)

	// sm.TransistionCondition(TransistionPair{Source: "reject", Destination: ""}, state.True)

	sm.Construct()
	sm.SetRoot("init")
	// go sm.Run()
	return sm
}
