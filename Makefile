CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.2
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
MOCKGEN=${BINDIR}/mockgen_${GOVER}

all: format build test lint
.PHONY: all

# ==============================================================================
# Main

build: bindir
	go build -o ${BINDIR}/app github.com/almostinf/order_delivery_service/cmd/app
.PHONY: build

run: swag-v1
	go run github.com/almostinf/order_delivery_service/cmd/app
.PHONY: run

# ==============================================================================
# Tests commands

test:
	go test -v -count=1 ./internal/...
.PHONY: test

test100:
	go test -v -count=100 ./internal/...
.PHONY: test100

# integration-test-local:
# 	go clean -testcache && go test -v ./integration_tests/...
# .PHONY: integration-test-local

integration-test-docker-up:
	docker-compose -f docker-compose.tests.yml up -d
.PHONY: integration-test-docker-up

integration-test-docker-down:
	docker-compose -f docker-compose.tests.yml down --remove-orphans --volumes
	docker rmi test
	docker rmi integration
.PHONY: integration-test-docker-down

race:
	go test -v -race -count=1 ./internal/...
.PHONY: race

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	xdg-open coverage.html
	rm coverage.out
	rm coverage.html
.PHONY: cover

# ==============================================================================
# Tools commands

swag-v1:
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

gen-mocks: install-mockgen
	${MOCKGEN} -source=internal/infrastructure/interfaces/courier.go -destination=internal/mocks/repo/courier_mocks.go
	${MOCKGEN} -source=internal/infrastructure/interfaces/order.go -destination=internal/mocks/repo/order_mocks.go
.PHONY: generate

install-mockgen: bindir
	test -f ${MOCKGEN} || \
		(GOBIN=${BINDIR} go install github.com/golang/mock/mockgen@v1.6.0 && \
		mv ${BINDIR}/mockgen ${MOCKGEN})
.PHONY: install-mockgen

lint: install-lint
	${LINTBIN} run
.PHONY: install-lint

precommit: format build test lint
	echo "OK"
.PHONY: precommit

bindir:
	mkdir -p ${BINDIR}
.PHONY: bindir

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks
.PHONY: format

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})
.PHONY: install-lint

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})
.PHONY: install-smartimports

# ==============================================================================
# Docker compose commands

compose-up: swag-v1
	docker-compose up -d
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

compose-down-debug:
	docker-compose down --remove-orphans
	docker rmi app
.PHONY: compose-down-debug
