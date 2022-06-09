package function

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

func init() {
	if err := setSecretsAsEnvs(); err != nil {
		log.Panic(err)
	}
}

func setSecretsAsEnvs() error {
	secretsDir := "/var/openfaas/secrets/"
	if v, ok := os.LookupEnv("OPENFAAS_SECRETS_PATH"); ok && len(v) > 0 {
		secretsDir = v
	}
	if _, err := os.Stat(secretsDir); err != nil && os.IsNotExist(err) {
		log.Printf("Secrets directory %s does not exist", secretsDir)
		return nil
	}

	s, err := os.ReadDir(secretsDir)
	if err != nil {
		return fmt.Errorf("unable to read secrets directory: %s, error: %s", secretsDir, err)
	}

	for _, f := range s {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), "..data") {
			secret := path.Join(secretsDir, f.Name())
			body, err := os.ReadFile(secret)
			if err != nil {
				return fmt.Errorf("unable to read secret file: %s, error: %s", secret, err)
			}
			if envName, ok := os.LookupEnv(f.Name()); ok && len(envName) > 0 {
				os.Setenv(envName, string(body))
			} else {
				log.Printf("Secret found without environment variable mapping: %s\n", f.Name())
			}
		}
	}

	return nil
}

// Example usage:

// faas-cli secret create s3-access-key --from-literal mykeyid
// faas-cli secret create s3-secret-key --from-literal myaccesskey

// environment:
//   s3-access-key: AWS_ACCESS_KEY_ID
//   s3-secret-key: AWS_SECRET_ACCESS_KEY

// Then read:
// os.Getenv("AWS_ACCESS_KEY_ID")
// os.Getenv("AWS_SECRET_ACCESS_KEY")
func Handle(w http.ResponseWriter, r *http.Request) {

	msg := fmt.Sprintf("AWS secrets from environment: %s, %s\n",
		os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"))

	log.Println(msg)
	w.Write([]byte(msg))
}
