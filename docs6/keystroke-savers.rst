..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Keystroke-savers
================

Short format specifiers
^^^^^^^^^^^^^^^^^^^^^^^

In our examples so far we've often made use of ``mlr --icsv --opprint`` or ``mlr --icsv --ojson``. These are such frequently occurring patterns that they have short options like **--c2p** and **--c2j**:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --c2p head -n 2 example.csv
    color  shape    flag index quantity rate
    yellow triangle true 11    43.6498  9.8870
    red    square   true 15    79.2778  0.0130

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --c2j head -n 2 example.csv
    {
      "color": "yellow",
      "shape": "triangle",
      "flag": true,
      "index": 11,
      "quantity": 43.6498,
      "rate": 9.8870
    }
    {
      "color": "red",
      "shape": "square",
      "flag": true,
      "index": 15,
      "quantity": 79.2778,
      "rate": 0.0130
    }

You can get the full list here (TODO:linkify).

File names up front
^^^^^^^^^^^^^^^^^^^

Already we saw that you can put the filename first using ``--from``. When you're interacting with your data at the command line, this makes it easier to up-arrow and append to the previous command:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --c2p --from example.csv sort -nr index then head -n 3
    color  shape  flag  index quantity rate
    purple square false 91    72.3735  8.2430
    yellow circle true  87    63.5058  8.3350
    yellow circle true  73    63.9785  4.2370

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --c2p --from example.csv sort -nr index then head -n 3 then cut -f shape,quantity
    shape  quantity
    square 72.3735
    circle 63.5058
    circle 63.9785

If there's more than one input file, you can use ``--mfrom``, then however many file names, then ``--`` to indicate the end of your input-file-name list:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --c2p --mfrom data/*.csv -- sort -n index

.mlrrc file
^^^^^^^^^^^

If you want the default file format for Miller to be CSV you can simply put ``--csv`` on a line by itself in your ``~/.mlrrc`` file. Then instead of ``mlr --csv cat example.csv`` you can just do ``mlr cat example.csv``. This is just the default, though, so ``mlr --opprint cat example.csv`` will still use default CSV format for input, and PPRINT (tabular) for output.

You can read more about this at the :doc:`customization` page.
