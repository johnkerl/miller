// Machine-readable (JSON) accessors over the command-line flag table, for
// `mlr help --json` and similar tooling. Unlike verb options (which are
// prose-only), the flag table is already fully structured since it drives the
// actual command-line parser -- so this is a thin serialization layer.

package cli

// FlagInfoForJSON is the structured view of a single command-line flag. Section
// is the human-readable flag-section name (e.g. "CSV-only flags"); Arg is the
// argument placeholder in curly braces (e.g. "{filename}") or "" for boolean
// flags; AltNames carries alternate spellings (e.g. "--csv" alongside "-c").
type FlagInfoForJSON struct {
	Section  string   `json:"section"`
	Name     string   `json:"name"`
	AltNames []string `json:"alt_names,omitempty"`
	Arg      string   `json:"arg,omitempty"`
	Help     string   `json:"help"`
}

func makeFlagInfoForJSON(sectionName string, flag *Flag) *FlagInfoForJSON {
	return &FlagInfoForJSON{
		Section:  sectionName,
		Name:     flag.name,
		AltNames: flag.altNames,
		Arg:      flag.arg,
		Help:     flag.help,
	}
}

// GetFlagInfosForJSON returns the full flag catalog, grouped by section in
// table order, flattened into a single list with each flag tagged by section.
func (ft *FlagTable) GetFlagInfosForJSON() []*FlagInfoForJSON {
	infos := make([]*FlagInfoForJSON, 0)
	for _, section := range ft.sections {
		for i := range section.flags {
			infos = append(infos, makeFlagInfoForJSON(section.name, &section.flags[i]))
		}
	}
	return infos
}

// GetFlagInfoForJSON returns the structured view of a single flag by name
// (matching either its primary name or any alternate spelling), or nil if there
// is no such flag.
func (ft *FlagTable) GetFlagInfoForJSON(name string) *FlagInfoForJSON {
	for _, section := range ft.sections {
		for i := range section.flags {
			if section.flags[i].Owns(name) {
				return makeFlagInfoForJSON(section.name, &section.flags[i])
			}
		}
	}
	return nil
}
