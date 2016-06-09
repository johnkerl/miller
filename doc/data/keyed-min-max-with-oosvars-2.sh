mlr --opprint --from data/medium put -q '
  @output[$a]["x_min"] = min(@output[$a]["x_min"], $x);
  @output[$a]["x_max"] = max(@output[$a]["x_max"], $x);
  end{
    emit @output, "a"
  }
'
