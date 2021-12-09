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
# How original is Miller?

It isn't. Miller is one of many, many participants in the online-analytical-processing culture. Other key participants include `awk`, SQL, spreadsheets, etc. etc.  etc.  Far from being an original concept, Miller explicitly strives to imitate several existing tools:

**The Unix toolkit**: Intentional similarities as described in [Unix-toolkit Context](unix-toolkit-context.md).

Recipes abound for command-line data analysis using the Unix toolkit. Here are just a couple of my favorites:

* [http://en.wikibooks.org/wiki/Ad_Hoc_Data_Analysis_From_The_Unix_Command_Line](http://en.wikibooks.org/wiki/Ad_Hoc_Data_Analysis_From_The_Unix_Command_Line)
* [http://www.gregreda.com/2013/07/15/unix-commands-for-data-science](http://www.gregreda.com/2013/07/15/unix-commands-for-data-science)
* [https://github.com/dbohdan/structured-text-tools](https://github.com/dbohdan/structured-text-tools)

**RecordStream**: Miller owes particular inspiration to [RecordStream](https://github.com/benbernard/RecordStream). The key difference is that RecordStream is a Perl-based tool for manipulating JSON (including requiring it to separately manipulate other formats such as CSV into and out of JSON), while Miller is fast Go which handles its formats natively.  The similarities include the `sort`, `stats1` (analog of RecordStream's `collate`), and `delta` operations, as well as `filter` and `put`, and pretty-print formatting.

**stats_m**: A third source of lineage is my Python [stats_m](https://github.com/johnkerl/scripts-math/tree/master/stats) module.  This includes simple single-pass algorithms which form Miller's `stats1` and `stats2` subcommands.

**SQL**: Fourthly, Miller's `group-by` command name is from SQL, as is the term `aggregate`.

**Added value**: Miller's added values include:

* Name-indexing, compared to the Unix toolkit's positional indexing.
* Raw speed, compared to `awk`, RecordStream, `stats_m`, or various other kinds of Python/Ruby/etc. scripts one can easily create.
* Compact keystroking for many common tasks, with a decent amount of flexibility.
* Ability to handle text files on the Unix pipe, without need for creating database tables, compared to SQL databases.
* Various file formats, and on-the-fly format conversion.

**jq**: Miller does for name-indexed text what [jq](https://stedolan.github.io/jq/) does for JSON. If you're not already familiar with `jq`, please check it out!.

**What about similar tools?**

Here's a comprehensive list: [https://github.com/dbohdan/structured-text-tools](https://github.com/dbohdan/structured-text-tools).  Last I knew it doesn't mention [rows](https://github.com/turicas/rows) so here's a plug for that as well.  As it turns out, I learned about most of these after writing Miller.

**What about DOTADIW?** One of the key points of the [Unix philosophy](http://en.wikipedia.org/wiki/Unix_philosophy) is that a tool should do one thing and do it well.  Hence `sort` and `cut` do just one thing. Why does Miller put `awk`-like processing, a few SQL-like operations, and statistical reduction all into one tool?  This is a fair question. First note that many standard tools, such as `awk` and `perl`, do quite a few things -- as does `jq`.  But I could have pushed for putting format awareness and name-indexing options into `cut`, `awk`, and so on (so you could do `cut -f hostname,uptime` or `awk '{sum += $x*$y}END{print sum}'`).  Patching `cut`, `sort`, etc. on multiple operating systems is a non-starter in terms of uptake.  Moreover, it makes sense for me to have Miller be a tool which collects together format-aware record-stream processing into one place, with good reuse of Miller-internal library code for its various features.

**Why not use Perl/Python/Ruby etc.?** Maybe you should. With those tools you'll get far more expressive power, and sufficiently quick turnaround time for small-to-medium-sized data.  Using Miller you'll get something less than a complete programming language, but which is fast, with moderate amounts of flexibility and much less keystroking.

When I was first developing Miller I made a survey of several languages. Using low-level implementation languages like C, Go, Rust, and Nim, I'd need to create my own domain-specific language (DSL) which would always be less featured than a full programming language, but I'd get better performance.  Using high-level interpreted languages such as Perl/Python/Ruby I'd get the language's `eval` for free and I wouldn't need a DSL; Miller would have mainly been a set of format-specific I/O hooks. If I'd gotten good enough performance from the latter I'd have done it without question and Miller would be far more flexible.  But high-level languages win the performance criteria by a landslide so we have Miller in Go with a custom DSL.

**No, really, why one more command-line data-manipulation tool?** I wrote Miller because I was frustrated with tools like `grep`, `sed`, and so on being *line-aware* without being *format-aware*. The single most poignant example I can think of is seeing people grep data lines out of their CSV files and sadly losing their header lines.  While some lighter-than-SQL processing is very nice to have, at core I wanted the format-awareness of [RecordStream](https://github.com/benbernard/RecordStream) combined with the raw speed of the Unix toolkit. Miller does precisely that.
