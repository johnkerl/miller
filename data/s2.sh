#!/bin/bash
mlr --opprint --right "$@" stats2 -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2,x2,x2 -g a,b ../data/mediumwide
echo
mlr --oxtab --right "$@" stats2 -a linreg-ols,linreg-pca,r2,corr,cov -f x,y,xy,y2,x2,x2 ../data/mediumwide
