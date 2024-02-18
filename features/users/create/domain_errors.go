package create

type ValidationError struct {
	Msg string
}

func (e *ValidationError) Code() string {
	return "validation_error"
}

func (e *ValidationError) Error() string {
	return e.Msg
}

func NewValidationError(msg string) error {
	return &ValidationError{Msg: msg}
}
