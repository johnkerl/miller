# What is Miller?

**Miller is like awk, sed, cut, join, and sort for data formats such as CSV, TSV, tabular JSON and positionally-indexed.**

# What can Miller do for me?

With Miller, you get to use named fields without needing to count positional
indices, using familiar formats such as CSV, TSV, JSON, and
positionally-indexed.  Then, on the fly, you can add new fields which are
functions of existing fields, drop fields, sort, aggregate statistically,
pretty-print, and more.

![cover-art](./docs/src/coverart/cover-combined.png)

* Miller operates on **key-value-pair data** while the familiar
Unix tools operate on integer-indexed fields: if the natural data structure for
the latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.

* Miller handles a **variety of data formats**,
including but not limited to the familiar **CSV**, **TSV**, and **JSON**.
(Miller can handle **positionally-indexed data** too!)

In the above image you can see how Miller embraces the common themes of
key-value-pair data in a variety of data formats.

# Getting started

* [Miller in 10 minutes](https://miller.readthedocs.io/en/latest/10min)
* [A quick tutorial on Miller](https://www.ict4g.net/adolfo/notes/data-analysis/miller-quick-tutorial.html)
* [Tools to manipulate CSV files from the Command Line](https://www.ict4g.net/adolfo/notes/data-analysis/tools-to-manipulate-csv.html)
* [www.togaware.com/linux/survivor/CSV_Files.html](https://www.togaware.com/linux/survivor/CSV_Files.html)
* [MLR for CSV manipulation](https://guillim.github.io/terminal/2018/06/19/MLR-for-CSV-manipulation.html)
* [Linux Magazine: Process structured text files with Miller](https://www.linux-magazine.com/Issues/2016/187/Miller)
* [Miller: Command Line CSV File Processing](https://onepointzero.app/posts/miller-command-line-csv-file-processing/)

# More documentation links

* [**Full documentation**](https://miller.readthedocs.io/)
* [Miller's license is two-clause BSD](https://github.com/johnkerl/miller/blob/main/LICENSE.txt)
* [Notes about issue-labeling in the Github repo](https://github.com/johnkerl/miller/wiki/Issue-labeling)
* [Active issues](https://github.com/johnkerl/miller/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc)

# Installing

There's a good chance you can get Miller pre-built for your system:

[![Ubuntu](https://img.shields.io/badge/distros-ubuntu-db4923.svg)](https://launchpad.net/ubuntu/+source/miller)
[![Ubuntu 16.04 LTS](https://img.shields.io/badge/distros-ubuntu1604lts-db4923.svg)](https://launchpad.net/ubuntu/xenial/+package/miller)
[![Fedora](https://img.shields.io/badge/distros-fedora-173b70.svg)](https://apps.fedoraproject.org/packages/miller)
[![Debian](https://img.shields.io/badge/distros-debian-c70036.svg)](https://packages.debian.org/stable/miller)
[![Gentoo](https://img.shields.io/badge/distros-gentoo-4e4371.svg)](https://packages.gentoo.org/packages/sys-apps/miller)

[![Pro-Linux](https://img.shields.io/badge/distros-prolinux-3a679d.svg)](http://www.pro-linux.de/cgi-bin/DBApp/check.cgi?ShowApp..20427.100)
[![Arch Linux](https://img.shields.io/badge/distros-archlinux-1792d0.svg)](https://aur.archlinux.org/packages/miller-git)

[![NetBSD](https://img.shields.io/badge/distros-netbsd-f26711.svg)](http://pkgsrc.se/textproc/miller)
[![FreeBSD](https://img.shields.io/badge/distros-freebsd-8c0707.svg)](https://www.freshports.org/textproc/miller/)

[![Anaconda](https://img.shields.io/badge/distros-anaconda-63ad41.svg)](https://anaconda.org/conda-forge/miller/)
[![Homebrew/MacOSX](https://img.shields.io/badge/distros-homebrew-ba832b.svg)](https://formulae.brew.sh/formula/miller)
[![MacPorts/MacOSX](https://img.shields.io/badge/distros-macports-1376ec.svg)](https://www.macports.org/ports.php?by=name&substr=miller)
[![Chocolatey](https://img.shields.io/badge/distros-chocolatey-red.svg)](https://chocolatey.org/packages/miller)

|OS|Installation command|
|---|---|
|Linux|`yum install miller`<br/> `apt-get install miller`|
|Mac|`brew install miller`<br/>`port install miller`|
|Windows|`choco install miller`|

See also [building from source](https://miller.readthedocs.io/en/latest/build.html).

# Build status

[![Go-port multi-platform build status](https://github.com/johnkerl/miller/actions/workflows/go.yml/badge.svg)](https://github.com/johnkerl/miller/actions)

# Building from source

* `make`: takes just a few seconds and produces the Miller executable, which is `./mlr` (or `.\mlr.exe` on Windows).
  * Without `make`: `go build github.com/johnkerl/miller/cmd/mlr`
* `make check` runs tests.
  * Without `make`: `go test github.com/johnkerl/miller/internal/pkg/...` and `mlr regtest`
* `make install` installs executable `/usr/local/bin/mlr` and manual page `/usr/local/share/man/man1/mlr.1` (so you can do `man mlr`).
  * You can instead do `./configure --prefix=/some/install/path` followed by `make install` if you want to install somewhere other than `/usr/local`.
  * Without make: `go install github.com/johnkerl/miller/cmd/mlr` will install to _GOPATH_`/bin/mlr`
* See also the doc page on [building from source](https://miller.readthedocs.io/en/latest/build).
* For more developer information please see [README-go-port.md](./README-go-port.md).

# License

[License: BSD2](https://github.com/johnkerl/miller/blob/main/LICENSE.txt)

# Community

* Discussion forum: https://github.com/johnkerl/miller/discussions
* Feature requests / bug reports: https://github.com/johnkerl/miller/issues

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-40-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

# Features

* Miller is **multi-purpose**: it's useful for **data cleaning**,
**data reduction**, **statistical reporting**, **devops**, **system
administration**, **log-file processing**, **format conversion**, and
**database-query post-processing**.

* You can use Miller to snarf and munge **log-file data**, including selecting
out relevant substreams, then produce CSV format and load that into
all-in-memory/data-frame utilities for further statistical and/or graphical
processing.

* Miller complements **data-analysis tools** such as **R**, **pandas**, etc.:
you can use Miller to **clean** and **prepare** your data. While you can do
**basic statistics** entirely in Miller, its streaming-data feature and
single-pass algorithms enable you to **reduce very large data sets**.

* Miller complements SQL **databases**: you can slice, dice, and reformat data
on the client side on its way into or out of a database. You can also reap some
of the benefits of databases for quick, setup-free one-off tasks when you just
need to query some data in disk files in a hurry.

* Miller also goes beyond the classic Unix tools by stepping fully into our
modern, **no-SQL** world: its essential record-heterogeneity property allows
Miller to operate on data where records with different schema (field names) are
interleaved.

* Miller is **streaming**: most operations need only a single record in
memory at a time, rather than ingesting all input before producing any output.
For those operations which require deeper retention (`sort`, `tac`, `stats1`),
Miller retains only as much data as needed. This means that whenever
functionally possible, you can operate on files which are larger than your
system&rsquo;s available RAM, and you can use Miller in **tail -f** contexts.

* Miller is **pipe-friendly** and interoperates with the Unix toolkit.

* Miller's I/O formats include **tabular pretty-printing**, **positionally
  indexed** (Unix-toolkit style), CSV, JSON, and others.

* Miller does **conversion** between formats.

* Miller's **processing is format-aware**: e.g. CSV `sort` and `tac` keep header lines first.

* Miller has high-throughput **performance** on par with the Unix toolkit.

* Miller is written in portable, modern Go, with **zero runtime dependencies**.
You can download or compile a single binary, `scp` it to a faraway machine,
and expect it to work.

# What people are saying about Miller

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Today I discovered Miller—it&#39;s like jq but for CSV: <a href="https://t.co/pn5Ni241KM">https://t.co/pn5Ni241KM</a><br><br>Also, &quot;Miller complements data-analysis tools such as R, pandas, etc.: you can use Miller to clean and prepare your data.&quot; <a href="https://twitter.com/GreatBlueC?ref_src=twsrc%5Etfw">@GreatBlueC</a> <a href="https://twitter.com/nfmcclure?ref_src=twsrc%5Etfw">@nfmcclure</a></p>&mdash; Adrien Trouillaud (@adrienjt) <a href="https://twitter.com/adrienjt/status/1308963056592891904?ref_src=twsrc%5Etfw">September 24, 2020</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Underappreciated swiss-army command-line chainsaw.<br><br>&quot;Miller is like awk, sed, cut, join, and sort for [...] CSV, TSV, and [...] JSON.&quot; <a href="https://t.co/TrQqSUK3KK">https://t.co/TrQqSUK3KK</a></p>&mdash; Dirk Eddelbuettel (@eddelbuettel) <a href="https://twitter.com/eddelbuettel/status/836555980771061760?ref_src=twsrc%5Etfw">February 28, 2017</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Miller looks like a great command line tool for working with CSV data. Sed, awk, cut, join all rolled into one: <a href="http://t.co/9BBb6VCZ6Y">http://t.co/9BBb6VCZ6Y</a></p>&mdash; Mike Loukides (@mikeloukides) <a href="https://twitter.com/mikeloukides/status/632885317389950976?ref_src=twsrc%5Etfw">August 16, 2015</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Miller is like sed, awk, cut, join, and sort for name-indexed data such as CSV: <a href="http://t.co/1zPbfg6B2W">http://t.co/1zPbfg6B2W</a> - handy tool!</p>&mdash; Ilya Grigorik (@igrigorik) <a href="https://twitter.com/igrigorik/status/635134857283153920?ref_src=twsrc%5Etfw">August 22, 2015</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Btw, I think Miller is the best CLI tool to deal with CSV. I used to use this when I need to preprocess too big CSVs to load into R (now we have vroom, so such cases might be rare, though...)<a href="https://t.co/kUjrSSGJoT">https://t.co/kUjrSSGJoT</a></p>&mdash; Hiroaki Yutani (@yutannihilat_en) <a href="https://twitter.com/yutannihilat_en/status/1252392795676934144?ref_src=twsrc%5Etfw">April 21, 2020</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Miller: a *format-aware* data munging tool By <a href="https://twitter.com/__jo_ker__?ref_src=twsrc%5Etfw">@__jo_ker__</a> to overcome limitations with *line-aware* workshorses like awk, sed et al <a href="https://t.co/LCyPkhYvt9">https://t.co/LCyPkhYvt9</a><br><br>The project website is a fantastic example of good software documentation!!</p>&mdash; Donny Daniel (@dnnydnl) <a href="https://twitter.com/dnnydnl/status/1038883999391932416?ref_src=twsrc%5Etfw">September 9, 2018</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">Holy holly data swiss army knife batman! How did no one suggest Miller <a href="https://t.co/JGQpmRAZLv">https://t.co/JGQpmRAZLv</a> for solving database cleaning / ETL issues to me before <br><br>Congrats to <a href="https://twitter.com/__jo_ker__?ref_src=twsrc%5Etfw">@__jo_ker__</a> for amazingly intuitive tool for critical data management tasks!<a href="https://twitter.com/hashtag/DataScienceandLaw?src=hash&amp;ref_src=twsrc%5Etfw">#DataScienceandLaw</a> <a href="https://twitter.com/hashtag/ComputationalLaw?src=hash&amp;ref_src=twsrc%5Etfw">#ComputationalLaw</a></p>&mdash; James Miller (@japanlawprof) <a href="https://twitter.com/japanlawprof/status/1006547451409518597?ref_src=twsrc%5Etfw">June 12, 2018</a></blockquote>

<blockquote class="twitter-tweet"><p lang="en" dir="ltr">🤯<a href="https://twitter.com/__jo_ker__?ref_src=twsrc%5Etfw">@__jo_ker__</a>&#39;s Miller easily reads, transforms, + writes all sorts of tabular data. It&#39;s standalone, fast, and built for streaming data (operating on one line at a time, so you can work on files larger than memory).<br><br>And the docs are dream. I&#39;ve been reading them all morning! <a href="https://t.co/Be2pGPZK6t">https://t.co/Be2pGPZK6t</a></p>&mdash; Benjamin Wolfe (he/him) (@BenjaminWolfe) <a href="https://twitter.com/BenjaminWolfe/status/1435966268499128324?ref_src=twsrc%5Etfw">September 9, 2021</a></blockquote>

## Contributors ✨

Thanks to all the fine people who help make Miller better ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/aborruso"><img src="https://avatars.githubusercontent.com/u/30607?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Andrea Borruso</b></sub></a><br /><a href="#ideas-aborruso" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-aborruso" title="Design">🎨</a></td>
    <td align="center"><a href="https://sjackman.ca/"><img src="https://avatars.githubusercontent.com/u/291551?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Shaun Jackman</b></sub></a><br /><a href="#ideas-sjackman" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="http://www.fredtrotter.com/"><img src="https://avatars.githubusercontent.com/u/83133?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Fred Trotter</b></sub></a><br /><a href="#ideas-ftrotter" title="Ideas, Planning, & Feedback">🤔</a> <a href="#design-ftrotter" title="Design">🎨</a></td>
    <td align="center"><a href="https://github.com/Komosa"><img src="https://avatars.githubusercontent.com/u/10688154?v=4?s=50" width="50px;" alt=""/><br /><sub><b>komosa</b></sub></a><br /><a href="#ideas-Komosa" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/jungle-boogie"><img src="https://avatars.githubusercontent.com/u/1111743?v=4?s=50" width="50px;" alt=""/><br /><sub><b>jungle-boogie</b></sub></a><br /><a href="#ideas-jungle-boogie" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/0-wiz-0"><img src="https://avatars.githubusercontent.com/u/2221844?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Thomas Klausner</b></sub></a><br /><a href="#infra-0-wiz-0" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a></td>
    <td align="center"><a href="https://github.com/skitt"><img src="https://avatars.githubusercontent.com/u/2128935?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Stephen Kitt</b></sub></a><br /><a href="#platform-skitt" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="http://leahneukirchen.org/"><img src="https://avatars.githubusercontent.com/u/139?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Leah Neukirchen</b></sub></a><br /><a href="#ideas-leahneukirchen" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/lgbaldoni"><img src="https://avatars.githubusercontent.com/u/1450716?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Luigi Baldoni</b></sub></a><br /><a href="#platform-lgbaldoni" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://yutani.rbind.io/"><img src="https://avatars.githubusercontent.com/u/1978793?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Hiroaki Yutani</b></sub></a><br /><a href="#ideas-yutannihilation" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://3e.org/"><img src="https://avatars.githubusercontent.com/u/41439?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Daniel M. Drucker</b></sub></a><br /><a href="#ideas-dmd" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/NikosAlexandris"><img src="https://avatars.githubusercontent.com/u/7046639?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Nikos Alexandris</b></sub></a><br /><a href="#ideas-NikosAlexandris" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/kundeng"><img src="https://avatars.githubusercontent.com/u/89032?v=4?s=50" width="50px;" alt=""/><br /><sub><b>kundeng</b></sub></a><br /><a href="#platform-kundeng" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="http://victorsergienko.com/"><img src="https://avatars.githubusercontent.com/u/151199?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Victor Sergienko</b></sub></a><br /><a href="#platform-singalen" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/gromgit"><img src="https://avatars.githubusercontent.com/u/215702?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Adrian Ho</b></sub></a><br /><a href="#design-gromgit" title="Design">🎨</a></td>
    <td align="center"><a href="https://github.com/Zachp"><img src="https://avatars.githubusercontent.com/u/1316442?v=4?s=50" width="50px;" alt=""/><br /><sub><b>zachp</b></sub></a><br /><a href="#platform-Zachp" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://dsel.net/"><img src="https://avatars.githubusercontent.com/u/921669?v=4?s=50" width="50px;" alt=""/><br /><sub><b>David Selassie</b></sub></a><br /><a href="#ideas-davidselassie" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="http://www.joelparkerhenderson.com/"><img src="https://avatars.githubusercontent.com/u/27145?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Joel Parker Henderson</b></sub></a><br /><a href="#ideas-joelparkerhenderson" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/divtiply"><img src="https://avatars.githubusercontent.com/u/5359679?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Michel Ace</b></sub></a><br /><a href="#ideas-divtiply" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="http://fuco1.github.io/sitemap.html"><img src="https://avatars.githubusercontent.com/u/2664959?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Matus Goljer</b></sub></a><br /><a href="#ideas-Fuco1" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/terorie"><img src="https://avatars.githubusercontent.com/u/21371810?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Richard Patel</b></sub></a><br /><a href="#platform-terorie" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://blog.kub1x.org/"><img src="https://avatars.githubusercontent.com/u/1833840?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Jakub Podlaha</b></sub></a><br /><a href="#design-kub1x" title="Design">🎨</a></td>
    <td align="center"><a href="https://goo.gl/ZGZynx"><img src="https://avatars.githubusercontent.com/u/85767?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Miodrag Milić</b></sub></a><br /><a href="#platform-majkinetor" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/derekmahar"><img src="https://avatars.githubusercontent.com/u/6047?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Derek Mahar</b></sub></a><br /><a href="#ideas-derekmahar" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/spmundi"><img src="https://avatars.githubusercontent.com/u/38196185?v=4?s=50" width="50px;" alt=""/><br /><sub><b>spmundi</b></sub></a><br /><a href="#ideas-spmundi" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/koernepr"><img src="https://avatars.githubusercontent.com/u/24551942?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Peter Körner</b></sub></a><br /><a href="#security-koernepr" title="Security">🛡️</a></td>
    <td align="center"><a href="https://github.com/rubyFeedback"><img src="https://avatars.githubusercontent.com/u/46686565?v=4?s=50" width="50px;" alt=""/><br /><sub><b>rubyFeedback</b></sub></a><br /><a href="#ideas-rubyFeedback" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/rbolsius"><img src="https://avatars.githubusercontent.com/u/2106964?v=4?s=50" width="50px;" alt=""/><br /><sub><b>rbolsius</b></sub></a><br /><a href="#platform-rbolsius" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/awildturtok"><img src="https://avatars.githubusercontent.com/u/1553491?v=4?s=50" width="50px;" alt=""/><br /><sub><b>awildturtok</b></sub></a><br /><a href="#ideas-awildturtok" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/agguser"><img src="https://avatars.githubusercontent.com/u/1206106?v=4?s=50" width="50px;" alt=""/><br /><sub><b>agguser</b></sub></a><br /><a href="#ideas-agguser" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/jganong"><img src="https://avatars.githubusercontent.com/u/2783890?v=4?s=50" width="50px;" alt=""/><br /><sub><b>jganong</b></sub></a><br /><a href="#ideas-jganong" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://www.linkedin.com/in/fulvio-scapin"><img src="https://avatars.githubusercontent.com/u/69568?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Fulvio Scapin</b></sub></a><br /><a href="#ideas-trantor" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/torbiak"><img src="https://avatars.githubusercontent.com/u/109347?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Jordan Torbiak</b></sub></a><br /><a href="#ideas-torbiak" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/Andy1978"><img src="https://avatars.githubusercontent.com/u/240064?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Andreas Weber</b></sub></a><br /><a href="#ideas-Andy1978" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/vapniks"><img src="https://avatars.githubusercontent.com/u/174330?v=4?s=50" width="50px;" alt=""/><br /><sub><b>vapniks</b></sub></a><br /><a href="#platform-vapniks" title="Packaging/porting to new platform">📦</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://github.com/89z"><img src="https://avatars.githubusercontent.com/u/73562167?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Zombo</b></sub></a><br /><a href="#platform-89z" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/BEFH"><img src="https://avatars.githubusercontent.com/u/3386600?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Brian Fulton-Howard</b></sub></a><br /><a href="#platform-BEFH" title="Packaging/porting to new platform">📦</a></td>
    <td align="center"><a href="https://github.com/ChCyrill"><img src="https://avatars.githubusercontent.com/u/2165604?v=4?s=50" width="50px;" alt=""/><br /><sub><b>ChCyrill</b></sub></a><br /><a href="#ideas-ChCyrill" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/jauderho"><img src="https://avatars.githubusercontent.com/u/13562?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Jauder Ho</b></sub></a><br /><a href="https://github.com/johnkerl/miller/commits?author=jauderho" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/psacawa"><img src="https://avatars.githubusercontent.com/u/21274063?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Paweł Sacawa</b></sub></a><br /><a href="https://github.com/johnkerl/miller/issues?q=author%3Apsacawa" title="Bug reports">🐛</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

<a href="https://github.com/johnkerl/miller/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=johnkerl/miller" />
</a>

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind are welcome!
