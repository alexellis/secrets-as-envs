# secrets-as-envs

The example in this repository shows how to map OpenFaaS secret files into environment variables.

> Disclaimer: OpenFaaS Ltd does not recommend using environment variables to store or inject secrets.

## The preferred way to handle secrets

The design of OpenFaaS supports consuming secrets as files, which is the preferred and recommended way to work with confidential configuration.

Why do you want to inject environment variables as secrets?

If you are wanting to inject secrets via environment variables in order to use the AWS SDK, try using this approach instead:

```yaml
  environment:
      AWS_SHARED_CREDENTIALS_FILE: "/var/openfaas/secrets/aws-ses-credentials"
```

If you're using the Datadog agent, and need a Node IP then you should try the operator:

See: [How to get DD_AGENT_HOST for functions for the Datadog agent?](https://github.com/openfaas/faas-netes/issues/913)

Is it because you're used to this feature from the Pod spec in Kubernetes? [Using Secrets as environment variables](https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets-as-environment-variables)

OpenFaaS uses files, which you can read with os.ReadFile in init(), or whenever you need them in Handler() using `os.ReadFile("/var/openfaas/secret/NAME")`:

```golang
func init() {
	file, err := ioutil.ReadFile("/var/openfaas/secrets/key")

	if err != nil {
		return err.Error()
	}
	log.Printf("value: %s", string(file))
}
```

## The technique

A `setSecretsAsEnvs` function is called via `init()` before any other code loads for the function's handler. It scans a set folder such as /var/openfaas/secrets and reads each file, mapping the contents of the file into an environment variable using `os.Setenv`.

## Use the example in your own functions

First off, create the secrets that you need:

```bash
faas-cli secret create s3-access-key \
  --from-literal mykeyid
faas-cli secret create s3-secret-key \
  --from-literal myaccesskey
```

Then reference them in stack.yml:

```yaml
  secrets-env:
    lang: golang-middleware
    handler: ./secrets-env
    image: alexellis2/secrets-env:0.1.3
    secrets:
      - s3-access-key
      - s3-secret-key
```

For each file you want to read and map into an environment variable, add an environment variable with the secret filename on the left, and the destination environment variable on the right.

```yaml
    environment:
      s3-access-key: AWS_ACCESS_KEY_ID
      s3-secret-key: AWS_SECRET_ACCESS_KEY
```

In the first example, we see that the value read from the `s3-access-key` file will be read and set into the environment variable `AWS_ACCESS_KEY_ID`.

Finally, copy in the code from handler.go.

In your Handle method, you can now access any mapped secret as an environment variable:

```golang
// Then read:
// os.Getenv("AWS_ACCESS_KEY_ID")
// os.Getenv("AWS_SECRET_ACCESS_KEY")
func Handle(w http.ResponseWriter, r *http.Request) {

	msg := fmt.Sprintf("AWS secrets from environment: %s, %s\n",
		os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"))

	log.Println(msg)
	w.Write([]byte(msg))
}
```

## Run the tests

The test use canned data in the "test-files" folder.

```bash
cd secrets-env

OPENFAAS_SECRETS_PATH=./test-files go test -v ./
```

For local development, you can also set the `OPENFAAS_SECRETS_PATH` folder to another path.

