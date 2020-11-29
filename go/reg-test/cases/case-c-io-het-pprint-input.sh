run_mlr --ipprint --odkvp cat $indir/a.pprint
run_mlr --ipprint --odkvp cat $indir/b.pprint
run_mlr --ipprint --odkvp cat $indir/c.pprint
run_mlr --ipprint --odkvp cat $indir/d.pprint
run_mlr --ipprint --odkvp cat $indir/e.pprint
run_mlr --ipprint --odkvp cat $indir/f.pprint
run_mlr --ipprint --odkvp cat $indir/g.pprint

run_mlr --ipprint --odkvp cat $indir/a.pprint $indir/a.pprint
run_mlr --ipprint --odkvp cat $indir/b.pprint $indir/b.pprint
run_mlr --ipprint --odkvp cat $indir/c.pprint $indir/c.pprint
run_mlr --ipprint --odkvp cat $indir/d.pprint $indir/d.pprint
run_mlr --ipprint --odkvp cat $indir/e.pprint $indir/e.pprint
run_mlr --ipprint --odkvp cat $indir/f.pprint $indir/f.pprint
run_mlr --ipprint --odkvp cat $indir/g.pprint $indir/g.pprint

run_mlr --ipprint --odkvp cat $indir/a.pprint $indir/b.pprint
run_mlr --ipprint --odkvp cat $indir/b.pprint $indir/c.pprint
run_mlr --ipprint --odkvp cat $indir/c.pprint $indir/d.pprint
run_mlr --ipprint --odkvp cat $indir/d.pprint $indir/e.pprint
run_mlr --ipprint --odkvp cat $indir/e.pprint $indir/f.pprint
run_mlr --ipprint --odkvp cat $indir/f.pprint $indir/g.pprint

run_mlr --ipprint --odkvp cat $indir/a.pprint $indir/b.pprint \
  $indir/c.pprint $indir/d.pprint $indir/e.pprint $indir/f.pprint $indir/g.pprint
