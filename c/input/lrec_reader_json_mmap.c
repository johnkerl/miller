// mmap: easy pointer math
// stdio from file: stat, alloc, read
// stdio from stdin: realloc w/ page-size fread

// note @ mlr -h: no streaing for JSON input. No records are processed until EOF is seen.

// paginated:
//   json parse || error msg
// produce sllv of items

// sllv processing:
//   insist sllv.length == 1 & is array & each array item is an object,
//   or each sllv item is an object
// for each item:
//   loop over k/v pairs in the object and insist on level-1 only.
