package unittest

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/giantswarm/apiextensions/v3/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/yaml"
)

type Config struct {
	OutputDir            string
	T                    *testing.T
	Marshaller           Marshaller
	TestFunc             TestFunc
	TestFuncReturnsBytes bool
	Update               bool
}

type Marshaller func(o interface{}) ([]byte, error)
type TestFunc func(interface{}) (interface{}, error)

// Runner is used to run unit test for a specific resource.
// It does so by running TestFunc with different input and compare the result
// with expected outputs.
//
// TestFunc is a function which takes the observed kubernetes object as input
// (e.g. AWSConfig) and returns another kubernetes object (e.g. Service).
//
// OutputDir holds yaml files, representing the yaml version of the object
// returned by TestFunc.
// Files are mapped 1 to 1 from input to output directory.
// e.g. when a file called `foo` is placed in the input directory, a
// corresponding file named `foo` must be placed in the output directory.
//
// Input directory is harcoded as the input directory in this package.
type Runner struct {
	OutputDir            string
	T                    *testing.T
	Marshaller           Marshaller
	TestFunc             TestFunc
	TestFuncReturnsBytes bool
	Update               bool

	files    []os.FileInfo
	inputDir string
}

// Value represents a test case.
// Name is the name of the test case.
// Input is the input kubernetes object.
// Output is the expected output.
type Value struct {
	Name   string
	Input  pkgruntime.Object
	Output []byte
}

// NewRunner creates a new Runner given a Config.
func NewRunner(config Config) (*Runner, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, microerror.Mask(executionError)
	}

	inputDir, err := filepath.Abs(filepath.Join(path.Dir(filename), "input"))
	if err != nil {
		return nil, microerror.Mask(err)
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var marshaller Marshaller = yaml.Marshal
	if config.Marshaller != nil {
		marshaller = config.Marshaller
	}

	r := &Runner{
		OutputDir:            config.OutputDir,
		T:                    config.T,
		Marshaller:           marshaller,
		TestFunc:             config.TestFunc,
		TestFuncReturnsBytes: config.TestFuncReturnsBytes,
		Update:               config.Update,
		inputDir:             inputDir,
		files:                files,
	}

	return r, nil
}

// Run execute all the test using testing/T.Run function.
func (r *Runner) Run() error {
	for _, file := range r.files {
		r.T.Run(file.Name(), func(t *testing.T) {
			input, err := inputValue(filepath.Join(r.inputDir, file.Name()))
			if err != nil {
				t.Fatal(err)
			}

			result, err := r.TestFunc(input)
			if err != nil {
				t.Fatal(err)
			}

			var testResult []byte
			if r.TestFuncReturnsBytes {
				testResult = result.([]byte)
			} else {
				testResult, err = r.Marshaller(result)
				if err != nil {
					t.Fatal(err)
				}
			}

			outputFile := filepath.Join(r.OutputDir, file.Name())
			if r.Update {
				err := ioutil.WriteFile(outputFile, testResult, 0644) // #nosec
				if err != nil {
					t.Fatal(err)
				}
			}

			output, err := ioutil.ReadFile(outputFile)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(testResult, output) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(output), string(testResult)))
			}
		})
	}

	return nil
}

// inputValue decode the input file as a kubernetes object and returns it.
func inputValue(inputFile string) (pkgruntime.Object, error) {
	// Read the file.
	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Create a decoder capable of decoding kubernetes objects but also
	// Giant Swarm objects.
	s := pkgruntime.NewScheme()
	err = scheme.AddToScheme(s)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	err = capi.AddToScheme(s)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	err = v1alpha1.AddToScheme(s)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	codecs := serializer.NewCodecFactory(s)
	deserializer := codecs.UniversalDeserializer()

	// Do the acutal decoding.
	input, err := pkgruntime.Decode(deserializer, inputData)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return input, nil
}

// ProjectRoot returns absolute path to prometheus-meta-operator root directory.
// This comes in handy when you need to access files in this repository from a test.
func ProjectRoot() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot get current filename")
	}

	return path.Join(path.Dir(filename), "../..")
}
