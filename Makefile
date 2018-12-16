WHAT := loki_exporter

PWD ?= $(shell pwd)

VERSION   ?= $(shell git describe --tags)
REVISION  ?= $(shell git rev-parse HEAD)
BRANCH    ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDUSER ?= $(shell id -un)
BUILDTIME ?= $(shell date '+%Y%m%d-%H:%M:%S')

DOCKER_REPO       ?= ricoberger
DOCKER_IMAGE_NAME ?= loki_exporter
DOCKER_IMAGE_TAG  ?= $(shell git describe --tags)

.PHONY: build build-darwin-amd64 build-linux-amd64 build-windows-amd64 clean release docker docker-publish docker-tag-latest

build:
	for target in $(WHAT); do \
		go build -ldflags "-X github.com/prometheus/common/version.Version=${VERSION} \
			-X github.com/prometheus/common/version.Revision=${REVISION} \
			-X github.com/prometheus/common/version.Branch=${BRANCH} \
			-X github.com/prometheus/common/version.BuildUser=${BUILDUSER} \
			-X github.com/prometheus/common/version.BuildDate=${BUILDTIME}" \
			-o ./bin/$$target ./cmd/$$target; \
	done

build-darwin-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -a -installsuffix cgo -ldflags "-X github.com/prometheus/common/version.Version=${VERSION} \
			-X github.com/prometheus/common/version.Revision=${REVISION} \
			-X github.com/prometheus/common/version.Branch=${BRANCH} \
			-X github.com/prometheus/common/version.BuildUser=${BUILDUSER} \
			-X github.com/prometheus/common/version.BuildDate=${BUILDTIME}" \
			-o ./bin/loki_exporter-${VERSION}-darwin-amd64/$$target ./cmd/$$target; \
	done

build-linux-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -ldflags "-X github.com/prometheus/common/version.Version=${VERSION} \
			-X github.com/prometheus/common/version.Revision=${REVISION} \
			-X github.com/prometheus/common/version.Branch=${BRANCH} \
			-X github.com/prometheus/common/version.BuildUser=${BUILDUSER} \
			-X github.com/prometheus/common/version.BuildDate=${BUILDTIME}" \
			-o ./bin/loki_exporter-${VERSION}-linux-amd64/$$target ./cmd/$$target; \
	done

build-windows-amd64:
	for target in $(WHAT); do \
		CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -a -installsuffix cgo -ldflags "-X github.com/prometheus/common/version.Version=${VERSION} \
			-X github.com/prometheus/common/version.Revision=${REVISION} \
			-X github.com/prometheus/common/version.Branch=${BRANCH} \
			-X github.com/prometheus/common/version.BuildUser=${BUILDUSER} \
			-X github.com/prometheus/common/version.BuildDate=${BUILDTIME}" \
			-o ./bin/loki_exporter-${VERSION}-windows-amd64/$$target.exe ./cmd/$$target; \
	done

clean:
	for target in $(WHAT); do \
		rm -rf ./bin/$$target*; \
	done

release: clean build-darwin-amd64 build-linux-amd64 build-windows-amd64
	cp ${PWD}/LICENSE ${PWD}/bin/loki_exporter-${VERSION}-darwin-amd64
	cp ${PWD}/LICENSE ${PWD}/bin/loki_exporter-${VERSION}-linux-amd64
	cp ${PWD}/LICENSE ${PWD}/bin/loki_exporter-${VERSION}-windows-amd64
	cp ${PWD}/config.yml ${PWD}/bin/loki_exporter-${VERSION}-darwin-amd64
	cp ${PWD}/config.yml ${PWD}/bin/loki_exporter-${VERSION}-linux-amd64
	cp ${PWD}/config.yml ${PWD}/bin/loki_exporter-${VERSION}-windows-amd64
	cd ${PWD}/bin; tar cfvz loki_exporter-${VERSION}-darwin-amd64.tar.gz ./loki_exporter-${VERSION}-darwin-amd64
	cd ${PWD}/bin; tar cfvz loki_exporter-${VERSION}-linux-amd64.tar.gz ./loki_exporter-${VERSION}-linux-amd64
	cd ${PWD}/bin; tar cfvz loki_exporter-${VERSION}-windows-amd64.tar.gz ./loki_exporter-${VERSION}-windows-amd64

docker: clean build-linux-amd64
	docker build -t "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

docker-publish:
	docker push "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME)"

docker-tag-latest:
	docker tag "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" "$(DOCKER_REPO)/$(DOCKER_IMAGE_NAME):latest"
