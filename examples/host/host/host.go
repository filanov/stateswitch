package host

import (
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

type hostApi struct {
	sm          stateswitch.StateMachine
	db          *gorm.DB
	hwValidator hardware.Validator
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
		db:          db,
		hwValidator: hwValidator,
		sm:          NewHostStateMachine(th),
	}
}

func (h *hostApi) Register(host *models.Host) error {
	return h.sm.Run(TransitionTypeRegister, &internalHost{host}, &TransitionArgsRegister{
		host: host,
	})
}

func (h *hostApi) SetHwInfo(host *models.Host, hw bool) error {
	return h.sm.Run(TransitionTypeSetHwInfo, &internalHost{host}, &TransitionArgsSetHwInfo{
		host:   host,
		hwInfo: hw,
	})
	//updates := map[string]interface{}{"status": host.Status, "hw_info": host.HwInfo}
	//return h.db.Model(host).Where("id = ?", host.ID).Updates(updates).Error
}

func (h *hostApi) List() ([]*models.Host, error) {
	var reply []*models.Host
	if err := h.db.Find(&reply).Error; err != nil {
		return nil, err
	}
	return reply, nil
}
