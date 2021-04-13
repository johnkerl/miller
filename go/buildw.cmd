go build mlr.go

go test -v miller/src/...
# 'go test' (with no arguments) is the same as 'mlr regtest'

mlr regtest regtest/cases
