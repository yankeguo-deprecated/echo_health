package echo_health

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"sync/atomic"
)

const (
	PathLiveness  = "/health/alive"
	PathReadiness = "/health/ready"
)

const (
	DefaultReadinessCascadeCount = 5
)

func runChecks(cks []Check, fn func(ck Check) error) string {
	sb := &strings.Builder{}
	for _, ck := range cks {
		if err := fn(ck); err != nil {
			if sb.Len() > 0 {
				sb.WriteString("; ")
			}
			sb.WriteString(err.Error())
		}
	}
	return sb.String()
}

func New(cks ...Check) echo.MiddlewareFunc {
	var readinessFailed int64
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().URL.Path == PathLiveness {
				// cascading readiness failure
				if readinessFailed > DefaultReadinessCascadeCount {
					return c.String(http.StatusServiceUnavailable, "cascading readiness failure")
				}
				// run checks
				res := runChecks(cks, func(ck Check) error {
					return ck.CheckLiveness()
				})
				if len(res) > 0 {
					return c.String(http.StatusServiceUnavailable, res)
				}
				// default to OK
				return c.String(http.StatusOK, "OK")
			} else if c.Request().URL.Path == PathReadiness {
				// run checks
				res := runChecks(cks, func(ck Check) error {
					return ck.CheckReadiness()
				})
				if len(res) > 0 {
					atomic.AddInt64(&readinessFailed, 1)
					return c.String(http.StatusServiceUnavailable, res)
				}
				atomic.StoreInt64(&readinessFailed, 0)
				// default to OK
				return c.String(http.StatusOK, "OK")
			}
			return next(c)
		}
	}
}
