package commo_test

import (
	"embed"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rotationalio/confire"
	"github.com/stretchr/testify/require"
	"go.rtnl.ai/commo/commo"
)

func TestLiveEmails(t *testing.T) {
	// Load local .env if it exists to make setting envvars easier.
	godotenv.Load("testdata/.env")

	// This test will send actual emails to an account as configured by the environment.
	// The $TEST_LIVE_EMAILS environment variable must be set to true to not skip.
	SkipByEnvVar(t, "TEST_LIVE_EMAILS")
	CheckEnvVars(t, "TEST_LIVE_EMAIL_RECIPIENT")

	// Configure email sending from the environment. See .env.template for requirements.
	conf := commo.Config{}
	err := confire.Process("commo_email", &conf)
	require.NoError(t, err, "environment not setup to send live emails; see .env.template")
	require.True(t, conf.Available(), "no backend setup to send live emails; see .env.template")

	err = commo.Initialize(conf, loadTestTemplates())
	require.NoError(t, err, "could not configure email sending")

	recipient := os.Getenv("TEST_LIVE_EMAIL_RECIPIENT")

	t.Run("TestEmail", func(t *testing.T) {
		data := struct{ ContactName string }{ContactName: "User Name"}

		email, err := commo.New(recipient, "Test Subject", "test_email", data)
		require.NoError(t, err, "could not create reset password email")

		err = email.Send()
		require.NoError(t, err, "could not send reset password email")
	})
}

// ############################################################################
// Helpers
// ############################################################################

func CheckEnvVars(t *testing.T, envs ...string) {
	for _, env := range envs {
		require.NotEmpty(t, os.Getenv(env), "required environment variable $%s not set", env)
	}
}

func SkipByEnvVar(t *testing.T, env string) {
	val := strings.ToLower(strings.TrimSpace(os.Getenv(env)))
	switch val {
	case "1", "t", "true":
		return
	default:
		t.Skipf("this test depends on the $%s envvar to run", env)
	}
}

// ############################################################################
// Test Template Loading
// ############################################################################

const (
	templatesDir = "testdata/templates"
	partialsDir  = "partials"
)

var (
	//go:embed testdata/templates/*.html testdata/templates/*.txt testdata/templates/partials/*html
	files embed.FS
)

// Load templates
func loadTestTemplates() map[string]*template.Template {
	var (
		err           error
		templateFiles []fs.DirEntry
	)

	templates := make(map[string]*template.Template)
	if templateFiles, err = fs.ReadDir(files, templatesDir); err != nil {
		panic(err)
	}

	// Each template needs to be parsed independently to ensure that define directives
	// are not overriden if they have the same name; e.g. to use the base template.
	for _, file := range templateFiles {
		if file.IsDir() {
			continue
		}

		// Each template will be accessible by its base name in the global map
		patterns := make([]string, 0, 2)
		patterns = append(patterns, filepath.Join(templatesDir, file.Name()))
		switch filepath.Ext(file.Name()) {
		case ".html":
			patterns = append(patterns, filepath.Join(templatesDir, partialsDir, "*.html"))
		}

		templates[file.Name()] = template.Must(template.ParseFS(files, patterns...))
	}
	return templates
}
