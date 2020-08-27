package unittest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"unicode"

	"github.com/giantswarm/apiextensions/v2/pkg/apis/provider/v1alpha1"
	"github.com/giantswarm/microerror"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/cluster-api/api/v1alpha2"
)

// NormalizeFileName converts all non-digit, non-letter runes in input string to
// dash ('-'). Coalesces multiple dashes into one.
func NormalizeFileName(s string) string {
	var result []rune
	for _, r := range s {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			result = append(result, r)
		} else {
			l := len(result)
			if l > 0 && result[l-1] != '-' {
				result = append(result, rune('-'))
			}
		}
	}
	return string(result)
}

type Runner struct {
	OutputDir string

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

func NewRunner(outputDir string) (*Runner, error) {
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
		OutputDir: outputDir,
		inputDir:  inputDir,
		files:     files,
		current:   -1,
		err:       nil,
	}

	return r, nil
}

func (r *Runner) Next() bool {
	if r.err != nil {
		return false
	}

	r.current++
	if len(r.files) <= r.current {
		return false
	}

	return true
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

func InputFiles() ([]os.FileInfo, string, error) {
	root, err := filepath.Abs("./input")
	if err != nil {
		return nil, "", microerror.Mask(err)
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, "", microerror.Mask(err)
	}

	return files, root, nil
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
