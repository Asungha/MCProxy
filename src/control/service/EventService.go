package service

import (
	"log"
	"time"

	proto "mc_reverse_proxy/src/control/controlProto"

	"github.com/google/uuid"
)

type EventData struct {
	MetricData  *proto.MetricData
	CommandData *proto.CommandData
}

type EventService struct {
	Subscriber *ThreadSafeMap[*map[string]chan EventData]
	buffer     int
}

func (e *EventService) Subscribe(topic string) (string, chan EventData) {
	sub := make(chan EventData, e.buffer)
	uid := uuid.New().String()
	if hasTopic := e.Subscriber.Contain(topic); hasTopic {
		e.Subscriber.Modify(func(m map[string]*map[string]chan EventData) map[string]*map[string]chan EventData {
			(*m[topic])[uid] = sub
			log.Printf("Put new subscriber: %s", uid)
			return m
		})
	} else {
		e.Subscriber.Set(topic, &map[string]chan EventData{uid: sub})
	}
	log.Printf("Put new subscriber: %s", uid)
	return uid, sub
}

func (e *EventService) Unsubscribe(uid string, topic string) {
	e.Subscriber.Modify(func(m map[string]*map[string]chan EventData) map[string]*map[string]chan EventData {
		delete((*m[topic]), uid)
		if len(*m[topic]) == 0 {
			delete(m, topic)
		}
		return m
	})
}

func (e *EventService) Publish(topic string, data EventData) {
	if !e.Subscriber.Contain(topic) {
		e.Subscriber.Set(topic, &map[string]chan EventData{})
	}
	for uid, channel := range *e.Subscriber.Get(topic) {
		log.Printf("Sendding data to %s", uid)
		select {
		case <-time.After(10 * time.Millisecond):
			log.Printf("[Event Service] Warning: sending data to %s timeout", uid)
		case channel <- data:
		}
	}
}

func NewEventService(buffer int) *EventService {
	return &EventService{
		Subscriber: NewThreadSafeMap[*map[string]chan EventData](),
		buffer:     buffer,
	}
}
