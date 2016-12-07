#ifndef FREE_FLAGS_H
#define FREE_FLAGS_H

// For use by various data structures including slls_t and lrec_t.  Some
// keys/values are dynamically allocated and should be freed the container's
// destructor, and some should not. Examples of the former include strduped
// keys/values; examples of the latter include data from string literals, or
// from mmapped file-input data.

#define NO_FREE          0x00
#define FREE_ENTRY_KEY   0x40
#define FREE_ENTRY_VALUE 0x04

#endif // FREE_FLAGS_H
