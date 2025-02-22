# Makefile
# Build toolskel
# By J. Stuart McMurray
# Created 20250202
# Last Modified 20250222

BINNAME       != basename $$(pwd)
GOBUILDFLAGS   = -trimpath -ldflags "-w -s"
GOTESTFLAGS   += -timeout 3s

.PHONY: all build test gotest provetest help install clean

all: test build ## Build ALL the things (default)

${BINNAME}:
	go build ${GOBUILDFLAGS} -o ${BINNAME}

build: ${BINNAME}

test: gotest provetest ## Run ALL the tests

gotest: ## Run go-specific tests
	go test ${GOBUILDFLAGS} ${GOTESTFLAGS} ./...
	go vet ${GOBUILDFLAGS} ./...
	staticcheck ./...
	go run ${GOBUILDFLAGS} . -h 2>&1 |\
	awk '\
		/^Options:$$|MQD DEBUG PACKAGE LOADED$$/\
			{ exit }\
		/^Usage: /\
			{ sub(/^Usage: [^[:space:]]+\//, "Usage: ") }\
		/.{80,}/\
			{ print "Long usage line: " $0; exit 1 }\
	'

provetest: ## Run tests with prove(1) if ./t exists
.if exists(./t/)
	prove -It --directives
.endif

install: ## Install to GOBIN ($GOPATH/bin or $HOME/go/bin)
	go install ${GOBUILDFLAGS}

clean: ## Remove built things
	rm -rf ${BINNAME}

help: .NOTMAIN ## This help
	@perl -ne '/^(\S+?):+.*?##\s*(.*)/&&print"$$1\t-\t$$2\n"' \
		${MAKEFILE_LIST} | column -ts "$$(printf "\t")"
