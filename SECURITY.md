# Security policy

## Reporting a vulnerability

Please report security vulnerabilities privately, through GitHub's private vulnerability reporting:

https://github.com/johnkerl/miller/security/advisories/new

Please don't open a public issue or pull request for a suspected vulnerability.

A useful report includes:

- the Miller version (`mlr --version`) and operating system
- a minimal command line and input that reproduce the problem
- what happens, and what you expected to happen

I'll acknowledge reports as soon as I can, work with you on a fix, and credit you in the published advisory unless you'd prefer not to be credited.

## Supported versions

Security fixes go into the latest Miller 6 release. Miller 5 (the C implementation) is no longer maintained.

## Scope

Miller can run external commands by design: the DSL `system` and `exec` functions, piped output redirects such as `tee | "command", $*`, and `--prepipe`/`--prepipex`. Running a Miller command line or DSL expression that you didn't write is equivalent to running a shell command you didn't write, so this isn't a vulnerability in itself.

To run Miller on command lines or DSL expressions from an untrusted source, use `--no-shell`, or set the `MLR_NO_SHELL` environment variable, which disables all of the above.

Things that are in scope include:

- data, filenames, or other input that causes Miller to run commands or access files the user didn't ask for
- settings in `.mlrrc` that lead to unexpected command execution
- ways to bypass `--no-shell` / `MLR_NO_SHELL`
- crashes or resource exhaustion triggered by input data
