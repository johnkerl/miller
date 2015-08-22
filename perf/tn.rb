#!/usr/bin/ruby

require 'time'

# ================================================================
# Notes:
# * mlr --onidx --ofs ' ' filter '$experiment=="sort1b"' then cut -x -f experiment tn.out |pgr -nc -title cut -xlabel nlines -ylabel seconds -legend 'a b' -lop
#
# ...
# experiment=cat,nlines=100,mlrcat_seconds=3,awkcat_seconds=4
# experiment=cat,nlines=200,mlrcat_seconds=3,awkcat_seconds=4
# experiment=cat,nlines=300,mlrcat_seconds=3,awkcat_seconds=4
# experiment=cat,nlines=400,mlrcat_seconds=3,awkcat_seconds=4
# experiment=cut,nlines=100,mlrcut_seconds=3,awkcut_seconds=4
# experiment=cut,nlines=200,mlrcut_seconds=3,awkcut_seconds=4
# experiment=cut,nlines=300,mlrcut_seconds=3,awkcut_seconds=4
# experiment=cut,nlines=400,mlrcut_seconds=3,awkcut_seconds=4
# ...

# ================================================================
def run_progs_for_linecounts(experiment, linecounts, progdesc_map)
  linecounts.each do |linecount|
    run_progs_for_linecount(experiment, linecount, progdesc_map)
  end
end

def run_progs_for_linecount(experiment, linecount, progdesc_map)
  timing_fields = progdesc_map.collect do |progdesc, prog|
    seconds = run_prog(experiment, linecount, progdesc, prog)
    "#{progdesc}=#{seconds}"
  end
  puts "experiment=#{experiment},linecount=#{linecount},#{timing_fields.join(',')}"
end

def run_prog(experiment, linecount, progdesc, prog)
  filename = "data/nlines/#{linecount}"
  cmd = "#{prog} < #{filename} > /dev/null"
  seconds = run_cmd(progdesc, cmd)
  seconds
end

def run_cmd(desc, cmd)
	t1 = Time.new
	system(cmd)
  status = $?
	t2 = Time.new
	secs = t2.to_f - t1.to_f
  if status.to_i == 0
    secs
  else
	  'error'
  end
end

# ----------------------------------------------------------------
linecounts = [1, 100000, 200000, 300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000]

run_progs_for_linecounts('cat', linecounts, {
  "cat"    => "cat",
  "catc0"  => "catc0",
  "catc"   => "catc",
  "awkcat" => "awk -F , '{print}'",
  "mlrcat" => "mlr cat"
})

run_progs_for_linecounts('cut', linecounts, {
  "cut"    => "cut -d , -f 1,4",
  "awkcut" => "awk -F , '{print $1,$4}'",
  "mlrcat" => "mlr cut -f a,x"
})

run_progs_for_linecounts('add', linecounts, {
  "awkadd" => "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}'",
  "mlradd" => "mlr put '$z=$x+$y'"
})

run_progs_for_linecounts('sort2', linecounts, {
  "sort"    => "sort -t, -k 1,4",
  "mlrsort" => "mlr sort a,x"
})

run_progs_for_linecounts('sort1s', linecounts, {
  "sort"    => "sort -t, -k 1",
  "mlrsort" => "mlr sort a"
})

run_progs_for_linecounts('sort1b', linecounts, {
  "sort"    => "sort -t, -k 4",
  "mlrsort" => "mlr sort x"
})
