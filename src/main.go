package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	pac "mc_reverse_proxy/src/packet"
	"mc_reverse_proxy/src/utils"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Status struct {
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			Name string `json:"name"`
			ID   string `json:"id"`
		} `json:"sample"`
	} `json:"players"`
	Description interface{} `json:"description"`
	Modinfo     struct {
		Type    string   `json:"type"`
		ModList []string `json:"modList"`
	} `json:"modinfo"`
}

func (s *Status) JSONString() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

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

func GetPlayerStatus(host map[string]map[string]string) {
	for {
		for k, v := range host {
			server, err := net.Dial("tcp", v["host"]+":"+v["port"])
			if err != nil {
				log.Printf("Failed to connect to backend server: %v", err)
				continue
			}

			// make ping packet
			p := pac.NewPacket(&pac.Handshake{
				ProtocolVersion: 21,
				Hostname:        k,
				Port:            25565,
				NextState:       1,
			})

			p.ID = 0x00
			server.SetWriteDeadline(time.Now().Add(time.Second * 5))
			// log.Printf("Sending handshake to %s", k)
			log.Printf("Payload: %s", p.Encode())
			_, err = server.Write(p.Encode())
			if err != nil {
				log.Printf("Failed to send handshake data to backend: %v", err)
			}

			//read response
			backend_buffer := make([]byte, 40960)
			server.SetReadDeadline(time.Now().Add(time.Second * 20))
			n, err := server.Read(backend_buffer)
			if err != nil {
				log.Printf("Failed to read data from backend: %v readed %d", err, n)
			}
			log.Printf("Read %d bytes from backend %s : %x", n, server.RemoteAddr().String(), backend_buffer[:n])
		}
		time.Sleep(5 * time.Second)
	}
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

	backends := map[string]string{}
	for k, v := range host {
		backends[k] = v["host"] + ":" + v["port"]
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

			go handleClientIntercept(&clientConn, logChann, listenAddr, commanChann, backends)
		}
	}()
}

func handleClientIntercept(clientConn *net.Conn, logChann chan string, listenAddr string, commanChann map[string]chan string, backends map[string]string) {
	buf := make([]byte, 8192)

	n, err := (*clientConn).Read(buf)

	var hostname string

	if err != nil {
		log.Printf("Failed to read data from client: %v", err)
	} else {
		strPayload := fmt.Sprintf("%x", buf[5:][:n])
		port, _ := strconv.ParseInt(strings.Split(listenAddr, ":")[1], 10, 64)
		portHex := fmt.Sprintf("%x", port)

		nestStateFlag := "0"

		hostnameHex := strings.Split(strPayload, portHex+nestStateFlag)[0]
		strHostname, err := hex.DecodeString(hostnameHex)
		if err != nil {
			panic(err)
		}
		re := regexp.MustCompile(`\b([a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+)\b`)
		matches := re.FindStringSubmatch(string(strHostname))
		if len(matches) > 1 {
			hostname = matches[1]
		}
	}
	if len(hostname) == 0 {
		if strings.Contains(utils.ByteToStr(buf[:n]), "#") && strings.Contains(utils.ByteToStr(buf[:n]), "|") {
			hostname = strings.Split(utils.ByteToStr(buf[:n]), "#")[1]
		}
	}
	logChann <- fmt.Sprintf("Accepted client connection from %s connecting to %s", (*clientConn).RemoteAddr().String(), hostname)
	if backends[hostname] == "" {
		logChann <- fmt.Sprintf("Failed to find backend server for hostname: %s", hostname)
		(*clientConn).Close()
		return
	}
	backendConn, err := net.Dial("tcp", backends[hostname])
	if err != nil {
		logChann <- fmt.Sprintf("Failed to connect to backend server: %v", err)
		(*clientConn).Close()
		return
	}
	_, err = backendConn.Write(buf[:n])
	if err != nil {
		log.Printf("Failed to write data to backend: %v", err)
	}
	commanChann[(*clientConn).RemoteAddr().String()] = make(chan string)
	go handleBackendIntercept(clientConn, &backendConn, logChann, commanChann[(*clientConn).RemoteAddr().String()])

	for {
		client_buffer := make([]byte, 512)
		n, err := (*clientConn).Read(client_buffer)
		if err != nil {
			break
		}
		backendConn.Write(client_buffer[:n])
	}
	(*clientConn).Close()
	backendConn.Close()
}

func handleBackendIntercept(clientConn *net.Conn, backendConn *net.Conn, logChann chan string, cmd chan string) {
	data := make(chan []byte)
	errs := make(chan error)
	getData := func() {
		for {
			backend_buffer := make([]byte, 1024*1024)
			n, err := (*backendConn).Read(backend_buffer)
			if err != nil {
				errs <- err
				break
			}
			data <- backend_buffer[:n]
			backend_buffer = nil
		}
	}
	go getData()
	for {
		select {
		case backend_buffer := <-data:
			p := pac.NewPacket(&pac.Raw{})
			err := p.Decode(&backend_buffer, len(backend_buffer))
			if err != nil {
				log.Printf("Failed to decode packet: %v", err)
			} else {
				// log.Printf("Server Data length: %d, ID %x", p.Length(), p.ID)
			}
			(*clientConn).Write(backend_buffer)
		case command := <-cmd:
			switch strings.Split(command, " ")[0] {
			case "disconnect":
				logChann <- fmt.Sprintf("Disconnecting client %s", (*clientConn).RemoteAddr().String())
				(*backendConn).Close()
				cmd <- "done"
				return
			}
		case <-errs:
			(*clientConn).Close()
			(*backendConn).Close()
			break
		}
	}
}

func logger(logChann chan string) {
	for {
		log.Printf(<-logChann)
	}
}
