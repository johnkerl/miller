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
