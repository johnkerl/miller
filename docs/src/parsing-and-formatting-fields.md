<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Parsing and formatting fields

Miller offers several ways to split strings into pieces (parsing them), and to put things together
into a string (formatting them).

## Splitting and joining with the same separator

One pattern we often have is items separated by the same separator, e.g. a field with value
`1;2;3;4` -- with a `;` between every pair of items. There are several useful
[DSL](miller-programming-language.md) [functions](reference-dsl-builtin-functions.md) for splitting
a string into pieces, and joining pieces into a string.

For example, suppose we have a CSV file like this:

<pre class="pre-highlight-in-pair">
<b>cat data/split1.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name,nicknames,codes
Alice,"Allie,Skater","1,3,5"
Robert,"Bob,Bobby,Biker","2,4,6"
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson cat data/split1.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": "Allie,Skater",
  "codes": "1,3,5"
},
{
  "name": "Robert",
  "nicknames": "Bob,Bobby,Biker",
  "codes": "2,4,6"
}
]
</pre>

Then we can use the [`splita`](reference-dsl-builtin-functions.md#splita) function to split the
`nicknames` string into an array of strings:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split1.csv put '$nicknames = splita($nicknames, ",")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": ["Allie", "Skater"],
  "codes": "1,3,5"
},
{
  "name": "Robert",
  "nicknames": ["Bob", "Bobby", "Biker"],
  "codes": "2,4,6"
}
]
</pre>

Likewise we can split the `codes` field. Since these look like numbers, we can again use `splita`
which tries to type-infer ints and floats when it finds them -- or, we can use
[splitax](reference-dsl-builtin-functions.md#splitax) to ask for the string to be split up into
substrings, with no type inference:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split1.csv put '$codes = splita($codes, ",")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": "Allie,Skater",
  "codes": [1, 3, 5]
},
{
  "name": "Robert",
  "nicknames": "Bob,Bobby,Biker",
  "codes": [2, 4, 6]
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split1.csv put '$codes = splitax($codes, ",")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": "Allie,Skater",
  "codes": ["1", "3", "5"]
},
{
  "name": "Robert",
  "nicknames": "Bob,Bobby,Biker",
  "codes": ["2", "4", "6"]
}
]
</pre>

We can do operations on the array, then use [joinv](reference-dsl-builtin-functions.md#joinv) to put them
back together:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split1.csv put '</b>
<b>  $codes = splita($codes, ",");                       # split into array of integers</b>
<b>  $codes = apply($codes, func(e) { return e * 100 }); # do math on the array of integers</b>
<b>  $codes = joinv($codes, ",");                        # join the updated array back into a string</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": "Allie,Skater",
  "codes": "100,300,500"
},
{
  "name": "Robert",
  "nicknames": "Bob,Bobby,Biker",
  "codes": "200,400,600"
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/split1.csv put '</b>
<b>  $codes = splita($codes, ",");                       # split into array of integers</b>
<b>  $codes = apply($codes, func(e) { return e * 100 }); # do math on the array of integers</b>
<b>  $codes = joinv($codes, ",");                        # join the updated array back into a string</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name,nicknames,codes
Alice,"Allie,Skater","100,300,500"
Robert,"Bob,Bobby,Biker","200,400,600"
</pre>

The full list of split functions includes
[splita](reference-dsl-builtin-functions.md#splita),
[splitax](reference-dsl-builtin-functions.md#splitax),
[splitkv](reference-dsl-builtin-functions.md#splitkv),
[splitkvx](reference-dsl-builtin-functions.md#splitkvx),
[splitnv](reference-dsl-builtin-functions.md#splitnv), and
[splitnx](reference-dsl-builtin-functions.md#splitx). The flavors have to to with what the output is
-- arrays or maps -- and whether or not type-inference is done.

The full list of join functions includes [joink](reference-dsl-builtin-functions.md#joink),
[joinv](reference-dsl-builtin-functions.md#joinv), and
[joinkv](reference-dsl-builtin-functions.md#joinkv). Here the flavors have to do with whether we put
array/map keys, values, or both into the resulting string.

## Example: shortening hostnames

Suppose you want to just keep the first two components of the hostnames:

<pre class="pre-highlight-in-pair">
<b>cat data/hosts.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
host,status
xy01.east.acme.org,up
ab02.west.acme.org,down
ac91.west.acme.org,up
</pre>

Using the [`splita`](reference-dsl-builtin-functions.md#splita) and
[`joinv`](reference-dsl-builtin-functions.md#joinv) functions, along with
[array slicing](reference-main-arrays.md#slicing), we get

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/hosts.csv put '$host = joinv(splita($host, ".")[1:2], ".")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
host,status
xy01.east,up
ab02.west,down
ac91.west,up
</pre>

## Flatten/unflatten: representing arrays in CSV

In the above examples, when we split a string field into an array, we used JSON output. That's
because JSON permits nested data structures. For CSV output, Miller uses, by default, a
_flatten/unflatten strategy_: array-valued fields are turned into multiple CSV columns. For example:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split1.csv put '$codes = splitax($codes, ",")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "Alice",
  "nicknames": "Allie,Skater",
  "codes": ["1", "3", "5"]
},
{
  "name": "Robert",
  "nicknames": "Bob,Bobby,Biker",
  "codes": ["2", "4", "6"]
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/split1.csv put '$codes = splitax($codes, ",")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name,nicknames,codes.1,codes.2,codes.3
Alice,"Allie,Skater",1,3,5
Robert,"Bob,Bobby,Biker",2,4,6
</pre>

See the [flatten/unflatten: converting between JSON and tabular formats¶](flatten-unflatten.md)
for more on this default behavior, including how to override it when you prefer.

## Splitting and joining with different separators

The above is well and good when a string contains pieces with multiple instances of the same
separator.  However sometimes we have input like `5-18:53:20`. Here we can use the more flexible
[unformat](reference-dsl-builtin-functions.md#unformat) and
[format](reference-dsl-builtin-functions.md#format) DSL functions.  (As above, there's an
[unformatx](reference-dsl-builtin-functions.md#unformatx) function if you want Miller to just split
the string into string pieces without trying to type-infer them.)

<pre class="pre-highlight-in-pair">
<b>cat data/split2.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
stamp,event
5-18:53:20,open
5-18:53:22,close
5-19:07:34,open
5-19:07:56,close
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/split2.csv put '$pieces = unformat("{}-{}:{}:{}", $stamp)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "stamp": "5-18:53:20",
  "event": "open",
  "pieces": [5, 18, 53, 20]
},
{
  "stamp": "5-18:53:22",
  "event": "close",
  "pieces": [5, 18, 53, 22]
},
{
  "stamp": "5-19:07:34",
  "event": "open",
  "pieces": [5, 19, "07", 34]
},
{
  "stamp": "5-19:07:56",
  "event": "close",
  "pieces": [5, 19, "07", 56]
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from data/split2.csv put '</b>
<b>  pieces = unformat("{}-{}:{}:{}", $stamp);</b>
<b>  $description = format("{} day(s) {} hour(s) {} minute(s) {} seconds(s)", pieces[1], pieces[2], pieces[3], pieces[4]);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
stamp      event description
5-18:53:20 open  5 day(s) 18 hour(s) 53 minute(s) 20 seconds(s)
5-18:53:22 close 5 day(s) 18 hour(s) 53 minute(s) 22 seconds(s)
5-19:07:34 open  5 day(s) 19 hour(s) 07 minute(s) 34 seconds(s)
5-19:07:56 close 5 day(s) 19 hour(s) 07 minute(s) 56 seconds(s)
</pre>

## Using regular expressions and capture groups

If you prefer [regular expressions](reference-main-regular-expressions.md), those can be used in this context as well:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from data/split2.csv put '</b>
<b>  if ($stamp =~ "(\d+)-(\d+):(\d+):(\d+)") {</b>
<b>    $description = "\1 day(s) \2 hour(s) \3 minute(s) \4 seconds(s)";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
stamp      event description
5-18:53:20 open  5 day(s) 18 hour(s) 53 minute(s) 20 seconds(s)
5-18:53:22 close 5 day(s) 18 hour(s) 53 minute(s) 22 seconds(s)
5-19:07:34 open  5 day(s) 19 hour(s) 07 minute(s) 34 seconds(s)
5-19:07:56 close 5 day(s) 19 hour(s) 07 minute(s) 56 seconds(s)
</pre>

## Special case: timestamps

Timestamps are complex enough to merit their own handling: see the
[DSL datetime/timezone functions page](reference-dsl-time.md). in particular the
[strptime](reference-dsl-builtin-functions.md#strptime)
and
[strftime](reference-dsl-builtin-functions.md#strftime)
functions.

## Special case: dhms and seconds

For historical reasons, Miller has a way to represent seconds in a more human-readable format, using days,
hours, minutes, and seconds. For example:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/sec2dhms.csv put '$dhms = sec2dhms($sec)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sec     dhms
1       1s
100     1m40s
10000   2h46m40s
1000000 11d13h46m40s
</pre>

Please see
[sec2dhms](reference-dsl-builtin-functions.md#sec2dhms)
and
[dhms2sec](reference-dsl-builtin-functions.md#sec2dhms)

## Special case: financial values

One way to handle currencies is to sub out the currency marker (like `$`) as well as commas:

<pre class="pre-highlight-in-pair">
<b>echo 'd=$1234.56' | mlr put '$d = float(gsub(ssub($d, "$", ""), ",", ""))'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
d=1234.56
</pre>

## Nesting and unnesting fields

Sometimes we want not to split strings into arrays, but rather, to use them to create multiple records.

For example:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p cat data/split1.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name   nicknames       codes
Alice  Allie,Skater    1,3,5
Robert Bob,Bobby,Biker 2,4,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p nest --evar , -f nicknames data/split1.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name   nicknames codes
Alice  Allie     1,3,5
Alice  Skater    1,3,5
Robert Bob       2,4,6
Robert Bobby     2,4,6
Robert Biker     2,4,6
</pre>

See [documentation on the nest verb](reference-verbs.md#nest) for general information on how to do this.
