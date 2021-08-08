map \d :w<C-m>:!clear;echo Building ...; echo; go build<C-m>
map \f :w<C-m>:!clear;echo Building ...; echo; build; echo; main<C-m>
map \T :w<C-m>:!clear;go test mlr/src/lib/...<C-m>
