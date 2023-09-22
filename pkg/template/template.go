package template

import (
	"bytes"
	"html/template"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/giantswarm/microerror"
)

// RenderTemplate render all template files matching templateLocation glob filter, using templateData.
// template files are using [[ and ]] as delimiters.
// There's an additional 'include' template function provided taken from helm: https://helm.sh/docs/howto/charts_tips_and_tricks/#using-the-include-function
func RenderTemplate(templateData interface{}, templateLocation string) ([]byte, error) {
	tpl := template.New("_base").Delims("[[", "]]")

	var funcMap template.FuncMap = map[string]interface{}{}
	// copied from: https://github.com/helm/helm/blob/8648ccf5d35d682dcd5f7a9c2082f0aaf071e817/pkg/engine/engine.go#L147-L154
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	funcMap["hasPrefix"] = strings.HasPrefix

	tpl, err := tpl.Funcs(sprig.FuncMap()).Funcs(funcMap).ParseGlob(templateLocation)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var b bytes.Buffer
	for _, t := range tpl.Templates() {
		if strings.HasPrefix(t.Name(), "_") {
			continue
		}
		err := t.Delims("[[", "]]").Execute(&b, templateData)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}
	return b.Bytes(), nil
}
