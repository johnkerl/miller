<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# How to contribute

## Community

You can ask questions -- or answer them! -- following the links at :doc:`community`.

## Documentation improvements

Pre-release Miller documentation is at https://github.com/johnkerl/miller/tree/main/docs6.

Clone https://github.com/johnkerl/miller and `cd` into `docs6`.

After ``sudo pip install sphinx`` (or ``pip3``) you should be able to do ``make html``.

Edit ``*.md.in`` files, then ``make html`` to generate ``*.md``, then run the Sphinx document-generator.

Open ``_build/html/index.html`` in your browser, e.g. ``file:////Users/yourname/git/miller/docs6/_build/html/contributing.html``, to verify.

PRs are welcome at https://github.com/johnkerl/miller.

Once PRs are merged, readthedocs creates https://miller.readthedocs.io using the following configs:

* https://readthedocs.org/projects/miller/
* https://readthedocs.org/projects/miller/builds/
* https://github.com/johnkerl/miller/settings/hooks

## Testing

As of Miller-6's current pre-release status, the best way to test is to either build from source via :doc:`build`, or by getting a recent binary at https://github.com/johnkerl/miller/actions, then click latest build, then *Artifacts*. Then simply use Miller for whatever you do, and create an issue at https://github.com/johnkerl/miller/issues.

Do note that as of 2021-06-17 a few things have not been ported to Miller 6 -- most notably, including regex captures and localtime DSL functions.

## Feature development

Issues: https://github.com/johnkerl/miller/issues

Developer notes: https://github.com/johnkerl/miller/blob/main/go/README.md

PRs which pass regression test (https://github.com/johnkerl/miller/blob/main/go/regtest/README.md) are always welcome!
