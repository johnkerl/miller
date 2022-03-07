<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Building from source

Please also see [Installation](installing-miller.md) for information about pre-built executables.

You will need to first install Go version 1.15 or higher: please see [https://go.dev](https://go.dev).

## Miller license

Two-clause BSD license [https://github.com/johnkerl/miller/blob/master/LICENSE.txt](https://github.com/johnkerl/miller/blob/master/LICENSE.txt).

## From release tarball

* Obtain `mlr-i.j.k.tar.gz` from [https://github.com/johnkerl/miller/tags](https://github.com/johnkerl/miller/tags), replacing `i.j.k` with the desired release, e.g. `6.1.0`.
* `tar zxvf mlr-i.j.k.tar.gz`
* `cd mlr-i.j.k`
* `cd go`
* `make` creates the `./mlr` (or `.\mlr.exe` on Windows) executable
    * Without `make`: `go build github.com/johnkerl/miller/cmd/mlr`
* `make check` runs tests
    * Without `make`: `go test github.com/johnkerl/miller/internal/pkg/...` and `mlr regtest`
* `make install` installs the `mlr` executable and the `mlr` manpage
    * Without make: `go install github.com/johnkerl/miller/cmd/mlr` will install to _GOPATH_`/bin/mlr`

## From git clone

* `git clone https://github.com/johnkerl/miller`
* `make`/`go build github.com/johnkerl/miller/cmd/mlr` as above

## In case of problems

If you have any build errors, feel free to open an issue with "New Issue" at [https://github.com/johnkerl/miller/issues](https://github.com/johnkerl/miller/issues).

## Dependencies

### Required external dependencies

These are necessary to produce the `mlr` executable.

* Go version 1.15 or higher: please see [https://go.dev](https://go.dev)
* Others packaged within `go.mod` and `go.sum` which you don't need to deal with manually -- the Go build process handles them for us

### Optional external dependencies

This documentation pageset is built using [https://www.mkdocs.org/](MkDocs). Please see [https://github.com/johnkerl/miller/blob/main/docs/README.md](https://github.com/johnkerl/miller/blob/main/docs/README.md) for details.

## Creating a new release: for developers

This is my checklist for making new releases.

In this example I am using version 6.1.0 to 6.2.0; of course that will change for subsequent revisions.

* Update version found in `mlr --version` and `man mlr`:

    * Edit `internal/pkg/version/version.go` from `6.1.0-dev` to `6.2.0`.
    * Edit `miller.spec`: `Version`, and `changelog` entry
    * Run `make dev` in the Miller repo base directory
    * The ordering in this makefile rule is important: the first build creates `mlr`; the second runs `mlr` to create `manpage.txt`; the third includes `manpage.txt` into one of its outputs.
    * Commit and push.

* Create the release tarball:

    * `make release_tarball`
    * This creates `miller-6.2.0.tar.gz` which we'll upload to GitHub, the URL of which will be in our `miller.spec`
    * Prepare the source RPM following [README-RPM.md](https://github.com/johnkerl/miller/blob/main/README-RPM.md).

* Create the Github release tag:

    * Don't forget the `v` in `v6.2.0`
    * Write the release notes
    * Get binaries from latest successful build from [https://github.com/johnkerl/miller/actions](https://github.com/johnkerl/miller/actions), or, build them on buildboxes. Note that thanks to [PR 822](https://github.com/johnkerl/miller/pull/822) which introduces [goreleaser](https://github.com/johnkerl/miller/blob/main/.goreleaser.yml) there are versions for many platforms auto-built and auto-attached to the GitHub release (below). The only exception is for Windows: goreleaser makes a `.tar.gz` file but it's nice to attach a `.zip` from GitHub actions for the benefit of Windows users.
    * Attach the release tarball, Windows `.zip`, and SRPM. Double-check assets were successfully uploaded.
    * Publish the release. Note that gorelease will create and attach the rest of the binaries.

* Check the release-specific docs:

    * Look at [https://miller.readthedocs.io](https://miller.readthedocs.io) for new-version docs, after a few minutes' propagation time. Note this won't work until Miller 6 is released.

* Notify:

    * Submit `brew` pull request; notify any other distros which don't appear to have autoupdated since the previous release (notes below)
    * Similarly for `macports`: [https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile](https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile)
    * See also [README-versions.md](https://github.com/johnkerl/miller/blob/main/README-versions.md) -- distros usually catch up over time but some contacts/pings never hurt to kick-start processes after owners move on from a project they started.
    * Social-media updates.

<pre class="pre-non-highlight-non-pair">
# brew notes:
git remote add upstream https://github.com/Homebrew/homebrew-core # one-time setup only
git fetch upstream
git rebase upstream/master
git checkout -b miller-6.1.0
shasum -a 256 /path/to/mlr-6.1.0.tar.gz
edit Formula/miller.rb
# Test the URL from the line like
#   url "https://github.com/johnkerl/miller/releases/download/v6.1.0/mlr-6.1.0.tar.gz"
# in a browser for typos.
# A '@BrewTestBot Test this please' comment within the homebrew-core pull request
# will restart the homebrew travis build.
git add Formula/miller.rb
git commit -m 'miller 6.1.0'
git push -u origin miller-6.1.0
(submit the pull request)
</pre>

* Afterwork:

    * Edit `internal/pkg/version/version.go` to change version from `6.2.0` to `6.2.0-dev`.
    * `make dev`
    * Commit and push.

## Misc. development notes

I use terminal width 120 and tabwidth 4. Miller documents use the Oxford comma: not _red, yellow and green_, but rather _red, yellow, and green_.
