<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Flatten/unflatten: JSON vs. tabular formats

TODO

* JSON-to-JSON
* JSON-to-not
* not-to-JSON
* not-to-not
* 'concatening keys'
* to-array heuristic
* no-flatten / no-unflatten options

## TBF

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
