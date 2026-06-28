// Machine-readable (JSON) accessors over the built-in-function catalog. These
// mirror the human-readable listings elsewhere in this file's neighbors (e.g.
// ListBuiltinFunctionUsages) but return structured data for `mlr help --json`
// and any other tooling -- AI agents in particular -- which needs to model
// Miller's function surface without scraping prose.

package cst

// FunctionInfoForJSON is the structured, marshalable view of a single built-in
// function. Field contents match what the text help shows: Arity is the same
// string produced by describeNargs (e.g. "1", "2,3", "1-4", "variadic"), and
// Help is JoinHelp()'d so source-code newlines are collapsed.
type FunctionInfoForJSON struct {
	Name     string   `json:"name"`
	Class    string   `json:"class"`
	Arity    string   `json:"arity"`
	Help     string   `json:"help"`
	Examples []string `json:"examples,omitempty"`
}

func makeFunctionInfoForJSON(info *BuiltinFunctionInfo) *FunctionInfoForJSON {
	return &FunctionInfoForJSON{
		Name:     info.name,
		Class:    string(info.class),
		Arity:    describeNargs(info),
		Help:     info.JoinHelp(),
		Examples: info.examples,
	}
}

// GetFunctionInfosForJSON returns the full function catalog in source-table
// (insertion) order, matching the human-readable `mlr help usage-functions`.
func (mgr *BuiltinFunctionManager) GetFunctionInfosForJSON() []*FunctionInfoForJSON {
	infos := make([]*FunctionInfoForJSON, 0, len(*mgr.lookupTable))
	for i := range *mgr.lookupTable {
		infos = append(infos, makeFunctionInfoForJSON(&(*mgr.lookupTable)[i]))
	}
	return infos
}

// GetFunctionInfoForJSON returns the structured view of a single function, or
// nil if there is no such function.
func (mgr *BuiltinFunctionManager) GetFunctionInfoForJSON(name string) *FunctionInfoForJSON {
	info := mgr.LookUp(name)
	if info == nil {
		return nil
	}
	return makeFunctionInfoForJSON(info)
}
