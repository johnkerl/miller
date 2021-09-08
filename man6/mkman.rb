#!/usr/bin/env ruby

# ================================================================
# This is a manpage autogenerator for Miller. There are various tools out there
# for creating xroff-formatted manpages, but I wanted something with minimal
# external dependencies which would also automatically generate most of its
# output from the mlr executable itself.  It turns out it's easy enough to get
# this in just a few lines of Ruby.
#
# Note for dev-viewing of the output:
# ./mkman.rb | groff -man -Tascii | less
# ================================================================

# ----------------------------------------------------------------
def main
  # In case the user running this has a .mlrrc
  ENV['MLRRC'] = '__none__'

  print make_top

  print make_section('NAME', [
    "miller \\\- like awk, sed, cut, join, and sort for name-indexed data such as CSV and tabular JSON."
  ])

  print make_section('SYNOPSIS', [
`mlr --help`
  ])

  print make_section('DESCRIPTION', [
"""Miller operates on key-value-pair data while the familiar Unix tools operate
on integer-indexed fields: if the natural data structure for the latter is the
array, then Miller's natural data structure is the insertion-ordered hash map.
This encompasses a variety of data formats, including but not limited to the
familiar CSV, TSV, and JSON.  (Miller can handle positionally-indexed data as
a special case.) This manpage documents #{`mlr --version`.chomp}."""
  ])

  print make_section('EXAMPLES', [
    ""
  ])
	print make_code_block(`mlr help basic-examples`)

	print make_section('DATA FORMATS', [])
	print make_code_block(`mlr help data-formats`)

	print make_section('HELP OPTIONS', [])
	print make_code_block(`mlr help topics`)

	print make_section('VERB LIST', [])
	print make_code_block(`mlr help list-verbs-as-paragraph`)

	print make_section('FUNCTION LIST', [])
	print make_code_block(`mlr help list-functions-as-paragraph`)

  section_names = `mlr help list-flag-sections`.split("\n")
  for section_name in section_names
	  print make_section(section_name.upcase, [""])
	  print make_code_block(`mlr help show-help-for-section '#{section_name}'`)
  end

	print make_section('AUXILIARY COMMANDS', [])
	print make_code_block(`mlr aux-list`)

  print make_section('MLRRC', [])
  print make_code_block(`mlr help mlrrc`)

  print make_section('REPL', [])
  print make_code_block(`mlr repl -h`)

  verbs = `mlr help list-verbs`
  print make_section('VERBS', [
    ""
  ])
  verbs = verbs.strip.split("\n")
  for verb in verbs
    print make_subsection(verb, [])
    print make_code_block(`mlr #{verb} -h`)
  end

  functions = `mlr help list-functions`
  print make_section('FUNCTIONS FOR FILTER/PUT', [
    ""
  ])
  functions = functions.strip.split("\n").uniq
  for function in functions
    print make_subsection(function, [])
    text = `mlr help function '#{function}'`
    text = text.sub(function + ' ', '')
    print make_code_block(text)
  end

  keywords = `mlr help list-keywords`
  print make_section('KEYWORDS FOR PUT AND FILTER', [
    ""
  ])
  keywords = keywords.strip.split("\n").uniq
  for keyword in keywords
    print make_subsection(keyword, [])
    text = `mlr help keyword '#{keyword}'`
    print make_code_block(text)
  end

  print make_section('AUTHOR', [
    "Miller is written by John Kerl <kerl.john.r@gmail.com>.",
    "This manual page has been composed from Miller's help output by Eric MSP Veith <eveith@veith-m.de>."
  ])
  print make_section('SEE ALSO', [
    "awk(1), sed(1), cut(1), join(1), sort(1), RFC 4180: Common Format and MIME Type for " +
    "Comma-Separated Values (CSV) Files, the Miller docsite https://miller.readthedocs.io"
  ])
end

# ================================================================
def make_top()
  t = Time::new
  stamp = t.gmtime.strftime("%Y-%m-%d")

  # Portability definitions thanks to some asciidoc output

"""'\\\" t
.\\\"     Title: mlr
.\\\"    Author: [see the \"AUTHOR\" section]
.\\\" Generator: #{$0}
.\\\"      Date: #{stamp}
.\\\"    Manual: \\ \\&
.\\\"    Source: \\ \\&
.\\\"  Language: English
.\\\"
.TH \"MILLER\" \"1\" \"#{stamp}\" \"\\ \\&\" \"\\ \\&\"
.\\\" -----------------------------------------------------------------
.\\\" * Portability definitions
.\\\" ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
.\\\" http://bugs.debian.org/507673
.\\\" http://lists.gnu.org/archive/html/groff/2009-02/msg00013.html
.\\\" ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
.ie \\n(.g .ds Aq \(aq
.el       .ds Aq '
.\\\" -----------------------------------------------------------------
.\\\" * set default formatting
.\\\" -----------------------------------------------------------------
.\\\" disable hyphenation
.nh
.\\\" disable justification (adjust text to left margin only)
.ad l
.\\\" -----------------------------------------------------------------
"""
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
def make_subsection(title, paragraphs)
  retval = ".SS \"#{title}\"\n"
  paragraphs.each do |paragraph|
    retval += ".sp\n"
    retval += groff_encode(paragraph) + "\n"
  end
  retval
end

# ----------------------------------------------------------------
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
  # In case the line starts with a dot:
  retval += block.gsub('\\', '\e').gsub(/^\./){'\&.'}
  # In case the line starts with a single quote:
  retval = retval.gsub(/^'/, '\(cq')
  retval += ".fi\n"
  retval += ".if n \\{\\\n"
  retval += ".RE\n"
end

# ----------------------------------------------------------------
def groff_encode(line)
  #line = line.gsub(/'/, '\(cq')
  #line = line.gsub(/"/, '\(dq')
  #line = line.gsub(/\./, '\&')
  #line = line.gsub(/-/, '\-')
  line = line.gsub(/\\([^-])/, '\e\1')
  line = line.gsub(/^\./){'\&.'}
  line
end

# ================================================================
main
