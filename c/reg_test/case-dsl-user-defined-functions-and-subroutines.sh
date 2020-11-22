run_mlr --opprint --from $indir/abixy put 'subr s(a,b) { $[a] = b } call s("W", 999)'

run_mlr --opprint --from $indir/abixy put '
  func f(x,y,z) {
    return x + y + z
  }
  subr s(a,b) {
      $[a] = b;
      $DID = "YES";
  }
  $o = f($x, $y, $i);
  call s("W", NR);
'

run_mlr --opprint --from $indir/abixy put '
  func f(x,y,z) {
    return x + y + z
  }
  subr s(a,b) {
      $[a] = b;
      return;
      $DID = "YES";
  }
  $o = f($x, $y, $i);
  call s("W", NR);
'

mlr_expect_fail --from $indir/abixy put '
  func f(x,y,z) {
    return x + y + z
  }
  subr s(a,b) {
      $[a] = b;
      return 1 # subr must not return value
  }
  $o = f($x, $y, $i);
  call s("W", NR);
'

mlr_expect_fail --from $indir/abixy put '
  func f(x,y,z) {
    return # func must return value
  }
  subr s(a,b) {
      $[a] = b;
  }
  $o = f($x, $y, $i);
  call s("W", NR);
'

# Test fencing: function f should not have access to boundvar k from the callsite.
run_mlr --from $indir/abixy --opprint put 'func f(x) {return k} for (k,v in $*) {$o=f($x)}'
run_mlr --from $indir/abixy --opprint put 'subr foo() {print "k is [".k."]"} for (k,v in $*) {call foo()}'

# Test overriding built-ins
mlr_expect_fail --opprint --from $indir/abixy put 'func log(x) { return 0 } $o = log($x)'

run_mlr --from $indir/abixy put 'subr log() { print "hello record  ". NR } call log()'

# No nesting of top-levels
mlr_expect_fail --from $indir/abixy --opprint put 'func f(x) { begin {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'func f(x) { end {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s(x) { begin {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s(x) { end {} }'

mlr_expect_fail --from $indir/abixy --opprint put 'func f(x) { func g(y) {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'func f(x) { subr t(y) {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s(x) { func g(y) {} }'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s(x) { subr t(y) {} }'

mlr_expect_fail --from $indir/abixy --opprint filter 'func f(x) { begin {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'func f(x) { end {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'subr s(x) { begin {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'subr s(x) { end {} }; true'

mlr_expect_fail --from $indir/abixy --opprint filter 'func f(x) { func g(y) {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'func f(x) { subr t(y) {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'subr s(x) { func g(y) {} }; true'
mlr_expect_fail --from $indir/abixy --opprint filter 'subr s(x) { subr t(y) {} }; true'

# Disallow redefines
mlr_expect_fail --from $indir/abixy --opprint put 'func log(x) { return $x + 1 }'
mlr_expect_fail --from $indir/abixy --opprint put 'func f(x) { return $x + 1 } func f(x) { return $x + 1}'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s() { } subr s() { }'
mlr_expect_fail --from $indir/abixy --opprint put 'subr s() { } subr s(x) { }'
run_mlr --from $indir/abixy --opprint put 'subr log(text) { print "TEXT IS ".text } call log("NR is ".NR)'

# scoping within distinct begin/end blocks
run_mlr --from $indir/abixy put -v '
    func f(x) {
        return x**2
    }
    func g(y) {
        return y+1
    }
    subr s(a,b,c) {
        print a.b.c
    }
    begin {
        @a = 0;
        var ell = 1;
        print "local1 = ".ell;
    }
    end {
        emit @a;
        var emm = 2;
        print "local2 = ".emm;
    }
    @a += 1;
    begin {
        @b = 0;
        @c = 0;
        print "local3 = ".ell;
    }
    @b += 2;
    @c += 3;
    end {
        emit @b;
        emit @c;
        print "local4 = ".emm;
    }
'

# print/dump from subr/func;  no tee/emit from func
run_mlr --from $indir/abixy --opprint put 'subr log(text) { print "TEXT IS ".text } call log("NR is ".NR)'
run_mlr --from $indir/abixy --opprint put 'func f(text) { print "TEXT IS ".text; return text.text } $o = f($a)'
run_mlr --from $indir/abixy put 'begin{@x=1} func f(x) { dump; print "hello"                 } $o=f($i)'
mlr_expect_fail --from $indir/abixy put 'begin{@x=1} func f(x) { dump; print "hello"; tee  > "x", $* } $o=f($i)'
mlr_expect_fail --from $indir/abixy put 'begin{@x=1} func f(x) { dump; print "hello"; emit > "x", @* } $o=f($i)' 
