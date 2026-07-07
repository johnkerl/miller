# Turn src/*.md.in into respective src/*.md. Then, generate static HTML into ./site.
build:
	make -C src genmds
	mkdocs build

# Optional: build a single-file epub of the documentation for offline
# reading. Requires pandoc. Not part of the default build; Read the Docs runs
# this via .readthedocs.yaml (see docs/build-epub.sh for details).
epub:
	./build-epub.sh

.PHONY: build epub
