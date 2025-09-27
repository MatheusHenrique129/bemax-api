package mocks

import (
	"github.com/MatheusHenrique129/bemax-api/internal/core/ports"
)

type mockLogger struct{}

func (m *mockLogger) Info(message string, tags ...interface{})             {}
func (m *mockLogger) Warn(message string, tags ...interface{})             {}
func (m *mockLogger) Debug(message string, tags ...interface{})            {}
func (m *mockLogger) Fatal(message string, tags ...interface{})            {}
func (m *mockLogger) Error(message string, err error, tags ...interface{}) {}

func NewMockLogger() ports.Logger {
	return &mockLogger{}
}
