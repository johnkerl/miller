# ----------------------------------------------------------------
run_mlr repl <<EOF
EOF

# ----------------------------------------------------------------
run_mlr repl <<EOF
x=1
y=2
print x+y
EOF

# ----------------------------------------------------------------
run_mlr repl <<EOF
\$x=3
\$*
EOF

# ----------------------------------------------------------------
run_mlr repl <<EOF
<
begin {
  print "In the beginning:"
}
end {
  print "At the end:"
}
# Populates the main block
print "In ...";
print "... the middle!"
>

begin { print "HELLO" }
begin { print "WORLD" }
end   { print "GOODBYE" }
end   { print "WORLD" }

# Immediately executed
print "HOW ARE THINGS?"

:blocks

:begin
:main
:end
EOF

# ----------------------------------------------------------------
run_mlr repl <<EOF

:open $indir/medium.dkvp
:context

:skip 10
:context
\$*

:process 10
:context
\$*

:skip until NR == 30
:context
\$*

:process until NR == 40
:context
\$*
EOF

# ----------------------------------------------------------------
run_mlr repl --j2x $indir/flatten-input-2.json <<EOF
:rw
:rw
EOF

run_mlr repl --x2j $indir/unflatten-input.xtab <<EOF
:rw
:rw
EOF

run_mlr repl --xtab $indir/unflatten-input.xtab <<EOF
:rw
:rw
EOF

run_mlr repl --json $indir/flatten-input-2.json <<EOF
:rw
:rw
EOF
