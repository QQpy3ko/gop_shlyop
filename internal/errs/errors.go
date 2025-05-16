package errs

type ErrBadRequest struct {
	Message string
}

func (e *ErrBadRequest) Error() string {
	if e.Message == "" {
		return "bad request"
	}
	return e.Message
}

func NewErrBadRequest(message string) error {
	return &ErrBadRequest{Message: message}
}

type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	if e.Message == "" {
		return "not found"
	}
	return e.Message
}

func NewErrNotFound(message string) error {
	return &ErrNotFound{Message: message}
}
