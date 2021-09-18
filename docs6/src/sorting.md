<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# Sorting

## TBD

## TBD

## TBD

```
mlr --from s --ojson put '
  func forward(a,b) {
    return a <=> b
  }
  func reverse(a,b) {
    return b <=> a
  }
  func even_then_odd(a,b) {
    ax = a % 2;
    bx = b % 2;
    if (ax == bx) {
      return a <=> b
    } elif (bx == 1) {
      return -1
    } else {
      return 1
    }
  }
  $aa = sortaf([5,2,4,1,3], "forward");
  $bb = sortaf([5,2,4,1,3], "reverse");
  $cc = sortaf([$a, $b],  "forward");
  $dd = sortaf([$a, $b],  "reverse");
  $ee = sortaf([7, 0, 4, 2, 1, 9, 3, 8, 6, 5], "even_then_odd")
'
```

docpg: sort-in-row as well as sort-records-at-end.

* doc sorta, sortmk, sortaf, sortmf
* why call-UDF-by-name (no 1st-class for now, too much work -- so sortaf(a, "f") not sortaf(a, f))

also: example for sortaf/sortmf of structs ...

## TBD

```
mlr --opprint --from s put -q '
  begin {
    @records = []
  }
  func f(a,b) {
    return a.x+a.y <=> b.x+b.y
  }

  @records[NR] = $*;
  end {
    @records = sortaf(@records, "f");
    for (_, record in @records) {
      emit record;
    }
  }
' then put '$z=$x+$y'
```

```
mlr -n put -q '
  func f1(ak,av,bk,bv) {
    return ak <=> bk
  }
  func f2(ak,av,bk,bv) {
    return bk <=> ak
  }
  func f3(ak,av,bk,bv) {
    return av <=> bv
  }
  func f4(ak,av,bk,bv) {
    return bv <=> av
  }
  end {
    x = {
      "c":1,
      "a":3,
      "b":2,
    };
    dump x;

    print sortmf(x, "f1");
    print sortmf(x, "f2");
    print sortmf(x, "f3");
    print sortmf(x, "f4");
  }
'
```
