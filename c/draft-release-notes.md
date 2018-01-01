## Features:

* **Comment strings in data files:** `mlr --skip-comments` allows you to
filter out input lines starting with `#`, for all file formats.  Likewise, `mlr
--skip-comments-with X` lets you specify the comment-string `X`.  Comments are
only supported at start of data line.

* The [**count-similar**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-verbs.html#count-similar)
verb lets you compute cluster sizes by cluster labels.

* While Miller DSL arithmetic gracefully overflows from 64-integer to
double-precision float (see also
[**here**](http://johnkerl.org/miller/doc/reference.html#Arithmetic), there
are now the **integer-preserving arithmetic operators** ``.+` `.-` `.*` `./` `.//`
for those times when you want integer overflow.

* There is a new [**bitcount**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#bitcount) function: for example, `echo x=0xf0000206 | mlr put '$y=bitcount($x)'` produces `x=0xf0000206,y=7`.

* [**Issue 158**](https://github.com/johnkerl/miller/issues/158): `mlr -T` is
an alias for `--nidx --fs tab`, and `mlr -t` is an alias for `mlr
--tsvlite`.

* The mathematical constants &pi; and <i>e</i> have been renamed from `PI` and `E` to `M_PI` and `M_E`, respectively. (It's annoying to get a syntax error when you try to define a variable named `E` in the DSL, when `A` through `D` work just fine.) This is a backward incompatibility, but not enough of us to justify calling this release Miller 6.0.0.

## Documentation:

* As noted
[**here**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/reference-dsl.html#A_note_on_the_complexity_of_Miller’s_expression_language), since Miller has its own DSL there will always be things better expressible in a general-purpose language. The new page
[**Sharing data with other languages**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/data-sharing.html) shows how to seamlessly share data back and forth between Miller, Ruby, and Python.

* [**SQL-input examples**](http://johnkerl.org/miller-releases/miller-5.2.0/doc/10-min.html#SQL-input_examples) contains detailed information the interplay between Miller and SQL.

* [**Issue 150**](https://github.com/johnkerl/miller/issues/150) raised a
question about suppressing numeric conversion. This resulted in a new FAQ entry
[**How do I suppress numeric conversion?**](http://johnkerl.org/miller/doc/faq.html#How_do_I_suppress_numeric_conversion?), as well as the
longer-term follow-on [**issue 151**](https://github.com/johnkerl/miller/issues/151) which will make numeric conversion happen on a just-in-time basis.

* To my surprise, **csvlite format options** weren&rsquo;t listed in `mlr --help` or the manpage. This has been fixed.

## Bugfixes:

* xxx mmap/madvise

* xxx nonesuch-lash https://github.com/johnkerl/miller/issues/162

* [**Issue 159**](https://github.com/johnkerl/miller/issues/159): fix regex-match of literal dot.
