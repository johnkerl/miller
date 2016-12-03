#ifndef LEMON_REPORT_H
#define LEMON_REPORT_H

#include "lemon_structs.h"

void Reprint(struct lemon *);
void ReportOutput(struct lemon *);
void ReportTable(struct lemon *, int mhflag, int suppress_line_directives);
void ReportHeader(struct lemon *);
void CompressTables(struct lemon *);

#endif // LEMON_REPORT_H
