package logger

import (
	"strings"
	"time"

	"github.com/MatheusHenrique129/bemax-backend/internal/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLevel  = "DEBUG"
	infoLevel   = "INFO"
	warnLevel   = "WARN"
	errorLevel  = "ERROR"
	dpanicLevel = "DPANIC"
	panicLevel  = "PANIC"
	fatalLevel  = "FATAL"
)

type loggerAdapter struct {
	logger *zap.SugaredLogger
}

func (l loggerAdapter) Info(message string, tags ...interface{}) {
	l.logger.Infow(message, tags...)
	_ = l.logger.Sync()
}

func (l loggerAdapter) Warn(message string, tags ...interface{}) {
	l.logger.Warnw(message, tags...)
	_ = l.logger.Sync()
}

func (l loggerAdapter) Debug(message string, args ...interface{}) {
	l.logger.Debugw(message, args...)
	_ = l.logger.Sync()
}

func (l loggerAdapter) Fatal(message string, args ...interface{}) {
	l.logger.Fatalw(message, args...)
	_ = l.logger.Sync()
}

func (l loggerAdapter) Error(message string, err error, tags ...interface{}) {
	tags = append(tags, zap.NamedError("error", err))
	l.logger.Errorw(message, tags...)
	_ = l.logger.Sync()
}

func LogLevel(level string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case debugLevel:
		return zapcore.DebugLevel
	case infoLevel:
		return zap.InfoLevel
	case warnLevel:
		return zap.WarnLevel
	case errorLevel:
		return zapcore.ErrorLevel
	case dpanicLevel:
		return zap.DPanicLevel
	case panicLevel:
		return zap.PanicLevel
	case fatalLevel:
		return zap.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewZapLoggerAdapter(logLevel zapcore.Level) ports.Logger {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(logLevel)
	config.EncoderConfig.EncodeTime = rfc3399NanoTimeEncoder

	zapLogger, _ := config.Build(zap.AddCallerSkip(1))
	sugar := zapLogger.Sugar()

	return &loggerAdapter{logger: sugar}
}

// rfc3399NanoTimeEncoder serializes a time.Time to an RFC3399-formatted string
// with microsecond precision padded with zeroes to make it fixed width.
func rfc3399NanoTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const RFC3339Micro = "2006-01-02T15:04:05.000000Z07:00"
	enc.AppendString(t.UTC().Format(RFC3339Micro))
}
