mlr --opprint --from reg-test/input/s.dkvp --from reg-test/input/t.dkvp put '@idx = NR % 5; @idx = @idx == 0 ? 5 : @idx; $[@idx] = "NEW"'
