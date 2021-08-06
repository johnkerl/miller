go build

go test -v mlr/src/...
# 'go test' (with no arguments) is the same as 'mlr regtest'

mlr regtest regtest/cases
