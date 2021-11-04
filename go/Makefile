# TODO: take from ../configure
INSTALLDIR=/usr/local/bin

build:
	go build

check:
	# Unit tests (small number)
	go test -v mlr/src/...
	# Regression tests (large number)
	go test

install: build
	cp mlr $(INSTALLDIR)

.PHONY: build

