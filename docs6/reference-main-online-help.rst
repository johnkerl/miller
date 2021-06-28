..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: online help
================================================================

TODO: expand this section

Examples:

.. code-block:: none
   :emphasize-lines: 1-1

    mlr --help
    Please run "mlr --help" for detailed usage information.

.. code-block:: none
   :emphasize-lines: 1-1

    mlr sort --help
    Usage: mlr sort {flags}
    Sorts records primarily by the first specified field, secondarily by the second
    field, and so on.  (Any records not having all specified sort keys will appear
    at the end of the output, in the order they were encountered, regardless of the
    specified sort order.) The sort is stable: records that compare equal will sort
    in the order they were encountered in the input record stream.
    
    Options:
    -f  {comma-separated field names}  Lexical ascending
    -n  {comma-separated field names}  Numerical ascending; nulls sort last
    -nf {comma-separated field names}  Same as -n
    -r  {comma-separated field names}  Lexical descending
    -nr {comma-separated field names}  Numerical descending; nulls sort first
    -h|--help Show this message.
    
    Example:
      mlr sort -f a,b -nr x,y,z
    which is the same as:
      mlr sort -f a -f b -nr x -nr y -nr z
