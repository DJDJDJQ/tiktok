package pkg

import (
	"errors"
	"fmt"
)

const (
	SuccessCode               = 0
	ServiceErrCode            = 10001
	TokenInvalidCode          = 10002
	ParamErrCode              = 10003
	LoginErrCode              = 10004
	RecordNotExistErrCode     = 10005
	RecordAlreadyExistErrCode = 10006
)

type ErrNo struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.StatusCode, e.StatusMsg)
}

func NewErrNo(code int32, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.StatusMsg = msg
	return e
}

var (
	Success               = NewErrNo(SuccessCode, "")
	ServiceErr            = NewErrNo(ServiceErrCode, "Service is unable to start successfully")
	TokenInvalidErr       = NewErrNo(TokenInvalidCode, "Token is invalid")
	ParamErr              = NewErrNo(ParamErrCode, "Wrong Parameter has been given")
	LoginErr              = NewErrNo(LoginErrCode, "Wrong username or password")
	RecordNotExistErr     = NewErrNo(RecordNotExistErrCode, "The record does not exists")
	RecordAlreadyExistErr = NewErrNo(RecordAlreadyExistErrCode, "The record already exists")
)

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.StatusMsg = err.Error()
	return s
}
