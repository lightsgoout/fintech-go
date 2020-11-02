package service

import "errors"

var (
	ErrAccountDoesNotExist  = errors.New("account does not exist")
	ErrIncompatibleCurrency = errors.New("incompatible currency")
	ErrBadAccountID         = errors.New("bad account id")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrAccountAlreadyExists = errors.New("account already exists")
)

type ErrInternal struct {
	err error
}

func NewErrInternal(err error) ErrInternal {
	return ErrInternal{
		err: err,
	}
}

func (e ErrInternal) Error() string {
	return "internal error: " + e.err.Error()
}

func (e ErrInternal) Unwrap() error {
	return e.err
}
