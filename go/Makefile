# TODO: take from ../configure
INSTALLDIR=/usr/local/bin

build:
	go build

check:
	go test -v mlr/src/...
	./mlr regtest -s 5

install: build
	cp mlr $(INSTALLDIR)

.PHONY: build

