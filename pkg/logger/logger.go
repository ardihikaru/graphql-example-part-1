// Package logger provides functions to set up a new logger
package logger

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/ardihikaru/graphql-example-part-1/pkg/config"
	enum "github.com/ardihikaru/graphql-example-part-1/pkg/enum/channel"
	"github.com/ardihikaru/graphql-example-part-1/pkg/enum/loglevel"
	mw "github.com/ardihikaru/graphql-example-part-1/pkg/middleware"
	rmqinterface "github.com/ardihikaru/graphql-example-part-1/pkg/rabbitmq/interface"
)

var prefix string
var reqId uint64

const (
	logFormatText    = "text"
	logFormatConsole = "console"
)

// Logger is a small wrapper around a zap.Logger.
type Logger struct {
	*zap.Logger
	logHttpRequest bool
	cfg            *config.LogPublisher
	publisher      rmqinterface.Publisher
}

// params defines log parameters
type params struct {
	Type     string  `json:"type"`
	Caller   string  `json:"caller"`
	Level    string  `json:"level"`
	Hostname *string `json:"hostname"`

	// specific only for API request logs
	RequestId     *string        `json:"request_id"`
	Proto         *string        `json:"proto"`
	Method        *string        `json:"method"`
	Path          *string        `json:"path"`
	Message       *string        `json:"msg"`
	LatencyInStr  *string        `json:"latency_in_str"`
	LatencyInTime *time.Duration `json:"latency_in_time"`
}

// init initializes prefix value
// FYI: this function adopts go-chi middleware
func init() {
	var err error
	var hostname string

	hostname, err = os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		_, err := rand.Read(buf[:])
		if err != nil {
			return
		}
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}

	prefix = fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

// New creates a new Logger with given logLevel and logFormat as part of a permanent field of the logger.
func New(logLevel, logFormat string, logHttpRequest bool, cfg *config.LogPublisher) (*Logger, error) {
	if logFormat == logFormatText {
		logFormat = logFormatConsole
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Encoding = logFormat

	var level zapcore.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, err
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("could not build logger: %w", err)
	}

	zap.ReplaceGlobals(logger)

	return &Logger{Logger: logger, logHttpRequest: logHttpRequest, cfg: cfg}, nil
}

// getRequestId is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
// FYI: this function adopts go-chi middleware
func getRequestId(r *http.Request) string {
	requestId := r.Header.Get(mw.RequestId)
	if requestId == "" {
		myId := atomic.AddUint64(&reqId, 1)
		requestId = fmt.Sprintf("%s-%06d", prefix, myId)
	}

	return requestId
}

// SetLogger returns a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return.
// Inspired by https://github.com/treastech/logger
// func SetLogger(l *Logger, publisher publisher) func(next http.Handler) http.Handler {
func SetLogger(logger *Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			logParams := buildParams(r, time.Now())
			defer func() {
				logger.Notice("served request",
					zap.String("proto", *logParams.Proto),
					zap.String("method", *logParams.Method),
					zap.String("path", *logParams.Path),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.Duration("lat", *logParams.LatencyInTime),
					zap.String("latStr", *logParams.LatencyInStr),
					zap.String("reqId", getRequestId(r)))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// buildParams builds log parameters
func buildParams(r *http.Request, timeNow time.Time) *params {
	latency := time.Since(timeNow)
	latencyStr := (time.Since(timeNow)).String()
	message := "served request"
	requestId := getRequestId(r)

	return &params{
		Type:          enum.RequestLog,
		Caller:        "",
		Level:         "",
		Proto:         &r.Proto,
		Method:        &r.Method,
		Path:          &r.URL.Path,
		Message:       &message,
		LatencyInStr:  &latencyStr,
		LatencyInTime: &latency,
		RequestId:     &requestId,
	}
}

// Notice wraps console message log in info level
func (logger *Logger) Notice(msg string, field ...zap.Field) {
	logger.Logger.Info(msg, field...)

	// when log publisher exists, publish log to the designated channel
	// FYI: ALWAYS enabled
	if logger.cfg.Notice {
		logger.publishLog(toMap(loglevel.Notice, msg, field))
	}
}

// Info wraps console message log in info level
func (logger *Logger) Info(msg string, field ...zap.Field) {
	logger.Logger.Info(msg, field...)

	// when log publisher exists, publish log to the designated channel
	// FYI: maybe disabled
	if logger.cfg.Info {
		logger.publishLog(toMap(loglevel.Info, msg, field))
	}
}

// Warn wraps console message log in warn level
func (logger *Logger) Warn(msg string, field ...zap.Field) {
	logger.Logger.Warn(msg, field...)

	// when log publisher exists, publish log to the designated channel
	// FYI: maybe disabled
	if logger.cfg.Warn {
		logger.publishLog(toMap(loglevel.Warn, msg, field))
	}
}

// Error wraps console message log in error level
func (logger *Logger) Error(msg string, field ...zap.Field) {
	logger.Logger.Error(msg, field...)

	// when log publisher exists, publish log to the designated channel
	// FYI: ALWAYS enabled
	if logger.cfg.Error {
		logger.publishLog(toMap(loglevel.Error, msg, field))
	}
}

// Debug wraps console message log in debug mode
func (logger *Logger) Debug(msg string, field ...zap.Field) {
	logger.Logger.Debug(msg, field...)

	// when log publisher exists, publish log to the designated channel
	// FYI: maybe disabled
	if logger.cfg.Debug {
		logger.publishLog(toMap(loglevel.Debug, msg, field))
	}
}

// toMap casts zap.Field into a map
func toMap(logLevel, msg string, fields []zap.Field) map[string]interface{} {
	zapAsMap := make(map[string]interface{})
	zapAsMap["message"] = msg
	zapAsMap["level"] = logLevel

	if fields != nil {
		for _, field := range fields {
			if field.Integer != 0 {
				zapAsMap[field.Key] = field.Key
			} else if field.Interface != nil {
				zapAsMap[field.Key] = field.Interface
			} else {
				// default
				zapAsMap[field.Key] = field.String
			}
		}
	}

	return zapAsMap
}

// publishLog publishes log
func (logger *Logger) publishLog(logData map[string]interface{}) {
	if !logger.logHttpRequest {
		// do nothing here
		logger.Debug("FYI: log publication is currently disabled. Do not log the request activity")

		return
	}

	if logger.logHttpRequest && logger.publisher == nil {
		// got an error
		logger.Error("FYI: log publication is ENABLED, but the publisher service is missing")

		// TODO: store the log to the persistence storage

		return
	}

	logger.Debug("publishing the captured log ...")

	responseBytes, err := json.Marshal(logData)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to convert log params to bytes: %s", err))

		// TODO: store into a persistence storage

		return
	}

	err = logger.publisher.PublishLogToElasticSearch(enum.RequestLog, responseBytes)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to publish log: %s", err))

		// TODO: store into a persistence storage

		return
	}
}
