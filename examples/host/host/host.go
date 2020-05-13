package host

import (
	"github.com/filanov/stateswitch"
	"github.com/filanov/stateswitch/examples/host/hardware"
	"github.com/filanov/stateswitch/examples/host/models"
	"github.com/jinzhu/gorm"
)

type stateHost struct{ host *models.Host }

func (sh stateHost) State() stateswitch.State {
	return stateswitch.State(sh.host.Status)
}

func (sh *stateHost) SetState(state stateswitch.State) error {
	sh.host.Status = string(state)
	return nil
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
	return h.sm.Run(TransitionTypeRegister, &stateHost{host: host}, nil)
}

func (h *hostApi) SetHwInfo(host *models.Host, hw bool) error {
	return h.sm.Run(TransitionTypeSetHwInfo, &stateHost{host: host}, &TransitionArgsSetHwInfo{hwInfo: hw})
}

func (h *hostApi) List() ([]*models.Host, error) {
	var reply []*models.Host
	if err := h.db.Find(&reply).Error; err != nil {
		return nil, err
	}
	return reply, nil
}
