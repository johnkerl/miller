#ifndef NETBSD_STRPTIME_H
#define NETBSD_STRPTIME_H

#include <time.h>
char* netbsd_strptime(const char *buf, const char *fmt, struct tm *tm);

#endif // NETBSD_STRPTIME_H
