package common

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

var (
	ServerError     = errs.Class("Server")
	ValidationError = errs.Class("Validation")
	MarshalError    = errs.Class("MarshalJSON")
	UnmarshalError  = errs.Class("UnmarshalJSON")
)

func logServerError(err error) {
	logrus.WithError(err).Error("server error")
}

func logValidationError(err error) {
	logrus.WithError(err).Warn("validation error")
}

func logMarshalError(err error) {
	logrus.WithError(err).Warn("JSON decoding error")
}

func logUnmarshalError(err error) {
	logrus.WithError(err).Warn("JSON encoding error")
}

func HasServerError(err error) bool {
	return ServerError.Has(err)
}

func HasValidationError(err error) bool {
	return ValidationError.Has(err)
}

func HasUnmarshalError(err error) bool {
	return UnmarshalError.Has(err)
}

func HasMarshalError(err error) bool {
	return MarshalError.Has(err)
}

func HandleMarshalError(w http.ResponseWriter, err error) {
	logMarshalError(err)
	http.Error(w, "json encoding error", statusFromError(err))
}

func HandleUnmarshalError(w http.ResponseWriter, err error) {
	logUnmarshalError(err)
	http.Error(w, "json decoding error", statusFromError(err))
}

func HandleServerError(w http.ResponseWriter, err error) {
	logServerError(err)
	http.Error(w, err.Error(), statusFromError(err))
}

func HandleValidationError(w http.ResponseWriter, err error) {
	logValidationError(err)
	http.Error(w, err.Error(), statusFromError(err))
}

func HandleDefaultError(w http.ResponseWriter, err error) {
	logrus.Errorf("%+v", err)
	http.Error(w, "unknown error", statusFromError(err))
}

func HandleCitationError(w http.ResponseWriter, err error) {
	switch {
	case HasServerError(err):
		HandleServerError(w, err)

	case HasValidationError(err):
		HandleValidationError(w, err)

	default:
		HandleDefaultError(w, err)
	}
}

func statusFromError(err error) int {
	switch {
	case HasServerError(err):
		return http.StatusInternalServerError

	case HasValidationError(err):
		return http.StatusBadRequest

	case HasUnmarshalError(err):
		return http.StatusBadRequest

	case HasMarshalError(err):
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}
