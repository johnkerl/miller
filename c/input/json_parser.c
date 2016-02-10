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
typedef struct _json_state_{
	unsigned long used_memory;

	unsigned int uint_max;
	unsigned long ulong_max;

	json_settings_t settings;
	int first_pass;

	const json_char * ptr;
	unsigned int cur_line, cur_col;

} json_parser_state_t;

// ----------------------------------------------------------------
static void * default_alloc(size_t size, int zero, void * user_data) {
	return zero ? calloc(1, size) : malloc(size);
}

// ----------------------------------------------------------------
static void default_free(void * ptr, void * user_data) {
	free(ptr);
}

// ----------------------------------------------------------------
static void * json_alloc(json_parser_state_t * pstate, unsigned long size, int zero) {
	if ((pstate->ulong_max - pstate->used_memory) < size)
		return 0;

	if (pstate->settings.max_memory
			&& (pstate->used_memory += size) > pstate->settings.max_memory)
	{
		return 0;
	}

	return pstate->settings.mem_alloc(size, zero, pstate->settings.user_data);
}

// ----------------------------------------------------------------
static int new_value(
	json_parser_state_t * pstate,
	json_value_t ** ptop,
	json_value_t ** proot,
	json_value_t ** palloc,
	json_type_t type)
{
	json_value_t * value;
	int values_size;

	if (!pstate->first_pass) {
		value = *ptop = *palloc;
		*palloc = (*palloc)->_reserved.next_alloc;

		if (!*proot)
			*proot = value;

		switch (value->type) {
			case JSON_ARRAY:
				if (value->u.array.length == 0)
					break;

				if (! (value->u.array.values = (json_value_t **) json_alloc
					(pstate, value->u.array.length * sizeof(json_value_t *), 0)) )
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
						(pstate, values_size + ((unsigned long) value->u.object.p.values), 0)) )
				{
					return 0;
				}

				value->_reserved.p.pobject_mem = (*(char **) &value->u.object.p.mem) + values_size;

				value->u.object.length = 0;
				break;

			case JSON_STRING:
				if (! (value->u.string.ptr = (json_char *) json_alloc
					(pstate, (value->u.string.length + 1) * sizeof (json_char), 0)) )
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
			(pstate, sizeof (json_value_t) + pstate->settings.value_extra, 1)))
	{
		return 0;
	}

	if (!*proot)
		*proot = value;

	value->type = type;
	value->parent = *ptop;

	#ifdef JSON_TRACK_SOURCE
		value->line = pstate->cur_line;
		value->col = pstate->cur_col;
	#endif

	if (*palloc)
		(*palloc)->_reserved.next_alloc = value;

	*palloc = *ptop = value;

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
	FLAG_NEXT             = 1 << 0,
	FLAG_REPROC           = 1 << 1,
	FLAG_NEED_COMMA       = 1 << 2,
	FLAG_SEEK_VALUE       = 1 << 3,
	FLAG_ESCAPED          = 1 << 4,
	FLAG_STRING           = 1 << 5,
	FLAG_NEED_COLON       = 1 << 6,
	FLAG_DONE             = 1 << 7,
	FLAG_NUM_NEGATIVE     = 1 << 8,
	FLAG_NUM_ZERO         = 1 << 9,
	FLAG_NUM_E            = 1 << 10,
	FLAG_NUM_E_GOT_SIGN   = 1 << 11,
	FLAG_NUM_E_NEGATIVE   = 1 << 12,
	FLAG_LINE_COMMENT     = 1 << 13,
	FLAG_BLOCK_COMMENT    = 1 << 14;

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
	json_value_t * ptop, * proot, * palloc = 0;
	json_parser_state_t state = { 0 };
	long flags;
	long num_digits = 0, num_e = 0;
	json_int_t num_fraction = 0;
	*pprename_me = NULL; // xxx rename

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

		ptop = proot = 0;
		flags = FLAG_SEEK_VALUE;

		state.cur_line = 1;

		for (state.ptr = json ;; ++state.ptr) {
			json_char* pb = (json_char*)((state.ptr == end) ? NULL : state.ptr);
			json_char   b = (state.ptr == end) ? 0 : *state.ptr;

			if (flags & FLAG_STRING) {
				if (!b) {
					sprintf(error, "Unexpected EOF in string (at %d:%d)", LINE_AND_COL);
					goto e_failed;
				}

				if (string_length > state.uint_max)
					goto e_overflow;

				if (flags & FLAG_ESCAPED) {
					flags &= ~ FLAG_ESCAPED;

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
							STRING_ADD(b);
					};

					continue;
				}

				if (b == '\\') {
					flags |= FLAG_ESCAPED;
					continue;
				}

				if (b == '"') {
					if (!state.first_pass)
						string [string_length] = 0;

					flags &= ~ FLAG_STRING;
					string = 0;

					switch (ptop->type) {
						case JSON_STRING:
							ptop->u.string.length = string_length;
							flags |= FLAG_NEXT;
							break;

						case JSON_OBJECT:

							if (state.first_pass) {
								(*(json_char **) &ptop->u.object.p.mem) += string_length + 1;
							} else {
								ptop->u.object.p.values [ptop->u.object.length].name
									= (json_char *) ptop->_reserved.p.pobject_mem;

								ptop->u.object.p.values [ptop->u.object.length].name_length = string_length;

								(*(json_char **) &ptop->_reserved.p.pobject_mem) += string_length + 1;
							}

							flags |= FLAG_SEEK_VALUE | FLAG_NEED_COLON;
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
				if (flags & (FLAG_LINE_COMMENT | FLAG_BLOCK_COMMENT)) {
					if (flags & FLAG_LINE_COMMENT) {
						if (b == '\r' || b == '\n' || !b) {
							flags &= ~ FLAG_LINE_COMMENT;
							--state.ptr;  /* so null can be reproc'd */
						}

						continue;
					}

					if (flags & FLAG_BLOCK_COMMENT) {
						if (!b) {
							sprintf(error, "%d:%d: Unexpected EOF in block comment", LINE_AND_COL);
							goto e_failed;
						}

						if (b == '*' && state.ptr < (end - 1) && state.ptr [1] == '/') {
							flags &= ~ FLAG_BLOCK_COMMENT;
							++state.ptr;  /* skip closing sequence */
						}

						continue;
					}
				} else if (b == '/') {
					if (! (flags & (FLAG_SEEK_VALUE | FLAG_DONE)) && ptop->type != JSON_OBJECT) {
						sprintf(error, "%d:%d: Comment not allowed here", LINE_AND_COL);
						goto e_failed;
					}

					if (++state.ptr == end) {
						sprintf(error, "%d:%d: EOF unexpected", LINE_AND_COL);
						goto e_failed;
					}

					switch (b = *state.ptr) {
						case '/':
							flags |= FLAG_LINE_COMMENT;
							continue;

						case '*':
							flags |= FLAG_BLOCK_COMMENT;
							continue;

						default:
							sprintf(error, "%d:%d: Unexpected `%c` in comment opening sequence", LINE_AND_COL, b);
							goto e_failed;
					};
				}
			}

			if (flags & FLAG_DONE) {
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

			if (flags & FLAG_SEEK_VALUE) {
				switch (b) {
					WHITESPACE:
						continue;

					case ']':
						if (ptop && ptop->type == JSON_ARRAY) {
							flags = (flags & ~ (FLAG_NEED_COMMA | FLAG_SEEK_VALUE)) | FLAG_NEXT;
						} else {
							sprintf (error, "%d:%d: Unexpected ]", LINE_AND_COL);
							goto e_failed;
						}

						break;

					default:
						if (flags & FLAG_NEED_COMMA) {
							if (b == ',') {
								flags &= ~ FLAG_NEED_COMMA;
								continue;
							} else {
								sprintf(error, "%d:%d: Expected , before %c", state.cur_line, state.cur_col, b);
								goto e_failed;
							}
						}

						if (flags & FLAG_NEED_COLON) {
							if (b == ':') {
								flags &= ~ FLAG_NEED_COLON;
								continue;
							} else {
								sprintf(error, "%d:%d: Expected : before %c", state.cur_line, state.cur_col, b);
								goto e_failed;
							}
						}

						flags &= ~ FLAG_SEEK_VALUE;

						switch (b) {
							case '{':
								if (!new_value(&state, &ptop, &proot, &palloc, JSON_OBJECT))
									goto e_alloc_failure;
								continue;

							case '[':
								if (!new_value(&state, &ptop, &proot, &palloc, JSON_ARRAY))
									goto e_alloc_failure;
								flags |= FLAG_SEEK_VALUE;
								continue;

							case '"':
								if (!new_value(&state, &ptop, &proot, &palloc, JSON_STRING))
									goto e_alloc_failure;
								flags |= FLAG_STRING;
								string = ptop->u.string.ptr;
								string_length = 0;
								continue;

							case 't':
								if ((end - state.ptr) < 3 || *(++state.ptr) != 'r' ||
									 *(++state.ptr) != 'u' || *(++state.ptr) != 'e')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &ptop, &proot, &palloc, JSON_BOOLEAN))
									goto e_alloc_failure;

								// xxx
								ptop->u.boolean.nval = 1;

								flags |= FLAG_NEXT;
								break;

							case 'f':

								if ((end - state.ptr) < 4 || *(++state.ptr) != 'a' ||
									 *(++state.ptr) != 'l' || *(++state.ptr) != 's' ||
									 *(++state.ptr) != 'e')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &ptop, &proot, &palloc, JSON_BOOLEAN))
									goto e_alloc_failure;

								flags |= FLAG_NEXT;
								break;

							case 'n':
								if ((end - state.ptr) < 3 || *(++state.ptr) != 'u' ||
									 *(++state.ptr) != 'l' || *(++state.ptr) != 'l')
								{
									goto e_unknown_value;
								}

								if (!new_value(&state, &ptop, &proot, &palloc, JSON_NULL))
									goto e_alloc_failure;

								flags |= FLAG_NEXT;
								break;

							default:
								if (isdigit (b) || b == '-') {
									// xxx start
									if (!new_value(&state, &ptop, &proot, &palloc, JSON_INTEGER))
										goto e_alloc_failure;

									if (!state.first_pass) {
										while (isdigit(b) || b == '+' || b == '-' || b == 'e' || b == 'E' || b == '.') {
											if ((++state.ptr) == end) {
												b = 0;
												break;
											}

											b = *state.ptr;
										}

										flags |= FLAG_NEXT | FLAG_REPROC;
										break;
									}

									flags &= ~ (FLAG_NUM_NEGATIVE | FLAG_NUM_E |
													 FLAG_NUM_E_GOT_SIGN | FLAG_NUM_E_NEGATIVE | FLAG_NUM_ZERO);

									num_digits = 0;
									num_fraction = 0;
									num_e = 0;

									if (b != '-') {
										flags |= FLAG_REPROC;
										break;
									}

									flags |= FLAG_NUM_NEGATIVE;
									continue;
								} else {
									sprintf(error, "%d:%d: Unexpected `%c` when seeking value", LINE_AND_COL, b);
									goto e_failed;
								}
						};
				};
			} else {
				switch (ptop->type) {
				case JSON_OBJECT:

					switch (b) {
						WHITESPACE:
							continue;

						case '"':
							if (flags & FLAG_NEED_COMMA) {
								sprintf(error, "%d:%d: Expected , before \"", LINE_AND_COL);
								goto e_failed;
							}

							flags |= FLAG_STRING;

							string = (json_char *) ptop->_reserved.p.pobject_mem;
							string_length = 0;

							break;

						case '}':
							flags = (flags & ~ FLAG_NEED_COMMA) | FLAG_NEXT;
							break;

						case ',':
							if (flags & FLAG_NEED_COMMA) {
								flags &= ~ FLAG_NEED_COMMA;
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

						if (ptop->type == JSON_INTEGER || flags & FLAG_NUM_E) {
							if (! (flags & FLAG_NUM_E)) {
								if (flags & FLAG_NUM_ZERO) {
									sprintf(error, "%d:%d: Unexpected `0` before `%c`", LINE_AND_COL, b);
									goto e_failed;
								}

								if (num_digits == 1 && b == '0')
									flags |= FLAG_NUM_ZERO;
							} else {
								flags |= FLAG_NUM_E_GOT_SIGN;
								num_e = (num_e * 10) + (b - '0');
								continue;
							}

							// xxx
							ptop->u.integer.nval = (ptop->u.integer.nval * 10) + (b - '0');
							continue;
						}

						num_fraction = (num_fraction * 10) + (b - '0');
						continue;
					}

					if (b == '+' || b == '-') {
						if ( (flags & FLAG_NUM_E) && !(flags & FLAG_NUM_E_GOT_SIGN)) {
							flags |= FLAG_NUM_E_GOT_SIGN;

							if (b == '-')
								flags |= FLAG_NUM_E_NEGATIVE;

							continue;
						}
					} else if (b == '.' && ptop->type == JSON_INTEGER) {
						if (!num_digits) {
							sprintf (error, "%d:%d: Expected digit before `.`", LINE_AND_COL);
							goto e_failed;
						}

						ptop->type = JSON_DOUBLE;
						// xxx
						ptop->u.dbl.nval = (double) ptop->u.integer.nval;

						num_digits = 0;
						continue;
					}

					if (! (flags & FLAG_NUM_E)) {
						if (ptop->type == JSON_DOUBLE) {
							if (!num_digits) {
								sprintf(error, "%d:%d: Expected digit after `.`", LINE_AND_COL);
								goto e_failed;
							}

							ptop->u.dbl.nval += ((double) num_fraction) / (pow(10.0, (double) num_digits));
						}

						if (b == 'e' || b == 'E') {
							flags |= FLAG_NUM_E;

							if (ptop->type == JSON_INTEGER) {
								ptop->type = JSON_DOUBLE;
								ptop->u.dbl.nval = (double) ptop->u.integer.nval;
							}

							num_digits = 0;
							flags &= ~ FLAG_NUM_ZERO;

							continue;
						}
					} else {
						if (!num_digits) {
							sprintf(error, "%d:%d: Expected digit after `e`", LINE_AND_COL);
							goto e_failed;
						}

						ptop->u.dbl.nval *= pow(10.0, (double) (flags & FLAG_NUM_E_NEGATIVE ? - num_e : num_e));
					}

					if (flags & FLAG_NUM_NEGATIVE) {
						// xxx
						if (ptop->type == JSON_INTEGER)
							ptop->u.integer.nval = - ptop->u.integer.nval;
						else
							ptop->u.dbl.nval = - ptop->u.dbl.nval;
					}

					flags |= FLAG_NEXT | FLAG_REPROC;
					break;

				default:
					break;
				};
			}

			if (flags & FLAG_REPROC) {
				flags &= ~ FLAG_REPROC;
				--state.ptr;
			}

			if (flags & FLAG_NEXT) {
				flags = (flags & ~ FLAG_NEXT) | FLAG_NEED_COMMA;

				if (!ptop->parent) {
					/* root value done */

					flags |= FLAG_DONE;
					continue;
				}

				if (ptop->parent->type == JSON_ARRAY)
					flags |= FLAG_SEEK_VALUE;

				if (!state.first_pass) {
					json_value_t * parent = ptop->parent;

					switch (parent->type) {
						case JSON_OBJECT:
							parent->u.object.p.values[parent->u.object.length].pvalue = ptop;
							break;

						case JSON_ARRAY:
							parent->u.array.values[parent->u.array.length] = ptop;
							break;

						default:
							break;
					};
				}

				if ((++ptop->parent->u.array.length) > state.uint_max)
					goto e_overflow;

				ptop = ptop->parent;

				continue;
			}
		}

		palloc = proot;
	}

	return proot;

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
		palloc = proot;

	while (palloc) {
		ptop = palloc->_reserved.next_alloc;
		state.settings.mem_free(palloc, state.settings.user_data);
		palloc = ptop;
	}

	if (!state.first_pass)
		json_value_free_ex(&state.settings, proot);

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
