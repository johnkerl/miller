# Intended to be invoked by "." from reg_test/run

run_mlr having-fields --at-least  a,b         $indir/abixy
run_mlr having-fields --at-least  a,c         $indir/abixy
run_mlr having-fields --at-least  a,b,i,x,y   $indir/abixy
run_mlr having-fields --which-are a,b,i,x     $indir/abixy
run_mlr having-fields --which-are a,b,i,x,y   $indir/abixy
run_mlr having-fields --which-are a,b,i,y,x   $indir/abixy
run_mlr having-fields --which-are a,b,i,x,w   $indir/abixy
run_mlr having-fields --which-are a,b,i,x,y,z $indir/abixy
run_mlr having-fields --at-most   a,c         $indir/abixy
run_mlr having-fields --at-most   a,b,i,x,y   $indir/abixy
run_mlr having-fields --at-most   a,b,i,x,y,z $indir/abixy

run_mlr having-fields --all-matching  '"^[a-z][a-z][a-z]$"'  $indir/having-fields-regex.dkvp
run_mlr having-fields --any-matching  '"^[a-z][a-z][a-z]$"'  $indir/having-fields-regex.dkvp
run_mlr having-fields --none-matching '"^[a-z][a-z][a-z]$"'  $indir/having-fields-regex.dkvp
run_mlr having-fields --all-matching  '"^[a-z][a-z][a-z]$"i' $indir/having-fields-regex.dkvp
run_mlr having-fields --any-matching  '"^[a-z][a-z][a-z]$"i' $indir/having-fields-regex.dkvp
run_mlr having-fields --none-matching '"^[a-z][a-z][a-z]$"i' $indir/having-fields-regex.dkvp
