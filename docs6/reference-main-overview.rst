..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: Miller commands
===============================

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
