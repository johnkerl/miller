This is a relatively minor release, containing feature requests.

**Features:**

* There is a new DSL functions [**mapexcept**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#mapexcept) which returns a copy of the argument with speicified key(s), if any, unset. This resolves https://github.com/johnkerl/miller/issues/137.

**Bugfixes:**

* CRLF line-endings were not being correctly autodetected when I/O formats were specified using <tt>--c2j</tt> et al.
