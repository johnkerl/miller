run_mlr --irs auto --ors lf cat $indir/line-term-lf.dkvp
run_mlr --irs auto --ors lf cat $indir/line-term-crlf.dkvp
run_mlr cat $indir/line-term-lf.dkvp
run_mlr cat $indir/line-term-crlf.dkvp

mention nidx
run_mlr --irs auto --ors lf --nidx --fs comma cat $indir/line-term-lf.dkvp
run_mlr --irs auto --ors lf --nidx --fs comma cat $indir/line-term-crlf.dkvp
run_mlr --nidx --fs comma cat $indir/line-term-lf.dkvp
run_mlr --nidx --fs comma cat $indir/line-term-crlf.dkvp


mention csvlite
run_mlr --irs auto --ors lf --csvlite cat $indir/line-term-lf.csv
run_mlr --irs auto --ors lf --csvlite cat $indir/line-term-crlf.csv
run_mlr --csvlite cat $indir/line-term-lf.csv
run_mlr --csvlite cat $indir/line-term-crlf.csv


mention pprint
run_mlr --irs auto --ors lf --pprint cat $indir/line-term-lf.csv
run_mlr --irs auto --ors lf --pprint cat $indir/line-term-crlf.csv
run_mlr --pprint cat $indir/line-term-lf.csv
run_mlr --pprint cat $indir/line-term-crlf.csv


mention xtab
run_mlr --ifs auto --xtab cat $indir/line-term-lf.xtab
run_mlr --ifs auto --xtab cat $indir/line-term-crlf.xtab
run_mlr --fs  auto --xtab cat $indir/line-term-lf.xtab
run_mlr --fs  auto --xtab cat $indir/line-term-crlf.xtab

mention xtab
run_mlr --ifs auto --xtab cat $indir/line-term-lf.xtab
run_mlr --ifs auto --xtab cat $indir/line-term-crlf.xtab
run_mlr --fs  auto --xtab cat $indir/line-term-lf.xtab
run_mlr --fs  auto --xtab cat $indir/line-term-crlf.xtab


mention csv
run_mlr --irs auto --ors lf --csv cat $indir/line-term-lf.csv
run_mlr --irs auto --ors lf --csv cat $indir/line-term-crlf.csv
run_mlr --csv cat $indir/line-term-lf.csv
run_mlr --csv cat $indir/line-term-crlf.csv


mention json nowrap nostack
run_mlr --irs auto --ors lf --json cat $indir/line-term-lf.json
run_mlr --irs auto --ors lf --json cat $indir/line-term-crlf.json
run_mlr --json cat $indir/line-term-lf.json
run_mlr --json cat $indir/line-term-crlf.json


mention json yeswrap nostack
run_mlr --irs auto --ors lf --jlistwrap --json cat $indir/line-term-lf-wrap.json
run_mlr --irs auto --ors lf --jlistwrap --json cat $indir/line-term-crlf-wrap.json
run_mlr --jlistwrap --json cat $indir/line-term-lf-wrap.json
run_mlr --jlistwrap --json cat $indir/line-term-crlf-wrap.json


mention json nowrap yesstack
run_mlr --irs auto --json --jvstack cat $indir/line-term-lf.json
run_mlr --irs auto --ors lf --json --jvstack cat $indir/line-term-crlf.json
run_mlr --json --jvstack cat $indir/line-term-lf.json
run_mlr --json --jvstack cat $indir/line-term-crlf.json


mention json yeswrap yesstack
run_mlr --irs auto --ors lf --jlistwrap --json --jvstack cat $indir/line-term-lf-wrap.json
run_mlr --irs auto --ors lf --jlistwrap --json --jvstack cat $indir/line-term-crlf-wrap.json
run_mlr --jlistwrap --json --jvstack cat $indir/line-term-lf-wrap.json
run_mlr --jlistwrap --json --jvstack cat $indir/line-term-crlf-wrap.json

mention json nowrap nostack
run_mlr --irs auto --ors lf --json cat $indir/line-term-lf.json
run_mlr --irs auto --ors lf --json cat $indir/line-term-crlf.json
run_mlr --json cat $indir/line-term-lf.json
run_mlr --json cat $indir/line-term-crlf.json


mention json yeswrap nostack
run_mlr --irs auto --ors lf --jlistwrap --json cat $indir/line-term-lf-wrap.json
run_mlr --irs auto --ors lf --jlistwrap --json cat $indir/line-term-crlf-wrap.json
run_mlr --jlistwrap --json cat $indir/line-term-lf-wrap.json
run_mlr --jlistwrap --json cat $indir/line-term-crlf-wrap.json


mention json nowrap yesstack
run_mlr --irs auto --ors lf --json --jvstack cat $indir/line-term-lf.json
run_mlr --irs auto --ors lf --json --jvstack cat $indir/line-term-crlf.json
run_mlr --json --jvstack cat $indir/line-term-lf.json
run_mlr --json --jvstack cat $indir/line-term-crlf.json


mention json yeswrap yesstack
run_mlr --irs auto --ors lf --jlistwrap --json --jvstack cat $indir/line-term-lf-wrap.json
run_mlr --irs auto --ors lf --jlistwrap --json --jvstack cat $indir/line-term-crlf-wrap.json
run_mlr --jlistwrap --json --jvstack cat $indir/line-term-lf-wrap.json
run_mlr --jlistwrap --json --jvstack cat $indir/line-term-crlf-wrap.json
