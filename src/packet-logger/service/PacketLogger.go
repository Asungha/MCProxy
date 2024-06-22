package service

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	config "mc_reverse_proxy/src/configuration/service"

	"time"
)

type PacketLoggable interface {
	Register(*PacketLogger) error
}

type PacketType string

const (
	UNKNOWN      PacketType = "unknown"
	MC_HANDSHAKE PacketType = "mc_handshake"
	MC_LOGIN     PacketType = "mc_login"
	MC_GAMMPLAY  PacketType = "mc_gameplay"
	MC_OTHER     PacketType = "mc_other"
	HTTP         PacketType = "http"
)

type PacketLog struct {
	Type      PacketType `json:"type" bson:"type"`
	IP        string     `json:"ip" bson:"ip"`
	Port      string     `json:"port" bson:"port"`
	Timestamp time.Time  `json:"timestamp" bson:"timestamp"`
	Data      string     `json:"data" bson:"data"`
}

var packetLogger *PacketLogger

type PacketLogger struct {
	client *mongo.Client
	config *config.ConfigurationService
}

func InitPacketLogger(config *config.ConfigurationService) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.LoggerMongoDBAddress).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return err
	}
	log.Println("[Packet Logger] Connected to mongodb")
	packetLogger = &PacketLogger{client: client, config: config}
	return nil
}

func Send(data PacketLog) {
	if packetLogger != nil {
		ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
		_, err := packetLogger.client.Database(packetLogger.config.LoggerMongoDBName).Collection(packetLogger.config.LoggerMongoColName).InsertOne(ctx, data)
		if err != nil {
			log.Printf("[Packet Logger] Inserting error : %v", err)
		}
	} else {
		log.Printf("[Packet Logger] Packet Logger not available")
	}
}
