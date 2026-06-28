package help

import (
	"encoding/json"
	"testing"
)

// TestFullCatalogIsValidJSON guards the `mlr help --as-json` contract: the
// assembled catalog must round-trip through encoding/json and carry a non-empty
// entry in each of the four sections.
func TestFullCatalogIsValidJSON(t *testing.T) {
	catalog := buildFullCatalog()

	bytes, err := json.Marshal(catalog)
	if err != nil {
		t.Fatalf("catalog failed to marshal: %v", err)
	}

	var roundTrip CatalogForJSON
	if err := json.Unmarshal(bytes, &roundTrip); err != nil {
		t.Fatalf("catalog failed to unmarshal: %v", err)
	}

	if catalog.MlrVersion == "" {
		t.Error("catalog mlr_version is empty")
	}
	if catalog.CatalogSchemaVersion == 0 {
		t.Error("catalog catalog_schema_version is zero")
	}
	if len(catalog.Verbs) == 0 {
		t.Error("catalog has no verbs")
	}
	if len(catalog.Functions) == 0 {
		t.Error("catalog has no functions")
	}
	if len(catalog.Flags) == 0 {
		t.Error("catalog has no flags")
	}
	if len(catalog.Keywords) == 0 {
		t.Error("catalog has no keywords")
	}
}

// TestCatalogEntriesArePopulated checks that the structured fields agents rely
// on are actually filled in, not just present.
func TestCatalogEntriesArePopulated(t *testing.T) {
	catalog := buildFullCatalog()

	for _, verb := range catalog.Verbs {
		if verb.Name == "" {
			t.Error("found a verb with an empty name")
		}
		if verb.UsageText == "" {
			t.Errorf("verb %q has empty usage_text", verb.Name)
		}
	}

	for _, function := range catalog.Functions {
		if function.Name == "" {
			t.Error("found a function with an empty name")
		}
		if function.Class == "" {
			t.Errorf("function %q has empty class", function.Name)
		}
		if function.Arity == "" {
			t.Errorf("function %q has empty arity", function.Name)
		}
		if function.Help == "" {
			t.Errorf("function %q has empty help", function.Name)
		}
	}

	for _, flag := range catalog.Flags {
		if flag.Name == "" {
			t.Error("found a flag with an empty name")
		}
		if flag.Section == "" {
			t.Errorf("flag %q has empty section", flag.Name)
		}
	}

	for _, keyword := range catalog.Keywords {
		if keyword.Name == "" {
			t.Error("found a keyword with an empty name")
		}
		if keyword.Help == "" {
			t.Errorf("keyword %q has empty help", keyword.Name)
		}
	}
}

// TestPerTopicLookups exercises the single-entry lookups used by
// `mlr help <topic> <name> --json`, including the not-found path.
func TestPerTopicLookups(t *testing.T) {
	if got := collectVerbs([]string{"cat"}); len(got) != 1 || got[0].Name != "cat" {
		t.Errorf("verb lookup for cat: got %+v", got)
	}
	if got := collectVerbs([]string{"no-such-verb"}); len(got) != 0 {
		t.Errorf("verb lookup for bogus name should be empty, got %+v", got)
	}
	if got := collectFunctions([]string{"splitax"}); len(got) != 1 || got[0].Name != "splitax" {
		t.Errorf("function lookup for splitax: got %+v", got)
	}
	if got := collectKeywords([]string{"ENV"}); len(got) != 1 || got[0].Name != "ENV" {
		t.Errorf("keyword lookup for ENV: got %+v", got)
	}
	if got := collectFlags([]string{"--ifs"}); len(got) != 1 || got[0].Name != "--ifs" {
		t.Errorf("flag lookup for --ifs: got %+v", got)
	}
}
