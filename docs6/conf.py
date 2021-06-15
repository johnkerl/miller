# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# https://www.sphinx-doc.org/en/master/usage/configuration.html

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
# import os
# import sys
# sys.path.insert(0, os.path.abspath('.'))


# -- Project information -----------------------------------------------------

project = 'Miller'
copyright = '2020, John Kerl'
author = 'John Kerl'

# The full version, including alpha/beta/rc tags
release = '6.0.0-alpha'

# -- General configuration ---------------------------------------------------
master_doc = 'index'

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
]

# Add any paths that contain templates here, relative to this directory.
templates_path = ['_templates']

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = ['_build', 'Thumbs.db', '.DS_Store']

# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
#html_theme = 'alabaster'
#html_theme = 'classic'
#html_theme = 'sphinxdoc'
#html_theme = 'nature'
html_theme = 'scrolls'

# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named "default.css" will overwrite the builtin "default.css".
html_static_path = ['_static']

# ----------------------------------------------------------------
# Include code-sample files in the Sphinx build tree.
# See also https://github.com/johnkerl/miller/issues/560.
#
# There is a problem, and an opportunity for a hack.
#
# * Our data files are in ./data/* (and a few other subdirs of .).
#
# * We want them copied to ./_build/data/* so that we can symlink from our doc
#   files (written in ./*.rst. autogenned to HTML in ./_build/html/*.html) to
#   relative paths like ./data/a.csv.
#
# * If we use html_extra_path = ['data'] then the files like ./data/a.csv
#   are copied to _build/html/a.csv -- one directory 'down'. This means that
#   example Miller commands are shown in the generated HTML using 'mlr --csv
#   cat data/a.csv' but 'data/a.csv' doesn't exist relative to _build/html.
#   This is bad enough for local Sphinx builds but worse for readthedocs
#   (https://miller.readthedocs.io) since *only* _build/html files are put into
#   readthedocs.
#
# * In our Makefile it's easy enough to do some cp -a commands from ./data
#   to ./_build_html/data etc. for local Sphinx builds -- however, readthedocs
#   doesn't use the Makefile at all, only this conf.py file.
#
# * Hence the hack: we have a subdir ./sphinx-hack which has a symlink
#   ./sphinx-hack/data pointing to ./data. So when the Sphinx build executes
#   html_extra_path and removes one directory level, it's an 'extra' level we
#   can do without.
#
# * This all relies on symlinks being propagated through GitHub version
#   control, readthedocs, and Sphinx build at readthedocs.

html_extra_path = [
  'sphinx-hack',
  '10-1.sh',
  '10-2.sh',
  'circle.csv',
  'commas.csv',
  'dates.csv',
  'example.csv',
  'expo-sample.sh',
  'log.txt',
  'make.bat',
  'manpage.txt',
  'oosvar-example-ewma.sh',
  'oosvar-example-sum-grouped.sh',
  'oosvar-example-sum.sh',
  'sample_mlrrc',
  'square.csv',
  'triangle.csv',
  'variance.mlr',
  'verb-example-ewma.sh',
]
