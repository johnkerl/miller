# If you edit this file, make sure to keep the space after the b.
#
# Note that echoing input piped to run_mr would be much nicer, but, I had
# issues with the $num_invocations_attempted not tracking correctly in
# reg-test/run when run_mlr is at the end of a pipe :(. This no-echo rule is
# enforced by reg-test/run.

run_mlr --ojson put '$y = strip($x)' <<EOF
x= a     b 
EOF
run_mlr --ojson put '$y = lstrip($x)' <<EOF
x= a     b 
EOF
run_mlr --ojson put '$y = rstrip($x)' <<EOF
x= a     b 
EOF
run_mlr --ojson put '$y = collapse_whitespace($x)' <<EOF
x= a     b 
EOF
run_mlr --ojson put '$y = clean_whitespace($x)' <<EOF
x= a     b 
EOF
