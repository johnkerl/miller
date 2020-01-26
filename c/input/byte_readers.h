#ifndef BYTE_READERS_H
#define BYTE_READERS_H
#include "input/byte_reader.h"

byte_reader_t* string_byte_reader_alloc();
byte_reader_t* stdio_byte_reader_alloc();

void string_byte_reader_free(byte_reader_t* pbr);
void stdio_byte_reader_free(byte_reader_t* pbr);

#endif // BYTE_READERS_H
