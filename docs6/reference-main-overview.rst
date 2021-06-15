..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: Miller commands
===============================

Overview
----------------------------------------------------------------

TODO: push this into reference-verbs overview, and make this page an overview of miller invocations

Whereas the Unix toolkit is made of the separate executables ``cat``, ``tail``, ``cut``,
``sort``, etc., Miller has subcommands, or **verbs**, invoked as follows:

.. code-block:: none

    mlr tac *.dat
    mlr cut --complement -f os_version *.dat
    mlr sort -f hostname,uptime *.dat

These fall into categories as follows:

* Analogs of their Unix-toolkit namesakes, discussed below as well as in :doc:`feature-comparison`: :ref:`reference-verbs-cat`, :ref:`reference-verbs-cut`, :ref:`reference-verbs-grep`, :ref:`reference-verbs-head`, :ref:`reference-verbs-join`, :ref:`reference-verbs-sort`, :ref:`reference-verbs-tac`, :ref:`reference-verbs-tail`, :ref:`reference-verbs-top`, :ref:`reference-verbs-uniq`.

* ``awk``-like functionality: :ref:`reference-verbs-filter`, :ref:`reference-verbs-put`, :ref:`reference-verbs-sec2gmt`, :ref:`reference-verbs-sec2gmtdate`, :ref:`reference-verbs-step`, :ref:`reference-verbs-tee`.

* Statistically oriented: :ref:`reference-verbs-bar`, :ref:`reference-verbs-bootstrap`, :ref:`reference-verbs-decimate`, :ref:`reference-verbs-histogram`, :ref:`reference-verbs-least-frequent`, :ref:`reference-verbs-most-frequent`, :ref:`reference-verbs-sample`, :ref:`reference-verbs-shuffle`, :ref:`reference-verbs-stats1`, :ref:`reference-verbs-stats2`.

* Particularly oriented toward :doc:`record-heterogeneity`, although all Miller commands can handle heterogeneous records: :ref:`reference-verbs-group-by`, :ref:`reference-verbs-group-like`, :ref:`reference-verbs-having-fields`.

* These draw from other sources (see also :doc:`originality`): :ref:`reference-verbs-count-distinct` is SQL-ish, and :ref:`reference-verbs-rename` can be done by ``sed`` (which does it faster: see :doc:`performance`. Verbs: :ref:`reference-verbs-check`, :ref:`reference-verbs-count-distinct`, :ref:`reference-verbs-label`, :ref:`reference-verbs-merge-fields`, :ref:`reference-verbs-nest`, :ref:`reference-verbs-nothing`, :ref:`reference-verbs-regularize`, :ref:`reference-verbs-rename`, :ref:`reference-verbs-reorder`, :ref:`reference-verbs-reshape`, :ref:`reference-verbs-seqgen`.


Verbs vs DSL
----------------------------------------------------------------

When you type ``mlr {something} myfile.dat``, the ``{something}`` part is called a **verb**. It specifies how you want to transform your data. (See also :doc:`reference-main-overview` for a breakdown.) The following is an alphabetical list of verbs with their descriptions.

The verbs ``put`` and ``filter`` are special in that they have a rich expression language (domain-specific language, or "DSL"). More information about them can be found at :doc:`reference-dsl`.

Here's a comparison of verbs and ``put``/``filter`` DSL expressions:

Example:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr stats1 -a sum -f x -g a data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* Verbs are coded in Go
* They run a bit faster
* They take fewer keystrokes
* There is less to learn
* Their customization is limited to each verb's options

Example:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr  put -q '@x_sum[$a] += $x; end{emit @x_sum, "a"}' data/small
    a=pan,x_sum=0.3467901443380824
    a=eks,x_sum=1.1400793586611044
    a=wye,x_sum=0.7778922255683036

* You get to write your own DSL expressions
* They run a bit slower
* They take more keystrokes
* There is more to learn
* They are highly customizable
