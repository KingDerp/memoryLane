package webutil

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

var (
	ServerError     = errs.Class("Server")
	ValidationError = errs.Class("Validation")
)

func logServerError(err error) {
	logrus.WithError(err).Error("server error")
}

func logValidationError(err error) {
	logrus.WithError(err).Warn("validation error")
}

func HasServerError(err error) bool {
	return ServerError.Has(err)
}

func HasValidationError(err error) bool {
	return ValidationError.Has(err)
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

func HandleError(w http.ResponseWriter, err error) {
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

	default:
		return http.StatusInternalServerError
	}
}

func RenderJSON(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		HandleError(w, err)
		return
	}

	//set header
	h := w.Header()
	h.Set("Content-Type", "application/json")
	w.Write(b)
}
