# How to run the profiler

Running Miller:

```
mlr --cpuprofile cpu.pprof put -f u/example.mlr then nothing ~/tmp/huge > /dev/null
```

(or whatever command-line flags).

# How to view the profiling results


## Text mode

```
go tool pprof mlr cpu.pprof
top10
```

## PDF mode

```
go tool pprof --pdf mlr cpu.pprof > mlr-call-graph.pdf
mv mlr-call-graph.pdf ~/Desktop
```

## Flame-graph mode

One-time setup:


```
export GOPATH=$HOME/go
mkdir -p $HOME/go
```

```
go get -u github.com/google/pprof
ll ~/go/bin/pprof
go get -u github.com/uber/go-torch
```

```
mkdir -p ~/git/brendangregg
cd ~/git/brendangregg
git clone https://github.com/brendangregg/FlameGraph
```

Per run:

```
cd /path/to/mlr/go
export PATH=${PATH}:~/git/brendangregg/FlameGraph/
go-torch cpu.pprof
mv torch.svg ~/Desktop/
```

# How to control garbage collection

```
# Note 100 is the default
# Raise the bar for GC threshold:
GOGC=200  GODEBUG=gctrace=1 mlr -n put -q -f u/mand.mlr 1> /dev/null

# Raise the bar higher for GC threshold:
GOGC=1000 GODEBUG=gctrace=1 mlr -n put -q -f u/mand.mlr 1> /dev/null

# Turn off GC entirely and see where time is spent:
GOGC=off  GODEBUG=gctrace=1 mlr -n put -q -f u/mand.mlr 1> /dev/null
```

# Findings from the profiler as of 2021-02

* GC on: lots of GC
* GC off: lots of `duffcopy` and `madvise`
* From benchmark `u/mand.mlr`: issue is allocation created by expressions
  * Things like `type BinaryFunc func(input1 *Mlrval, input2 *Mlrval) (output Mlrval)`
  * `z = 2 + x + 4 + y * 3` results in AST, mapped to a CST, with a malloc on the output of every unary/binary/ternary function
  * Idea: replace with `type BinaryFunc func(output* Mlrval, input1 *Mlrval, input2 *Mlrval)`: no allocation at all.
    * Breaks the Fibonacci test since the binary-operator node is no longer re-entrant
  * Idea: replace with `type BinaryFunc func(input1 *Mlrval, input2 *Mlrval) (output *Mlrval)`: better.
    * Makes possible zero-copy eval of literal nodes, etc.

```
for i in 100 200 300 400 500 600 700 800 900 1000 ; do
  for j in 1 2 3 4 5 ; do
    echo $i;
    justtime GOGC=$i mlr -n put -q -f u/mand.mlr > /dev/null
  done
done
```

```
 100 23.754
 100 23.883
 100 24.021
 100 24.022
 100 24.305
 200 20.864
 200 20.211
 200 19.980
 200 20.251
 200 20.691
 300 19.140
 300 18.610
 300 18.793
 300 19.111
 300 19.027
 400 18.067
 400 18.274
 400 18.344
 400 18.378
 400 18.250
 500 17.791
 500 17.644
 500 17.814
 500 18.064
 500 18.403
 600 17.878
 600 17.892
 600 18.034
 600 18.125
 600 18.008
 700 18.153
 700 18.286
 700 17.342
 700 21.136
 700 20.729
 800 19.585
 800 19.116
 800 17.170
 800 18.549
 800 18.236
 900 16.950
 900 17.883
 900 17.532
 900 17.551
 900 17.804
1000 20.076
1000 20.745
1000 19.657
1000 18.733
1000 18.560
```

Sweet spot around 500. Note https://golang.org/pkg/runtime/debug/#SetGCPercent.
