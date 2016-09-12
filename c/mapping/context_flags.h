#ifndef CONTEXT_FLAGS_H
#define CONTEXT_FLAGS_H

// The grammar permits certain statements which are syntactically invalid, (a) because it's awkward to handle
// there, and (b) because we get far better control over error messages here (vs. 'syntax error').
// The following flags are used as the CST is built from the AST for CST-build-time validation.

#define IN_BINDABLE     0x00000100 // boundvars are only OK inside a bindable, e.g. (recursively) inside a for-loop
#define IN_BREAKABLE    0x00000200 // break/continue are only OK (recursively) inside for/while/do-while
#define IN_BEGIN_OR_END 0x00000400 // $stuff is not OK (recursively) inside begin/end
#define IN_FUNC_DEF     0x00000800 // $stuff, @stuff is not OK inside user-defined functions
#define IN_MLR_FILTER   0x00004000 // mlr filter takes only a single boolean; no @-vars; no looping; etc.

#endif // CONTEXT_FLAGS_H
