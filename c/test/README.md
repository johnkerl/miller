There are two classes of testing for Miller:

* C source-file names starting with `test_` use MinUnit to **unit-test** various subsystems of interest.  These are separate executables built and run by the build framework.

* `test/run` runs the main `mlr` executable with canned inputs, comparing actual to canned outputs, to **regression-test** Miller's end-to-end operation.
