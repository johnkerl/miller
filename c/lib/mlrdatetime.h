#ifndef MLRDATETIME_H
#define MLRDATETIME_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

double get_systime();

// portable timegm replacement
time_t mlr_timegm (struct tm *ptm);

char* mlr_alloc_time_string_from_seconds(time_t seconds, char* format);
time_t mlr_seconds_from_time_string(char* string, char* format);

#endif // MLRDATETIME_H
