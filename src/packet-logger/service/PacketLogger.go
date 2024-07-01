package service

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "mc_reverse_proxy/src/common"
	config "mc_reverse_proxy/src/configuration/service"
	utils "mc_reverse_proxy/src/utils"

	"time"
)

type PacketLoggable interface {
	Register(*PacketLogger) error
}

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
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}
	splitted := strings.Split(config.LoggerMongoDBAddress, "@")
	utils.FLog.PacketLogger("Connected to MongoDB at %s using database \"%s\" and collection \"%s\"", strings.ReplaceAll(splitted[len(splitted)-1], "/", ""), config.LoggerMongoDBName, config.LoggerMongoColName)
	packetLogger = &PacketLogger{client: client, config: config}
	return nil
}

func Send(data PacketLog) {
	if packetLogger != nil {
		go func() {
			ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)
			_, err := packetLogger.client.Database(packetLogger.config.LoggerMongoDBName).Collection(packetLogger.config.LoggerMongoColName).InsertOne(ctx, data)
			if err != nil {
				utils.FLogErr.PacketLogger("Inserting error: %v", err)
			}
		}()
	}
}
