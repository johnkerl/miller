..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Then-chaining
=============

How do I examine then-chaining?
----------------------------------------------------------------

Then-chaining found in Miller is intended to function the same as Unix pipes, but with less keystroking. You can print your data one pipeline step at a time, to see what intermediate output at one step becomes the input to the next step.

First, look at the input data:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/then-example.csv
    Status,Payment_Type,Amount
    paid,cash,10.00
    pending,debit,20.00
    paid,cash,50.00
    pending,credit,40.00
    paid,debit,30.00

Next, run the first step of your command, omitting anything from the first ``then`` onward:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --icsv --opprint count-distinct -f Status,Payment_Type data/then-example.csv
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

After that, run it with the next ``then`` step included:

.. code-block:: none
   :emphasize-lines: 1-3

    $ mlr --icsv --opprint count-distinct -f Status,Payment_Type \
      then sort -nr count \
      data/then-example.csv
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

Now if you use ``then`` to include another verb after that, the columns ``Status``, ``Payment_Type``, and ``count`` will be the input to that verb.

Note, by the way, that you'll get the same results using pipes:

.. code-block:: none
   :emphasize-lines: 1-2

    $ mlr --csv count-distinct -f Status,Payment_Type data/then-example.csv \
    | mlr --icsv --opprint sort -nr count
    Status  Payment_Type count
    paid    cash         2
    pending debit        1
    pending credit       1
    paid    debit        1

NR is not consecutive after then-chaining
----------------------------------------------------------------

Given this input data:

.. code-block:: none
   :emphasize-lines: 1-1

    $ cat data/small
    a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
    a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729

why don't I see ``NR=1`` and ``NR=2`` here??

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$x > 0.5' then put '$NR = NR' data/small
    a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,NR=2
    a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,NR=5

The reason is that ``NR`` is computed for the original input records and isn't dynamically updated. By contrast, ``NF`` is dynamically updated: it's the number of fields in the current record, and if you add/remove a field, the value of ``NF`` will change:

.. code-block:: none
   :emphasize-lines: 1-1

    $ echo x=1,y=2,z=3 | mlr put '$nf1 = NF; $u = 4; $nf2 = NF; unset $x,$y,$z; $nf3 = NF'
    nf1=3,u=4,nf2=5,nf3=3

``NR``, by contrast (and ``FNR`` as well), retains the value from the original input stream, and records may be dropped by a ``filter`` within a ``then``-chain. To recover consecutive record numbers, you can use out-of-stream variables as follows:

.. code-block:: none
   :emphasize-lines: 1-11

    $ mlr --opprint --from data/small put '
      begin{ @nr1 = 0 }
      @nr1 += 1;
      $nr1 = @nr1
    ' \
    then filter '$x>0.5' \
    then put '
      begin{ @nr2 = 0 }
      @nr2 += 1;
      $nr2 = @nr2
    '
    a   b   i x                  y                  nr1 nr2
    eks pan 2 0.7586799647899636 0.5221511083334797 2   1
    wye pan 5 0.5732889198020006 0.8636244699032729 5   2

Or, simply use ``mlr cat -n``:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr filter '$x > 0.5' then cat -n data/small
    n=1,a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
    n=2,a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
