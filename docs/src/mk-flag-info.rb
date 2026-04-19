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
The letters `c`, `t`, `j`, `l`, `d`, `n`, `x`, `p`, `m`, and `y` refer to formats CSV, TSV, JSON, JSON Lines,
DKVP, NIDX, XTAB, PPRINT, markdown, and YAML, respectively. DCF is also supported (use `--dcf` for DCF in and out). Note that markdown format is available for
output only.

| In \\ out  | **CSV** | **TSV** | **JSON** | **JSONL** | **DKVP** | **NIDX** | **XTAB** | **PPRINT** | **Markdown** | **YAML** |
|------------|---------|---------|----------|-----------|----------|----------|----------|------------|--------------|----------|
| **CSV**    |         | `--c2t` | `--c2j`  | `--c2l`   | `--c2d`  | `--c2n`  | `--c2x`  | `--c2p`    | `--c2m`      | `--c2y`  |
| **TSV**    | `--t2c` |         | `--t2j`  | `--t2l`   | `--t2d`  | `--t2n`  | `--t2x`  | `--t2p`    | `--t2m`      | `--t2y`  |
| **JSON**   | `--j2c` | `--j2t` |          | `--j2l`   | `--j2d`  | `--j2n`  | `--j2x`  | `--j2p`    | `--j2m`      | `--j2y`  |
| **JSONL**  | `--l2c` | `--l2t` | `--l2j`  |           | `--l2d`  | `--l2n`  | `--l2x`  | `--l2p`    | `--l2m`      | `--l2y`  |
| **DKVP**   | `--d2c` | `--d2t` | `--d2j`  | `--d2l`   |          | `--d2n`  | `--d2x`  | `--d2p`    | `--d2m`      | `--d2y`  |
| **NIDX**   | `--n2c` | `--n2t` | `--n2j`  | `--n2l`   | `--n2d`  |          | `--n2x`  | `--n2p`    | `--n2m`      | `--n2y`  |
| **XTAB**   | `--x2c` | `--x2t` | `--x2j`  | `--x2l`   | `--x2d`  | `--x2n`  |          | `--x2p`    | `--x2m`      | `--x2y`  |
| **PPRINT** | `--p2c` | `--p2t` | `--p2j`  | `--p2l`   | `--p2d`  | `--p2n`  | `--p2x`  |            | `--p2m`      | `--p2y`  |
| **Markdown** | `--m2c` | `--m2t` | `--m2j`  | `--m2l`   | `--m2d`  | `--m2n`  | `--m2x`  | `--m2p`    |              | `--m2y`  |
| **YAML**   | `--y2c` | `--y2t` | `--y2j`  | `--y2l`   | `--y2d`  | `--y2n`  | `--y2x`  | `--y2p`    | `--y2m`      | `--y2y`  |

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
