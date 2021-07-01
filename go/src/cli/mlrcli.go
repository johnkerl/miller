package cli

// ================================================================
// Stuff needing to be ported from C
// ================================================================

//// ----------------------------------------------------------------
//#define DEFAULT_OFMT                     "%lf"
//#define DEFAULT_OQUOTING                 QUOTE_MINIMAL
//#define DEFAULT_OOSVAR_FLATTEN_SEPARATOR ":"
//#define DEFAULT_COMMENT_STRING           "#"
//
//// ASCII 1f and 1e
//#define ASV_FS "\x1f"
//#define ASV_RS "\x1e"
//
//#define ASV_FS_FOR_HELP "0x1f"
//#define ASV_RS_FOR_HELP "0x1e"
//
//// Unicode code points U+241F and U+241E, encoded as UTF-8.
//#define USV_FS "\xe2\x90\x9f"
//#define USV_RS "\xe2\x90\x9e"
//
//#define USV_FS_FOR_HELP "U+241F (UTF-8 0xe2909f)"
//#define USV_RS_FOR_HELP "U+241E (UTF-8 0xe2909e)"

// ----------------------------------------------------------------
//static lhmss_t* singleton_pdesc_to_chars_map = nil;
//static lhmss_t* get_desc_to_chars_map() {
//	if (singleton_pdesc_to_chars_map == nil) {
//		singleton_pdesc_to_chars_map = lhmss_alloc();
//		lhmss_put(singleton_pdesc_to_chars_map, "cr",        "\r",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crcr",      "\r\r",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lf",        "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "lflf",      "\n\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlf",      "\r\n",     NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "crlfcrlf",  "\r\n\r\n", NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "tab",       "\t",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "space",     " ",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "comma",     ",",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "newline",   "\n",       NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "pipe",      "|",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "slash",     "/",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "colon",     ":",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "semicolon", ";",        NO_FREE);
//		lhmss_put(singleton_pdesc_to_chars_map, "equals",    "=",        NO_FREE);
//	}
//	return singleton_pdesc_to_chars_map;
//}

//// For displaying the default separators in on-line help
//static char* rebackslash(char* sep) {
//	if sep == "\r"))
//		return "\\r";
//	else if sep == "\n"))
//		return "\\n";
//	else if sep == "\r\n"))
//		return "\\r\\n";
//	else if sep == "\t"))
//		return "\\t";
//	else if sep == " "))
//		return "space";
//	else
//		return sep;
//}
