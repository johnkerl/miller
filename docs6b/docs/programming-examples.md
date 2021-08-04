<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Programming-language examples

Here are a few things focusing on Miller's DSL as a programming language per se, outside of its normal use for streaming record-processing.

## Sieve of Eratosthenes

The `Sieve of Eratosthenes <http://en.wikipedia.org/wiki/Sieve_of_Eratosthenes>`_ is a standard introductory programming topic. The idea is to find all primes up to some *N* by making a list of the numbers 1 to *N*, then striking out all multiples of 2 except 2 itself, all multiples of 3 except 3 itself, all multiples of 4 except 4 itself, and so on. Whatever survives that without getting marked is a prime. This is easy enough in Miller. Notice that here all the work is in ``begin`` and ``end`` statements; there is no file input (so we use ``mlr -n`` to keep Miller from waiting for input data).

<pre>
<b>cat programs/sieve.mlr</b>
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

<pre>
<b>mlr -n put -f programs/sieve.mlr</b>
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
</pre>

## Mandelbrot-set generator

The `Mandelbrot set <http://en.wikipedia.org/wiki/Mandelbrot_set>`_ is also easily expressed. This isn't an important case of data-processing in the vein for which Miller was designed, but it is an example of Miller as a general-purpose programming language -- a test case for the expressiveness of the language.

The (approximate) computation of points in the complex plane which are and aren't members is just a few lines of complex arithmetic (see the Wikipedia article); how to render them is another task.  Using graphics libraries you can create PNG or JPEG files, but another fun way to do this is by printing various characters to the screen:

<pre>
<b>cat programs/mand.mlr</b>
# Mandelbrot set generator: simple example of Miller DSL as programming language.
begin {
  # Set defaults
  @rcorn     = -2.0;
  @icorn     = -2.0;
  @side      = 4.0;
  @iheight   = 50;
  @iwidth    = 100;
  @maxits    = 100;
  @levelstep = 5;
  @chars     = "@X*o-."; # Palette of characters to print to the screen.
  @verbose   = false;
  @do_julia  = false;
  @jr        = 0.0;      # Real part of Julia point, if any
  @ji        = 0.0;      # Imaginary part of Julia point, if any
}

# Here, we can override defaults from an input file (if any).  In Miller's
# put/filter DSL, absent-null right-hand sides result in no assignment so we
# can simply put @rcorn = $rcorn: if there is a field in the input like
# 'rcorn = -1.847' we'll read and use it, else we'll keep the default.
@rcorn     = $rcorn;
@icorn     = $icorn;
@side      = $side;
@iheight   = $iheight;
@iwidth    = $iwidth;
@maxits    = $maxits;
@levelstep = $levelstep;
@chars     = $chars;
@verbose   = $verbose;
@do_julia  = $do_julia;
@jr        = $jr;
@ji        = $ji;

end {
  if (@verbose) {
    print "RCORN     = ".@rcorn;
    print "ICORN     = ".@icorn;
    print "SIDE      = ".@side;
    print "IHEIGHT   = ".@iheight;
    print "IWIDTH    = ".@iwidth;
    print "MAXITS    = ".@maxits;
    print "LEVELSTEP = ".@levelstep;
    print "CHARS     = ".@chars;
  }

  # Iterate over a matrix of rows and columns, printing one character for each cell.
  for (int ii = @iheight-1; ii >= 0; ii -= 1) {
    num pi = @icorn + (ii/@iheight) * @side;
    for (int ir = 0; ir < @iwidth; ir += 1) {
      num pr = @rcorn + (ir/@iwidth) * @side;
      printn get_point_plot(pr, pi, @maxits, @do_julia, @jr, @ji);
    }
    print;
  }
}

# This is a function to approximate membership in the Mandelbrot set (or Julia
# set for a given Julia point if do_julia == true) for a given point in the
# complex plane.
func get_point_plot(pr, pi, maxits, do_julia, jr, ji) {
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
    # The // operator is Miller's (pythonic) integer-division operator
    int level = (iti // @levelstep) % strlen(@chars);
    return substr(@chars, level, level);
  }
}
</pre>

At standard resolution this makes a nice little ASCII plot:

<pre>
<b>mlr -n put -f ./programs/mand.mlr</b>
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXX.XXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXooXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX**o..*XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXX*-....-oXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX@XXXXXXXXXX*......o*XXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXX**oo*-.-........oo.XXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX....................X..o-XXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXXXX*oo......................oXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@XXX*XXXXXXXXXXXX**o........................*X*X@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@XXXXXXooo***o*.*XX**X..........................o-XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@XXXXXXXX*-.......-***.............................oXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@XXXXXXXX*@..........Xo............................*XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@XXXX@XXXXXXXX*o@oX...........@...........................oXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
.........................................................o*XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@XXXXXXXXX*-.oX...........@...........................oXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@XXXXXXXXXX**@..........*o............................*XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@XXXXXXXXXXXXX-........***.............................oXXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@XXXXXXXXXXXXoo****o*.XX***@..........................o-XXXXXXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@XXXXX*XXXX*XXXXXXX**-........................***XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX*o*.....................@o*XXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXX*....................*..o-XX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX*ooo*-.o........oo.X*XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX**@.....*XXXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX*o....-o*XXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXo*o..*XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXXX*o*XXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXXXXX@XXXXXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXXXXXX@@XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@XXXXX@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
</pre>

But using a very small font size (as small as my Mac will let me go), and by choosing the coordinates to zoom in on a particular part of the complex plane, we can get a nice little picture:

<pre>
#!/bin/bash
# Get the number of rows and columns from the terminal window dimensions
iheight=$(stty size | mlr --nidx --fs space cut -f 1)
iwidth=$(stty size | mlr --nidx --fs space cut -f 2)
echo "rcorn=-1.755350,icorn=+0.014230,side=0.000020,maxits=10000,iheight=$iheight,iwidth=$iwidth" \
| mlr put -f programs/mand.mlr
</pre>

.. image:: pix/mand.png
