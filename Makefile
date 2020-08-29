GIT ?= git
DOCKER_IMAGE := highwayproject/highway
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")


.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o bin/highway ./cmd/highway/main.go

.PHONY: run
run:
	./bin/highway

.PHONY: test
test:
	go test -cover -race -v ./...

.PHONY: docker-build
docker-build:
	docker build -f Dockerfile -t ${DOCKER_IMAGE}:${VERSION} .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 -v $(shell pwd)/config.yml:/config.yml ${DOCKER_IMAGE}:${VERSION}

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION)

.PHONY: docker-push-latest
docker-push-latest:
	docker rmi ${DOCKER_IMAGE}:latest || true
	docker tag ${DOCKER_IMAGE}:${VERSION} ${DOCKER_IMAGE}:latest
	docker push ${DOCKER_IMAGE}:latest
