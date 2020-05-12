package host

import "github.com/filanov/stateswitch/examples/host/models"

type API interface {
	List() ([]*models.Host, error)
	Register(*models.Host) error
	SetHwInfo(*models.Host, bool) error
}
