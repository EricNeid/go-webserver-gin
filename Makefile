DIR := ${CURDIR}
GO_IMAGE := golang:1.21-alpine
LINTER_IMAGE := golangci/golangci-lint:v1.54-alpine

.PHONY: build-windows
build-windows:
	docker run -it --rm \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/webserver/


.PHONY: build-linux
build-linux:
	docker run -it --rm \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go build -o ./out/ ./cmd/webserver/


.PHONY: test
test:
	docker run -it --rm \
		-w /app -v ${DIR}:/app \
		${GO_IMAGE} \
		go test ./...


.PHONY: test-integration
test-integration:
	go test -tags=integration ./...


.PHONY: lint
lint:
	docker run -it --rm \
		-e CGO_ENABLED=0 \
		-w /app -v ${DIR}:/app \
		${LINTER_IMAGE} \
		golangci-lint run ./...


.PHONY: clean
clean:
	rm -rf out/
	rm -rf logs/
	rm -rf cmd/mapprovider/logs/
