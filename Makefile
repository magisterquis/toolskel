# Makefile
# Build toolskel
# By J. Stuart McMurray
# Created 20250202
# Last Modified 20260407

BINNAME       != basename $$(pwd)
GOBUILDFLAGS   = -trimpath -ldflags "-w -s"
GOTESTFLAGS   += -timeout 3s
SHMORESUBR     = t/shmore.subr
SHMOREURL      = https://raw.githubusercontent.com/magisterquis/shmore/refs/heads/master/shmore.subr

all: test build ## Build ALL the things (default)
.PHONY: all

${BINNAME}:
	go build ${GOBUILDFLAGS} -o ${BINNAME}

build: ${BINNAME}
.PHONY: build

test: gotest provetest ## Run ALL the tests
.PHONY: test

gotest: ## Run go-specific tests
	go test ${GOBUILDFLAGS} ${GOTESTFLAGS} ./...
	go vet ${GOBUILDFLAGS} ./...
	! which staticcheck >/dev/null || staticcheck ./...
	[[ -z "$$(go fix -diff ./... | tee /dev/stderr)" ]]
	go run ${GOBUILDFLAGS} . -h 2>&1 |\
	awk '\
		/^Options:$$|MQD DEBUG PACKAGE LOADED$$/\
			{ exit }\
		/^Usage: /\
			{ sub(/^Usage: [^[:space:]]+\//, "Usage: ") }\
		/.{80,}/\
			{ print "Long usage line: " $$0; exit 1 }\
	'
.PHONY: gotest

provetest: ## Run tests with prove(1) if ./t exists
.if exists(./t/)
	prove -It --directives
.endif
.PHONY: provetest

update: ## Fetch the latest Shmore and up-to-date Go things
	curl\
		--fail\
		--show-error\
		--silent\
		--output ${SHMORESUBR}.new\
		${SHMOREURL}
	diff -q ${SHMORESUBR} ${SHMORESUBR}.new >/dev/null &&\
		rm ${SHMORESUBR}.new ||\
		mv ${SHMORESUBR}.new ${SHMORESUBR}
	go get -t -u go ./...
	go mod tidy
.PHONY: update

install: ## Install to GOBIN ($GOPATH/bin or $HOME/go/bin)
	go install ${GOBUILDFLAGS}
.PHONY: install

clean: ## Remove built things
	rm -rf ${BINNAME}
.PHONY: clean

help: .NOTMAIN ## This help
	@perl -ne '/^(\S+?):+.*?##\s*(.*)/&&print"$$1\t-\t$$2\n"' \
		${MAKEFILE_LIST} | column -ts "$$(printf "\t")"
.PHONY: help
