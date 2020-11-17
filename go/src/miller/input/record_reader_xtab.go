package input

// typedef struct _lrec_reader_stdio_xtab_state_t {
// 	char*  ifs;
// 	char*  ips;
// 	int    ifslen;
// 	int    ipslen;
// 	int    allow_repeat_ips;
// 	int    do_auto_line_term;
// 	int    at_eof;
// 	comment_handling_t comment_handling;
// 	char*  comment_string;
// 	size_t line_length;
// } lrec_reader_stdio_xtab_state_t;

// // ----------------------------------------------------------------
// lrec_reader_t* lrec_reader_stdio_xtab_alloc(char* ifs, char* ips, int allow_repeat_ips,
// 	comment_handling_t comment_handling, char* comment_string)
// {
// 	lrec_reader_t* plrec_reader = mlr_malloc_or_die(sizeof(lrec_reader_t));
// 
// 	lrec_reader_stdio_xtab_state_t* pstate = mlr_malloc_or_die(sizeof(lrec_reader_stdio_xtab_state_t));
// 	pstate->ifs               = ifs;
// 	pstate->ips               = ips;
// 	pstate->ifslen            = strlen(ifs);
// 	pstate->ipslen            = strlen(ips);
// 	pstate->allow_repeat_ips  = allow_repeat_ips;
// 	pstate->do_auto_line_term = FALSE;
// 	pstate->at_eof            = FALSE;
// 	pstate->comment_handling  = comment_handling;
// 	pstate->comment_string    = comment_string;
// 	// This is used to track nominal line length over the file read. Bootstrap with a default length.
// 	pstate->line_length       = MLR_ALLOC_READ_LINE_INITIAL_SIZE;
// 
// 	if (streq(ifs, "auto")) {
// 		pstate->do_auto_line_term = TRUE;
// 		pstate->ifs = "\n";
// 		pstate->ifslen = 1;
// 	}
// 
// 	plrec_reader->pvstate       = (void*)pstate;
// 	plrec_reader->popen_func    = file_reader_stdio_vopen;
// 	plrec_reader->pclose_func   = file_reader_stdio_vclose;
// 	plrec_reader->pprocess_func = lrec_reader_stdio_xtab_process;
// 	plrec_reader->psof_func     = lrec_reader_stdio_xtab_sof;
// 	plrec_reader->pfree_func    = lrec_reader_stdio_xtab_free;
// 
// 	return plrec_reader;
// }

// static void lrec_reader_stdio_xtab_sof(void* pvstate, void* pvhandle) {
// 	lrec_reader_stdio_xtab_state_t* pstate = pvstate;
// 	pstate->at_eof = FALSE;
// }

// // ----------------------------------------------------------------
// static lrec_t* lrec_reader_stdio_xtab_process(void* pvstate, void* pvhandle, context_t* pctx) {
// 	FILE* input_stream = pvhandle;
// 	lrec_reader_stdio_xtab_state_t* pstate = pvstate;
// 
// 	if (pstate->at_eof)
// 		return NULL;
// 
// 	slls_t* pxtab_lines = slls_alloc();
// 
// 	while (TRUE) {
// 		char* line = NULL;
// 
// 		if (pstate->comment_handling == COMMENTS_ARE_DATA) {
// 			line = mlr_alloc_read_line_multiple_delimiter(input_stream, pstate->ifs, pstate->ifslen,
// 				&pstate->line_length);
// 		} else {
// 			line = mlr_alloc_read_line_multiple_delimiter_stripping_comments(input_stream,
// 				pstate->ifs, pstate->ifslen, &pstate->line_length,
// 				pstate->comment_handling, pstate->comment_string);
// 		}
// 
// 		if (line == NULL) { // EOF
// 			// EOF or blank line terminates the stanza.
// 			pstate->at_eof = TRUE;
// 			if (pxtab_lines->length == 0) {
// 				slls_free(pxtab_lines);
// 				return NULL;
// 			} else {
// 				return
// 					lrec_parse_stdio_xtab_multi_ips(pxtab_lines, pstate->ips, pstate->ipslen, pstate->allow_repeat_ips);
// 			}
// 
// 		} else if (*line == '\0') {
// 			free(line);
// 			if (pxtab_lines->length > 0) {
// 				return
// 					lrec_parse_stdio_xtab_multi_ips(pxtab_lines, pstate->ips, pstate->ipslen, pstate->allow_repeat_ips);
// 			}
// 
// 		} else {
// 			slls_append_with_free(pxtab_lines, line);
// 		}
// 	}
// }

// lrec_t* lrec_parse_stdio_xtab_multi_ips(slls_t* pxtab_lines, char* ips, int ipslen, int allow_repeat_ips) {
// 	lrec_t* prec = lrec_xtab_alloc(pxtab_lines);
// 
// 	for (sllse_t* pe = pxtab_lines->phead; pe != NULL; pe = pe->pnext) {
// 		char* line = pe->value;
// 		char* p = line;
// 		char* key = p;
// 
// 		while (*p != 0 && !streqn(p, ips, ipslen))
// 			p++; // Advance by only 1 in case of subsequent match
// 		if (*p == 0) {
// 			lrec_put(prec, key, "", NO_FREE);
// 		} else {
// 			while (*p != 0 && streqn(p, ips, ipslen)) {
// 				*p = 0;
// 				p += ipslen;
// 			}
// 			lrec_put(prec, key, p, NO_FREE);
// 		}
// 	}
// 
// 	return prec;
// }
