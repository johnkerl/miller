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

#include "json_parser.h"

const struct _json_value_t json_value_none;

#include <stdio.h>
#include <string.h>
#include <ctype.h>
#include <math.h>

typedef unsigned int json_uchar;

static unsigned char hex_value(json_char c) {
	if (isdigit(c))
		return c - '0';

	switch (c) {
		case 'a': case 'A': return 0x0A;
		case 'b': case 'B': return 0x0B;
		case 'c': case 'C': return 0x0C;
		case 'd': case 'D': return 0x0D;
		case 'e': case 'E': return 0x0E;
		case 'f': case 'F': return 0x0F;
		default: return 0xFF;
	}
}

// ----------------------------------------------------------------
typedef struct {
	unsigned long used_memory;

	unsigned int uint_max;
	unsigned long ulong_max;

	json_settings_t settings;
	int first_pass;

	const json_char * ptr;
	unsigned int cur_line, cur_col;

} json_state;

// ----------------------------------------------------------------
static void * default_alloc(size_t size, int zero, void * user_data) {
	return zero ? calloc(1, size) : malloc(size);
}

// ----------------------------------------------------------------
static void default_free(void * ptr, void * user_data) {
	free(ptr);
}

// ----------------------------------------------------------------
static void * json_alloc(json_state * state, unsigned long size, int zero) {
	if ((state->ulong_max - state->used_memory) < size)
		return 0;

	if (state->settings.max_memory
			&& (state->used_memory += size) > state->settings.max_memory)
	{
		return 0;
	}

	return state->settings.mem_alloc(size, zero, state->settings.user_data);
}

// ----------------------------------------------------------------
static int new_value(
	json_state * state,
	json_value_t ** top, json_value_t ** root,
	json_value_t ** alloc,
	json_type_t type)
{
	json_value_t * value;
	int values_size;

	if (!state->first_pass) {
		value = *top = *alloc;
		*alloc = (*alloc)->_reserved.next_alloc;

		if (!*root)
			*root = value;

		switch (value->type) {
			case JSON_ARRAY:
				if (value->u.array.length == 0)
					break;

				if (! (value->u.array.values = (json_value_t **) json_alloc
					(state, value->u.array.length * sizeof(json_value_t *), 0)) )
				{
					return 0;
				}

				value->u.array.length = 0;
				break;

			case JSON_OBJECT:
				if (value->u.object.length == 0)
					break;

				values_size = sizeof(*value->u.object.p.values) * value->u.object.length;

				if (! (value->u.object.p.values = (json_object_entry_t *) json_alloc
						(state, values_size + ((unsigned long) value->u.object.p.values), 0)) )
				{
					return 0;
				}

				value->_reserved.p.pobject_mem = (*(char **) &value->u.object.p.mem) + values_size; // xxx pun

				value->u.object.length = 0;
				break;

			case JSON_STRING:
				if (! (value->u.string.ptr = (json_char *) json_alloc
					(state, (value->u.string.length + 1) * sizeof (json_char), 0)) )
				{
					return 0;
				}

				value->u.string.length = 0;
				break;

			default:
				break;
		};

		return 1;
	}

	if (! (value = (json_value_t *) json_alloc
			(state, sizeof (json_value_t) + state->settings.value_extra, 1)))
	{
		return 0;
	}

	if (!*root)
		*root = value;

	value->type = type;
	value->parent = *top;

	#ifdef JSON_TRACK_SOURCE
		value->line = state->cur_line;
		value->col = state->cur_col;
	#endif

	if (*alloc)
		(*alloc)->_reserved.next_alloc = value;

	*alloc = *top = value;

	return 1;
}

// ----------------------------------------------------------------
#define WHITESPACE \
	case '\n': ++state.cur_line;  state.cur_col = 0; \
	case ' ': case '\t': case '\r'

#define STRING_ADD(b)  \
	do { if (!state.first_pass) string [string_length] = b;  ++string_length; } while (0);

#define LINE_AND_COL \
	state.cur_line, state.cur_col

static const long
	flag_next             = 1 << 0,
	flag_reproc           = 1 << 1,
	flag_need_comma       = 1 << 2,
	flag_seek_value       = 1 << 3,
	flag_escaped          = 1 << 4,
	flag_string           = 1 << 5,
	flag_need_colon       = 1 << 6,
	flag_done             = 1 << 7,
	flag_num_negative     = 1 << 8,
	flag_num_zero         = 1 << 9,
	flag_num_e            = 1 << 10,
	flag_num_e_got_sign   = 1 << 11,
	flag_num_e_negative   = 1 << 12,
	flag_line_comment     = 1 << 13,
	flag_block_comment    = 1 << 14;

// ================================================================
json_value_t * json_parse(const json_char * json, size_t length, char* error_buf) {
	json_settings_t settings = { 0 };
	json_char* prename_me = 0;
	return json_parse_ex(json, length, error_buf, &prename_me, &settings);
}

// ----------------------------------------------------------------
json_value_t * json_parse_ex(
	const json_char * json,
	size_t length,
	char * error_buf,
	json_char** pprename_me,
	json_settings_t* settings)
{
	json_char error [JSON_ERROR_MAX];
	const json_char * end;
	json_value_t * top, * root, * alloc = 0;
	json_state state = { 0 };
	long flags;
	long num_digits = 0, num_e = 0;
	json_int_t num_fraction = 0;
	*pprename_me = NULL;

	// Skip UTF-8 BOM
	if (length >= 3 && ((unsigned char) json [0]) == 0xEF
					&& ((unsigned char) json [1]) == 0xBB
					&& ((unsigned char) json [2]) == 0xBF)
	{
		json += 3;
		length -= 3;
	}

	error[0] = '\0';
	end = (json + length);

	memcpy(&state.settings, settings, sizeof(json_settings_t));

	if (!state.settings.mem_alloc)
		state.settings.mem_alloc = default_alloc;

	if (!state.settings.mem_free)
		state.settings.mem_free = default_free;

	memset(&state.uint_max, 0xFF, sizeof(state.uint_max));
	memset(&state.ulong_max, 0xFF, sizeof(state.ulong_max));

	state.uint_max -= 8; /* limit of how much can be added before next check */
	state.ulong_max -= 8;

	for (state.first_pass = 1; state.first_pass >= 0; --state.first_pass) {
		json_uchar uchar;
		unsigned char uc_b1, uc_b2, uc_b3, uc_b4;
		json_char * string = 0;
		unsigned int string_length = 0;

		top = root = 0;
		flags = flag_seek_value;

		state.cur_line = 1;

		for (state.ptr = json ;; ++state.ptr) {
			json_char* pb = (json_char*)((state.ptr == end) ? NULL : state.ptr);
			json_char   b = (state.ptr == end) ? 0 : *state.ptr;

			if (flags & flag_string) {
				if (!b) {
					sprintf(error, "Unexpected EOF in string (at %d:%d)", LINE_AND_COL);
					goto e_failed;
				}

				if (string_length > state.uint_max)
					goto e_overflow;

				if (flags & flag_escaped) {
					flags &= ~ flag_escaped;

					switch (b) {
						case 'b':  STRING_ADD('\b');  break;
						case 'f':  STRING_ADD('\f');  break;
						case 'n':  STRING_ADD('\n');  break;
						case 'r':  STRING_ADD('\r');  break;
						case 't':  STRING_ADD('\t');  break;
						case 'u':

						  if (end - state.ptr < 4 ||
								(uc_b1 = hex_value (*++state.ptr)) == 0xFF ||
								(uc_b2 = hex_value (*++state.ptr)) == 0xFF ||
								(uc_b3 = hex_value (*++state.ptr)) == 0xFF ||
								(uc_b4 = hex_value (*++state.ptr)) == 0xFF)
						  {
								sprintf(error, "Invalid character value `%c` (at %d:%d)", b, LINE_AND_COL);
								goto e_failed;
						  }

						  uc_b1 = (uc_b1 << 4) | uc_b2;
						  uc_b2 = (uc_b3 << 4) | uc_b4;
						  uchar = (uc_b1 << 8) | uc_b2;

						  if ((uchar & 0xF800) == 0xD800) {
								json_uchar uchar2;

								if (end - state.ptr < 6 || (*++state.ptr) != '\\' || (*++state.ptr) != 'u' ||
									 (uc_b1 = hex_value (*++state.ptr)) == 0xFF ||
									 (uc_b2 = hex_value (*++state.ptr)) == 0xFF ||
									 (uc_b3 = hex_value (*++state.ptr)) == 0xFF ||
									 (uc_b4 = hex_value (*++state.ptr)) == 0xFF)
								{
									 sprintf(error, "Invalid character value `%c` (at %d:%d)", b, LINE_AND_COL);
									 goto e_failed;
								}

								uc_b1 = (uc_b1 << 4) | uc_b2;
								uc_b2 = (uc_b3 << 4) | uc_b4;
								uchar2 = (uc_b1 << 8) | uc_b2;

								uchar = 0x010000 | ((uchar & 0x3FF) << 10) | (uchar2 & 0x3FF);
						  }

						  if (sizeof(json_char) >= sizeof(json_uchar) || (uchar <= 0x7F)) {
							  STRING_ADD((json_char) uchar);
							  break;
						  }

						  if (uchar <= 0x7FF) {
								if (state.first_pass) {
									string_length += 2;
								} else {
									string [string_length ++] = 0xC0 | (uchar >> 6);
									string [string_length ++] = 0x80 | (uchar & 0x3F);
								}

								break;
						  }

						  if (uchar <= 0xFFFF) {
								if (state.first_pass) {
									string_length += 3;
								} else {
									string [string_length ++] = 0xE0 | (uchar >> 12);
									string [string_length ++] = 0x80 | ((uchar >> 6) & 0x3F);
									string [string_length ++] = 0x80 | (uchar & 0x3F);
								}

								break;
						  }

						  if (state.first_pass) {
							  string_length += 4;
						  } else {  string [string_length ++] = 0xF0 | (uchar >> 18);
							  string [string_length ++] = 0x80 | ((uchar >> 12) & 0x3F);
							  string [string_length ++] = 0x80 | ((uchar >> 6) & 0x3F);
							  string [string_length ++] = 0x80 | (uchar & 0x3F);
						  }

						  break;

						default:
							STRING_ADD (b);
					};

					continue;
				}

				if (b == '\\') {
					flags |= flag_escaped;
					continue;
				}

				if (b == '"') {
					if (!state.first_pass)
						string [string_length] = 0;

					flags &= ~ flag_string;
					string = 0;

					switch (top->type) {
						case JSON_STRING:
							top->u.string.length = string_length;
							flags |= flag_next;
							break;

						case JSON_OBJECT:

							if (state.first_pass) {
								(*(json_char **) &top->u.object.p.mem) += string_length + 1; // xxx pun
							} else {
								top->u.object.p.values [top->u.object.length].name
									= (json_char *) top->_reserved.p.pobject_mem;

								top->u.object.p.values [top->u.object.length].name_length = string_length;

								(*(json_char **) &top->_reserved.p.pobject_mem) += string_length + 1; // xxx pun
							}

							flags |= flag_seek_value | flag_need_colon;
							continue;

						default:
							break;
					};
				} else {
					STRING_ADD(b);
					continue;
				}
			}

			if (state.settings.setting_flags & JSON_ENABLE_COMMENTS) {
				if (flags & (flag_line_comment | flag_block_comment)) {
					if (flags & flag_line_comment) {
						if (b == '\r' || b == '\n' || !b) {
							flags &= ~ flag_line_comment;
							--state.ptr;  /* so null can be reproc'd */
						}

						continue;
					}

					if (flags & flag_block_comment) {
						if (!b) {
							sprintf(error, "%d:%d: Unexpected EOF in block comment", LINE_AND_COL);
							goto e_failed;
						}

						if (b == '*' && state.ptr < (end - 1) && state.ptr [1] == '/') {
							flags &= ~ flag_block_comment;
							++state.ptr;  /* skip closing sequence */
						}

						continue;
					}
				} else if (b == '/') {
					if (! (flags & (flag_seek_value | flag_done)) && top->type != JSON_OBJECT) {
						sprintf(error, "%d:%d: Comment not allowed here", LINE_AND_COL);
						goto e_failed;
					}

					if (++state.ptr == end) {
						sprintf(error, "%d:%d: EOF unexpected", LINE_AND_COL);
						goto e_failed;
					}

					switch (b = *state.ptr) {
						case '/':
							flags |= flag_line_comment;
							continue;

						case '*':
							flags |= flag_block_comment;
							continue;

						default:
							sprintf(error, "%d:%d: Unexpected `%c` in comment opening sequence", LINE_AND_COL, b);
							goto e_failed;
					};
				}
			}

			if (flags & flag_done) {
				if (!b)
					break;
				if (state.settings.setting_flags & JSON_ENABLE_SEQUENTIAL_OBJECTS) {
					*pprename_me = pb + 1;
					break;
				}

				switch (b) {
					WHITESPACE:
						continue;

					default:
						sprintf(error, "%d:%d: Trailing text: `%c`", state.cur_line, state.cur_col, b);
						goto e_failed;
				};
			}

			if (flags & flag_seek_value) {
				switch (b) {
					WHITESPACE:
						continue;

					case ']':
						if (top && top->type == JSON_ARRAY) {
							flags = (flags & ~ (flag_need_comma | flag_seek_value)) | flag_next;
						} else {
							sprintf (error, "%d:%d: Unexpected ]", LINE_AND_COL);
							goto e_failed;
						}

						break;

					default:
						if (flags & flag_need_comma) {
							if (b == ',') {
								flags &= ~ flag_need_comma;
								continue;
							} else {
								sprintf(error, "%d:%d: Expected , before %c", state.cur_line, state.cur_col, b);
								goto e_failed;
							}
						}

						if (flags & flag_need_colon) {
							if (b == ':') {
								flags &= ~ flag_need_colon;
								continue;
							} else {
								sprintf(error, "%d:%d: Expected : before %c", state.cur_line, state.cur_col, b);
								goto e_failed;
							}
						}

						flags &= ~ flag_seek_value;

						switch (b) {
							case '{':
								if (!new_value(&state, &top, &root, &alloc, JSON_OBJECT))
									goto e_alloc_failure;
								continue;

							case '[':
								if (!new_value(&state, &top, &root, &alloc, JSON_ARRAY))
									goto e_alloc_failure;
								flags |= flag_seek_value;
								continue;

							case '"':
								if (!new_value(&state, &top, &root, &alloc, JSON_STRING))
									goto e_alloc_failure;
								flags |= flag_string;
								string = top->u.string.ptr;
								string_length = 0;
								continue;

							case 't':
								if ((end - state.ptr) < 3 || *(++state.ptr) != 'r' ||
									 *(++state.ptr) != 'u' || *(++state.ptr) != 'e')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &top, &root, &alloc, JSON_BOOLEAN))
									goto e_alloc_failure;

								top->u.boolean = 1;

								flags |= flag_next;
								break;

							case 'f':

								if ((end - state.ptr) < 4 || *(++state.ptr) != 'a' ||
									 *(++state.ptr) != 'l' || *(++state.ptr) != 's' ||
									 *(++state.ptr) != 'e')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &top, &root, &alloc, JSON_BOOLEAN))
									goto e_alloc_failure;

								flags |= flag_next;
								break;

							case 'n':
								if ((end - state.ptr) < 3 || *(++state.ptr) != 'u' ||
									 *(++state.ptr) != 'l' || *(++state.ptr) != 'l')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &top, &root, &alloc, JSON_NULL))
									goto e_alloc_failure;

								flags |= flag_next;
								break;

							default:
								if (isdigit (b) || b == '-') {
									if (!new_value(&state, &top, &root, &alloc, JSON_INTEGER))
										goto e_alloc_failure;

									if (!state.first_pass) {
										while (isdigit(b) || b == '+' || b == '-' || b == 'e' || b == 'E' || b == '.') {
											if ((++state.ptr) == end) {
												b = 0;
												break;
											}

											b = *state.ptr;
										}

										flags |= flag_next | flag_reproc;
										break;
									}

									flags &= ~ (flag_num_negative | flag_num_e |
													 flag_num_e_got_sign | flag_num_e_negative |
														 flag_num_zero);

									num_digits = 0;
									num_fraction = 0;
									num_e = 0;

									if (b != '-') {
										flags |= flag_reproc;
										break;
									}

									flags |= flag_num_negative;
									continue;
								} else {
									sprintf(error, "%d:%d: Unexpected `%c` when seeking value", LINE_AND_COL, b);
									goto e_failed;
								}
						};
				};
			} else {
				switch (top->type) {
				case JSON_OBJECT:

					switch (b) {
						WHITESPACE:
							continue;

						case '"':
							if (flags & flag_need_comma) {
								sprintf(error, "%d:%d: Expected , before \"", LINE_AND_COL);
								goto e_failed;
							}

							flags |= flag_string;

							string = (json_char *) top->_reserved.p.pobject_mem;
							string_length = 0;

							break;

						case '}':
							flags = (flags & ~ flag_need_comma) | flag_next;
							break;

						case ',':
							if (flags & flag_need_comma) {
								flags &= ~ flag_need_comma;
								break;
							}

						default:
							sprintf(error, "%d:%d: Unexpected `%c` in object", LINE_AND_COL, b);
							goto e_failed;
					};

					break;

				case JSON_INTEGER:
				case JSON_DOUBLE:
					if (isdigit(b)) {
						++num_digits;

						if (top->type == JSON_INTEGER || flags & flag_num_e) {
							if (! (flags & flag_num_e)) {
								if (flags & flag_num_zero) {
									sprintf(error, "%d:%d: Unexpected `0` before `%c`", LINE_AND_COL, b);
									goto e_failed;
								}

								if (num_digits == 1 && b == '0')
									flags |= flag_num_zero;
							} else {
								flags |= flag_num_e_got_sign;
								num_e = (num_e * 10) + (b - '0');
								continue;
							}

							top->u.integer = (top->u.integer * 10) + (b - '0');
							continue;
						}

						num_fraction = (num_fraction * 10) + (b - '0');
						continue;
					}

					if (b == '+' || b == '-') {
						if ( (flags & flag_num_e) && !(flags & flag_num_e_got_sign)) {
							flags |= flag_num_e_got_sign;

							if (b == '-')
								flags |= flag_num_e_negative;

							continue;
						}
					} else if (b == '.' && top->type == JSON_INTEGER) {
						if (!num_digits) {
							sprintf (error, "%d:%d: Expected digit before `.`", LINE_AND_COL);
							goto e_failed;
						}

						top->type = JSON_DOUBLE;
						top->u.dbl = (double) top->u.integer;

						num_digits = 0;
						continue;
					}

					if (! (flags & flag_num_e)) {
						if (top->type == JSON_DOUBLE) {
							if (!num_digits) {
								sprintf(error, "%d:%d: Expected digit after `.`", LINE_AND_COL);
								goto e_failed;
							}

							top->u.dbl += ((double) num_fraction) / (pow(10.0, (double) num_digits));
						}

						if (b == 'e' || b == 'E') {
							flags |= flag_num_e;

							if (top->type == JSON_INTEGER) {
								top->type = JSON_DOUBLE;
								top->u.dbl = (double) top->u.integer;
							}

							num_digits = 0;
							flags &= ~ flag_num_zero;

							continue;
						}
					} else {
						if (!num_digits) {
							sprintf(error, "%d:%d: Expected digit after `e`", LINE_AND_COL);
							goto e_failed;
						}

						top->u.dbl *= pow(10.0, (double) (flags & flag_num_e_negative ? - num_e : num_e));
					}

					if (flags & flag_num_negative) {
						if (top->type == JSON_INTEGER)
							top->u.integer = - top->u.integer;
						else
							top->u.dbl = - top->u.dbl;
					}

					flags |= flag_next | flag_reproc;
					break;

				default:
					break;
				};
			}

			if (flags & flag_reproc) {
				flags &= ~ flag_reproc;
				--state.ptr;
			}

			if (flags & flag_next) {
				flags = (flags & ~ flag_next) | flag_need_comma;

				if (!top->parent) {
					/* root value done */

					flags |= flag_done;
					continue;
				}

				if (top->parent->type == JSON_ARRAY)
					flags |= flag_seek_value;

				if (!state.first_pass) {
					json_value_t * parent = top->parent;

					switch (parent->type) {
						case JSON_OBJECT:
							parent->u.object.p.values[parent->u.object.length].pvalue = top;
							break;

						case JSON_ARRAY:
							parent->u.array.values[parent->u.array.length] = top;
							break;

						default:
							break;
					};
				}

				if ((++top->parent->u.array.length) > state.uint_max)
					goto e_overflow;

				top = top->parent;

				continue;
			}
		}

		alloc = root;
	}

	return root;

e_unknown_value:

	sprintf(error, "%d:%d: Unknown value", LINE_AND_COL);
	goto e_failed;

e_alloc_failure:

	strcpy(error, "Memory allocation failure");
	goto e_failed;

e_overflow:

	sprintf(error, "%d:%d: Too long (caught overflow)", LINE_AND_COL);
	goto e_failed;

e_failed:

	if (error_buf) {
		if (*error)
			strcpy(error_buf, error);
		else
			strcpy(error_buf, "Unknown error");
	}

	if (state.first_pass)
		alloc = root;

	while (alloc) {
		top = alloc->_reserved.next_alloc;
		state.settings.mem_free(alloc, state.settings.user_data);
		alloc = top;
	}

	if (!state.first_pass)
		json_value_free_ex(&state.settings, root);

	return 0;
}

// ----------------------------------------------------------------
void json_value_free_ex(json_settings_t * settings, json_value_t * pvalue) {
	json_value_t * cur_value;

	if (!pvalue)
		return;

	pvalue->parent = 0;

	while (pvalue) {
		switch (pvalue->type) {
			case JSON_ARRAY:

				if (!pvalue->u.array.length) {
					settings->mem_free(pvalue->u.array.values, settings->user_data);
					break;
				}

				pvalue = pvalue->u.array.values [--pvalue->u.array.length];
				continue;

			case JSON_OBJECT:
				if (!pvalue->u.object.length) {
					settings->mem_free(pvalue->u.object.p.values, settings->user_data);
					break;
				}

				pvalue = pvalue->u.object.p.values [--pvalue->u.object.length].pvalue;
				continue;

			case JSON_STRING:
				settings->mem_free(pvalue->u.string.ptr, settings->user_data);
				break;

			default:
				break;
		};

		cur_value = pvalue;
		pvalue = pvalue->parent;
		settings->mem_free(cur_value, settings->user_data);
	}
}

// ----------------------------------------------------------------
void json_value_free(json_value_t * pvalue) {
	json_settings_t settings = { 0 };
	settings.mem_free = default_free;
	json_value_free_ex(&settings, pvalue);
}

// ----------------------------------------------------------------
char* json_describe_type(json_type_t type) {
	switch(type) {
	case JSON_NONE:    return "JSON_NONE";    break;
	case JSON_OBJECT:  return "JSON_OBJECT";  break;
	case JSON_ARRAY:   return "JSON_ARRAY";   break;
	case JSON_INTEGER: return "JSON_INTEGER"; break;
	case JSON_DOUBLE:  return "JSON_DOUBLE";  break;
	case JSON_STRING:  return "JSON_STRING";  break;
	case JSON_BOOLEAN: return "JSON_BOOLEAN"; break;
	case JSON_NULL:    return "JSON_NULL";    break;
	default:           return "???";          break;
	}
}
