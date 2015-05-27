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
run("CATMLR",    "mlr --no-mmap cat     data/big > /dev/null")
run("CATMLRM",   "mlr --mmap cat        data/big > /dev/null")
puts

run("(catv)",    "cat                       < data/big.csv > /dev/null")
run("catc0v",    "tools/catc0               < data/big.csv > /dev/null")
run("catc0v",    "tools/catc0                 data/big.csv > /dev/null")
run("catcv",     "tools/catc                < data/big.csv > /dev/null")
run("catcv",     "tools/catc                  data/big.csv > /dev/null")
run("catmv",     "tools/catm                  data/big.csv > /dev/null")
run("catawkv",   "awk -F, '{print}'     data/big.csv > /dev/null")
run("CATMLRV",   "mlr --no-mmap --csv cat      data/big.csv > /dev/null")
run("CATMLRVM",  "mlr --mmap --csv cat         data/big.csv > /dev/null")
puts
puts

run("cutcut",    "cut -d , -f 1,4              data/big > /dev/null")
run("cutawk",    "awk -F, '{print $1,$4}'      data/big > /dev/null")
run("CUTMLR",    "mlr --no-mmap cut -f a,x     data/big > /dev/null")
run("CUTMLRM",   "mlr --mmap cut -f a,x        data/big > /dev/null")
run("CUTMLRX",   "mlr --no-mmap cut -x -f a,x  data/big > /dev/null")
run("CUTMLRXM",  "mlr --mmap cut -x -f a,x     data/big > /dev/null")
puts

run("cutcutv",   "cut -d , -f 1,4                 data/big.csv > /dev/null")
run("cutawkv",   "awk -F, '{print $1,$4}'         data/big.csv > /dev/null")
run("CUTMLRV",   "mlr --no-mmap --csv cut -f a,x            data/big.csv > /dev/null")
run("CUTMLRVM",  "mlr --mmap    --csv cut -f a,x            data/big.csv > /dev/null")
run("CUTMLRXV",  "mlr --no-mmap --csv cut -x -f a,x         data/big.csv > /dev/null")
run("CUTMLRXVM", "mlr --mmap    --csv cut -x -f a,x         data/big.csv > /dev/null")
puts
puts

run("rensed",    "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/'   data/big > /dev/null")
run("RENMLR",    "mlr --no-mmap rename x,EKS,b,BEE      data/big > /dev/null")
run("RENMLRM",   "mlr --mmap rename x,EKS,b,BEE         data/big > /dev/null")
puts

run("rensedv",   "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/'     data/big.csv > /dev/null")
run("RENMLRV",   "mlr --no-mmap --csv rename x,EKS,b,BEE  data/big.csv > /dev/null")
run("RENMLRVM",  "mlr --mmap --csv rename x,EKS,b,BEE     data/big.csv > /dev/null")
puts
puts

run("addawk",    "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' data/big > /dev/null")
run("ADDMLR",    "mlr --no-mmap put '$z=$x+$y'            data/big > /dev/null")
run("ADDMLRM",   "mlr --mmap put '$z=$x+$y'               data/big > /dev/null")
puts

run("addawkv",   "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' data/big.csv > /dev/null")
run("ADDMLRV",   "mlr --no-mmap --csv put '$z=$x+$y'      data/big.csv > /dev/null")
run("ADDMLRVM",  "mlr --mmap --csv put '$z=$x+$y'         data/big.csv > /dev/null")
puts
puts

run("MEANMLR",    "mlr --no-mmap stats1 -a mean                  -f x,y -g a,b   data/big > /dev/null")
run("MEANMLRM",   "mlr --mmap    stats1 -a mean                  -f x,y -g a,b   data/big > /dev/null")
run('CORRMLR',    "mlr --no-mmap stats2 -a corr                  -f x,y -g a,b   data/big > /dev/null")
run('CORRMLRM',   "mlr --mmap    stats2 -a corr                  -f x,y -g a,b   data/big > /dev/null")
run('LINREGMLR',  "mlr --no-mmap stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   data/big > /dev/null")
run('LINREGMLRM', "mlr --mmap    stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   data/big > /dev/null")
puts

run("MEANMLRV",    "mlr --no-mmap --csv stats1 -a mean                  -f x,y -g a,b   data/big.csv > /dev/null")
run("MEANMLRVM",   "mlr --mmap    --csv stats1 -a mean                  -f x,y -g a,b   data/big.csv > /dev/null")
run('CORRMLRV',    "mlr --no-mmap --csv stats2 -a corr                  -f x,y -g a,b   data/big.csv > /dev/null")
run('CORRMLRVM',   "mlr --mmap    --csv stats2 -a corr                  -f x,y -g a,b   data/big.csv > /dev/null")
run('LINREGMLRV',  "mlr --no-mmap --csv stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   data/big.csv > /dev/null")
run('LINREGMLRVM', "mlr --mmap    --csv stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   data/big.csv > /dev/null")
puts
puts
puts

run("sortsort1", "sort -t= -k 1,2            data/big > /dev/null")
run("SORTMLR1",  "mlr --no-mmap sort -f a,b  data/big > /dev/null")
run("SORTMLR1M", "mlr --mmap sort -f a,b     data/big > /dev/null")
puts

run("sortsort2", "sort -t,    -k 4,5         data/big > /dev/null")
run("SORTMLR2",  "mlr --no-mmap sort -n x,y  data/big > /dev/null")
run("SORTMLR2M", "mlr --mmap sort -n x,y     data/big > /dev/null")
puts
puts

run("sortsort1v",  "sort -t, -k 1,4                  data/big.csv > /dev/null")
run("SORTMLR1V",   "mlr --no-mmap  --csv sort -f a,x data/big.csv > /dev/null")
run("SORTMLR1VM",  "mlr --mmap     --csv sort -f a,x data/big.csv > /dev/null")
puts

run("sortsort2v",  "sort -t, -n -k 4,5                data/big.csv > /dev/null")
run("SORTMLR2V",   "mlr  --no-mmap  --csv sort -n x,y data/big.csv > /dev/null")
run("SORTMLR2VM",  "mlr  --mmap     --csv sort -n x,y data/big.csv > /dev/null")
