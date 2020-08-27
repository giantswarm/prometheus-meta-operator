package unittest

import (
	"bytes"
	"fmt"
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

type Runner struct {
	OutputDir string
	T         *testing.T
	TestFunc  func(interface{}) (metav1.Object, error)

	inputDir string
	files    []os.FileInfo
	current  int
	err      error
}

type Value struct {
	Name   string
	Input  pkgruntime.Object
	Output []byte
}

func NewRunner(config Config) (*Runner, error) {
	_, filename, _, ok := runtime.Caller(0)
	fmt.Println(path.Dir(filename), ok)
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

func (r *Runner) Next() bool {
	if r.err != nil {
		return false
	}

	r.current++
	return len(r.files) > r.current
}

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

func (r *Runner) inputValue() (pkgruntime.Object, error) {
	inputFile := filepath.Join(r.inputDir, r.files[r.current].Name())
	inputData, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return nil, microerror.Mask(err)
	}

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
	input, err := pkgruntime.Decode(deserializer, inputData)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	return input, nil
}

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
