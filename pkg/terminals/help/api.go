// Exported accessors over the machine-readable help catalog, for callers
// outside this package -- in particular the `mlr mcp` server, which serves
// the same catalog/index/search over MCP that `mlr help --as-json` and
// `mlr which` serve over the command line.

package help

import (
	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/transformers"
)

// BuildFullCatalog returns the entire help catalog: the same document emitted
// by `mlr help --as-json`.
func BuildFullCatalog() *CatalogForJSON {
	return buildFullCatalog()
}

// CatalogSchemaVersion returns the current catalog-shape version; see the
// catalogSchemaVersion const.
func CatalogSchemaVersion() int {
	return catalogSchemaVersion
}

// BuildIndex returns the lightweight capability index: the same list emitted
// by `mlr help --as-json --index`.
func BuildIndex() []IndexEntryForJSON {
	return buildIndex()
}

// CollectVerbs returns the structured views of the named verbs, or all verbs
// if names is empty. Unknown names are skipped.
func CollectVerbs(names []string) []*transformers.VerbInfoForJSON {
	return collectVerbs(names)
}

// CollectFunctions is CollectVerbs for DSL builtin functions.
func CollectFunctions(names []string) []*cst.FunctionInfoForJSON {
	return collectFunctions(names)
}

// CollectFlags is CollectVerbs for main-level flags.
func CollectFlags(names []string) []*cli.FlagInfoForJSON {
	return collectFlags(names)
}

// CollectKeywords is CollectVerbs for DSL keywords.
func CollectKeywords(names []string) []*cst.KeywordInfoForJSON {
	return collectKeywords(names)
}

// WhichSearch scores every catalog item against the query and returns matches
// in descending score order, along with the same confidence signal `mlr
// which` reports via its exit code: true when the top result's name contains
// a query token.
func WhichSearch(query string) (results []WhichResultEntry, confident bool) {
	tokens := whichTokenize(query)
	if len(tokens) == 0 {
		return nil, false
	}
	results = whichSearch(tokens)
	confident = len(results) > 0 && results[0].nameHit
	return results, confident
}
