#ifndef MLRESCAPE_H
#define MLRESCAPE_H

// Avoids shell-injection cases by replacing single-quote with backslash single-quote,
// then wrapping the entire result in initial and final single-quote.
char* alloc_file_name_escaped_for_popen(char* filename);

#endif // MLRESCAPE_H
