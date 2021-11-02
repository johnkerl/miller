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

// Miller to-do items:
// 1. Modify the data structures so sval & length are shared between boolean, integer, & double.
// 2. Just append a byte-range (start to end pointer) rather than making repeated per-character calls
//    to these functions.

#include "lib/mlrutil.h"
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

	int first_pass;

	const json_char * ptr;
	unsigned int cur_line, cur_col;

} json_parser_state_t;

// ----------------------------------------------------------------
static void * default_alloc(size_t size, int zero) {
	void* ptr = mlr_malloc_or_die(size);
	if (zero)
		memset(ptr, 0, size);
	return ptr;
}

// ----------------------------------------------------------------
static void * json_alloc(json_parser_state_t * pstate, unsigned long size, int zero) {
	if ((pstate->ulong_max - pstate->used_memory) < size)
		return 0;

	return default_alloc(size, zero);
}

// ----------------------------------------------------------------
static int new_value(
	json_parser_state_t * pstate,
	json_value_t ** ptop,
	json_value_t ** proot,
	json_value_t ** palloc,
	json_type_t type)
{
	json_value_t * pvalue;
	int values_size;

	if (!pstate->first_pass) {
		pvalue = *ptop = *palloc;
		*palloc = (*palloc)->_reserved.next_alloc;

		if (!*proot)
			*proot = pvalue;

		switch (pvalue->type) {
			case JSON_ARRAY:
				if (pvalue->u.array.length == 0)
					break;
				if (! (pvalue->u.array.values = (json_value_t **) json_alloc
					(pstate, pvalue->u.array.length * sizeof(json_value_t *), 0)) )
				{
					return 0;
				}
				pvalue->u.array.length = 0;
				break;

			case JSON_OBJECT:
				if (pvalue->u.object.length == 0)
					break;
				values_size = sizeof(*pvalue->u.object.p.values) * pvalue->u.object.length;
				if (! (pvalue->u.object.p.values = (json_object_entry_t *) json_alloc
						(pstate, values_size + ((unsigned long) pvalue->u.object.p.values), 0)) )
				{
					return 0;
				}
				pvalue->_reserved.p.pobject_mem = (*(char **) &pvalue->u.object.p.mem) + values_size;
				pvalue->u.object.length = 0;
				break;

			case JSON_STRING:
				if (! (pvalue->u.string.ptr = (json_char *) json_alloc
					(pstate, (pvalue->u.string.length + 1) * sizeof (json_char), 0)) )
				{
					return 0;
				}
				pvalue->u.string.length = 0;
				break;

			case JSON_BOOLEAN:
				if (! (pvalue->u.boolean.sval = (json_char *) json_alloc
					(pstate, (pvalue->u.boolean.length + 1) * sizeof (json_char), 1)) )
				{
					return 0;
				}
				pvalue->u.boolean.length = 0;
				break;

			case JSON_INTEGER:
				if (! (pvalue->u.integer.sval = (json_char *) json_alloc
					(pstate, (pvalue->u.integer.length + 5) * sizeof (json_char), 1)) )
				{
					return 0;
				}
				pvalue->u.integer.length = 0;
				break;

			case JSON_DOUBLE:
				if (! (pvalue->u.dbl.sval = (json_char *) json_alloc
					(pstate, (pvalue->u.dbl.length + 5) * sizeof (json_char), 1)) )
				{
					return 0;
				}
				pvalue->u.dbl.length = 0;
				break;

			default:
				break;
		}

		return 1;
	}

	if (! (pvalue = (json_value_t *) json_alloc(pstate, sizeof (json_value_t), 1))) {
		return 0;
	}

	if (!*proot)
		*proot = pvalue;

	pvalue->type = type;
	pvalue->parent = *ptop;

	pvalue->line = pstate->cur_line;
	pvalue->col = pstate->cur_col;

	if (*palloc)
		(*palloc)->_reserved.next_alloc = pvalue;

	*palloc = *ptop = pvalue;

	return 1;
}

// ----------------------------------------------------------------
#define WHITESPACE \
	case '\n': ++state.cur_line;  state.cur_col = 0; \
	case ' ': case '\t': case '\r'

#define STRING_ADD(b)  \
	do { if (!state.first_pass) string[string_length] = b;  ++string_length; } while (0);

static inline void boolean_sval_add(json_parser_state_t* pstate, json_value_t* ptop, char b) {
	if (!pstate->first_pass) {
		ptop->u.boolean.sval[ptop->u.boolean.length++] = b;
	} else {
		ptop->u.boolean.length++;
	}
}
static inline void integer_sval_add(json_parser_state_t* pstate, json_value_t* ptop, char b) {
	if (!pstate->first_pass) {
		ptop->u.integer.sval[ptop->u.integer.length++] = b;
	} else {
		ptop->u.integer.length++;
	}
}
static inline void dbl_sval_add(json_parser_state_t* pstate, json_value_t* ptop, char b) {
	if (!pstate->first_pass) {
		ptop->u.dbl.sval[ptop->u.dbl.length++] = b;
	} else {
		ptop->u.dbl.length++;
	}
}


static inline void boolean_sval_end(json_parser_state_t* pstate, json_value_t* ptop) {
	if (!pstate->first_pass) {
		ptop->u.boolean.sval[ptop->u.boolean.length] = 0;
	}
}
static inline void integer_sval_end(json_parser_state_t* pstate, json_value_t* ptop) {
	if (!pstate->first_pass) {
		ptop->u.integer.sval[ptop->u.integer.length] = 0;
	}
}
static inline void dbl_sval_end(json_parser_state_t* pstate, json_value_t* ptop) {
	if (!pstate->first_pass) {
		ptop->u.dbl.sval[ptop->u.dbl.length] = 0;
	}
}

static const long
	FLAG_NEXT             = 1 << 0,
	FLAG_REPROC           = 1 << 1,
	FLAG_NEED_COMMA       = 1 << 2,
	FLAG_SEEK_VALUE       = 1 << 3,
	FLAG_ESCAPED          = 1 << 4,
	FLAG_IN_STRING        = 1 << 5,
	FLAG_NEED_COLON       = 1 << 6,
	FLAG_DONE             = 1 << 7,
	FLAG_NUM_NEGATIVE     = 1 << 8,
	FLAG_NUM_ZERO         = 1 << 9,
	FLAG_NUM_E            = 1 << 10,
	FLAG_NUM_E_GOT_SIGN   = 1 << 11,
	FLAG_NUM_E_NEGATIVE   = 1 << 12;

// ================================================================
json_value_t * json_parse(
	const json_char * json,
	size_t length,
	char * error_buf,
	json_char** ppend_of_item,
	long long *pline_number) // should be set to 0 by the caller before 1st call
{
	json_char error[JSON_ERROR_MAX];
	const json_char * end;
	json_value_t * ptop, * proot, * palloc = 0;
	json_parser_state_t state = { 0 };
	long flags;
	long num_digits = 0, num_e = 0;
	json_int_t num_fraction = 0;
	*ppend_of_item = NULL;

	// Skip UTF-8 BOM
	if (length >= 3 && ((unsigned char) json[0]) == 0xEF
					&& ((unsigned char) json[1]) == 0xBB
					&& ((unsigned char) json[2]) == 0xBF)
	{
		json += 3;
		length -= 3;
	}

	error[0] = '\0';
	end = (json + length);

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

		state.cur_line = *pline_number + 1;

		for (state.ptr = json ;; ++state.ptr) {
			json_char* pb = (json_char*)((state.ptr == end) ? NULL : state.ptr);
			json_char   b = (state.ptr == end) ? 0 : *state.ptr;

			if (flags & FLAG_IN_STRING) {
				if (!b) {
					sprintf(error, "Unexpected EOF in string at line %d column %d.",
						state.cur_line, state.cur_col);
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
								sprintf(error, "Invalid character value `%c` at line %d column %d.",
									b, state.cur_line, state.cur_col);
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
									sprintf(error, "Invalid character value `%c` at line %d column %d.",
										b, state.cur_line, state.cur_col);
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
									string[string_length ++] = 0xC0 | (uchar >> 6);
									string[string_length ++] = 0x80 | (uchar & 0x3F);
								}

								break;
						  }

						  if (uchar <= 0xFFFF) {
								if (state.first_pass) {
									string_length += 3;
								} else {
									string[string_length ++] = 0xE0 | (uchar >> 12);
									string[string_length ++] = 0x80 | ((uchar >> 6) & 0x3F);
									string[string_length ++] = 0x80 | (uchar & 0x3F);
								}

								break;
						  }

						  if (state.first_pass) {
							  string_length += 4;
						  } else {  string[string_length ++] = 0xF0 | (uchar >> 18);
							  string[string_length ++] = 0x80 | ((uchar >> 12) & 0x3F);
							  string[string_length ++] = 0x80 | ((uchar >> 6) & 0x3F);
							  string[string_length ++] = 0x80 | (uchar & 0x3F);
						  }

						  break;

						default:
							STRING_ADD(b);
					}

					continue;
				}

				if (b == '\\') {
					flags |= FLAG_ESCAPED;
					continue;
				}

				if (b == '"') {
					if (!state.first_pass)
						string[string_length] = 0;

					flags &= ~ FLAG_IN_STRING;
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
								ptop->u.object.p.values[ptop->u.object.length].name
									= (json_char *) ptop->_reserved.p.pobject_mem;

								ptop->u.object.p.values[ptop->u.object.length].name_length = string_length;

								(*(json_char **) &ptop->_reserved.p.pobject_mem) += string_length + 1;
							}

							flags |= FLAG_SEEK_VALUE | FLAG_NEED_COLON;
							continue;

						default:
							break;
					}
				} else {
					STRING_ADD(b);
					continue;
				}
			}

			if (flags & FLAG_DONE) {
				if (!b)
					break;

				*ppend_of_item = pb + 1;
				break;

				state.cur_col++;
				switch (b) {
					WHITESPACE:
						continue;

					default:
						sprintf(error, "Line %d column %d: Trailing text: `%c`",
							state.cur_line, state.cur_col, b);
						goto e_failed;
				}
			}

			if (flags & FLAG_SEEK_VALUE) {
				state.cur_col++;
				switch (b) {
					WHITESPACE:
						continue;

					case ']':
						if (ptop && ptop->type == JSON_ARRAY) {
							flags = (flags & ~ (FLAG_NEED_COMMA | FLAG_SEEK_VALUE)) | FLAG_NEXT;
						} else {
							sprintf (error, "Line %d column %d: Unexpected ]",
								state.cur_line, state.cur_col);
							goto e_failed;
						}

						break;

					default:
						if (flags & FLAG_NEED_COMMA) {
							if (b == ',') {
								flags &= ~ FLAG_NEED_COMMA;
								continue;
							} else {
								sprintf(error, "Line %d column %d: Expected , before %c",
									state.cur_line, state.cur_col, b);
								goto e_failed;
							}
						}

						if (flags & FLAG_NEED_COLON) {
							if (b == ':') {
								flags &= ~ FLAG_NEED_COLON;
								continue;
							} else {
								sprintf(error, "Line %d column %d: Expected : before %c",
									state.cur_line, state.cur_col, b);
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
								flags |= FLAG_IN_STRING;
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

								boolean_sval_add(&state, ptop, 't');
								boolean_sval_add(&state, ptop, 'r');
								boolean_sval_add(&state, ptop, 'u');
								boolean_sval_add(&state, ptop, 'e');
								boolean_sval_end(&state, ptop);

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
								boolean_sval_add(&state, ptop, 'f');
								boolean_sval_add(&state, ptop, 'a');
								boolean_sval_add(&state, ptop, 'l');
								boolean_sval_add(&state, ptop, 's');
								boolean_sval_add(&state, ptop, 'e');
								boolean_sval_end(&state, ptop);

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
									// Start of new number
									if (!new_value(&state, &ptop, &proot, &palloc, JSON_INTEGER))
										goto e_alloc_failure;

									if (!state.first_pass) {
										while (isdigit(b) || b == '+' || b == '-' || b == 'e' || b == 'E' || b == '.') {

											switch (ptop->type) {
											case JSON_INTEGER:
												integer_sval_add(&state, ptop, b);
												break;
											case JSON_DOUBLE:
												dbl_sval_add(&state, ptop, b);
												break;
											default:
												break;
											}

											if ((++state.ptr) == end) {

												switch (ptop->type) {
												case JSON_INTEGER:
													integer_sval_end(&state, ptop);
													break;
												case JSON_DOUBLE:
													dbl_sval_end(&state, ptop);
													break;
												default:
													break;
												}

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
									sprintf(error, "Line %d column %d: Unexpected `0x%02x` when seeking value",
										state.cur_line, state.cur_col, (unsigned)b);
									goto e_failed;
								}
						}
				}
			} else {
				switch (ptop->type) {
				case JSON_OBJECT:

					state.cur_col++;
					switch (b) {
						WHITESPACE:
							continue;

						case '"':
							if (flags & FLAG_NEED_COMMA) {
								sprintf(error, "Line %d column %d: Expected , before \"",
									state.cur_line, state.cur_col);
								goto e_failed;
							}

							flags |= FLAG_IN_STRING;

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
							sprintf(error, "Line %d column %d: Unexpected `%c` in object",
								state.cur_line, state.cur_col, b);
							goto e_failed;
					}

					break;

				case JSON_INTEGER:
				case JSON_DOUBLE:
					if (isdigit(b)) {
						++num_digits;

						if (ptop->type == JSON_INTEGER || flags & FLAG_NUM_E) {
							if (! (flags & FLAG_NUM_E)) {
								if (flags & FLAG_NUM_ZERO) {
									sprintf(error, "Line %d column %d: Unexpected `0` before `%c`",
										state.cur_line, state.cur_col, b);
									goto e_failed;
								}

								if (num_digits == 1 && b == '0')
									flags |= FLAG_NUM_ZERO;
							} else {
								flags |= FLAG_NUM_E_GOT_SIGN;
								num_e = (num_e * 10) + (b - '0');
								continue;
							}

							integer_sval_add(&state, ptop, b);
							continue;
						}

						integer_sval_add(&state, ptop, b);
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
							sprintf (error, "Line %d column %d: Expected digit before `.`",
								state.cur_line, state.cur_col);
							goto e_failed;
						}

						ptop->type = JSON_DOUBLE;
						ptop->u.dbl.length = ptop->u.integer.length;
						ptop->u.dbl.sval = ptop->u.integer.sval;
						dbl_sval_add(&state, ptop, b);

						num_digits = 0;
						continue;
					}

					if (! (flags & FLAG_NUM_E)) {
						if (ptop->type == JSON_DOUBLE) {
							if (!num_digits) {
								sprintf(error, "Line %d column %d: Expected digit after `.`",
									state.cur_line, state.cur_col);
								goto e_failed;
							}
						}

						if (b == 'e' || b == 'E') {
							flags |= FLAG_NUM_E;

							if (ptop->type == JSON_INTEGER) {
								ptop->type = JSON_DOUBLE;
								ptop->u.dbl.length = ptop->u.integer.length;
								ptop->u.dbl.sval = ptop->u.integer.sval;
								dbl_sval_add(&state, ptop, b);
							}

							num_digits = 0;
							flags &= ~ FLAG_NUM_ZERO;

							continue;
						}
					} else {
						if (!num_digits) {
							sprintf(error, "Line %d column %d: Expected digit after `e`",
								state.cur_line, state.cur_col);
							goto e_failed;
						}
					}

					flags |= FLAG_NEXT | FLAG_REPROC;
					break;

				default:
					break;
				}
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
					}
				}

				if ((++ptop->parent->u.array.length) > state.uint_max)
					goto e_overflow;

				ptop = ptop->parent;

				continue;
			}
		}

		palloc = proot;
	}

	*pline_number = state.cur_line;
	return proot;

e_unknown_value:

	sprintf(error, "Line %d column %d: Unknown value",
		state.cur_line, state.cur_col);
	goto e_failed;

e_alloc_failure:

	strcpy(error, "Memory allocation failure");
	goto e_failed;

e_overflow:

	sprintf(error, "Line %d column %d: Too long (caught overflow)",
		state.cur_line, state.cur_col);
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
		free(palloc);
		palloc = ptop;
	}

	if (!state.first_pass)
		json_free_value(proot);

	*pline_number = state.cur_line;
	return 0;
}

json_value_t * json_parse_for_unit_test(
	const json_char * json,
	json_char** ppend_of_item)
{
	json_char error_buf[JSON_ERROR_MAX];
	long long line_number = 0;
	return json_parse(json, strlen(json), error_buf, ppend_of_item, &line_number);
}

// ----------------------------------------------------------------
void json_free_value(json_value_t * pvalue) {
	json_value_t * cur_value;

	if (!pvalue)
		return;

	pvalue->parent = 0;

	while (pvalue) {
		switch (pvalue->type) {
			case JSON_ARRAY:
				if (!pvalue->u.array.length) {
					free(pvalue->u.array.values);
					break;
				}
				pvalue = pvalue->u.array.values[--pvalue->u.array.length];
				continue;

			case JSON_OBJECT:
				if (!pvalue->u.object.length) {
					free(pvalue->u.object.p.values);
					break;
				}
				pvalue = pvalue->u.object.p.values[--pvalue->u.object.length].pvalue;
				continue;

			case JSON_STRING:
				free(pvalue->u.string.ptr);
				break;

			case JSON_BOOLEAN:
				free(pvalue->u.boolean.sval);
				break;

			case JSON_INTEGER:
				free(pvalue->u.integer.sval);
				break;

			case JSON_DOUBLE:
				free(pvalue->u.dbl.sval);
				break;

			default:
				break;
		}

		cur_value = pvalue;
		pvalue = pvalue->parent;
		free(cur_value);
	}
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

// ----------------------------------------------------------------
const char* leader = "  ";
static void leader_print(int depth) {
	for (int i = 0; i < depth; i++)
		printf("%s", leader);
}

static void json_print_non_recursive_aux(json_value_t* pvalue, int depth) {
	leader_print(depth);
	if (pvalue == NULL) {
		printf("pvalue=NULL\n");
		return;
	}
	printf("type=%s", json_describe_type(pvalue->type));
	switch(pvalue->type) {
	case JSON_NONE:
		break;
	case JSON_OBJECT:
		printf(",length=%d", pvalue->u.object.length);
		break;
	case JSON_ARRAY:
		printf(",length=%d", pvalue->u.object.length);
		break;
	case JSON_INTEGER:
		printf(",length=%d", pvalue->u.integer.length);
		printf(",sval=\"%s\"", pvalue->u.integer.sval);
		break;
	case JSON_DOUBLE:
		printf(",length=%d", pvalue->u.dbl.length);
		printf(",sval=\"%s\"", pvalue->u.dbl.sval);
		break;
	case JSON_STRING:
		printf(",length=%d", pvalue->u.string.length);
		printf(",ptr=\"%s\"", pvalue->u.string.ptr);
		break;
	case JSON_BOOLEAN:
		printf(",length=%d", pvalue->u.boolean.length);
		printf(",sval=\"%s\"", pvalue->u.boolean.sval);
		break;
	case JSON_NULL:
		break;
	}
	printf("\n");
}

static void json_print_recursive_aux(json_value_t* pvalue, int depth) {
	json_print_non_recursive_aux(pvalue, depth);
	if (pvalue == NULL)
		return;
	if (pvalue->type == JSON_OBJECT) {
		for (int i = 0; i < pvalue->u.object.length; i++) {
			leader_print(depth+1);
			printf("key=\"%s\"\n", pvalue->u.object.p.values[i].name);
			json_print_recursive_aux(pvalue->u.object.p.values[i].pvalue, depth+1);
		}
	} else if (pvalue->type == JSON_ARRAY) {
		for (int i = 0; i < pvalue->u.array.length; i++) {
			leader_print(depth+1);
			printf("index=%d\n", i);
			json_print_recursive_aux(pvalue->u.array.values[i], depth+1);
		}
	}
}

// ----------------------------------------------------------------
void json_print_non_recursive(json_value_t* pvalue) {
	json_print_non_recursive_aux(pvalue, 0);
}

void json_print_recursive(json_value_t* pvalue) {
	json_print_recursive_aux(pvalue, 0);
}
