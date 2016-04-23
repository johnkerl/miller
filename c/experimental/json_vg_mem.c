// ================================================================
// For use with valgrind --leak-check=full
// ================================================================

#include <stdio.h>
#include <string.h>
#include <lib/mlrutil.h>
#include <input/json_parser.h>

// ----------------------------------------------------------------
static void parse(char* pinput) {
	json_char* json_input = (json_char*)pinput;
	json_value_t* parsed_top_level_json;
	json_char error_buf[JSON_ERROR_MAX];

	json_char* item_start = json_input;
	size_t length = strlen(pinput);

	while (TRUE) {
		parsed_top_level_json = json_parse(item_start, length, error_buf, &item_start);

		if (parsed_top_level_json == NULL) {
			fprintf(stderr, "Unable to parse JSON data: %s\n", error_buf);
			exit(1);
		}

		printf("\n");
		printf("----------------------------------------------------------------\n");
		printf("INPUT:\n");
		printf("%s\n", pinput);
		printf("\n");
		printf("OUTPUT:\n");
		json_print_recursive(parsed_top_level_json);

		json_free_value(parsed_top_level_json);

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
	parse("[3]");
	parse("[[3]]");
	parse("{\"a\":3}");
	parse("[{\"a\":3},{\"b\":4}]");

	parse(
		"{\n"
		"  \"z\":  {\n"
		"    \"pan\":  {\n"
		"      \"1\":  0.726803,\n"
		"      \"0\":  0.952618\n"
		"    },\n"
		"    \"eks\":  {\n"
		"      \"0\":  0.134189,\n"
		"      \"1\":  0.187885\n"
		"    },\n"
		"    \"wye\":  {\n"
		"      \"1\":  0.863624\n"
		"    },\n"
		"    \"zee\":  {\n"
		"      \"0\":  0.976181\n"
		"    },\n"
		"    \"hat\":  {\n"
		"      \"1\":  0.749551\n"
		"    }\n"
		"  }\n"
		"}\n"
	);

	parse(
		"[\n"
		"{ \"a\": \"pan\", \"b\": \"pan\", \"i\": 1, \"x\": 0.3467901443380824, \"y\": 0.7268028627434533 }\n"
		",{ \"a\": \"eks\", \"b\": \"pan\", \"i\": 2, \"x\": 0.7586799647899636, \"y\": 0.5221511083334797 }\n"
		",{ \"a\": \"wye\", \"b\": \"wye\", \"i\": 3, \"x\": 0.20460330576630303, \"y\": 0.33831852551664776 }\n"
		",{ \"a\": \"eks\", \"b\": \"wye\", \"i\": 4, \"x\": 0.38139939387114097, \"y\": 0.13418874328430463 }\n"
		",{ \"a\": \"wye\", \"b\": \"pan\", \"i\": 5, \"x\": 0.5732889198020006, \"y\": 0.8636244699032729 }\n"
		",{ \"a\": \"zee\", \"b\": \"pan\", \"i\": 6, \"x\": 0.5271261600918548, \"y\": 0.49322128674835697 }\n"
		",{ \"a\": \"eks\", \"b\": \"zee\", \"i\": 7, \"x\": 0.6117840605678454, \"y\": 0.1878849191181694 }\n"
		",{ \"a\": \"zee\", \"b\": \"wye\", \"i\": 8, \"x\": 0.5985540091064224, \"y\": 0.976181385699006 }\n"
		",{ \"a\": \"hat\", \"b\": \"wye\", \"i\": 9, \"x\": 0.03144187646093577, \"y\": 0.7495507603507059 }\n"
		",{ \"a\": \"pan\", \"b\": \"wye\", \"i\": 10, \"x\": 0.5026260055412137, \"y\": 0.9526183602969864 }\n"
		"]\n"
	);

	return 0;
}
