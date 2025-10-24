package commo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.rtnl.ai/commo"
)

func TestEmailValidate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []*commo.Email{
			{
				"admin@server.com",
				[]string{"test@example.com"},
				"This is a test email",
				"test",
				nil,
			},
			{
				"admin@server.com",
				[]string{"test@example.com"},
				"This is a test email",
				"test",
				map[string]any{"count": 4},
			},
		}

		for i, email := range testCases {
			require.NoError(t, email.Validate(), "test case %d failed", i)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := []struct {
			email *commo.Email
			err   error
		}{
			{
				&commo.Email{
					Sender:   "admin@server.com",
					To:       []string{"test@example.com"},
					Template: "test",
					Data:     nil,
				},
				commo.ErrMissingSubject,
			},
			{
				&commo.Email{
					To:       []string{"test@example.com"},
					Subject:  "This is a test email",
					Template: "test",
					Data:     nil,
				},
				commo.ErrMissingSender,
			},
			{
				&commo.Email{
					Sender:   "admin@server.com",
					Subject:  "This is a test email",
					Template: "test",
					Data:     nil,
				},
				commo.ErrMissingRecipient,
			},
			{
				&commo.Email{
					Sender:  "admin@server.com",
					To:      []string{"test@example.com"},
					Subject: "This is a test email",
					Data:    nil,
				},
				commo.ErrMissingTemplate,
			},
			{
				&commo.Email{
					Sender:   "admin@@server",
					To:       []string{"test@example.com"},
					Subject:  "This is a test email",
					Template: "test",
					Data:     nil,
				},
				commo.ErrIncorrectEmail,
			},
			{
				&commo.Email{
					Sender:   "admin@server.com",
					To:       []string{"@example.com"},
					Subject:  "This is a test email",
					Template: "test",
					Data:     nil,
				},
				commo.ErrIncorrectEmail,
			},
		}

		for i, tc := range testCases {
			require.ErrorIs(t, tc.email.Validate(), tc.err, "test case %d failed", i)
		}
	})

}
