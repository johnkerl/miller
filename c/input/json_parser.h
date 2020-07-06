// ================================================================
// Copyright (C) 2012, 2013, 2014 James McLaughlin et al.  All rights reserved.
// https://github.com/udp/json-parser
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
// 1. Redistributions of source code must retain the above copyright
//   notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.
// ================================================================

#ifndef JSON_PARSER_H
#define JSON_PARSER_H

#ifndef json_char
	#define json_char char
#endif

#ifndef json_int_t
	#ifndef _MSC_VER
		#include <inttypes.h>
		#define json_int_t int64_t
	#else
		#define json_int_t __int64
	#endif
#endif

#include <stdlib.h>

// ----------------------------------------------------------------
typedef enum {
	JSON_NONE,
	JSON_OBJECT,
	JSON_ARRAY,
	JSON_INTEGER,
	JSON_DOUBLE,
	JSON_STRING,
	JSON_BOOLEAN,
	JSON_NULL
} json_type_t;

extern const struct _json_value_t json_value_none;

typedef struct _json_object_entry_t {
	 json_char * name;
	 unsigned int name_length;
	 struct _json_value_t * pvalue;
} json_object_entry_t;

typedef struct _json_value_t {
	struct _json_value_t * parent;

	json_type_t type;

	union {
		struct {
			unsigned int length;
			union {
				json_object_entry_t * values;
				char* mem;
			} p;
		} object;

		struct {
			unsigned int length;
			struct _json_value_t ** values;
		} array;

		// For Miller we want floating-point numbers to be preserved as-is, with however many decimal places
		// the user's input does or does not have, until/unless we do any math which modifies values.
		struct {
			unsigned int length;
			char* sval;
		} boolean;

		struct {
			unsigned int length;
			char* sval;
		} integer;

		struct {
			unsigned int length;
			char* sval;
		} dbl;

		struct {
			unsigned int length;
			json_char * ptr; /* null-terminated */
		} string;

	} u;

	union {
		struct _json_value_t * next_alloc;
		union {
			void * pvobject_mem;
			char * pobject_mem;
		} p;
	} _reserved;

	// Location of the value in the source JSON
	unsigned int line, col;

} json_value_t;

#define JSON_ERROR_MAX 128

// The end-of-item returned pointer enables us to handle input of the form
//
//   { "a" : 1 }
//   { "b" : 2 }
//   { "c" : 3 }
//
// in addition to
//
// [
//   { "a" : 1 }
//   { "b" : 2 }
//   { "c" : 3 }
// ]
//
// This is in line with what jq can handle. In this case, json_parse will return
// once for each top-level item and will give us back a pointer to the start of
// the rest of the input stream, so we can call json_parse on the rest until it is
// all exhausted.

json_value_t * json_parse(
	const json_char * json,
	size_t length,
	char * error_buf,
	json_char** ppend_of_item,
	long long *pline_number); // should be set to 0 by the caller before 1st call

json_value_t * json_parse_for_unit_test(
	const json_char * json,
	json_char** ppend_of_item);

void json_free_value(json_value_t *);

char* json_describe_type(json_type_t type);

void json_print_non_recursive(json_value_t* pvalue);
void json_print_recursive(json_value_t* pvalue);

#endif // JSON_PARSER_H
