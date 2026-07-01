// Migration-progress tracker for structured verb options (Tier-2 catalog).
// Analogous to FLAG_TABLE.NilCheck() for flags, but non-fatal because verb
// Options are migrated incrementally across PRs.

package transformers

import "fmt"

// VerbOptionsNilCheck prints a migration summary: how many verbs have been
// given structured Options, and which verbs still have nil Options (i.e.
// are not yet migrated to Tier-2). Intended to be invoked from the help
// terminal and from regression tests to track progress.
func VerbOptionsNilCheck() {
	var unmigrated []string
	for i := range TRANSFORMER_LOOKUP_TABLE {
		if TRANSFORMER_LOOKUP_TABLE[i].Options == nil {
			unmigrated = append(unmigrated, TRANSFORMER_LOOKUP_TABLE[i].Verb)
		}
	}
	total := len(TRANSFORMER_LOOKUP_TABLE)
	migrated := total - len(unmigrated)
	fmt.Printf("Verb options migration: %d/%d migrated.\n", migrated, total)
	if len(unmigrated) > 0 {
		fmt.Printf("Unmigrated verbs (Options == nil):\n")
		for _, verb := range unmigrated {
			fmt.Printf("  %s\n", verb)
		}
	}
}
