#!/usr/bin/env ruby

# ================================================================
# Markdown autogen for Miller command-line flags.
# Invoked by files in *.md.in.
# ================================================================

section_names = `mlr help list-flag-sections`.split("\n")

for section_name in section_names
  puts "## #{section_name}"
  puts

  # TODO: split out info-print from per-flag
  if section_name =~ /.*conversion.*keystroke-saver*/
    # The markdown in this section looks a lot better when hand-crafted (thanks Nikos!).
    puts <<EOF
The letters `c`, `t`, `j`, `d`, `n`, `x`, `p`, and `m` refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,
PPRINT, and markdown, respectively. Note that markdown format is available for
output only.

| In \ out   | **CSV** | **TSV** | **JSON** | **DKVP** | **NIDX** | **XTAB** | **PPRINT** | **Markdown** |
|------------|---------|---------|----------|----------|----------|----------|------------|--------------|
| **CSV**    |         | `--c2t` | `--c2j`  | `--c2d`  | `--c2n`  | `--c2x`  | `--c2p`    | `--c2m`      |
| **TSV**    | `--t2c` |         | `--t2j`  | `--t2d`  | `--t2n`  | `--t2x`  | `--t2p`    | `--t2m`      |
| **JSON**   | `--j2c` | `--j2t` |          | `--j2d`  | `--j2n`  | `--j2x`  | `--j2p`    | `--j2m`      |
| **DKVP**   | `--d2c` | `--d2t` | `--d2j`  |          | `--d2n`  | `--d2x`  | `--d2p`    | `--d2m`      |
| **NIDX**   | `--n2c` | `--n2t` | `--n2j`  | `--n2d`  |          | `--n2x`  | `--n2p`    | `--n2m`      |
| **XTAB**   | `--x2c` | `--x2t` | `--x2j`  | `--x2d`  | `--x2n`  |          | `--x2p`    | `--x2m`      |
| **PPRINT** | `--p2c` | `--p2t` | `--p2j`  | `--p2d`  | `--p2n`  | `--p2x`  |            | `--p2m`      |

Additionally:

* `-p` is a keystroke-saver for `--nidx --fs space --repifs`.
* `-T` is a keystroke-saver for `--nidx --fs tab`.
EOF
  else
    # TODO: better formatting
    system("mlr help print-info-for-section '#{section_name}'")

    puts
    puts "**Flags:**"
    puts

    flags = `mlr help list-flags-for-section '#{section_name}'`.split
    for flag in flags
      headline = `mlr help show-headline-for-flag '#{flag}'`
      help = `mlr help show-help-for-flag '#{flag}'`
      puts "* `#{headline.chomp}`: #{help}"
    end
  end

  puts
end
