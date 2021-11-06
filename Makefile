# TODO: 'cp go/mlr .' or 'copy go\mlr.exe .' with reliable platform detection
# and no confusing error messages.

build:
	make -C go build
	@echo Miller executable is: go/mlr

check:
	make -C go check

install:
	make -C go install
	make -C man install

# Go does its own dependency management, outside of make.
.PHONY: build
