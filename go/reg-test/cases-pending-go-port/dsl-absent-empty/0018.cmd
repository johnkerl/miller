mlr --ofs tab put '$osum=@sum; $ostype=typeof(@sum);$xtype=typeof($x);@sum+=$x; $nstype=typeof(@sum);$nsum=@sum; end { emit @sum }' ./reg-test/cases-pending-go-port/dsl-absent-empty/0018.input
