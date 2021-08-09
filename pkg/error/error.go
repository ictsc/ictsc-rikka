package error

type BadRequestError struct {
	Message string
}

type UnauthorizedError struct {
	Message string
}

type ForbiddenError struct {
	Message string
}

type NotFoundError struct {
	Message string
}

type InternalServerError struct {
	Err error
}

var _ error = (*BadRequestError)(nil)
var _ error = (*UnauthorizedError)(nil)
var _ error = (*ForbiddenError)(nil)
var _ error = (*NotFoundError)(nil)
var _ error = (*InternalServerError)(nil)

func NewBadRequestError(msg string) *BadRequestError {
	return &BadRequestError{
		Message: msg,
	}
}

func NewUnauthorizedError(msg string) *UnauthorizedError {
	return &UnauthorizedError{
		Message: msg,
	}
}

func NewForbiddenError(msg string) *ForbiddenError {
	return &ForbiddenError{
		Message: msg,
	}
}

func NewNotFoundError(msg string) *NotFoundError {
	return &NotFoundError{
		Message: msg,
	}
}

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{
		Err: err,
	}
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (V *InternalServerError) Error() string {
	return "internal server error"
}
