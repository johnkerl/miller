This is a relatively minor release, containing feature requests.

**Features:**

* There is a new DSL function [**mapexcept**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#mapexcept) which returns a copy of the argument with specified key(s), if any, unset.  Likewise, [**mapselect**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#mapselect) returns a copy of the argument with only specified key(s), if any, set.  This resolves https://github.com/johnkerl/miller/issues/137.

* xxx min/max functions and stats1/merge-fields min/max/percentile mix int and string. esp. string-only order statistics. doclink for mixed case. interpolation obv nonsensical.

* A new **-u** option for [**count-distinct**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-verbs.html#count-distinct) allows unlashed counts for multiple field names. For example, with `-f a,b` and
without `-u`, `count-distinct` computes counts for distinct pairs of `a` and `b` field values. With `-f a,b` and with `-u`, it computes counts
for distinct `a` field values and counts for distinct `b` field values separately.

* xxx `./configure` vs. `autoreconf -fiv` 1st, and which issue is resolved by this.

* xxx UTF-8 BOM strip for CSV files; resolves xxx

* For `put` and `filter` with `-S`, program literals such as the `6` in `$x = 6` were being parsed as strings. This is not sensible, since the `-S` option for `put` and `filter` is intended to suppress numeric conversion of record data, not program literals. To get string `6` one may use `$x = "6"`.

**Documentation:**

* Suppose you have counters in a SQL database with different values in successive queries.  A new cookbook example shows [**how to compute differences between successive queries**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook.html#Showing_differences_between_successive_queries).

* Another new cookbook example shows [**how to compute interquartile ranges**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook2.html#Computing_interquartile_ranges)

**Bugfixes:**

* CRLF line-endings were not being correctly autodetected when I/O formats were specified using <tt>--c2j</tt> et al.

* Integer division by zero was causing a fatal runtime exception, rather than computing <tt>inf</tt> or <tt>nan</tt> as in the floating-point case.
