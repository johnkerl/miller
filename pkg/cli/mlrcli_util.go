package cli

// CheckArgCount is for flags with values, e.g. ["-n" "10"], while we're
// looking at the "-n": this let us see if the "10" slot exists. A non-nil
// return is a (multi-line) user-facing error, ready for printing as-is by the
// entrypoint layer.
func CheckArgCount(args []string, argi int, argc int, n int) error {
	if (argc - argi) < n {
		return FlagErrorf(
			"%s: option \"%s\" missing argument(s).\nPlease run \"%s --help\" for detailed usage information.",
			"mlr", args[argi], "mlr")
	}
	return nil
}

// SeparatorFromArg is for letting people do things like `--ifs pipe`
// rather than `--ifs '|'`.
func SeparatorFromArg(name string) string {
	sep, ok := SEPARATOR_NAMES_TO_VALUES[name]
	if ok {
		return sep
	}
	return name
}

// SeparatorRegexFromArg is for letting people do things like `--ifs-regex whitespace`
// rather than `--ifs '([ \t])+'`.
func SeparatorRegexFromArg(name string) string {
	sep, ok := SEPARATOR_REGEX_NAMES_TO_VALUES[name]
	if ok {
		return sep
	}
	return name
}
