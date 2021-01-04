#PREJECT?=github.com/kaatinga/dockerhomework1

GOOS?=linux
GOARCH?=amd64

RELEASE?=1.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y +%m +%d-%H:%M:%S')
LDFLAGS = -ldflags="-X 'main.Version=${RELEASE}' -X 'main.Commit=${COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'"

.PHONY : build run test

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build ${LDFLAGS} -o /app .

test:
	go test

run: test
	PORT=8080 go run ${LDFLAGS} .

