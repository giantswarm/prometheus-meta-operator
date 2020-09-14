package unittest

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/microerror"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

type Config struct {
	OutputDir string
	T         *testing.T
	TestFunc  func(interface{}) (metav1.Object, error)
}

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
	OutputDir string
	T         *testing.T
	TestFunc  func(interface{}) (metav1.Object, error)

	inputDir string
	files    []os.FileInfo
	current  int
	err      error
}

// Value represent a test case.
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

	r := &Runner{
		OutputDir: config.OutputDir,
		T:         config.T,
		TestFunc:  config.TestFunc,
		inputDir:  inputDir,
		files:     files,
		current:   -1,
		err:       nil,
	}

	return r, nil
}

// Run execute all the test using testing/T.Run function.
func (r *Runner) Run() error {
	for r.Next() {
		value := r.Value()
		r.T.Run(value.Name, func(t *testing.T) {
			namespace, err := r.TestFunc(value.Input)
			if err != nil {
				t.Fatal(err)
			}

			testResult, err := yaml.Marshal(namespace)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(testResult, value.Output) {
				t.Fatalf("\n\n%s\n", cmp.Diff(string(value.Output), string(testResult)))
			}
		})
	}
	if err := r.Err(); err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// Next return true when there is more test cases to run.
// There is 1 test case per input file.
func (r *Runner) Next() bool {
	if r.err != nil {
		return false
	}

	r.current++
	return len(r.files) > r.current
}

// Value returns the current test case values.
func (r *Runner) Value() Value {
	input, err := r.inputValue()
	if err != nil {
		r.err = microerror.Mask(err)
		return Value{}
	}

	output, err := r.outputValue()
	if err != nil {
		r.err = microerror.Mask(err)
		return Value{}
	}

	v := Value{
		Name:   r.files[r.current].Name(),
		Input:  input,
		Output: output,
	}

	return v
}

// inputValue decode the input file as a kubernetes object and returns it.
func (r *Runner) inputValue() (pkgruntime.Object, error) {
	// Read the file.
	inputFile := filepath.Join(r.inputDir, r.files[r.current].Name())
	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// Create a decoder capable of decoding kubernetes objects but also
	// Giant Swarm objects.
	scheme := pkgruntime.NewScheme()
	err = v1alpha2.AddToScheme(scheme)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	err = v1alpha1.AddToScheme(scheme)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	codecs := serializer.NewCodecFactory(scheme)
	deserializer := codecs.UniversalDeserializer()

	// Do the acutal decoding.
	input, err := pkgruntime.Decode(deserializer, inputData)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return input, nil
}

// outputValue return the expected output by reading the output file.
func (r *Runner) outputValue() ([]byte, error) {
	outputFile := filepath.Join(r.OutputDir, r.files[r.current].Name())
	output, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return output, nil
}

func (r *Runner) Err() error {
	return r.err
}
