PKGS=./...

.PHONY: default
default: test

.PHONY: ci
ci: test benchmark

.PHONY: test
test:
	$(MAKE) test-arch GOARCH=amd64
	$(MAKE) test-arch GOARCH=386

.PHONY: test-coverage
test-coverage:
	$(MAKE) test-coverage-arch GOARCH=amd64
	$(MAKE) test-coverage-arch GOARCH=386

.PHONY: test-coverage-arch
test-coverage-arch:
	go test -cover -coverprofile=coverage$(GOARCH).out github.com/tigrannajaryan/govariant/variant

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
	sed -f benchmark/patch_results.sed benchmark/benchmark.log > benchmark/benchmark$(GOARCH).log
