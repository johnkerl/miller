<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Why?

Someone asked me the other day about design, tradeoffs, thought process, why I felt it necessary to build Miller, etc. Here are some answers.

## Who is Miller for?

For background, I'm a software engineer, with a heavy devops bent and a non-trivial amount of data-engineering in my career. **Initially I wrote Miller mainly for myself:** I'm coder-friendly (being a coder); I'm Github-friendly; most of my data are well-structured or easily structurable (TSV-formatted SQL-query output, CSV files, log files, JSON data structures); I care about interoperability between all the various formats Miller supports (I've encountered them all); I do all my work on Linux or OS X.

But now there's this neat little tool **which seems to be useful for people in various disciplines**. I don't even know entirely *who*. I can click through Github starrers and read a bit about what they seem to do, but not everyone that uses Miller is even *on* Github (or stars things). I've gotten a lot of feature requests through Github -- but only from people who are Github users.  Not everyone's a coder (it seems like a lot of Miller's Github starrers are devops folks like myself, or data-science-ish people, or biology/genomics folks.) A lot of people care 100% about CSV. And so on.

So the reason for the [Miller User Survey](https://github.com/johnkerl/miller/discussions/542) is to answer questions such as: does Miller do what you need? Do you use it for all sorts of things, or just one or two nice things? Are there things you wish it did but it doesn't? Is it almost there, or just nowhere near what you want? Are there not enough features or way too many? Are the docs too complicated; do you have a hard time finding out how to do what you want? Should I think differently about what this tool even *is* in the first place? Should I think differently about who it's for?

## What was Miller created to do?

First: there are tools like `xsv` which handles CSV marvelously and `jq` which handles JSON marvelously, and so on -- but I over the years of my career in the software industry I've found myself, and others, doing a lot of ad-hoc things which really were fundamentally the same *except* for format. So the number one thing about Miller is doing common things while supporting **multiple formats**: (a) ingest a list of records where a record is a list of key-value pairs (however represented in the input files); (b) transform that stream of records; (c) emit the transformed stream -- either in the same format as input, or in a different format.

Second thing, a lot like the first: just as I didn't want to build something only for a single file format, I didn't want to build something only for one problem domain. In my work doing software engineering, devops, data engineering, etc. I saw a lot of commonalities and I wanted to **solve as many problems simultaneously as possible**.

Third: it had to be **streaming**. As time goes by and we (some of us, sometimes) have machines with tens or hundreds of GB of RAM, it's maybe less important, but I'm unhappy with tools which ingest all data, then do stuff, then emit all data. One reason is to be able to handle files bigger than available RAM. Another reason is to be able to handle input which trickles in, e.g.  you have some process emitting data now and then and you can pipe it to Miller and it will emit transformed records one at a time.

Fourth: it had to be **fast**. This precludes all sorts of very nice things written in Ruby, for example. I love Ruby as a very expressive language, and I have several very useful little utility scripts written in Ruby. But a few years ago I ported over some of my old tried-and-true C programs and the lines-of-code count was a *lot* lower -- it was great! Until I ran them on multi-GB files and realized they took 60x as long to complete.  So I couldn't write Miller in Ruby, or in languages like it. I was going to have to do something in a low-level language in order to make it performant.

Fifth thing: I wanted Miller to be **pipe-friendly and interoperate with other command-line tools**.  Since the basic paradigm is ingest records, transform records, emit records -- where the input and output formats can be the same or different, and the transform can be complex, or just pass-through -- this means you can use it to transform data, or re-format it, or both. So if you just want to do data-cleaning/prep/formatting and do all the "real" work in R, you can. If you just want a little glue script between other tools you can get that. And if you want to do non-trivial data-reduction in Miller you can.

Sixth thing: Must have **comprehensive documentation and unit-test**. Since Miller handles a lot of formats and solves a lot of problems, there's a lot to test and a lot to keep working correctly as I add features or optimize. And I wanted it to be able to explain itself -- not only through web docs like the one you're reading but also through `man mlr` and `mlr --help`, `mlr sort --help`, etc.

Seventh thing: **Must have a domain-specific language** (DSL) **but also must let you do common things without it**. All those little verbs Miller has to help you *avoid* having to write for-loops are great. I use them for keystroke-saving: `mlr stats1 -a mean,stddev,min,max -f quantity`, for example, without you having to write for-loops or define accumulator variables. But you also have to be able to break out of that and write arbitrary code when you want to: `mlr put '$distance = $rate * $time'` or anything else you can think up. In Perl/AWK/etc.  it's all DSL. In xsv et al.  it's all verbs. In Miller I like having the combination.

Eighth thing: It's an **awful lot of fun to write**. In my experience I didn't find any tools which do multi-format, streaming, efficient, multi-purpose, with DSL and non-DSL, so I wrote one. But I don't guarantee it's unique in the world. It fills a niche in the world (people use it) but it also fills a niche in my life.

## Tradeoffs

Miller is command-line-only by design. People who want a graphical user interface won't find it here.  This is in part (a) accommodating my personal preferences, and in part (b) guided by my experience/belief that the command line is very expressive. Steeper learning curve than a GUI, yes. I consider that price worth paying for the tool-niche which Miller occupies.

Another tradeoff: supporting lists of records keeps me supporting only what can be expressed in *all* of those formats. For example, `[1,2,3,4,5]` is valid but unmillerable JSON: the list elements are not records.  So Miller can't (and won't) handle arbitrary JSON -- because Miller only handles tabular data which can be expressed in a variety of formats. 

A third tradeoff is doing build-from-scratch in a low-level language. It'd be quicker to write (but slower to run) if written in a high-level language. If Miller were written in Python, it would be implemented in significantly fewer lines of code than its current Go implementation. The DSL would just be an `eval` of Python code. And it would run slower, but maybe not enough slower to be a problem for most folks. Later I found out about the [rows](https://github.com/turicas/rows) tool -- if you find Miller useful, you should check out `rows` as well.

A fourth tradeoff is in the DSL (more visibly so in 5.0.0 but already in pre-5.0.0): how much to make it dynamically typed -- so you can just say `y=x+1` with a minimum number of keystrokes -- vs. having it do a good job of telling you when you've made a typo. This is a common paradigm across *all* languages.  Some like Ruby you don't declare anything and they're quick to code little stuff in but programs of even a few thousand lines (which isn't large in the software world) become insanely unmanageable.  Then, Java at the other extreme, does scale and is very typesafe -- but you have to type in a lot of punctuation, angle brackets, datatypes, repetition, etc. just to be able to get anything done. And some in the middle like Go are typesafe but with type-inference which aim to do the best of both. In the Miller (5.0.0) DSL you get `y=x+1` by default but you can have things like `int y = x+1` etc. so the typesafety is opt-in. See also the [Type-checking page](reference-dsl-variables.md#type-checking) for more information on this.

## Related tools

Here's a comprehensive list: [https://github.com/dbohdan/structured-text-tools](https://github.com/dbohdan/structured-text-tools). It doesn't mention [rows](https://github.com/turicas/rows) so here's a plug for that as well.

## Moving forward

I originally aimed Miller at people who already know what `sed`/`awk`/`cut`/`sort`/`join` are and wanted some options. But as time goes by I realize that tools like this can be useful to folks who *don't* know what those things are; people who aren't primarily coders; people who are scientists, or data scientists. These days some journalists do data analysis.  So moving forward in terms of docs, I am working on having more cookbook, follow-by-example stuff in addition to the existing language-reference kinds of stuff.  And continuing to seek out input from people who use Miller on where to go next.
