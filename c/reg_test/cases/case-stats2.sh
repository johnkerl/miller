run_mlr --opprint stats2       -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2        $indir/abixy-wide
run_mlr --opprint stats2       -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2 -g a,b $indir/abixy-wide
run_mlr --oxtab   stats2 -s    -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2        $indir/abixy-wide-short
run_mlr --oxtab   stats2 -s    -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2 -g a,b $indir/abixy-wide-short
run_mlr --opprint stats2 --fit -a linreg-ols,linreg-pca             -f x,y,xy,y2        $indir/abixy-wide-short
run_mlr --opprint stats2 --fit -a linreg-ols,linreg-pca             -f x,y,xy,y2 -g a   $indir/abixy-wide-short

run_mlr --opprint stats2    -a logireg -f x,y      $indir/logi.dkvp
run_mlr --opprint stats2    -a logireg -f x,y -g g $indir/logi.dkvp

run_mlr --oxtab   stats2 -a cov -f x,y      $indir/abixy-het
run_mlr --oxtab   stats2 -a cov -f x,y -g a $indir/abixy-het
