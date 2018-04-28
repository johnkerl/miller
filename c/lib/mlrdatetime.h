#ifndef MLRDATETIME_H
#define MLRDATETIME_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

typedef enum _time_from_seconds_choice_t {
	TIME_FROM_SECONDS_GMT,
	TIME_FROM_SECONDS_LOCAL
} time_from_seconds_choice_t;

// Seconds since the epoch.
double get_systime();

// These use the gmtime/timegm and strftime/strptime standard-library functions, with the addition
// of support for floating-point seconds since the epoch.
char* mlr_alloc_time_string_from_seconds(double seconds_since_the_epoch, char* format,
	time_from_seconds_choice_t time_from_seconds_choice);
double mlr_seconds_from_time_string(char* string, char* format);

#endif // MLRDATETIME_H
