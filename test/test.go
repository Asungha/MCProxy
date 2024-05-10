package main

import (
	"encoding/hex"
	"log"
	pac "mc_reverse_proxy/src/packet"
	"net"
)

func main() {
	// Start the proxy
	hs := pac.NewPacket(pac.NewHandshake())
	hs.Data.ProtocolVersion = 765
	hs.Data.Hostname = "mc.billybacker.cc"
	hs.Data.NextState = 2
	hs.Data.Port = 25565

	p_data := pac.PlayerData{}
	p_data.Name = "test"
	uuid, _ := hex.DecodeString("2d0af11c88644d01a9d255358ab060ff")
	p_data.UUID = uuid
	packet := pac.NewPacket(&p_data)
	hs.Data.PlayerData = &packet

	log.Printf("Handshake packet: %s", hs.Data.String())

	upstreamConn, err := net.Dial("tcp", "127.0.0.1:25565")
	if err != nil {
		log.Printf("Failed to connect to upstream server: %v", err)
	}
	defer upstreamConn.Close()
	b, err := hs.Encode()
	upstreamConn.Write(b)

	for {
		buffer := make([]byte, 40960)
		n, err := upstreamConn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from upstream connection: %v", err)
			return
		}
		log.Printf("Read %d bytes from upstream: %x", n, buffer[:n])
	}
}
