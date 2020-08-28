# Readme

Package `md` is used to extract the text in code segments of a markdown file.

Code segments are enclosed in triple backticks: "```", eg:

```
A : B b | c ;
```

`md.GetSource(input string) string` returns a string with the same number of bytes as the input with:

- All characters inside code segments preserved in place;
- All space characters (`' ', '\n', '\r', '\t'`) outside code segments preserved in place;
- All non-space characters outside code segments replace by `' '` (space).
- The enclosing backticks of code segments replaced with spaces, i.e.: ` "```" ` replaced with `" "`.
