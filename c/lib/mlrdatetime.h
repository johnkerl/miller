#ifndef MLRDATETIME_H
#define MLRDATETIME_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

double get_systime();

char* mlr_alloc_time_string_from_seconds(double seconds_since_the_epoch, char* format);
time_t mlr_seconds_from_time_string(char* string, char* format);

#endif // MLRDATETIME_H
