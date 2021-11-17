#ifndef LEMON_ASSERT_H
#define LEMON_ASSERT_H

void lemon_assert(char *file, int line);
#ifndef MLR_DSL_NDEBUG
#  define assert(X) if(!(X))lemon_assert(__FILE__,__LINE__)
#else
#  define assert(X)
#endif

#endif // LEMON_ASSERT_H
