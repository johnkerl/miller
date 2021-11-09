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
precommit:
	make -C go fmt
	make -C go build
	make -C go check
	make -C man build
	make -C docs
	echo DONE

# Keystroke-saver
itso: check build install

# Go does its own dependency management, outside of make.
.PHONY: build precommit
