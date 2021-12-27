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
run("(cat)",     "cat                 < ../data/big.dkvp > /dev/null")
run("catc0",     "catc0               < ../data/big.dkvp > /dev/null")
run("catc0",     "catc0                 ../data/big.dkvp > /dev/null")
run("catc",      "catc                < ../data/big.dkvp > /dev/null")
run("catc",      "catc                  ../data/big.dkvp > /dev/null")
run("catm",      "catm                  ../data/big.dkvp > /dev/null")
run("catgo",     "catgo                  ../data/big.dkvp > /dev/null")
run("catawk",    "awk -F, '{print}'     ../data/big.dkvp > /dev/null")
run("catmawk",   "mawk -F, '{print}'     ../data/big.dkvp > /dev/null")
run("CATMLR",    "../c/mlr --no-mmap cat     ../data/big.dkvp > /dev/null")
run("CATMLRM",   "../c/mlr --mmap cat        ../data/big.dkvp > /dev/null")
run("CATMLR",    "../go/mlr --no-mmap cat     ../data/big.dkvp > /dev/null")
run("CATMLRM",   "../go/mlr --mmap cat        ../data/big.dkvp > /dev/null")
puts

run("(catv)",    "cat                        < ../data/big.csv > /dev/null")
run("catc0v",    "catc0                      < ../data/big.csv > /dev/null")
run("catc0v",    "catc0                        ../data/big.csv > /dev/null")
run("catcv",     "catc                       < ../data/big.csv > /dev/null")
run("catcv",     "catc                         ../data/big.csv > /dev/null")
run("catmv",     "catm                         ../data/big.csv > /dev/null")
run("catgov",    "catgo                         ../data/big.csv > /dev/null")
run("catawkv",   "awk -F, '{print}'            ../data/big.csv > /dev/null")
run("catmawkv",  "mawk -F, '{print}'            ../data/big.csv > /dev/null")
run("CATMLRV",   "../c/mlr --no-mmap --csv --rs lf cat  ../data/big.csv > /dev/null")
run("CATMLRVM",  "../c/mlr --mmap --csv --rs lf cat     ../data/big.csv > /dev/null")
run("CATMLRV",   "../c/mlr --no-mmap --csvlite cat      ../data/big.csv > /dev/null")
run("CATMLRVM",  "../c/mlr --mmap --csvlite cat         ../data/big.csv > /dev/null")
puts
puts

run("cutcut",    "cut -d , -f 1,4              ../data/big.dkvp > /dev/null")
run("cutawk",    "awk -F, '{print $1,$4}'      ../data/big.dkvp > /dev/null")
run("cutmawk",   "mawk -F, '{print $1,$4}'      ../data/big.dkvp > /dev/null")
run("CUTMLR",    "mlr --no-mmap cut -f a,x     ../data/big.dkvp > /dev/null")
run("CUTMLRM",   "mlr --mmap cut -f a,x        ../data/big.dkvp > /dev/null")
run("CUTMLRX",   "mlr --no-mmap cut -x -f a,x  ../data/big.dkvp > /dev/null")
run("CUTMLRXM",  "mlr --mmap cut -x -f a,x     ../data/big.dkvp > /dev/null")
puts

run("cutcutv",   "cut -d , -f 1,4                 ../data/big.csv > /dev/null")
run("cutawkv",   "awk -F, '{print $1,$4}'         ../data/big.csv > /dev/null")
run("cutmawkv",  "mawk -F, '{print $1,$4}'         ../data/big.csv > /dev/null")
run("CUTMLRV",   "mlr --no-mmap --csv --rs lf cut -f a,x            ../data/big.csv > /dev/null")
run("CUTMLRVM",  "mlr --mmap    --csv --rs lf cut -f a,x            ../data/big.csv > /dev/null")
run("CUTMLRXV",  "mlr --no-mmap --csv --rs lf cut -x -f a,x         ../data/big.csv > /dev/null")
run("CUTMLRXVM", "mlr --mmap    --csv --rs lf cut -x -f a,x         ../data/big.csv > /dev/null")
run("CUTMLRV",   "mlr --no-mmap --csvlite cut -f a,x            ../data/big.csv > /dev/null")
run("CUTMLRVM",  "mlr --mmap    --csvlite cut -f a,x            ../data/big.csv > /dev/null")
run("CUTMLRXV",  "mlr --no-mmap --csvlite cut -x -f a,x         ../data/big.csv > /dev/null")
run("CUTMLRXVM", "mlr --mmap    --csvlite cut -x -f a,x         ../data/big.csv > /dev/null")
puts
puts

run("rensed",    "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/'   ../data/big.dkvp > /dev/null")
run("RENMLR",    "mlr --no-mmap rename x,EKS,b,BEE      ../data/big.dkvp > /dev/null")
run("RENMLRM",   "mlr --mmap rename x,EKS,b,BEE         ../data/big.dkvp > /dev/null")
puts

run("rensedv",   "sed -e 's/x=/EKS=/' -e 's/b=/BEE=/'     ../data/big.csv > /dev/null")
run("RENMLRV",   "mlr --no-mmap --csv --rs lf rename x,EKS,b,BEE  ../data/big.csv > /dev/null")
run("RENMLRVM",  "mlr --mmap --csv --rs lf rename x,EKS,b,BEE     ../data/big.csv > /dev/null")
run("RENMLRV",   "mlr --no-mmap --csvlite rename x,EKS,b,BEE  ../data/big.csv > /dev/null")
run("RENMLRVM",  "mlr --mmap --csvlite rename x,EKS,b,BEE     ../data/big.csv > /dev/null")
puts
puts

run("addawk",    "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' ../data/big.dkvp > /dev/null")
run("addmawk",   "mawk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' ../data/big.dkvp > /dev/null")
run("ADDMLR",    "mlr --no-mmap put '$z=$x+$y'            ../data/big.dkvp > /dev/null")
run("ADDMLRM",   "mlr --mmap put '$z=$x+$y'               ../data/big.dkvp > /dev/null")
puts

run("addawkv",   "awk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' ../data/big.csv > /dev/null")
run("addmawkv",  "mawk -F, '{gsub(\"x=\",\"\",$4);gsub(\"y=\",\"\",$5);print $4+$5}' ../data/big.csv > /dev/null")
run("ADDMLRV",   "mlr --no-mmap --csv --rs lf put '$z=$x+$y'      ../data/big.csv > /dev/null")
run("ADDMLRVM",  "mlr --mmap --csv --rs lf put '$z=$x+$y'         ../data/big.csv > /dev/null")
run("ADDMLRV",   "mlr --no-mmap --csvlite  put '$z=$x+$y'      ../data/big.csv > /dev/null")
run("ADDMLRVM",  "mlr --mmap --csvlite  put '$z=$x+$y'         ../data/big.csv > /dev/null")
puts
puts

run("MEANMLR",    "mlr --no-mmap stats1 -a mean                  -f x,y -g a,b   ../data/big.dkvp > /dev/null")
run("MEANMLRM",   "mlr --mmap    stats1 -a mean                  -f x,y -g a,b   ../data/big.dkvp > /dev/null")
run('CORRMLR',    "mlr --no-mmap stats2 -a corr                  -f x,y -g a,b   ../data/big.dkvp > /dev/null")
run('CORRMLRM',   "mlr --mmap    stats2 -a corr                  -f x,y -g a,b   ../data/big.dkvp > /dev/null")
run('LINREGMLR',  "mlr --no-mmap stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.dkvp > /dev/null")
run('LINREGMLRM', "mlr --mmap    stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.dkvp > /dev/null")
puts

run("MEANMLRV",    "mlr --no-mmap --csv --rs lf stats1 -a mean                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run("MEANMLRVM",   "mlr --mmap    --csv --rs lf stats1 -a mean                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('CORRMLRV',    "mlr --no-mmap --csv --rs lf stats2 -a corr                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('CORRMLRVM',   "mlr --mmap    --csv --rs lf stats2 -a corr                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('LINREGMLRV',  "mlr --no-mmap --csv --rs lf stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.csv > /dev/null")
run('LINREGMLRVM', "mlr --mmap    --csv --rs lf stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.csv > /dev/null")
run("MEANMLRV",    "mlr --no-mmap --csvlite stats1 -a mean                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run("MEANMLRVM",   "mlr --mmap    --csvlite stats1 -a mean                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('CORRMLRV',    "mlr --no-mmap --csvlite stats2 -a corr                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('CORRMLRVM',   "mlr --mmap    --csvlite stats2 -a corr                  -f x,y -g a,b   ../data/big.csv > /dev/null")
run('LINREGMLRV',  "mlr --no-mmap --csvlite stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.csv > /dev/null")
run('LINREGMLRVM', "mlr --mmap    --csvlite stats2 -a linreg-ols,linreg-pca -f x,y -g a,b   ../data/big.csv > /dev/null")
puts
puts
puts

run("sortsort1", "sort -t= -k 1,2            ../data/big.dkvp > /dev/null")
run("SORTMLR1",  "mlr --no-mmap sort -f a,b  ../data/big.dkvp > /dev/null")
run("SORTMLR1M", "mlr --mmap sort -f a,b     ../data/big.dkvp > /dev/null")
puts

run("sortsort2", "sort -t,    -k 4,5         ../data/big.dkvp > /dev/null")
run("SORTMLR2",  "mlr --no-mmap sort -n x,y  ../data/big.dkvp > /dev/null")
run("SORTMLR2M", "mlr --mmap sort -n x,y     ../data/big.dkvp > /dev/null")
puts
puts

run("sortsort1v",  "sort -t, -k 1,4                  ../data/big.csv > /dev/null")
run("SORTMLR1V",   "mlr --no-mmap  --csv --rs lf sort -f a,x ../data/big.csv > /dev/null")
run("SORTMLR1VM",  "mlr --mmap     --csv --rs lf sort -f a,x ../data/big.csv > /dev/null")
run("SORTMLR1V",   "mlr --no-mmap  --csvlite sort -f a,x ../data/big.csv > /dev/null")
run("SORTMLR1VM",  "mlr --mmap     --csvlite sort -f a,x ../data/big.csv > /dev/null")
puts

run("sortsort2v",  "sort -t, -n -k 4,5                ../data/big.csv > /dev/null")
run("SORTMLR2V",   "mlr  --no-mmap  --csv --rs lf sort -n x,y ../data/big.csv > /dev/null")
run("SORTMLR2VM",  "mlr  --mmap     --csv --rs lf sort -n x,y ../data/big.csv > /dev/null")
run("SORTMLR2V",   "mlr  --no-mmap  --csvlite sort -n x,y ../data/big.csv > /dev/null")
run("SORTMLR2VM",  "mlr  --mmap     --csvlite sort -n x,y ../data/big.csv > /dev/null")
