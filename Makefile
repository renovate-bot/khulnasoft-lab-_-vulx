VERSION := $(shell git describe --tags)
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION)"

GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GOSRC=$(GOPATH)/src

MKDOCS_IMAGE := khulnasoft/mkdocs-material:dev
MKDOCS_PORT := 8000

u := $(if $(update),-u)

$(GOBIN)/wire:
	GO111MODULE=off go get github.com/google/wire/cmd/wire

.PHONY: wire
wire: $(GOBIN)/wire
	wire gen ./pkg/...

.PHONY: mock
mock: $(GOBIN)/mockery
	mockery -all -inpkg -case=snake -dir $(DIR)

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

$(GOBIN)/golangci-lint:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOBIN) v1.41.1

.PHONY: test
test:
	go test -v -short -coverprofile=coverage.txt -covermode=atomic ./...

integration/testdata/fixtures/images/*.tar.gz:
	git clone https://github.com/khulnasoft-lab/vul-test-images.git integration/testdata/fixtures/images

.PHONY: test-integration
test-integration: integration/testdata/fixtures/images/*.tar.gz
	go test -v -tags=integration ./integration/...

.PHONY: lint
lint: $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run --timeout 5m

.PHONY: fmt
fmt:
	find ./ -name "*.proto" | xargs clang-format -i

.PHONY: build
build:
	go build $(LDFLAGS) ./cmd/vul

.PHONY: protoc
protoc:
	find ./rpc/ -name "*.proto" -type f -exec protoc --proto_path=$(GOSRC):. --twirp_out=. --twirp_opt=paths=source_relative --go_out=. --go_opt=paths=source_relative {} \;

.PHONY: install
install:
	go install $(LDFLAGS) ./cmd/vul

.PHONY: clean
clean:
	rm -rf integration/testdata/fixtures/images

$(GOBIN)/labeler:
	go install github.com/knqyf263/labeler@latest

.PHONY: label
label: $(GOBIN)/labeler
	labeler apply misc/triage/labels.yaml -r khulnasoft-lab/vul -l 5

.PHONY: mkdocs-serve
## Runs MkDocs development server to preview the documentation page
mkdocs-serve:
	docker build -t $(MKDOCS_IMAGE) -f docs/build/Dockerfile docs/build
	docker run --name mkdocs-serve --rm -v $(PWD):/docs -p $(MKDOCS_PORT):8000 $(MKDOCS_IMAGE)
