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

#ifndef _JSON_H
#define _JSON_H

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
#define JSON_ENABLE_COMMENTS           0x01
// xxx comment
#define JSON_ENABLE_SEQUENTIAL_OBJECTS 0x02

typedef struct {
	int setting_flags;
	unsigned long max_memory;

	// Custom allocator support (leave null to use malloc/free)
	void * (* mem_alloc) (size_t, int zero, void * user_data);
	void (* mem_free) (void *, void * user_data);
	void * user_data;  /* will be passed to mem_alloc and mem_free */
	size_t value_extra;  /* how much extra space to allocate for values? */
} json_settings_t;

typedef enum {
	json_none,
	json_object,
	json_array,
	json_integer,
	json_double,
	json_string,
	json_boolean,
	json_null

} json_type_t;

extern const struct _json_value_t json_value_none;

typedef struct _json_object_entry_t {
	 json_char * name;
	 unsigned int name_length;
	 struct _json_value_t * value;
} json_object_entry_t;

typedef struct _json_value_t {
	struct _json_value_t * parent;

	json_type_t type;

	union {
		int boolean;
		json_int_t integer;
		double dbl;

		struct {
			unsigned int length;
			json_char * ptr; /* null-terminated */
		} string;

		struct {
			unsigned int length;
			json_object_entry_t * values;
		} object;

		struct {
			unsigned int length;
			struct _json_value_t ** values;
		} array;

	} u;

	union {
		struct _json_value_t * next_alloc;
		void * object_mem;
	} _reserved;

	#ifdef JSON_TRACK_SOURCE
		// Location of the value in the source JSON
		unsigned int line, col;
	#endif

} json_value_t;

#define JSON_ERROR_MAX 128
json_value_t * json_parse(
	const json_char * json,
	size_t length,
	char*  error_buf);

json_value_t * json_parse_ex(
	const json_char * json,
	size_t length,
	char * error_buf,
	json_char** pprename_me,
	json_settings_t * settings);

void json_value_free(json_value_t *);

// Not usually necessary, unless you used a custom mem_alloc and now want to
// use a custom mem_free.
void json_value_free_ex(
	json_settings_t * settings,
	json_value_t *);

#endif // _JSON_H
