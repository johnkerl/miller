mlr --oxtab stats1 -f x -a p25,p75 \
    then put '$x_iqr = $x_p75 - $x_p25' \
    data/medium 
