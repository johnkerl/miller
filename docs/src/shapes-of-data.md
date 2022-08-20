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
# Shapes of data

## No output at all

Try `od -xcv` and/or `cat -e` on your file to check for non-printable characters.

If you're using Miller version less than 5.0.0 (try `mlr --version` on your system to find out), when the line-ending-autodetect feature was introduced, please see [http://johnkerl.org/miller-releases/miller-4.5.0/doc/index.html](http://johnkerl.org/miller-releases/miller-4.5.0/doc/index.html).

## Fields not selected

Check the field-separators of the data, e.g. with the command-line `head` program. Example: for CSV, Miller's default record separator is comma; if your data is tab-delimited, e.g. `aTABbTABc`, then Miller won't find three fields named `a`, `b`, and `c` but rather just one named `aTABbTABc`.  Solution in this case: `mlr --fs tab {remaining arguments ...}`.

Also try `od -xcv` and/or `cat -e` on your file to check for non-printable characters.

## Diagnosing delimiter specifications

Use the `file` command to see if there are CR/LF terminators (in this case, there are not):

<pre class="pre-highlight-in-pair">
<b>file data/colours.csv </b>
</pre>
<pre class="pre-non-highlight-in-pair">
data/colours.csv: Unicode text, UTF-8 text
</pre>

Look at the file to find names of fields:

<pre class="pre-highlight-in-pair">
<b>cat data/colours.csv </b>
</pre>
<pre class="pre-non-highlight-in-pair">
KEY;DE;EN;ES;FI;FR;IT;NL;PL;TO;TR
masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz
masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
</pre>

Extract a few fields:

<pre class="pre-highlight-non-pair">
<b>mlr --csv cut -f KEY,PL,TO data/colours.csv </b>
</pre>

Use XTAB output format to get a sharper picture of where records/fields are being split:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --oxtab cat data/colours.csv </b>
</pre>
<pre class="pre-non-highlight-in-pair">
KEY;DE;EN;ES;FI;FR;IT;NL;PL;TO;TR masterdata_colourcode_1;Weiß;White;Blanco;Valkoinen;Blanc;Bianco;Wit;Biały;Alb;Beyaz

KEY;DE;EN;ES;FI;FR;IT;NL;PL;TO;TR masterdata_colourcode_2;Schwarz;Black;Negro;Musta;Noir;Nero;Zwart;Czarny;Negru;Siyah
</pre>

Using XTAB output format makes it clearer that `KEY;DE;...;TR` is being treated as a single field name in the CSV header, and likewise each subsequent line is being treated as a single field value. This is because the default field separator is a comma but we have semicolons here.  Use XTAB again with different field separator (`--fs semicolon`):

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ifs semicolon --oxtab cat data/colours.csv </b>
</pre>
<pre class="pre-non-highlight-in-pair">
KEY masterdata_colourcode_1
DE  Weiß
EN  White
ES  Blanco
FI  Valkoinen
FR  Blanc
IT  Bianco
NL  Wit
PL  Biały
TO  Alb
TR  Beyaz

KEY masterdata_colourcode_2
DE  Schwarz
EN  Black
ES  Negro
FI  Musta
FR  Noir
IT  Nero
NL  Zwart
PL  Czarny
TO  Negru
TR  Siyah
</pre>

Using the new field-separator, retry the cut:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --fs semicolon cut -f KEY,PL,TO data/colours.csv </b>
</pre>
<pre class="pre-non-highlight-in-pair">
KEY;PL;TO
masterdata_colourcode_1;Biały;Alb
masterdata_colourcode_2;Czarny;Negru
</pre>

## I assigned $9 and it's not 9th

Miller records are ordered lists of key-value pairs. For NIDX format, DKVP format when keys are missing, or CSV/CSV-lite format with `--implicit-csv-header`, Miller will sequentially assign keys of the form `1`, `2`, etc. But these are not integer array indices: they're just field names taken from the initial field ordering in the input data, when it was originally read from the input file(s).

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --dkvp cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1=x,2=y,3=z
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --dkvp put '$6="a";$4="b";$55="cde"'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1=x,2=y,3=z,6=a,4=b,55=cde
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --nidx cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x,y,z
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --csv --implicit-csv-header cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1,2,3
x,y,z
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --dkvp rename 2,999</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1=x,999=y,3=z
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --dkvp rename 2,newname</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1=x,newname=y,3=z
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x,y,z | mlr --csv --implicit-csv-header reorder -f 3,1,2</b>
</pre>
<pre class="pre-non-highlight-in-pair">
3,1,2
z,x,y
</pre>

## Why doesn't mlr cut put fields in the order I want?

Example: columns `rate,shape,flag` were requested but they appear here in the order `shape,flag,rate`:

<pre class="pre-highlight-in-pair">
<b>cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
purple,triangle,false,5,51,81.2290,8.5910
red,square,false,6,64,77.1991,9.5310
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv cut -f rate,shape,flag example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape,flag,rate
triangle,true,9.8870
square,true,0.0130
circle,true,2.9010
square,false,7.4670
triangle,false,8.5910
square,false,9.5310
triangle,false,5.8240
circle,true,4.2370
circle,true,8.3350
square,false,8.2430
</pre>

The issue is that Miller's `cut`, by default, outputs cut fields in the order they appear in the input data. This design decision was made intentionally to parallel the Unix/Linux system `cut` command, which has the same semantics.

The solution is to use the `-o` option:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cut -o -f rate,shape,flag example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
rate,shape,flag
9.8870,triangle,true
0.0130,square,true
2.9010,circle,true
7.4670,square,false
8.5910,triangle,false
9.5310,square,false
5.8240,triangle,false
4.2370,circle,true
8.3350,circle,true
8.2430,square,false
</pre>

## Numbering and renumbering records

The `awk`-like built-in variable `NR` is incremented for each input record:

<pre class="pre-highlight-in-pair">
<b>cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,triangle,true,1,11,43.6498,9.8870
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
red,square,false,4,48,77.5542,7.4670
purple,triangle,false,5,51,81.2290,8.5910
red,square,false,6,64,77.1991,9.5310
purple,triangle,false,7,65,80.1405,5.8240
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv put '$nr = NR' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate,nr
yellow,triangle,true,1,11,43.6498,9.8870,1
red,square,true,2,15,79.2778,0.0130,2
red,circle,true,3,16,13.8103,2.9010,3
red,square,false,4,48,77.5542,7.4670,4
purple,triangle,false,5,51,81.2290,8.5910,5
red,square,false,6,64,77.1991,9.5310,6
purple,triangle,false,7,65,80.1405,5.8240,7
yellow,circle,true,8,73,63.9785,4.2370,8
yellow,circle,true,9,87,63.5058,8.3350,9
purple,square,false,10,91,72.3735,8.2430,10
</pre>

However, this is the record number within the original input stream -- not after any filtering you may have done:

<pre class="pre-highlight-in-pair">
<b>mlr --csv filter '$color == "yellow"' then put '$nr = NR' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate,nr
yellow,triangle,true,1,11,43.6498,9.8870,1
yellow,circle,true,8,73,63.9785,4.2370,8
yellow,circle,true,9,87,63.5058,8.3350,9
</pre>

There are two good options here. One is to use the `cat` verb with `-n`:

<pre class="pre-highlight-in-pair">
<b>mlr --csv filter '$color == "yellow"' then cat -n example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n,color,shape,flag,k,index,quantity,rate
1,yellow,triangle,true,1,11,43.6498,9.8870
2,yellow,circle,true,8,73,63.9785,4.2370
3,yellow,circle,true,9,87,63.5058,8.3350
</pre>

The other is to keep your own counter within the `put` DSL:

<pre class="pre-highlight-in-pair">
<b>mlr --csv filter '$color == "yellow"' then put 'begin {@n = 1} $n = @n; @n += 1' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate,n
yellow,triangle,true,1,11,43.6498,9.8870,1
yellow,circle,true,8,73,63.9785,4.2370,2
yellow,circle,true,9,87,63.5058,8.3350,3
</pre>

The difference is a matter of taste (although `mlr cat -n` puts the counter first).

## Options for dealing with duplicate rows

If your data has records appearing multiple times, you can use [uniq](reference-verbs.md#uniq) to show and/or count the unique records.

If you want to look at partial uniqueness -- for example, show only the first record for each unique combination of the `account_id` and `account_status` fields -- you might use `mlr head -n 1 -g account_id,account_status`. Please also see [head](reference-verbs.md#head).

## Rectangularizing data

Suppose you have a method (in whatever language) which is printing things of the form

<pre class="pre-non-highlight-non-pair">
outer=1
outer=2
outer=3
</pre>

and then calls another method which prints things of the form

<pre class="pre-non-highlight-non-pair">
middle=10
middle=11
middle=12
middle=20
middle=21
middle=30
middle=31
</pre>

and then, perhaps, that second method calls a third method which prints things of the form

<pre class="pre-non-highlight-non-pair">
inner1=100,inner2=101
inner1=120,inner2=121
inner1=200,inner2=201
inner1=210,inner2=211
inner1=300,inner2=301
inner1=312
inner1=313,inner2=314
</pre>

with the result that your program's output is

<pre class="pre-non-highlight-non-pair">
outer=1
middle=10
inner1=100,inner2=101
middle=11
middle=12
inner1=120,inner2=121
outer=2
middle=20
inner1=200,inner2=201
middle=21
inner1=210,inner2=211
outer=3
middle=30
inner1=300,inner2=301
middle=31
inner1=312
inner1=313,inner2=314
</pre>

The idea here is that middles starting with a 1 belong to the outer value of 1, and so on.  (For example, the outer values might be account IDs, the middle values might be invoice IDs, and the inner values might be invoice line-items.) If you want all the middle and inner lines to have the context of which outers they belong to, you can modify your software to pass all those through your methods. Alternatively, don't refactor your code just to handle some ad-hoc log-data formatting -- instead, use the following to [rectangularize the data](record-heterogeneity.md). The idea is to use an out-of-stream variable to accumulate fields across records. Clear that variable when you see an outer ID; accumulate fields; emit output when you see the inner IDs.

<pre class="pre-highlight-in-pair">
<b>mlr --from data/rect.txt put -q '</b>
<b>  is_present($outer) {</b>
<b>    unset @r</b>
<b>  }</b>
<b>  for (k, v in $*) {</b>
<b>    @r[k] = v</b>
<b>  }</b>
<b>  is_present($inner1) {</b>
<b>    emit @r</b>
<b>  }'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
outer=1,middle=10,inner1=100,inner2=101
outer=1,middle=12,inner1=120,inner2=121
outer=2,middle=20,inner1=200,inner2=201
outer=2,middle=21,inner1=210,inner2=211
outer=3,middle=30,inner1=300,inner2=301
outer=3,middle=31,inner1=312,inner2=301
outer=3,middle=31,inner1=313,inner2=314
</pre>

See also the [record-heterogeneity page](record-heterogeneity.md); see in
particular the [`regularize` verb](reference-verbs.md#regularize) for a way to
do this with much less keystroking.
