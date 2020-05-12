package hardware

type Validator interface {
	IsSufficient(bool) bool
}

type validator struct{}

func New() *validator {
	return &validator{}
}

func (v *validator) IsSufficient(hw bool) bool {
	return hw
}
