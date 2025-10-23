package commo_test

import (
	"os"
	"testing"

	"github.com/rotationalio/confire"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/commo/commo"
)

var testEnv = map[string]string{
	"EMAIL_SENDER":            "Jane Szack <jane@example.com>",
	"EMAIL_SENDER_NAME":       "Jane Szack",
	"EMAIL_TESTING":           "true",
	"EMAIL_SMTP_HOST":         "smtp.example.com",
	"EMAIL_SMTP_PORT":         "25",
	"EMAIL_SMTP_USERNAME":     "jszack",
	"EMAIL_SMTP_PASSWORD":     "supersecret",
	"EMAIL_SMTP_USE_CRAM_MD5": "true",
	"EMAIL_SMTP_POOL_SIZE":    "16",
	"EMAIL_SENDGRID_API_KEY":  "sg:fakeapikey",
}

func TestConfig(t *testing.T) {
	// Set required environment variables and cleanup after the test is complete.
	t.Cleanup(cleanupEnv())
	setEnv()

	// NOTE: no validation is run while creating the config from the environment
	conf, err := config()
	require.Equal(t, testEnv["EMAIL_SENDER"], conf.Sender)
	require.Equal(t, testEnv["EMAIL_SENDER_NAME"], conf.SenderName)
	require.True(t, conf.Testing)
	require.Equal(t, testEnv["EMAIL_SMTP_HOST"], conf.SMTP.Host)
	require.Equal(t, uint16(25), conf.SMTP.Port)
	require.Equal(t, testEnv["EMAIL_SMTP_USERNAME"], conf.SMTP.Username)
	require.Equal(t, testEnv["EMAIL_SMTP_PASSWORD"], conf.SMTP.Password)
	require.True(t, conf.SMTP.UseCRAMMD5)
	require.Equal(t, 16, conf.SMTP.PoolSize)
	require.Equal(t, testEnv["EMAIL_SENDGRID_API_KEY"], conf.SendGrid.APIKey)
	require.NoError(t, err, "could not process configuration from the environment")
}

func TestConfigAvailable(t *testing.T) {
	testCases := []struct {
		conf   commo.Config
		assert require.BoolAssertionFunc
	}{
		{
			commo.Config{},
			require.False,
		},
		{
			commo.Config{
				SMTP: commo.SMTPConfig{Host: "email.example.com"},
			},
			require.True,
		},
		{
			commo.Config{
				SendGrid: commo.SendGridConfig{APIKey: "sg:fakeapikey"},
			},
			require.True,
		},
	}

	for i, tc := range testCases {
		tc.assert(t, tc.conf.Available(), "test case %d failed", i)
	}
}

func TestConfigValidation(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []commo.Config{
			{
				Testing: false,
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SMTP: commo.SMTPConfig{
					Host:       "smtp.example.com",
					Port:       587,
					Username:   "admin",
					Password:   "supersecret",
					UseCRAMMD5: false,
					PoolSize:   4,
				},
			},
			{
				Testing: true,
			},
			{
				Sender:  "Peony Quarterdeck <peony@example.com>",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
			{
				Sender:  "peony@example.com",
				Testing: false,
				SendGrid: commo.SendGridConfig{
					APIKey: "sg:fakeapikey",
				},
			},
		}

		for i, conf := range testCases {
			require.NoError(t, conf.Validate(), "test case %d failed", i)
		}

	})

	t.Run("Invalid", func(t *testing.T) {
		testCases := []struct {
			conf commo.Config
			err  error
		}{
			{
				commo.Config{
					Testing: false,
					SMTP:    commo.SMTPConfig{Host: "email.example.com"},
				},
				commo.ErrConfigMissingSender,
			},
			{
				commo.Config{
					Testing:  false,
					SendGrid: commo.SendGridConfig{APIKey: "sg:fakeapikey"},
				},
				commo.ErrConfigMissingSender,
			},
			{
				commo.Config{
					Sender:  "foo",
					Testing: false,
					SMTP:    commo.SMTPConfig{Host: "smtp.example.com"},
				},
				commo.ErrConfigInvalidSender,
			},
			{
				commo.Config{
					Sender:  "orchid@example.com",
					Testing: false,
					SMTP:    commo.SMTPConfig{Host: "smtp.example.com"},
				},
				commo.ErrConfigInvalidSupport,
			},
			{
				commo.Config{
					Sender:  "orchid@example.com",
					Testing: false,
					SMTP: commo.SMTPConfig{
						Host: "smtp.example.com",
					},
					SendGrid: commo.SendGridConfig{
						APIKey: "sg:fakeapikey",
					},
				},
				commo.ErrConfigConflict,
			},
			{
				commo.Config{
					Sender:  "orchid@example.com",
					Testing: false,
					SMTP: commo.SMTPConfig{
						Host: "smtp.example.com",
						Port: 0,
					},
				},
				commo.ErrConfigMissingPort,
			},
			{
				commo.Config{
					Sender:  "orchid@example.com",
					Testing: false,
					SMTP: commo.SMTPConfig{
						Host: "smtp.example.com",
						Port: 527,
					},
				},
				commo.ErrConfigPoolSize,
			},
			{
				commo.Config{
					Sender:  "orchid@example.com",
					Testing: false,
					SMTP: commo.SMTPConfig{
						Host:       "smtp.example.com",
						Port:       527,
						PoolSize:   4,
						UseCRAMMD5: true,
					},
				},
				commo.ErrConfigCRAMMD5Auth,
			},
		}

		for i, tc := range testCases {
			require.ErrorIs(t, tc.conf.Validate(), tc.err, "test case %d failed", i)
		}
	})
}

func TestSMTPConfig(t *testing.T) {
	t.Run("Addr", func(t *testing.T) {
		conf := commo.SMTPConfig{
			Host: "smtp.example.com",
			Port: 527,
		}
		require.Equal(t, "smtp.example.com:527", conf.Addr())
	})
}

func TestGetSenderName(t *testing.T) {
	testCases := []struct {
		conf     commo.Config
		expected string
	}{
		{
			commo.Config{
				Sender: "Jane Szack <jane@example.com>",
			},
			"Jane Szack",
		},
		{
			commo.Config{
				Sender: "jane@example.com",
			},
			"",
		},
		{
			commo.Config{
				Sender:     "",
				SenderName: "Jane Szack",
			},
			"Jane Szack",
		},
		{
			commo.Config{
				Sender:     "",
				SenderName: "",
			},
			"",
		},
		{
			commo.Config{
				Sender:     "foo",
				SenderName: "",
			},
			"",
		},
		{
			commo.Config{
				Sender:     "John Doe <john.doe@example.com>",
				SenderName: "Jane Szack",
			},
			"Jane Szack",
		},
		{
			commo.Config{
				Sender:     "john.doe@example.com",
				SenderName: "John Doe",
			},
			"John Doe",
		},
	}

	for i, tc := range testCases {
		require.Equal(t, tc.expected, tc.conf.GetSenderName(), "test case %d failed", i)
	}
}

// Creates a new email config from the current environment.
func config() (conf commo.Config, err error) {
	if err = confire.Process("email", &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

// Returns the current environment for the specified keys, or if no keys are specified
// then it returns the current environment for all keys in the testEnv variable.
func curEnv(keys ...string) map[string]string {
	env := make(map[string]string)
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := os.LookupEnv(key); ok {
				env[key] = val
			}
		}
	} else {
		for key := range testEnv {
			env[key] = os.Getenv(key)
		}
	}

	return env
}

// Sets the environment variables from the testEnv variable. If no keys are specified,
// then this function sets all environment variables from the testEnv.
func setEnv(keys ...string) {
	if len(keys) > 0 {
		for _, key := range keys {
			if val, ok := testEnv[key]; ok {
				os.Setenv(key, val)
			}
		}
	} else {
		for key, val := range testEnv {
			os.Setenv(key, val)
		}
	}
}

// Cleanup helper function that can be run when the tests are complete to reset the
// environment back to its previous state before the test was run.
func cleanupEnv(keys ...string) func() {
	prevEnv := curEnv(keys...)
	return func() {
		for key, val := range prevEnv {
			if val != "" {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
