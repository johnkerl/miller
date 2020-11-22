run_mlr --opprint --from $indir/ten.dkvp put -s a=0 -s b=1 '
  @c = @a + @b;
  $fa = @a;
  $fb = @b;
  $fc = @c;
  @a = @b;
  @b = @c;
'
run_mlr --opprint --from $indir/s.dkvp put          '@sum += 1; $z=@sum'
run_mlr --opprint --from $indir/s.dkvp put -s sum=0 '@sum += 1; $z=@sum'
run_mlr --opprint --from $indir/s.dkvp put -s sum=8 '@sum += 1; $z=@sum'
