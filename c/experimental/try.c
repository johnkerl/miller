#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>

#define ISO8601_TIME_FORMAT "%Y-%m-%dT%H:%M:%SZ"
#define NZBUFLEN 63
//char *strptime(const char *s, const char *format, struct tm *ptm);

// ----------------------------------------------------------------
int main(int argc, char** argv) {
	time_t strftime_input = 1491399761;

	struct tm tm = *gmtime(&strftime_input);
	char* strftime_output = malloc(NZBUFLEN + 1);
	int written_length = strftime(strftime_output, NZBUFLEN, ISO8601_TIME_FORMAT, &tm);
	if (written_length > NZBUFLEN || written_length == 0) {
		strcpy(strftime_output, "strftime failed");
	}
	printf("strftime: %d -> %s\n", (int)strftime_input, strftime_output);

	char* strptime_input = "2017-03-04T05:06:07Z";
	char* strptime_output = strptime(strptime_input, ISO8601_TIME_FORMAT, &tm);
	if (strptime_output == NULL) {
		printf("Could not strptime(\"%s\", \"%s\").\n", strptime_input, ISO8601_TIME_FORMAT);
	} else {
		printf("strptime: %s ->\n", strptime_input);
		printf("  tm_sec    = %d\n",  tm.tm_sec);
		printf("  tm_min    = %d\n",  tm.tm_min);
		printf("  tm_hour   = %d\n",  tm.tm_hour);
		printf("  tm_mday   = %d\n",  tm.tm_mday);
		printf("  tm_mon    = %d\n",  tm.tm_mon);
		printf("  tm_year   = %d\n",  tm.tm_year);
		printf("  tm_wday   = %d\n",  tm.tm_wday);
		printf("  tm_yday   = %d\n",  tm.tm_yday);
		printf("  tm_isdst  = %d\n",  tm.tm_isdst);
		printf("  tm_zone   = %s\n",  tm.tm_zone);
		printf("  tm_gmtoff = %ld\n", tm.tm_gmtoff);
		printf("  remainder = \"%s\"\n", strptime_output);
	}

	return 0;
}
