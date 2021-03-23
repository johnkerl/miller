# Minimal makefile for Sphinx documentation
#

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?=
SPHINXBUILD   ?= sphinx-build
SOURCEDIR     = .
BUILDDIR      = _build

# Respective MANPATH entries would include /usr/local/share/man or $HOME/man.
INSTALLDIR=/usr/local/share/man/man1
INSTALLHOME=$(HOME)/man/man1

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	./genrst
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

# These targets are only for local dev work.
install:
	mkdir -p $(INSTALLDIR)
	cp mlr.1 $(INSTALLDIR)/mlr.1
installhome:
	mkdir -p $(INSTALLHOME)
	cp mlr.1 $(INSTALLHOME)/mlr.1

