package state

import (
	"log"
	Logger "mc_reverse_proxy/src/logger"
)

type InitState struct {
	sm         *StateMachine
	ServerList map[string]map[string]string
}

func (init *InitState) ImplState() {}

func (init *InitState) Enter(sm *StateMachine) error {
	// log.Printf("Enter init state")
	init.sm = sm
	// log.Printf("Wait for connection")
	err := init.sm.Conn.WaitClientConnection()
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	log.Printf("Creating log")
	log := Logger.NewLog(init.sm.Conn.ClientAddress, Logger.CONNECT, "Init", nil, nil, nil)
	init.sm.PushLog(log)
	// log.Printf("got connection")
	// init.sm.Conn.WaitGroup.Add(1)
	go init.sm.Conn.ListenClient()
	// log.Printf("%v | %v", init.sm.Conn, init.sm.Conn.ClientConn)
	return nil
}
func (init *InitState) Exit() IState {
	// log.Printf("Exit init state")
	return &HandshakeState{}
}

func (init *InitState) Action() error {
	return nil
}
