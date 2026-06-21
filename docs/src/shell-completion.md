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
# Shell completion

Miller can generate tab-completion scripts for `bash` and `zsh`. Once installed,
pressing <b>TAB</b> completes Miller's main flags, verb names, subcommands like
`help` and `version`, each verb's own flags, the `then` keyword, and filenames
-- and it does so in a way that understands Miller's
[then-chains](reference-main-then-chaining.md).

## Why this is more than the usual flag completion

Most command-line tools have a single set of flags: `prog --flag1 --flag2 file`.

Miller's command line, by contrast, is instead a [sequence of contexts](#reference-main-overview.md):

<pre class="pre-non-highlight-non-pair">
mlr {main flags} verb1 {verb1 flags} then verb2 {verb2 flags} {filenames}
</pre>

So the same word can mean different things depending on where it sits. Miller's
completion walks the command line left-to-right and offers candidates
appropriate to the cursor's position:

* Before the first verb: main flags (e.g. `--icsv`), verb names (e.g. `cat`),
  and subcommands (e.g. `help`, `version`, `repl`).
* Inside a verb: that verb's own flags, plus `then` and filenames.
* Right after `then`: verb names.
* As the argument to a flag that takes one (e.g. `mlr --ifs`): the flag's
  values where these are a known set, otherwise filenames.

## Installing for bash

Add this to your `~/.bashrc`:

<pre class="pre-non-highlight-non-pair">
eval "$(mlr completion bash)"
</pre>

Or install it system-wide (loaded by the `bash-completion` package):

<pre class="pre-non-highlight-non-pair">
mlr completion bash > /etc/bash_completion.d/mlr
</pre>

Prefer `eval "$(mlr completion bash)"` over `source <(mlr completion bash)`. The
latter silently does nothing on the bash 3.2 that ships with macOS, where
sourcing from a process-substitution file descriptor can read no data.

## Installing for zsh

Add this to your `~/.zshrc`:

<pre class="pre-non-highlight-non-pair">
eval "$(mlr completion zsh)"
</pre>

Or place the script on your `$fpath` so zsh autoloads it:

<pre class="pre-non-highlight-non-pair">
mlr completion zsh > "${fpath[1]}/_mlr"
</pre>

The generated script initializes zsh's completion system (`compinit`) if your
startup files have not already done so, so it works even with a minimal
`~/.zshrc`.

## What completion looks like

Before the first verb, <b>TAB</b> offers verb names along with subcommands like
`help` and `version`:

<pre class="pre-non-highlight-non-pair">
mlr <b>TAB</b>
altkv      cat        completion     count      ...      help      version
</pre>

The top-level help and version flags are offered too:

<pre class="pre-non-highlight-non-pair">
mlr --v<b>TAB</b>
--value-color   --version   --vflatsep
</pre>

The `help` subcommand completes its topics, and topics that take a name
argument complete that too:

<pre class="pre-non-highlight-non-pair">
mlr help <b>TAB</b>
flags   list-verbs   verb   function   keyword   ...

mlr help verb <b>TAB</b>
altkv   bar   cat   count   cut   ...

mlr help function strl<b>TAB</b>
strlen
</pre>

A leading dash offers main flags, including the format-conversion
keystroke-savers (`--c2j`, `--x2y`, and so on):

<pre class="pre-non-highlight-non-pair">
mlr --c<b>TAB</b>
--c2b   --c2c   --c2d   --c2j   --c2p   ...   --csv   --csvlite
</pre>

Inside a verb, <b>TAB</b> offers that verb's flags:

<pre class="pre-non-highlight-non-pair">
mlr cat -<b>TAB</b>
-n   -N   -g   --filename   --filenum
</pre>

After `then`, <b>TAB</b> offers verb names again:

<pre class="pre-non-highlight-non-pair">
mlr --icsv cat -n then head -n 10 then <b>TAB</b>
altkv      cat        count          cut        ...
</pre>

For flags whose argument is a known set of values, <b>TAB</b> offers those
values. Format flags (`-i`, `-o`, `--io`) offer file-format names:

<pre class="pre-non-highlight-non-pair">
mlr -i <b>TAB</b>
csv   csvlite   dcf   dkvp   dkvpx   gen   json   markdown   nidx   pprint   tsv   xtab   yaml
</pre>

and separator flags (`--ifs`, `--ofs`, `--ips`, and so on) offer the named
separator aliases:

<pre class="pre-non-highlight-non-pair">
mlr --ifs <b>TAB</b>
comma   pipe   semicolon   space   tab   ...
</pre>

## Generating the scripts

The `mlr completion` command prints the scripts, and `mlr completion --help`
describes the options:

<pre class="pre-highlight-in-pair">
<b>mlr completion --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr completion {bash|zsh}
Generates a shell tab-completion script for Miller.

Bash:
  Add to your ~/.bashrc:
    eval "$(mlr completion bash)"
  Or install system-wide:
    mlr completion bash > /etc/bash_completion.d/mlr
  Note: prefer 'eval' over 'source <(mlr completion bash)'. The latter
  silently fails on the bash 3.2 that ships with macOS, where sourcing from a
  process-substitution FIFO can read nothing.

Zsh:
  Add to your ~/.zshrc:
    eval "$(mlr completion zsh)"
  Or place the output on your $fpath, e.g.:
    mlr completion zsh > "${fpath[1]}/_mlr"
  The script initializes zsh's completion system (compinit) if your startup
  files have not done so already.

Completion is context-aware across Miller's then-chains: it offers main flags
and verb names before the first verb, the current verb's flags inside a verb,
verb names after 'then', and filenames where appropriate.
</pre>

## Notes

* Completion candidates are produced by Miller itself: the shell scripts simply
  forward the current command-line words to `mlr` and render what it returns.
  This means completion stays in sync with Miller's flags and verbs
  automatically -- there is no separate list to maintain.

* Value completion is offered for flags with a known set of values (file
  formats for `-i`/`-o`/`--io`, separator aliases for `--ifs` and friends).
  Other arg-taking flags fall back to filename completion; per-verb argument
  values (such as field names) are not yet completed.

* The `mlr completion complete ...` subcommand is an internal interface used by
  the generated scripts; it is not intended to be run directly.
