This is a relatively minor release, containing feature requests.

**Features:**

* There is a new DSL function [**mapexcept**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#mapexcept) which returns a copy of the argument with specified key(s), if any, unset.  Likewise, [**mapselect**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#mapselect) returns a copy of the argument with only specified key(s), if any, set.  This resolves https://github.com/johnkerl/miller/issues/137.

* xxx min/max functions and stats1/merge-fields min/max/percentile mix int and string. esp. string-only order statistics. doclink for mixed case. interpolation obv nonsensical.

* xxx `./configure` vs. `autoreconf -fiv` 1st, and which issue is resolved by this.

* xxx UTF-8 BOM strip for CSV files; resolves xxx

**Documentation:**

* xxx cookbook example [**Showing differences between successive queries**](http://www.johnkerl.org/miller-releases/miller-5.2.0/doc/cookbook.html#Showing_differences_between_successive_queries)

**Bugfixes:**

* CRLF line-endings were not being correctly autodetected when I/O formats were specified using <tt>--c2j</tt> et al.

* Integer division by zero was causing a fatal runtime exception, rather than computing <tt>inf</tt> or <tt>nan</tt> as in the floating-point case.
