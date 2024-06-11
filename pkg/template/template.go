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

	// Clone the func map because we are adding context-specific functions.
	var funcMap template.FuncMap = map[string]interface{}{}
	for k, v := range sprig.FuncMap() {
		funcMap[k] = v
	}

	// copied from: https://github.com/helm/helm/blob/8648ccf5d35d682dcd5f7a9c2082f0aaf071e817/pkg/engine/engine.go#L147-L154
	funcMap["include"] = func(name string, data interface{}) (string, error) {
		buf := bytes.NewBuffer(nil)
		if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	// We add mimir specific functions to the template
	funcMap["grafanaExploreURL"] = grafanaExploreURL
	funcMap["queryFromGeneratorURL"] = queryFromGeneratorURL

	tpl, err := tpl.Funcs(funcMap).ParseGlob(templateLocation)
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
