package service

import (
	"reflect"
	"testing"
	"time"
)

func TestEventService_Subscribe(t *testing.T) {
	type args struct {
		topic string
	}
	tests := []struct {
		name string
		e    *EventService
		args args
	}{
		{
			name: "No available topic",
			args: args{
				topic: "test",
			},
			e: NewEventService(0),
		},
		{
			name: "topic available (created by other subscriber)",
			args: args{
				topic: "test",
			},
			e: func() *EventService {
				s := NewEventService(0)
				s.Subscribe("test")
				return s
			}(),
		},
		{
			name: "topic available 1 (created by publisher)",
			args: args{
				topic: "test",
			},
			e: func() *EventService {
				s := NewEventService(0)
				s.Publish("test", EventData{})
				return s
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)
			_, channel := tt.e.Subscribe(tt.args.topic)
			go func() {
				select {
				case <-time.After(1 * time.Second):
					t.Errorf("EventService.Subscribe() timeout")
					done <- true
				case <-channel:
					done <- true
				}
			}()
			tt.e.Publish(tt.args.topic, EventData{})
			<-done
		})
	}
}

func TestEventService_Unsubscribe(t *testing.T) {
	type args struct {
		uid   string
		topic string
	}
	tests := []struct {
		name string
		e    *EventService
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.Unsubscribe(tt.args.uid, tt.args.topic)
		})
	}
}

func TestEventService_Publish(t *testing.T) {
	type args struct {
		topic string
		data  EventData
	}
	tests := []struct {
		name            string
		e               *EventService
		args            args
		timeoutExpected bool
	}{
		{
			name: "No available topic",
			e:    NewEventService(0),
			args: args{
				topic: "test",
				data:  EventData{},
			},
			timeoutExpected: false,
		},
		{
			name: "topic available (created by subscriber)",
			e: func() *EventService {
				s := NewEventService(0)
				return s
			}(),
			args: args{
				topic: "test",
				data:  EventData{},
			},
			timeoutExpected: false,
		},
		{
			name: "subscriber sub-unsub",
			e: func() *EventService {
				s := NewEventService(0)
				uid, _ := s.Subscribe("test")
				s.Unsubscribe(uid, "test")
				return s
			}(),
			args: args{
				topic: "test",
				data:  EventData{},
			},
			timeoutExpected: false,
		},
		{
			name: "subscriber not listen on unbuffered chan",
			e: func() *EventService {
				s := NewEventService(0)
				s.Subscribe("test")
				return s
			}(),
			args: args{
				topic: "test",
				data:  EventData{},
			},
			timeoutExpected: true,
		},
		{
			name: "subscriber not listen on buffered chan",
			e: func() *EventService {
				s := NewEventService(1)
				s.Subscribe("test")
				return s
			}(),
			args: args{
				topic: "test",
				data:  EventData{},
			},
			timeoutExpected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)
			go func() {
				tt.e.Publish(tt.args.topic, tt.args.data)
				done <- true
			}()
			select {
			case <-time.After(10 * time.Millisecond):
				if !tt.timeoutExpected {
					t.Errorf("EventService.Subscribe() timeout unexpected")
				}
			case <-done:
				if tt.timeoutExpected {
					t.Errorf("EventService.Subscribe() timeout expected, but no timeout occured")
				}
			}
		})
	}
}

func TestNewEventService(t *testing.T) {
	tests := []struct {
		name string
		want *EventService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventService(0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEventService() = %v, want %v", got, tt.want)
			}
		})
	}
}
