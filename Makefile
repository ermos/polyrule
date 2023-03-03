build/bin:
	@go build -o ${GOPATH}/bin/polyrule ./

build: test
	goreleaser --snapshot --skip-publish --rm-dist

release:
	test $(VERSION)
	git tag -a v$(VERSION) -m "$(VERSION)"
	goreleaser --snapshot --skip-publish --rm-dist

delete-tag:
	test $(VERSION)
	git tag --delete v$(VERSION)