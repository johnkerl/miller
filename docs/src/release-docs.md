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
# Documents for releases

If your `mlr version` says something like `mlr 6.0.0-dev`, with the `-dev` suffix, you're likely building from source, or you've obtained a recent artifact from GitHub Actions -- 
the page [https://miller.readthedocs.io](https://miller.readthedocs.io) contains information for the latest contributions to the [Miller repository](https://github.com/johnkerl/miller).

If your `mlr version` says something like `Miller v5.10.2` or `mlr 6.0.0`, without the `-dev` suffix, you're likely using a Miller executable from a package manager -- please see below for the documentation for Miller as of the release you're using.

| Release | Docs                                                                | Release notes |
|---------|---------------------------------------------------------------------|---------------|
6.5.0     | [Miller 6.5.0](https://johnkerl.org/miller-docs-by-release/6.5.0)  | [Bugfixes and memory-reduction optimizations](https://github.com/johnkerl/miller/releases/tag/v6.5.0) |
6.4.0     | [Miller 6.4.0](https://johnkerl.org/miller-docs-by-release/6.4.0)  | [ mlr summary verb, exec() function, mlr cat --filename, multiline string literals, and more](https://github.com/johnkerl/miller/releases/tag/v6.4.0) |
6.3.0     | [Miller 6.3.0](https://johnkerl.org/miller-docs-by-release/6.3.0)  | [Windows terminal colors, Latin-1, and more](https://github.com/johnkerl/miller/releases/tag/v6.3.0) |
6.2.0     | [Miller 6.2.0](https://johnkerl.org/miller-docs-by-release/6.2.0)  | [Restore --tsvlite; add gssub and expand dhms functions](https://github.com/johnkerl/miller/releases/tag/v6.2.0) |
6.1.0     | [Miller 6.1.0](https://johnkerl.org/miller-docs-by-release/6.1.0)  | [Natural sort, true TSV, sliding-window averages, and more](https://github.com/johnkerl/miller/releases/tag/v6.1.0) |
6.0.0     | [Miller 6.0.0](https://johnkerl.org/miller-docs-by-release/6.0.0)  | [Miller 6](https://github.com/johnkerl/miller/releases/tag/v6.0.0) |
5.10.4    | [Miller 5.10.4](https://johnkerl.org/miller-docs-by-release/5.10.4) | [5.10 bugfix for issue #1108](https://github.com/johnkerl/miller/releases/tag/v5.10.4) |
5.10.3    | [Miller 5.10.3](https://johnkerl.org/miller-docs-by-release/5.10.3) | [Address Conda-build issue](https://github.com/johnkerl/miller/releases/tag/v5.10.3) |
5.10.2    | [Miller 5.10.2](https://johnkerl.org/miller-docs-by-release/5.10.2) | [Restore mlr manpage to distro file](https://github.com/johnkerl/miller/releases/tag/v5.10.2) |
5.10.1    |                                                                     | [Bugfixes](https://github.com/johnkerl/miller/releases/tag/v5.10.1) |
5.10.0    | [Miller 5.10.0](https://johnkerl.org/miller-docs-by-release/5.10.0) | [sort-within-records, unsparsify -f, misc updates; Go-port beta](https://github.com/johnkerl/miller/releases/tag/v5.10.0) |
5.9.1     |                                                                     | [Security update: disallow --prepipe in .mlrrc](https://github.com/johnkerl/miller/releases/tag/v5.9.1) |
5.9.0     | [Miller 5.9.0](https://johnkerl.org//miller-docs-by-release/5.9.0/) | [.mlrrc feature, and fix Windows build](https://github.com/johnkerl/miller/releases/tag/v5.9.0) |
5.8.0     | [Miller 5.8.0](https://johnkerl.org//miller-docs-by-release/5.8.0/) | [Better environment-variable support, new 'count' verb, bugfixes](https://github.com/johnkerl/miller/releases/tag/v5.8.0) |
5.7.0     | [Miller 5.7.0](https://johnkerl.org//miller-docs-by-release/5.7.0/) | [Ports, bugfixes, and keystroke-savers](https://github.com/johnkerl/miller/releases/tag/v5.7.0) |
5.6.2     | [Miller 5.6.2](https://johnkerl.org//miller-docs-by-release/5.6.2/) | [Bug fix for CSV/TSV with many files](https://github.com/johnkerl/miller/releases/tag/v5.6.2) |
5.6.1     |                                                                     | [Mobile-friendly docs](https://github.com/johnkerl/miller/releases/tag/v5.6.1) |
5.6.0     | [Miller 5.6.0](https://johnkerl.org//miller-docs-by-release/5.6.0/) | [System calls / external commands, ASV/USV support, and bulk numeric formatting](https://github.com/johnkerl/miller/releases/tag/v5.6.0) |
5.5.0     | [Miller 5.5.0](https://johnkerl.org//miller-docs-by-release/5.5.0/) | [Positional indexing and other data-cleaning features](https://github.com/johnkerl/miller/releases/tag/v5.5.0) |
5.4.0     | [Miller 5.4.0](https://johnkerl.org//miller-docs-by-release/5.4.0/) | [New data-cleaning features, Windows mlr.exe, limited localtime support, and bugfixes](https://github.com/johnkerl/miller/releases/tag/5.4.0) |
5.3.0     | [Miller 5.3.0](https://johnkerl.org//miller-docs-by-release/5.3.0/) | [Data comments, documentation improvements, and bug fixes](https://github.com/johnkerl/miller/releases/tag/v5.3.0) |
5.2.2     |                                                                     | [Bug-fix release: 64-bit aggregators](https://github.com/johnkerl/miller/releases/tag/v5.2.2) |
5.2.1     |                                                                     | [Fix non-x86/gcc7 build error](https://github.com/johnkerl/miller/releases/tag/v5.2.1) |
5.2.0     | [Miller 5.2.0](https://johnkerl.org//miller-docs-by-release/5.2.0/) | [stats across regexed field names, string/num stats, CSV UTF BOM strip](https://github.com/johnkerl/miller/releases/tag/v5.2.0) |
5.1.0     |                                                                     | [JSON-array support, fractional seconds in strptime/strftime, and other minor features](https://github.com/johnkerl/miller/releases/tag/v5.1.0) |
5.1.0w    |                                                                     | [MLR.EXE: Windows beta](https://github.com/johnkerl/miller/releases/tag/v5.1.0w) |
5.0.1     |                                                                     | [Two minor bugfixes](https://github.com/johnkerl/miller/releases/tag/v5.0.1) |
5.0.0     | [Miller 5.0.0](https://johnkerl.org//miller-docs-by-release/5.0.0/) | [Autodetected line-endings, in-place mode, user-defined functions, and more](https://github.com/johnkerl/miller/releases/tag/v5.0.0) |
4.5.0     |                                                                     | [Customizable output format for redirected output](https://github.com/johnkerl/miller/releases/tag/v4.5.0) |
4.4.0     |                                                                     | [Redirected output, row-value shift, and other features](https://github.com/johnkerl/miller/releases/tag/v4.4.0) |
4.3.0     |                                                                     | [Interpolated percentiles, markdown-tabular output format, CSV-quote preservation](https://github.com/johnkerl/miller/releases/tag/v4.3.0) |
4.2.0     |                                                                     | [Multi-emit](https://github.com/johnkerl/miller/releases/tag/v4.2.0) |
4.1.0     |                                                                     | [for/if/while and various features](https://github.com/johnkerl/miller/releases/tag/v4.1.0) |
4.0.0     | [Miller 4.0.0](https://johnkerl.org//miller-docs-by-release/4.0.0/) | [Variables, begin/end blocks, pattern-action blocks](https://github.com/johnkerl/miller/releases/tag/v4.0.0) |
3.5.0     |                                                                     | [New data-rearrangers: nest, shuffle, repeat; misc. features](https://github.com/johnkerl/miller/releases/tag/v3.5.0) |
3.4.0     |                                                                     | [JSON, reshape, regex captures, and more](https://github.com/johnkerl/miller/releases/tag/v3.4.0) |
3.3.2     |                                                                     | [Bootstrap sampling, EWMA, merge-fields, isnull/isnotnull functions](https://github.com/johnkerl/miller/releases/tag/v3.3.2) |
3.2.2     |                                                                     | [Performance improvements, compressed I/O, and variable-name escaping](https://github.com/johnkerl/miller/releases/tag/v3.2.2) |
3.1.2     |                                                                     | [Bugfix for stats1 max](https://github.com/johnkerl/miller/releases/tag/v3.1.2) |
3.1.1     |                                                                     | [Fix regression tests for i386](https://github.com/johnkerl/miller/releases/tag/v3.1.1) |
3.1.0     |                                                                     | [Minor feature enhancements, and portability](https://github.com/johnkerl/miller/releases/tag/v3.1.0) |
3.0.1     | [Miller 3.0.1](https://johnkerl.org//miller-docs-by-release/3.0.1/) | [Allow scientific notation in DSL literals; mlr bar --auto](https://github.com/johnkerl/miller/releases/tag/v3.0.1) |
3.0.0     |                                                                     | [Integer and float arithmetic, improved documentation, minor feature enhancements](https://github.com/johnkerl/miller/releases/tag/v3.0.0) |
2.3.2     |                                                                     | [Iterative stats, exclude-filter, implicit-CSV-header, and other features](https://github.com/johnkerl/miller/releases/tag/v2.3.2) |
2.3.1     |                                                                     | [Bug fix for mlr top -a](https://github.com/johnkerl/miller/releases/tag/v2.3.1) |
2.3.0     |                                                                     | [Regex support, gsub, reservoir sampling, iterative stats, and other features](https://github.com/johnkerl/miller/releases/tag/v2.3.0) |
2.2.1     |                                                                     | [Autoconfig support](https://github.com/johnkerl/miller/releases/tag/v2.2.1) |
2.2.0     |                                                                     | [Multi-character RS,FS,PS](https://github.com/johnkerl/miller/releases/tag/v2.2.0) |
2.1.4     |                                                                     | [Improved read performance for RFC4180 CSV](https://github.com/johnkerl/miller/releases/tag/v2.1.4) |
2.1.4     |                                                                     | [Improved read performance for RFC4180 CSV](https://github.com/johnkerl/miller/releases/tag/v2.1.4) |
2.1.3     |                                                                     | [Reduce tar-file size](https://github.com/johnkerl/miller/releases/tag/v2.1.3) |
2.1.1     |                                                                     | [Incremental read-performance increase for CSV format](https://github.com/johnkerl/miller/releases/tag/v2.1.1) |
2.1.0     |                                                                     | [Minor enhancements and bug fixes](https://github.com/johnkerl/miller/releases/tag/v2.1.0) |
2.0.0     | [Miller 2.0.0](https://johnkerl.org//miller-docs-by-release/2.0.0/) | [RFC4180-compliant CSV](https://github.com/johnkerl/miller/releases/tag/v2.0.0) |
1.0.1     |                                                                     | [Add INSTALLDIR Makefile option for Homebrew](https://github.com/johnkerl/miller/releases/tag/v1.0.1) |
1.0.0     | [Miller 1.0.0](https://johnkerl.org//miller-docs-by-release/1.0.0/) | [Initial public release](https://github.com/johnkerl/miller/releases/tag/v1.0.0) |
