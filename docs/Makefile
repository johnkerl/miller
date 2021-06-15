# Minimal makefile for Sphinx documentation
#
# Note: run this after make in the ../c directory and make in the ../man directory
# since ../c/mlr is used to autogenerate ../man/manpage.txt which is used in this directory.
# See also https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers

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
	$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
