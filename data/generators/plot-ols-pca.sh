#!/bin/bash

file="$1"

# Set x_y_pca_m and x_y_pca_b as shell variables
mlr --ofs newline stats2 -a linreg-ols,linreg-pca -f x,y $file
eval $(mlr --ofs newline stats2 -a linreg-ols,linreg-pca -f x,y $file)

# In addition to x and y, make a new yfit which is the line fit. Plot using your favorite tool.
mlr --onidx put "\$olsfit=($x_y_ols_m*\$x)+$x_y_ols_b;\$pcafit=($x_y_pca_m*\$x)+$x_y_pca_b" \
  then cut -x -f a,b,i $file \
  | pgr -p -ms 2 -title 'linreg example' -xmin -0.1 -xmax 1.1 -ymin -0.1 -ymax 1.1 -legend 'y yfit'
