run_mlr cat $indir/nest-explode.dkvp
run_mlr cat $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --values --across-fields  -f x $indir/nest-explode.dkvp
run_mlr nest --explode --values --across-fields  -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --values --across-fields  -f x then nest --implode --values --across-fields  -f x $indir/nest-explode.dkvp
run_mlr nest --explode --values --across-fields  -f x --nested-fs pipe --nested-ps = then nest --implode --values --across-fields  -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --values --across-records -f x $indir/nest-explode.dkvp
run_mlr nest --explode --values --across-records -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --values --across-records -f x then nest --implode --values --across-records -f x $indir/nest-explode.dkvp
run_mlr nest --explode --values --across-records -f x --nested-fs pipe --nested-ps = then nest --implode --values --across-records -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --pairs  --across-fields  -f x $indir/nest-explode.dkvp
run_mlr nest --explode --pairs  --across-fields  -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp

run_mlr nest --explode --pairs  --across-records -f x $indir/nest-explode.dkvp
run_mlr nest --explode --pairs  --across-records -f x --nested-fs pipe --nested-ps = $indir/nest-explode-vary-fs-ps.dkvp
