/*
An email rendering and sending package that can be configured to use either SendGrid or SMTP.

Usage Example:

	// Load configuration from a .env file
	conf := commo.Config{}
	err := confire.Process("commo_email", &conf)
	checkErr(err)

	// Load templates
	var templates map[string]*template.Template
	templates = ... // not shown here; see `commo/commo_test.go` for full example

	// Initialize commo
	err = commo.Initialize(conf, templates)
	checkErr(err)

	// Data is a struct that has all of the required fields for the template being used
	data := struct{ ContactName string }{ContactName: "User Name"}

	// Create the email
	email, err := commo.New("Test User <test@example.com>", "Email Subject", "template_name_no_ext", data)
	checkErr(err)

	// Send the email
	err = email.Send()
	checkErr(err)

See the test TestLiveEmails() in commo_test.go for a full working example.
*/
package commo

import (
	"context"
	"errors"
	"html/template"

	"github.com/jordan-wright/email"
	"go.rtnl.ai/x/backoff"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Package level variables
var (
	initialized bool
	config      Config
	templs      map[string]*template.Template
	pool        *email.Pool
	sg          *sendgrid.Client
)

// Config for [backoff.ExponentialBackOff]
const (
	multiplier          = 2.0
	randomizationFactor = 0.45
)

// Initialize the package to start sending emails. If there is no valid email
// configuration available then configuration is gracefully ignored without error.
func Initialize(conf Config, templates map[string]*template.Template) (err error) {
	// Do not configure email if it is not available but also do not return an error.
	if !conf.Available() {
		return nil
	}

	if err = conf.Validate(); err != nil {
		return err
	}

	// TODO: if in testing mode create a mock for sending emails.

	if conf.SMTP.Enabled() {
		if pool, err = conf.SMTP.Pool(); err != nil {
			return err
		}
	}

	if conf.SendGrid.Enabled() {
		sg = conf.SendGrid.Client()
	}

	config = conf
	templs = templates
	initialized = true
	return nil
}

// Loads templates into commo's internal template storage. Useful for testing.
func WithTemplates(templates map[string]*template.Template) {
	templs = templates
}

// Send an email using the configured send methodology. Uses exponential backoff to
// retry multiple times on error with an increasing delay between attempts.
func Send(email *Email) (err error) {
	// The package must be initialized to send.
	if !initialized {
		return ErrNotInitialized
	}

	// Select the send function to deliver the email with.
	var send sender
	switch {
	case config.SMTP.Enabled():
		send = sendSMTP
	case config.SendGrid.Enabled():
		send = sendSendGrid
	case config.Testing:
		send = sendMock
	default:
		panic("unhandled send email method")
	}

	exponential := backoff.ExponentialBackOff{
		InitialInterval:     config.Backoff.InitialInterval,
		RandomizationFactor: randomizationFactor,
		Multiplier:          multiplier,
		MaxInterval:         config.Backoff.MaxInterval,
	}

	// Attempt to send the message with multiple retries.
	if _, err = backoff.Retry(context.Background(), func() (any, serr error) {
		serr = send(email)
		return nil, serr
	},
		backoff.WithBackOff(&exponential),
		backoff.WithMaxElapsedTime(config.Backoff.MaxElapsedTime),
	); err != nil {
		return err
	}

	return nil

}

type sender func(*Email) error

func sendSMTP(e *Email) (err error) {
	var msg *email.Email
	if msg, err = e.ToSMTP(); err != nil {
		return err
	}

	if err = pool.Send(msg, config.Backoff.Timeout); err != nil {
		return err
	}
	return nil
}

func sendSendGrid(e *Email) (err error) {
	var msg *sgmail.SGMailV3
	if msg, err = e.ToSendGrid(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Backoff.Timeout)
	defer cancel()

	var rep *rest.Response
	if rep, err = sg.SendWithContext(ctx, msg); err != nil {
		return err
	}

	if rep.StatusCode < 200 || rep.StatusCode >= 300 {
		return errors.New(rep.Body)
	}

	return nil
}

func sendMock(*Email) (err error) {
	return errors.New("not implemented")
}
