<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# Installing Miller

You can install Miller for various platforms as follows.

* Miller 6 is in pre-release, and is described by the docs you're reading ([https://johnkerl.org/miller6](https://johnkerl.org/miller6)).
    * You can get latest Miller 6 builds for Linux, MacOS, and Windows by visiting [https://github.com/johnkerl/miller/actions](https://github.com/johnkerl/miller/actions), selecting the latest build, and clicking _Artifacts_. (These are retained for 5 days after each commit.)
    * See also the [build page](build.md) if you prefer -- in particular, if your platform's package manager doesn't have the latest release.
* Miller 5 is released, and is described by [https://miller.readthedocs.io](https://miller.readthedocs.io).
    * Linux: `yum install miller` or `apt-get install miller` depending on your flavor of Linux, or [Homebrew](https://docs.brew.sh/linux).
    * MacOS: `brew install miller` or `port install miller` depending on your preference of [Homebrew](https://brew.sh) or [MacPorts](https://macports.org).
    * Windows: `choco install miller` using [Chocolatey](https://chocolatey.org).

As a first check, you should be able to run `mlr --version` at your system's command prompt and see something like the following:

<pre class="pre-highlight-in-pair">
<b>mlr --version</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Miller v6.0.0-dev
</pre>

As a second check, given [example.csv](./example.csv) you should be able to do

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
purple,triangle,false,5,51,81.2290,8.5910
red,square,false,6,64,77.1991,9.5310
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow triangle true  1  11    43.6498  9.8870
red    square   true  2  15    79.2778  0.0130
red    circle   true  3  16    13.8103  2.9010
red    square   false 4  48    77.5542  7.4670
purple triangle false 5  51    81.2290  8.5910
red    square   false 6  64    77.1991  9.5310
purple triangle false 7  65    80.1405  5.8240
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
purple square   false 10 91    72.3735  8.2430
</pre>

If you run into issues on these checks, please check out the resources on the [community page](community.md) for help.

Otherwise, let's go on to [Miller in 10 minutes](10min.md)!
