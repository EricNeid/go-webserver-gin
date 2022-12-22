# About

A simple extendable server to access with go and gin.

## Quickstart

The simples way to compile this application is to use the provide makefile.
It provides cross compilation to linux and windows and makes use of docker.

Docker:

```bash
make build-windows
make build-linux
```

Manual and without docker:

```bash
go build -o ./out/ ./cmd/webserver/
```

## Options

Application can be configure using command line arguments or
environment variables or a combination of both.

* listen-addr/LISTEN_ADDR - listing address, ie. ":5000"
* base-path/BASE_PATH - base path to serve application, ie "/custom"
* serve-static/SERVE-STATIC - folder containing static html and the path to serve it, ie. "public=>/dashboard"

Example:

```bash
./mapprovider -base-path mapprovider-0.1.0 -listen-addr :8080
```

## Endpoints

The following endpoints are available, assuming
the configuration is not changed.

* <http://localhost:5000>
* <http://localhost:5000/logs>

## Testing

To run tests:

```bash
make test
```
