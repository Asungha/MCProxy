package statemachine

import (
	"context"
	"errors"
	state "mc_reverse_proxy/src/statemachine/dto"
	statemachineDTO "mc_reverse_proxy/src/statemachine/dto"
)

// type State string
type TransistionPair struct {
	Source      string
	Destination string
}
type Condition struct {
	TransistionPair
	Fx func() bool
}

type IStateMachine interface {
	RegisterState(statemachineDTO.IState) error
	State(statemachineDTO.State) statemachineDTO.IState
	TransistionCondition(pair TransistionPair, condition func() bool)
	Construct() error
	Destruct() error
	SetRoot(statemachineDTO.State)
	Run() error
}

type StateMachine struct {
	currectState statemachineDTO.IState

	States     map[string]statemachineDTO.IState
	Conditions []Condition
	Ctx        context.Context
	Cancle     context.CancelCauseFunc
	DeferFunc  statemachineDTO.DeferFunction
}

func (sm *StateMachine) RegisterState(stateName string, s statemachineDTO.IState) error {
	if sm.States == nil {
		sm.States = make(map[string]statemachineDTO.IState)
		sm.States[""] = nil
	}
	sm.States[string(stateName)] = s
	return nil
}

func (sm *StateMachine) TransistionCondition(pair TransistionPair, condition statemachineDTO.ConditionFunction) {
	sm.Conditions = append(sm.Conditions, Condition{Fx: condition, TransistionPair: pair})
}

func (sm *StateMachine) TransistionFunction(source string, fx statemachineDTO.TransistionFunction) {
	(sm.States[string(source)]).AddTransistionFunction(fx)
}

func (sm *StateMachine) Construct() error {
	for _, condition := range sm.Conditions {
		c := condition
		var fx statemachineDTO.TransistionFunction = func(i statemachineDTO.IState) (statemachineDTO.IState, error) {
			if c.Fx() {
				return sm.States[c.Destination], nil
			}
			return nil, nil
		}
		(sm.States[c.Source]).AddTransistionFunction(fx)
	}
	return nil
}

func (sm *StateMachine) Destruct() error {
	sm.Cancle(errors.New("Desturct Called"))
	for _, state := range sm.States {
		if state == nil {
			continue
		}
		(state).Destruct()
	}
	if sm.DeferFunc != nil {
		sm.DeferFunc()
	}
	sm.Conditions = nil
	sm.States = nil
	return nil
}

func (sm *StateMachine) SetRoot(state string) {
	sm.currectState = sm.States[string(state)]
}

func (sm *StateMachine) State(s string) state.IState {
	return sm.States[string(s)]
}

func (sm *StateMachine) Run() error {
	defer sm.Destruct()
	defer sm.DeferFunc()
	for {
		err := (sm.currectState).Enter()
		if err != nil {
			return err
		}

		nextState, err := (sm.currectState).Transition()
		if err != nil {
			return err
		}
		if nextState == nil {
			return nil
		}
		sm.currectState = nextState
	}
}
