package templates

import (
	"bytes"
	"text/template"

	"github.com/giantswarm/microerror"
)

func Render(content string, data interface{}) (string, error) {
	var err error

	main := template.New("main")
	main, err = main.Parse(content)
	if err != nil {
		return "", microerror.Mask(err)
	}

	var b bytes.Buffer
	err = main.ExecuteTemplate(&b, "main", data)
	if err != nil {
		return "", microerror.Mask(err)
	}

	return b.String(), nil
}
