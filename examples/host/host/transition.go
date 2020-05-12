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

func (th *transitionHandler) SetHwInfo(i *internalHost, args stateswitch.TransitionArgs) error {
	hwInfo, ok := args.(bool)
	if !ok {
		return errors.Errorf("invalid argument")
	}
	i.host.HwInfo = &hwInfo
	return nil
}

func (th *transitionHandler) IsSufficient(i *internalHost, args stateswitch.TransitionArgs) (bool, error) {
	hwInfo, ok := args.(bool)
	if !ok {
		return false, errors.Errorf("invalid argument")
	}
	return th.hwValidator.IsSufficient(hwInfo), nil
}

func (th *transitionHandler) IsInsufficient(i *internalHost, args stateswitch.TransitionArgs) (bool, error) {
	reply, err := th.IsSufficient(i, args)
	return !reply, err
}

func (th *transitionHandler) PostSetHwInfo(i *internalHost, args stateswitch.TransitionArgs) error {
	updates := map[string]interface{}{"status": i.host.Status, "hw_info": i.host.HwInfo}
	return th.db.Model(i.host).Where("id = ?", i.host.ID).Updates(updates).Error
}

///////////////////////////////////////////////////////////////////////////////////////////////////
// Register transition
///////////////////////////////////////////////////////////////////////////////////////////////////

type TransitionArgsRegister struct {
	host *models.Host
}

func (th *transitionHandler) RegisterNew(i *internalHost, args stateswitch.TransitionArgs) error {
	return th.db.Create(i.host).Error
}

func (th *transitionHandler) RegisterAgain(i *internalHost, args stateswitch.TransitionArgs) error {
	updates := map[string]interface{}{"status": i.host.Status, "hw_info": nil}
	return th.db.Model(i.host).Where("id = ?", i.host.ID).Updates(updates).Error
}
