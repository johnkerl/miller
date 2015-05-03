#ifdef __GEN_MAIN__
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
int main(void) {
    char* names[] = { "hat", "pan", "eks", "wye", "zee" };
    int num_names = sizeof(names) / sizeof(names[0]);
    srand(time(0) ^ getpid());
	for (int i = 0; ; i++) {
        int ai = rand() % num_names;
        int bi = rand() % num_names;
        char* a = names[ai];
        char* b = names[bi];
        double x = (double)rand() / (double)RAND_MAX;
        double y = (double)rand() / (double)RAND_MAX;
		printf("a=%s,b=%s,i=%d,x=%.18lf,y=%.18lf\n", a, b, i, x, y);
    }
	return 0;
}
#endif // __GEN_MAIN__

// a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
// a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
// a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
// a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
