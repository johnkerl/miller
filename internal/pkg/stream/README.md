The streamer uses Go channels to pipe together file-reads, to record-reading/parsing, to a chain of record-transformers, to record-writing/formatting, to terminal standard output.

This is the high-level sketch of Miller.
