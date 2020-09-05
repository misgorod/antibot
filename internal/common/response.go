package common

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

/*type Logger interface {
	Debug(args ...interface{})
	Tracef(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
}*/

func RespondOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func RespondBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
}

func RespondInternalServerError(log *logrus.Logger, w http.ResponseWriter, err error) {
	log.WithError(err).Error("internal error")
	w.WriteHeader(http.StatusInternalServerError)
}
