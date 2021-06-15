..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

What's new in Miller 6
================================================================

[doc is WIP]

* Completely reworked documentation (here) and on-line help (``mlr --help``)
* Arrays in the ``put``/``filter`` DSL
* JSON:

  * Improved JSON support
  * Streamable JSON parsing

* Full* support for Windows

  * Make a Windows docpage

* Build artifacts (binaries) using GitHub Actions
* In-process support for compressed input
* Input-preservation -- find a way to describe this -- link to the issue ...
* REPL TBD
* :doc:`output-colorization`
* Minor:

  * Getoptish (#467)
  * ``--mfrom``, ``--load``, ``--mload``
  * Better syntax-error messages for the DSL, including line number
  * Completely reworked regression-testing

* Dev: ported to Go

  * Developer notes: https://github.com/johnkerl/miller/blob/main/go/README.md

See also https://github.com/johnkerl/miller/issues/372
