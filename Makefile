VERSION := 0.1.5
TAR_FILE := retrier-v$(VERSION)-darwin-arm.tar.gz
SHA256 := $(shell shasum -a 256 $(TAR_FILE) | awk '{print $$1}')

.PHONY: all build tar bump-patch update-sha

all: bump-patch build tar update-sha

build:
    go build -o retrier .

tar:
    @echo "Creating tarball for version $(VERSION)"
    tar -czvf $(TAR_FILE) retrier

bump-patch:
    rm -rf $(TAR_FILE)
    @NEW_VERSION=$(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}') && \
    sed -i '' "s/^VERSION := .*/VERSION := $${NEW_VERSION}/" Makefile && \
    sed -i '' "s/$(VERSION)/$${NEW_VERSION}/g" main.go retrier.rb && \
    $(eval VERSION := $${NEW_VERSION}) && \
    $(eval TAR_FILE := retrier-v$${NEW_VERSION}-darwin-arm.tar.gz) && \
    echo "Version bumped to $${NEW_VERSION}"

update-sha: tar
    @NEW_SHA256=$(SHA256) && \
    sed -i '' "s/sha256 \".*\"/sha256 \"$${NEW_SHA256}\"/" retrier.rb && \
    echo "SHA256 updated to $${NEW_SHA256}"