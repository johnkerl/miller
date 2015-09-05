#ifndef MLR_TEST_UTIL_H
#define MLR_TEST_UTIL_H

// Returns the name of the temp file
char* write_temp_file_or_die(char* contents);

void unlink_file_or_die(char* path);

#endif // MLR_TEST_UTIL_H
