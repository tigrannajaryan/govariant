PKGS=./...

.PHONY: default
default: test

.PHONY: ci
ci: test benchmark

.PHONY: test
test:
	$(MAKE) test-arch GOARCH=amd64
	$(MAKE) test-arch GOARCH=386

.PHONY: test-arch
test-arch:
	@echo ============================== Testing GOARCH=$(GOARCH) ==============================
	go test -v ./...

.PHONY: benchmark
benchmark:
	$(MAKE) benchmark-arch GOARCH=amd64
	$(MAKE) benchmark-arch GOARCH=386

.PHONY: benchmark-arch
benchmark-arch:
	@echo ============================== Benchmarking GOARCH=$(GOARCH) =========================
	go test -bench . -benchmem $(PKGS) $(BENCHARGS) | tee benchmark/benchmark.log
	sed "s/BenchmarkVariant//" benchmark/benchmark.log | sed "s+ns/op++" | sed "s+\(pkg: github\.com/tigrannajaryan/govariant/\)\(.*\)+\1\2\t\t\2+" | sed "s+B/op++" > benchmark/benchmark$(GOARCH).log
