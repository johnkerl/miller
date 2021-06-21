..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: auxiliary commands
================================================================

There are a few nearly-standalone programs which have nothing to do with the rest of Miller, do not participate in record streams, and do not deal with file formats. They might as well be little standalone executables but they're delivered within the main Miller executable for convenience.

.. code-block:: none
   :emphasize-lines: 1-1

    mlr aux-list
    Available subcommands:
      aux-list
      hex
      lecat
      termcvt
      unhex
      help
      regtest
      repl
    For more information, please invoke mlr {subcommand} --help.

.. code-block:: none
   :emphasize-lines: 1-1

    mlr lecat --help
    Usage: mlr lecat [options] {zero or more file names}
    Simple hex-dump.
    If zero file names are supplied, standard input is read.
    Options:
    -r: print only raw hex without leading offset indicators or trailing ASCII dump.
    -h or --help: print this message

.. code-block:: none
   :emphasize-lines: 1-1

    mlr termcvt --help
    Usage: mlr termcvt [option] {zero or more file names}
    Option (exactly one is required):
    --cr2crlf
    --lf2crlf
    --crlf2cr
    --crlf2lf
    --cr2lf
    --lf2cr
    -I in-place processing (default is to write to stdout)
    -h or --help: print this message
    Zero file names means read from standard input.
    Output is always to standard output; files are not written in-place.

.. code-block:: none
   :emphasize-lines: 1-1

    mlr hex --help
    Usage: mlr hex [options] {zero or more file names}
    Simple hex-dump.
    If zero file names are supplied, standard input is read.
    Options:
    -r: print only raw hex without leading offset indicators or trailing ASCII dump.
    -h or --help: print this message

.. code-block:: none
   :emphasize-lines: 1-1

    mlr unhex --help
    Usage: mlr unhex [options] {zero or more file names}
    Simple hex-dump.
    If zero file names are supplied, standard input is read.
    Options:
    -r: print only raw hex without leading offset indicators or trailing ASCII dump.
    -h or --help: print this message

Examples:

.. code-block:: none
   :emphasize-lines: 1-1

    echo 'Hello, world!' | mlr lecat --mono
    Hello, world![LF]

.. code-block:: none
   :emphasize-lines: 1-1

    echo 'Hello, world!' | mlr termcvt --lf2crlf | mlr lecat --mono
    Hello, world![CR][LF]

.. code-block:: none
   :emphasize-lines: 1-1

    mlr hex data/budget.csv
    00000000: 23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 |# Asana -- here |
    00000010: 61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 |are the budget f|
    00000020: 69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 |igures you asked|
    00000030: 20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 | for!.type,quant|
    00000040: 69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 |ity.purple,456.7|
    00000050: 38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f |8.green,678.12.o|
    00000060: 72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          |range,123.45.|

.. code-block:: none
   :emphasize-lines: 1-1

    mlr hex -r data/budget.csv
    23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 
    61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 
    69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 
    20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 
    69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 
    38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f 
    72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          

.. code-block:: none
   :emphasize-lines: 1-1

    mlr hex -r data/budget.csv | sed 's/20/2a/g' | mlr unhex
    #*Asana*--*here*are*the*budget*figures*you*asked*for!
    type,quantity
    purple,456.78
    green,678.12
    orange,123.45
