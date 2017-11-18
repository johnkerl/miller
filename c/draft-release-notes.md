Features:

* comment-skipper
* `count-similar` verb
* `.*` int math
* `popcount` function

```seq 1 32 | mlr --nidx put '$2=bitcount(NR)' then sort -n 2```

Documentation:

* ruby/python/etc. dkvp-reader/writers, and example code
* 'How do I suppress numeric conversion?' @ mlh w/ pagelink; xref to followup task (#150 ...)
* csvlite at mlh/man

Bugfixes:

* fix https://github.com/johnkerl/miller/issues/158: short option for `--nidx --fs tab`
* Fix https://github.com/johnkerl/miller/issues/159: regex-match of literal dot

