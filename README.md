## secrets-as-envs

Load a list of pre-defined secrets from openfaas files into environment variables.

Running the tests:

```
cd secrets-env

OPENFAAS_SECRETS_PATH=./test-files go test -v ./
```