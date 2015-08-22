#!/usr/bin/ruby

require 'time'

# ----------------------------------------------------------------
def run(desc, cmd)
	t1 = Time.new
	system(cmd)
  status = $?
	t2 = Time.new
	secs = t2.to_f - t1.to_f
  if status.to_i == 0
	  puts("%-12s %6.3f        %s" % [desc, secs, cmd])
  else
	  puts("%-12s %6s        %s" % [desc, "ERROR", cmd])
  end
end

# cutcut  real 0.38
# cutmlr  real 3.23
# cat     real 0.09
# catc    real 0.55

# ----------------------------------------------------------------


run("sortsort1", "sort -t, -k 1,2      < data/big > /dev/null")
#run("OSORTMLR1", "../c/mlr sort a,b    < data/big > /dev/null")
run("SORTMLR1",  "mlr sort a,b         < data/big > /dev/null")
puts

run("sortsort1v",  "sort -t, -k 1,4            < data/big.csv > /dev/null")
#run("OSORTMLR1V",  "../c/mlr --csv sort a,x    < data/big.csv > /dev/null")
run("SORTMLR1V",   "mlr      --csv sort a,x    < data/big.csv > /dev/null")
puts

run("sortsort2", "sort -t, -k 4,5      < data/big > /dev/null")
#run("OSORTMLR2", "../c/mlr sort x,y    < data/big > /dev/null")
run("SORTMLR2",  "mlr sort x,y         < data/big > /dev/null")
puts
puts

run("sortsort2v",  "sort -t, -k 4,5            < data/big.csv > /dev/null")
#run("OSORTMLR2V",  "../c/mlr --csv sort x,y    < data/big.csv > /dev/null")
run("SORTMLR2V",   "mlr      --csv sort x,y    < data/big.csv > /dev/null")
