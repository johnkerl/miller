mlr --opprint put '$FIELD =~ "([A-Z]+)[^0-9]*([0-9]+)"'  then put '$F1="\1"; $F2="\2"; $F3="\3"' regtest/input/capture.dkvp
