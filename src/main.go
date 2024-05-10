package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	pac "mc_reverse_proxy/src/packet"
	"mc_reverse_proxy/src/utils"
	"net"
	"os"
	"strings"
	"sync"
)

type Connector struct {
	client *net.Conn
	server *net.Conn
}

func (c *Connector) handleClient() {
	_, err := io.Copy(*c.server, *c.client)
	if err != nil {
		log.Printf("Failed to copy data from client to backend: %v", err)
	}

	(*c.client).Close()
	(*c.server).Close()
}

func (c *Connector) handleServer() {
	_, err := io.Copy(*c.client, *c.server)
	if err != nil {
		log.Printf("Failed to copy data from backend to client: %v", err)
	}

	(*c.client).Close()
	(*c.server).Close()
}

func (c *Connector) start() {
	go c.handleClient()
	go c.handleServer()
}

func NewConnector(client *net.Conn, server *net.Conn) *Connector {
	return &Connector{
		client: client,
		server: server,
	}
}

type ServerStatus struct {
	ServerName    string `json:"server_name"`
	OnlinePlayers int    `json:"online_players"`
	MaxPlayers    int    `json:"max_players"`
}

type MinecraftPacket struct {
	PacketID   int    `json:"packet_id"`
	Payload    string `json:"payload"`
	PlayerName string `json:"player_name"`
}

type HandshakePacket struct {
	ProtocolVersion int    `json:"protocol_version"`
	ServerAddress   string `json:"server_address"`
	ServerPort      int    `json:"server_port"`
	NextState       int    `json:"next_state"`
}

type ConnectionPair struct {
	Client *net.Conn
	Server *net.Conn
}

// func GetPlayerStatus(host map[string]map[string]string) {
// 	for {
// 		for k, v := range host {
// 			server, err := net.Dial("tcp", v["host"]+":"+v["port"])
// 			if err != nil {
// 				log.Printf("Failed to connect to backend server: %v", err)
// 				continue
// 			}

// 			// make ping packet
// 			p := pac.NewPacket(&pac.Handshake{
// 				ProtocolVersion: 21,
// 				Hostname:        k,
// 				Port:            25565,
// 				NextState:       1,
// 			})

// 			p.ID = 0x00
// 			server.SetWriteDeadline(time.Now().Add(time.Second * 5))
// 			// log.Printf("Sending handshake to %s", k)
// 			log.Printf("Payload: %s", p.Encode())
// 			_, err = server.Write(p.Encode())
// 			if err != nil {
// 				log.Printf("Failed to send handshake data to backend: %v", err)
// 			}

// 			//read response
// 			backend_buffer := make([]byte, 40960)
// 			server.SetReadDeadline(time.Now().Add(time.Second * 20))
// 			n, err := server.Read(backend_buffer)
// 			if err != nil {
// 				log.Printf("Failed to read data from backend: %v readed %d", err, n)
// 			}
// 			log.Printf("Read %d bytes from backend %s : %x", n, server.RemoteAddr().String(), backend_buffer[:n])
// 		}
// 		time.Sleep(5 * time.Second)
// 	}
// }

func HandleData(client *net.Conn, server *net.Conn, done chan bool) {
	stop := make(chan bool)
	clientData := make(chan []byte)
	serverData := make(chan []byte)
	lock := sync.Mutex{}

	// listen for client data
	go func() {
		for {
			buf := make([]byte, 512)
			n, err := (*client).Read(buf)
			if err != nil {
				log.Printf("Failed to read data from client: %v", err)
				stop <- true
				break
			}
			lock.Lock()
			clientData <- buf[:n]
			lock.Unlock()
			buf = nil
		}
		log.Printf("Client data handler closed")
	}()

	// listen for server data
	go func() {
		for {
			buf := make([]byte, 512)
			n, err := (*server).Read(buf)
			if err != nil {
				log.Printf("Failed to read data from server: %v", err)
				stop <- true
				break
			}
			lock.Lock()
			serverData <- buf[:n]
			lock.Unlock()
			buf = nil
		}
		log.Printf("Server data handler closed")
	}()

	// Data handler
	go func() {
		for {
			select {
			case client_buffer := <-clientData:
				// log.Printf("Client data: %x", client_buffer)
				_, err := (*server).Write(client_buffer)
				if err != nil {
					log.Printf("Failed to write data to server: %v", err)
					stop <- true
					return
				}
			case server_buffer := <-serverData:
				// log.Printf("Server data: %x", server_buffer)
				_, err := (*client).Write(server_buffer)
				if err != nil {
					log.Printf("Failed to write data to client: %v", err)
					stop <- true
					return
				}
			case <-stop:
				return
			}
		}
	}()
	<-stop
	log.Printf("Data handler closed")
	stop = nil
	clientData = nil
	serverData = nil
	done <- true
}

func main() {
	//read json file
	host_file, err := os.Open("host.json")
	if err != nil {
		log.Fatalf("Failed to open host config file: %v", err)
	}
	defer host_file.Close()

	host := make(map[string]map[string]string)
	decoder := json.NewDecoder(host_file)
	err = decoder.Decode(&host)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	backends := map[string]map[string]string{}
	for k, v := range host {
		backends[k] = make(map[string]string)
		backends[k]["target"] = v["ip"] + ":" + v["port"]
		backends[k]["hostname"] = v["hostname"]
	}

	config := map[string]string{}
	config_file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer config_file.Close()

	decoder = json.NewDecoder(config_file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	listenAddr := config["listen"] + ":" + config["port"]

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	defer listener.Close()

	log.Printf("Proxy server listening on %s", listenAddr)

	logChann := make(chan string)
	go logger(logChann)
	lock := sync.Mutex{}
	commanChann := make(map[string]chan string)

	log.Printf("Server started")
	func() {
		for {
			// Accept incoming client connections
			clientConn, err := listener.Accept()
			if err != nil {
				log.Printf("Failed to accept client connection: %v", err)
				continue
			}
			go handleClientIntercept(&clientConn, logChann, listenAddr, commanChann, &lock, backends)
		}
	}()
}

func handleClientIntercept(clientConn *net.Conn, logChann chan string, listenAddr string, commanChann map[string]chan string, lock *sync.Mutex, backends map[string]map[string]string) {
	buf := make([]byte, 512)
	// remove buffer after done
	unregister := func() {
		lock.Lock()
		delete(commanChann, (*clientConn).RemoteAddr().String())
		lock.Unlock()
	}
	defer unregister()

	n, err := (*clientConn).Read(buf)

	var hostname string

	if err != nil {
		log.Printf("Failed to read data from client: %v", err)
	} else {
		hs := pac.Handshake{}
		pac := pac.NewPacket(&hs)
		defer pac.Destroy()
		err := pac.Decode(&buf, n)
		if err != nil {
			log.Printf("Failed to decode packet: %v", err)
		} else {
			hostname = pac.Data.Hostname
		}
		pac.Destroy()
	}
	var ip string
	var newHostname string
	if len(hostname) == 0 {
		if strings.Contains(utils.ByteToStr(buf[:n]), "#") && strings.Contains(utils.ByteToStr(buf[:n]), "|") {
			hostname = strings.Split(utils.ByteToStr(buf[:n]), "#")[1]
		}
	} else {
		ip = backends[hostname]["target"]
		if v, ok := backends[hostname]; ok {
			newHostname = v["hostname"]
		} else {
			newHostname = hostname
		}

	}
	logChann <- fmt.Sprintf("Accepted client connection from %s connecting to %s", (*clientConn).RemoteAddr().String(), hostname)
	if _, ok := backends[hostname]; !ok {
		if _, ok := backends["default"]; ok {
			logChann <- fmt.Sprintf("Using default backend server %s for hostname: %s", backends["default"]["hostname"], hostname)
			ip = backends["default"]["target"]
			newHostname = backends["default"]["hostname"]
		} else {
			logChann <- fmt.Sprintf("Failed to find backend server for hostname: %s", hostname)
			(*clientConn).Close()
			clientConn = nil
			return
		}
	}
	backendConn, err := net.Dial("tcp", ip)
	if err != nil {
		logChann <- fmt.Sprintf("Failed to connect to backend server: %v", err)
		(*clientConn).Close()
		clientConn = nil
		return
	}
	r := pac.Raw{}
	packet := pac.NewPacket(&r)
	err = packet.Decode(&buf, n)
	if err != nil {
		log.Printf("Failed to decode packet: %v", err)
		_, err = backendConn.Write(buf[:n])
		if err != nil {
			log.Printf("Failed to write data to backend: %v", err)
		}
	} else {
		log.Printf("New Hostname: %s", newHostname)
		b := buf[:n]
		hs := pac.Handshake{}
		p := pac.NewPacket(&hs)
		err := p.Decode(&b, n)
		if err != nil {
			log.Printf("Failed to decode packet: %v", err)
		}
		p.Data.Hostname = newHostname
		data, err := p.Encode()
		if err != nil {
			log.Printf("Failed to encode packet: %v", err)
		}
		n_w, err := backendConn.Write(data)
		if err != nil {
			log.Printf("Failed to write data to backend: %v", err)
		} else {
			log.Printf("%d Packet sent to backend", n_w)
		}
		p.Destroy()
	}
	buf = nil

	lock.Lock()
	commanChann[(*clientConn).RemoteAddr().String()] = make(chan string)
	lock.Unlock()
	log.Printf("Client connected to backend")

	done := make(chan bool)
	go HandleData(clientConn, &backendConn, done)
	<-done
	(*clientConn).Close()
	backendConn.Close()
	log.Printf("Client disconnected from backend")
}

// func handleBackendIntercept(clientConn *net.Conn, backendConn *net.Conn, logChann chan string, cmd chan string) {
// 	data := make(chan []byte)
// 	errs := make(chan error)
// 	defer func() {
// 		data = nil
// 		errs = nil
// 	}()
// 	backend_buffer := make([]byte, 1024*1024)
// 	getData := func() {
// 		for {
// 			n, err := (*backendConn).Read(backend_buffer)
// 			if err != nil {
// 				errs <- err
// 				logChann <- fmt.Sprintf("Backend %s disconnected err %s", (*backendConn).RemoteAddr().String(), err.Error())
// 				break
// 			}
// 			data <- backend_buffer[:n]
// 		}
// 	}
// 	go getData()
// 	for {
// 		select {
// 		case backend_buffer := <-data:
// 			// log.Printf("Data from backend: %x", backend_buffer)
// 			// p := pac.NewPacket(&pac.Raw{})
// 			// _ := p.Decode(&backend_buffer, len(backend_buffer))
// 			// if err != nil {
// 			// 	log.Printf("Failed to decode packet: %v", err)
// 			// } else {
// 			// 	// log.Printf("Server Data length: %d, ID %x", p.Length(), p.ID)
// 			// }
// 			(*clientConn).Write(backend_buffer)
// 		case command := <-cmd:
// 			switch strings.Split(command, " ")[0] {
// 			case "disconnect":
// 				logChann <- fmt.Sprintf("Disconnecting client %s", (*clientConn).RemoteAddr().String())
// 				cmd <- "done"
// 				backend_buffer = nil
// 				return
// 			}
// 		case <-errs:
// 			backend_buffer = nil
// 			return
// 		}
// 	}
// }

func logger(logChann chan string) {
	for {
		log.Printf(<-logChann)
	}
}
