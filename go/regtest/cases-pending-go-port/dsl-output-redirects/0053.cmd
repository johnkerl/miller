mlr --from regtest/input/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > ENV["outdir"]."/foo.dat", (@x, @y), "a"'
