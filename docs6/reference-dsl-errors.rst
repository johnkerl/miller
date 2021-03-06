..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference: errors and transparency
======================================

As soon as you have a programming language, you start having the problem *What is my code doing, and why?* This includes getting syntax errors -- which are always annoying -- as well as the even more annoying problem of a program which parses without syntax error but doesn't do what you expect.

The ``syntax error`` message is cryptic: it says ``syntax error at `` followed by the next symbol it couldn't parse. This is good, but (as of 5.0.0) it doesn't say things like ``syntax error at line 17, character 22``. Here are some common causes of syntax errors:

* Don't forget ``;`` at end of line, before another statement on the next line.

* Miller's DSL lacks the ``++`` and ``--`` operators.

* Curly braces are required for the bodies of ``if``/``while``/``for`` blocks, even when the body is a single statement.

Now for transparency:

* As in any language, you can do (see :ref:`reference-dsl-print-statements`) ``print`` (or ``eprint`` to print to stderr). See also :ref:`reference-dsl-dump-statements` and :ref:`reference-dsl-emit-statements`.

* The ``-v`` option to ``mlr put`` and ``mlr filter`` prints abstract syntax trees for your code. While not all details here will be of interest to everyone, certainly this makes questions such as operator precedence completely unambiguous.

* The ``-T`` option prints a trace of each statement executed.

* The ``-t`` and ``-a`` options show low-level details for the parsing process and for stack-variable-index allocation, respectively. These will likely be of interest to people who enjoy compilers, and probably less useful for a more general audience.

* Please see :ref:`reference-dsl-type-checking` for type declarations and type-assertions you can use to make sure expressions and the data flowing them are evaluating as you expect.  I made them optional because one of Miller's important use-cases is being able to say simple things like ``mlr put '$y = $x + 1' myfile.dat`` with a minimum of punctuational bric-a-brac -- but for programs over a few lines I generally find that the more type-specification, the better.
