package function

import (
	"os"
	"path"
	"testing"
)

func Test_Handler(t *testing.T) {
	if _, ok := os.LookupEnv("OPENFAAS_SECRETS_PATH"); !ok {
		return
	}

	secretsDir := os.Getenv("OPENFAAS_SECRETS_PATH")

	secretFile := path.Join(secretsDir, "s3-access-key")
	data, err := os.ReadFile(secretFile)
	if err != nil {
		t.Fatal(secretFile, err)
	}

	os.Setenv("s3-access-key", "AWS_ACCESS_KEY_ID")

	if err := setSecretsAsEnvs(); err != nil {
		t.Error(err)
	}

	v := os.Getenv("AWS_ACCESS_KEY_ID")

	if len(v) == 0 {
		t.Error("AWS_ACCESS_KEY_ID is not set")
	}

	if v != string(data) {
		t.Fatal("Expected a value in AWS_ACCESS_KEY_ID")
	}
}
