PREFIX=/usr/local
INSTALLDIR=$(PREFIX)/bin

build:
	go build github.com/johnkerl/miller/cmd/mlr
mlr:
	go build github.com/johnkerl/miller/cmd/mlr

# For interactive use, 'mlr regtest' offers more options and transparency.
check: unit_test regression_test

# Unit tests (small number)
unit_test:
	go test github.com/johnkerl/miller/internal/pkg/...

# Keystroke-savers
unbackslash_test:
	go test $(ls internal/pkg/lib/*.go|grep -v test) internal/pkg/lib/unbackslash_test.go
mlrval_functions_test:
	go test internal/pkg/types/mlrval_functions_test.go $(ls internal/pkg/types/*.go | grep -v test)
mlrval_format_test:
	go test internal/pkg/types/mlrval_format_test.go $(ls internal/pkg/types/*.go|grep -v test)
regex_test:
	go test internal/pkg/lib/regex_test.go internal/pkg/lib/regex.go

# Regression tests (large number)
#
# See ./regression_test.go for information on how to get more details
# for debugging.  TL;DR is for CI jobs, we have 'go test -v'; for
# interactive use, instead of 'go test -v' simply use 'mlr regtest
# -vvv' or 'mlr regtest -s 20'. See also internal/pkg/auxents/regtest.
regression_test:
	go test -v regression_test.go

# DESTDIR is for package installs; nominally blank when this is run interactively.
# See also https://www.gnu.org/prep/standards/html_node/DESTDIR.html
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
	make -C docs/src forcebuild
	make -C docs
	@echo DONE

# Keystroke-savers
it: build check
so: install

# Please see comments in ./create-release-tarball as well as
# https://miller.readthedocs.io/en/latest/build/#creating-a-new-release-for-developers
release_tarball: build check
	./create-release-tarball

# Go does its own dependency management, outside of make.
.PHONY: build mlr check unit_test regression_test fmt dev
