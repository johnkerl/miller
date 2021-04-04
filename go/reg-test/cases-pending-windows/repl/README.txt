The Miller REPL startup banner is only printed if the input is detected to be a
terminal -- i.e. yes for interactive use, no for scripted use (such as
regression tests).

However the "is a terminal" check is (as of April 2021 anyway) non-functional
on Windows so the startup banner is always printed.

This means that if you develop on Linux/MacOS any tests using `mlr repl` will
pass but will fail on Windows.

Be sure to use `mlr repl -q` so the Miller REPL startup banner will be
suppressed on all platforms.
