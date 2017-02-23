NAME     := ct2stimer
VERSION  := v0.2.0
REVISION := $(shell git rev-parse --short HEAD)

SRCS      := $(shell find . -type f -name '*.go')
TEMPLATES := $(shell find . -type f -name '*.tmpl')
LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

DIST_DIRS := find * -type d -exec

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS) $(TEMPLATES)
	$(MAKE) generate
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: ci-test
ci-test:
	echo "" > coverage.txt
	for d in `glide novendor`; do \
		go test -coverprofile=profile.out -covermode=atomic -v $$d || break;  \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi; \
	done

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf vendor/*

.PHONY: cross-build
cross-build: generate
	for os in linux; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: deps
deps: glide
	go get -u github.com/jteeuwen/go-bindata/...
	glide install

.PHONY: dist
dist:
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf $(NAME)-$(VERSION)-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r $(NAME)-$(VERSION)-{}.zip {} \; && \
	cd ..

.PHONY: generate
generate: $(TEMPLATES)
	go generate -x ./...

.PHONY: glide
glide:
ifeq ($(shell command -v glide 2> /dev/null),)
	curl https://glide.sh/get | sh
endif

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test: generate
	go test -cover -race -v `glide novendor`

.PHONY: update-deps
update-deps: glide
	glide update
