..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Running shell commands
======================

TODO: while-read example from issues

The :ref:`reference-dsl-system` DSL function allows you to run a specific shell command and put its output -- minus the final newline -- into a record field. The command itself is any string, either a literal string, or a concatenation of strings, perhaps including other field values or what have you.

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --opprint put '$o = system("echo hello world")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  hello world
    eks pan 2 0.7586799647899636  0.5221511083334797  hello world
    wye wye 3 0.20460330576630303 0.33831852551664776 hello world
    eks wye 4 0.38139939387114097 0.13418874328430463 hello world
    wye pan 5 0.5732889198020006  0.8636244699032729  hello world

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --opprint put '$o = system("echo {" . NR . "}")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  {1}
    eks pan 2 0.7586799647899636  0.5221511083334797  {2}
    wye wye 3 0.20460330576630303 0.33831852551664776 {3}
    eks wye 4 0.38139939387114097 0.13418874328430463 {4}
    wye pan 5 0.5732889198020006  0.8636244699032729  {5}

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --opprint put '$o = system("echo -n ".$a."| sha1sum")' data/small
    a   b   i x                   y                   o
    pan pan 1 0.3467901443380824  0.7268028627434533  f29c748220331c273ef16d5115f6ecd799947f13  -
    eks pan 2 0.7586799647899636  0.5221511083334797  456d988ecb3bf1b75f057fc6e9fe70db464e9388  -
    wye wye 3 0.20460330576630303 0.33831852551664776 eab0de043d67f441c7fd1e335f0ca38708e6ebf7  -
    eks wye 4 0.38139939387114097 0.13418874328430463 456d988ecb3bf1b75f057fc6e9fe70db464e9388  -
    wye pan 5 0.5732889198020006  0.8636244699032729  eab0de043d67f441c7fd1e335f0ca38708e6ebf7  -

Note that running a subprocess on every record takes a non-trivial amount of time. Comparing asking the system ``date`` command for the current time in nanoseconds versus computing it in process:

..
    hard-coded, not live-code, since %N doesn't exist on all platforms

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$t=system("date +%s.%N")' then step -a delta -f t data/small
    a   b   i x                   y                   t                    t_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.513903817 0
    eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.514722876 0.000819
    wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.515618046 0.000895
    eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.516547441 0.000929
    wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.517518828 0.000971

.. code-block:: none
   :emphasize-lines: 1,1

    $ mlr --opprint put '$t=systime()' then step -a delta -f t data/small
    a   b   i x                   y                   t                 t_delta
    pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.518699 0
    eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.518717 0.000018
    wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.518723 0.000006
    eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.518727 0.000004
    wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.518730 0.000003
