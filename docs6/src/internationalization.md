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
# Internationalization

Miller handles ASCII and UTF-8 strings. (I have no plans to support UTF-16 or ISO-8859-1.)

Support for internationalization includes:

* Tabular output formats such pprint and xtab (see [File Formats](file-formats.md)) are aligned correctly.
* The [strlen](reference-dsl-builtin-functions.md#strlen) function correctly counts UTF-8 codepoints rather than bytes.
* The [toupper](reference-dsl-builtin-functions.md#toupper), [tolower](reference-dsl-builtin-functions.md#tolower), and [capitalize](reference-dsl-builtin-functions.md#capitalize) DSL functions operate within the capabilities of the Go libraries.
* While Miller's function names, verb names, online help, etc. are all in English, you can write field names, string literals, variable names, etc in UTF-8.

<pre class="pre-highlight-in-pair">
<b>cat παράδειγμα.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
χρώμα,σχήμα,σημαία,κ,δείκτης,ποσότητα,ρυθμός
κίτρινο,τρίγωνο,αληθινό,1,11,43.6498,9.8870
κόκκινο,τετράγωνο,αληθινό,2,15,79.2778,0.0130
κόκκινο,κύκλος,αληθινό,3,16,13.8103,2.9010
κόκκινο,τετράγωνο,ψευδές,4,48,77.5542,7.4670
μοβ,τρίγωνο,ψευδές,5,51,81.2290,8.5910
κόκκινο,τετράγωνο,ψευδές,6,64,77.1991,9.5310
μοβ,τρίγωνο,ψευδές,7,65,80.1405,5.8240
κίτρινο,κύκλος,αληθινό,8,73,63.9785,4.2370
κίτρινο,κύκλος,αληθινό,9,87,63.5058,8.3350
μοβ,τετράγωνο,ψευδές,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p filter '$σχήμα == "κύκλος"' παράδειγμα.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
χρώμα   σχήμα  σημαία  κ δείκτης ποσότητα ρυθμός
κόκκινο κύκλος αληθινό 3 16      13.8103  2.9010
κίτρινο κύκλος αληθινό 8 73      63.9785  4.2370
κίτρινο κύκλος αληθινό 9 87      63.5058  8.3350
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p sort -f σημαία παράδειγμα.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
χρώμα   σχήμα     σημαία  κ  δείκτης ποσότητα ρυθμός
κίτρινο τρίγωνο   αληθινό 1  11      43.6498  9.8870
κόκκινο τετράγωνο αληθινό 2  15      79.2778  0.0130
κόκκινο κύκλος    αληθινό 3  16      13.8103  2.9010
κίτρινο κύκλος    αληθινό 8  73      63.9785  4.2370
κίτρινο κύκλος    αληθινό 9  87      63.5058  8.3350
κόκκινο τετράγωνο ψευδές  4  48      77.5542  7.4670
μοβ     τρίγωνο   ψευδές  5  51      81.2290  8.5910
κόκκινο τετράγωνο ψευδές  6  64      77.1991  9.5310
μοβ     τρίγωνο   ψευδές  7  65      80.1405  5.8240
μοβ     τετράγωνο ψευδές  10 91      72.3735  8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p put '$форма = toupper($форма); $длина = strlen($цвет)' пример.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
цвет       форма       флаг   κ  индекс количество скорость длина
желтый     ТРЕУГОЛЬНИК истина 1  11     43.6498    9.8870   6
красный    КВАДРАТ     истина 2  15     79.2778    0.0130   7
красный    КРУГ        истина 3  16     13.8103    2.9010   7
красный    КВАДРАТ     ложь   4  48     77.5542    7.4670   7
фиолетовый ТРЕУГОЛЬНИК ложь   5  51     81.2290    8.5910   10
красный    КВАДРАТ     ложь   6  64     77.1991    9.5310   7
фиолетовый ТРЕУГОЛЬНИК ложь   7  65     80.1405    5.8240   10
желтый     КРУГ        истина 8  73     63.9785    4.2370   6
желтый     КРУГ        истина 9  87     63.5058    8.3350   6
фиолетовый КВАДРАТ     ложь   10 91     72.3735    8.2430   10
</pre>
