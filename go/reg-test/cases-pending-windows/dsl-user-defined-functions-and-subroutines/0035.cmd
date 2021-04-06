mlr --from reg-test/input/abixy put 'begin{@x=1} func f(x) { dump; print "hello"; tee  > ENV["outdir"]."/udf-x", $* } $o=f($i)'
