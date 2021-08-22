<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Installation

## Prebuilt executables via package managers

*Note: Miller 6 is in pre-release so the above commands will get you Miller 5.
Once Miller 6 is released, the commands in this section will install Miller 6 for you.
Until then, please see the following sections for how to get Miller 6.*

[Homebrew](https://brew.sh/) installation support for OS X is available via

<pre class="pre-highlight-non-pair">
<b>brew update && brew install miller</b>
</pre>

... and also via [MacPorts](https://www.macports.org/):

<pre class="pre-highlight-non-pair">
<b>sudo port selfupdate && sudo port install miller</b>
</pre>

You may already have the `mlr` executable available in your platform's package manager on NetBSD, Debian Linux, Ubuntu Xenial and upward, Arch Linux, or perhaps other distributions. For example, on various Linux distributions you might do one of the following:

<pre class="pre-highlight-non-pair">
<b>sudo apt-get install miller</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>sudo apt install miller</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>sudo yum install miller</b>
</pre>

On Windows, Miller is available via [Chocolatey](https://chocolatey.org/):

<pre class="pre-highlight-non-pair">
<b>choco install miller</b>
</pre>

## Prebuilt executables via GitHub per release

Please see [https://github.com/johnkerl/miller/releases](https://github.com/johnkerl/miller/releases) where there are builds for OS X Yosemite, Linux x86-64 (dynamically linked), and Windows.

## Prebuilt executables via GitHub per commit

Miller is [autobuilt for **Linux**, **MacOS**, and **Windows** using **GitHub Actions** on every commit](https://github.com/johnkerl/miller/actions): select the latest build and click _Artifacts_. (These are retained for 5 days after each commit.)

## Building from source

Please see [Building from source](build.md).
