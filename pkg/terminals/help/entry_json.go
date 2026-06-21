// Machine-readable (JSON) help, for `mlr help --json` and friends.
//
// This assembles the structured catalogs exposed by the verb, function, flag,
// and keyword registries into a single document, so AI agents and other tooling
// can model Miller's surface without scraping the human-readable prose. The
// plain (non-`--json`) help behavior is unchanged; `--json` only switches the
// rendering.

package help

import (
	"encoding/json"
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/transformers"
	"github.com/johnkerl/miller/v6/pkg/version"
)

// CatalogForJSON is the top-level document emitted by `mlr help --json` with no
// further topic: the entire help catalog in one machine-readable object.
type CatalogForJSON struct {
	MlrVersion string                          `json:"mlr_version"`
	Verbs      []*transformers.VerbInfoForJSON `json:"verbs"`
	Functions  []*cst.FunctionInfoForJSON      `json:"functions"`
	Flags      []*cli.FlagInfoForJSON          `json:"flags"`
	Keywords   []*cst.KeywordInfoForJSON       `json:"keywords"`
}

// extractJSONFlag removes any "--json" token from args, returning whether one
// was present along with the remaining args. The flag may appear anywhere
// (e.g. `mlr help --json` or `mlr help verb cat --json`).
func extractJSONFlag(args []string) (bool, []string) {
	jsonMode := false
	kept := make([]string, 0, len(args))
	for _, arg := range args {
		if arg == "--json" {
			jsonMode = true
		} else {
			kept = append(kept, arg)
		}
	}
	return jsonMode, kept
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

// helpJSON dispatches `mlr help --json [topic [names...]]`. With no topic it
// emits the full catalog; with a topic (verb/function/flag/keyword) it emits
// just those entries -- all of them if no names are given, or the named ones.
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
		fmt.Printf("mlr help --json: unsupported topic \"%s\".\n", topic)
		fmt.Printf("Supported: (no topic) for the full catalog, or one of: verb, function, flag, keyword.\n")
		return 1
	}
}

func buildFullCatalog() *CatalogForJSON {
	return &CatalogForJSON{
		MlrVersion: version.STRING,
		Verbs:      transformers.GetVerbInfosForJSON(),
		Functions:  cst.BuiltinFunctionManagerInstance.GetFunctionInfosForJSON(),
		Flags:      cli.FLAG_TABLE.GetFlagInfosForJSON(),
		Keywords:   cst.GetKeywordInfosForJSON(),
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
