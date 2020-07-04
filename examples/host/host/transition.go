package host

import (
	"github.com/filanov/stateswitch"
	"github.com/filanov/stateswitch/examples/host/hardware"
	"github.com/filanov/stateswitch/examples/host/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type transitionHandler struct {
	db          *gorm.DB
	hwValidator hardware.Validator
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// SetHwInfo transition
///////////////////////////////////////////////////////////////////////////////////////////////////

type TransitionArgsSetHwInfo struct {
	hwInfo bool
}

func (th *transitionHandler) SetHwInfo(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.Errorf("incompatible type of StateSwitch")
	}
	params, ok := args.(*TransitionArgsSetHwInfo)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	sHost.host.HwInfo = &params.hwInfo
	return nil
}

func (th *transitionHandler) IsSufficient(_ stateswitch.StateSwitch, args stateswitch.TransitionArgs) (bool, error) {
	params, ok := args.(*TransitionArgsSetHwInfo)
	if !ok {
		return false, errors.Errorf("invalid argument")
	}
	return th.hwValidator.IsSufficient(params.hwInfo), nil
}

func (th *transitionHandler) IsConnected(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) (bool, error) {
	// Always connected
	return true, nil
}

func (th *transitionHandler) PostSetHwInfo(sw stateswitch.StateSwitch, _ stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.Errorf("incompatible type of StateSwitch")
	}
	updates := map[string]interface{}{"status": sHost.host.Status, "hw_info": sHost.host.HwInfo}
	return th.db.Model(sHost.host).Where("id = ?", sHost.host.ID).Updates(updates).Error
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Register transition
///////////////////////////////////////////////////////////////////////////////////////////////////

type TransitionArgsRegister struct {
	host *models.Host
}

func (th *transitionHandler) RegisterNew(sw stateswitch.StateSwitch, args stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.Errorf("incompatible type of StateSwitch")
	}
	return th.db.Create(sHost.host).Error
}

func (th *transitionHandler) RegisterAgain(sw stateswitch.StateSwitch, _ stateswitch.TransitionArgs) error {
	sHost, ok := sw.(*stateHost)
	if !ok {
		return errors.Errorf("incompatible type of StateSwitch")
	}
	updates := map[string]interface{}{"status": sHost.host.Status, "hw_info": nil}
	return th.db.Model(sHost.host).Where("id = ?", sHost.host.ID).Updates(updates).Error
}
