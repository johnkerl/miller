<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Keystroke-savers

## Short format specifiers

In our examples so far we've often made use of `mlr --icsv --opprint` or `mlr --icsv --ojson`. These are such frequently occurring patterns that they have short options like **--c2p** and **--c2j**:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p head -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag index quantity rate
yellow triangle true 11    43.6498  9.8870
red    square   true 15    79.2778  0.0130
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2j head -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "shape": "triangle",
  "flag": true,
  "index": 11,
  "quantity": 43.6498,
  "rate": 9.8870
}
{
  "color": "red",
  "shape": "square",
  "flag": true,
  "index": 15,
  "quantity": 79.2778,
  "rate": 0.0130
}
</pre>

You can get the full list here (TODO:linkify).

## File names up front

Already we saw that you can put the filename first using `--from`. When you're interacting with your data at the command line, this makes it easier to up-arrow and append to the previous command:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv sort -nr index then head -n 3</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  index quantity rate
purple square false 91    72.3735  8.2430
yellow circle true  87    63.5058  8.3350
yellow circle true  73    63.9785  4.2370
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv sort -nr index then head -n 3 then cut -f shape,quantity</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape  quantity
square 72.3735
circle 63.5058
circle 63.9785
</pre>

If there's more than one input file, you can use `--mfrom`, then however many file names, then `--` to indicate the end of your input-file-name list:

<pre class="pre-highlight-non-pair">
<b>mlr --c2p --mfrom data/*.csv -- sort -n index</b>
</pre>

## .mlrrc file

If you want the default file format for Miller to be CSV you can simply put `--csv` on a line by itself in your `~/.mlrrc` file. Then instead of `mlr --csv cat example.csv` you can just do `mlr cat example.csv`. This is just the default, though, so `mlr --opprint cat example.csv` will still use default CSV format for input, and PPRINT (tabular) for output.

You can read more about this at the [Customization](customization.md) page.
