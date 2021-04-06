mlr --from reg-test/input/abixy-het put -q '@x={"a":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
