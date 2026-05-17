.PHONY: init,proto,build,run,stop,clean,test,lint,fmt,ent,docker,deploy

GOPATH := $(shell go env GOPATH)
PROJECT_NAME := feishu-office-suite
VERSION := v1.0.0
BUILD_TIME := $(shell date -u)
GO_VERSION := $(shell go version | awk '{print $$3}')

API_DIR := api
PROTO_FILES := $(shell find $(API_DIR) -name "*.proto")
GO_PROTO_FILES := $(patsubst %.proto,%.pb.go,$(PROTO_FILES))

DOCKER_IMAGE := $(PROJECT_NAME):$(VERSION)
K8S_NAMESPACE := feishu-suite

init:
	@echo "Initializing Go modules..."
	go mod tidy
	@echo "Installing dependencies..."
	@which protoc || (echo "Please install protoc: https://grpc.io/docs/protoc-installation/" && exit 1)
	@which protoc-gen-go || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@which protoc-gen-go-http || go install github.com/go-kratos/grpc-gateway/v2/cmd/protoc-gen-go-http/v2@latest
	@which protoc-gen-kratos || go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	@which protoc-gen-openapi || go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	@which ent || go ent/cmd/ent/...@latest
	@which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Init completed!"

proto: $(GO_PROTO_FILES)

%.pb.go: %.proto
	@echo "Generating protobuf code for $<..."
	protoc --proto_path=. \
		--proto_path=$(GOPATH)/src \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		$<

ent:
	@echo "Generating Ent code..."
	cd internal/data && ent generate ./ent/schema

build:
	@echo "Building $(PROJECT_NAME)..."
	go build -o bin/server ./cmd/server
	go build -o bin/worker ./cmd/worker

run:
	@echo "Running $(PROJECT_NAME)..."
	docker-compose -f deployments/docker/docker-compose.yaml up -d
	./bin/server -conf ./configs

stop:
	docker-compose -f deployments/docker/docker-compose.yaml down

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf internal/data/ent/gen/
	find . -name "*.pb.go" -delete

test:
	@echo "Running tests..."
	go test -v -cover ./...

lint:
	@echo "Running linter..."
	golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...
	imports_reorder -w .

docker:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -f deployments/docker/Dockerfile -t $(DOCKER_IMAGE) .

docker-push:
	docker push $(DOCKER_IMAGE)

k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/k8s/namespace.yaml
	kubectl apply -f deployments/k8s/configmap.yaml
	kubectl apply -f deployments/k8s/secret.yaml
	kubectl apply -f deployments/k8s/deployment.yaml
	kubectl apply -f deployments/k8s/service.yaml

k8s-delete:
	@echo "Deleting from Kubernetes..."
	kubectl delete -f deployments/k8s/

help:
	@echo "Available targets:"
	@echo "  init         - Initialize project and install dependencies"
	@echo "  proto        - Generate protobuf code"
	@echo "  ent          - Generate Ent code"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application with Docker Compose"
	@echo "  stop         - Stop Docker Compose"
	@echo "  clean        - Clean generated files"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  docker       - Build Docker image"
	@echo "  docker-push  - Push Docker image to registry"
	@echo "  k8s-deploy   - Deploy to Kubernetes"
	@echo "  k8s-delete   - Delete from Kubernetes"