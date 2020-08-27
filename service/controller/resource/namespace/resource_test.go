package namespace

import (
	"bytes"
	"flag"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/prometheus-meta-operator/pkg/unittest"
	"github.com/google/go-cmp/cmp"
)

var update = flag.Bool("update", false, "update .golden CF template file")

func TestNamespace(t *testing.T) {
	outputDir, err := filepath.Abs("./output")
	if err != nil {
		t.Fatal(err)
	}

	runner, err := unittest.NewRunner(outputDir)
	if err != nil {
		t.Fatal(err)
	}

	for runner.Next() {
		tc := runner.Value()
		t.Run(tc.Name, func(t *testing.T) {
			namespace, err := toNamespace(tc.Input)
			if err != nil {
				t.Fatal(err)
			}

			testResult, err := yaml.Marshal(namespace)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(testResult, tc.Output) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(tc.Output), string(testResult)))
			}
		})
	}
	if err := runner.Err(); err != nil {
		t.Fatal(err)
	}
}
