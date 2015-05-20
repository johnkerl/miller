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
run("(cat)",     "cat                 < data/big > /dev/null")
run("catc0",     "tools/catc0         < data/big > /dev/null")
run("catc0",     "tools/catc0           data/big > /dev/null")
run("catc",      "tools/catc          < data/big > /dev/null")
run("catc",      "tools/catc            data/big > /dev/null")
run("catm",      "tools/catm            data/big > /dev/null")
run("catawk",    "awk -F, '{print}'     data/big > /dev/null")
run("CATMLR",    "mlr cat               data/big > /dev/null")
run("CATMLRM",   "mlr --mmap cat        data/big > /dev/null")
puts

run("(catv)",    "cat                 < data/big.csv > /dev/null")
run("catc0v",    "tools/catc0               < data/big.csv > /dev/null")
run("catc0v",    "tools/catc0                 data/big.csv > /dev/null")
run("catcv",     "tools/catc                < data/big.csv > /dev/null")
run("catcv",     "tools/catc                  data/big.csv > /dev/null")
run("catmv",     "tools/catm                  data/big.csv > /dev/null")
run("catawkv",   "awk -F, '{print}'   < data/big.csv > /dev/null")
run("CATMLRV",   "mlr --csv cat       < data/big.csv > /dev/null")
puts
puts

run("cutcut",    "cut -d , -f 1,4            data/big > /dev/null")
run("cutawk",    "awk -F, '{print $1,$4}'    data/big > /dev/null")
run("CUTMLR",    "mlr cut -f a,x             data/big > /dev/null")
run("CUTMLRM",   "mlr --mmap cut -f a,x      data/big > /dev/null")
run("CUTMLRX",   "mlr cut -x -f a,x          data/big > /dev/null")
run("CUTMLRXM",  "mlr --mmap cut -x -f a,x   data/big > /dev/null")
puts

run("cutcutv",   "cut -d , -f 1,4               < data/big.csv > /dev/null")
run("cutawkv",   "awk -F, '{print $1,$4}'       < data/big.csv > /dev/null")
run("CUTMLRV",   "mlr --csv cut -f a,x          < data/big.csv > /dev/null")
run("CUTMLRXV",  "mlr --csv cut -x -f a,x       < data/big.csv > /dev/null")
puts
puts

run("rensed",    "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/' < data/big > /dev/null")
run("RENMLR",    "mlr rename x,EKS,b,BEE              < data/big > /dev/null")
puts

run("rensedv",   "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/' < data/big.csv > /dev/null")
run("RENMLRV",   "mlr --csv rename x,EKS,b,BEE        < data/big.csv > /dev/null")
puts
puts

run("addawk",    "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' < data/big > /dev/null")
run("ADDMLR",    "mlr put '$z=$x+$y'                  < data/big > /dev/null")
puts

run("addawkv",   "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' < data/big.csv > /dev/null")
run("ADDMLRV",   "mlr --csv put '$z=$x+$y'            < data/big.csv > /dev/null")
puts
puts

run("MEANMLR",   "mlr       stats1 -a mean -f x,y -g a,b < data/big > /dev/null")
run('CORRMLR',   "mlr       stats2 -a corr -f x,y -g a,b < data/big > /dev/null")
run('LINREGMLR', "mlr       stats2 -a linreg-ols,linreg-pca -f x,y -g a,b < data/big > /dev/null")
puts

run("MEANMLRV",  "mlr      --csv  stats1 -a mean -f x,y -g a,b < data/big.csv > /dev/null")
run('CORRMLRV',  "mlr      --csv  stats2 -a corr -f x,y -g a,b < data/big.csv > /dev/null")
run('LINREGMLRV', "mlr     --csv  stats2 -a linreg-ols,linreg-pca -f x,y -g a,b < data/big.csv > /dev/null")
puts
puts

run("sortsort1", "sort -t= -k 1,2      < data/big > /dev/null")
run("SORTMLR1",  "mlr sort -f a,b      < data/big > /dev/null")
puts

run("sortsort2", "sort -t,    -k 4,5   < data/big > /dev/null")
run("SORTMLR2",  "mlr sort -n x,y      < data/big > /dev/null")
puts
puts

run("sortsort1v",  "sort -t, -k 1,4            < data/big.csv > /dev/null")
run("SORTMLR1V",   "mlr      --csv sort -f a,x < data/big.csv > /dev/null")
puts

run("sortsort2v",  "sort -t, -n -k 4,5         < data/big.csv > /dev/null")
run("SORTMLR2V",   "mlr      --csv sort -n x,y < data/big.csv > /dev/null")
