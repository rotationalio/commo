package commo

import (
	"bytes"
	"fmt"
)

// Render returns the text and html executed templates for the specified name
// and data. Ensure that the extension is not supplied to the render method.
func Render(name string, data any) (text, html []byte, err error) {
	if text, err = render(name+".txt", data); err != nil {
		return nil, nil, err
	}

	if html, err = render(name+".html", data); err != nil {
		return nil, nil, err
	}

	return text, html, nil
}

// Render returns the text and html executed templates as strings for the
// specified name and data. Ensure that the extension is not supplied to the
// render method.
func RenderString(name string, data any) (text, html string, err error) {
	var (
		tb []byte
		hb []byte
	)

	if tb, hb, err = Render(name, data); err != nil {
		return "", "", nil
	}

	return string(tb), string(hb), nil
}

func render(name string, data any) (_ []byte, err error) {
	if templs == nil {
		return nil, ErrTemplatesNotLoaded
	}

	t, ok := templs[name]
	if !ok {
		return nil, fmt.Errorf("could not find %q in templates", name)
	}

	buf := &bytes.Buffer{}
	if err = t.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
