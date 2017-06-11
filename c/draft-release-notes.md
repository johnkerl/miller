This release contains mostly feature requests.

**Features:**

* There is a new DSL function
[**mapexcept**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#mapexcept) which returns a
copy of the argument with specified key(s), if any, unset.  The motivating use-case is to split records to multiple
filenames depending on particular field value, which is omitted from the output: `mlr --from f.dat put 'tee >
"/tmp/data-".$a, mapexcept($*, "a")'` Likewise,
[**mapselect**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#mapselect) returns a copy of the
argument with only specified key(s), if any, set.  This resolves https://github.com/johnkerl/miller/issues/137.

* The [**min**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#min)
and [**max**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#max) DSL functions, and the
min/max/percentile aggregators for the
[**stats1**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-verbs.html#stats1) and
[**merge-fields**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-verbs.html#merge-fields) verbs, now
support numeric as well as string field values. (For mixed string/numeric fields, numbers compare before strings.) This
means in particular that order statistics are now possible on string-only fields. Interpolation is obviously nonsensical
for strings, so interpolated percentiles such as `mlr stats1 -a p50 -f a -i` yields an error for string-only fields.
Likewise, any other aggregations requiring arithmetic, such as <tt>mean</tt>, also produce an error on string-valued
input.

* A new **-u** option for [**count-distinct**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-verbs.html#count-distinct) allows unlashed counts for multiple field names. For example, with `-f a,b` and
without `-u`, `count-distinct` computes counts for distinct pairs of `a` and `b` field values. With `-f a,b` and with `-u`, it computes counts
for distinct `a` field values and counts for distinct `b` field values separately.

* If you [build from source](http://johnkerl.org/miller-releases/miller-5.2.0/doc/build.html), you can now
do `./configure` without first doing `autoreconf -fiv`. This resolves https://github.com/johnkerl/miller/issues/xxx.
**xxx to do**: figure out and fix the timestamp issue.
**xxx to do**: update the build.html page.

* The UTF-8 BOM sequence `0xef` `0xbb` `0xbf` is now automatically ignored from the start of CSV files. (The same is
already done for JSON files.) This resolves https://github.com/johnkerl/miller/issues/xxx.

* For `put` and `filter` with `-S`, program literals such as the `6` in `$x = 6` were being parsed as strings. This is not sensible, since the `-S` option for `put` and `filter` is intended to suppress numeric conversion of record data, not program literals. To get string `6` one may use `$x = "6"`.

**Documentation:**

* A new cookbook example shows [**how to compute differences between successive
queries**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook.html#Showing_differences_between_successive_queries),
e.g. to find out what changed in time-varying data when you run and rerun a SQL query.

* Another new cookbook example shows [**how to compute interquartile ranges**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook2.html#Computing_interquartile_ranges).

* A third new cookbook example shows [**how to compute weighted means**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook2.html#Computing_weighted_means).

**Bugfixes:**

* CRLF line-endings were not being correctly autodetected when I/O formats were specified using <tt>--c2j</tt> et al.

* Integer division by zero was causing a fatal runtime exception, rather than computing <tt>inf</tt> or <tt>nan</tt> as in the floating-point case.
