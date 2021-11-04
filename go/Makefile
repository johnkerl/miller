# TODO: take from ../configure
INSTALLDIR=/usr/local/bin

build: .always
	go build

check: build
	go test -v mlr/src/...
	./mlr regtest -s 5

install: build
	cp mlr $(INSTALLDIR)

.always:
	true
