package namespace

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

var update = flag.Bool("update", false, "update .golden CF template file")

func TestNamespace(t *testing.T) {
	root := "../../../../pkg/unittest/input"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			var input runtime.Object
			{
				inputFile := filepath.Join(root, file.Name())
				inputData, err := ioutil.ReadFile(inputFile)
				if err != nil {
					t.Fatal(err)
				}

				scheme := runtime.NewScheme()
				v1alpha2.AddToScheme(scheme)
				v1alpha1.AddToScheme(scheme)
				codecs := serializer.NewCodecFactory(scheme)
				deserializer := codecs.UniversalDeserializer()
				input, err = runtime.Decode(deserializer, inputData)
				if err != nil {
					t.Fatal(err)
				}
			}

			namespace, err := toNamespace(input)
			if err != nil {
				t.Fatal(err)
			}

			var testResult []byte
			{
				testResult, err = yaml.Marshal(namespace)
				if err != nil {
					t.Fatal(err)
				}
			}

			var expectedOutput []byte
			{
				outputFile := filepath.Join("output", file.Name())
				if *update {
					err := ioutil.WriteFile(outputFile, testResult, 0644)
					if err != nil {
						t.Fatal(err)
					}
				}
				expectedOutput, err = ioutil.ReadFile(outputFile)
				if err != nil {
					t.Fatal(err)
				}
			}

			if !bytes.Equal(testResult, expectedOutput) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(expectedOutput), string(testResult)))
			}
		})
	}
}
