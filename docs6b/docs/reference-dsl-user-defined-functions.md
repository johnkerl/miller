<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DSL reference: user-defined functions

As of Miller 5.0.0 you can define your own functions, as well as subroutines.

## User-defined functions

Here's the obligatory example of a recursive function to compute the factorial function:

<pre>
<b>mlr --opprint --from data/small put '</b>
<b>    func f(n) {</b>
<b>        if (is_numeric(n)) {</b>
<b>            if (n > 0) {</b>
<b>                return n * f(n-1);</b>
<b>            } else {</b>
<b>                return 1;</b>
<b>            }</b>
<b>        }</b>
<b>        # implicitly return absent-null if non-numeric</b>
<b>    }</b>
<b>    $ox = f($x + NR);</b>
<b>    $oi = f($i);</b>
<b>'</b>
a   b   i x                   y                   ox                  oi
pan pan 1 0.3467901443380824  0.7268028627434533  0.46705354854811026 1
eks pan 2 0.7586799647899636  0.5221511083334797  3.680838410072862   2
wye wye 3 0.20460330576630303 0.33831852551664776 1.7412511955594865  6
eks wye 4 0.38139939387114097 0.13418874328430463 18.588348778962008  24
wye pan 5 0.5732889198020006  0.8636244699032729  211.38730958519247  120
</pre>

Properties of user-defined functions:

* Function bodies start with ``func`` and a parameter list, defined outside of ``begin``, ``end``, or other ``func`` or ``subr`` blocks. (I.e. the Miller DSL has no nested functions.)

* A function (uniqified by its name) may not be redefined: either by redefining a user-defined function, or by redefining a built-in function. However, functions and subroutines have separate namespaces: you can define a subroutine ``log`` which does not clash with the mathematical ``log`` function.

* Functions may be defined either before or after use (there is an object-binding/linkage step at startup).  More specifically, functions may be either recursive or mutually recursive. Functions may not call subroutines.

* Functions may be defined and called either within ``mlr put`` or ``mlr put``.

* Functions have read access to ``$``-variables and ``@``-variables but may not modify them. See also :ref:`cookbook-memoization-with-oosvars` for an example.

* Argument values may be reassigned: they are not read-only.

* When a return value is not implicitly returned, this results in a return value of absent-null. (In the example above, if there were records for which the argument to ``f`` is non-numeric, the assignments would be skipped.) See also the section on :doc:`reference-main-null-data`.

* See the section on :ref:`reference-dsl-local-variables` for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on :ref:`reference-dsl-expressions-from-files` for information on the use of ``-f`` and ``-e`` flags.

## User-defined subroutines

Example:

<pre>
<b>mlr --opprint --from data/small put -q '</b>
<b>  begin {</b>
<b>    @call_count = 0;</b>
<b>  }</b>
<b>  subr s(n) {</b>
<b>    @call_count += 1;</b>
<b>    if (is_numeric(n)) {</b>
<b>      if (n > 1) {</b>
<b>        call s(n-1);</b>
<b>      } else {</b>
<b>        print "numcalls=" . @call_count;</b>
<b>      }</b>
<b>    }</b>
<b>  }</b>
<b>  print "NR=" . NR;</b>
<b>  call s(NR);</b>
<b>'</b>
NR=1
numcalls=1
NR=2
numcalls=3
NR=3
numcalls=6
NR=4
numcalls=10
NR=5
numcalls=15
</pre>

Properties of user-defined subroutines:

* Subroutine bodies start with ``subr`` and a parameter list, defined outside of ``begin``, ``end``, or other ``func`` or ``subr`` blocks. (I.e. the Miller DSL has no nested subroutines.)

* A subroutine (uniqified by its name) may not be redefined. However, functions and subroutines have separate namespaces: you can define a subroutine ``log`` which does not clash with the mathematical ``log`` function.

* Subroutines may be defined either before or after use (there is an object-binding/linkage step at startup).  More specifically, subroutines may be either recursive or mutually recursive. Subroutines may call functions.

* Subroutines may be defined and called either within ``mlr put`` or ``mlr put``.

* Subroutines have read/write access to ``$``-variables and ``@``-variables.

* Argument values may be reassigned: they are not read-only.

* See the section on :ref:`reference-dsl-local-variables` for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on :ref:`reference-dsl-expressions-from-files` for information on the use of ``-f`` and ``-e`` flags.
