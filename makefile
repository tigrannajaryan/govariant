PKGS=./...

default: test

ci: test benchmark

test:
	$(MAKE) test-arch GOARCH=amd64
	$(MAKE) test-arch GOARCH=386

test-arch:
	@echo ============================== Testing GOARCH=$(GOARCH) ==============================
	go test -v ./...

benchmark:
	$(MAKE) benchmark-arch GOARCH=amd64
	$(MAKE) benchmark-arch GOARCH=386

benchmark-arch:
	@echo ============================== Benchmarking GOARCH=$(GOARCH) =========================
	go test -bench . -benchmem $(PKGS) $(BENCHARGS) | tee benchmark.log
	sed "s/BenchmarkVariant//" benchmark.log | sed "s+ns/op++" | sed "s+\(pkg: github\.com/tigrannajaryan/govariant/\)\(.*\)+\1\2\t\t\2+" | sed "s+B/op++" > benchmark$(GOARCH).log
