GIT ?= git
DOCKER_IMAGE := highway
COMMIT := $(shell $(GIT) rev-parse HEAD)
VERSION ?= $(shell $(GIT) describe --tags ${COMMIT} 2> /dev/null || echo "$(COMMIT)")


.PHONY: build
build:
	go build -a -o bin/highway ./cmd/highway/main.go

.PHONY: run
run: build
	./bin/highway

.PHONY: test
test:
	go test -cover -race -v ./...

.PHONY: docker-build
docker-build:
	docker build -f Dockerfile -t $(DOCKER_IMAGE):$(VERSION) .;\

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE):$(VERSION);\

.PHONY: docker-push-latest
docker-push-latest:
	docker rmi $(DOCKER_IMAGE):latest || true;\
	docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest;\
	docker push $(DOCKER_IMAGE):latest;\
