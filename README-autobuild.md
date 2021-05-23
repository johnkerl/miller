# Status as of May 2021

## Go port

GitHub actions are fine for Linux/MacOS/Windows at https://github.com/johnkerl/miller/actions.

## C implementation

### Travis CI for non-Windows

* Travis was in place for years now
* As of May 2021 `travis-ci.org` is moving to `travis-ci.com`
* I am getting 404s in the Webhook setup on the GitHub side but builds are trigger on the Travis side which is baffling
* This is chewing through credits at a furious rate regardless
* Also, almost all commits these days are on the Go code so trigger a build of the C code on each commit is not a good use of resources
* In summary: Travis CI is not currently working, and isn't worth fixing

### GitHub actions for non-Windows

* Incompatible with `autoconf` as described at https://www.preining.info/blog/2018/12/git-and-autotools-a-hate-relation/
* In summary `autoconf` does not work with GitHub Actions, and this does not appear to be a forward path

### Conclusion for non-Windows

* For those rare commits which do involve C code -- until the Go port is complete -- I'll runs C makes manually.
* For releases, I'll run them manually -- which is the current process [as defined here](https://miller.readthedocs.io/en/latest/build.html#creating-a-new-release-for-developers).

### Windows builds

This used:

* [appveyor.yml](appveyor.yml)
* https://ci.appveyor.com/project/johnkerl/miller
* https://github.com/johnkerl/miller/settings/hooks

Unfortunately, I understand next to nothing about what I'm doing here -- whenever the AppVeyor build breaks (and the Travis build doesn't) I end up googling for various things in the https://ci.appveyor.com/project/johnkerl/miller build-log output, then iteratively updating `appveyor.yml` until I can get a build again.

As of May 2021 I've disabled Appveyor builds. Moving forward, for the C implementation I'll build Windows executables on a personal laptop and upload `mlr.exe` when I cut a release.
