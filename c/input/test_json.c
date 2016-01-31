// ================================================================
// Copyright (C) 2015 Mirko Pasqualetti  All rights reserved.
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

#include <stdio.h>
#include <stdlib.h>
#include <sys/stat.h>

#include "json.h"

// Test for json.c
//
// Compile (static linking) with
//         gcc -o test_json -I.. test_json.c ../json.c -lm
//
// Compile (dynamic linking) with
//         gcc -o test_json -I.. test_json.c -lm -ljsonparser
//
// USAGE: ./test_json <json_file>

// ================================================================
static void print_depth_shift(int depth) {
	int j;
	for (j = 0; j < depth; j++) {
		printf(" ");
	}
}

static void process_value_aux(json_value_t* value, int depth);

static void process_object(json_value_t* value, int depth) {
	int length, x;
	if (value == NULL) {
		return;
	}
	length = value->u.object.length;
	for (x = 0; x < length; x++) {
		print_depth_shift(depth);
		printf("object[%d].name = %s\n", x, value->u.object.values[x].name);
		process_value_aux(value->u.object.values[x].value, depth+1);
	}
}

static void process_array(json_value_t* value, int depth) {
	int length, x;
	if (value == NULL) {
		return;
	}
	length = value->u.array.length;
	printf("array\n");
	for (x = 0; x < length; x++) {
		process_value_aux(value->u.array.values[x], depth);
	}
}

// xxx Miller top-levels will be JSON "object" (i.e. Miller hashmap)
// or JSON "array" (i.e. Miller list of hashmap). In the latter case
// all second-level objects must be objects.
static void process_value_aux(json_value_t* value, int depth) {
	if (value == NULL) {
		return;
	}
	if (value->type != json_object) {
		print_depth_shift(depth);
	}
	switch (value->type) {
		case json_none:
			printf("none\n");
			break;
		case json_object:
			process_object(value, depth+1);
			break;
		case json_array:
			process_array(value, depth+1);
			break;
		case json_integer:
			printf("int: %10" PRId64 "\n", value->u.integer);
			break;
		case json_double:
			printf("double: %f\n", value->u.dbl);
			break;
		case json_string:
			printf("string: %s\n", value->u.string.ptr);
			break;
		case json_boolean:
			printf("bool: %d\n", value->u.boolean);
			break;
		case json_null:
			printf("null\n");
			break;
	}
}
static void process_value(json_value_t* value) {
	process_value_aux(value, 0);
}

// ================================================================
int notmain(int argc, char** argv) {
	char* filename;
	FILE *fp;
	struct stat filestatus;
	int file_size;
	char* file_contents;
	json_char* json;
	json_value_t* value;

	if (argc != 2) {
		fprintf(stderr, "%s <file_json>\n", argv[0]);
		return 1;
	}
	filename = argv[1];

	if (stat(filename, &filestatus) != 0) {
		fprintf(stderr, "File %s not found\n", filename);
		return 1;
	}
	file_size = filestatus.st_size;
	file_contents = (char*)malloc(filestatus.st_size);
	if (file_contents == NULL) {
		fprintf(stderr, "Memory error: unable to allocate %d bytes\n", file_size);
		return 1;
	}

	fp = fopen(filename, "rt");
	if (fp == NULL) {
		fprintf(stderr, "Unable to open %s\n", filename);
		fclose(fp);
		free(file_contents);
		return 1;
	}
	if (fread(file_contents, file_size, 1, fp) != 1 ) {
		fprintf(stderr, "Unable t read content of %s\n", filename);
		fclose(fp);
		free(file_contents);
		return 1;
	}
	fclose(fp);

	json_char error_buf[JSON_ERROR_MAX];

	json = (json_char*)file_contents;

	json_settings_t settings = {
		.setting_flags = JSON_ENABLE_SEQUENTIAL_OBJECTS,
		.max_memory = 0
	};
	json_char* prename_me = json;
	int length = file_size;

	while (1) {

		//printf("--------------------------------\n\n");
		//printf("[%s]\n", prename_me);
		//printf("--------------------------------\n\n");
		value = json_parse_ex(prename_me, length, error_buf, &prename_me, &settings);

		if (value == NULL) {
			fprintf(stderr, "Unable to parse data: \"%s\"\n", error_buf);
			free(file_contents);
			exit(1);
		}

		process_value(value);

		json_value_free(value);

		if (prename_me == NULL)
			break;
		if (*prename_me == 0)
			break;
		//printf("\n");
		//printf("PAGE\n");
		//printf("\n");
		length -= (prename_me - json);
		json = prename_me;
	}

	free(file_contents);
	return 0;
}

#ifdef TEST_JSON_MAIN
int main(int argc, char** argv) {
	return notmain(argc, argv);
}
#endif
