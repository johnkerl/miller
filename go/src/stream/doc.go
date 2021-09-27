// Package stream is the high-level view of Miller. It uses Go channels to pipe
// together file-reads, to record-reading/parsing, to a chain of
// record-transformers, to record-writing/formatting, to terminal output.
package stream
