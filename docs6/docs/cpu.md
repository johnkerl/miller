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
# CPU/multicore usage

Miller 6 is written in [Go](https://golang.org/) which supports multicore programming.

Miller uses Go's _channel_ concept. The following are all separate goroutines:

* One channel for input-record reader: parsing input file(s) to record objects
* One channel for each [verb](reference-verbs.md) in [then-chains](reference-main-then-chaining.md)
* One channel for output-record writer: formatting output records as text
* One controller channel which coordinates all these, without much work to do.

For example, `mlr --csv cut -f somefield then sort -f otherfield then put '$z =
$x + $y' a.csv b.csv c.csv` will have 6 goroutines running: input-reader,
`cut`, `sort`, `put`, output-writer, controller.

If all the verbs in the chain are [streaming](streaming-and-memory.md) --
operating on each record as it arrives, then passing it on -- then all verbs in
the chain will be active at once. On the other hand, if there is a
non-streaming verb in the chain, which produces output only after receiving all
input -- for example, `sort` -- then we would expect verbs after that in the
chain to sit idle until the end of the input stream is reached, the `sort` does
its computation, then sends its output to downstream verbs.

In practice, profiling has shown that the input-reader uses the most CPU of all
the above. This means CPUs running verbs may not be 100% utilized, since they
are likely to be spending some of their time waiting for input data.

Running Miller on a machine with more CPUs than active channels (as listed
above) won't speed up a given invocation of Miller. However, of course, you'll
be able to run more invocations of Miller at the same time if you like.

You can set the Go-standard environment variable `GOMAXPROCS` if you like. If
you don't, Miller will (as is standard for Go programs in Go 1.16 and above) up
to all available CPUs.

If you set `$GOMAXPROCS=1` in the environment, that's fine -- the Go runtime
will multiplex different channel-handling goroutines onto the same CPU.
