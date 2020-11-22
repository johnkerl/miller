echo '{"x":1}'                                           | run_mlr --json cat
echo '{"x":[1,2,3]}'                                     | run_mlr --json cat
echo '{"x":[1,[2,3,4],5]}'                               | run_mlr --json cat
echo '{"x":[1,[2,[3,4,5],6],7]}'                         | run_mlr --json cat

echo '{"x":{}}'                                          | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3}}'                         | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}'       | run_mlr --json cat

echo '{"x":{},"y":1}'                                    | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3},"y":4}'                   | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}' | run_mlr --json cat

echo '{"x":1}'                                           | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,2,3]}'                                     | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,[2,3,4],5]}'                               | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,[2,[3,4,5],6],7]}'                         | run_mlr --json cat | run_mlr --json cat

echo '{"x":{}}'                                          | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3}}'                         | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}'       | run_mlr --json cat | run_mlr --json cat

echo '{"x":{},"y":1}'                                    | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3},"y":4}'                   | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}' | run_mlr --json cat | run_mlr --json cat
