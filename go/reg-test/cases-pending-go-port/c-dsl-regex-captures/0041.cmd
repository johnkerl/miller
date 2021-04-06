mlr --inidx --odkvp put '$1 =~ "ab(c)?d(..)g" { $c1 = "\1"; $c2 = "\2"}' ./reg-test/cases-pending-go-port/c-dsl-regex-captures/0041.input
