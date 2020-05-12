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
	host   *models.Host
	hwInfo bool
}

func (th *transitionHandler) SetHwInfo(args stateswitch.TransitionArgs) error {
	params, ok := args.(*TransitionArgsSetHwInfo)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	params.host.HwInfo = &params.hwInfo
	return nil
}

func (th *transitionHandler) IsSufficient(args stateswitch.TransitionArgs) (bool, error) {
	params, ok := args.(*TransitionArgsSetHwInfo)
	if !ok {
		return false, errors.Errorf("invalid argument")
	}
	return th.hwValidator.IsSufficient(params.hwInfo), nil
}

func (th *transitionHandler) IsInsufficient(args stateswitch.TransitionArgs) (bool, error) {
	reply, err := th.IsSufficient(args)
	return !reply, err
}

func (th *transitionHandler) PostSetHwInfo(args stateswitch.TransitionArgs) error {
	params, ok := args.(*TransitionArgsSetHwInfo)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	updates := map[string]interface{}{"status": params.host.Status, "hw_info": params.host.HwInfo}
	return th.db.Model(params.host).Where("id = ?", params.host.ID).Updates(updates).Error
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Register transition
///////////////////////////////////////////////////////////////////////////////////////////////////

type TransitionArgsRegister struct {
	host *models.Host
}

func (th *transitionHandler) RegisterNew(args stateswitch.TransitionArgs) error {
	params, ok := args.(*TransitionArgsRegister)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	return th.db.Create(params.host).Error
}

func (th *transitionHandler) RegisterAgain(args stateswitch.TransitionArgs) error {
	params, ok := args.(*TransitionArgsRegister)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	updates := map[string]interface{}{"status": params.host.Status, "hw_info": nil}
	return th.db.Model(params.host).Where("id = ?", params.host.ID).Updates(updates).Error
}
