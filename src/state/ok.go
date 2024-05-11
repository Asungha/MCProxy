package state

import (
	"errors"
	pac "mc_reverse_proxy/src/packet"
)

type OKState struct {
	data pac.Packet[*pac.Raw]
}

func (o *OKState) ImplState() {}

func (o *OKState) Enter(data pac.Packet[*pac.Raw]) error {
	res := pac.Packet[*pac.Raw]{}
	err := pac.CastPacket(data, &res)
	if err != nil {
		return err
	}
	o.data = res
	return nil
}

func (o *OKState) Exit() {}

func (o *OKState) Validate() error {
	return nil
}

func (o *OKState) Update(sm *StateMachine, data pac.Packet[pac.IPacketData]) error {
	return errors.New("State Done")
}
