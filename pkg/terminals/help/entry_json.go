// Machine-readable (JSON) help, for `mlr help --as-json` and friends.
//
// This assembles the structured catalogs exposed by the verb, function, flag,
// and keyword registries into a single document, so AI agents and other tooling
// can model Miller's surface without scraping the human-readable prose. The
// plain (non-`--as-json`) help behavior is unchanged; `--as-json` only switches
// the rendering.
//
// Two equivalent ways to opt in:
//   - Per-call flag `--as-json` anywhere on a `mlr help ...` command line.
//   - Env var MLR_HELP_JSON set to a truthy value (1, true, yes).

package help

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/transformers"
	"github.com/johnkerl/miller/v6/pkg/version"
)

// catalogSchemaVersion is bumped whenever the shape of the JSON catalog
// changes. Agents and tools can use this (together with mlr_version) as a
// cache key: re-fetch only when either value changes.
const catalogSchemaVersion = 1

// CatalogForJSON is the top-level document emitted by `mlr help --as-json`
// with no further topic: the entire help catalog in one machine-readable
// object.
type CatalogForJSON struct {
	MlrVersion           string                          `json:"mlr_version"`
	CatalogSchemaVersion int                             `json:"catalog_schema_version"`
	Verbs                []*transformers.VerbInfoForJSON `json:"verbs"`
	Functions            []*cst.FunctionInfoForJSON      `json:"functions"`
	Flags                []*cli.FlagInfoForJSON          `json:"flags"`
	Keywords             []*cst.KeywordInfoForJSON       `json:"keywords"`
}

// wantJSONOutput returns true when the caller has opted in to JSON output via
// either the --as-json flag or a truthy MLR_HELP_JSON env var.
func wantJSONOutput(args []string) (bool, []string) {
	if isTruthyEnv(os.Getenv("MLR_HELP_JSON")) {
		// Env var wins; still strip any --as-json tokens so dispatch is clean.
		_, rest := extractAsJSONFlag(args)
		return true, rest
	}
	return extractAsJSONFlag(args)
}

// isTruthyEnv returns true for non-empty strings commonly used as boolean
// env-var truthy values: "1", "true", "yes" (case-insensitive).
func isTruthyEnv(v string) bool {
	switch v {
	case "1", "true", "True", "TRUE", "yes", "Yes", "YES":
		return true
	}
	return false
}

// extractAsJSONFlag removes any "--as-json" token from args, returning whether
// one was present along with the remaining args. The flag may appear anywhere
// (e.g. `mlr help --as-json` or `mlr help verb cat --as-json`).
func extractAsJSONFlag(args []string) (bool, []string) {
	found := false
	kept := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--as-json" {
			found = true
		} else {
			kept = append(kept, arg)
		}
	}
	return found, kept
}

// printAsJSON marshals v as indented JSON to stdout. Returns a process exit
// code.
func printAsJSON(v any) int {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("mlr help: could not render JSON: %v\n", err)
		return 1
	}
	fmt.Println(string(bytes))
	return 0
}

// helpJSON dispatches `mlr help --as-json [topic [names...]]`. With no topic
// it emits the full catalog; with a topic (verb/function/flag/keyword) it
// emits just those entries -- all of them if no names are given, or the named
// ones.
func helpJSON(args []string) int {
	if len(args) == 0 {
		return printAsJSON(buildFullCatalog())
	}

	topic := args[0]
	names := args[1:]

	switch topic {
	case "verb", "verbs":
		return printAsJSON(collectVerbs(names))
	case "function", "functions":
		return printAsJSON(collectFunctions(names))
	case "flag", "flags":
		return printAsJSON(collectFlags(names))
	case "keyword", "keywords":
		return printAsJSON(collectKeywords(names))
	default:
		fmt.Printf("mlr help --as-json: unsupported topic \"%s\".\n", topic)
		fmt.Printf("Supported: (no topic) for the full catalog, or one of: verb, function, flag, keyword.\n")
		return 1
	}
}

func buildFullCatalog() *CatalogForJSON {
	return &CatalogForJSON{
		MlrVersion:           version.STRING,
		CatalogSchemaVersion: catalogSchemaVersion,
		Verbs:                transformers.GetVerbInfosForJSON(),
		Functions:            cst.BuiltinFunctionManagerInstance.GetFunctionInfosForJSON(),
		Flags:                cli.FLAG_TABLE.GetFlagInfosForJSON(),
		Keywords:             cst.GetKeywordInfosForJSON(),
	}
}

func collectVerbs(names []string) []*transformers.VerbInfoForJSON {
	if len(names) == 0 {
		return transformers.GetVerbInfosForJSON()
	}
	infos := make([]*transformers.VerbInfoForJSON, 0, len(names))
	for _, name := range names {
		if info := transformers.GetVerbInfoForJSON(name); info != nil {
			infos = append(infos, info)
		}
	}
	return infos
}

func collectFunctions(names []string) []*cst.FunctionInfoForJSON {
	if len(names) == 0 {
		return cst.BuiltinFunctionManagerInstance.GetFunctionInfosForJSON()
	}
	infos := make([]*cst.FunctionInfoForJSON, 0, len(names))
	for _, name := range names {
		if info := cst.BuiltinFunctionManagerInstance.GetFunctionInfoForJSON(name); info != nil {
			infos = append(infos, info)
		}
	}
	return infos
}

func collectFlags(names []string) []*cli.FlagInfoForJSON {
	if len(names) == 0 {
		return cli.FLAG_TABLE.GetFlagInfosForJSON()
	}
	infos := make([]*cli.FlagInfoForJSON, 0, len(names))
	for _, name := range names {
		if info := cli.FLAG_TABLE.GetFlagInfoForJSON(name); info != nil {
			infos = append(infos, info)
		}
	}
	return infos
}

func collectKeywords(names []string) []*cst.KeywordInfoForJSON {
	if len(names) == 0 {
		return cst.GetKeywordInfosForJSON()
	}
	infos := make([]*cst.KeywordInfoForJSON, 0, len(names))
	for _, name := range names {
		if info := cst.GetKeywordInfoForJSON(name); info != nil {
			infos = append(infos, info)
		}
	}
	return infos
}
