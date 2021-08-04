<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Intro to Miller's programming language

In the [Miller in 10 minutes](10min.md) page we took a tour of some of Miller's most-used verbs including ``cat``, ``head``, ``tail``, ``cut``, and ``sort``. These are analogs of familiar system commands, but empowered by field-name indexing and file-format awareness: the system ``sort`` command only knows about lines and column names like ``1,2,3,4``, while ``mlr sort`` knows about CSV/TSV/JSON/etc records, and field names like ``color,shape,flag,index``.

We also caught a glimpse of Miller's ``put`` and ``filter`` verbs. These two are special since they let you express statements using Miller's programming language. It's a *embedded domain-specific language* since it's inside Miller: often referred to simply as the *Miller DSL*.

In the [DSL reference](reference-dsl.md) page we have a complete reference to Miller's programming language. For now, let's take a quick look at key features -- you can use as few or as many features as you like.

## Records and fields

Let's keep using the sample [example.csv](./example.csv). When we type

<pre>
<b>mlr --c2p put '$cost = $quantity * $rate' example.csv</b>
color  shape    flag  index quantity rate   cost
yellow triangle true  11    43.6498  9.8870 431.5655726
red    square   true  15    79.2778  0.0130 1.0306114
red    circle   true  16    13.8103  2.9010 40.063680299999994
red    square   false 48    77.5542  7.4670 579.0972113999999
purple triangle false 51    81.2290  8.5910 697.8383389999999
red    square   false 64    77.1991  9.5310 735.7846221000001
purple triangle false 65    80.1405  5.8240 466.738272
yellow circle   true  73    63.9785  4.2370 271.0769045
yellow circle   true  87    63.5058  8.3350 529.3208430000001
purple square   false 91    72.3735  8.2430 596.5747605000001
</pre>

a few things are happening:

* We refer to fields in the input data using a dollar sign and then the field name, e.g. ``$quantity``. (If a field name has special characters like a dot or slash, just use curly braces: ``${field.name}``.)
* The expression ``$cost = $quantity * $rate`` is executed once per record of the data file. Our [example.csv](./example.csv) has 10 records so this expression was executed 10 times, with the field names ``$quantity`` and ``$rate`` bound to the current record's values for those fields.
* On the left-hand side we have the new field name ``$cost`` which didn't come from the input data. Assignments to new variables result in a new field being placed after all the other ones. If we'd assigned to an existing field name, it would have been updated in-place.
* The entire expression is surrounded by single quotes, to get it past the system shell. Inside those, only double quotes have meaning in Miller's programming language.

## Multi-line statements, and statements-from-file

You can use more than one statement, separating them with semicolons, and optionally putting them on lines of their own:

<pre>
<b>mlr --c2p put '$cost = $quantity * $rate; $index = $index * 100'  example.csv</b>
color  shape    flag  index quantity rate   cost
yellow triangle true  1100  43.6498  9.8870 431.5655726
red    square   true  1500  79.2778  0.0130 1.0306114
red    circle   true  1600  13.8103  2.9010 40.063680299999994
red    square   false 4800  77.5542  7.4670 579.0972113999999
purple triangle false 5100  81.2290  8.5910 697.8383389999999
red    square   false 6400  77.1991  9.5310 735.7846221000001
purple triangle false 6500  80.1405  5.8240 466.738272
yellow circle   true  7300  63.9785  4.2370 271.0769045
yellow circle   true  8700  63.5058  8.3350 529.3208430000001
purple square   false 9100  72.3735  8.2430 596.5747605000001
</pre>

<pre>
<b>mlr --c2p put '</b>
<b>  $cost = $quantity * $rate;</b>
<b>  $index *= 100</b>
<b>' example.csv</b>
color  shape    flag  index quantity rate   cost
yellow triangle true  1100  43.6498  9.8870 431.5655726
red    square   true  1500  79.2778  0.0130 1.0306114
red    circle   true  1600  13.8103  2.9010 40.063680299999994
red    square   false 4800  77.5542  7.4670 579.0972113999999
purple triangle false 5100  81.2290  8.5910 697.8383389999999
red    square   false 6400  77.1991  9.5310 735.7846221000001
purple triangle false 6500  80.1405  5.8240 466.738272
yellow circle   true  7300  63.9785  4.2370 271.0769045
yellow circle   true  8700  63.5058  8.3350 529.3208430000001
purple square   false 9100  72.3735  8.2430 596.5747605000001
</pre>

One of Miller's key features is the ability to express data-transformation right there at the keyboard, interactively. But if you find yourself using expressions repeatedly, you can put everything between the single quotes into a file and refer to that using ``put -f``:

<pre>
<b>cat dsl-example.mlr</b>
$cost = $quantity * $rate;
$index *= 100
</pre>

<pre>
<b>mlr --c2p put -f dsl-example.mlr example.csv</b>
color  shape    flag  index quantity rate   cost
yellow triangle true  1100  43.6498  9.8870 431.5655726
red    square   true  1500  79.2778  0.0130 1.0306114
red    circle   true  1600  13.8103  2.9010 40.063680299999994
red    square   false 4800  77.5542  7.4670 579.0972113999999
purple triangle false 5100  81.2290  8.5910 697.8383389999999
red    square   false 6400  77.1991  9.5310 735.7846221000001
purple triangle false 6500  80.1405  5.8240 466.738272
yellow circle   true  7300  63.9785  4.2370 271.0769045
yellow circle   true  8700  63.5058  8.3350 529.3208430000001
purple square   false 9100  72.3735  8.2430 596.5747605000001
</pre>

This becomes particularly important on Windows. Quite a bit of effort was put into making Miller on Windows be able to handle the kinds of single-quoted expressions we're showing here, but if you get syntax-error messages on Windows using examples in this documentation, you can put the parts between single quotes into a file and refer to that using ``mlr put -f``.

## Out-of-stream variables, begin, and end

Above we saw that your expression is executed once per record -- if a file has a million records, your expression will be executed a million times, once for each record. But you can mark statements to only be executed once, either before the record stream begins, or after the record stream is ended. If you know about [AWK](https://en.wikipedia.org/wiki/AWK), you might have noticed that Miller's programming language is loosely inspired by it, including the ``begin`` and ``end`` statements.

Above we also saw that names like ``$quantity`` are bound to each record in turn.

To make ``begin`` and ``end`` statements useful, we need somewhere to put things that persist across the duration of the record stream, and a way to emit them. Miller uses **out-of-stream variables** (or **oosvars** for short) whose names start with an ``@`` sigil, and the **emit** keyword to write them into the output record stream:

<pre>
<b>mlr --c2p --from example.csv put 'begin { @sum = 0 } @sum += $quantity; end {emit @sum}'</b>
color  shape    flag  index quantity rate
yellow triangle true  11    43.6498  9.8870
red    square   true  15    79.2778  0.0130
red    circle   true  16    13.8103  2.9010
red    square   false 48    77.5542  7.4670
purple triangle false 51    81.2290  8.5910
red    square   false 64    77.1991  9.5310
purple triangle false 65    80.1405  5.8240
yellow circle   true  73    63.9785  4.2370
yellow circle   true  87    63.5058  8.3350
purple square   false 91    72.3735  8.2430

sum
652.7185
</pre>

If you want the end-block output to be the only output, and not include the input data, you can use ``mlr put -q``:

<pre>
<b>mlr --c2p --from example.csv put -q 'begin { @sum = 0 } @sum += $quantity; end {emit @sum}'</b>
sum
652.7185
</pre>

<pre>
<b>mlr --c2j --from example.csv put -q 'begin { @sum = 0 } @sum += $quantity; end {emit @sum}'</b>
{
  "sum": 652.7185
}
</pre>

<pre>
<b>mlr --c2j --from example.csv put -q '</b>
<b>  begin { @count = 0; @sum = 0 }</b>
<b>  @count += 1;</b>
<b>  @sum += $quantity;</b>
<b>  end {emit (@count, @sum)}</b>
<b>'</b>
{
  "count": 10,
  "sum": 652.7185
}
</pre>

We'll see in the documentation for :ref:`reference-verbs-stats1` that there's a lower-keystroking way to get counts and sums of things -- so, take this sum/count example as an indication of the kinds of things you can do using Miller's programming language.

## Context variables

Also inspired by [AWK](https://en.wikipedia.org/wiki/AWK), the Miller DSL has the following special **context variables**:

* ``FILENAME`` -- the filename the current record came from. Especially useful in things like ``mlr ... *.csv``.
* ``FILENUM`` -- similarly, but integer 1,2,3,... rather than filenam.e
* ``NF`` -- the number of fields in the current record. Note that if you assign ``$newcolumn = some value`` then ``NF`` will increment.
* ``NR`` -- starting from 1, counter of how many records processed so far.
* ``FNR`` -- similar, but resets to 1 at the start of each file.

<pre>
<b>cat context-example.mlr</b>
$nf       = NF;
$nr       = NR;
$fnr      = FNR;
$filename = FILENAME;
$filenum  = FILENUM;
$newnf    = NF;
</pre>

<pre>
<b>mlr --c2p put -f context-example.mlr data/a.csv data/b.csv</b>
a b c nf nr fnr filename   filenum newnf
1 2 3 3  1  1   data/a.csv 1       8
4 5 6 3  2  2   data/a.csv 1       8
7 8 9 3  3  1   data/b.csv 2       8
</pre>

## Functions and local variables

You can define your own functions:

<pre>
<b>cat factorial-example.mlr</b>
func factorial(n) {
  if (n <= 1) {
    return n
  } else {
    return n * factorial(n-1)
  }
}
</pre>

<pre>
<b>mlr --c2p --from example.csv put -f factorial-example.mlr -e '$fact = factorial(NR)'</b>
color  shape    flag  index quantity rate   fact
yellow triangle true  11    43.6498  9.8870 1
red    square   true  15    79.2778  0.0130 2
red    circle   true  16    13.8103  2.9010 6
red    square   false 48    77.5542  7.4670 24
purple triangle false 51    81.2290  8.5910 120
red    square   false 64    77.1991  9.5310 720
purple triangle false 65    80.1405  5.8240 5040
yellow circle   true  73    63.9785  4.2370 40320
yellow circle   true  87    63.5058  8.3350 362880
purple square   false 91    72.3735  8.2430 3628800
</pre>

Note that here we used the ``-f`` flag to ``put`` to load our function
definition, and also the ``-e`` flag to add another statement on the command
line. (We could have also put ``$fact = factorial(NR)`` inside
``factorial-example.mlr`` but that would have made that file less flexible for our
future use.)

## If-statements, loops, and local variables

Suppose you want to only compute sums conditionally -- you can use an ``if`` statement:

<pre>
<b>cat if-example.mlr</b>
begin {
  @count_of_red = 0;
  @sum_of_red = 0
}

if ($color == "red") {
  @count_of_red += 1;
  @sum_of_red += $quantity;
}

end {
  emit (@count_of_red, @sum_of_red)
}
</pre>

<pre>
<b>mlr --c2p --from example.csv put -q -f if-example.mlr</b>
count_of_red sum_of_red
4            247.84139999999996
</pre>

Miller's else-if is spelled ``elif``.

As we'll see more of in section (TODO:linkify), Miller has a few kinds of
for-loops. In addition to the usual 3-part ``for (i = 0; i < 10; i += 1)`` kind
that many programming languages have, Miller also lets you loop over arrays and
hashmaps. We haven't encountered arrays and hashmaps yet in this introduction,
but for now it suffices to know that ``$*`` is a special variable holding the
current record as a hashmap:

<pre>
<b>cat for-example.mlr</b>
for (k, v in $*) {
  print "KEY IS ". k . " VALUE IS ". v;
}
print
</pre>

<pre>
<b>mlr --csv cat data/a.csv</b>
a,b,c
1,2,3
4,5,6
</pre>

<pre>
<b>mlr --csv --from data/a.csv put -qf for-example.mlr</b>
KEY IS a VALUE IS 1
KEY IS b VALUE IS 2
KEY IS c VALUE IS 3

KEY IS a VALUE IS 4
KEY IS b VALUE IS 5
KEY IS c VALUE IS 6
</pre>

Here we used the local variables ``k`` and ``v``. Now we've seen four kinds of variables:

* Record fields like ``$shape``
* Out-of-stream variables like ``@sum``
* Local variables like ``k``
* Built-in context variables like ``NF`` and ``NR``

If you're curious about scope and extent of local variables, you can read more in (TODO:linkify) the section on variables.

## Arithmetic

Numbers in Miller's programming language are intended to operate with the principle of least surprise:

* Internally, numbers are either 64-bit signed integers or double-precision floating-point.
* Sums, differences, and products of integers are also integers (so ``2*3=6`` not ``6.0``) -- unless the result of the operation would overflow a 64-bit signed integer in which case the result is automatically converted to float. (If you ever want integer-to-integer arithmetic, use ``x .+ y``, ``x .* y``, etc.)
* Quotients of integers are integers if the division is exact, else floating-point:  so ``6/2=3`` but ``7/2=3.5``.

You can read more about this at (TODO:linkify).

## Absent data

In addition to types including string, number (int/float), arrays, and hashmaps, Miller varibles can also be **absent**. This is when a variable never had a value assigned to it. Miller's treatment of absent data is intended to make it easy for you to handle non-heterogeneous data. We'll see more in section (TODO:linkify) but the basic idea is:

* Adding a number to absent gives the number back. This means you don't have to put ``@sum = 0`` in your ``begin`` blocks.
* Any variable which has the absent value is not assigned. This means you don't have to check presence of things from one record to the next.

For example, you can sum up all the ``$a`` values across records without having to check whether they're present or not:

<pre>
<b>mlr --json cat absent-example.json</b>
{
  "a": 1,
  "b": 2
}
{
  "c": 3
}
{
  "a": 4,
  "b": 5
}
</pre>

<pre>
<b>mlr --json put '@sum_of_a += $a; end {emit @sum_of_a}' absent-example.json</b>
{
  "a": 1,
  "b": 2
}
{
  "c": 3
}
{
  "a": 4,
  "b": 5
}
{
  "sum_of_a": 5
}
</pre>
