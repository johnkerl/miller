// ================================================================
// strftime.c
// Original Author:	G. Haley
// Additions from:	Eric Blake, Corinna Vinschen
// Changes to allow dual use as wcstime, also:	Craig Howland
//
// Places characters into the array pointed to by s as controlled by the string
// pointed to by format. If the total number of resulting characters including
// the terminating null character is not more than maxsize, returns the number
// of characters placed into the array pointed to by s (not including the
// terminating null character); otherwise zero is returned and the contents of
// the array indeterminate.
// ================================================================

#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>

static size_t strftime_impl(
	char* s,
	size_t maxsize,
	const char* format,
	const struct tm *tim_p/*,
	struct __locale_t *locale*/);

//static int iso_year_adjust(const struct tm *tim_p);

// ----------------------------------------------------------------
size_t mlr_strftime(
	char *restrict s,
	size_t maxsize,
	const char *restrict format,
	const struct tm *restrict tim_p)
{
	return strftime_impl(s, maxsize, format, tim_p/*, locale*/);
}

//#ifdef _WANT_C99_TIME_FORMATS
//  era_info_t *era_info = NULL;
//  alt_digits_t *alt_digits = NULL;
//  size_t ret = __strftime(s, maxsize, format, tim_p, __get_current_locale(), &era_info, &alt_digits);
//  if (era_info)
//    era_info_free(era_info);
//  if (alt_digits)
//    alt_digits_free(alt_digits);
//  return ret;
//#else /* !_WANT_C99_TIME_FORMATS */
//  return __strftime(s, maxsize, format, tim_p, __get_current_locale(),
//		     NULL, NULL);
//#endif /* !_WANT_C99_TIME_FORMATS */

#define CHECK_LENGTH()	if (len < 0 || (count += len) >= maxsize) return 0

///* Enforce the coding assumptions that YEAR_BASE is positive.  (%C, %Y, etc.) */
#define YEAR_BASE 1900 // xxx temp
#define UINT_MAX 0xffffffff
//#if YEAR_BASE < 0
//#  error "YEAR_BASE < 0"
//#endif

// ----------------------------------------------------------------
static size_t strftime_impl(
	char* s,
	size_t maxsize,
	const char* format,
	const struct tm *tim_p/*,
	struct __locale_t *locale*/)
{
	size_t count = 0;

	int len = 0;
//	const char *ctloc;
//	size_t i, ctloclen;
	size_t i;
	size_t ctloclen;
	char alt;
	char pad;
	unsigned long width;
//	int tzset_called = 0;
//	const struct lc_time_T *_CurrentTimeLocale = __get_time_locale(locale);

	const char* format_walker = format;
	for (;;) {
		while (*format_walker && *format_walker != '%') {
			if (count < maxsize - 1)
				s[count++] = *format_walker++;
			else
				return 0;
		}
		if (*format_walker == '\0')
			break;

		format_walker++;
		pad = '\0';
		width = 0;

		// POSIX-1.2008 feature: '0' and '+' modifiers require 0-padding with
		// slightly different semantics.
		if (*format_walker == '0' || *format_walker == '+')
			pad = *format_walker++;

		// POSIX-1.2008 feature: A minimum field width can be specified.
		if (*format_walker >= '1' && *format_walker <= '9') {
			char *fp;
			width = strtoul(format_walker, &fp, 10);
			format_walker = fp;
		}

		alt = '\0';
		if (*format_walker == 'E') {
			alt = *format_walker++;
#ifdef _WANT_C99_TIME_FORMATS
			if (!*era_info && *_CurrentTimeLocale->era)
				*era_info = get_era_info(tim_p, _CurrentTimeLocale->era);
#endif /* _WANT_C99_TIME_FORMATS */

		} else if (*format_walker == 'O') {
			alt = *format_walker++;
#ifdef _WANT_C99_TIME_FORMATS
			if (!*alt_digits && *_CurrentTimeLocale->alt_digits)
				*alt_digits = alt_digits_alloc(_CurrentTimeLocale->alt_digits);
#endif /* _WANT_C99_TIME_FORMATS */
		}

		switch (*format_walker) {
//		case 'a':
//			// _ctloc(wday[tim_p->tm_wday]);
//			for (i = 0; i < ctloclen; i++) {
//				if (count < maxsize - 1)
//					s[count++] = ctloc[i];
//				else
//					return 0;
//				}
//				break;

//		case 'A':
//			// _ctloc(weekday[tim_p->tm_wday]);
//			for (i = 0; i < ctloclen; i++) {
//				if (count < maxsize - 1)
//					s[count++] = ctloc[i];
//				else
//					return 0;
//			}
//			break;

//		case 'b':
//		case 'h':
//			// _ctloc(mon[tim_p->tm_mon]);
//			for (i = 0; i < ctloclen; i++) {
//				if (count < maxsize - 1)
//					s[count++] = ctloc[i];
//				else
//					return 0;
//			}
//			break;

//		case 'B':
//			_ctloc(month[tim_p->tm_mon]);
//			for (i = 0; i < ctloclen; i++) {
//				if (count < maxsize - 1)
//					s[count++] = ctloc[i];
//				else
//					return 0;
//			}
//			break;

//      case 'c':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'E' && *era_info && *_CurrentTimeLocale->era_d_t_fmt)
//          _ctloc(era_d_t_fmt);
//        else
//#endif /* _WANT_C99_TIME_FORMATS */
//          _ctloc(c_fmt);
//        goto recurse;

//      case 'r':
//        _ctloc(ampm_fmt);
//        goto recurse;

//      case 'x':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'E' && *era_info && *_CurrentTimeLocale->era_d_fmt)
//          _ctloc(era_d_fmt);
//        else
//#endif /* _WANT_C99_TIME_FORMATS */
//          _ctloc(x_fmt);
//        goto recurse;

//      case 'X':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'E' && *era_info && *_CurrentTimeLocale->era_t_fmt)
//          _ctloc(era_t_fmt);
//        else
//#endif /* _WANT_C99_TIME_FORMATS */
//          _ctloc(X_fmt);
//recurse:
//        if (*ctloc) {
//            /* Recurse to avoid need to replicate %Y formation. */
//            len = __strftime(&s[count], maxsize - count, ctloc, tim_p,
//                              locale, era_info, alt_digits);
//            if (len > 0)
//              count += len;
//            else
//              return 0;
//          }
//        break;

//	case 'C':
//        {
//          /* Examples of (tm_year + YEAR_BASE) that show how %Y == %C%y
//             with 32-bit int.
//             %Y               %C              %y
//             2147485547       21474855        47
//             10000            100             00
//             9999             99              99
//             0999             09              99
//             0099             00              99
//             0001             00              01
//             0000             00              00
//             -001             -0              01
//             -099             -0              99
//             -999             -9              99
//             -1000            -10             00
//             -10000           -100            00
//             -2147481748      -21474817       48
//
//             Be careful of both overflow and sign adjustment due to the
//             asymmetric range of years.
//          */
//#ifdef _WANT_C99_TIME_FORMATS
//          if (alt == 'E' && *era_info)
//            len = snprintf(&s[count], maxsize - count, "%s", (*era_info)->era_C);
//          else
//#endif /* _WANT_C99_TIME_FORMATS */
//            {
//              char *fmt = "%s%.*d";
//              char *pos = "";
//              int neg = tim_p->tm_year < -YEAR_BASE;
//              int century = tim_p->tm_year >= 0
//                ? tim_p->tm_year / 100 + YEAR_BASE / 100
//                : abs(tim_p->tm_year + YEAR_BASE) / 100;
//              if (pad) /* '0' or '+' */
//                {
//                  fmt = "%s%0.*d";
//                  if (century >= 100 && pad == '+')
//                    pos = "+";
//                }
//              if (width < 2)
//                width = 2;
//              len = snprintf(&s[count], maxsize - count, fmt,
//                              neg ? "-" : pos, width - neg, century);
//            }
//            CHECK_LENGTH();
//        }
//		break;

		case 'd':
		case 'e':
#ifdef _WANT_C99_TIME_FORMATS
			if (alt == 'O' && *alt_digits) {
				if (tim_p->tm_mday < 10) {
					if (*format_walker == 'd') {
						if (maxsize - count < 2)
							return 0;
						len = conv_to_alt_digits(&s[count], maxsize - count, 0, *alt_digits);
						CHECK_LENGTH();
					}
					if (*format_walker == 'e' || len == 0)
						s[count++] = ' ';
				}
				len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_mday, *alt_digits);
				CHECK_LENGTH();
				if (len > 0)
					break;
			}
#endif /* _WANT_C99_TIME_FORMATS */
			len = snprintf(&s[count], maxsize - count, *format_walker == 'd' ? "%.2d" : "%2d", tim_p->tm_mday);
			CHECK_LENGTH();
			break;

		case 'D':
			/* %m/%d/%y */
			len = snprintf(&s[count], maxsize - count,
				"%.2d/%.2d/%.2d",
				tim_p->tm_mon + 1, tim_p->tm_mday,
				tim_p->tm_year >= 0 ? tim_p->tm_year % 100
				: abs(tim_p->tm_year + YEAR_BASE) % 100);
			CHECK_LENGTH();
			break;

//	case 'F':
//        { /* %F is equivalent to "%+4Y-%m-%d", flags and width can change
//             that.  Recurse to avoid need to replicate %Y formation. */
//          char fmtbuf[32], *fmt = fmtbuf;
//
//          *fmt++ = '%';
//          if (pad) /* '0' or '+' */
//            *fmt++ = pad;
//          else
//            *fmt++ = '+';
//          if (!pad)
//            width = 10;
//          if (width < 6)
//            width = 6;
//          width -= 6;
//          if (width) {
//              len = snprintf(fmt, fmtbuf + 32 - fmt, "%lu", width);
//              if (len > 0)
//                fmt += len;
//            }
//          strcpy(fmt, "Y-%m-%d");
//          len = __strftime (&s[count], maxsize - count, fmtbuf, tim_p, locale, era_info, alt_digits);
//          if (len > 0)
//            count += len;
//          else
//            return 0;
//        }
//		break;

//	case 'g':
//        /* Be careful of both overflow and negative years, thanks to
//               the asymmetric range of years.  */
//        {
//          int adjust = iso_year_adjust(tim_p);
//          int year = tim_p->tm_year >= 0 ? tim_p->tm_year % 100
//              : abs(tim_p->tm_year + YEAR_BASE) % 100;
//          if (adjust < 0 && tim_p->tm_year <= -YEAR_BASE)
//              adjust = 1;
//          else if (adjust > 0 && tim_p->tm_year < -YEAR_BASE)
//              adjust = -1;
//          len = snprintf(&s[count], maxsize - count, "%.2d",
//                          ((year + adjust) % 100 + 100) % 100);
//            CHECK_LENGTH();
//        }
//		break;

//      case 'G':
//        {
//          /* See the comments for 'C' and 'Y'; this is a variable length
//             field.  Although there is no requirement for a minimum number
//             of digits, we use 4 for consistency with 'Y'.  */
//          int sign = tim_p->tm_year < -YEAR_BASE;
//          int adjust = iso_year_adjust(tim_p);
//          int century = tim_p->tm_year >= 0
//            ? tim_p->tm_year / 100 + YEAR_BASE / 100
//            : abs(tim_p->tm_year + YEAR_BASE) / 100;
//          int year = tim_p->tm_year >= 0 ? tim_p->tm_year % 100
//            : abs(tim_p->tm_year + YEAR_BASE) % 100;
//          if (adjust < 0 && tim_p->tm_year <= -YEAR_BASE)
//            sign = adjust = 1;
//          else if (adjust > 0 && sign)
//            adjust = -1;
//          year += adjust;
//          if (year == -1)
//            {
//              year = 99;
//              --century;
//            }
//          else if (year == 100)
//            {
//              year = 0;
//              ++century;
//            }
//          char fmtbuf[10], *fmt = fmtbuf;
//          /* int potentially overflows, so use unsigned instead.  */
//          unsigned p_year = century * 100 + year;
//          if (sign)
//            *fmt++ = '-';
//          else if (pad == '+' && p_year >= 10000)
//            {
//              *fmt++ = '+';
//              sign = 1;
//            }
//          if (width && sign)
//            --width;
//          *fmt++ = '%';
//          if (pad)
//            *fmt++ = '0';
//          strcpy(fmt, ".*u");
//          len = snprintf(&s[count], maxsize - count, fmtbuf, width, p_year);
//            if (len < 0  ||  (count+=len) >= maxsize)
//              return 0;
//        }
//          break;

		case 'H':
#ifdef _WANT_C99_TIME_FORMATS
			if (alt == 'O' && *alt_digits) {
				len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_hour, *alt_digits);
				CHECK_LENGTH();
				if (len > 0)
					break;
			}
#endif /* _WANT_C99_TIME_FORMATS */
			/*FALLTHRU*/

		case 'k':   /* newlib extension */
			len = snprintf(&s[count], maxsize - count,
				*format_walker == 'k' ? "%2d" : "%.2d",
				tim_p->tm_hour);
			CHECK_LENGTH();
			break;

//      case 'l':   /* newlib extension */
//        if (alt == 'O')
//          alt = '\0';
//        /*FALLTHRU*/

//      case 'I':
//        {
//          register int  h12;
//          h12 = (tim_p->tm_hour == 0 || tim_p->tm_hour == 12) ? 12 : tim_p->tm_hour % 12;
//#ifdef _WANT_C99_TIME_FORMATS
//          if (alt != 'O' || !*alt_digits
//              || !(len = conv_to_alt_digits(&s[count], maxsize - count, h12, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
//            len = snprintf(&s[count], maxsize - count, *format_walker == 'I' ? "%.2d" : "%2d", h12);
//          CHECK_LENGTH();
//        }
//        break;

		case 'j':
			len = snprintf(&s[count], maxsize - count, "%.3d", tim_p->tm_yday + 1);
			CHECK_LENGTH();
			break;

		case 'm':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt != 'O' || !*alt_digits
//            || !(len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_mon + 1, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
			len = snprintf(&s[count], maxsize - count, "%.2d", tim_p->tm_mon + 1);
			CHECK_LENGTH();
			break;

		case 'M':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt != 'O' || !*alt_digits
//            || !(len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_min, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
			len = snprintf(&s[count], maxsize - count, "%.2d", tim_p->tm_min);
			CHECK_LENGTH();
			break;

		case 'n':
			if (count < maxsize - 1)
				s[count++] = '\n';
			else
				return 0;
			break;

//      case 'p':
//      case 'P':
//        _ctloc(am_pm[tim_p->tm_hour < 12 ? 0 : 1]);
//        for (i = 0; i < ctloclen; i++) {
//            if (count < maxsize - 1)
//              s[count++] = (*format_walker == 'P' ? TOLOWER(ctloc[i])
//                                               : ctloc[i]);
//            else
//              return 0;
//          }
//        break;

		case 'R':
			len = snprintf(&s[count], maxsize - count, "%.2d:%.2d",
				tim_p->tm_hour, tim_p->tm_min);
			CHECK_LENGTH();
			break;

//      case 's':
///*
// * From:
// * The Open Group Base Specifications Issue 7
// * IEEE Std 1003.1, 2013 Edition
// * Copyright (c) 2001-2013 The IEEE and The Open Group
// * XBD Base Definitions
// * 4. General Concepts
// * 4.15 Seconds Since the Epoch
// * A value that approximates the number of seconds that have elapsed since the
// * Epoch. A Coordinated Universal Time name (specified in terms of seconds
// * (tm_sec), minutes (tm_min), hours (tm_hour), days since January 1 of the year
// * (tm_yday), and calendar year minus 1900 (tm_year)) is related to a time
// * represented as seconds since the Epoch, according to the expression below.
// * If the year is <1970 or the value is negative, the relationship is undefined.
// * If the year is >=1970 and the value is non-negative, the value is related to a
// * Coordinated Universal Time name according to the C-language expression, where
// * tm_sec, tm_min, tm_hour, tm_yday, and tm_year are all integer types:
// * tm_sec + tm_min*60 + tm_hour*3600 + tm_yday*86400 +
// *     (tm_year-70)*31536000 + ((tm_year-69)/4)*86400 -
// *     ((tm_year-1)/100)*86400 + ((tm_year+299)/400)*86400
// * OR
// * ((((tm_year-69)/4 - (tm_year-1)/100 + (tm_year+299)/400 +
// *         (tm_year-70)*365 + tm_yday)*24 + tm_hour)*60 + tm_min)*60 + tm_sec
// */
///* modified from %z case by hoisting offset outside if block and initializing */
//        {
//          long offset = 0;    /* offset < 0 => W of GMT, > 0 => E of GMT:
//                                 subtract to get UTC */
//
//          if (tim_p->tm_isdst >= 0)
//            {
//              TZ_LOCK;
//              if (!tzset_called)
//                {
//                  _tzset_unlocked ();
//                  tzset_called = 1;
//                }
//
//#if defined (__CYGWIN__)
//              /* Cygwin must check if the application has been built with or
//                 without the extra tm members for backward compatibility, and
//                 then use either that or the old method fetching from tzinfo.
//                 Rather than pulling in the version check infrastructure, we
//                 just call a Cygwin function. */
//              extern long __cygwin_gettzoffset (const struct tm *tmp);
//              offset = __cygwin_gettzoffset (tim_p);
//#elif defined (__TM_GMTOFF)
//              offset = tim_p->__TM_GMTOFF;
//#else
//              __tzinfo_type *tz = __gettzinfo ();
//              /* The sign of this is exactly opposite the envvar TZ.  We
//                 could directly use the global _timezone for tm_isdst==0,
//                 but have to use __tzrule for daylight savings.  */
//              offset = -tz->__tzrule[tim_p->tm_isdst > 0].offset;
//#endif
//              TZ_UNLOCK;
//            }
//          len = snprintf(&s[count], maxsize - count, "%lld",
//                          (((((long long)tim_p->tm_year - 69)/4
//                              - (tim_p->tm_year - 1)/100
//                              + (tim_p->tm_year + 299)/400
//                              + (tim_p->tm_year - 70)*365 + tim_p->tm_yday)*24
//                            + tim_p->tm_hour)*60 + tim_p->tm_min)*60
//                          + tim_p->tm_sec - offset);
//          CHECK_LENGTH();
//        }
//          break;

		case 'S':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt != 'O' || !*alt_digits
//            || !(len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_sec, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
			len = snprintf(&s[count], maxsize - count, "%.2d", tim_p->tm_sec);
			CHECK_LENGTH();
			break;

		case 't':
			if (count < maxsize - 1)
				s[count++] = '\t';
			else
				return 0;
			break;

		case 'T':
			len = snprintf(&s[count], maxsize - count, "%.2d:%.2d:%.2d",
				tim_p->tm_hour, tim_p->tm_min, tim_p->tm_sec);
			CHECK_LENGTH();
			break;

//      case 'u':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'O' && *alt_digits) {
//            len = conv_to_alt_digits(&s[count], maxsize - count,
//                                      tim_p->tm_wday == 0 ? 7 : tim_p->tm_wday,
//                                      *alt_digits);
//            CHECK_LENGTH();
//            if (len > 0)
//              break;
//          }
//#endif /* _WANT_C99_TIME_FORMATS */
//          if (count < maxsize - 1) {
//              if (tim_p->tm_wday == 0)
//                s[count++] = '7';
//              else
//                s[count++] = '0' + tim_p->tm_wday;
//            }
//          else
//            return 0;
//          break;

//      case 'U':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt != 'O' || !*alt_digits
//            || !(len = conv_to_alt_digits(&s[count], maxsize - count,
//                                           (tim_p->tm_yday + 7 -
//                                            tim_p->tm_wday) / 7,
//                                           *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
//          len = snprintf(&s[count], maxsize - count, "%.2d",
//                       (tim_p->tm_yday + 7 -
//                        tim_p->tm_wday) / 7);
//          CHECK_LENGTH();
//        break;

//      case 'V':
//        {
//          int adjust = iso_year_adjust (tim_p);
//          int wday = (tim_p->tm_wday) ? tim_p->tm_wday - 1 : 6;
//          int week = (tim_p->tm_yday + 10 - wday) / 7;
//          if (adjust > 0)
//              week = 1;
//          else if (adjust < 0)
//              /* Previous year has 53 weeks if current year starts on
//                 Fri, and also if current year starts on Sat and
//                 previous year was leap year.  */
//              week = 52 + (4 >= (wday - tim_p->tm_yday
//                                 - isleap (tim_p->tm_year
//                                           + (YEAR_BASE - 1
//                                              - (tim_p->tm_year < 0
//                                                 ? 0 : 2000)))));
//#ifdef _WANT_C99_TIME_FORMATS
//          if (alt != 'O' || !*alt_digits
//              || !(len = conv_to_alt_digits(&s[count], maxsize - count, week, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
//            len = snprintf(&s[count], maxsize - count, "%.2d", week);
//            CHECK_LENGTH();
//        }
//          break;

//      case 'w':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'O' && *alt_digits)
//          {
//            len = conv_to_alt_digits(&s[count], maxsize - count, tim_p->tm_wday, *alt_digits);
//            CHECK_LENGTH();
//            if (len > 0)
//              break;
//          }
//#endif /* _WANT_C99_TIME_FORMATS */
//        if (count < maxsize - 1)
//            s[count++] = '0' + tim_p->tm_wday;
//        else
//          return 0;
//        break;

//      case 'W':
//        {
//          int wday = (tim_p->tm_wday) ? tim_p->tm_wday - 1 : 6;
//          wday = (tim_p->tm_yday + 7 - wday) / 7;
//#ifdef _WANT_C99_TIME_FORMATS
//          if (alt != 'O' || !*alt_digits
//              || !(len = conv_to_alt_digits(&s[count], maxsize - count, wday, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
//            len = snprintf(&s[count], maxsize - count, "%.2d", wday);
//            CHECK_LENGTH();
//        }
//        break;

//      case 'y':
//          {
//#ifdef _WANT_C99_TIME_FORMATS
//            if (alt == 'E' && *era_info)
//              len = snprintf(&s[count], maxsize - count, "%d",
//                              (*era_info)->year);
//            else
//#endif /* _WANT_C99_TIME_FORMATS */
//              {
//                /* Be careful of both overflow and negative years, thanks to
//                   the asymmetric range of years.  */
//                int year = tim_p->tm_year >= 0 ? tim_p->tm_year % 100
//                           : abs(tim_p->tm_year + YEAR_BASE) % 100;
//#ifdef _WANT_C99_TIME_FORMATS
//                if (alt != 'O' || !*alt_digits
//                    || !(len = conv_to_alt_digits(&s[count], maxsize - count, year, *alt_digits)))
//#endif /* _WANT_C99_TIME_FORMATS */
//                  len = snprintf(&s[count], maxsize - count, "%.2d",
//                                  year);
//              }
//              CHECK_LENGTH();
//          }
//        break;

		case 'Y':
//#ifdef _WANT_C99_TIME_FORMATS
//        if (alt == 'E' && *era_info) {
//            ctloc = (*era_info)->era_Y;
//            goto recurse;
//          }
//        else
//#endif /* _WANT_C99_TIME_FORMATS */
		{
			char fmtbuf[10], *fmt = fmtbuf;
			int sign = tim_p->tm_year < -YEAR_BASE;
			/* int potentially overflows, so use unsigned instead.  */
			register unsigned year = (unsigned) tim_p->tm_year + (unsigned) YEAR_BASE;
			if (sign) {
				*fmt++ = '-';
				year = UINT_MAX - year + 1;
			} else if (pad == '+' && year >= 10000) {
				*fmt++ = '+';
				sign = 1;
			}
			if (width && sign)
				--width;
			*fmt++ = '%';
			if (pad)
				*fmt++ = '0';
			strcpy(fmt, ".*u");
			len = snprintf(&s[count], maxsize - count, fmtbuf, width, year);
			CHECK_LENGTH();
		}
			break;

//      case 'z':
//          if (tim_p->tm_isdst >= 0)
//            {
//            long offset;
//
//            TZ_LOCK;
//            if (!tzset_called)
//              {
//                _tzset_unlocked ();
//                tzset_called = 1;
//              }
//
//#if defined (__CYGWIN__)
//            /* Cygwin must check if the application has been built with or
//               without the extra tm members for backward compatibility, and
//               then use either that or the old method fetching from tzinfo.
//               Rather than pulling in the version check infrastructure, we
//               just call a Cygwin function. */
//            extern long __cygwin_gettzoffset (const struct tm *tmp);
//            offset = __cygwin_gettzoffset (tim_p);
//#elif defined (__TM_GMTOFF)
//            offset = tim_p->__TM_GMTOFF;
//#else
//            __tzinfo_type *tz = __gettzinfo ();
//            /* The sign of this is exactly opposite the envvar TZ.  We
//               could directly use the global _timezone for tm_isdst==0,
//               but have to use __tzrule for daylight savings.  */
//            offset = -tz->__tzrule[tim_p->tm_isdst > 0].offset;
//#endif
//            TZ_UNLOCK;
//            len = snprintf(&s[count], maxsize - count, "%+03ld%.2ld",
//                            offset / SECSPERHOUR,
//                            labs(offset / SECSPERMIN) % 60L);
//              CHECK_LENGTH();
//            }
//          break;

//      case 'Z':
//        if (tim_p->tm_isdst >= 0)
//          {
//            size_t size;
//            const char *tznam = NULL;
//
//            TZ_LOCK;
//            if (!tzset_called)
//              {
//                _tzset_unlocked ();
//                tzset_called = 1;
//              }
//#if defined (__CYGWIN__)
//            /* See above. */
//            extern const char *__cygwin_gettzname (const struct tm *tmp);
//            tznam = __cygwin_gettzname (tim_p);
//#elif defined (__TM_ZONE)
//            tznam = tim_p->__TM_ZONE;
//#endif
//            if (!tznam)
//              tznam = _tzname[tim_p->tm_isdst > 0];
//            /* Note that in case of wcsftime this loop only works for
//               timezone abbreviations using the portable codeset (aka ASCII).
//               This seems to be the case, but if that ever changes, this
//               loop needs revisiting. */
//            size = strlen (tznam);
//            for (i = 0; i < size; i++)
//              {
//                if (count < maxsize - 1)
//                  s[count++] = tznam[i];
//                else
//                  {
//                    TZ_UNLOCK;
//                    return 0;
//                  }
//              }
//            TZ_UNLOCK;
//          }
//        break;

		case '%':
			if (count < maxsize - 1)
				s[count++] = '%';
			else
				return 0;
			break;

		default:
			return 0;
		}

		if (*format_walker)
			format_walker++;
		else
			break;
	}
	if (maxsize)
		s[count] = '\0';

	return count;
}

// ================================================================
//#include "newlib.h"
//#include <sys/config.h>
//#include <stddef.h>
//#include <stdio.h>
//#include <time.h>
//#include <string.h>
//#include <stdlib.h>
//#include <limits.h>
//#include <ctype.h>
//#include <wctype.h>
//#include "local.h"
//#include "setlocale.h"

//#define _ctloc(x) (ctloclen = strlen (ctloc = _CurrentTimeLocale->x), ctloc)
//#define TOLOWER(c)	tolower((int)(unsigned char)(c))


//static _CONST int dname_len[7] = /{6, 6, 7, 9, 8, 6, 8};

//// Using the tm_year, tm_wday, and tm_yday components of TIM_P, return
//// -1, 0, or 1 as the adjustment to add to the year for the ISO week
//// numbering used in "%g%G%V", avoiding overflow.
//static int iso_year_adjust(const struct tm *tim_p) {
//	/* Account for fact that tm_year==0 is year 1900.  */
//	int leap = isleap (tim_p->tm_year + (YEAR_BASE - (tim_p->tm_year < 0 ? 0 : 2000)));
//
//  // Pack the yday, wday, and leap year into a single int since there are so many disparate cases.
//#define PACK(yd, wd, lp) (((yd) << 4) + (wd << 1) + (lp))
//
//  switch (PACK (tim_p->tm_yday, tim_p->tm_wday, leap)) {
//    case PACK (0, 5, 0): /* Jan 1 is Fri, not leap.  */
//    case PACK (0, 6, 0): /* Jan 1 is Sat, not leap.  */
//    case PACK (0, 0, 0): /* Jan 1 is Sun, not leap.  */
//    case PACK (0, 5, 1): /* Jan 1 is Fri, leap year.  */
//    case PACK (0, 6, 1): /* Jan 1 is Sat, leap year.  */
//    case PACK (0, 0, 1): /* Jan 1 is Sun, leap year.  */
//    case PACK (1, 6, 0): /* Jan 2 is Sat, not leap.  */
//    case PACK (1, 0, 0): /* Jan 2 is Sun, not leap.  */
//    case PACK (1, 6, 1): /* Jan 2 is Sat, leap year.  */
//    case PACK (1, 0, 1): /* Jan 2 is Sun, leap year.  */
//    case PACK (2, 0, 0): /* Jan 3 is Sun, not leap.  */
//    case PACK (2, 0, 1): /* Jan 3 is Sun, leap year.  */
//      return -1; /* Belongs to last week of previous year.  */
//    case PACK (362, 1, 0): /* Dec 29 is Mon, not leap.  */
//    case PACK (363, 1, 1): /* Dec 29 is Mon, leap year.  */
//    case PACK (363, 1, 0): /* Dec 30 is Mon, not leap.  */
//    case PACK (363, 2, 0): /* Dec 30 is Tue, not leap.  */
//    case PACK (364, 1, 1): /* Dec 30 is Mon, leap year.  */
//    case PACK (364, 2, 1): /* Dec 30 is Tue, leap year.  */
//    case PACK (364, 1, 0): /* Dec 31 is Mon, not leap.  */
//    case PACK (364, 2, 0): /* Dec 31 is Tue, not leap.  */
//    case PACK (364, 3, 0): /* Dec 31 is Wed, not leap.  */
//    case PACK (365, 1, 1): /* Dec 31 is Mon, leap year.  */
//    case PACK (365, 2, 1): /* Dec 31 is Tue, leap year.  */
//    case PACK (365, 3, 1): /* Dec 31 is Wed, leap year.  */
//      return 1; /* Belongs to first week of next year.  */
//    }
//  return 0; /* Belongs to specified year.  */
//#undef PACK
//}

// ================================================================
//#ifdef _WANT_C99_TIME_FORMATS
//typedef struct _era_info_t {
//	int   year;
//	char *era_C;
//	char *era_Y;
//} era_info_t;

//static era_info_t * get_era_info (const struct tm *tim_p, const char *era) {
//	char *c;
//	const char *dir;
//	long offset;
//	struct tm stm, etm;
//	era_info_t *ei;
//
//	ei = (era_info_t *) calloc (1, sizeof (era_info_t));
//	if (!ei)
//		return NULL;
//
//	stm.tm_isdst = etm.tm_isdst = 0;
//	while (era) {
//
//		dir = era;
//		era += 2;
//		offset = strtol(era, &c, 10);
//		era = c + 1;
//		stm.tm_year = strtol(era, &c, 10) - YEAR_BASE;
//		// Adjust offset for negative gregorian dates.
//		if (stm.tm_year <= -YEAR_BASE)
//			++stm.tm_year;
//		stm.tm_mon = strtol(c + 1, &c, 10) - 1;
//		stm.tm_mday = strtol(c + 1, &c, 10);
//		stm.tm_hour = stm.tm_min = stm.tm_sec = 0;
//		era = c + 1;
//		if (era[0] == '-' && era[1] == '*') {
//			etm = stm;
//			stm.tm_year = INT_MIN;
//			stm.tm_mon = stm.tm_mday = stm.tm_hour = stm.tm_min = stm.tm_sec = 0;
//			era += 3;
//
//		} else if (era[0] == '+' && era[1] == '*') {
//			etm.tm_year = INT_MAX;
//			etm.tm_mon = 11;
//			etm.tm_mday = 31;
//			etm.tm_hour = 23;
//			etm.tm_min = etm.tm_sec = 59;
//			era += 3;
//
//		} else {
//			etm.tm_year = strtol(era, &c, 10) - YEAR_BASE;
//			// Adjust offset for negative gregorian dates.
//			if (etm.tm_year <= -YEAR_BASE)
//				++etm.tm_year;
//			etm.tm_mon = strtol(c + 1, &c, 10) - 1;
//			etm.tm_mday = strtol(c + 1, &c, 10);
//			etm.tm_mday = 31;
//			etm.tm_hour = 23;
//			etm.tm_min = etm.tm_sec = 59;
//			era = c + 1;
//		}
//		if ((tim_p->tm_year > stm.tm_year
//			|| (tim_p->tm_year == stm.tm_year
//					&& (tim_p->tm_mon > stm.tm_mon
//							|| (tim_p->tm_mon == stm.tm_mon
//								&& tim_p->tm_mday >= stm.tm_mday))))
//			&& (tim_p->tm_year < etm.tm_year
//					|| (tim_p->tm_year == etm.tm_year
//							&& (tim_p->tm_mon < etm.tm_mon
//									|| (tim_p->tm_mon == etm.tm_mon
//											&& tim_p->tm_mday <= etm.tm_mday)))))
//		{
//			/* Gotcha */
//			size_t len;
//
//			/* year */
//			if (*dir == '+' && stm.tm_year != INT_MIN)
//				ei->year = tim_p->tm_year - stm.tm_year + offset;
//			else
//				ei->year = etm.tm_year - tim_p->tm_year + offset;
//			/* era_C */
//			c = strchr(era, ':');
//			len = c - era;
//			ei->era_C = (char *) malloc ((len + 1) * sizeof (char));
//			if (!ei->era_C) {
//				free (ei);
//				return NULL;
//			}
//			strncpy(ei->era_C, era, len);
//			era += len;
//			ei->era_C[len] = '\0';
//			/* era_Y */
//			++era;
//			c = strchr(era, ';');
//			if (!c)
//				c = strchr(era, '\0');
//			len = c - era;
//			ei->era_Y = (char *) malloc ((len + 1) * sizeof (char));
//			if (!ei->era_Y) {
//				free (ei->era_C);
//				free (ei);
//				return NULL;
//			}
//			strncpy(ei->era_Y, era, len);
//			era += len;
//			ei->era_Y[len] = '\0';
//			return ei;
//		}
//		else
//			era = strchr(era, ';');
//		if (era)
//			++era;
//	}
//	return NULL;
//}

//static void era_info_free(era_info_t *ei) {
//  free(ei->era_C);
//  free(ei->era_Y);
//  free(ei);
//}

// ================================================================
//typedef struct _alt_digits_t {
//	size_t num;
//	char** digit;
//	char*  buffer;
//} alt_digits_t;

// ----------------------------------------------------------------
//static alt_digits_t* alt_digits_alloc(const char *alt_digits) {
//	alt_digits_t *adi;
//	const char *a, *e;
//	char *aa, *ae;
//	size_t len;
//
//	adi = (alt_digits_t *) calloc (1, sizeof (alt_digits_t));
//	if (!adi)
//		return NULL;
//
//	/* Compute number of alt_digits. */
//	adi->num = 1;
//	for (a = alt_digits; (e = strchr (a, ';')) != NULL; a = e + 1)
//		++adi->num;
//	// Allocate the `digit' array, which is an array of `num' pointers into `buffer'.
//	adi->digit = (char **) calloc (adi->num, sizeof (char *));
//	if (!adi->digit) {
//		free (adi);
//		return NULL;
//	}
//	// Compute memory required for `buffer'.
//	len = strlen(alt_digits);
//	// Allocate it.
//	adi->buffer = (char *) malloc ((len + 1) * sizeof (char));
//	if (!adi->buffer) {
//		free (adi->digit);
//		free (adi);
//		return NULL;
//	}
//	// Store digits in it.
//	strcpy(adi->buffer, alt_digits);
//	/* Store the pointers into `buffer' into the appropriate `digit' slot. */
//	for (len = 0, aa = adi->buffer; (ae = strchr(aa, ';')) != NULL; ++len, aa = ae + 1) {
//		*ae = '\0';
//		adi->digit[len] = aa;
//	}
//	adi->digit[len] = aa;
//	return adi;
//}

// ----------------------------------------------------------------
//static void alt_digits_free(alt_digits_t *adi) {
//	free (adi->digit);
//	free (adi->buffer);
//	free (adi);
//}

// ----------------------------------------------------------------
//// Return 0 if no alt_digit is available for a number.
//// Return -1 if buffer size isn't sufficient to hold alternative digit.
//// Return length of new digit otherwise. */
//static int conv_to_alt_digits(char *buf, size_t bufsiz, unsigned num, alt_digits_t *adi) {
//	if (num < adi->num) {
//		size_t len = strlen(adi->digit[num]);
//		if (bufsiz < len)
//			return -1;
//		strcpy(buf, adi->digit[num]);
//		return (int) len;
//	}
//	return 0;
//}

// ================================================================
//size_t strftime_l(
//	char *__restrict s,
//	size_t maxsize,
//	const char *__restrict format,
//	const struct tm *__restrict tim_p,
//	struct __locale_t *locale)
//{
//#ifdef _WANT_C99_TIME_FORMATS
//  era_info_t *era_info = NULL;
//  alt_digits_t *alt_digits = NULL;
//  size_t ret = __strftime (s, maxsize, format, tim_p, locale,
//			   &era_info, &alt_digits);
//  if (era_info)
//    era_info_free(era_info);
//  if (alt_digits)
//    alt_digits_free(alt_digits);
//  return ret;
//#else /* !_WANT_C99_TIME_FORMATS */
//  return __strftime(s, maxsize, format, tim_p, locale, NULL, NULL);
//#endif /* !_WANT_C99_TIME_FORMATS */
//}
