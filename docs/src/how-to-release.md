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
# How to create a new release

This is my checklist for making new releases.

In this example I am using version 6.2.0 to 6.3.0; of course that will change for subsequent revisions.

* Update version found in `mlr --version` and `man mlr`:

    * Edit `internal/pkg/version/version.go` from `6.2.0-dev` to `6.3.0`.
    * Edit `miller.spec`: `Version`, and `changelog` entry
    * Run `make dev` in the Miller repo base directory
    * The ordering in this makefile rule is important: the first build creates `mlr`; the second runs `mlr` to create `manpage.txt`; the third includes `manpage.txt` into one of its outputs.
    * Commit and push.

* Create the release tarball:

    * `make release_tarball`
    * This creates `miller-6.3.0.tar.gz` which we'll upload to GitHub, the URL of which will be in our `miller.spec`
    * Prepare the source RPM following [README-RPM.md](https://github.com/johnkerl/miller/blob/main/README-RPM.md).

* Create the Github release tag:

    * Don't forget the `v` in `v6.3.0`
    * Write the release notes -- save as a pre-release until below
        * Be sure the commit being used is the (non-`main`) PR commit containing the new version, or, `main` after that PR is merged back to `main`. (Otherwise, the release will be tagging the commit _before_ the changes, and `mlr version` will not show the new release number.)
    * Thanks to [PR 822](https://github.com/johnkerl/miller/pull/822) which introduces [goreleaser](https://github.com/johnkerl/miller/blob/main/.goreleaser.yml) there are versions for many platforms auto-built and auto-attached to the GitHub release.
    * Attach the release tarball and SRPM. Double-check assets were successfully uploaded.
    * Publish the release in pre-release mode, until all CI jobs finish successfully. Note that gorelease will create and attach the rest of the binaries.
    * Before marking the release as public, download an executable from among the generated binaries and make sure its `mlr version` prints what you expect -- else, restart this process.
    * Then mark the release as public.

* Check the release-specific docs:

    * Look at [https://miller.readthedocs.io](https://miller.readthedocs.io) for new-version docs, after a few minutes' propagation time. Note this won't work until Miller 6 is released.

* Notify:

    * Submit `brew` pull request; notify any other distros which don't appear to have autoupdated since the previous release (notes below)
    * Similarly for `macports`: [https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile](https://github.com/macports/macports-ports/blob/master/textproc/miller/Portfile)
    * See also [README-versions.md](https://github.com/johnkerl/miller/blob/main/README-versions.md) -- distros usually catch up over time but some contacts/pings never hurt to kick-start processes after owners move on from a project they started.
    * Social-media updates.
    * Brew notes:
      * [How to submit a version upgrade](https://github.com/Homebrew/homebrew-core/blob/HEAD/CONTRIBUTING.md#to-submit-a-version-upgrade-for-the-foo-formula)
      * `brew bump-formula-pr --force --strict miller --url https://github.com/johnkerl/miller/releases/download/v6.2.0/miller-6.2.0.tar.gz --sha256 xxx` with `xxx` from `shasum -a 256 miller-6.2.0.tar.gz`.

* Afterwork:

    * Edit `internal/pkg/version/version.go` to change version from `6.3.0` to `6.3.0-dev`.
    * `make dev`
    * Commit and push.
