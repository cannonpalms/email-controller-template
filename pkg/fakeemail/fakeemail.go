package fakeemail

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
)

// ErrInvalidEmailAddress represents an error related to an invalid email address.
type ErrInvalidEmailAddress struct {
	Email string
}

func (e *ErrInvalidEmailAddress) Error() string {
	return fmt.Sprintf("Invalid email address: %s", e.Email)
}

// ErrEmailBounced represents an error indicating that the email bounced.
type ErrEmailBounced struct {
	Email string
}

func (e *ErrEmailBounced) Error() string {
	return fmt.Sprintf("Email bounced: %s", e.Email)
}

// ErrEmailBlocked represents an error indicating that the email is blocked.
type ErrEmailBlocked struct {
	Email string
}

func (e *ErrEmailBlocked) Error() string {
	return fmt.Sprintf("Email blocked: %s", e.Email)
}

// Config holds configuration options for the fakeemail module.
type Config struct {
	BounceRate float64     // Rate for simulating bounced emails (0.0 to 1.0).
	BlockRate  float64     // Rate for simulating blocked emails (0.0 to 1.0).
	Logger     *log.Logger // Logger for logging messages.
}

// EmailService represents an email service with configurable options.
type EmailService struct {
	Config
}

// NewEmailService creates a new EmailService with the given configuration.
func NewEmailService(config Config) (*EmailService, error) {
	if err := validateRate(config.BounceRate); err != nil {
		return nil, fmt.Errorf("Bounce rate validation error: %w", err)
	}

	if err := validateRate(config.BlockRate); err != nil {
		return nil, fmt.Errorf("Block rate validation error: %w", err)
	}

	return &EmailService{
		Config: config,
	}, nil
}

// DefaultEmailService creates a new EmailService with default configuration.
func DefaultEmailService() *EmailService {
	defaultConfig := Config{
		BounceRate: 0.25,                      // Default bounce rate is set to 0.25 (25%).
		BlockRate:  0.1,                       // Default block rate is set to 0.1 (10%).
		Logger:     log.New(os.Stdout, "", 0), // Default logger is log.New.
	}
	es, _ := NewEmailService(defaultConfig)
	return es
}

// Send logs the email's body to the console and may simulate email bouncing or blocking.
// Now accepts an Email and always returns an EmailIdentifier and an error.
func (es *EmailService) Send(email Email) (EmailIdentifier, error) {
	// Always compute the EmailIdentifier using the Hash method of the Email type.
	// We return the EmailIdentifier even if the destination address is invalid, the email is blocked,
	// or the email bounces
	id := email.ID()

	if !isValidEmail(email) {
		if es.Logger != nil {
			es.Logger.Printf("Invalid email address: %s\n", email.DestinationAddress)
		}
		return id, &ErrInvalidEmailAddress{Email: email.DestinationAddress}
	}

	if shouldSimulate(es.BounceRate) {
		if es.Logger != nil {
			es.Logger.Printf("Simulating email bounce: %s\n", email.DestinationAddress)
		}
		return id, &ErrEmailBounced{Email: email.DestinationAddress}
	}

	if shouldSimulate(es.BlockRate) {
		if es.Logger != nil {
			es.Logger.Printf("Simulating email block: %s\n", email.DestinationAddress)
		}
		return id, &ErrEmailBlocked{Email: email.DestinationAddress}
	}

	logEmail(es.Logger, email.DestinationAddress, email.Subject, email.Body)
	return id, nil
}

// isValidEmail checks if the given destination email address is valid.
func isValidEmail(email Email) bool {
	// Basic email validation using a regular expression.
	// This is a simple example, and you may want to use a more robust approach in a real-world scenario.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email.DestinationAddress)
}

// logEmail logs the email details to the console.
func logEmail(logger *log.Logger, destination, subject, body string) {
	logger.Printf("Sending email to: %s\n", destination)
	logger.Printf("Subject: %s\n", subject)
	logger.Printf("Body: %s\n", body)
}

// shouldSimulate determines whether to simulate an event based on the rate.
func shouldSimulate(rate float64) bool {
	return rand.Float64() < rate
}

// validateRate ensures that the rate is between 0 and 1.
func validateRate(rate float64) error {
	if rate < 0 || rate > 1 {
		return errors.New("Rate must be between 0 and 1")
	}
	return nil
}
