package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestLogLevel(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedLevel zapcore.Level
	}{
		{
			name:          "default",
			input:         "",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "debug",
			input:         "debug",
			expectedLevel: zapcore.DebugLevel,
		},
		{
			name:          "info",
			input:         "info",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "warn",
			input:         "warn",
			expectedLevel: zapcore.WarnLevel,
		},
		{
			name:          "error",
			input:         "error",
			expectedLevel: zapcore.ErrorLevel,
		},
		{
			name:          "dpanic",
			input:         "dpanic",
			expectedLevel: zapcore.DPanicLevel,
		},
		{
			name:          "panic",
			input:         "panic",
			expectedLevel: zapcore.PanicLevel,
		},
		{
			name:          "fatal",
			input:         "fatal",
			expectedLevel: zapcore.FatalLevel,
		},
		{
			name:          "unknown",
			input:         "asdasd",
			expectedLevel: zapcore.InfoLevel,
		},
		{
			name:          "spaces",
			input:         "   debug  ",
			expectedLevel: zapcore.DebugLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level := LogLevel(tc.input)
			assert.Equal(t, tc.expectedLevel, level)
		})
	}
}

func TestNewZapLogger_Smoke(t *testing.T) {
	level := LogLevel("info")

	logger := NewZapLoggerAdapter(level)
	logger.Info("info: %s", "test")
	logger.Warn("warn: %d", 123)
	logger.Debug("debug: %v", struct{ A int }{A: 3})
	logger.Error("test error: %v", assert.AnError)
}
