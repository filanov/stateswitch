package host

import (
	"fmt"
	"github.com/filanov/stateswitch"
	"github.com/filanov/stateswitch/examples/host/hardware"
	"github.com/filanov/stateswitch/examples/host/models"
	"github.com/jinzhu/gorm"
)

type internalHost struct{ host *models.Host }

func (i internalHost) State() stateswitch.State {
	return stateswitch.State(i.host.Status)
}

func (i *internalHost) SetState(state stateswitch.State) error {
	i.host.Status = string(state)
	return nil
}

func (s *internalHost) RunCondition(ifn stateswitch.Condition, args stateswitch.TransitionArgs) (bool, error) {
	fn, ok := ifn.(func (*internalHost, stateswitch.TransitionArgs) (bool, error))
	if !ok {
		return false, fmt.Errorf("Condition function type is not applicable ...")
	}
	return fn(s, args)
}

func (s *internalHost) RunTransition(ifn stateswitch.Transition, args stateswitch.TransitionArgs) error {
	fn, ok := ifn.(func(*internalHost, stateswitch.TransitionArgs) error)
	if !ok {
		return fmt.Errorf("Transition function type is not applicable ...")
	}
	return fn(s, args)
}

type hostApi struct {
	sm stateswitch.StateMachine
	db *gorm.DB
}

const (
	StateDiscovering  = "discovering"
	StateKnown        = "known"
	StateInsufficient = "insufficient"
)

const (
	TransitionTypeSetHwInfo = "SetHwInfo"
	TransitionTypeRegister  = "Register"
)

func NewHostStateMachine(th *transitionHandler) stateswitch.StateMachine {
	sm := stateswitch.NewStateMachine()

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetHwInfo,
		SourceStates:     stateswitch.States{StateDiscovering, StateKnown, StateInsufficient},
		DestinationState: StateKnown,
		Condition:        th.IsSufficient,
		Transition:       th.SetHwInfo,
		PostTransition:   th.PostSetHwInfo,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeSetHwInfo,
		SourceStates:     stateswitch.States{StateDiscovering, StateKnown, StateInsufficient},
		DestinationState: StateInsufficient,
		Condition:        th.IsInsufficient,
		Transition:       th.SetHwInfo,
		PostTransition:   th.PostSetHwInfo,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeRegister,
		SourceStates:     stateswitch.States{""},
		DestinationState: StateDiscovering,
		Condition:        nil,
		Transition:       nil,
		PostTransition:   th.RegisterNew,
	})

	sm.AddTransition(stateswitch.TransitionRule{
		TransitionType:   TransitionTypeRegister,
		SourceStates:     stateswitch.States{StateDiscovering, StateKnown, StateInsufficient},
		DestinationState: StateDiscovering,
		Condition:        nil,
		Transition:       nil,
		PostTransition:   th.RegisterAgain,
	})

	return sm
}

func New(db *gorm.DB, hwValidator hardware.Validator) API {
	th := &transitionHandler{
		db:          db,
		hwValidator: hwValidator,
	}
	return &hostApi{
		db: db,
		sm: NewHostStateMachine(th),
	}
}

func (h *hostApi) Register(host *models.Host) error {
	return h.sm.Run(TransitionTypeRegister, &internalHost{host}, nil)
}

func (h *hostApi) SetHwInfo(host *models.Host, hw bool) error {
	return h.sm.Run(TransitionTypeSetHwInfo, &internalHost{host}, hw)
}

func (h *hostApi) List() ([]*models.Host, error) {
	var reply []*models.Host
	if err := h.db.Find(&reply).Error; err != nil {
		return nil, err
	}
	return reply, nil
}
