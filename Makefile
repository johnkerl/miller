build:
	make -C go build

check:
	make -C go check

install:
	make -C go install
	make -C man install

# Go does its own dependency management, outside of make.
.PHONY: build
