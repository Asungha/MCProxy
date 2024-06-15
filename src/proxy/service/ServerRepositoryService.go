package service

import (
	proto "mc_reverse_proxy/src/control/controlProto"
	controlService "mc_reverse_proxy/src/control/service"
	"time"
)

type ServerMetadata struct {
	AdvertisedHostname string
}

type ChallengeSessionData struct {
	Requester string
	Challenge string
	TTL       time.Time
	Signature string
}

func (d *ChallengeSessionData) Valid() bool {
	return true
}

type ControlPluginSupport interface {
	*controlService.EventService

	NewChallengeSession(requester string) ChallengeSessionData
	VerifyChallenge(proto.AuthRequest) bool

	GetEventChannel(topic string, hostname string) chan controlService.EventData

	SendCommand(hostname string, data *proto.CommandData) error
	GetMetricChannel(hostname string) chan controlService.EventData
}
type ServerRepositoryService interface {
	Load() error
	Resolve(hostname string) (string, error)
}

type UpdatableRepositoryService interface {
	Upsert(id int, hostname string, address string) error
	Insert(hostname string, address string) error
	Delete(id int) error
	Count() (int, error)
	Destroy()
	ServerRepositoryService
}

type ServerList struct {
	ID       int
	Hostname string
	Address  string
}
type ListableRepositoryService interface {
	List() ([]ServerList, error)
	ServerRepositoryService
}
