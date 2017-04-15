This is a relatively minor release, containing feature requests and bugfixes while I've been working on the Windows port (which is nearly complete).

**Features:**
* **JSON arrays**: as described [here](http://johnkerl.org/miller-releases/miller-5.1.0/doc/file-formats.html#Tabular_JSON), Miller being a tabular data processor can't and won't support arbitrary JSON. But as of 5.1.0, arrays are converted to maps with integer keys, which are then at least processable using Miller.

xxx

```
--json-map-arrays-on-input    JSON arrays are unmillerable. --json-map-arrays-on-input
--json-skip-arrays-on-input   is the default: arrays are converted to integer-indexed
--json-fatal-arrays-on-input  maps. The other two options cause them to be skipped, or
```

* The new [**mlr fraction**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#fraction) verb makes possible in a few keystrokes what was only possible before using two-pass DSL logic: here you can turn numerical values down a column into their fractional/percentage contribution to column totals, optionally grouped by other key columns.

* The DSL functions
[**strptime**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#strptime) and
[**strftime**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#strftime) DSL functions now
handle fractional seconds. For parsing, use <b>%S</b> format as always; for formatting, there are now
<b>%1S</b> through
<b>%9S</b> which allow you to configure a specified number of decimal places.  The return value from <tt>strptime</tt>
is now floating-point, not integer, which is a minor backward incompatibility not worth labeling this release as 6.0.0.
You can work around this using <tt>int(strptime(...))</tt>.  The DSL functions
[**gmt2sec**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#gmt2sec) and
[**sec2gmt**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#sec2gmt), which are
keystroke-savers for <tt>strptime</tt> and <tt>strftime</tt>, are similarly modified, as is the
[**sec2gmt**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#sec2gmt) verb.
This resolves https://github.com/johnkerl/miller/issues/125.

* A few nearly-standalone programs -- which do not have anything to do with record streams -- are packaged within the Miller. (For example, hex-dump, unhex, and show-line-endings commands.) These are described [**here**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference.html#Auxiliary_commands).

* The [**stats1**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#stats1)
and [**merge-fields**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#merge-fields)
verbs now support an **antimode** aggregator, in addition to the existing mode aggregator.

* The [**join**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#join) verb
now by default does not require sorted input, which is the more common use case. (Memory-parsimonious joins which require sorted input, while no longer the default, are available using <tt>-s</tt>.) This another minor backward incompatibility not worth making a 6.0.0 over.
This resolves https://github.com/johnkerl/miller/issues/134.

* [**mlr nest**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-verbs.html#nest) has a keystroke-saving <tt>--evar</tt> option for a common use case, namely, exploding a field by value across records.

**Documentation:**
* The DSL reference now has [**per-function descriptions**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/reference-dsl.html#Built-in_functions_for_filter_and_put).
* There is a new [**feature-counting example**](http://johnkerl.org/miller-releases/miller-5.1.0/doc/cookbook.html#Feature-counting) in the cookbook.

**Bugfixes:**
* **mlr join -j -l** was not functioning correctly. This resolves https://github.com/johnkerl/miller/issues/136.

* **JSON escapes on output** (`\t` and so on) were incorrect. This resolves
https://github.com/johnkerl/miller/issues/135.

