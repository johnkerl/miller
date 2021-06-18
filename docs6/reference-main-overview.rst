..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: Miller commands
===============================

Overview
----------------------------------------------------------------

The outline of an invocation of Miller is

* ``mlr``
* Options controlling input/output formatting, etc. (:doc:`reference-main-io-options`).
* One or more verbs (such as ``cut``, ``sort``, etc.) (:doc:`reference-verbs`) -- chained together using ``then`` (:doc:`reference-main-then-chaining`). You use these to transform your data.
* Zero or more filenames, with input taken from standard input if there are no filenames present.

For example, reading from a file:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint head -n 2 then sort -f shape example.csv
    color  shape    flag index quantity rate
    red    square   true 15    79.2778  0.0130
    yellow triangle true 11    43.6498  9.8870

Reading from standard input:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat example.csv | mlr --icsv --opprint head -n 2 then sort -f shape
    color  shape    flag index quantity rate
    red    square   true 15    79.2778  0.0130
    yellow triangle true 11    43.6498  9.8870

The rest of this reference section gives you full information on each of these parts of the command line.

Verbs vs DSL
----------------------------------------------------------------

When you type ``mlr {something} myfile.dat``, the ``{something}`` part is called a **verb**. It specifies how you want to transform your data. Most of the verbs are counterparts of built-in system tools like ``cut`` and ``sort`` -- but with file-format awareness, and giving you the ability to refer to fields by name.

The verbs ``put`` and ``filter`` are special in that they have a rich expression language (domain-specific language, or "DSL"). More information about them can be found at :doc:`reference-dsl`.

Here's a comparison of verbs and ``put``/``filter`` DSL expressions:

Example of using a verb for data processing:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr stats1 -a sum -f x -g a data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* Verbs are coded in Go
* They run a bit faster
* They take fewer keystrokes
* There's less to learn
* Their customization is limited to each verb's options

Example of doing the same thing using a DSL expression:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr  put -q '@x_sum[$a] += $x; end{emit @x_sum, "a"}' data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* You get to write your own expressions in Miller's programming language
* They run a bit slower
* They take more keystrokes
* There's more to learn
* They're highly customizable
