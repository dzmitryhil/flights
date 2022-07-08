PACKAGES=$(shell go list ./... )
VERSION := $(shell git describe --abbrev=6 --dirty --always --tags)
COMMIT := $(shell git log -1 --format='%H')
IMPORT_PREFIX=github.com/dzmitryhil
SCAN_FILES := $(shell find . -type f -name '*.go' -not -name '*mock.go' -not -name '*_gen.go' -not -path "*/vendor/*")

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

.PHONY: all
all: lint test install

.PHONY: build
build: go.sum
		go build $(BUILD_FLAGS) ./cmd/flightsd

.PHONY: install
install: go.sum
		go install $(BUILD_FLAGS) ./cmd/flightsd

.PHONY: go.sum
go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

.PHONY: test
test:
	@go test -mod=readonly $(PACKAGES)

.PHONY: lint
lint:
	golangci-lint -c .golangci.yml run
	gofmt -d -s $(SCAN_FILES)

.PHONY: format
format:
	gofumpt -lang=1.6 -extra -s -w $(SCAN_FILES)
	gogroup -order std,other,prefix=$(IMPORT_PREFIX) -rewrite $(SCAN_FILES)

.PHONY: swagger-generate
swagger-generate:
	swagger generate spec -o ./docs/openapi/gen.yaml
	swagger mixin ./docs/openapi/rules.yaml ./docs/openapi/gen.yaml -o ./docs/openapi/openapi.yml --format=yaml
	rm docs/openapi/gen.yaml

###############################################################################
###                      Docker wrapped commands                            ###
###############################################################################

.PHONY: in-docker
in-docker:
	docker build -t flights-dev-utils ./dev/tools -f dev/tools/devtools.Dockerfile
	docker run -i --rm \
		-v ${PWD}:/go/src/github.com/dzmitryhil/flights:delegated \
		--mount source=dev-tools-cache,destination=/cache/,consistency=delegated flights-dev-utils bash -x -c "\
		$(ARGS)"

.PHONY: lint-in-docker
lint-in-docker:
	make in-docker ARGS="make lint"

.PHONY: format-in-docker
format-in-docker:
	make in-docker ARGS="make format"

.PHONY: all-in-docker
all-in-docker:
	make in-docker ARGS="make all"

.PHONY: swagger-generate-in-docker
swagger-generate-in-docker:
	touch docs/openapi/openapi.yml # to preserve files ownership
	make in-docker ARGS="make swagger-generate"
