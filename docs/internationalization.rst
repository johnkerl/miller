..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Internationalization
================================================================

Miller handles strings with any characters other than 0x00 or 0xff, using explicit UTF-8-friendly string-length computations.  (I have no plans to support UTF-16 or ISO-8859-1.)

By and large, Miller treats strings as sequences of non-null bytes without need to interpret them semantically. Intentional support for internationalization includes:

* <a href="file-formats.html">Tabular output formats</a> (pprint and xtab) are aligned correctly.
* The <a href="reference-dsl.html#strlen">strlen</a> function correctly counts UTF-8 codepoints rather than bytes.
* The <a href="reference-dsl.html#toupper">toupper</a>, <a href="reference-dsl.html#tolower">tolower</a>, and <a href="reference-dsl.html#capitalize">capitalize</a> DSL functions within the capabilities of <a href="https://github.com/sheredom/utf8.h">https://github.com/sheredom/utf8.h</a>.

Meanwhile, regular expressions and the <a href="reference-dsl.html#sub">sub</a> and <a href="reference-dsl.html#gsub">gsub</a> function correctly, albeit without explicit intentional support.

Please file an issue at <a href="https://github.com/johnkerl/miller">https://github.com/johnkerl/miller</a> if you encounter bugs related to internationalization (or anything else for
that matter).
