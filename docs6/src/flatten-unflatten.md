<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# Flatten/unflatten: converting between JSON and tabular formats

Miller has long supported reading and writing multiple [file
formats](file-formats.md) including CSV and JSON, as well as converting back
and forth between them. Two things new in [Miller 6](new-in-miller-6-md),
though, are that [arrays are now fully supported](reference-main-arrays.md),
and that [record values are typed](new-in-miller-6.md#improved-numeric-conversion)
throughout Miller's processing chain from input through [verbs](reference-verbs.md)
to output -- which includes improved handling for [maps](reference-main-maps.md) and
[arrays](reference-main-arrays.md) as record values.

This raises the question, though, of how to handle maps and arrays as record values.
For [JSON files](file-formats.md#json), this is easy -- JSON is a nested format where values
can be maps or arrays, which can contain other maps or arrays, and so on, with the nesting
happily indicated by curly braces:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {"x": 2, "y": 3}
}
{
  "a": 4,
  "b": {"x": 5, "y": 6}
}
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/map-values-nested.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {"s": {"w": 2, "x": 3}, "t": {"y": 4, "z": 5}}
}
{
  "a": 6,
  "b": {"s": {"w": 7, "x": 8}, "t": {"y": 9, "z": 10}}
}
</pre>

How can we represent these in CSV files?

Miller's [non-JSON formats](file-formats.md), such as CSV, are all non-nested -- a
cell in a CSV row can't contain another entire row. As we'll see in this
section, there are two main ways to **flatten** nested data structures down to
individual CSV cells -- either by _key-spreading_ (which is the default), or by
_JSON-stringifying_:

* **Key-spreading** is when the single map-valued field
`b={"x": 2, "y": 3}` spreads into multiple fields `b.x=2,b.y=3`;
* **JSON-stringifying** is when the single map-valued field `"b": {"x": 2, "y": 3}` becomes the single string-valued field `b="{\"x\":2,\"y\":3}"`.

Miller intends to provide intuitive default behavior for these conversions, while also
providing you with more control when you need it.

## Converting maps between JSON and non-JSON

Let's first look at the default behavior with map-valued fields. Miller's
default behavior is to spread the map values into multiple keys -- using
Miller's `flatsep` separator, which defaults to `.` -- to join the original
record key with map keys:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {"x": 2, "y": 3}
}
{
  "a": 4,
  "b": {"x": 5, "y": 6}
}
</pre>

Flattened to CSV format:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b.x,b.y
1,2,3
4,5,6
</pre>

Flattened to pretty-print format:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b.x b.y
1 2   3
4 5   6
</pre>

Using flatten-separator `:` instead of the default `.`:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint --flatsep : cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b:x b:y
1 2   3
4 5   6
</pre>

If the maps are more deeply nested, each level of map keys is joined in:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values-nested.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {"s": {"w": 2, "x": 3}, "t": {"y": 4, "z": 5}}
}
{
  "a": 6,
  "b": {"s": {"w": 7, "x": 8}, "t": {"y": 9, "z": 10}}
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/map-values-nested.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b.s.w b.s.x b.t.y b.t.z
1 2     3     4     5
6 7     8     9     10
</pre>

**Unflattening** is simply the reverse -- from non-JSON back to JSON:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {"x": 2, "y": 3}
}
{
  "a": 4,
  "b": {"x": 5, "y": 6}
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b.x,b.y
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/map-values.json | mlr --icsv --ojson cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {
    "x": 2,
    "y": 3
  }
}
{
  "a": 4,
  "b": {
    "x": 5,
    "y": 6
  }
}
</pre>

## Converting arrays between JSON and non-JSON

If the input data contains arrays, these are also flattened similarly: the
[1-up array indices](reference-main-arrays.md#1-up-indexing) `1,2,3,...` become string keys
`"1","2","3",...`:

<pre class="pre-highlight-in-pair">
<b>cat data/array-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": [2, 3]
}
{
  "a": 4,
  "b": [5, 6]
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/array-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b.1 b.2
1 2   3
4 5   6
</pre>

If the arrays are more deeply nested, each level of arrays keys is joined in:

<pre class="pre-highlight-in-pair">
<b>cat data/array-values-nested.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": [[2, 3], [4, 5]]
}
{
  "a": 6,
  "b": [[7, 8], [9, 10]]
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/array-values-nested.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b.1.1 b.1.2 b.2.1 b.2.2
1 2     3     4     5
6 7     8     9     10
</pre>

In the nested-data examples shown here, nested map values are shown containing
maps, and nested array values are shown containing arrays -- of course (even
though not shown here) nested map values can contain arrays, and vice versa.

**Unflattening** arrays is, again, simply the reverse -- from non-JSON back to JSON:

<pre class="pre-highlight-in-pair">
<b>cat data/array-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": [2, 3]
}
{
  "a": 4,
  "b": [5, 6]
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/array-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b.1,b.2
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv cat data/array-values.json | mlr --icsv --ojson cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": [2, 3]
}
{
  "a": 4,
  "b": [5, 6]
}
</pre>

## Auto-inferencing of arrays on unflatten

Note that the CSV field names `b.x` and `b.y` aren't too different from `b.1`
and `b.2`.  Miller has the heuristic that if it's unflattening and gets a map
with keys `"1"`, `"2"`, etc.  -- starting with `"1"`, consecutively, and with
no gaps -- it turns that back into an array.  This is precisely to undo the
flatten conversion. However, it may (or may not) be surprising:

<pre class="pre-highlight-in-pair">
<b>cat data/consecutive.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a.1,a.2,a.3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2j cat data/consecutive.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": [4, 5, 6]
}
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/non-consecutive.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a.1,a.3,a.5
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2j cat data/non-consecutive.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": {
    "1": 4,
    "3": 5,
    "5": 6
  }
}
</pre>

## TODO: w/o defaults

* list the auto-flatten/a-uf rules & what they do in terms of appending to the chain.
* also: csv-to-csv w/ $z=[1,2,3] (or, better, splita) example
* show behavior w/ these being used.
* show json-parse, json-stringify.
* xref to JSON-in-CSV

