# Status of the Go port

* This is not necessarily a full Go port of Miller. At the moment, it's a little spot for some experimentation. Things are very rough and very iterative.
* One reason Miller exists is to be a useful tool for myself and others; another is it's fun to write. At bare minimum, I'll re-teach myself some Go.
* In all likelihood though this will turn into a full port which will someday become Miller 6.0.
* Benefits:
  * The lack of a streaming (record-by-record) JSON reader in the C implementation (https://github.com/johnkerl/miller/issues/99) is immediately solved in the Go implementation.
  * The quoted-DKVP feature from https://github.com/johnkerl/miller/issues/266 will be easily addressed.
  * String/number-formatting issues in https://github.com/johnkerl/miller/issues/211 https://github.com/johnkerl/miller/issues/178 https://github.com/johnkerl/miller/issues/151 https://github.com/johnkerl/miller/issues/259 will be fixed during the Go port.
  * I think some DST/timezone issues such as https://github.com/johnkerl/miller/issues/359 will be easier to fix using the Go datetime library than using the C datetime library
  * The code will be easier to read and, I hope, easier for others to contribute to.
* In the meantime I will still keep fixing bugs, doing some features, in C on Miller 5.x.

