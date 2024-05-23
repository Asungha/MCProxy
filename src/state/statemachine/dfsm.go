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
	Fx func() bool
}

type IStateMachine interface {
	RegisterState(state.IState) error
	State(State) state.IState
	TransistionCondition(pair TransistionPair, condition func() bool)
	Construct() error
	Destruct() error
	SetRoot(State)
	Run() error
}

type AStateMachine struct {
	currectState state.IState
	States       map[string]state.IState
	Conditions   []Condition
	Ctx          context.Context
	Cancle       context.CancelCauseFunc
	IStateMachine
}

func (sm *AStateMachine) RegisterState(stateName State, s state.IState) error {
	if sm.States == nil {
		sm.States = make(map[string]state.IState)
		sm.States[""] = nil
	}
	// log.Printf(string(stateName))
	sm.States[string(stateName)] = s
	// log.Printf("%v", sm.States)
	return nil
}

func (sm *AStateMachine) TransistionCondition(pair TransistionPair, condition state.ConditionFunction) {
	// c := condition
	sm.Conditions = append(sm.Conditions, Condition{Fx: condition, TransistionPair: pair})
}

func (sm *AStateMachine) TransistionFunction(source State, fx state.TransistionFunction) {
	(sm.States[string(source)]).AddTransistionFunction(fx)
	// log.Printf("TFunc %s %v", source, sm.States[string(source)].GetTFunc())
}

func (sm *AStateMachine) Construct() error {
	for _, condition := range sm.Conditions {
		c := condition
		var fx state.TransistionFunction = func(i state.IState) (state.IState, error) {
			if c.Fx() {
				return sm.States[c.Destination], nil
			}
			return nil, nil
		}
		(sm.States[c.Source]).AddTransistionFunction(fx)
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

func (sm *AStateMachine) State(s State) state.IState {
	return sm.States[string(s)]
}

func (sm *AStateMachine) Run() error {
	defer log.Println("[statemachine worker] Thread exit")
	log.Println("[statemachine worker] start")
	// defer func() {
	// 	log.Printf("[State machine] Exit")
	// 	// err := sm.Halt()
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// 	sm.Cancle(nil)
	// }()
	// c := make(chan error)
	// defer close(c)
	for {
		// var ctx context.Context
		// var cancle context.CancelFunc
		// if sm.currectState.UseTimeout() {
		// 	ctx, cancle = context.WithDeadline(context.Background(), time.Now().Add(sm.currectState.GetTimeout()))
		// } else {
		// 	ctx, cancle = context.WithCancel(context.Background())
		// }
		// log.Printf("run start")
		// go func() {
		// 	// defer log.Println("[state worker] Thread exit")
		// 	c <- (sm.currectState).Enter()
		// }()
		// select {
		// case err := <-c:
		// 	if err != nil {
		// 		log.Printf("run err %v", err.Error())
		// 		// log.Println("e")
		// 		return err
		// 		// panic(err)
		// 	}
		// 	log.Printf("run ok")
		// case <-sm.Ctx.Done():
		// 	// log.Printf("Halted")
		// 	log.Printf("run halt")
		// 	return errors.New("Halted")
		// }
		err := (sm.currectState).Enter()
		if err != nil {
			log.Printf("run err %v", err.Error())
			// log.Println("e")
			return err
			// panic(err)
		}
		// err := (sm.currectState).Enter()

		nextState, err := (sm.currectState).Transition()
		if err != nil {
			// log.Printf("%v %v", nextState, err)
			// log.Println("t")
			// log.Println(err)
			return err
			// panic(err)
		}
		if nextState == nil {
			// log.Println("done")
			return nil
		}
		sm.currectState = nextState
	}
}
