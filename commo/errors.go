package commo

import "errors"

var (
	ErrIncorrectEmail     = errors.New("could not parse email address")
	ErrMissingRecipient   = errors.New("missing email recipient(s)")
	ErrMissingSender      = errors.New("missing email sender")
	ErrMissingSubject     = errors.New("missing email subject")
	ErrMissingTemplate    = errors.New("missing email template name")
	ErrNotInitialized     = errors.New("email sending method has not been configured")
	ErrTemplatesNotLoaded = errors.New("templates have not been loaded yet")
)

var (
	ErrConfigConflict       = errors.New("invalid configuration: cannot specify configuration for both smtp and sendgrid")
	ErrConfigCRAMMD5Auth    = errors.New("invalid configuration: smtp cram-md5 requires username and password")
	ErrConfigInvalidSender  = errors.New("invalid configuration: could not parse sender email address")
	ErrConfigInvalidSupport = errors.New("invalid configuration: could not parse support email address")
	ErrConfigMissingPort    = errors.New("invalid configuration: smtp port is required")
	ErrConfigMissingSender  = errors.New("invalid configuration: sender email is required")
	ErrConfigPoolSize       = errors.New("invalid configuration: smtp connections pool size must be greater than zero")
)
