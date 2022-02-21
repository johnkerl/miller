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
for processing datetimes and timezones in your data.

See also the [section on time-related
functions](reference-dsl-builtin-functions.md#time-functions) for
information auto-generated from Miller's online-help strings.

## Epoch seconds

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
seconds, are common in some contexts, particularly JavaScript. If you ever
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

The [systimeint](reference-dsl-builtin-functions.md#systimeint) DSL function
is nothing more than a keystroke-saver for `int(systime())`.

## UTC times with standard format

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

## Local times with standard format; specifying timezones

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

## Custom formats: strptime and strftime

The to-string and from-string functions we've seen so far are low-keystroking:
with a little bit of typing you can convert datetimes to/from epoch seconds.
The minus, however, is flexibility. This is where the
[strftime](reference-dsl-builtin-functions.md#strftime) and
[strptime](reference-dsl-builtin-functions.md#strptime) functions come into play.

Notes:

* The names `strftime` and `strptime` far predate Miller; they were chosen for familiarity. The `f` is for _format_: from epoch-seconds to human-readable string. The `p` is for _parse_: for doing the reverse.
* Even though Miller is written in Go as of Miller 6, it still largely preserves [C-like](https://en.wikipedia.org/wiki/C_date_and_time_functions#strftime) `strftime` and `strptime` semantics. As noted below, not all format strings used by the C library are recognized.
  * For `strftime`, this is thanks to [https://github.com/lestrrat-go/strftime](https://github.com/lestrrat-go/strftime), with a Miller-specific modification for fractional seconds.
  * For `strftime`, this is thanks to [https://github.com/pbnjay/strptime](https://github.com/pbnjay/strptime), with Miller-specific modifications.

Available format strings for `strftime`, taken directly from [https://github.com/lestrrat-go/strftime](https://github.com/lestrrat-go/strftime):

| Pattern | Description |
|---------|-------------|
| `%A` | national representation of the full weekday name |
| `%a` | national representation of the abbreviated weekday |
| `%B` | national representation of the full month name |
| `%b` | national representation of the abbreviated month name |
| `%C` | (year / 100) as decimal number; single digits are preceded by a zero |
| `%c` | national representation of time and date |
| `%D` | equivalent to `%m/%d/%y` |
| `%d` | day of the month as a decimal number (01-31) |
| `%e` | the day of the month as a decimal number (1-31); single digits are preceded by a blank |
| `%F` | equivalent to `%Y-%m-%d` |
| `%H` | the hour (24-hour clock) as a decimal number (00-23) |
| `%h` | same as `%b` |
| `%I` | the hour (12-hour clock) as a decimal number (01-12) |
| `%j` | the day of the year as a decimal number (001-366) |
| `%k` | the hour (24-hour clock) as a decimal number (0-23); single digits are preceded by a blank |
| `%l` | the hour (12-hour clock) as a decimal number (1-12); single digits are preceded by a blank |
| `%M` | the minute as a decimal number (00-59) |
| `%m` | the month as a decimal number (01-12) |
| `%n` | a newline |
| `%p` | national representation of either "ante meridiem" (a.m.) or "post meridiem" (p.m.) as appropriate. |
| `%R` | equivalent to `%H:%M` |
| `%r` | equivalent to `%I:%M:%S %p` |
| `%S` | the second as a decimal number (00-60) |
| `%1S`, ..., `%9S` | the second as a decimal number (00-60) with 1..9 decimal places, respectively |
| `%T` | equivalent to `%H:%M:%S` |
| `%t` | a tab |
| `%U` | the week number of the year (Sunday as the first day of the week) as a decimal number (00-53) |
| `%u` | the weekday (Monday as the first day of the week) as a decimal number (1-7) |
| `%V` | the week number of the year (Monday as the first day of the week) as a decimal number (01-53) |
| `%v` | equivalent to `%e-%b-%Y` |
| `%W` | the week number of the year (Monday as the first day of the week) as a decimal number (00-53) |
| `%w` | the weekday (Sunday as the first day of the week) as a decimal number (0-6) |
| `%X` | national representation of the time |
| `%x` | national representation of the date |
| `%Y` | the year with century as a decimal number |
| `%y` | the year without century as a decimal number (00-99) |
| `%Z` | the time zone name |
| `%z` | the time zone offset from UTC |
| `%%` | a `%` |

Available format strings for `strptime`:

| Pattern | Description |
|---------|-------------|
| `%%` |  A literal '%' character. |
| `%b` |  Month as locale’s abbreviated name. |
| `%B` |  Month as locale’s full name. |
| `%d` |  Day of the month as a zero-padded decimal number. |
| `%f` |  Microsecond as a decimal number, zero-padded on the left. |
| `%H` |  Hour (24-hour clock) as a zero-padded decimal number. |
| `%I` |  Hour (12-hour clock) as a zero-padded decimal number. |
| `%j` |  Three-digit day of year, like 004 or 363. |
| `%m` |  Month as a zero-padded decimal number. |
| `%M` |  Minute as a zero-padded decimal number. |
| `%p` |  Locale’s equivalent of either AM or PM. |
| `%S` |  Second as a zero-padded decimal number. |
| `%y` |  Year without century as a zero-padded decimal number. |
| `%Y` |  Year with century as a decimal number. |
| `%z` |  UTC offset in the form +HHMM or -HHMM. |
| `%Z` |  Time zone name. UTC, EST, CST -- only if you're in that timezone. |

Examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print strftime(0, "%Y-%m-%dT%H:%M:%SZ");</b>
<b>  print strftime(0, "%FT%TZ");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-01T00:00:00Z
1970-01-01T00:00:00Z
</pre>

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

Unfortunately, names from `%A` and `%B` are only available in English, as an artifact of a design
choice in the Go `time` library which Miller (and its `strftime` / `strptime` supporting packages as
noted above) rely on.

## A note on timezones

A note on timezones for `strptime`:

* Three-letter timezone names such as `CST` are recognized _only if you're in them_. (`UTC` is an exception.) This is because these aren't globally unique: `CST` can stand for `Central Standard Time`, `_Cuba Standard Time`, `_China Standard Time`, etc.
* Timezone specifiers which _are_ globally unique are of the form `-0400` and `+0500`.
* Specifiers like `-04:30`, `UTC-8`, and `Asia/Istanbul` were not supported in Miller 5 (which used the C `strptime` library), and are likewise not supported in Miller 6. See however the `TZ` environment-variable examples below.
* If you wish to match a final `Z` in the input, use a final `Z` in the format string. For example (see [ISO8601](https://en.wikipedia.org/wiki/ISO_8601)) you can match the timestamp `1970-01-01T00:00:00Z` using the format string `%FT%TZ`.

## Fractional seconds

For historical reasons, Miller's `strftime` and `strptime` use different format specifications for fractional seconds. Examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end {</b>
<b>  print strftime(123456.789, "%Y-%m-%d %H:%M:%S");</b>
<b>  print strftime(123456.789, "%Y-%m-%d %H:%M:%1S");</b>
<b>  print strftime(123456.789, "%Y-%m-%d %H:%M:%3S");</b>
<b>  print strftime(123456.789, "%Y-%m-%d %H:%M:%6S");</b>
<b>  print strptime("1970-01-02 10:17:36.789000", "%Y-%m-%d %H:%M:%S");</b>
<b>  print strptime("1970-01-02 10:17:36.789000", "%Y-%m-%d %H:%M:%S.%f");</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1970-01-02 10:17:36
1970-01-02 10:17:36.7
1970-01-02 10:17:36.789
1970-01-02 10:17:36.789000
(error)
123456.789
</pre>

## strptime_local and strftime_local

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

## Relative times

You can get the seconds since the Miller process start using
[uptime](reference-dsl-builtin-functions.md#uptime):

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '$u=uptime()'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   u
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

## References

* Non-Miller-specific list of formatting characters for `strftime` and `strptime`: [https://devhints.io/strftime](https://devhints.io/strftime)
* List of valid timezone names: [https://en.wikipedia.org/wiki/List_of_tz_database_time_zones](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones)
