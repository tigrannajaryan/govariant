PKGS=./...

default: test

ci: test benchmark

test:
	$(MAKE) test-arch GOARCH=amd64
	$(MAKE) test-arch GOARCH=386

test-arch:
	go test -v ./...

benchmark:
	$(MAKE) benchmark-arch GOARCH=amd64
	$(MAKE) benchmark-arch GOARCH=386

benchmark-arch:
	go test -bench . -benchmem $(PKGS)
