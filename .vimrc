map \d :w<C-m>:!clear;echo Building ...; echo; make mlr<C-m>
map \f :w<C-m>:!clear;echo Building ...; echo; make ut<C-m>
map \r :w<C-m>:!clear;echo Building ...; echo; make ut-scan ut-mlv<C-m>
map \t :w<C-m>:!clear;go test github.com/johnkerl/miller/internal/pkg/transformers/...<C-m>
