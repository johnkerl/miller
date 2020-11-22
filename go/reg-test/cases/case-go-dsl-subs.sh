echo x=eeee | run_mlr put '$y=ssub($x, "e", "X")'
echo x=eeee | run_mlr put '$y=sub($x, "e", "X")'
echo x=eeee | run_mlr put '$y=gsub($x, "e", "X")'
