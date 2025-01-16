VERSION = 0.1.0
COMMIT_ID = $(shell git rev-list -1 HEAD)
TESTS := $(wildcard tests/*.sh)
build:
	@go build -ldflags "-X main.version=${VERSION}-${COMMIT_ID}"
	@ln -sf rvld ld
test: build
	$(MAKE) $(TESTS)
	@printf '\e[32mPassed all tests\e[0m\n'
$(TESTS):
	@echo 'Testing' $@
	@./$@
	@printf '\e[32mOK\n\e[0m'
clean:
	go clean
	rm -rf out/
	rm  ld

.PHONY: build clean $(TESTS)