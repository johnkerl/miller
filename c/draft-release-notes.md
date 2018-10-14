# New data-cleaning features, Windows mlr.exe, limited localtime support, and bugfixes

## Features:

* The new [**clean-whitespace**](http://johnkerl.org/miller/doc/reference-verbs.html#clean-whitespace) verb resolves
https://github.com/johnkerl/miller/issues/190 from @aborruso.
Along with the new functions
[**strip**](http://johnkerl.org/miller/doc/reference-dsl.html#strip),
[**lstrip**](http://johnkerl.org/miller/doc/reference-dsl.html#lstrip),
[**rstrip**](http://johnkerl.org/miller/doc/reference-dsl.html#rstrip),
[**collapse_whitespace**](http://johnkerl.org/miller/doc/reference-dsl.html#collapse_whitespace), and
[**clean_whitespace**](http://johnkerl.org/miller/doc/reference-dsl.html#clean_whitespace), there is
coarser-grained and finer-grained control over whitespace within field names and/or values.
See the linked-to documentation for examples.

* The new [**altkv**](http://johnkerl.org/miller/doc/reference-verbs.html#altkv) verb resolves
https://github.com/johnkerl/miller/issues/184 which was originally opened via an email request. This supports mapping
value-lists such as `a,b,c,d` to alternating key-value pairs such as `a=b,c=d`.

* The new [**fill-down**](http://johnkerl.org/miller/doc/reference-verbs.html#fill-down) verb resolves
https://github.com/johnkerl/miller/issues/189
by
@aborruso
See the linked-to documentation for examples.

* The [**uniq**](http://johnkerl.org/miller/doc/reference-verbs.html#verb) verb now has a **uniq -a**
which resolves https://github.com/johnkerl/miller/issues/168 from @sjackman.

* The new
[**regextract**](http://johnkerl.org/miller/doc/reference-dsl.html#regextract) and
[**regextract_or_else**](http://johnkerl.org/miller/doc/reference-dsl.html#regextract_or_else)
functions resolve
https://github.com/johnkerl/miller/issues/183
by @aborruso.
xxx.

* The new [**ssub**](http://johnkerl.org/miller/doc/reference-dsl.html#ssub) function arises from
https://github.com/johnkerl/miller/issues/171
by @dohse, as a simplified way to avoid escaping characters which are special to regular-expression parsers.

* There are [**localtime**] functions in response to
https://github.com/johnkerl/miller/issues/170 by @sitaramc, as follows. However note that
as discussed on https://github.com/johnkerl/miller/issues/170 these do not undo one another in all
circumstances.
This is a non-issue for timezones which do not do DST. Otherwise, please use with disclaimers.
  * [**localdate**](http://johnkerl.org/miller/doc/reference-dsl.html#localdate)
  * [**localtime2sec**](http://johnkerl.org/miller/doc/reference-dsl.html#localtime2sec)
  * [**sec2localdate**](http://johnkerl.org/miller/doc/reference-dsl.html#sec2localdate)
  * [**sec2localtime**](http://johnkerl.org/miller/doc/reference-dsl.html#sec2localtime)
  * [**strftime_local**](http://johnkerl.org/miller/doc/reference-dsl.html#strftime_local)
  * [**strptime_local**](http://johnkerl.org/miller/doc/reference-dsl.html#strptime_local)

## Builds:

* Windows build-artifacts are now available in Appveyor at
https://ci.appveyor.com/project/johnkerl/miller/build/artifacts, and will be attached to this and future releases. This
reseolvs https://github.com/johnkerl/miller/issues/167, https://github.com/johnkerl/miller/issues/148, and
https://github.com/johnkerl/miller/issues/109.

* Travis builds at https://travis-ci.org/johnkerl/miller/builds now run on OSX as well as Linux.

* An Ubuntu 17 build issue was fixed by @singalen on https://github.com/johnkerl/miller/issues/164.

## Documentation:

* <tt>put</tt>/<tt>filter</tt> documentation was confusing as reported by @NikosAlexandris on
https://github.com/johnkerl/miller/issues/169.

* The new FAQ entry
http://johnkerl.org/miller-releases/miller-head/doc/faq.html#How_to_rectangularize_after_joins_with_unpaired?
resolves
https://github.com/johnkerl/miller/issues/193
by @aborruso.

* The new cookbook entry
http://johnkerl.org/miller/doc/cookbook.html#Options_for_dealing_with_duplicate_rows arises from
https://github.com/johnkerl/miller/issues/168 from @sjackman.

* The <tt>unsparsify</tt> documentation had some words missing as reported by
@tst2005 on https://github.com/johnkerl/miller/issues/194.

* There was a typo in the cookpage page http://johnkerl.org/miller/doc/cookbook.html#Full_field_renames_and_reassigns
as fixed by @tst2005 in https://github.com/johnkerl/miller/pull/192.

## Bugfixes: 

* There was a memory leak for TSV-format files only as reported by @treynr on https://github.com/johnkerl/miller/issues/181.

* Dollar sign in regular expressions were not being escaped properly as reported by @dohse on
https://github.com/johnkerl/miller/issues/171.
