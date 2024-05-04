package state

import (
	"log"
	packet "mc_reverse_proxy/packet"
	service "mc_reverse_proxy/service"
)

type RequestStatusState struct {
	connData     *ConnectionData
	action       int
	originalData *[]byte
	modifiedData []byte
}

func (r *RequestStatusState) ImplState() {}

func (r *RequestStatusState) Enter(connData *ConnectionData) {
	r.connData = connData
	select {
	case clientData := <-r.connData.ClientData:
		r.originalData = clientData.bytes
		p := packet.NewPacket(&packet.Status{})
		err := p.Decode(clientData.bytes, clientData.length)
		if err != nil {
			log.Printf("Error decoding packet in Request status state: %v", err)
			r.action = ACTION_TRANSPARENT
			return
		}
		_, err = service.ModifyStatusMessage(p.Data)
		if err != nil {
			log.Printf("Error modifying status message in Request status state: %v", err)
			r.action = ACTION_TRANSPARENT
			return
		}
		r.modifiedData = p.Encode()
		if err != nil {
			log.Printf("Error encoding packet in Request status state: %v", err)
			r.action = ACTION_TRANSPARENT
			return
		}
		r.action = ACTION_ACCEPT
	case <-r.connData.ctx.Done():
		r.action = ACTION_CANCLE
	}
}

func (r *RequestStatusState) Exit() {}

func (r *RequestStatusState) Validate() error {
	return nil
}

func (r *RequestStatusState) Update(sm *StateMachine) error {
	switch r.action {
	case ACTION_TRANSPARENT:
		sm.setState(&ResponseStatusState{IsTranparent: true})
		r.connData.ServerData <- Data{bytes: r.originalData, length: len(r.modifiedData)}
	case ACTION_ACCEPT:
		sm.setState(&ResponseStatusState{IsTranparent: false})
		// r.connData.ServerData <- Data{bytes: &r.modifiedData, length: len(r.modifiedData)}
	case ACTION_CANCLE:
		sm.setState(&RejectState{err: nil})
	}
	return nil
}

type ResponseStatusState struct {
	connData     *ConnectionData
	action       int
	err          error
	IsTranparent bool
}

func (r *ResponseStatusState) ImplState() {}

func (r *ResponseStatusState) Enter(connData *ConnectionData) {
	r.connData = connData
	if r.IsTranparent { // If the data is transparent, just forward it to the client
		select {
		case serverData := <-r.connData.ServerData:
			r.connData.ClientData <- serverData
			r.action = ACTION_ACCEPT
		case <-r.connData.ctx.Done():
			r.action = ACTION_CANCLE
			return
		}
	} else { // Otherwise, send a custom response to the client

	}
}

func (r *ResponseStatusState) Exit() {}

func (r *ResponseStatusState) Validate() error {
	return nil
}

func (r *ResponseStatusState) Update(sm *StateMachine) error {
	switch r.action {
	case ACTION_CANCLE:
		sm.setState(&RejectState{err: nil})
	case ACTION_ACCEPT:
		sm.setState(&OKState{})
	}
	return nil
}
