package common

import "github.com/zeebo/errs"

var (
	ServerError     = errs.Class("Server")
	ValidationError = errs.Class("Validation")
)
