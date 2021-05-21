## Wndows builds

This uses/used:

* [appveyor.yml](appveyor.yml)
* https://ci.appveyor.com/project/johnkerl/miller.
* https://github.com/johnkerl/miller/settings/hooks

Unfortunately, I understand next to nothing about what I'm doing here -- whenever the AppVeyor build breaks (and the Travis build doesn't) I end up googling for various things in the https://ci.appveyor.com/project/johnkerl/miller build-log output, then iteratively updating `appveyor.yml` until I can get a build again.

As of May 2021 I've disabled Appveyor builds. Moving forward:

* For the C implementation I'll build Windows executables on a personal laptop and upload `mlr.exe` when I cut a release
* The [Go port](https://github.com/johnkerl/miller/blob/main/go/README.md) builds for Windows without problem via GitHub Actions at https://github.com/johnkerl/miller/actions
