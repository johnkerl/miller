..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference: overview
================================================================

Overview
----------------------------------------------------------------

Here's comparison of verbs and ``put``/``filter`` DSL expressions:

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

Please see :doc:`reference-verbs` for information on verbs other than ``put`` and ``filter``.

The essential usages of ``mlr filter`` and ``mlr put`` are for record-selection and record-updating expressions, respectively. For example, given the following input data:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

you might retain only the records whose ``a`` field has value ``eks``:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$a == "eks"' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463

or you might add a new field which is a function of existing fields:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr put '$ab = $a . "_" . $b ' data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,ab=pan_pan
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,ab=eks_pan
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,ab=wye_wye
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,ab=eks_wye
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,ab=wye_pan

The two verbs ``mlr filter`` and ``mlr put`` are essentially the same. The only differences are:

* Expressions sent to ``mlr filter`` must end with a boolean expression, which is the filtering criterion;

* ``mlr filter`` expressions may not reference the ``filter`` keyword within them; and

* ``mlr filter`` expressions may not use ``tee``, ``emit``, ``emitp``, or ``emitf``.

All the rest is the same: in particular, you can define and invoke functions and subroutines to help produce the final boolean statement, and record fields may be assigned to in the statements preceding the final boolean statement.

There are more details and more choices, of course, as detailed in the following sections.

