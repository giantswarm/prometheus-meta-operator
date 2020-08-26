package namespace

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

var update = flag.Bool("update", false, "update .golden CF template file")

func TestNamespace(t *testing.T) {
	root := "../../../../pkg/unittest/input"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range files {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			o := filepath.Join(root, file.Name())
			inputFile, err := ioutil.ReadFile(o)
			if err != nil {
				t.Fatal(err)
			}

			var obj *v1alpha2.Cluster
			err = yaml.Unmarshal(inputFile, &obj)
			if err != nil {
				t.Fatal(err)
			}
			namespace, err := toNamespace(obj)
			if err != nil {
				t.Fatal(err)
			}

			data, err := yaml.Marshal(namespace)
			if err != nil {
				t.Fatal(err)
			}

			p := filepath.Join("output", file.Name())

			if *update {
				err := ioutil.WriteFile(p, data, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}
			outputFile, err := ioutil.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(data, outputFile) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(outputFile), string(data)))
			}
		})
	}
}
