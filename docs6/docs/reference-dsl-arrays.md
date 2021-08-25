<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Arrays

TODO

<pre class="pre-highlight-in-pair">
<b>mlr --json cat data/array-example.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "key": "ax04",
  "samples": [45, 67, 22]
}
{
  "key": "cz09",
  "samples": [11, 29, 84, 91]
}
</pre>

## TBF

[Arrays](reference-dsl-arrays.md) are supported [as of Miller 6](new-in-miller-6.md).

Suppose we have arrays like this in our input data:

<pre class="pre-highlight-in-pair">
<b>cat data/json-example-3.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "label": "orange",
  "values": [12.2, 13.8, 17.2]
}
{
  "label": "purple",
  "values": [27.0, 32.4]
}
</pre>

comment:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --oxtab cat data/json-example-3.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
label    orange
values.1 12.2
values.2 13.8
values.3 17.2

label    purple
values.1 27.0
values.2 32.4
</pre>

comment:

<pre class="pre-highlight-in-pair">
<b>mlr --json --jvstack cat data/json-example-3.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "label": "orange",
  "values": [12.2, 13.8, 17.2]
}
{
  "label": "purple",
  "values": [27.0, 32.4]
}
</pre>

* 1-up and why
* 1-var/2-var for-loops
* auto-extend and null-gaps
* x[1]=2 is map not array if x doesn't exist -- xlink to maps page
* POLS mentions
