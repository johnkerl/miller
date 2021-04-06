mlr --opprint put -v 'end{@mean=@sum/@count; emitf @mean}; begin {@count=0; @sum=0.0}; @count=@count+1; @sum=@sum+$x' reg-test/input/abixy
