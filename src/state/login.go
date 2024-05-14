package state

import (
	"errors"
	"log"
	pac "mc_reverse_proxy/src/packet"
	"mc_reverse_proxy/src/utils"
)

type LoginState struct {
	sm      *StateMachine
	conn    *Connection
	Data    pac.Packet[*pac.PlayerData]
	OldData []byte
}

func (h *LoginState) ImplState() {}

func (h *LoginState) Enter(sm *StateMachine) error {
	log.Printf("[state machine: Login] Enter login state")
	h.sm = sm
	// log.Printf("%v", &h.sm.Conn.ClientData)
	// log.Printf("%v | %v", h.sm.Conn, h.sm.Conn.ClientConn)
	select {
	case rawData := <-h.sm.Conn.ClientData:
		// log.Printf("%x", rawData)
		pd := pac.PlayerData{}
		pd_pac := pac.Packet[*pac.PlayerData]{Data: &pd}
		err := pd_pac.Decode(&rawData, len(rawData))
		if err != nil {
			return err
		}
		h.Data = pd_pac
		// log.Printf("%s %s", h.hostname, data.String())
		return nil
	case <-h.sm.ctx.Done():
		return errors.New("Context Done")
	}
	// log.Printf("Data %x", rawData)
}

func (h *LoginState) Action() error {
	h.sm.StateChangeLock.Lock()
	defer h.sm.StateChangeLock.Unlock()
	data, err := h.Data.Encode()
	if err != nil {
		return err
	}
	log.Printf("[state machine: Login Debug] OldData: %v", h.OldData)
	log.Printf("[state machine: Login Debug] Player Data: %v", data)
	return h.sm.Conn.WriteServer(utils.Concat(h.OldData, data))
}

func (h *LoginState) Exit() IState {
	log.Printf("[state machine: Login] Change to pass through state")
	return &PassthroughState{}
	// return &RejectState{Message: "Unexpected condition"}
}
