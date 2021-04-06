# TODO: prepare manually

run_mlr --opprint --from $indir/ten.dkvp put -e 'begin {@a=0}' -e 'begin {@b=1}' -e '
  @c = @a + @b;
  $fa = @a;
  $fb = @b;
  $fc = @c;
  @a = @b;
  @b = @c;
'
