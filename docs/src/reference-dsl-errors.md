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
# DSL errors and transparency

# Handling for data errors

By default, Miller doesn't stop data processing for a single cell error. For example:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data-error.csv cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x
1
2
3
text
4
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data-error.csv put '$y = log10($x)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x,y
1,0
2,0.3010299956639812
3,0.4771212547196624
text,(error)
4,0.6020599913279624
</pre>

If you do want to stop processing, though, you have three options. The first is the `mlr -x` flag:

<pre class="pre-highlight-in-pair">
<b>mlr -x --csv --from data-error.csv put '$y = log10($x)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x,y
1,0
2,0.3010299956639812
3,0.4771212547196624
mlr: data error at NR=4 FNR=4 FILENAME=data-error.csv
mlr: field y: log10: unacceptable type string with value "text"
mlr: exiting due to data error.
</pre>

The second is to put `-x` into your [`~/.mlrrc file`](customization.md).

The third is to set the `MLR_FAIL_ON_DATA_ERROR` environment variable, which makes `-x` implicit.

# Common causes of syntax errors

As soon as you have a [programming language](miller-programming-language.md), you start having the problem *What is my code doing, and why?* This includes getting syntax errors -- which are always annoying -- as well as the even more annoying problem of a program which parses without syntax error but doesn't do what you expect.

The syntax-error message gives you line/column position for the syntax that couldn't be parsed. The cause may be clear from that information, or perhaps not.  Here are some common causes of syntax errors:

* Don't forget `;` at end of line, before another statement on the next line.

* Miller's DSL lacks the `++` and `--` operators.

* Curly braces are required for the bodies of `if`/`while`/`for` blocks, even when the body is a single statement.

# Transparency

As for transparency:

* As in any language, you can do `print`, or `eprint` to print to stderr.  See [Print statements](reference-dsl-output-statements.md#print-statements); see also [Dump statements](reference-dsl-output-statements.md#dump-statements) and [Emit statements](reference-dsl-output-statements.md#emit-statements).

* The `-v` option to `mlr put` and `mlr filter` prints abstract syntax trees for your code. While not all details here will be of interest to everyone, certainly this makes questions such as operator precedence completely unambiguous.

* Please see [type-checking](reference-dsl-variables.md#type-checking) for type declarations and type-assertions you can use to make sure expressions and the data flowing them are evaluating as you expect.  I made them optional because one of Miller's important use-cases is being able to say simple things like `mlr put '$y = $x + 1' myfile.dat` with a minimum of punctuational bric-a-brac -- but for programs over a few lines long, I generally find that the more type-specification, the better.
