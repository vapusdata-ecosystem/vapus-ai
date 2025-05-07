package dmerrors

import "errors"

type Error struct {
	libErr error
	appErr error
}

func (e Error) LibError() error {
	return e.libErr
}

func (e Error) AppError() error {
	return e.appErr
}

func DMError(apperror, liberr error) error {
	return Error{
		libErr: liberr,
		appErr: apperror,
	}
}

func (e Error) Error() string {
	if e.libErr == nil && e.appErr == nil {
		return ""
	}
	if e.libErr == nil {
		return e.appErr.Error()
	}
	if e.appErr == nil {
		return e.libErr.Error()
	}
	return errors.Join(e.libErr, e.appErr).Error()
}

// func DMErrorChain(err)
