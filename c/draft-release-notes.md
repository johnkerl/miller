# Title here

## Features:

* The new [**--jofmt**](http://johnkerl.org/miller/doc/file-formats.html#Formatting_JSON_options) option allows you to apply
the global <tt>--ofmt</tt> to floating-point values in JSON output &mdash; with caveats as noted [**here**](http://johnkerl.org/miller/doc/file-formats.html#Formatting_JSON_options).

## Documentation:

* ...

## Bugfixes: 

* The [**toupper**](http://johnkerl.org/miller/doc/reference-dsl.html#toupper), [**tolower**](http://johnkerl.org/miller/doc/reference-dsl.html#tolower), and [**capitalize**](http://johnkerl.org/miller/doc/reference-dsl.html#capitalize) DSL functions are now UTF-8 aware, thanks to @sheredom's marvelous https://github.com/sheredom/utf8.h. The [**internationalization page**]((http://johnkerl.org/miller/doc/internationalization.html) has also been expanded.

* https://github.com/johnkerl/miller/issues/250 fixes a bug using [**in-place mode**](https://johnkerl.org/miller/doc/reference.html#In-place_mode) in conjunction with verbs (such as [**rename**](http://johnkerl.org/miller/doc/reference-dsl.html#rename) or [**sort**](http://johnkerl.org/miller/doc/reference-dsl.html#sort)) which take field-name lists as arguments.

* https://github.com/johnkerl/miller/issues/253 fixes a bug in the [**label**](http://johnkerl.org/miller/doc/reference-verbs.html#label) when one or more names are common between old and new.

* https://github.com/johnkerl/miller/issues/251 fixes a corner-case bug when (a) input is CSV; (b) the last field ends with a comma and no newline; (c) input is from standard input and/or <tt>--no-mmap</tt> is supplied.
