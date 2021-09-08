<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
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
# Maps

Miller data types are listed on the [Data types](reference-main-data-types.md) page; here we focus specifically on maps.

On the whole, maps are as in most other programming languages. However, following the
[Principle of Least Surprise](https://en.wikipedia.org/wiki/Principle_of_least_astonishment)
and aiming to reduce keystroking for Miller's most-used streaming-record-processing model,
there are a few differences as noted below.

## Types of maps

_Map literals_ are written in curly braces with string keys any [Miller data type](reference-main-data-types.md) (including other maps, or arrays) as values. Also, integers may be given as keys although they'll be stored as strings.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = {"a": 1, "b": {"x": 2, "y": [3,4,5]}, 99: true};</b>
<b>    dump x;</b>
<b>    print x[99];</b>
<b>    print x["99"];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": {
    "x": 2,
    "y": [3, 4, 5]
  },
  "99": true
}
true
true
</pre>

As with arrays and argument-lists, trailing commas are supported:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = {</b>
<b>      "a" : 1,</b>
<b>      "b" : 2,</b>
<b>      "c" : 3,</b>
<b>    };</b>
<b>    print x;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": 2,
  "c": 3
}
</pre>

The current record, accessible using `$*`, is a map.

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv head -n 2 then put -q '</b>
<b>  dump $*;</b>
<b>  print "Color is", $*["color"];</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "shape": "triangle",
  "flag": "true",
  "k": 1,
  "index": 11,
  "quantity": 43.6498,
  "rate": 9.8870
}
Color is yellow
{
  "color": "red",
  "shape": "square",
  "flag": "true",
  "k": 2,
  "index": 15,
  "quantity": 79.2778,
  "rate": 0.0130
}
Color is red
</pre>

The collection of all [out-of-stream variables](reference-dsl-variables.md#out-of-stream0variables), `@*`, is a map.

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv put -q '</b>
<b>  begin {</b>
<b>    @last_rates = {};</b>
<b>  }</b>
<b>  @last_rates[$shape] = $rate;</b>
<b>  @last_color = $color;</b>
<b>  end {</b>
<b>    dump @*;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "last_rates": {
    "triangle": 5.8240,
    "square": 8.2430,
    "circle": 8.3350
  },
  "last_color": "purple"
}
</pre>

Also note that several [built-in functions](reference-dsl-builtin-functions.md) operate on maps and/or return maps.

## Insertion order is preserved

Miller maps preserve insertion order. So if you write `@m["y"]=7` and then `@m["x"]=3` then any loop over
the map `@m` will give you the kays `"y"` and `"x"` in that order.

## String keys, with conversion from/to integer

All Miller map keys are strings. If a map is indexed with an integer for either
read or write (i.e. on either the right-hand side or left-hand side of an
assignment) then the integer will be converted to/from string, respectively. So
`@m[3]` is the same as `@m["3"]`. The reason for this is for situations like
[operating on all records](operating-on-all-records.md) where it's important to
let people do `@records[NR] = $*`.

## Auto-create

Indexing any as-yet-assigned local variable or out-of-stream variable results
in **auto-create** of that variable as a map variable:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv put -q '</b>
<b>  # You can do this but you do not need to:</b>
<b>  # begin { @last_rates = {} }</b>
<b>  @last_rates[$shape] = $rate;</b>
<b>  end {</b>
<b>    dump @last_rates;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "triangle": 5.8240,
  "square": 8.2430,
  "circle": 8.3350
}
</pre>

*This also means that auto-create results in maps, not arrays, even if keys are integers.*
If you want to auto-extend an [array](reference-main-arrays.md), initialize it explicitly to `[]`.

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv head -n 4 then put -q '</b>
<b>  begin {</b>
<b>    @my_array = [];</b>
<b>  }</b>
<b>  @my_array[NR] = $quantity;</b>
<b>  @my_map[NR] = $rate;</b>
<b>  end {</b>
<b>    dump</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "my_array": [43.6498, 79.2778, 13.8103, 77.5542],
  "my_map": {
    "1": 9.8870,
    "2": 0.0130,
    "3": 2.9010,
    "4": 7.4670
  }
}
</pre>

## Auto-deepen

Similarly, maps are **auto-deepened**: you can put `@m["a"]["b"]["c"]=3`
without first setting `@m["a"]={}` and `@m["a"]["b"]={}`. The reason for this
is for doing data aggregations: for example if you want compute keyed sums, you
can do that with a minimum of keystrokes.

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put -q '</b>
<b>  @quantity_sum[$color][$shape] += $rate;</b>
<b>  end {</b>
<b>    emit @quantity_sum, "color", "shape";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    quantity_sum
yellow triangle 9.8870
yellow circle   12.572000000000001
red    square   17.011
red    circle   2.9010
purple triangle 14.415
purple square   8.2430
</pre>

## Looping

See [single-variable for-loops](reference-dsl-control-structures.md#single-variable-for-loops) and [key-value for-loops](reference-dsl-control-structures.md#key-value-for-loops).
