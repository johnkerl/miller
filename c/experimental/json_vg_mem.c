#include <stdio.h>

// For use with valgrind --leak-check=full

// ----------------------------------------------------------------
static void parse(void* pvinput) {
	json_char* json_input = (json_char*)pvinput;
	json_value_t* parsed_top_level_json;
	json_char error_buf[JSON_ERROR_MAX];

	json_char* item_start = json_input;
	int length = phandle->eof - phandle->sol;

	while (TRUE) {
		parsed_top_level_json = json_parse(item_start, length, error_buf, &item_start);

		if (parsed_top_level_json == NULL) {
			fprintf(stderr, "Unable to parse JSON data: %s\n", error_buf);
			exit(1);
		}

		json_print_recursive(json_value_t* pvalue);

		json_free_recursive(parsed_top_level_json);

		if (item_start == NULL)
			break;
		if (*item_start == 0)
			break;
		length -= (item_start - json_input);
		json_input = item_start;
	}
}

// ----------------------------------------------------------------
int main(int argc, char** argv) {
	parse("3");
	return 0;
}
