# Miller Sphinx docs

## Why use Sphinx

* Connects to readthedocs.com so people can get their docmods onto the web instead of the self-hosted https://johnkerl.org/miller/doc. Thanks to @pabloab for the great advice!
* More standard look and feel -- lots of people use readthedocs for other things so this should feel familiar
* We get a Search feature for free

## Contributing

* You need `pip install sphinx` (or `pip3 install sphinx`)
* The docs include lots of live code examples which will be invoked using `mlr` which must be somewhere in your `$PATH`
* Clone https://github.com/johnkerl/miller and cd into `docs/` within your clone
* Editing loop:
  * Edit `*.rst.in`
  * Run `make html`
  * Either `open _build/html/index.html` (MacOS) or point your browser to `file:///path/to/your/clone/of/miller/docs/_build/html/index.html`
* Submitting:
  * `git add` your modified files, `git commit`, `git push`, and submit a PR at https://github.com/johnkerl/miller
* A nice markup reference: https://www.sphinx-doc.org/en/1.8/usage/restructuredtext/basics.html

## Notes

* CSS:
  * I used the Sphinx Classic theme which I like a lot except the colors -- it's a blue scheme and Miller has never been blue.
  * Files are in `docs/_static/*.css` where I marked my mods with `/* CHANGE ME */`.
  * If you modify the CSS you must run `make clean html` (not just `make html`) then reload in your browser.
* Live code:
  * I didn't find a way to include non-Python live-code examples within Sphinx so I adapted the pre-Sphinx Miller-doc strategy which is to have a generator script read a template file (here, `foo.rst.in`), run the marked lines, and generate the output file (`foo.rst`).
  * Edit the `*.rst.in` files, not `*.rst` directly.
  * Within the `*.rst.in` files are lines like `POKI_RUN_COMMAND`. These will be run, and their output included, by `make html` which calls the `genrst` script for you.
* readthedocs:
  * https://readthedocs.org/
  * https://readthedocs.org/projects/miller/
  * https://readthedocs.org/projects/miller/builds/
  * https://miller.readthedocs.io/en/latest/

## To do

* separate install from build; latter to reference section
* unix-toolkit context: needs a leading paragraph
* Let's all discuss if/how we want the v2 docs to be structured better than the v1 docs.
* !! cross-references all need work !!
* Scan for hrefs and other non-ported markup
* Autogen the `manpage.txt`
* Get rid of `POKI_CARDIFY` -- just indent by 4
* Chocolatey to windows-install notes
