package cli

// Decide whether to insert a flatten or unflatten verb at the end of the
// chain.  See also repl/verbs.go which handles the same issue in the REPL.
//
// PROBLEM TO BE SOLVED:
//
// JSON and YAML have nested structures and CSV et al. do not. For example:
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
// APPROACH:
//
// Use the Principle of Least Surprise (POLS). Below, "JSON" stands for any
// format capable of representing nested structures natively -- currently
// JSON, JSON Lines, and YAML.
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
//
// * Overriding these: if the last verb the user has explicitly provided is
//   flatten, don't undo that by putting an unflatten right after.
//

// isNestable returns true for formats which can represent nested/array
// structures natively, and thus don't need auto-flatten/auto-unflatten.
func isNestable(format string) bool {
	return format == "json" || format == "jsonl" || format == "yaml"
}

func DecideFinalFlatten(writerOptions *TWriterOptions) bool {
	ofmt := writerOptions.OutputFileFormat
	if writerOptions.AutoFlatten {
		// JSON/YAML/JSON-Lines preserve nested/array structure natively, so
		// they never need flattening.
		//
		// DCF is excluded for a different reason: it's not nestable, but it
		// has its own hardcoded comma-list serialization for a fixed set of
		// field names (Depends, Recommends, etc. -- see
		// pkg/output/record_writer_dcf.go), which generic key-spreading
		// flatten would clobber.
		if !isNestable(ofmt) && ofmt != "dcf" {
			return true
		}
	}
	return false
}

func DecideFinalUnflatten(
	options *TOptions,
	verbSequences [][]string,
) bool {

	numVerbs := len(verbSequences)
	if numVerbs > 0 {
		lastVerbSequence := verbSequences[numVerbs-1]
		if len(lastVerbSequence) > 0 {
			lastVerbName := lastVerbSequence[0]
			if lastVerbName == "flatten" {
				return false
			}
		}
	}

	ifmt := options.ReaderOptions.InputFileFormat
	ofmt := options.WriterOptions.OutputFileFormat

	if options.WriterOptions.AutoUnflatten {
		if !isNestable(ifmt) {
			if isNestable(ofmt) {
				return true
			}
		}
	}
	return false
}
