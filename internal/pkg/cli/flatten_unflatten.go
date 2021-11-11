package cli

// ================================================================
// Decide whether to insert a flatten or unflatten verb at the end of the
// chain.  See also repl/verbs.go which handles the same issue in the REPL.
//
// ----------------------------------------------------------------
// PROBLEM TO BE SOLVED:
//
// JSON has nested structures and CSV et al. do not. For example:
// {
//   "req" : {
//     "method": "GET",
//     "path":   "api/check",
//   }
// }
//
// For CSV we flatten this down to
//
// {
//   "req.method": "GET",
//   "req.path":   "api/check"
// }
//
// ----------------------------------------------------------------
// APPROACH:
//
// Use the Principle of Least Surprise (POLS).
//
// * If input is JSON and output is JSON:
//   o Records can be nested from record-read
//   o They remain that way through the Miller record-processing stream
//   o They are nested on record-write
//   o No action needs to be taken
//
// * If input is JSON and output is non-JSON:
//   o Records can be nested from record-read
//   o They remain that way through the Miller record-processing stream
//   o On record-write, nested structures will be converted to string (carriage
//     returns and all) using json_stringify. People *might* want this but
//     (using POLS) we will (by default) AUTO-FLATTEN for them. There is a
//     --no-auto-unflatten CLI flag for those who want it.
//
// * If input is non-JSON and output is non-JSON:
//   o If there is a "req.method" field, people should be able to do
//     'mlr sort -f req.method' with no surprises. (Again, POLS.) Therefore
//     no auto-unflatten on input.  People can insert an unflatten verb
//     into their verb chain if they really want unflatten for non-JSON
//     files.
//   o The DSL can make nested data, so AUTO-FLATTEN at output.
//
// * If input is non-JSON and output is JSON:
//   o Default is to auto-unflatten at output.
//   o There is a --no-auto-unflatten for those who want it.
// ================================================================

func DecideFinalFlatten(writerOptions *TWriterOptions) bool {
	ofmt := writerOptions.OutputFileFormat
	if writerOptions.AutoFlatten {
		if ofmt != "json" {
			return true
		}
	}
	return false
}

func DecideFinalUnflatten(options *TOptions) bool {
	ifmt := options.ReaderOptions.InputFileFormat
	ofmt := options.WriterOptions.OutputFileFormat

	if options.WriterOptions.AutoUnflatten {
		if ifmt != "json" {
			if ofmt == "json" {
				return true
			}
		}
	}
	return false
}
