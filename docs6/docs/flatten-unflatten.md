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

Miller's [other formats](file-formats.md), though, such as CSV, are all non-nested -- a cell
in a CSV row can't contain another entire row. As we'll see in this section, there are two main
ways to **flatten** nested data structures down to individual CSV cells -- either by _key-spreading_
(which is the default), or by _JSON-stringifying):

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

## Flattening arrays from JSON to non-JSON

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

## Unflattening maps from non-JSON to JSON

Miller's default unflattening behavior from non-JSON to JSON formats is the opposite of the flattening
behavior:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values-spread.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b.x,b.y
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson cat data/map-values-spread.csv</b>
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

Here too the `--flatsep` flag can be used to specify the separator in the data if it's not the default `.`:

<pre class="pre-highlight-in-pair">
<b>cat data/map-values-spread-colon.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b:x,b:y
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --flatsep : cat data/map-values-spread-colon.csv</b>
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

## Unflattening arrays from non-JSON to JSON

Arrays are unflattened similarly:

TODO: check why auto-infer is NOT happening :(

<pre class="pre-highlight-in-pair">
<b>cat data/array-values-spread.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b.1,b.2
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson cat data/array-values-spread.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {
    "1": 2,
    "2": 3
  }
}
{
  "a": 4,
  "b": {
    "1": 5,
    "2": 6
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/array-values-spread-colon.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b:1,b:2
1,2,3
4,5,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --flatsep : cat data/array-values-spread-colon.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {
    "1": 2,
    "2": 3
  }
}
{
  "a": 4,
  "b": {
    "1": 5,
    "2": 6
  }
}
</pre>

xxx mention (aftre fixing!) the heuristic: if keys 1..n all integerable, and consecutive, then infer
array. else map. w/ genmd'ed example.

## TODO: xxx simple examples w/o flatten/unflatten

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
<b>mlr --j2c --no-auto-flatten cat data/map-values.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b
1,"{
  ""x"": 2,
  ""y"": 3
}"
4,"{
  ""x"": 5,
  ""y"": 6
}"
</pre>

and back

x 2 for the back -- w/ json-parse

## TODO

TODO: try out in-DSL array-create (e.g. splita) for CSV to CSV ...

## TODO

TODO: manual control w/ f/uf verb/func

```
$ mlr --csv put '$z=[1,2,3]' example.csv
color,shape,flag,k,index,quantity,rate,z.1,z.2,z.3
yellow,triangle,true,1,11,43.6498,9.8870,1,2,3
red,square,true,2,15,79.2778,0.0130,1,2,3
```
