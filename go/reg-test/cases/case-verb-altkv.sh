run_mlr altkv <<EOF
a,b,c,d,e,f
EOF

run_mlr altkv <<EOF
a,b,c,d,e,f,g
EOF

run_mlr --inidx --ifs comma altkv <<EOF
a,b,c,d,e,f
EOF

run_mlr --inidx --ifs comma altkv <<EOF
a,b,c,d,e,f,g
EOF
