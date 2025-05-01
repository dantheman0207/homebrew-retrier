VERSION := 0.1.18
TAR_FILE := retrier-v$(VERSION)-darwin-arm.tar.gz
SHA256 := $(shell [ -f $(TAR_FILE) ] && shasum -a 256 $(TAR_FILE) | awk '{print $$1}' || echo "")

.PHONY: all build tar bump-patch update-sha

all: bump-patch
	$(MAKE) reexec || $(MAKE) revert-version

reexec: build tar update-sha tag gh-release

build:
	go build -o retrier .

tar:
	@echo "Creating tarball for v$(VERSION)"
	tar -czvf $(TAR_FILE) ./*
	xattr -cr $(TAR_FILE)


tag:
	git add $(TAR_FILE)
	git status
	git commit -a -m "Release $(VERSION)"
	git tag -a v$(VERSION) -m "Release $(VERSION)" -f
	git push origin v$(VERSION) -f
	git push

gh-release:
	gh release create v$(VERSION) --title "Release $(VERSION)" --notes "Release $(VERSION)" --target main || true
	gh release upload v$(VERSION) $(TAR_FILE) retrier --clobber

bump-patch:
	rm -rf $(TAR_FILE)
	@NEW_VERSION=$(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3+1}') && \
	echo "Bumping version from $(VERSION) to $${NEW_VERSION}" && \
	sed -i '' "s/$(VERSION)/$${NEW_VERSION}/g" main.go retrier.rb Makefile && \
	echo "Version bumped to $${NEW_VERSION}"

revert-version:
	rm -rf $(TAR_FILE)
	@NEW_VERSION=$(shell echo $(VERSION) | awk -F. '{print $$1"."$$2"."$$3-1}') && \
	echo "Lowering version from $(VERSION) to $${NEW_VERSION}" && \
	sed -i '' "s/$(VERSION)/$${NEW_VERSION}/g" main.go retrier.rb Makefile && \
	echo "Version lowered to $${NEW_VERSION}"

update-sha:
	@NEW_SHA256=$(shell [ -f $(TAR_FILE) ] && shasum -a 256 $(TAR_FILE) | awk '{print $$1}' || echo "") && \
	if [ -z "$${NEW_SHA256}" ]; then \
        echo "Error: SHA256 checksum is empty. Exiting."; \
        exit 1; \
    fi && \
	sed -i '' "s/sha256 \".*\"/sha256 \"$${NEW_SHA256}\"/" retrier.rb && \
	echo "SHA256 updated to $${NEW_SHA256}"

