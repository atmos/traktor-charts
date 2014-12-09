.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

# Prepend our _vendor directory to the system GOPATH
# so that import path resolution will prioritize
# our third party snapshots.
GOPATH := ${PWD}/_vendor:${GOPATH}
BREW_PREFIX := $(shell brew --prefix 2>/dev/null)
GO := $(BREW_PREFIX)/bin/go
GOROOT := $(BREW_PREFIX)/Cellar/go/1.3.3/libexec
CGO_CFLAGS := -I$(BREW_PREFIX)/include
CGO_LDFLAGS := -L$(BREW_PREFIX)/lib
export CGO_CFLAGS CGO_LDFLAGS GOROOT GOPATH

default: build

build:
	/usr/bin/env CC=clang \
		$(GO) build -v -o ./bin/traktor-charts \
		./src/db.go ./src/output.go ./src/traktor.go ./src/traktor-charts.go

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	$(GO) fmt ./src/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

run: build
	./bin/traktor-charts

test:
	$(GO) test ./src/...

vendor_clean:
	rm -dRf ./_vendor/src

# We have to set GOPATH to just the _vendor
# directory to ensure that `go get` doesn't
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.
vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor $(GO) get -d -u -v \
	github.com/jpoehls/gophermail \
	github.com/codegangsta/martini \
	github.com/stretchr/testify \
	github.com/mattn/go-sqlite3

vendor_update: vendor_get
	rm -rf `find ./_vendor/src -type d -name .git` \
	&& rm -rf `find ./_vendor/src -type d -name .hg` \
	&& rm -rf `find ./_vendor/src -type d -name .bzr` \
	&& rm -rf `find ./_vendor/src -type d -name .svn`

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	$(GO) get code.google.com/p/go.tools/cmd/vet
	$(GO) vet ./src/...
