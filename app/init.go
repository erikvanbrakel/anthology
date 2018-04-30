package app

import (
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/access"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Init(logger *logrus.Logger) routing.Handler {
	return func(rc *routing.Context) error {
		now := time.Now()
		rc.Response = &access.LogResponseWriter{rc.Response, http.StatusOK, 0}

		ac := newRequestScope(now, logger, rc.Request)

		rc.Set("Context", ac)

		fault.Recovery(ac.Errorf, convertError)(rc)
		logAccess(rc, ac.Infof, ac.Now())

		return nil
	}
}

// GetRequestScope returns the RequestScope of the current request.
func GetRequestScope(c *routing.Context) RequestScope {
	return c.Get("Context").(RequestScope)
}

func logAccess(c *routing.Context, logFunc access.LogFunc, start time.Time) {
	rw := c.Response.(*access.LogResponseWriter)
	elapsed := float64(time.Now().Sub(start).Nanoseconds()) / 1e6
	requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.URL.Path, c.Request.Proto)
	logFunc(`[%.3fms] %s %d %d`, elapsed, requestLine, rw.Status, rw.BytesWritten)
}

func convertError(c *routing.Context, err error) error {
	return err
}
