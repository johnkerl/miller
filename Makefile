PREFIX=/usr/local
INSTALLDIR=$(PREFIX)/bin

build:
	go build

check:
	# Unit tests (small number)
	go test -v mlr/internal/pkg/...
	# Regression tests (large number)
	#
	# See ./regression_test.go for information on how to get more details
	# for debugging.  TL;DR is for CI jobs, we have 'go test -v'; for
	# interactive use, instead of 'go test -v' simply use 'mlr regtest
	# -vvv' or 'mlr regtest -s 20'. See also src/auxents/regtest.
	go test -v

install: build
	cp mlr $(DESTDIR)/$(INSTALLDIR)
	make -C man install

fmt:
	-go fmt ./...

# For developers before pushing to GitHub.
#
# These steps are done in a particular order:
# go:
# * builds the mlr executable
# man:
# * creates manpage mlr.1 and manpage.txt using mlr from the $PATH
# * copies the latter to docs/src
# docs:
# * turns *.md.in into *.md (live code samples), using mlr from the $PATH
# * note the man/manpage.txt becomes some of the HTML content
# * turns *.md into docs/site HTML and CSS files
dev:
	-make fmt
	make build
	make check
	make -C man build
	make -C docs
	@echo DONE

# Keystroke-saver
itso: build check install

# Please see comments in ./create-release-tarball as well as
# https://miller.readthedocs.io/en/latest/build/#creating-a-new-release-for-developers
release_tarball: build check
	./create-release-tarball

# Go does its own dependency management, outside of make.
.PHONY: build check fmt dev
