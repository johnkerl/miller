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
# Output colorization

As of version 6.0.0, Miller supports output colorization. This includes colors
per se, as well as bold, underline, and/or reverse. As of version 6.1.0, the
default is bold+underline for keys -- this works well as a default regardless
of your terminal's color settings -- with further customization options as
described below:

![pix/colorization.png](pix/colorization.png)

## Things having colors

* Keys in CSV header lines, JSON keys, etc
* Values in CSV data lines, JSON scalar values, etc
* "PASS" and "FAIL" in regression-test output
* Some [online-help](online-help.md) strings

## Rules for colorization

* By default, colorize output only if writing to stdout and stdout is a TTY.

    * Example: color: `mlr --csv cat foo.csv`
    * Example: no color: `mlr --csv cat foo.csv > bar.csv`
    * Example: no color: `mlr --csv cat foo.csv | less`

* NIDX output format (plain text, Unix-style, with implicitly positionally indexed fields) is not colored. This is because in other formats, keys are one color and values are another. For NIDX, all output is values.

## Mechanisms for colorization

* Miller uses ANSI escape sequences only.
* Doesn't try to check to see whether the terminal is capable of 256-color ANSI vs 16-color ANSI. Note that if colors are in the range 0..15 then 16-color ANSI escapes are used, so this is in the user's control.

## How you can control colorization

* Suppression/unsuppression:

    * `export MLR_NO_COLOR=true` means Miller won't color even when it normally would.
    * `export MLR_ALWAYS_COLOR=true` means Miller will color even when it normally would not. For example, you might want to use this when piping `mlr` output to `less -r`.
    * Command-line flags `--no-color` or `-M`, `--always-color` or `-C`.
    * On Windows, replace `export` with `set`

* Color choices can be specified by using environment variables or command-line flags, with values 0..255:

    * `export MLR_KEY_COLOR=208`
    * `export MLR_VALUE_COLOR=33`
    * On Windows, replace `export` with `set`
    * Likewise for `MLR_PASS_COLOR`, `MLR_FAIL_COLOR`, `MLR_HELP_COLOR`, `MLR_REPL_PS1_COLOR`, and `MLR_REPL_PS2_COLOR`.
    * Command-line flags `--key-color 208`, `--value-color 33`, etc., and likewise for `--pass-color`, `--fail-color`, `--repl-ps1-color`, `--repl-ps2-color`, and `--help-color`.
    * This is particularly useful if your terminal's background color clashes with current settings.

If environment-variable settings and command-line flags are both provided, the latter take precedence.

## How to see available color codes

Use `mlr --list-color-names` to see available names, like `orchid`. Colors can also be specified
using numbers in the range 0..255, like 170: please see `mlr --list-color-codes`.  You can also use
`bold`, `underline`, and/or `reverse`.  Additionally, combinations of all of these can be joined
with a `-`, like `red-bold`, `bold-170`, `bold-underline`, etc.

![pix/colorization2.png](pix/colorization2.png)
