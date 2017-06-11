mlr --oxtab stats1 --fr '[i-z]' -a p25,p75 \
    then put 'for (k,v in $*) {if (k =~ "(.*)_p25") {$["\1_iqr"] =$["\1_p75"] - $["\1_p25"]}}' \
    data/medium 
