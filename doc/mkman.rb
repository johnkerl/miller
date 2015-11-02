#!/usr/bin/env ruby

# xxx note about why: mindeps

# ----------------------------------------------------------------
def main

  # xxx emit autogen stuff: $0, hostname, date ...
  print make_top

  print make_section('NAME', [
    "Miller is like sed, awk, cut, join, and sort for name-indexed data such as CSV."
  ])

  print make_section('SYNOPSIS', [
"mlr [I/O options] {verb} [verb-dependent options ...] {zero or more file names}"
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
% mlr put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\\1:\\2\")' data/*
% mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*
% mlr stats2 -a linreg-pca -f u,v -g shape data/*
"""
  )

  print make_section('OPTIONS', [
"""In the following option flags, the version with \"i\" designates the input
stream, \"o\" the output stream, and the version without prefix sets the option
for both input and output stream. For example: --irs sets the input record
separator, --ors the output record separator, and --rs sets both the input and
output separator to the given value."""
  ])

  print make_code_block(`mlr -h`)

  # xxx do better than backtick -- trap $?
  verbs = `mlr --list-all-verbs-raw`
  print make_subsection('VERBS', [
    ""
  ])
  verbs = verbs.strip.split("\n")
  for verb in verbs
    print make_subsubsection(verb, [])
    print make_code_block(`mlr #{verb} -h`)
  end

  print make_section('AUTHOR', [
    "Miller is written by John Kerl <kerl.john.r@gmail.com>.",
    "This manual page has been composed from Miller's help output by Eric MSP Veith <eveith@veith-m.de>."
  ])
  print make_section('SEE ALSO', [
    "sed(1), awk(1), cut(1), join(1), sort(1), RFC 4180: Common Format and MIME Type for " +
    "Comma-Separated Values (CSV) Files, the miller website http://johnkerl.org/miller/doc"
  ])
end

# ================================================================
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
# xxx temp
def make_subsection(title, paragraphs)
  retval = ".SS \"#{title}\"\n"
  paragraphs.each do |paragraph|
    retval += ".sp\n"
    retval += groff_encode(paragraph) + "\n"
  end
  retval
end

# ----------------------------------------------------------------
# xxx temp
def make_subsubsection(title, paragraphs)
  retval  = ".sp\n";
  retval += "\\fB#{title}\\fR\n"
  paragraphs.each do |paragraph|
    retval += ".sp\n"
    retval += groff_encode(paragraph) + "\n"
  end
  retval
end

# ----------------------------------------------------------------
def make_code_block(block)
  retval  = ".if n \\{\\\n"
  retval += ".RS 0\n"
  retval += ".\\}\n"
  retval += ".nf\n"
  retval += block.gsub('\\', '\e')
  retval += ".fi\n"
  retval += ".if n \\{\\\n"
  retval += ".RE\n"
end

# ----------------------------------------------------------------
def groff_encode(line)
  #line = line.gsub(/'/, '\(cq')
  #line = line.gsub(/"/, '\(dq')
  line = line.gsub(/\./, '\&')
  #line = line.gsub(/-/, '\-')
  line = line.gsub(/\\/, '\e')
  line
end

# ================================================================
main
