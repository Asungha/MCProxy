package service

import (
	"encoding/json"
	"mc_reverse_proxy/src/packet"
)

func ModifyStatusMessage(packet *packet.Status) (*packet.Status, error) {
	jsonData := map[string]interface{}{}
	err := json.Unmarshal([]byte(packet.Json), &jsonData)
	if err != nil {
		return packet, err
	}
	jsonData["description"] = map[string]string{"text": "Proxy server"}
	jsonData["players"] = map[string]int{"max": 40, "online": 2}

	newJson, err := json.Marshal(jsonData)
	if err != nil {
		return packet, err
	}
	packet.Json = string(newJson)
	return packet, nil
}
