package remotewriteagentsecret

import (
	"context"
	"flag"
	"path/filepath"
	"testing"

	"github.com/giantswarm/prometheus-meta-operator/v2/pkg/unittest"
)

var update = flag.Bool("update", false, "update the ouput file")

func TestRemoteWriteSecret(t *testing.T) {
	outputDir, err := filepath.Abs("./test")
	if err != nil {
		t.Fatal(err)
	}

	c := unittest.Config{
		OutputDir: outputDir,
		T:         t,
		TestFunc: func(v interface{}) (interface{}, error) {
			return toSecret(context.TODO(), v, Config{
				PasswordManager: TestPasswordManager{},
			})
		},
		Update: *update,
	}
	runner, err := unittest.NewRunner(c)
	if err != nil {
		t.Fatal(err)
	}

	err = runner.Run()
	if err != nil {
		t.Fatal(err)
	}
}

type TestPasswordManager struct {
}

func (m TestPasswordManager) GeneratePassword(length int) (string, error) {
	return "password", nil
}

func (m TestPasswordManager) Hash(plaintext string) (string, error) {
	return "encrypted", nil
}
