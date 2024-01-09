# Docker Compose With Coder: How to

## Config

The default configuration is located in the `.env` file.
I'm working on a `arm64` architecture. So, ff you need to adapt the values of the environment variables, create an new env file (ex: `amd64.env`) and use the `--env-file` flag of Docker Compose.

- Add `.config` to `.gitignore`


## Build
```bash
docker compose --env-file .env build
```

## Run
```bash
# If you use .env
docker compose up -d
# If you use a specific .env file
docker compose --env-file ./amd64.env up  -d
```
Then: 
- Open: http://0.0.0.0:4000
- Open a terminal from the Web IDE
- Type this command `git config --global --add safe.directory /ide.simplism.cloud`
- And configure **git** (if necessary):
  ```bash
  git config --global user.name @your-handle
  git config --global user.email your@e.mail
  ```

### Connect to the container
```bash
set -o allexport; source .env; set +o allexport
docker exec --workdir /${WORKDIR} -it ${CONTAINER_NAME} \
/bin/bash
```

## Stop
```bash
docker compose down
```

## Docker in Docker

If you get this error message: `permission denied while trying to connect to the Docker daemon socket`, use this before : `sudo chmod 666 /var/run/docker.sock`.

## TLS certificates

- You can use TLS certificates
- You can use https://github.com/FiloSottile/mkcert

### Use TLS certificates

Update the `.env` file with the name of the certificate and the key:

For example:
```
TLS_CERT=ide.simplism.cloud.crt
TLS_CERT_KEY=ide.simplism.cloud.key
```

- Copy the certificate and the key to the `./certs` folder
- Give the appropriates rights to the files: `chmod 777 ide.simplism.cloud.*`
- Add this to your hosts file: `0.0.0.0 ide.simplism.cloud`
- Use this entrypoint in the `compose.yaml` file: `entrypoint: ["code-server", "--cert", "/${WORKDIR}/certs/${TLS_CERT}", "--cert-key", "/${WORKDIR}/certs/${TLS_CERT_KEY}", "--auth", "none", "--host", "0.0.0.0", "--port", "${CODER_HTTP_PORT}", "/${WORKDIR}"]`
- Restart: `docker compose up`

open https://ide.simplism.cloud:4000

### Use Mkcert

First install `mkcert` (https://github.com/FiloSottile/mkcert?tab=readme-ov-file#installation)
Then:

```bash
cd certs
mkcert -key-file ide.personal.faas.key -cert-file ide.personal.faas.crt personal.faas "*.personal.faas"
mkcert -install
chmod 777 ide.personal.faas.*
```

- Update the `.env` file with the name of the certificate and the key
- Add this to the `hosts` file: `0.0.0.0 ide.personal.faas`
- Use this entrypoint in the `compose.yaml` file: `entrypoint: ["code-server", "--cert", "/${WORKDIR}/certs/${TLS_CERT}", "--cert-key", "/${WORKDIR}/certs/${TLS_CERT_KEY}", "--auth", "none", "--host", "0.0.0.0", "--port", "${CODER_HTTP_PORT}", "/${WORKDIR}"]`
- Restart: `docker compose up`

open https://ide.personal.faas:4000
