package transitions

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// DO NOT EDIT!
// This code is generated with http://github.com/MrEhbr/gofsm tool

type (
	// OrderTransition is a state transition and all data are literal values that simplifies FSM usage and make it generic.
	OrderTransition struct {
		Event         string
		From          StateType
		To            StateType
		BeforeActions []string
		Actions       []string
	}
	// OrderHandle handles transitions action
	OrderHandleAction func(action string, fromState, toState StateType, obj *Order) error
	// Save state to external storage
	OrderPersistState func(obj *Order, state StateType) error
	// OrderStateMachine is a FSM that can handle transitions of a lot of objects. eventHandler and transitions are configured before use them.
	OrderStateMachine struct {
		transitions   []OrderTransition
		actionHandler OrderHandleAction
		persister     OrderPersistState
	}
)

var (
	ErrOrderFsmAction       = errors.New("OrderStateMachine action error")
	ErrOrderFsmBeforeAction = errors.New("OrderStateMachine before action error")
	// ErrOrderFsmSkip indicates that further processing not need
	// used in before_actions
	ErrOrderFsmSkip = errors.New("skip")
)

type Option func(*OrderStateMachine)

func WithActionHandler(h OrderHandleAction) Option {
	return func(fsm *OrderStateMachine) {
		fsm.actionHandler = h
	}
}

func WithPersiter(p OrderPersistState) Option {
	return func(fsm *OrderStateMachine) {
		fsm.persister = p
	}
}

func WithTransitions(tr []OrderTransition) Option {
	return func(fsm *OrderStateMachine) {
		fsm.transitions = tr
	}
}

// NewOrderStateMachine creates a new state machine.
func NewOrderStateMachine(opts ...Option) *OrderStateMachine {
	fsm := &OrderStateMachine{}
	for _, o := range opts {
		o(fsm)
	}

	return fsm
}

// ChangeState fires a event and if event succeeded then change state.
func (m *OrderStateMachine) ChangeState(event string, obj *Order) error {
	trans, ok := m.findTransMatching(obj.State, event)
	if !ok {
		return fmt.Errorf("cannot find transition for event [%s] when in state [%v]", event, obj.State)
	}

	if len(trans.BeforeActions) > 0 && m.actionHandler != nil {
		for _, action := range trans.BeforeActions {
			if err := m.actionHandler(action, trans.From, trans.To, obj); err != nil {
				if errors.Is(err, ErrOrderFsmSkip) {
					return nil
				}

				return fmt.Errorf("%w. action [%s] return error: %s", ErrOrderFsmBeforeAction, action, err)
			}
		}
	}

	if m.persister != nil {
		if err := m.persister(obj, trans.To); err != nil {
			return err
		}
	}

	obj.State = trans.To

	if len(trans.Actions) > 0 && m.actionHandler != nil {
		var errs error
		for _, action := range trans.Actions {
			if err := m.actionHandler(action, trans.From, trans.To, obj); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("%w. action [%s] return error: %s", ErrOrderFsmAction, action, err))
			}
		}

		if errs != nil {
			return errs
		}
	}

	return nil
}

func (m *OrderStateMachine) Can(state StateType, event string) bool {
	_, ok := m.findTransMatching(state, event)
	return ok
}

func (m *OrderStateMachine) FindTransitionForStates(from, to StateType) (OrderTransition, bool) {
	for _, v := range m.transitions {
		if v.From == from && v.To == to {
			return v, true
		}
	}
	return OrderTransition{}, false
}

// findTransMatching gets corresponding transition according to current state and event.
func (m *OrderStateMachine) findTransMatching(fromState StateType, event string) (OrderTransition, bool) {
	for _, v := range m.transitions {
		if v.From == fromState && v.Event == event {
			return v, true
		}
	}
	return OrderTransition{}, false
}

const (
	// Order state machine events
	OrderEventPlaceOrder   = "place_order"
	OrderEventFailOrder    = "fail_order"
	OrderEventSuccessOrder = "success_order"

	// Order state machine actions
	OrderActionBook              = "book"
	OrderActionCheckAvailability = "check_availability"
	OrderActionSendEmail         = "send_email"
)

// OrderTransitions generated from transitions.json
var OrderTransitions = []OrderTransition{
	{
		From: StateTypeUnknown,
		To:   Created,
	},
	{
		Event: OrderEventPlaceOrder,
		From:  Created,
		To:    Started,
		BeforeActions: []string{
			OrderActionCheckAvailability,
			OrderActionBook,
		},
	},
	{
		Event: OrderEventFailOrder,
		From:  Created,
		To:    Failed,
	},
	{
		Event: OrderEventSuccessOrder,
		From:  Started,
		To:    Finished,
		Actions: []string{
			OrderActionSendEmail,
		},
	},
}
