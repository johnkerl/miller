#ifndef AUX_ENTRIES_H
#define AUX_ENTRIES_H

// Handles 'mlr lecat' and any other one-off tool things which don't go through the record-streaming logic.
// If the argument after the basename (i.e. argv[1]) is recognized then this function doesn't return,
// invoking the code for that argument instead and exiting.
void do_aux_entries(int argc, char** argv);
void show_aux_entries(FILE* fp);

#endif // AUX_ENTRIES_H
