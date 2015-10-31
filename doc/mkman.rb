#!/usr/bin/env ruby

# ----------------------------------------------------------------
def main

  # xxx emit autogen stuff: $0, hostname, date ...
  print make_top

  print make_section('NAME', [
    "Miller is like sed, awk, cut, join, and sort for name-indexed data such as CSV."
  ])

  print make_section('DESCRIPTION', [
"""This is something the Unix toolkit always could have done, and arguably always
should have done.  It operates on key-value-pair data while the familiar
Unix tools operate on integer-indexed fields: if the natural data structure for
the latter is the array, then Miller's natural data structure is the
insertion-ordered hash map.  This encompasses a variety of data formats,
including but not limited to the familiar CSV.  (Miller can handle
positionally-indexed data as a special case.)"""
  ])

  print make_section('EXAMPLES', [
    ""
  ])

  print make_code_block(
"""
% mlr --csv cut -f hostname,uptime mydata.csv
% mlr --csv filter '$status != \"down\" && $upsec >= 10000' *.csv
% mlr --nidx put '$sum = $7 + 2.1*$8' *.dat
% grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group
% mlr join -j account_id -f accounts.dat then group-by account_name balances.dat
% mlr put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\1:\2\")' data/*
% mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*
% mlr stats2 -a linreg-pca -f u,v -g shape data/*
"""
  )

  print make_section('AUTHOR', [
    "Miller is written by John Kerl <kerl.john.r@gmail.com>.",
    "This manual page has been composed from Miller's help output by Eric MSP Veith <eveith@veith-m.de>."
  ])
  print make_section('SEE ALSO', [
    "sed(1), awk(1), cut(1), join(1), sort(1), RFC 4180: Common Format and MIME Type for " +
    "Comma-Separated Values (CSV) Files, the miller website http://johnkerl.org/miller/doc"
  ])
end

# ----------------------------------------------------------------
def make_top()
  # xxx format the data
  ".TH \"MILLER\" \"1\" \"09/14/2015\" \"\\ \\&\" \"\\ \\&\"\n"
end

# ----------------------------------------------------------------
def make_section(title, paragraphs)
  retval = ".SH \"#{title}\"\n"
  paragraphs.each do |paragraph|
    retval += ".sp\n"
    retval += groff_encode(paragraph) + "\n"
  end
  retval
end

# ----------------------------------------------------------------
def make_code_block(block)
  retval  = ".if n \\{\\\n"
  retval += ".RS 4\n"
  retval += ".\\}\n"
  retval += ".nf\n"
  retval += block
  retval += ".fi\n"
  retval += ".if n \\{\\\n"
  retval += ".RE\n"
end

# ----------------------------------------------------------------
def groff_encode(line)
  line = line.gsub(/'/, '\(cq')
  line = line.gsub(/"/, '\(dq')
  line = line.gsub(/\./, '\&')
  line = line.gsub(/-/, '\-')
  line
end

# ================================================================
main

# ================================================================

# .TH "MILLER" "1" "09/14/2015" "\ \&" "\ \&"
# .\" -----------------------------------------------------------------
# .\" * Define some portability stuff
# .\" -----------------------------------------------------------------
# .\" ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# .\" http://bugs.debian.org/507673
# .\" http://lists.gnu.org/archive/html/groff/2009-02/msg00013.html
# .\" ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
# .ie \n(.g .ds Aq \(aq
# .el       .ds Aq '
# .\" -----------------------------------------------------------------
# .\" * set default formatting
# .\" -----------------------------------------------------------------
# .\" disable hyphenation
# .nh
# .\" disable justification (adjust text to left margin only)
# .ad l

# .\" -----------------------------------------------------------------
# .\" * MAIN CONTENT STARTS HERE *
# .\" -----------------------------------------------------------------
# .SH "NAME"
# miller \- sed, awk, cut, join, and sort for name\-indexed data such as CSV

# .\" ----------------------------------------------------------------
# .SH "SYNOPSIS"
# .sp

# .\" ----------------------------------------------------------------
# .SH "DESCRIPTION"
# .sp
# Description text with backslash ampersand period\&. Also \-\- backslash dash.
# Also \(cq for single quote.
# Also \(dq for double quote.
# Also \*(Aq for single quote.

# .\" ----------------------------------------------------------------
# .SH "EXAMPLES"
# .sp
# .if n \{\
# .RS 4
# .\}
# .nf
# Here is a code block
# Here is a code block
# Here is a code block
# .fi
# .if n \{\
# .RE
# .\}

# .\" ----------------------------------------------------------------
# .SH "OPTIONS"
# .sp
# In the following option flags, the version with "i" designates the input stream, "o" the output stream, and the version without prefix sets the option for both input and output stream\&. For example: \-\-irs sets the input record separator, \-\-ors the output record separator, and \-\-rs sets both the input and output separator to the given value\&.
# .SS "SEPARATOR"
# .PP
# \-\-rs, \-\-irs, \-\-ors
# .RS 4
# Record separators, defaulting to newline
# .RE
# .PP
# \-\-fs, \-\-ifs, \-\-ofs, \-\-repifs
# .RS 4
# Field separators, defaulting to ","
# .RE
# .PP
# \-\-ps, \-\-ips, \-\-ops
# .RS 4
# Pair separators, defaulting to "="
# .RE
# .SS "DATA\-FORMAT"
# .PP
# \-\-dkvp, \-\-idkvp, \-\-odkvp
# .RS 4
# Delimited key\-value pairs, e\&.g "a=1,b=2" (default)
# .RE
# .PP
# \-\-nidx, \-\-inidx, \-\-onidx
# .RS 4
# Implicitly\-integer\-indexed fields (Unix\-toolkit style)
# .RE
# .PP
# \-\-csv, \-\-icsv, \-\-ocsv
# .RS 4
# Comma\-separated value (or tab\-separated with \-\-fs tab, etc\&.)
# .RE
# .PP
# \-\-pprint, \-\-ipprint, \-\-opprint, \-\-right
# .RS 4
# Pretty\-printed tabular (produces no output until all input is in)
# .RE
# .PP
# \-\-pprint, \-\-ipprint, \-\-opprint, \-\-right
# .RS 4
# Pretty\-printed tabular (produces no output until all input is in)
# .RE
# .sp
# \-p is a keystroke\-saver for \-\-nidx \-\-fs space \-\-repifs
# .SS "NUMERICAL FORMAT"
# .PP
# .RS 4
# Sets the numerical format given a printf\-style format string\&.
# .RE
# .SS "OTHER"
# .PP
# .RS 4
# Seeds the random number generator used for put/filter
# urand()
# with a number n of the form 12345678 or 0xcafefeed\&.
# .RE
# .SS "VERBS"
# .sp
# .it 1 an-trap
# .nr an-no-space-flag 1
# .nr an-break-flag 1
# .br
# .ps +1
# \fBcut\fR
# .RS 4
# .sp
# Usage: mlr cut [options]
# .sp
# Passes through input records with specified fields included/excluded\&.
# .PP
# \-f {a,b,c}
# .RS 4
# Field names to include for cut\&.
# .RE
# .PP
# \-o
# .RS 4
# Retain fields in the order specified here in the argument list\&. Default is to retain them in the order found in the input data\&.
# .RE
# .PP
# \-x|\-\-complement
# .RS 4
# Exclude, rather that include, field names specified by \-f\&.
# .RE
# .RE
# .sp
# .it 1 an-trap
# .nr an-no-space-flag 1
# .nr an-break-flag 1
# .br
# .ps +1
# \fBfilter\fR
# .RS 4
# .sp
# prints the AST (abstract syntax tree) for the expression, which gives full transparency on the precedence and associativity rules of Miller\(cqs grammar\&. Please use a dollar sign for field names and double\-quotes for string literals\&. Miller built\-in variables are NF, NR, FNR, FILENUM, FILENAME, PI, E\&.
# .sp
# Examples:
# .sp
# .if n \{\
# .RS 4
# .\}
# .nf
# mlr filter \*(Aqlog10($count) > 4\&.0\*(Aq
# mlr filter \*(AqFNR == 2          (second record in each file)\*(Aq
# mlr filter \*(Aqurand() < 0\&.001\*(Aq  (subsampling)
# mlr filter \*(Aq$color != "blue" && $value > 4\&.2\*(Aq
# mlr filter \*(Aq($x<\&.5 && $y<\&.5) || ($x>\&.5 && $y>\&.5)\*(Aq
# .fi
# .if n \{\
# .RE
# .\}
# .sp
# Please see http://johnkerl\&.org/miller/doc/reference\&.html for more information including function list\&.
# .RE
# .sp
# .it 1 an-trap
# .nr an-no-space-flag 1
# .nr an-break-flag 1
# .br
# .ps +1
# .\" ----------------------------------------------------------------
# .RE
# .RE

# .\" ----------------------------------------------------------------
# .SH "AUTHOR"
# .sp
# miller is written by John Kerl <kerl\&.john\&.r@gmail\&.com>\&.
# .sp
# This manual page has been composed from miller\(cqs help output by Eric MSP Veith <eveith@veith\-m\&.de>\&.
# .SH "SEE ALSO"
# .sp
# sed(1), awk(1), cut(1), join(1), sort(1), RFC 4180: Common Format and MIME Type for Comma\-Separated Values (CSV) Files, the miller website http://johnkerl\&.org/miller/doc
