#ifndef MLRDATETIME_H
#define MLRDATETIME_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

// Seconds since the epoch.
double get_systime();

// These use the gmtime/timegm and strftime/strptime standard-library functions, with the addition
// of support for floating-point seconds since the epoch.
char* mlr_alloc_time_string_from_seconds(double seconds_since_the_epoch, char* format);
double mlr_seconds_from_time_string(char* string, char* format);

#endif // MLRDATETIME_H
