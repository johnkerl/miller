mlr --opprint --from data/medium put -q '
  @min[$a] = min(@min[$a], $x);
  @max[$a] = max(@max[$a], $x);
  end{
    emit (@min, @max), "a";
  }
'
