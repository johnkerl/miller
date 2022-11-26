<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Programming-language examples

Here are a few things focusing on Miller's DSL as a programming language per se, outside of its normal use for streaming record-processing.

## Sieve of Eratosthenes

The [Sieve of Eratosthenes](http://en.wikipedia.org/wiki/Sieve_of_Eratosthenes) is a standard introductory programming topic. The idea is to find all primes up to some *N* by making a list of the numbers 1 to *N*, then striking out all multiples of 2 except 2 itself, all multiples of 3 except 3 itself, all multiples of 4 except 4 itself, and so on. Whatever survives that without getting marked is a prime. This is easy enough in Miller. Notice that here all the work is in `begin` and `end` statements; there is no file input (so we use `mlr -n` to keep Miller from waiting for input data).

<pre class="pre-highlight-in-pair">
<b>cat programs/sieve.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# ================================================================
# Sieve of Eratosthenes: simple example of Miller DSL as programming language.
# ================================================================

# Put this in a begin-block so we can do either
#   mlr -n put -q -f name-of-this-file.mlr
# or
#   mlr -n put -q -f name-of-this-file.mlr -e '@n = 200'
# i.e. 100 is the default upper limit, and another can be specified using -e.
begin {
  @n = 100;
}

end {
  for (int i = 0; i <= @n; i += 1) {
    @s[i] = true;
  }
  @s[0] = false; # 0 is neither prime nor composite
  @s[1] = false; # 1 is neither prime nor composite
  # Strike out multiples
  for (int i = 2; i <= @n; i += 1) {
    for (int j = i+i; j <= @n; j += i) {
      @s[j] = false;
    }
  }
  # Print survivors
  for (int i = 0; i <= @n; i += 1) {
    if (@s[i]) {
      print i;
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put -f programs/sieve.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
2
3
5
7
11
13
17
19
23
29
31
37
41
43
47
53
59
61
67
71
73
79
83
89
97
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## Mandelbrot-set generator

The [Mandelbrot set](http://en.wikipedia.org/wiki/Mandelbrot_set) is also easily expressed. This isn't an important case of data processing (the use-case Miller was designed for), but it is an example of Miller as a general-purpose programming language -- a test case for the expressiveness of the language.

The (approximate) computation of points in the complex plane which are and aren't members is just a few lines of complex arithmetic (see the [Wikipedia article](https://en.wikipedia.org/wiki/Mandelbrot_set)); how to render them visually is another task.  Using graphics libraries you can create PNG or JPEG files, but another fun way to do this is by printing various characters to the screen:

<pre class="pre-highlight-in-pair">
<b>cat programs/mand.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# Mandelbrot set generator: simple example of Miller DSL as programming language.
begin {
  # Set defaults. They can be overridden by e.g.
  #   mlr -n put -e 'begin{@maxits=200}' -f nameofthisfile.mlr
  # or
  #   mlr -n put -s maxits=200 -f nameofthisfile.mlr
  @rcorn     ??= -2.0;
  @icorn     ??= -2.0;
  @side      ??=  4.0;
  @iheight   ??=   50;
  @iwidth    ??=  100;
  @maxits    ??=  100;
  @levelstep ??=    5;
  @chars     ??= "@X*o-.";
  @silent    ??= false;
  @do_julia  ??= false;
  @jr        ??= 0.0;      # Real part of Julia point, if any
  @ji        ??= 0.0;      # Imaginary part of Julia point, if any
}

end {
  if (!@silent) {
    print "RCORN     = ".@rcorn;
    print "ICORN     = ".@icorn;
    print "SIDE      = ".@side;
    print "IHEIGHT   = ".@iheight;
    print "IWIDTH    = ".@iwidth;
    print "MAXITS    = ".@maxits;
    print "LEVELSTEP = ".@levelstep;
    print "CHARS     = ".@chars;
  }

  for (int ii = @iheight-1; ii >= 0; ii -= 1) {
    num ci = @icorn + (ii/@iheight) * @side;
    for (int ir = 0; ir < @iwidth; ir += 1) {
      num cr = @rcorn + (ir/@iwidth) * @side;
      str c = get_point_plot(cr, ci, @maxits, @do_julia, @jr, @ji);
      if (!@silent) {
        printn c
      }
    }
    if (!@silent) {
      print;
    }
  }
}

func get_point_plot(num pr, num pi, int maxits, bool do_julia, num jr, num ji): str {
  num zr = 0.0;
  num zi = 0.0;
  num cr = 0.0;
  num ci = 0.0;

  if (!do_julia) {
    zr = 0.0;
    zi = 0.0;
    cr = pr;
    ci = pi;
  } else {
    zr = pr;
    zi = pi;
    cr = jr;
    ci = ji;
  }

  int iti = 0;
  bool escaped = false;
  num zt = 0;
  for (iti = 0; iti < maxits; iti += 1) {
    num mag = zr*zr + zi+zi;
    if (mag > 4.0) {
        escaped = true;
        break;
    }
    # z := z^2 + c
    zt = zr*zr - zi*zi + cr;
    zi = 2*zr*zi + ci;
    zr = zt;
  }
  if (!escaped) {
    return ".";
  } else {
    int level = (iti // @levelstep) % strlen(@chars);
    return substr(@chars, level, level);
  }
}
</pre>

At standard resolution this makes a nice little ASCII plot:

<pre class="pre-highlight-in-pair">
<b>mlr -n put -s iheight=25 -s iwidth=50 -f ./programs/mand.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
RCORN     = -2
ICORN     = -2
SIDE      = 4.0
IHEIGHT   = 25
IWIDTH    = 50
MAXITS    = 100
LEVELSTEP = 5
CHARS     = @X*o-.
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@XX.XX@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@XX*o.XXX@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@XX@XXXXX...oXXXXX@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@XXXXXX..........X.-X@@@@@@@@@@@@@@@@@@
@@@@@@@@XXXXXXXX*o............XX@@@@@@@@@@@@@@@@@@
@@@@@@XXXX-...-*...............X@@@@@@@@@@@@@@@@@@
@XX@XXXXoo....................XX@@@@@@@@@@@@@@@@@@
@@@XXXXX-o....................XXX@@@@@@@@@@@@@@@@@
@@@@XXXXXX-....*...............XXXXX@@@@@@@@@@@@@@
@@@@@@@XXXXX*XXX*-............*XXX@@@@@@@@@@@@@@@@
@@@@@@@@@@@@XXXXXX..........*.-X@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@XXXX*@..*XXXX@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@XXXXXoo.XXXX@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@XXXXXX@XXXX@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

But using a very small font size (as small as my Mac will let me go), and by choosing the coordinates to zoom in on a particular part of the complex plane, we can get a nice little picture:

<pre class="pre-non-highlight-non-pair">
#!/bin/bash
# Get the number of rows and columns from the terminal window dimensions
iheight=$(stty size | mlr --nidx --fs space cut -f 1)
iwidth=$(stty size | mlr --nidx --fs space cut -f 2)
mlr -n put \
  -s rcorn=-1.755350 -s icorn=0.014230 -s side=0.000020 -s maxits=10000 -s iheight=$iheight -s iwidth=$iwidth \
  -f programs/mand.mlr
</pre>

![pix/mand.png](pix/mand.png)
