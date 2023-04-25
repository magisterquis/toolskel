# Makefile
# Build toolskel
# By J. Stuart McMurray
# Created 20230415
# Last Modified 20230425

.PHONY: all test vet staticcheck

PACKAGES=. ./internal/gencode

all: test vet staticcheck 

test:
	go test -tags testgoimports ${PACKAGES}

vet:
	go vet ${PACKAGES}

staticcheck:
	staticcheck ${PACKAGES}
