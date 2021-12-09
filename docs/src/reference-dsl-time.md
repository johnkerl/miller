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
# DSL datetime/timezone functions

Dates/times are not a separate data type; Miller uses ints for
[seconds since the epoch](https://en.wikipedia.org/wiki/Unix_time) and strings for formatted
date/times. In this page we take a look at what some of the various options are
for processing datetimes andd timezones in your data.

See also the [section on time-related
functions](reference-dsl-builtin-functions.md#time-functions) for
information auto-generated from Miller's online-help strings.

# Epoch seconds

[Seconds since the epoch](https://en.wikipedia.org/wiki/Unix_time), or _Unix
Time_, is seconds (positive, zero, or negative) since midnight January 1 1970
UTC. This representation has several advantages, and is quite common in the
computing world.

Since this is a [number](reference-main-arithmetic.md) in Miller -- 64-bit
signed integer or double-precision floating-point -- it can represent dates
billions of years into the past or future without worry of overflow.  (There is
no [year-2038 problem](https://en.wikipedia.org/wiki/Year_2038_problem) here.)
Being numbers, epoch-seconds are easy to store in databases, communicate over
networks in binary format, etc.  Another benefit of epoch-seconds is that
they're independent of timezone or daylight-savings time.

One minus is that, being just numbers, they're not particularly human-readable
-- hence the to-string and from-string functions described below.  Another
caveat (not really a minus) is that _epoch milliseconds_, rather than epoch
seconds, are common in some contexts, particulary JavaScript. If you ever
(anywhere) see a timestamp for the year 49,000-something -- probably someone is
treating epoch-milliseconds as epoch-seconds.

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print sec2gmt(1500000000);</b>
<b>  print sec2gmt(1500000000000);</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
2017-07-14T02:40:00Z
49503-02-10T02:40:00Z
</pre>

You can get the current system time, as epoch-seconds, using the
[systime](reference-dsl-builtin-functions.md#systime) DSL function:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '$t = systime()'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   t
yellow triangle true  1  11    43.6498  9.8870 1634784588.045347
red    square   true  2  15    79.2778  0.0130 1634784588.045385
red    circle   true  3  16    13.8103  2.9010 1634784588.045386
red    square   false 4  48    77.5542  7.4670 1634784588.045393
purple triangle false 5  51    81.2290  8.5910 1634784588.045394
red    square   false 6  64    77.1991  9.5310 1634784588.045417
purple triangle false 7  65    80.1405  5.8240 1634784588.045418
yellow circle   true  8  73    63.9785  4.2370 1634784588.045419
yellow circle   true  9  87    63.5058  8.3350 1634784588.045421
purple square   false 10 91    72.3735  8.2430 1634784588.045422
</pre>

The [systimeint](reference-dsl-builtin-functions.md#systimeint) DSL functions
is nothing more than a keystroke-saver for `int(systime())`.

# UTC times with standard format

One way to make epoch-seconds human-readable, while maintaining some of their
benefits such as being independent of timezone and daylight savings, is to use
the [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.  This was the
first (and initially only) human-readable date/time format supported by Miller
going all the way back to Miller 1.0.0.

You can get these from epoch-seconds using the 
[sec2gmt](reference-dsl-builtin-functions.md#sec2gmt) DSL function.
(Note that the terms _UTC_ and _GMT_ are used interchangeably in Miller.)
We also have [sec2gmtdate](reference-dsl-builtin-functions.md#sec2gmtdate) DSL function.

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print sec2gmt(0);</b>
<b>  print sec2gmt(1234567890.123);</b>
<b>  print sec2gmt(-1234567890.123);</b>
<b>  print;</b>
<b>  print sec2gmtdate(0);</b>
<b>  print sec2gmtdate(1234567890.123);</b>
<b>  print sec2gmtdate(-1234567890.123);</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01T00:00:00Z
2009-02-13T23:31:30Z
1930-11-18T00:28:29Z

1970-01-01
2009-02-13
1930-11-18
</pre>

# Local times with standard format; specifying timezones

You can use similar formatting for dates in your preferred timezone, not just UTC/GMT.
We have the
[sec2localtime](reference-dsl-builtin-functions.md#sec2localtime),
[sec2localdate](reference-dsl-builtin-functions.md#sec2localdate), and
[localtime2sec](reference-dsl-builtin-functions.md#localtime2sec) DSL functions.

You can specify the timezone using any of the following:

* An environment variable, e.g. `export TZ=Asia/Istanbul` at your system prompt (`set TZ=Asia/Istanbul` in Windows).
* Using the `--tz` flag. This sets the `TZ` environment variable, but only internally to the `mlr` process.
* Within a DSL expression, you can assign to `ENV["TZ"]`.
* By supplying an additional argument to any of the functions with `local` in their names.

Regardless, if you specify an invalid timezone, you'll be clearly notified:

<pre class="pre-highlight-in-pair">
<b>mlr --from example.csv --tz This/Is/A/Typo cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr: unknown time zone This/Is/A/Typo
</pre>

<pre class="pre-highlight-in-pair">
<b>export TZ=Asia/Istanbul</b>
<b>mlr -n put 'end { print sec2localtime(0) }'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01 02:00:00
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --tz America/Sao_Paulo -n put 'end { print sec2localtime(0) }'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1969-12-31 21:00:00
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  ENV["TZ"] = "Asia/Istanbul";</b>
<b>  print sec2localtime(0);</b>
<b>  print sec2localdate(0);</b>
<b>  print localtime2sec("2000-01-02 03:04:05");</b>
<b>  print;</b>
<b>  ENV["TZ"] = "America/Sao_Paulo";</b>
<b>  print sec2localtime(0);</b>
<b>  print sec2localdate(0);</b>
<b>  print localtime2sec("2000-01-02 03:04:05");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01 02:00:00
1970-01-01
946775045

1969-12-31 21:00:00
1969-12-31
946789445
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print sec2localtime(0, 0, "Asia/Istanbul");</b>
<b>  print sec2localdate(0, "Asia/Istanbul");</b>
<b>  print localtime2sec("2000-01-02 03:04:05", "Asia/Istanbul");</b>
<b>  print;</b>
<b>  print sec2localtime(0, 0, "America/Sao_Paulo");</b>
<b>  print sec2localdate(0, "America/Sao_Paulo");</b>
<b>  print localtime2sec("2000-01-02 03:04:05", "America/Sao_Paulo");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01 02:00:00
1970-01-01
946775045

1969-12-31 21:00:00
1969-12-31
946789445
</pre>

Note that for local times, Miller omits the `T` and the `Z` you see in GMT times.

We also have the 
[gmt2localtime](reference-dsl-builtin-functions.md#gmt2localtime) and
[localtime2gmt](reference-dsl-builtin-functions.md#localtime2gmt) convenience functions:

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  ENV["TZ"] = "Asia/Istanbul";</b>
<b>  print gmt2localtime("1970-01-01T00:00:00Z");</b>
<b>  print localtime2gmt("1970-01-01 00:00:00");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01 02:00:00
1969-12-31T22:00:00Z
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print gmt2localtime("1970-01-01T00:00:00Z", "America/Sao_Paulo");</b>
<b>  print gmt2localtime("1970-01-01T00:00:00Z", "Asia/Istanbul");</b>
<b>  print localtime2gmt("1970-01-01 00:00:00",  "America/Sao_Paulo");</b>
<b>  print localtime2gmt("1970-01-01 00:00:00",  "Asia/Istanbul");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1969-12-31 21:00:00
1970-01-01 02:00:00
1970-01-01T03:00:00Z
1969-12-31T22:00:00Z
</pre>

# GMT and local times with custom formats

The to-string and from-string functions we've seen so far are low-keystroking:
with a little bit of typing you can convert datetimes to/from epoch seconds.
The minus, however, is flexibility. This is where the
[strftime](reference-dsl-builtin-functions.md#strftime),
[strptime](reference-dsl-builtin-functions.md#strptime) functions come into play.

Notes:

* The names `strftime` and `strptime` far predate Miller; they were chosen for familiarity. The `f` is for _format_: from epoch-seconds to human-readable string. The `p` is for _parse_: for doing the reverse.
* Even though Miller is written in Go as of Miller 6, it still preserves [C-like](https://en.wikipedia.org/wiki/C_date_and_time_functions#strftime) `strftime` and `strptime` semantics.
  * For `strftime`, this is thanks to [https://github.com/lestrrat-go/strftime](https://github.com/lestrrat-go/strftime).
  * For `stpftime`, this is thanks to [https://github.com/pbnjay/strptime](https://github.com/pbnjay/strptime).
* See [https://devhints.io/strftime](https://devhints.io/strftime) for sample format strings you can use.

Some examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  ENV["TZ"] = "Asia/Istanbul";</b>
<b>  print strftime(0, "%Y-%m-%d %H:%M:%S");</b>
<b>  print strftime(0, "%Y-%m-%d %H:%M:%S %Z");</b>
<b>  print strftime(0, "%Y-%m-%d %H:%M:%S %z");</b>
<b>  print strftime(0, "%A, %B %e, %Y");</b>
<b>  print strftime(123456789, "%I:%M %p");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01 00:00:00
1970-01-01 00:00:00 UTC
1970-01-01 00:00:00 +0000
Thursday, January  1, 1970
09:33 PM
</pre>

Unfortunately, names from `%A` and `%B` are only available in English, as an
artifact of a design choice in the Go `time` library which Miller (and its
`strftime` / `strptime` supporting packages as noted above) rely on.

We also have
[strftimelocal](reference-dsl-builtin-functions.md#strftimelocal) and
[strptimelocal](reference-dsl-builtin-functions.md#strptimelocal):

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  ENV["TZ"] = "America/Anchorage";</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %Z");</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %z");</b>
<b>  print strftime_local(0, "%A, %B %e, %Y");</b>
<b>  print strptime_local("2020-03-01 00:00:00", "%Y-%m-%d %H:%M:%S");</b>
<b>  print;</b>
<b>  ENV["TZ"] = "Asia/Hong_Kong";</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %Z");</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %z");</b>
<b>  print strftime_local(0, "%A, %B %e, %Y");</b>
<b>  print strptime_local("2020-03-01 00:00:00", "%Y-%m-%d %H:%M:%S");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1969-12-31 14:00:00 AHST
1969-12-31 14:00:00 -1000
Wednesday, December 31, 1969
1583053200

1970-01-01 08:00:00 HKT
1970-01-01 08:00:00 +0800
Thursday, January  1, 1970
1582992000
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %Z", "America/Anchorage");</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %z", "America/Anchorage");</b>
<b>  print strftime_local(0, "%A, %B %e, %Y",        "America/Anchorage");</b>
<b>  print strptime_local("2020-03-01 00:00:00", "%Y-%m-%d %H:%M:%S", "America/Anchorage");</b>
<b>  print;</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %Z", "Asia/Hong_Kong");</b>
<b>  print strftime_local(0, "%Y-%m-%d %H:%M:%S %z", "Asia/Hong_Kong");</b>
<b>  print strftime_local(0, "%A, %B %e, %Y",        "Asia/Hong_Kong");</b>
<b>  print strptime_local("2020-03-01 00:00:00", "%Y-%m-%d %H:%M:%S", "Asia/Hong_Kong");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1969-12-31 14:00:00 AHST
1969-12-31 14:00:00 -1000
Wednesday, December 31, 1969
1583053200

1970-01-01 08:00:00 HKT
1970-01-01 08:00:00 +0800
Thursday, January  1, 1970
1582992000
</pre>

# Relative times

You can get the seconds since the Miller process start using
[uptime](reference-dsl-builtin-functions.md#uptime):


<pre class="pre-highlight-in-pair">
<b>color  shape    flag  k  index quantity rate   u</b>
</pre>
<pre class="pre-non-highlight-in-pair">
yellow triangle true  1  11    43.6498  9.8870 0.0011110305786132812
red    square   true  2  15    79.2778  0.0130 0.0011241436004638672
red    circle   true  3  16    13.8103  2.9010 0.0011250972747802734
red    square   false 4  48    77.5542  7.4670 0.0011301040649414062
purple triangle false 5  51    81.2290  8.5910 0.0011301040649414062
red    square   false 6  64    77.1991  9.5310 0.002481222152709961
purple triangle false 7  65    80.1405  5.8240 0.0024831295013427734
yellow circle   true  8  73    63.9785  4.2370 0.0024831295013427734
yellow circle   true  9  87    63.5058  8.3350 0.0024852752685546875
purple square   false 10 91    72.3735  8.2430 0.002485990524291992
</pre>

Time-differences can be done in seconds, of course; you can also use the following if you like:

<pre class="pre-highlight-in-pair">
<b>mlr -F | grep hms</b>
</pre>
<pre class="pre-non-highlight-in-pair">
dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000
dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000
fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"
fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"
hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000
hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000
sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"
sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"
</pre>

# References

* List of formatting characters for `strftime` and `strptime`: [https://devhints.io/strftime](https://devhints.io/strftime)
* List of valid timezone names: [https://en.wikipedia.org/wiki/List_of_tz_database_time_zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
