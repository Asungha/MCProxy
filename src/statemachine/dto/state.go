package state

import (
	"errors"
	"time"
)

type TransistionFunction func(IState) (IState, error)

func NewTransistionFunction(fx TransistionFunction) *TransistionFunction {
	return &fx
}

type Function func(IState) error
type ActionFunction func(IState) error
type ConditionFunction func() bool
type DeferFunction func()

func True() bool {
	return true
}

// type StateID string

type StateMetadata map[string]interface{}

// func CastMetadata[T any](data interface{}) *T {
// 	casted, ok := data.(T)
// 	if !ok {
// 		return nil
// 	}
// 	return &casted
// }

type IState interface {
	Init(Function) IState

	Enter() error
	Transition() (IState, error)

	AddTransistionFunction(TransistionFunction)

	SetTimeout(time.Duration)
	GetTimeout() time.Duration
	UseTimeout() bool

	Destruct() error
	GetTFunc() []TransistionFunction
}

type State struct {
	TransistionFunctions []TransistionFunction
	enter                Function
	timeout              time.Duration
	useTimeout           bool
}

func NewState(fx Function) *State {
	s := &State{}
	s.Init(fx)
	return s
}

func (s *State) GetTFunc() []TransistionFunction {
	return s.TransistionFunctions
}

func (s *State) Init(fx Function) IState {
	s.TransistionFunctions = []TransistionFunction{}
	s.enter = fx
	return s
}

func (s *State) Destruct() error {
	s.TransistionFunctions = nil
	return nil
}

func (s *State) AddTransistionFunction(fx TransistionFunction) {
	s.TransistionFunctions = append(s.TransistionFunctions, fx)
}

func (s *State) SetFunction(fx Function) {
	s.enter = fx
}

func (s *State) SetTimeout(d time.Duration) {
	s.useTimeout = true
	s.timeout = d
}

func (s *State) GetTimeout() time.Duration {
	return s.timeout
}

func (s *State) UseTimeout() bool {
	return s.useTimeout
}

func (s *State) Enter() error {
	return s.enter(s)
}

func (s *State) Transition() (IState, error) {
	for _, fx := range s.TransistionFunctions {
		if fx == nil {
			return nil, errors.New("no transistion function found")
		}
		res, err := fx(s)
		if err != nil {
			return nil, err
		} else if res != nil {
			return res, nil
		}
	}
	return nil, nil
}
