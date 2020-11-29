
run_mlr put -v 'begin{@a=1}; $e=2; $f==$g||$h==$i {};               $x=6; end{@z=9}' /dev/null
run_mlr put -v 'begin{@a=1}; $e=2; $f==$g||$h==$i {$s=1};           $x=6; end{@z=9}' /dev/null
run_mlr put -v 'begin{@a=1}; $e=2; $f==$g||$h==$i {$s=1;$t=2};      $x=6; end{@z=9}' /dev/null
run_mlr put -v 'begin{@a=1}; $e=2; $f==$g||$h==$i {$s=1;$t=2;$u=3}; $x=6; end{@z=9}' /dev/null
run_mlr put -v 'begin{@a=1}; $e=2; $f==$g||$h==$i {$s=1;@t["u".$5]=2;emit @v;emit @w; dump}; $x=6; end{@z=9}' /dev/null
run_mlr put -v 'begin{true{@x=1}}; true{@x=2}; end{true{@x=3}}' /dev/null
