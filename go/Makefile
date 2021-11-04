# TODO: take from ../configure
INSTALLDIR=/usr/local/bin

build:
	go build

check:
	# Unit tests (small number)
	go test -v mlr/src/...
	# Regression tests (large number)
	#
	# See ./regression_test.go for information on how to get more details
	# for debugging.  TL;DR is for CI jobs, we have 'go test -v'; for
	# interactive use, instead of 'go test -v' simply use 'mlr regtest
	# -vvv' or 'mlr regtest -s 20'. See also src/auxents/regtest.
	go test -v

install: build
	cp mlr $(INSTALLDIR)

# Go does its own dependency management, outside of make.
.PHONY: build
