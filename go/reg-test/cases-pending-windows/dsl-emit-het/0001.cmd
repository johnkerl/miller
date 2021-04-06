mlr --from reg-test/input/abixy put -q '@x={"a":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
