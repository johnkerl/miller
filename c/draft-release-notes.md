# Positional-indexing and other data-cleaning features

## Features:

* The [**positional-indexing feature**](http://johnkerl.org/miller/doc/reference-dsl.html#Positional_field_names) resolves https://github.com/johnkerl/miller/issues/236 from @aborruso. You can now get the name of the 3rd field of each record via <tt>$[[3]]</tt>, and the value by <tt>$[[[3]]]</tt>. These are both usable on either the left-hand or right-hand side of assignment statements, so you can more easily do things like renaming fields progrmatically within the DSL.

* There is a new [**capitalize**](http://johnkerl.org/miller/doc/reference-dsl.html#capitalize) DSL function, complementing the already-existing <tt>toupper</tt>. This stems from https://github.com/johnkerl/miller/issues/236.

* There is a new [**skip-trivial-records**](http://johnkerl.org/miller/doc/reference-verbs.html#skip-trivial-records) verb, resolving https://github.com/johnkerl/miller/issues/197. Similarly, there is a new [**remove-empty-columns**](http://johnkerl.org/miller/doc/reference-verbs.html#remove-empty-columns) verb, resolving https://github.com/johnkerl/miller/issues/206. Both are useful for **data-cleaning use-cases**.

* Another pair is https://github.com/johnkerl/miller/issues/181 and https://github.com/johnkerl/miller/issues/256. While Miller uses <tt>mmap</tt> internally (and invisibily) to get approximately a 20% performance boost over not using it, this can cause out-of-memory issues with reading either large files, or too many small ones. Now, Miller automatically avoids <tt>mmap</tt> in these cases. You can still use <tt>--mmap</tt> or <tt>--no-mmap</tt> if you want manual control of this.

* There is a new [**--ivar option for the nest verb**](http://johnkerl.org/miller/doc/reference-verbs.html#nest) which complements the already-existing <tt>--evar</tt>. This is from https://github.com/johnkerl/miller/pull/260 thanks to @jgreely.

* There is a new keystroke-saving [**urandrange**](http://johnkerl.org/miller/doc/reference-dsl.html#urand) DSL function: <tt>urandrange(low, high)</tt> is the same as <tt>low + (high - low) * urand()</tt>.

* There is a new [**-v option for the cat verb**](http://johnkerl.org/miller/doc/reference-verbs.html#cat) which writes a low-level record-structure dump to standard error.

* There is a new [**-N option for mlr**](http://johnkerl.org/miller/doc/manpage.html) which is a keystroke-saver for <tt>--implicit-csv-header --headerless-csv-output</tt>.

## Documentation:

* The new FAQ entry http://johnkerl.org/miller/doc/faq.html#How_to_escape_'?'_in_regexes resolves https://github.com/johnkerl/miller/issues/203.

* The new FAQ entry http://johnkerl.org/miller/doc/faq.html#How_can_I_filter_by_date resolves https://github.com/johnkerl/miller/issues/208.

* https://github.com/johnkerl/miller/issues/244 fixes a documentation issue while highlighting the need for https://github.com/johnkerl/miller/issues/241.

## Bugfixes: 

* There was a SEGV using `nest` within `then`-chains, fixed in response to https://github.com/johnkerl/miller/issues/220.

* Quotes and backslashes weren't being escaped in JSON output with <tt>--jvquoteall</tt>; reported on https://github.com/johnkerl/miller/issues/222.

## Note:

I've never code-named releases but if I were to code-name 5.5.0 I would call it "aborruso". Andrea has contributed many fantastic feature requests, as well as driving a huge volume of Miller-related discussions in StackExchange (https://github.com/johnkerl/miller/issues/212). Mille grazie al mio amico @aborruso!
