# Miller docs

## Why Mkdocs

- Connects to [https://miller.readthedocs.io](https://miller.readthedocs.io) so people can get their
  docmods onto the web instead of the self-hosted
  [https://johnkerl.org/miller/doc](https://johnkerl.org/miller/doc). Thanks to @pabloab for the
  great advice!
- More standard look and feel -- lots of people use readthedocs for other things so this should feel
  familiar.
- We get a Search feature for free.
- Mkdocs vs Sphinx: these are similar tools, but I find that I more easily get better desktop+mobile
  formatting using Mkdocs.

## Setup

- `pip install mkdocs` (or `pip3 install mkdocs`) as well as `pip install mkdocs-material`.
- The docs include lots of live code examples which are invoked using `mlr`, which must be somewhere
  in your `$PATH`.
- Clone [https://github.com/johnkerl/miller](https://github.com/johnkerl/miller) and cd into `docs/`
  within your clone.

## How the build works

- `docs/src` has `*.md.in` files containing markdown as well as directives for auto-generating code
  samples. Edit these files, not `*.md` directly.
- Within a `*.md.in` file, lines like `GENMD_RUN_COMMAND` mark a live code sample: the command gets
  run, and its output included, when the file is processed.
- The `genmds` script reads `docs/src/*.md.in` and writes `docs/src/*.md`.
- `mkdocs build` reads `docs/src/*.md` and writes HTML files in `docs/site`.
- Running `make` within the `docs` directory handles both of those steps.
- TL;DR: just `make docs` from the Miller base directory.

## Everyday editing loop

- In one terminal, cd to the `docs` directory and leave `mkdocs serve` running.
- In another terminal, cd to the `docs/src` subdirectory and edit `*.md.in`.
- Run `genmds` to re-create all the `*.md` files, or `genmds foo.md.in` to just re-create the
  `foo.md.in` file you just edited, or (simplest) just `make` within the `docs/src` subdirectory.
- In your browser, visit [http://127.0.0.1:8000](http://127.0.0.1:8000)
- This doesn't write HTML in `docs/site`; HTML is served up directly in the browser -- this is nice
  for previewing interactive edits.

## Adding a new page

- Create `docs/src/foo.md.in`.
- Add a nav entry in `docs/mkdocs.yml` pointing at `foo.md` -- the filename must match the basename
  of your new `.md.in` exactly, or the built page won't be linked into the site.
- First-time gotcha: `docs/src/Makefile`'s `genmds` target only rebuilds files already listed in
  `$(wildcard *.md)`. Since `foo.md` doesn't exist yet, a plain `make` within `docs/src` won't
  generate it. Two ways around this:
  - Simplest: just use `make docs` from the Miller base directory (or `make -C docs/src forcebuild`)
    -- these force-run `genmds` on every `*.md.in` unconditionally, new or not.
  - Or, from `docs/src`, run `./genmds foo.md.in` directly, once, to create the initial `foo.md`.
    After that it's picked up by ordinary `make` runs like any other page.
- Don't work around this by pre-creating an empty `foo.md` yourself (e.g. via `touch`): if its mtime
  ends up newer than `foo.md.in`, make's implicit `%.md: %.md.in` rule will see the target as up to
  date and skip regenerating it, silently leaving it empty.
- `git add` both `foo.md.in` and the generated `foo.md`.

## Publishing build

- cd to the `src` subdirectory of `docs` and edit `*.md.in`.
- `make -C ..`
- This does write HTML in `docs/site`.
- In your browser, visit `file:///your/path/to/miller/docs/site/index.html`
- Link-checking:
  - `sudo pip3 install git+https://github.com/linkchecker/linkchecker.git`
  - `cd site` and `linkchecker .`

## Submitting

- Do the publishing-build steps -- in particular, `docs/src/*.md.in` and `docs/src/*.md` are both
  checked in to source control.
  - TL;DR: edit `docs/src/foo.md.in` and run `make docs`.
  - If you don't want to do `pip install mkdocs` then feel free to put up a PR which edits a
    `foo.md.in` as well as its `foo.md`.
- `git add` your modified files (`*.md.in` as well as `*.md`), `git commit`, `git push`, and submit
  a PR at [https://github.com/johnkerl/miller](https://github.com/johnkerl/miller).

## Style notes

- Miller documents use the Oxford comma: not _red, yellow and green_, but rather _red, yellow, and
  green_.
- CSS: the Mkdocs "material" theme is used, customized via `docs/src/extra.css` for Miller
  coloring/branding.

## readthedocs

- Published to [https://miller.readthedocs.io/en/latest](https://miller.readthedocs.io/en/latest) on each commit to `main` in this repo.
- [https://readthedocs.org/](https://readthedocs.org/)
- [https://readthedocs.org/projects/miller/](https://readthedocs.org/projects/miller/)
- [https://readthedocs.org/projects/miller/builds/](https://readthedocs.org/projects/miller/builds/)
- [https://readthedocs.org/api/v2/webhook/miller/134065/](https://readthedocs.org/api/v2/webhook/miller/134065/)
