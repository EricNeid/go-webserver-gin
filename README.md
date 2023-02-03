<!--
SPDX-FileCopyrightText: 2023 Eric Neidhardt
SPDX-License-Identifier: CC-BY-4.0
-->
<!-- markdownlint-disable MD022 MD032 MD024-->
<!-- markdownlint-disable MD041-->
[![Go Report Card](https://goreportcard.com/badge/github.com/EricNeid/go-webserver-gin?style=flat-square)](https://goreportcard.com/report/github.com/EricNeid/go-webserver-gin)
![Test](https://github.com/EricNeid/go-webserver-gin/actions/workflows/tests.yml/badge.svg)
![Linting](https://github.com/EricNeid/go-webserver-gin/actions/workflows/linting.yml/badge.svg)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/EricNeid/go-webserver-gin)
[![Gitpod Ready-to-Code](https://img.shields.io/badge/Gitpod-Ready--to--Code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/EricNeid/go-webserver-gin)

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
