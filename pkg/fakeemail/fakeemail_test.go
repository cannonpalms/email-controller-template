package fakeemail

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestSendValidEmail(t *testing.T) {
	config := Config{
		BounceRate: 0.0, // No simulated bounce.
		BlockRate:  0.0, // No simulated block.
		Logger:     log.New(&testLogger{logs: &[]string{}}, "", 0),
	}

	emailService, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Error creating EmailService: %v", err)
	}

	err = emailService.Send("recipient@example.com", "Test Subject", "This is the email body.")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestSendInvalidEmail(t *testing.T) {
	config := Config{
		BounceRate: 0.0, // No simulated bounce.
		BlockRate:  0.0, // No simulated block.
		Logger:     newTestLogger(),
	}

	emailService, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Error creating EmailService: %v", err)
	}

	err = emailService.Send("invalid_email", "Test Subject", "This is the email body.")
	if _, ok := err.(*ErrInvalidEmailAddress); !ok {
		t.Errorf("Expected ErrInvalidEmailAddress, got: %v", err)
	}

	assertLogs(t, config.Logger, []string{})
}

func TestSendSimulatedBounce(t *testing.T) {
	config := Config{
		BounceRate: 1.0, // Simulate bounce.
		BlockRate:  0.0, // No simulated block.
		Logger:     newTestLogger(),
	}

	emailService, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Error creating EmailService: %v", err)
	}

	err = emailService.Send("recipient@example.com", "Test Subject", "This is the email body.")
	if _, ok := err.(*ErrEmailBounced); !ok {
		t.Errorf("Expected ErrEmailBounced, got: %v", err)
	}

	assertLogs(t, config.Logger, []string{"Simulating email bounce: recipient@example.com"})
}

func TestSendSimulatedBlock(t *testing.T) {
	config := Config{
		BounceRate: 0.0, // No simulated bounce.
		BlockRate:  1.0, // Simulate block.
		Logger:     newTestLogger(),
	}

	emailService, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Error creating EmailService: %v", err)
	}

	err = emailService.Send("recipient@example.com", "Test Subject", "This is the email body.")
	if _, ok := err.(*ErrEmailBlocked); !ok {
		t.Errorf("Expected ErrEmailBlocked, got: %v", err)
	}

	assertLogs(t, config.Logger, []string{"Simulating email block: recipient@example.com"})
}

func assertLogs(t *testing.T, logger *log.Logger, expectedLogs []string) {
	t.Helper()

	if len(expectedLogs) != 0 {
		// If we expect logs, convert actual logs to string slices.
		actualLogs := loggerStringSlice(logger)
		if len(actualLogs) != len(expectedLogs) {
			t.Errorf("Expected %d log entries, got %d", len(expectedLogs), len(actualLogs))
		}

		for i, logEntry := range actualLogs {
			if logEntry != expectedLogs[i] {
				t.Errorf("Expected log entry: \"%s\", got: \"%s\"", expectedLogs[i], logEntry)
			}
		}
	} else {
		// If we expect no logs, ensure the logger has no entries.
		actualLogs := loggerStringSlice(logger)
		if len(actualLogs) != 0 {
			t.Errorf("Expected no logs, but got %d log entries", len(actualLogs))
		}
	}
}

func loggerStringSlice(logger *log.Logger) []string {
	// Helper function to convert logger entries to a string slice.
	logs := logger.Writer().(*testLogger).logs
	trimmedLogs := []string{}
	for _, log := range *logs {
		trimmedLogs = append(trimmedLogs, strings.TrimSpace(log))
	}
	return trimmedLogs
}

type testLogger struct {
	logs *[]string
}

func newTestLogger() *log.Logger {
	return log.New(&testLogger{logs: &[]string{}}, "", 0)
}

func (tl *testLogger) Write(p []byte) (n int, err error) {
	*tl.logs = append(*tl.logs, string(p))
	return len(p), nil
}

func (tl *testLogger) Println(v ...interface{}) {
	*tl.logs = append(*tl.logs, fmt.Sprintln(v...))
}

func (tl *testLogger) Printf(format string, v ...interface{}) {
	*tl.logs = append(*tl.logs, fmt.Sprintf(format, v...))
}
