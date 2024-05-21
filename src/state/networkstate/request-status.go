package state

// import (
// 	"bytes"
// 	"log"
// 	pac "mc_reverse_proxy/src/packet"
// )

// type SendStatus struct {
// 	data          pac.Packet[*pac.Status]
// 	isOldProtocol bool
// }

// func (h *SendStatus) ImplState() {}

// func (h *SendStatus) Enter(data []byte) error {
// 	log.Printf("SendStatus Enter")
// 	if bytes.Equal(data, []byte{0x01, 0x00}) {
// 		// log.Printf("Old protocol detected")
// 		h.isOldProtocol = true
// 		return nil
// 	}
// 	res := pac.Packet[*pac.Status]{Data: &pac.Status{}}
// 	// res := pac.NewPacket(pac.NewHandshake())
// 	err := res.Decode(&data, len(data))
// 	if err != nil {
// 		return err
// 	}
// 	h.data = res
// 	return nil
// }
// func (h *SendStatus) Exit() {}

// func (h *SendStatus) Validate() error {
// 	return nil
// }

// func (h *SendStatus) Action() (*[]byte, error) {
// 	if h.isOldProtocol {
// 		return nil, nil // skip
// 	}
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

// 	// log.Printf("Status packet: %s", h.data.Data.String())
// 	// log.Printf("Status packet: %v", h.data.Data)
// 	json, _ := h.data.Data.JSON()
// 	// json.Players.Online = 0
// 	// json.Players.Max = 0
// 	json.Description.Extra = []pac.StatusDesExtra{}
// 	builder := pac.DescriptionBuilder{}
// 	builder.Add("§4M§cA§6I§2B§aO§bR§3K")
// 	json.Description.Extra = builder.Build()
// 	h.data.Data.SetJSON(json)

// 	// log.Printf("Status packet: %v", json)
// 	// h.data.Data.Json = strings.Replace(h.data.Data.Json, "\"online\":0", "\"online\":555555555", -1)
// 	data, err := h.data.Encode()
// 	return &data, err
// }

// func (h *SendStatus) Update(sm *StateMachine, data *map[string]string, event chan Event, ack chan bool) error {
// 	// if h.isOldProtocol {
// 	// 	return sm.setState(&SendStatus{})
// 	// } else {
// 	// 	return sm.setState(&PassthroughState{})
// 	// }
// 	return nil
// }
