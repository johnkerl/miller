# ----------------------------------------------------------------

run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[0]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[1]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[2]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[3]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[4]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[5]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[6]; dump $v'

run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[0]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[1]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[2]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[3]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[4]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[5]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[6]; dump @v }'

run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[0]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[1]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[2]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[3]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[4]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[5]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[6]; dump v }'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"]["x"];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1]["x"];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"][2];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1][2];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"]["x"];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1]["x"];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"][2];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1][2];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*;
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"];
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"]["x"];
    dump $*;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*;
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"];
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"][2];
    dump $*;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"]["x"];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1]["x"];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"][2];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1][2];
    dump v;
  }
'
