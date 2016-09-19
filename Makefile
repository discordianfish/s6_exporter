REVISION = $(shell git rev-parse --short HEAD)
VERSION  = $(shell git name-rev --tags --name-only $(REVISION))
BRANCH   = $(shell git symbolic-ref --short -q HEAD)
DATE     = $(shell date +%Y%m%d-%H:%M:%S)
build:
	go build -ldflags "\
		-X github.com/prometheus/common/version.Version=${VERSION} \
		-X github.com/prometheus/common/version.Revision=${REVISION} \
		-X github.com/prometheus/common/version.Branch=${BRANCH} \
		-X github.com/prometheus/common/version.BuildUser=${USER} \
		-X github.com/prometheus/common/version.BuildDate=${DATE} \
	" ./...

.PHONY: build
