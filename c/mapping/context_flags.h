#ifndef CONTEXT_FLAGS_H
#define CONTEXT_FLAGS_H

// The grammar permits certain statements which are syntactically invalid, (a) because it's awkward to handle
// there, and (b) because we get far better control over error messages here (vs. 'syntax error').
// The following flags are used as the CST is built from the AST for CST-build-time validation.

#define IN_BREAKABLE        0x00000200 // break/continue are only OK (recursively) inside for/while/do-while
#define IN_BEGIN_OR_END     0x00000400 // $stuff is not OK (recursively) inside begin/end
#define IN_FUNC_DEF         0x00000800 // local only valid in func/subr; no srec assignments in functions
#define IN_SUBR_DEF         0x00001000 // local only valid in func/subr
#define IN_MLR_FILTER       0x00002000 // Anywhere within mlr filter, the 'filter' keyword is invalid
#define IN_MLR_FINAL_FILTER 0x00004000 // mlr filter's final statement must be a bare boolean

#endif // CONTEXT_FLAGS_H
