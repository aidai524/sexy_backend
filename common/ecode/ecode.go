package ecode

const (
	OK           = 0
	OKNew        = 200
	ServerErr    = -500
	RequestErr   = -400
	Unauthorized = -401
	Forbidden    = -403

	ExternalError = 2001

	UnknownError       = 100000
	ParamLimitMaxError = 100001
)

type CommonErrorCode int

const (
	CommonError CommonErrorCode = iota + 500000
)
