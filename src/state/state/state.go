package state

import (
	"errors"
	"time"
)

type TransistionFunction func(IState) (IState, error)
type Function func(IState) error
type ActionFunction func(IState) error

// type StateID string

type StateMetadata map[string]interface{}

func CastMetadata[T any](data interface{}) *T {
	casted, ok := data.(T)
	if !ok {
		return nil
	}
	return &casted
}

type IState interface {
	Init(*Function) IState

	Enter() error
	Transition() (IState, error)

	AddTransistionFunction(*TransistionFunction)

	SetMetadata(string, any)
	GetMetadata(string) any

	SetTimeout(time.Duration)
	GetTimeout() time.Duration
	UseTimeout() bool

	Halt() error
	Destruct() error
	GetTFunc() []*TransistionFunction
}

type AState struct {
	TransistionFunctions        []*TransistionFunction
	defaultTransistionFunctions *TransistionFunction
	enter                       *Function
	haltChan                    chan bool
	timeout                     time.Duration
	useTimeout                  bool

	IState
}

func (s *AState) GetTFunc() []*TransistionFunction {
	return s.TransistionFunctions
}

func (s *AState) Init(fx *Function) IState {
	s.TransistionFunctions = []*TransistionFunction{}
	s.haltChan = make(chan bool)
	s.enter = fx
	return s
}

func (s *AState) Destruct() error {
	err := s.Halt()
	if err != nil {
		return nil
	}
	close(s.haltChan)
	s.TransistionFunctions = nil
	return nil
}

func (s *AState) Halt() error {
	select {
	case <-time.After(5 * time.Second):
		return errors.New("halt timeout")
	case s.haltChan <- true:
		return nil
	}
}

func (s *AState) AddTransistionFunction(fx *TransistionFunction) {
	s.TransistionFunctions = append(s.TransistionFunctions, fx)
}

func (s *AState) SetFunction(fx *Function) {
	s.enter = fx
}

func (s *AState) SetTimeout(d time.Duration) {
	s.useTimeout = true
	s.timeout = d
}

func (s *AState) GetTimeout() time.Duration {
	return s.timeout
}

func (s *AState) UseTimeout() bool {
	return s.useTimeout
}

func (s *AState) Enter() error {
	errs := make(chan error)
	go func() {
		errs <- (*s.enter)(s)
	}()
	select {
	case halt := <-s.haltChan:
		if halt {
			return errors.New("state halted")
		}
	case err := <-errs:
		return err
	}
	return nil
}

func (s *AState) Transition() (IState, error) {
	for _, fx := range s.TransistionFunctions {
		if fx == nil {
			return nil, errors.New("no transistion function found")
		}
		res, err := (*fx)(s)
		if err != nil {
			return nil, err
		} else if res != nil {
			return res, nil
		}
	}
	if s.defaultTransistionFunctions != nil {
		res, err := (*s.defaultTransistionFunctions)(s)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, nil
}
