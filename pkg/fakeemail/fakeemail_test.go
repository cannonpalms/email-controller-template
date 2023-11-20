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

	email := Email{
		DestinationAddress: "recipient@example.com",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	emailID, err := emailService.Send(email)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// You can now use emailID as needed.
	fmt.Printf("Email sent successfully. Email ID: %d\n", emailID)
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

	email := Email{
		DestinationAddress: "invalid_email",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	emailID, err := emailService.Send(email)
	if _, ok := err.(*ErrInvalidEmailAddress); !ok {
		t.Errorf("Expected ErrInvalidEmailAddress, got: %v", err)
	}

	// You can still use emailID, even if there was an error.
	fmt.Printf("Email ID: %d\n", emailID)

	// Assert that the logger contains the expected log entry
	assertLogs(t, config.Logger, []string{"Invalid email address: invalid_email"})
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

	email := Email{
		DestinationAddress: "recipient@example.com",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	emailID, err := emailService.Send(email)
	if _, ok := err.(*ErrEmailBounced); !ok {
		t.Errorf("Expected ErrEmailBounced, got: %v", err)
	}

	// You can still use emailID, even if there was an error.
	fmt.Printf("Email ID: %d\n", emailID)

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

	email := Email{
		DestinationAddress: "recipient@example.com",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	emailID, err := emailService.Send(email)
	if _, ok := err.(*ErrEmailBlocked); !ok {
		t.Errorf("Expected ErrEmailBlocked, got: %v", err)
	}

	// You can still use emailID, even if there was an error.
	fmt.Printf("Email ID: %d\n", emailID)

	assertLogs(t, config.Logger, []string{"Simulating email block: recipient@example.com"})
}

func TestSendWithInvalidEmailLogger(t *testing.T) {
	config := Config{
		BounceRate: 0.0, // No simulated bounce.
		BlockRate:  0.0, // No simulated block.
		Logger:     log.New(&testLogger{logs: &[]string{}}, "", 0),
	}

	emailService, err := NewEmailService(config)
	if err != nil {
		t.Fatalf("Error creating EmailService: %v", err)
	}

	// Set an invalid email address
	email := Email{
		DestinationAddress: "invalid_email",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	emailID, err := emailService.Send(email)
	if _, ok := err.(*ErrInvalidEmailAddress); !ok {
		t.Errorf("Expected ErrInvalidEmailAddress, got: %v", err)
	}

	// You can still use emailID, even if there was an error.
	fmt.Printf("Email ID: %d\n", emailID)

	// Assert that the logger contains the expected log entry
	assertLogs(t, config.Logger, []string{"Invalid email address: invalid_email"})
}

func TestEqualEmailsProduceEqualIdentifiers(t *testing.T) {
	// Create two identical emails
	email1 := Email{
		DestinationAddress: "recipient@example.com",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	email2 := Email{
		DestinationAddress: "recipient@example.com",
		Subject:            "Test Subject",
		Body:               "This is the email body.",
	}

	// Get the identifiers without sending emails
	emailID1 := email1.ID()
	emailID2 := email2.ID()

	// Ensure that the identifiers are equal
	if emailID1 != emailID2 {
		t.Errorf("Expected equal identifiers for identical emails, got %d and %d", emailID1, emailID2)
	}

	// You can use emailID1 or emailID2 as needed.
	fmt.Printf("Email IDs: %d and %d\n", emailID1, emailID2)
}

// Helper function to assert logs in the logger.
func assertLogs(t *testing.T, logger *log.Logger, expectedLogs []string) {
	t.Helper()

	actualLogs := loggerStringSlice(logger)
	defer func() {
		t.Logf("Actual logs: %v", actualLogs)
	}()

	if len(expectedLogs) != 0 {
		// If we expect logs, convert actual logs to string slices.
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
		if len(actualLogs) != 0 {
			t.Errorf("Expected no logs, but got %d log entries", len(actualLogs))
		}
	}
}

// Helper function to convert logger entries to a string slice.
func loggerStringSlice(logger *log.Logger) []string {
	// Helper function to convert logger entries to a string slice.
	logs := logger.Writer().(*testLogger).logs
	trimmedLogs := []string{}
	for _, log := range *logs {
		trimmedLogs = append(trimmedLogs, strings.TrimSpace(log))
	}
	return trimmedLogs
}

// Helper type for testing logging.
type testLogger struct {
	logs *[]string
}

// Helper function to create a new test logger.
func newTestLogger() *log.Logger {
	return log.New(&testLogger{logs: &[]string{}}, "", 0)
}

// Implementation of the Write method for the testLogger type.
func (tl *testLogger) Write(p []byte) (n int, err error) {
	*tl.logs = append(*tl.logs, string(p))
	return len(p), nil
}

// Implementation of the Println method for the testLogger type.
func (tl *testLogger) Println(v ...interface{}) {
	*tl.logs = append(*tl.logs, fmt.Sprintln(v...))
}

// Implementation of the Printf method for the testLogger type.
func (tl *testLogger) Printf(format string, v ...interface{}) {
	*tl.logs = append(*tl.logs, fmt.Sprintf(format, v...))
}
