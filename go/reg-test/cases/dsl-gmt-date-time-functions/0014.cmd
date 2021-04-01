mlr --opprint put '$hms=fsec2hms($sec);  $resec=hms2fsec($hms);  $diff=$resec-$sec' reg-test/input/fsec2xhms
