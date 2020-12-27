PREJECT?=github.com/kaatinga/dockerhomework1

GOOS?=darwin
GOARCH?=amd64

RELEASE?=v1.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y +%m +%d-%H:%M:%S')

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build \
	-ldflags="-s -w -X ${PROJECT}/version.Version=${VERSION} \
	-X ${PROJECT}/version.Commit=${GIT_COMMIT} \
	-X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
	-o /app .

