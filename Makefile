PREFIX=/usr/local
INSTALLDIR=$(PREFIX)/bin

# ================================================================
# General-use targets

# This must remain the first target in this file, which is what 'make' with no
# arguments will run.
build:
	go build github.com/johnkerl/miller/cmd/mlr
	@echo "Build complete. The Miller executable is ./mlr (or .\mlr.exe on Windows)."
	@echo "You can use 'make check' to run tests".

# For interactive use, 'mlr regtest' offers more options and transparency.
check: unit-test regression-test
	@echo "Tests complete. You can use 'make install' if you like, optionally preceded"
	@echo "by './configure --prefix=/your/install/path' if you wish to install to"
	@echo "somewhere other than /usr/local/bin -- the default prefix is /usr/local."

# DESTDIR is for package installs; nominally blank when this is run interactively.
# See also https://www.gnu.org/prep/standards/html_node/DESTDIR.html
install: build
	cp mlr $(DESTDIR)/$(INSTALLDIR)
	make -C man install

# ================================================================
# Dev targets

# ----------------------------------------------------------------
# Unit tests (small number)
unit-test ut:
	go test github.com/johnkerl/miller/internal/pkg/...

ut-lib:
	go test github.com/johnkerl/miller/internal/pkg/lib...
ut-mlv:
	go test github.com/johnkerl/miller/internal/pkg/mlrval/...
ut-bifs:
	go test github.com/johnkerl/miller/internal/pkg/bifs/...
ut-input:
	go test github.com/johnkerl/miller/internal/pkg/input/...

bench:
	go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/...
bench-mlv:
	go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/mlrval/...
bench-input:
	go test -run=nonesuch -bench=. github.com/johnkerl/miller/internal/pkg/input/...

# ----------------------------------------------------------------
# Regression tests (large number)
#
# See ./regression_test.go for information on how to get more details
# for debugging.  TL;DR is for CI jobs, we have 'go test -v'; for
# interactive use, instead of 'go test -v' simply use 'mlr regtest
# -vvv' or 'mlr regtest -s 20'. See also internal/pkg/auxents/regtest.
regression-test:
	go test -v regression_test.go

# go fmt ./... finds experimental C files which we want to ignore.
fmt:
	-go fmt ./cmd/...
	-go fmt ./internal/pkg/...
	-go fmt ./regression_test.go

# Needs first: go install honnef.co/go/tools/cmd/staticcheck@latest
# See also: https://staticcheck.io
staticcheck:
	staticcheck ./...

# ----------------------------------------------------------------
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

docs:
	make -C docs

# ----------------------------------------------------------------
# Keystroke-savers
it: build check
so: install
sure: build check
mlr:
	go build github.com/johnkerl/miller/cmd/mlr

# ----------------------------------------------------------------
# Please see comments in ./create-release-tarball as well as
# https://miller.readthedocs.io/en/latest/build/#creating-a-new-release-for-developers
release_tarball: build check
	./create-release-tarball

# ================================================================
# Go does its own dependency management, outside of make.
.PHONY: build mlr check unit_test regression_test bench fmt staticcheck dev docs
