package main

import (
	"fmt"
	"log"

	"github.com/cannonpalms/email-controller-template/pkg/fakeemail"
	"go.uber.org/zap"
)

func main() {

	// Initialize a logger compatible with log.Logger. Rather than initializing your own logger,
	// you can use the logger you derive from the context passed to the Reconcile function of any controller
	// you create.
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapLogger.Sync() // Flushes any buffered log entries before the program exits.

	logger := zap.NewStdLog(zapLogger)

	// Create a new EmailService with custom configuration and Zap logger.
	emailService, err := fakeemail.NewEmailService(fakeemail.Config{
		BounceRate: 0.2,    // Custom bounce rate is set to 0.2 (20%).
		BlockRate:  0.05,   // Custom block rate is set to 0.05 (5%).
		Logger:     logger, // Use the Zap logger.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Simulate sending an email.
	err = emailService.Send("recipient@example.com", "Hello", "This is a test email.")
	if err != nil {
		switch {
		case err.(*fakeemail.ErrInvalidEmailAddress) != nil:
			fmt.Println("Error: Invalid email address")
		case err.(*fakeemail.ErrEmailBounced) != nil:
			fmt.Println("Error: Email bounced")
		case err.(*fakeemail.ErrEmailBlocked) != nil:
			fmt.Println("Error: Email blocked")
		default:
			fmt.Printf("Error: %v\n", err)
		}
	}
}
