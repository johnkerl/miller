<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Special symbols and formatting

## How can I handle commas-as-data in various formats?

[CSV](file-formats.md) handles this well and by design:

<pre class="pre-highlight">
<b>cat commas.csv</b>
</pre>
<pre class="pre-non-highlight">
Name,Role
"Xiao, Lin",administrator
"Khavari, Darius",tester
</pre>

Likewise [JSON](file-formats.md#json):

<pre class="pre-highlight">
<b>mlr --icsv --ojson cat commas.csv</b>
</pre>
<pre class="pre-non-highlight">
{
  "Name": "Xiao, Lin",
  "Role": "administrator"
}
{
  "Name": "Khavari, Darius",
  "Role": "tester"
}
</pre>

For Miller's [XTAB](file-formats.md#xtab-vertical-tabular) there is no escaping for carriage returns, but commas work fine:

<pre class="pre-highlight">
<b>mlr --icsv --oxtab cat commas.csv</b>
</pre>
<pre class="pre-non-highlight">
Name Xiao, Lin
Role administrator

Name Khavari, Darius
Role tester
</pre>

But for [key-value-pairs](file-formats.md#dkvp-key-value-pairs) and [index-numbered](file-formats.md#nidx-index-numbered-toolkit-style) formats, commas are the default field separator. And -- as of Miller 5.4.0 anyway -- there is no CSV-style double-quote-handling like there is for CSV. So commas within the data look like delimiters:

<pre class="pre-highlight">
<b>mlr --icsv --odkvp cat commas.csv</b>
</pre>
<pre class="pre-non-highlight">
Name=Xiao, Lin,Role=administrator
Name=Khavari, Darius,Role=tester
</pre>

One solution is to use a different delimiter, such as a pipe character:

<pre class="pre-highlight">
<b>mlr --icsv --odkvp --ofs pipe cat commas.csv</b>
</pre>
<pre class="pre-non-highlight">
Name=Xiao, Lin|Role=administrator
Name=Khavari, Darius|Role=tester
</pre>

To be extra-sure to avoid data/delimiter clashes, you can also use control
characters as delimiters -- here, control-A:

<pre class="pre-highlight">
<b>mlr --icsv --odkvp --ofs '\001'  cat commas.csv | cat -v</b>
</pre>
<pre class="pre-non-highlight">
Name=Xiao, Lin\001Role=administrator
Name=Khavari, Darius\001Role=tester
</pre>

## How can I handle field names with special symbols in them?

Simply surround the field names with curly braces:

<pre class="pre-highlight">
<b>echo 'x.a=3,y:b=4,z/c=5' | mlr put '${product.all} = ${x.a} * ${y:b} * ${z/c}'</b>
</pre>
<pre class="pre-non-highlight">
x.a=3,y:b=4,z/c=5,product.all=60
</pre>

## How can I put single-quotes into strings?

This is a little tricky due to the shell's handling of quotes. For simplicity, let's first put an update script into a file:

<pre class="pre-non-highlight">
$a = "It's OK, I said, then 'for now'."
</pre>

<pre class="pre-highlight">
<b>echo a=bcd | mlr put -f data/single-quote-example.mlr</b>
</pre>
<pre class="pre-non-highlight">
a=It's OK, I said, then 'for now'.
</pre>

So, it's simple: Miller's DSL uses double quotes for strings, and you can put single quotes (or backslash-escaped double-quotes) inside strings, no problem.

Without putting the update expression in a file, it's messier:

<pre class="pre-highlight">
<b>echo a=bcd | mlr put '$a="It'\''s OK, I said, '\''for now'\''."'</b>
</pre>
<pre class="pre-non-highlight">
a=It's OK, I said, 'for now'.
</pre>

The idea is that the outermost single-quotes are to protect the `put` expression from the shell, and the double quotes within them are for Miller. To get a single quote in the middle there, you need to actually put it *outside* the single-quoting for the shell. The pieces are the following, all concatenated together:

* `$a="It`
* `\'`
* `s OK, I said,`
* `\'`
* `for now`
* `\'`
* `.`

## How to escape '?' in regexes?

One way is to use square brackets; an alternative is to use simple string-substitution rather than a regular expression.

<pre class="pre-highlight">
<b>cat data/question.dat</b>
</pre>
<pre class="pre-non-highlight">
a=is it?,b=it is!
</pre>
<pre class="pre-highlight">
<b>mlr --oxtab put '$c = gsub($a, "[?]"," ...")' data/question.dat</b>
</pre>
<pre class="pre-non-highlight">
a is it?
b it is!
c is it ...
</pre>
<pre class="pre-highlight">
<b>mlr --oxtab put '$c = ssub($a, "?"," ...")' data/question.dat</b>
</pre>
<pre class="pre-non-highlight">
a is it?
b it is!
c is it ...
</pre>

The `ssub` function exists precisely for this reason: so you don't have to escape anything.
