package statemachine

import (
	"context"
	"log"
	state "mc_reverse_proxy/src/state/state"
)

type State string
type TransistionPair struct {
	Source      string
	Destination string
}
type Condition struct {
	TransistionPair
	Fx *func() bool
}

type IStateMachine interface {
	RegisterState(state.IState) error
	TransistionCondition(pair TransistionPair, condition func() bool)
	Construct() error
	Destruct() error
	SetRoot(State)
	Run() error
	Halt() error
}

type AStateMachine struct {
	currectState state.IState
	States       map[string]state.IState
	Conditions   []*Condition
	Ctx          context.Context
	Cancle       context.CancelCauseFunc
	IStateMachine
}

func (sm *AStateMachine) RegisterState(stateName State, s state.IState) error {
	if sm.States == nil {
		sm.States = make(map[string]state.IState)
		sm.States[""] = nil
	}
	log.Printf(string(stateName))
	sm.States[string(stateName)] = s
	log.Printf("%v", sm.States)
	return nil
}

func (sm *AStateMachine) TransistionCondition(pair TransistionPair, condition *func() bool) {
	// c := condition
	sm.Conditions = append(sm.Conditions, &Condition{Fx: condition, TransistionPair: pair})
}

func (sm *AStateMachine) Construct() error {
	for _, condition := range sm.Conditions {
		c := condition
		var fx state.TransistionFunction = func(i state.IState) (state.IState, error) {
			if (*c.Fx)() {
				return sm.States[c.Destination], nil
			}
			return nil, nil
		}
		(sm.States[c.Source]).AddTransistionFunction(&fx)
		log.Printf("TFunc %s %v", c.Source, sm.States[c.Source].GetTFunc())
	}
	return nil
}

func (sm *AStateMachine) Destruct() error {
	for _, state := range sm.States {
		(state).Destruct()
	}
	sm.Conditions = nil
	sm.States = nil
	return nil
}

func (sm *AStateMachine) SetRoot(state State) {
	sm.currectState = sm.States[string(state)]
}

func (sm *AStateMachine) Run() error {
	for {
		// var ctx context.Context
		// var cancle context.CancelFunc
		// if sm.currectState.UseTimeout() {
		// 	ctx, cancle = context.WithDeadline(context.Background(), time.Now().Add(sm.currectState.GetTimeout()))
		// } else {
		// 	ctx, cancle = context.WithCancel(context.Background())
		// }
		err := (sm.currectState).Enter()
		if err != nil {
			log.Printf("%v", err.Error())
			log.Println("e")
			return err
			// panic(err)
		}
		nextState, err := (sm.currectState).Transition()
		if err != nil {
			// log.Printf("%v %v", nextState, err)
			log.Println("t")
			log.Println(err)
			return err
			// panic(err)
		}
		if nextState == nil {
			log.Println("done")
			return nil
		}
		sm.currectState = nextState
	}
}
