package logger

import (
	"fmt"
	"time"

	"github.com/MapMinder/mapminder_backend/internal/config"
	"github.com/MapMinder/mapminder_backend/internal/status"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes logger
func Init(zapConfig *config.ZapConfig) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapConfig.LogLevel),
		Development:      zapConfig.IsDevelopment,
		Encoding:         zapConfig.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	log, err = config.Build()
	if err != nil {
		log.Sugar().Fatalf("Failed to build logger: %s", err)
	}

	log.Info("Logger initialized")
}

func Get() *zap.Logger {
	return log
}

// MiddlewareLogger logger for gin to log for requests
func MiddlewareLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// create middleware logs for requests
		logger.Info(
			"Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

func Sync() {
	log.Sync()
}

// NOTE: 全てのログケースをカバーするにはすごく時間かかるため必要ある際に追加する

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func InfoStatus(s *status.Status, fields ...zap.Field) {
	customizedFields := append(fields, zap.Int("Code", s.Code))
	log.Info(s.Message, customizedFields...)
}

func Infow(msg string, fields ...interface{}) {
	log.Sugar().Infow(msg, fields...)
}

func Infof(templateString string, fields ...interface{}) {
	log.Sugar().Infof(templateString, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func WarnStatus(s *status.Status, fields ...zap.Field) {
	customizedFields := append(fields, zap.Int("Code", s.Code))
	log.Warn(s.Message, customizedFields...)
}

func Warnw(msg string, fields ...interface{}) {
	log.Sugar().Warnw(msg, fields...)
}

func Warnf(templateString string, fields ...interface{}) {
	log.Sugar().Warnf(templateString, fields...)
}

func Error(err error, fields ...zap.Field) {
	log.Error(err.Error(), fields...)
}

func ErrorStatus(s *status.Status, fields ...zap.Field) {
	customizedFields := append(fields, zap.Int("Code", s.Code))
	log.Error(s.Message, customizedFields...)
}

func Errorw(msg string, err error, fields ...interface{}) {
	log.Sugar().Errorw(fmt.Sprintf(msg+" %s", err.Error()), fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func Fatalf(templateString string, fields ...interface{}) {
	log.Sugar().Fatalf(templateString, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	log.Panic(msg, fields...)
}

func InitForTest() {
	log = zap.NewNop()
}
