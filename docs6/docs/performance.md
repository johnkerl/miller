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
# Performance

## Disclaimer

In a previous version of this page, I compared Miller to some items in the Unix toolkit in terms of run time. But such comparisons are very much not apples-to-apples:

* Miller's principal strength is that it handles **key-value data in various formats** while the system tools **do not**. So if you time `mlr sort` on a CSV file against system `sort`, it's not relevant to say which is faster by how many percent -- Miller will respect the header line, leaving it in place, while the system sort will move it, sorting it along with all the other header lines. This would be comparing the run times of two programs produce different outputs.  Likewise, `awk` doesn't respect header lines, although you can code up some CSV-handling using `if (NR==1) { ... } else { ... }`. And that's just CSV: I don't know any simple way to get `sort`, `awk`, etc. to handle DKVP, JSON, etc. -- which is the main reason I wrote Miller.

* **Implementations differ by platform**: one `awk` may be fundamentally faster than another, and `mawk` has a very efficient bytecode implementation -- which handles positionally indexed data far faster than Miller does.

* The system `sort` command will, on some systems, handle too-large-for-RAM datasets by spilling to disk; Miller (as of version 5.2.0, mid-2017) does not. Miller sorts are always stable; GNU supports stable and unstable variants.

* Etc.

## Summary

Miller can do many kinds of processing on key-value-pair data using elapsed time roughly of the same order of magnitude as items in the Unix toolkit can handle positionally indexed data. Specific results vary widely by platform, implementation details, and multi-core use (or not). Lastly, specific special-purpose non-record-aware processing will run far faster if implemented in `grep`, `sed`, etc.

## Some examples

This is some data from [https://community.opencellid.org](https://community.opencellid.org): approximately 40
million records, 0.9GB compressed, 3.2GB uncommpressed.

First we see that decompression is much cheaper than compression: 10 seconds vs. 2.5 minutes:

```
$ time gunzip cell_towers.csv.gz

real    0m10.431s
user    0m9.235s
sys     0m1.030s

$ ls -lh cell_towers.csv
-rw-r--r--  1 johnkerl  staff   3.2G Sep  8 17:13 cell_towers.csv

$ time gzip cell_towers.csv

real    2m30.171s
user    2m28.508s
sys     0m1.257s

$ ls -lh cell_towers.csv.gz
-rw-r--r--  1 johnkerl  staff   917M Sep  7 12:34 cell_towers.csv.gz

```

Next we look at system `cut` which needs to split on lines and fields. Since `cut` is in the
[Unix toolkit](unix-toolkit-context.md) it handles integer column names, starting with 1.

```
$ gunzip < cell_towers.csv.gz | head -n 5
radio,mcc,net,area,cell,unit,lon,lat,range,samples,changeable,created,updated,averageSignal
UMTS,262,2,801,86355,0,13.285512,52.522202,1000,7,1,1282569574,1300155341,0
GSM,262,2,801,1795,0,13.276907,52.525714,5716,9,1,1282569574,1300155341,0
GSM,262,2,801,1794,0,13.285064,52.524,6280,13,1,1282569574,1300796207,0
UMTS,262,2,801,211250,0,13.285446,52.521744,1000,3,1,1282569574,1299466955,0
UMTS,262,2,801,86353,0,13.293457,52.521515,1000,2,1,1282569574,1291380444,0
```

This takes about a minute and a half:

```
$ time cut -d, -f 1,2,12,13 cell_towers.csv > /dev/null

real    1m29.228s
user    1m26.347s
sys     0m2.426s
```

Columns `1,2,12,13` are the same as `radio,mcc,created,updated`. Since
decompression is quick, it's perhaps unsurprising that whether we decompress
and have Miller read uncompressed data, or have it [decompress
in-process](reference-main-compressed-data.md#automatic-detection-on-input), or
use an [external decompressor with
`--prepipe`](reference-main-compressed-data.md#external-decompressors-on-input),
the results are about the same. This is not as fast as `cut`, but it's in the ballpark.

```

$ gunzip cell_towers.csv.gz
$ time mlr --csv cut -f radio,mcc,created,updated cell_towers.csv > /dev/null

real    4m10.097s
user    8m56.975s
sys     4m40.046s

$ gzip cell_towers.csv
$ time mlr --csv cut -f radio,mcc,created,updated cell_towers.csv.gz > /dev/null

real    4m14.185s
user    9m5.044s
sys     4m23.886s

$ time mlr --csv --prepipe gunzip cut -f radio,mcc,created,updated cell_towers.csv.gz > /dev/null

real    4m13.614s
user    9m5.623s
sys     4m57.827s

```


