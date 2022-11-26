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
# DSL unset statements

You can clear a map key by assigning the empty string as its value: `$x=""` or `@x=""`. Using `unset` you can remove the key entirely. Examples:

<pre class="pre-highlight-in-pair">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put 'unset $x, $a' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
b=pan,i=1,y=0.726802
b=pan,i=2,y=0.522151
b=wye,i=3,y=0.338318
b=wye,i=4,y=0.134188
b=pan,i=5,y=0.863624
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

This can also be done, of course, using `mlr cut -x`. You can also clear out-of-stream or local variables, at the base name level, or at an indexed sublevel:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum; dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.346791
    },
    "eks": {
      "pan": 0.758679,
      "wye": 0.381399
    },
    "wye": {
      "wye": 0.204603,
      "pan": 0.573288
    }
  }
}
{}
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum["eks"]; dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.346791
    },
    "eks": {
      "pan": 0.758679,
      "wye": 0.381399
    },
    "wye": {
      "wye": 0.204603,
      "pan": 0.573288
    }
  }
}
{
  "sum": {
    "pan": {
      "pan": 0.346791
    },
    "wye": {
      "wye": 0.204603,
      "pan": 0.573288
    }
  }
}
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

If you use `unset all` (or `unset @*` which is synonymous), that will unset all out-of-stream variables which have been assigned up to that point.
