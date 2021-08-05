# Miller docs

## Why use Mkdocs

* Connects to https://miller.readthedocs.io so people can get their docmods onto the web instead of the self-hosted https://johnkerl.org/miller/doc. Thanks to @pabloab for the great advice!
* More standard look and feel -- lots of people use readthedocs for other things so this should feel familiar.
* We get a Search feature for free.
* Mkdocs vs Sphinx: these are similar tools, but I find that I more easily get better desktop+mobile formatting using Mkdocs.

## Contributing

* You need `pip install mkdocs` (or `pip3 install mkdocs`).
* The docs include lots of live code examples which will be invoked using `mlr` which must be somewhere in your `$PATH`.
* Clone https://github.com/johnkerl/miller and cd into `docs/` within your clone.
* Quick-editing loop:
  * In one terminal, cd to this directory and leave `mkdocs serve` running.
  * In another terminal, cd to the `docs` subdirectory and edit `*.md.in`.
  * Run `genmds` to re-create all the `*.md` files, or `genmds foo.md.in` to just re-create the `foo.md.in` file you just edited.
  * In your browser, visit http://127.0.0.1:8000
* Alternate editing loop:
  * Leave one terminal open as a place you will run `mkdocs build`
  * In one terminal, cd to the `docs` subdirectory and edit `*.md.in`.
  * Run `genmds` to re-create all the `*.md` files, or `genmds foo.md.in` to just re-create the `foo.md.in` file you just edited.
  * In the first terminal, run `mkdocs build` which will populate the `site` directory.
  * In your browser, visit `file:///your/path/to/miller/docs/site/index.html`
  * Link-checking:
    * `sudo pip3 install git+https://github.com/linkchecker/linkchecker.git`
    * `cd site` and `linkchecker .`
* Submitting:
  * `git add` your modified files, `git commit`, `git push`, and submit a PR at https://github.com/johnkerl/miller.

## Notes

* CSS:
  * I used the Mkdocs Readthedocs theme which I like a lot. I customized `docs/extra.css` for Miller coloring/branding.
* Live code:
  * I didn't find a way to include non-Python live-code examples within Mkdocs so I adapted the pre-Mkdocs Miller-doc strategy which is to have a generator script read a template file (here, `foo.md.in`), run the marked lines, and generate the output file (`foo.md`). This is `genmds`.
  * Edit the `*.md.in` files, not `*.md` directly.
  * Within the `*.md.in` files are lines like `GENMD_RUN_COMMAND`. These will be run, and their output included, by `genmds` which calls the `genmds` script for you.
* readthedocs:
  * https://readthedocs.org/
  * https://readthedocs.org/projects/miller/
  * https://readthedocs.org/projects/miller/builds/
  * https://miller.readthedocs.io/en/latest/
