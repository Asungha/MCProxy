package state

import (
	"errors"
	"fmt"
	"log"
	pac "mc_reverse_proxy/src/packet"
)

type HandshakeState struct {
	sm       *StateMachine
	Data     pac.Packet[*pac.Handshake]
	hostname string
}

func (h *HandshakeState) ImplState() {}

func (h *HandshakeState) Enter(sm *StateMachine) error {
	log.Printf("[state machine: Handshake] Enter handshake state")
	h.sm = sm
	select {
	case rawData := <-h.sm.Conn.ClientData:
		hs := pac.Handshake{}
		data := pac.Packet[*pac.Handshake]{Data: &hs}
		err := data.Decode(&rawData, len(rawData))
		if err != nil {
			return err
		}
		h.hostname = data.Data.Hostname
		h.Data = data
		return nil
	case <-h.sm.ctx.Done():
		return errors.New("Context Done")
	}
}

func (h *HandshakeState) Action() error {
	if data, ok := h.sm.Conn.ServerList[h.hostname]; ok {
		if target, ok := data["target"]; ok {
			err := h.sm.Conn.ConnectServer(target)
			if err != nil {
				h.sm.errorMetric.ServerConnectFailed += 1
				return err
			}
			err = h.sm.Conn.PreConditionCheck()
			if err != nil {
				log.Printf("[handshake state] Precondition failed %v", err)
				return err
			}
			go h.sm.Conn.ListenServer()
			hs_packet, err := h.Data.Encode()
			if err != nil {
				h.sm.errorMetric.PacketDeserializeFailed += 1
				log.Printf("[handshake state] Encode handshale failed %v", err)
				return err
			}
			h.sm.StateChangeLock.Lock()
			err = h.sm.Conn.WriteServer(hs_packet)
			if err != nil {
				log.Printf("[handshake state] Handshake packet send failed %v", err)
				return err
			}
			if h.Data.Data.Tail != nil {
				err = h.sm.Conn.WriteServer(h.Data.Data.Tail)
				if err != nil {
					log.Printf("[handshake state] additional packet send failed %v", err)
					return err
				}
			}
			h.sm.StateChangeLock.Unlock()
			return nil
		}
		return errors.New("[Handshake State] host config file malformed")
	}
	h.sm.errorMetric.HostnameResolveFailed += 1
	return errors.New(fmt.Sprintf("[Handshake State] Host %s not found", h.hostname))
}

func (h *HandshakeState) Exit() IState {
	err := h.sm.Conn.PreConditionCheck()
	if err != nil {
		log.Printf("[state machine: Handshake] Error: %v", err)
		return nil
	}
	if h.Data.Data.NextState == 0x01 {
		if h.hostname != "" {
			h.sm.proxyMetric.PlayerGetStatus += 1
			log.Printf("[state machine: Handshake] Handshake for status request done")
			return &PassthroughState{}
		} else {
			log.Printf("[state machine: Handshake] Error: Hostname not provided")
			return nil
		}
	} else {
		h.sm.proxyMetric.PlayerLogin += 1
		h.sm.proxyMetric.PlayerPlaying += 1
		h.sm.PlayerPlaying = true
		log.Printf("[state machine: Handshake] Handshake for login request done")
		return &PassthroughState{}
	}
}
