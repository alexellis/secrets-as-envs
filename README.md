## secrets-as-envs

Load a list of pre-defined secrets from openfaas files into environment variables.

Running the tests:

```
cd secrets-env

OPENFAAS_SECRETS_PATH=./test-files go test -v ./
```

Using the example in your own code:

First off, create the secrets that you need:

```
faas-cli secret create s3-access-key --from-literal mykeyid
faas-cli secret create s3-secret-key --from-literal myaccesskey
```

Then reference them in stack.yml:

```
  secrets-env:
    lang: golang-middleware
    handler: ./secrets-env
    image: alexellis2/secrets-env:0.1.3
    environment:
      s3-access-key: AWS_ACCESS_KEY_ID
      s3-secret-key: AWS_SECRET_ACCESS_KEY
    secrets:
      - s3-access-key
      - s3-secret-key
```

For each file you want to read and map into an environment variable, add an environment variable with the secret filename on the left, and the destination environment variable on the right.

```
s3-access-key: AWS_ACCESS_KEY_ID
```

Here, the textual value from the secret `s3-access-key` will be read and set into the environment variable `AWS_ACCESS_KEY_ID`.

Finally, copy in the code from handler.go.

In your Handle method, you can now access any mapped secret as an environment variable:

```
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
