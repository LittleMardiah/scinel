package service

// This file contains shared mocks for service tests.
// Individual service tests are in their own files.

type MockLogger struct{}

func (l *MockLogger) Info(msg string, fields ...interface{}) {}
func (l *MockLogger) Error(msg string, fields ...interface{}) {}
func (l *MockLogger) Warn(msg string, fields ...interface{})  {}
