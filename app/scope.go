package app

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type RequestScope interface {
	Logger
	Now() time.Time
}

type requestScope struct {
	Logger
	now       time.Time
	requestID string
}

func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request) RequestScope {
	log := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		log.SetField("RequestID", requestID)
	}
	return &requestScope{
		Logger:    log,
		now:       now,
		requestID: requestID,
	}
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}
