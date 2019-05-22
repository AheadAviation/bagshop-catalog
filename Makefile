SHELL := /bin/bash
GO := go
NAME := bagshop-catalog
REPO := aheadaviation
OS := $(shell uname)
MAIN_GO := cmd/catalogsvc/main.go
ROOT_PACKAGE := $(GIT_PROVIDER)/aheadaviation/$(NAME)
GO_VERSION := $(shell $(GO) version | sed -e 's/^[^0-9.]*\([0-9.]*\).*/\1/')
PACKAGE_DIRS := $(shell $(GO) list ./... | grep -v /vendor/)
PKGS := $(shell go list ./... | grep -v /vendor | grep -v generated)
PKGS := $(subst :,_,$(PKGS))
BUILDFLAGS := ''
CGO_ENABLED = 0
VENDOR_DIR=vendor

build:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_GO)

test:
	CGO_ENABLED=$(CGO_ENABLED) $(GO) test $(PACKAGE_DIRS) -test.v -covermode=count -coverprofile=coverage.out -json > report.json

install:
	GOBIN=${GOBIN}/bin $(GO) install -ldflags $(BUILDFLAGS) $(MAIN_GO)

fmt:
	@FORMATTED=`$(GO) fmt $(PACKAGE_DIRS)`
	@([[ ! -z "$(FORMATTED)" ]] && printf "Fixed unformatted files:\n$(FORMATTED)") || true

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) build -ldflags $(BUILDFLAGS) -o bin/$(NAME)-amd64-linux $(MAIN_GO)

.PHONY: release clean

dep: go mod download

docker:
	docker build -t $(REPO)/$(NAME) .
	docker push $(REPO)/$(NAME)

rm-dep:
	rm -rf vendor

cover:
	go test -covermode=count ./...

coverprofile:
	go get github.com/modocache/gover
	go test -v -covermode=count -coverprofile=profile.coverprofile
	go test -v -covermode=count -coverprofile=db.coverprofile ./db
	go test -v -covermode=count -coverprofile=mongo.coverprofile ./db/mysql
	go test -v -covermode=count -coverprofile=api.coverprofile ./api
	go test -v -covermode=count -coverprofile=users.coverprofile ./item
	gover
	mv gover.coverprofile cover.profile
	rm *.coverprofile

clean:
	rm -rf build release

FGT := $(GOPATH)/bin/fgt
$(FGT):
	go get github.com/GeertJohan/fgt

GOLINT := $(GOPATH)/bin/golint
$(GOLINT):
	go get github.com/golang/lint/golint

$(PKGS): $(GOLINT) $(FGT)
	@echo "LINTING"
	@$(FGT) $(GOLINT) $(GOPATH)/src/$@/*.go
	@echo "VETTING"
	@go vet -v $@
	@echo "TESTING"
	@go test -v $@

.PHONY: lint
lint: vendor | $(PKGS) $(GOLINT) # ❷
	@cd $(BASE) && ret=0 && for pkg in $(PKGS); do \
		test -z "$$($(GOLINT) $$pkg | tee /dev/stderr)" || ret=1 ; \
	done ; exit $$ret
