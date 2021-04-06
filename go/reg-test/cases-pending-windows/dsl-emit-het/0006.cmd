mlr --from reg-test/input/abixy-het put -q '@x={"a":NR}; @y={"b":-NR}; emit (@x, @y), "k"'
