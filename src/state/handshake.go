package state

import (
	"errors"
	"fmt"
	"log"
	pac "mc_reverse_proxy/src/packet"
)

type HandshakeState struct {
	sm       *StateMachine
	conn     *Connection
	Data     pac.Packet[*pac.Handshake]
	hostname string
}

func (h *HandshakeState) ImplState() {}

func (h *HandshakeState) Enter(sm *StateMachine) error {
	// log.Printf("Enter handshake state")
	h.sm = sm
	// log.Printf("%v", &h.sm.Conn.ClientData)
	// log.Printf("%v | %v", h.sm.Conn, h.sm.Conn.ClientConn)
	select {
	case rawData := <-h.sm.Conn.ClientData:
		// log.Printf("%x", rawData)
		pd := pac.PlayerData{}
		pd_pac := pac.Packet[*pac.PlayerData]{Data: &pd}
		hs := pac.Handshake{PlayerData: &pd_pac}
		data := pac.Packet[*pac.Handshake]{Data: &hs}
		err := data.Decode(&rawData, len(rawData))
		if err != nil {
			return err
		}
		h.hostname = data.Data.Hostname
		h.Data = data
		// log.Printf("%s %s", h.hostname, data.String())
		return nil
	case <-h.sm.ctx.Done():
		return errors.New("Context Done")
	}
	// log.Printf("Data %x", rawData)
}

func (h *HandshakeState) Action() error {
	if data, ok := h.sm.Conn.ServerList[h.hostname]; ok {
		if target, ok := data["target"]; ok {
			err := h.sm.Conn.ConnectServer(target)
			if err != nil {
				return err
			}
			err = h.sm.Conn.PreConditionCheck()
			if err != nil {
				log.Printf("[handshake state] Precondition failed %v", err)
				return err
			}
			h.sm.Conn.WaitGroup.Add(1)
			go h.sm.Conn.ListenServer()
			hs_packet, err := h.Data.Encode()
			if err != nil {
				log.Printf("[handshake state] Encode handshale failed %v", err)
				return err
			}
			err = h.sm.Conn.WriteServer(hs_packet)
			if err != nil {
				log.Printf("[handshake state] Handshake packet send failed %v", err)
				return err
			}
			return nil
		}
		return errors.New("[Handshake State] host config file malformed")
	}
	return errors.New(fmt.Sprintf("[Handshake State] Host %s not found", h.hostname))
}

func (h *HandshakeState) Exit() IState {
	// log.Printf("Exit handshake state")
	err := h.sm.Conn.PreConditionCheck()
	if err != nil {
		return &RejectState{Message: err.Error()}
	}
	defer func() {
		if h.Data.Data.PlayerData != nil {
			h.Data.Data.PlayerData.Destroy()
		}
		h.Data.Destroy()
	}()
	if h.Data.Data.NextState == 0x01 {
		if h.hostname != "" {
			// log.Printf("Send status")
			// d := pac.Packet[*pac.Status]{}
			return &PassthroughState{}
		}
		h.sm.Conn.TargetHostname = h.Data.Data.Hostname
		// log.Printf("HandshakeState hostname %v", h.Data.Data.Hostname)
		return h
	} else {
		// log.Printf("Login")
		return &PassthroughState{}
	}
	// return &RejectState{Message: "Unexpected condition"}
}

// func (h *HandshakeState) Update(sm *StateMachine, data *map[string]string, event chan Event, ack chan bool) error {
// 	log.Printf("HandshakeState Update %v", h.Data.Data.NextState)
// 	if h.Data.Data.NextState == 0x01 {
// 		if h.hostname != "" {
// 			log.Printf("Send status")
// 			// d := pac.Packet[*pac.Status]{}
// 			return sm.setState(&PassthroughState{})
// 		}
// 		// log.Printf("Request Status")
// 		if data != nil {
// 			(*data)["playername"] = "Unkonwn"
// 			(*data)["hostname"] = h.Data.Data.Hostname
// 		}
// 		log.Printf("HandshakeState hostname %v", h.Data.Data.Hostname)
// 		rawData, err := h.Data.Encode()
// 		if err != nil {
// 			return err
// 		}
// 		event <- Event{Type: "handshake_hostname", Data: map[string]string{"hostname": h.Data.Data.Hostname, "data": string(rawData)}}
// 		<-ack
// 		h.hostname = h.Data.Data.Hostname
// 		return sm.setState(h)
// 	} else {
// 		// log.Printf("Passthrough")
// 		// if data != nil {
// 		// 	(*data)["playername"] = (h.Data.Data.PlayerData.Data.Name)
// 		// }
// 		event <- Event{Type: "handshake_hostname", Data: map[string]string{"hostname": h.Data.Data.Hostname}}
// 		return sm.setState(&PassthroughState{})
// 	}
// }
