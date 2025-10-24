# commo

An email rendering and sending package that can be configured to use either SendGrid or SMTP.

* GitHub: <https://github.com/rotationalio/commo>
* Go Docs: <https://go.rtnl.ai/commo>

## Usage

First, add it to your module with `go get go.rtnl.ai/commo`.

```go
import "go.rtnl.ai/commo"

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

See the test `TestLiveEmails` in [`commo_test.go`](./commo_test.go) for a full working example.

## License

See [LICENSE](./LICENSE)

## Naming

See <https://en.wikipedia.org/wiki/Communications_officer> for information on why this library is named "COMMO".
