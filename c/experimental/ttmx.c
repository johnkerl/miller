#include <stdio.h>
#include <string.h>
#include <time.h>
char * strptimex(const char *buf, const char *fmt, struct tm *tm);
int main(void) {
	struct tm tm;
	char* foo = strptimex("2017-01-09T11:47:48Z", "%Y-%m-%dT%H:%M:%SZ", &tm);
	printf("[%s]\n", foo);
	return 0;
}
