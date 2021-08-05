..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Internationalization
================================================================

Miller handles ASCII and UTF-8 strings. (I have no plans to support UTF-16 or ISO-8859-1.)

Support for internationalization includes:

* Tabular output formats such pprint and xtab (see :doc:`file-formats`) are aligned correctly.
* The :ref:`reference-dsl-strlen` function correctly counts UTF-8 codepoints rather than bytes.
* The :ref:`reference-dsl-toupper`, :ref:`reference-dsl-tolower`, and :ref:`reference-dsl-capitalize` DSL functions operate within the capabilities of the Go libraries.

Please file an issue at https://github.com/johnkerl/miller if you encounter bugs related to internationalization (or anything else for that matter).
