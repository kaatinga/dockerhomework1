PREJECT?=github.com/kaatinga/dockerhomework1

GOOS?=linux
GOARCH?=amd64

RELEASE?=1.0.0
COMMIT := git-$(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u '+%Y +%m +%d-%H:%M:%S')

build:
	GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build \
	-ldflags="-X 'main.Version=${RELEASE}' \
	-X 'main.Commit=${COMMIT}' \
	-X 'main.BuildTime=${BUILD_TIME}'" \
	-o /app .

