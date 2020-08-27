JSON object parsing with preserving keys order
=============================================

![coverage](https://gitlab.com/c0b/go-ordered-json/badges/master/coverage.svg)

Refers

1. JSON and Go        https://blog.golang.org/json-and-go
2. Go-Ordered-JSON    https://github.com/virtuald/go-ordered-json
   from this thread [*Preserving key order in encoding/json*](https://groups.google.com/forum/#!topic/golang-dev/zBQwhm3VfvU)
   and the [*Abandoned 7930: encoding/json: Optionally preserve the key order of JSON objects*](https://go-review.googlesource.com/c/go/+/7930)
3. Python OrderedDict https://github.com/python/cpython/blob/2.7/Lib/collections.py#L38
   the Python's OrderedDict uses a double linked list internally, maintain a consistent public interface with `dict`

Disclaimer:

same as Go's default [map](https://blog.golang.org/go-maps-in-action), this OrderedMap is not safe for concurrent use, if need atomic access, may use a sync.Mutex to synchronize.
