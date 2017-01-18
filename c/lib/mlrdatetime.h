#ifndef MLRDATETIME_H
#define MLRDATETIME_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

double get_systime();

// portable timegm replacement
time_t mlr_timegm (struct tm *ptm);

#endif // MLRDATETIME_H
