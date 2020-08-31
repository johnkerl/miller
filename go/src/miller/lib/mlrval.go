package lib

// Two kinds of null: absent (key not present in a record) and void (key
// present with empty value).  // Note void is an acceptable string (empty
// string) but not an acceptable number.

// Void-valued mlrvals have u.strv = "".

// #define MT_ERROR    0 // E.g. error encountered in one eval & it propagates up the AST.
// #define MT_ABSENT   1 // No such key, e.g. $z in 'x=,y=2'
// #define MT_EMPTY    2 // Empty value, e.g. $x in 'x=,y=2'
// #define MT_STRING   3
// #define MT_INT      4
// #define MT_FLOAT    5
// #define MT_BOOLEAN  6
// #define MT_DIM      7

// typedef struct _mv_t {
// 	union {
// 		char*      strv;  // MT_STRING and MT_EMPTY
// 		long long  intv;  // MT_INT, and == 0 for MT_ABSENT and MT_ERROR
// 		double     fltv;  // MT_FLOAT
// 		int        boolv; // MT_BOOLEAN
// 	} u;
// 	unsigned char type;
// 	char free_flags;
// } mv_t;

type Mlrval struct {
}
