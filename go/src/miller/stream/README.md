The streamer uses Go channels to pipe together file-reads, to record-reading/parsing, to a chain of record-mappers, to record-writing/formatting, to terminal standard output.

This is the main sketch of Miller, invoked straight from `main()` within `mlr.go`.
