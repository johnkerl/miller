..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: Miller commands
===============================

Overview
----------------------------------------------------------------

TODO: overview of miller invocations

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

* You get to write your own expressions
* They run a bit slower
* They take more keystrokes
* There's more to learn
* They're highly customizable
