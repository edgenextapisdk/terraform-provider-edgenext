TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=edgenext
CHANGED_FILES=$$(git diff --name-only master -- $(PKG_NAME) | grep '.go$$')
PLATFORMS=darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm openbsd/amd64 openbsd/386 solaris/amd64 windows/386 windows/amd64
GO_VER ?= go

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	goimports -w ./$(PKG_NAME)
	gofmt -s -w ./$(PKG_NAME)

fmt-faster:
	@if [[ -z $(CHANGED_FILES) ]]; then \
		echo "skip the fmt cause the CHANGED_FILES is null."; \
		exit 0; \
	else \
		echo "==> [Faster]Fixing source code with gofmt...\n $(CHANGED_FILES) \n"; \
		goimports -w $(CHANGED_FILES); \
		gofmt -s -w $(CHANGED_FILES); \
	fi

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run --timeout=30m ./$(PKG_NAME)

tools:
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

test-build-all:
	@$(foreach platform, $(PLATFORMS), \
		echo GOOS=$(firstword $(subst /, ,$(platform))) GOARCH=$(lastword $(subst /, ,$(platform))) \
		go build -o terraform-provider-edgenext; \
		GOOS=$(firstword $(subst /, ,$(platform))) GOARCH=$(lastword $(subst /, ,$(platform))) \
		go build -o terraform-provider-edgenext; \
		rm -f terraform-provider-edgenext; \
	)

test-build: test-build-x64 test-build-x86

test-build-x64:
	GOARCH=amd64 go build -o terraform-provider-edgenext-amd64
	rm -f terraform-provider-edgenext-amd64

test-build-x86:
	GOARCH=386 go build -o terraform-provider-edgenext-386
	rm -f terraform-provider-edgenext-386

doc:
	cd gendoc && go run ./... -link-format=terraform && cd ..

doc-github:
	cd gendoc && go run ./... -link-format=github && cd ..

doc-faster:
	@echo "==> [Faster]Generating doc..."
	@if [ ! -f gendoc/gendoc ]; then \
		$(MAKE) doc-bin-build; \
	fi
	@$(MAKE) doc-with-bin

doc-with-bin:
	cd gendoc && ./gendoc ./... && cd ..

doc-bin-build:
	@echo "==> Building gendoc binary..."
	cd gendoc && go build ./... && cd ..

hooks: tools
	@find .git/hooks -type l -exec rm {} \;
	@find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;
	@echo "==> Install hooks done."

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/ || (echo; \
		echo "Unexpected mispelling found in website files."; \
		echo "To automatically fix the misspelling, run 'make website-lint-fix' and commit the changes."; \
		exit 1)

website-lint-fix:
	@echo "==> Applying automatic website linter fixes..."
	@misspell -w -source=text website/

ready: doc fmt-faster

.PHONY: build sweep test testacc fmt fmtcheck lint tools test-compile doc doc-github hooks website-lint website-lint-fix ready
