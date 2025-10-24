# commo

An email rendering and sending package that can be configured to use either SendGrid or SMTP.

## Usage

```go
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
```

See the test(s) in [commo/commo_test.go](./commo/commo_test.go) for a full working example.
