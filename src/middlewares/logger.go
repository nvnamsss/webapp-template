package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
	"webapp-template/src/logger"

	"github.com/gin-gonic/gin"
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
// It receives:
//  1. A time package format string (e.g. time.RFC3339).
//  2. A boolean stating whether to use UTC time zone or local.
func GinLogger(utc bool, opts ...OpOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx context.Context

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		ctx = logger.SetRqIDToCtx(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		var bodyBytes []byte
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
			if c.Request.Body != nil {
				bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			}
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		rqID := logger.GetRqIDFromCtx(c.Request.Context())

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Errorf(e)
			}
		} else {

			fields := logger.Fields{
				"status":     c.Writer.Status(),
				"method":     c.Request.Method,
				"path":       path,
				"query":      query,
				"ip":         c.ClientIP(),
				"user-agent": c.Request.UserAgent(),
				"time":       end.Unix(),
				"exec_time":  float64(latency) / 1000000,
				"rqID":       rqID,
				// "user_id":           ctxutil.GetUserIDStr(c.Request.Context()),
				"device_id":         c.Request.Header.Get("device_id"),
				"device_session_id": c.Request.Header.Get("device_session_id"),
				"x-forwarded-for":   c.Request.Header.Get("X-Forwarded-For"),
				"cf-connecting-ip":  c.Request.Header.Get("CF-Connecting-IP"),
			}

			if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
				buffer := new(bytes.Buffer)
				_ = json.Compact(buffer, bodyBytes)
				fields["data"] = buffer
			}

			if c.Writer.Status() != http.StatusOK {
				// fields["error_code"] = ctxutil.GetErrorCode(c.Request.Context())
			}

			if len(opts) > 0 {
				ret := LogOption{}
				ret.applyOpts(opts)

				switch {
				case len(ret.key) > 0:
					// for _, key := range ret.key {
					// 	fields[key] = ctxutil.GetCtxValueString(c.Request.Context(), key)
					// }
				}
			}

			logger.WithFields(fields).Infof("")
		}

	}
}

// RecoveryWithZap returns a gin.HandlerFunc (middleware)
// that recovers from any panics and logs requests using uber-go/zap.
// All errors are logged using zap.Error().
// stack means whether output the stack info.
// The stack info is easy to find where the error occurs but the stack info is too large.
func RecoveryWithLogger(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.WithFields(logger.Fields{
						"error":   err,
						"request": string(httpRequest),
					}).Infof("Recovery from panic")
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.WithFields(logger.Fields{
						"time":    time.Now().Unix(),
						"error":   err,
						"request": string(httpRequest),
						"stack":   string(debug.Stack()),
					}).Errorf("Recovery from panic")
				} else {
					logger.WithFields(logger.Fields{
						"time":    time.Now().Unix(),
						"error":   err,
						"request": string(httpRequest),
					}).Errorf("Recovery from panic")
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

type LogOption struct {
	key []string
}

// FieldOption configures logging options
type OpOption func(*LogOption)

// WithKey attaches a key added into log from context
func WithKey(key []string) OpOption {
	return func(op *LogOption) { op.key = key }
}

func (op *LogOption) applyOpts(opts []OpOption) {
	for _, opt := range opts {
		opt(op)
	}
}
