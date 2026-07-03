// In-process MCP tools: list_capabilities and which. These serve straight
// from the compiled-in help registries -- no subprocess -- via the exported
// accessors in pkg/terminals/help.

package mcp

import (
	"context"
	"fmt"

	"github.com/johnkerl/miller/v6/pkg/cli"
	"github.com/johnkerl/miller/v6/pkg/dsl/cst"
	"github.com/johnkerl/miller/v6/pkg/terminals/help"
	"github.com/johnkerl/miller/v6/pkg/transformers"
	"github.com/johnkerl/miller/v6/pkg/version"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type listCapabilitiesInput struct {
	Index bool     `json:"index,omitempty" jsonschema:"Return only {kind name summary} triples across the whole surface -- the cheap first call."`
	Kind  string   `json:"kind,omitempty" jsonschema:"Restrict to one catalog kind: verb / function / flag / keyword."`
	Names []string `json:"names,omitempty" jsonschema:"With kind: return full entries for just these names. Unknown names are skipped."`
}

// listCapabilitiesOutput is the one output shape for all three call styles:
// index-only, kind-filtered, or the full catalog. mlr_version and
// catalog_schema_version together form a cache key: the content can only
// change when one of them does.
type listCapabilitiesOutput struct {
	MlrVersion           string                          `json:"mlr_version"`
	CatalogSchemaVersion int                             `json:"catalog_schema_version"`
	Index                []help.IndexEntryForJSON        `json:"index,omitempty"`
	Verbs                []*transformers.VerbInfoForJSON `json:"verbs,omitempty"`
	Functions            []*cst.FunctionInfoForJSON      `json:"functions,omitempty"`
	Flags                []*cli.FlagInfoForJSON          `json:"flags,omitempty"`
	Keywords             []*cst.KeywordInfoForJSON       `json:"keywords,omitempty"`
}

func listCapabilitiesHandler(
	_ context.Context,
	_ *mcpsdk.CallToolRequest,
	input listCapabilitiesInput,
) (*mcpsdk.CallToolResult, listCapabilitiesOutput, error) {
	output := listCapabilitiesOutput{
		MlrVersion:           version.STRING,
		CatalogSchemaVersion: help.CatalogSchemaVersion(),
	}

	if input.Index {
		output.Index = help.BuildIndex()
		return nil, output, nil
	}

	switch input.Kind {
	case "":
		catalog := help.BuildFullCatalog()
		output.Verbs = catalog.Verbs
		output.Functions = catalog.Functions
		output.Flags = catalog.Flags
		output.Keywords = catalog.Keywords
	case "verb", "verbs":
		output.Verbs = help.CollectVerbs(input.Names)
	case "function", "functions":
		output.Functions = help.CollectFunctions(input.Names)
	case "flag", "flags":
		output.Flags = help.CollectFlags(input.Names)
	case "keyword", "keywords":
		output.Keywords = help.CollectKeywords(input.Names)
	default:
		return nil, output, fmt.Errorf(
			"unsupported kind %q: use verb, function, flag, or keyword", input.Kind)
	}

	return nil, output, nil
}

type whichInput struct {
	Query string `json:"query" jsonschema:"Natural-language intent e.g. \"join two files on a key\"."`
}

type whichOutput struct {
	Confident bool                    `json:"confident"`
	Results   []help.WhichResultEntry `json:"results"`
}

func whichHandler(
	_ context.Context,
	_ *mcpsdk.CallToolRequest,
	input whichInput,
) (*mcpsdk.CallToolResult, whichOutput, error) {
	if input.Query == "" {
		return nil, whichOutput{}, fmt.Errorf("query must be non-empty")
	}
	results, confident := help.WhichSearch(input.Query)
	if results == nil {
		results = []help.WhichResultEntry{}
	}
	return nil, whichOutput{Confident: confident, Results: results}, nil
}
