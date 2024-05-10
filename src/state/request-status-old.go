package state

// import (
// 	pac "mc_reverse_proxy/src/packet"
// )

// type OldReqStatus struct {
// 	data pac.OldStatusReq
// }

// func (h *OldReqStatus) ImplState() {}

// func (h *OldReqStatus) Enter(data []byte) error {
// 	// // res := pac.Packet[*pac.Handshake]{}
// 	// res := pac.NewPacket(pac.NewHandshake())
// 	// err := res.Decode(&data, len(data))
// 	// if err != nil {
// 	// 	return err
// 	// } else {
// 	// 	log.Printf("Handshake packet: %s", res.Data.String())
// 	// }
// 	// h.data = res
// 	// log.Printf("OldReqStatus Enter")
// 	var res pac.OldStatusReq
// 	err := res.Decode(data, len(data))
// 	if err != nil {
// 		return err
// 	}
// 	h.data = res
// 	return nil
// }
// func (h *OldReqStatus) Exit() {}

// func (h *OldReqStatus) Validate() error {
// 	return nil
// }

// func (h *OldReqStatus) Action() (*[]byte, error) {
// 	// res := pac.Packet[*pac.Raw]{Data: &pac.Raw{}}
// 	// err := res.Decode(&h.data, len(h.data))
// 	// if err != nil {
// 	// 	log.Printf("Error decoding packet in Passthrough state: %v", err)
// 	// } else {
// 	// 	log.Printf("[Action] Play packet: %s", res.Data.String())
// 	// }
// 	// err := pac.AsRawPacket(h.data, &res)
// 	// return res, nil
// 	// log.Printf("[Action] Handshake packet: %v", h.data.Data)
// 	// return h.data.Encode()
// 	data, err := h.data.Encode()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &data, nil
// }

// func (h *OldReqStatus) Update(sm *StateMachine, data map[string]string) error {
// 	return sm.setState(&PassthroughState{})

// }
