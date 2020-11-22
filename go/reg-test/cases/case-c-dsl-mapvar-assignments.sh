
# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO FULL SREC

run_mlr --from $indir/xyz2 put '$* = {"a":1, "b": {"x":8,"y":9}, "c":3}'
run_mlr --from $indir/xyz2 put '$* = $*'
run_mlr --from $indir/xyz2 put '@a = 1; $* = @a'
run_mlr --from $indir/xyz2 put '@b[1] = 2; $* = @b'
run_mlr --from $indir/xyz2 put '@c[1][2] = 3; $* = @c'
run_mlr --from $indir/xyz2 put '@a = 1; $* = @*'
run_mlr --from $indir/xyz2 put '@b[1] = 2; $* = @*'
run_mlr --from $indir/xyz2 put '@c[1][2] = 3; $* = @*'
run_mlr --from $indir/xyz2 put 'a = 1; $* = a'
run_mlr --from $indir/xyz2 put 'b[1] = 2; $* = b'
run_mlr --from $indir/xyz2 put 'c[1][2] = 3; $* = c'
run_mlr --from $indir/xyz2 put '$* = 3'

run_mlr --from $indir/xyz2 put '
  func map_valued_func() {
    return {"a":1,"b":2}
  }
  map m = map_valued_func();
  $* = m;
'

run_mlr --from $indir/xyz2 put '
  func map_valued_func() {
    return {"a":1,"b":2}
  }
  $* = map_valued_func();
'

# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO FULL OOSVAR

run_mlr --from $indir/xyz2 put '@* = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump'
run_mlr --from $indir/xyz2 put '@* = $*; dump'
run_mlr --from $indir/xyz2 put '@* = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump; @* = @*; dump'
run_mlr --from $indir/xyz2 put '@a = 1; @* = @a; dump'
run_mlr --from $indir/xyz2 put '@b[1] = 2; @* = @b; dump'
run_mlr --from $indir/xyz2 put '@c[1][2] = 3; @* = @c; dump'
run_mlr --from $indir/xyz2 put '@a = 1; @* = @*; dump'
run_mlr --from $indir/xyz2 put '@b[1] = 2; @* = @*; dump'
run_mlr --from $indir/xyz2 put '@c[1][2] = 3; @* = @*; dump'
run_mlr --from $indir/xyz2 put 'a = 1; @* = a; dump'
run_mlr --from $indir/xyz2 put 'b[1] = 2; @* = b; dump'
run_mlr --from $indir/xyz2 put 'c[1][2] = 3; @* = c; dump'
run_mlr --from $indir/xyz2 put '@* = 3'

run_mlr --from $indir/xyz2 put '
  func map_valued_func() {
    return {"a":1,"b":2}
  }
  map m = map_valued_func();
  @* = m;
  dump
'

run_mlr --from $indir/xyz2 put '
  func map_valued_func() {
    return {"a":1,"b":2}
  }
  @* = map_valued_func();
  dump
'

# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO OOSVAR

run_mlr --from $indir/xyz2 put -q '@o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump @o'
run_mlr --from $indir/xyz2 put -q '@o = @o; dump @o'
run_mlr --from $indir/xyz2 put -q '@o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump; @o = @*; dump'
run_mlr --from $indir/xyz2 put -q '@o = $*; dump @o'
run_mlr --from $indir/xyz2 put -q '@o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump @o; @o = @o; dump @o'
run_mlr --from $indir/xyz2 put -q '@o = {"a":1, "b": {"x":8,"y":9}, "c":3}; @o = @o; dump @o'
run_mlr --from $indir/xyz2 put -q '@a = 1; @o = @a; dump @o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; @o = @b; dump @o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; @o = @c; dump @o'
run_mlr --from $indir/xyz2 put -q '@a = 1; @o = @*; dump @o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; @o = @*; dump @o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; @o = @*; dump @o'
run_mlr --from $indir/xyz2 put -q 'a = 1; @o = a; dump @o'
run_mlr --from $indir/xyz2 put -q 'b[1] = 2; @o = b; dump @o'
run_mlr --from $indir/xyz2 put -q 'c[1][2] = 3; @o = c; dump @o'
run_mlr --from $indir/xyz2 put -q 'func map_valued_func() { return {"a":1,"b":2}} @o = map_valued_func(); dump @o'

# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO MAP-TYPED LOCAL

run_mlr         --from $indir/xyz2 put -q 'map o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump o'
mlr_expect_fail --from $indir/xyz2 put -q 'map o = @o; dump o'
run_mlr         --from $indir/xyz2 put -q 'map o = $*; dump o'
run_mlr         --from $indir/xyz2 put -q 'map o = {"a":1, "b": {"x":8,"y":9}, "c":3}; o = o; dump o'
mlr_expect_fail --from $indir/xyz2 put -q '@a = 1; map o = @a; dump o'
run_mlr         --from $indir/xyz2 put -q '@b[1] = 2; map o = @b; dump o'
run_mlr         --from $indir/xyz2 put -q '@c[1][2] = 3; map o = @c; dump o'
run_mlr         --from $indir/xyz2 put -q '@a = 1; map o = @*; dump o'
run_mlr         --from $indir/xyz2 put -q '@b[1] = 2; map o = @*; dump o'
run_mlr         --from $indir/xyz2 put -q '@c[1][2] = 3; map o = @*; dump o'
mlr_expect_fail --from $indir/xyz2 put -q 'a = 1; map o = a; dump o'
run_mlr         --from $indir/xyz2 put -q 'b[1] = 2; map o = b; dump o'
run_mlr         --from $indir/xyz2 put -q 'c[1][2] = 3; map o = c; dump o'
run_mlr --from $indir/xyz2 put -q 'func map_valued_func() { return {"a":1,"b":2}} map o = map_valued_func(); dump o'

# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO VAR-TYPED LOCAL

run_mlr --from $indir/xyz2 put -q 'var o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump o'
run_mlr --from $indir/xyz2 put -q 'var o = @o; dump o'
run_mlr --from $indir/xyz2 put -q 'var o = $*; dump o'
run_mlr --from $indir/xyz2 put -q 'var o = {"a":1, "b": {"x":8,"y":9}, "c":3}; o = o; dump o'
run_mlr --from $indir/xyz2 put -q '@a = 1; var o = @a; dump o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; var o = @b; dump o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; var o = @c; dump o'
run_mlr --from $indir/xyz2 put -q '@a = 1; var o = @*; dump o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; var o = @*; dump o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; var o = @*; dump o'
run_mlr --from $indir/xyz2 put -q 'a = 1; var o = a; dump o'
run_mlr --from $indir/xyz2 put -q 'b[1] = 2; var o = b; dump o'
run_mlr --from $indir/xyz2 put -q 'c[1][2] = 3; var o = c; dump o'
run_mlr --from $indir/xyz2 put -q 'func map_valued_func() { return {"a":1,"b":2}} var o = map_valued_func(); dump o'

# ----------------------------------------------------------------
announce MAPVAR ASSIGNMENTS TO UNTYPED LOCAL

run_mlr --from $indir/xyz2 put -q 'o = {"a":1, "b": {"x":8,"y":9}, "c":3}; dump o'
run_mlr --from $indir/xyz2 put -q 'o = @o; dump o'
run_mlr --from $indir/xyz2 put -q 'o = $*; dump o'
run_mlr --from $indir/xyz2 put -q 'o = {"a":1, "b": {"x":8,"y":9}, "c":3}; o = o; dump o'
run_mlr --from $indir/xyz2 put -q '@a = 1; o = @a; dump o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; o = @b; dump o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; o = @c; dump o'
run_mlr --from $indir/xyz2 put -q '@a = 1; o = @*; dump o'
run_mlr --from $indir/xyz2 put -q '@b[1] = 2; o = @*; dump o'
run_mlr --from $indir/xyz2 put -q '@c[1][2] = 3; o = @*; dump o'
run_mlr --from $indir/xyz2 put -q 'a = 1; o = a; dump o'
run_mlr --from $indir/xyz2 put -q 'b[1] = 2; o = b; dump o'
run_mlr --from $indir/xyz2 put -q 'c[1][2] = 3; o = c; dump o'
run_mlr --from $indir/xyz2 put -q 'func map_valued_func() { return {"a":1,"b":2}} o = map_valued_func(); dump o'
