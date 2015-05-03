// A C-program for MT19937, with initialization improved 2002/1/26.
// Coded by Takuji Nishimura and Makoto Matsumoto.
//
// Before using, initialize the state by using mtrand_init(seed)
// or mtrand_init_from_array(init_key, key_length).
//
// Copyright (C) 1997 - 2002, Makoto Matsumoto and Takuji Nishimura,
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   1. Redistributions of source code must retain the above copyright
//      notice, this list of conditions and the following disclaimer.
//
//   2. Redistributions in binary form must reproduce the above copyright
//      notice, this list of conditions and the following disclaimer in the
//      documentation and/or other materials provided with the distribution.
//
//   3. The names of its contributors may not be used to endorse or promote
//      products derived from this software without specific prior written
//      permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
//
// Any feedback is very welcome.
// http://www.math.sci.hiroshima-u.ac.jp/~m-mat/MT/emt.html
// email: m-mat @ math.sci.hiroshima-u.ac.jp (remove space)
// Real versions are due to Isaku Wada, added 2002-01-09.

#include <stdio.h>
#include <unistd.h>
#include <sys/time.h>
#include "lib/mtrand.h"

// ----------------------------------------------------------------
// Period parameters
#define N 624
#define M 397
#define MATRIX_A   0x9908b0df   // constant vector a
#define UPPER_MASK 0x80000000 // most significant w-r bits
#define LOWER_MASK 0x7fffffff // least significant r bits

static unsigned mt[N];     // the array for the state vector
static int mti=N+1;             // mti==N+1 means mt[N] is not initialized

// ----------------------------------------------------------------
void mtrand_init_default()
{
	struct timeval tv;
	(void)gettimeofday(&tv, NULL);
	mtrand_init((unsigned)tv.tv_sec ^ (unsigned)tv.tv_usec ^ (unsigned)getpid());
}

// ----------------------------------------------------------------
// Initializes mt[N] with a seed.
void mtrand_init(unsigned s)
{
	mt[0]= s & 0xffffffff;
	for (mti=1; mti<N; mti++) {
		mt[mti] =
	    (1812433253 * (mt[mti-1] ^ (mt[mti-1] >> 30)) + mti);
		// See Knuth TAOCP Vol2. 3rd Ed. P.106 for multiplier.  In the previous
		// versions, MSBs of the seed affect only MSBs of the array mt[].
		// 2002/01/09 modified by Makoto Matsumoto
		mt[mti] &= 0xffffffff; // for >32 bit machines
	}
}

// ----------------------------------------------------------------
// Initialize by an array with array-length:
// init_key is the array for initializing keys;
// key_length is its length.
void mtrand_init_from_array(unsigned init_key[], int key_length)
{
	int i, j, k;
	mtrand_init(19650218);
	i=1; j=0;
	k = (N>key_length ? N : key_length);
	for (; k; k--) {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1664525))
		  + init_key[j] + j; // non linear
		mt[i] &= 0xffffffff; // for WORDSIZE > 32 machines
		i++; j++;
		if (i>=N) { mt[0] = mt[N-1]; i=1; }
		if (j>=key_length) j=0;
	}
	for (k=N-1; k; k--) {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 30)) * 1566083941))
		  - i; // non linear
		mt[i] &= 0xffffffff; // for WORDSIZE > 32 machines
		i++;
		if (i>=N) { mt[0] = mt[N-1]; i=1; }
	}

	mt[0] = 0x80000000; // MSB is 1, ensuring non-zero initial array
}

// ----------------------------------------------------------------
// Generates a uniformly distributed 32-bit integer.
unsigned get_mtrand_int32(void)
{
	unsigned y;
	static unsigned mag01[2]={0x0, MATRIX_A};
	// mag01[x] = x * MATRIX_A  for x=0,1

	if (mti >= N) { // Generate N words at one time
		int kk;

		if (mti == N+1)   // If mtrand_init() has not been called,
			mtrand_init(5489); // a default initial seed is used.

		for (kk=0;kk<N-M;kk++) {
			y = (mt[kk]&UPPER_MASK)|(mt[kk+1]&LOWER_MASK);
			mt[kk] = mt[kk+M] ^ (y >> 1) ^ mag01[y & 1];
		}
		for (;kk<N-1;kk++) {
			y = (mt[kk]&UPPER_MASK)|(mt[kk+1]&LOWER_MASK);
			mt[kk] = mt[kk+(M-N)] ^ (y >> 1) ^ mag01[y & 1];
		}
		y = (mt[N-1]&UPPER_MASK)|(mt[0]&LOWER_MASK);
		mt[N-1] = mt[M-1] ^ (y >> 1) ^ mag01[y & 1];

		mti = 0;
	}

	y = mt[mti++];

	// Tempering
	y ^= (y >> 11);
	y ^= (y << 7) & 0x9d2c5680;
	y ^= (y << 15) & 0xefc60000;
	y ^= (y >> 18);

	return y;
}

// ----------------------------------------------------------------
// Generates a uniformly distributed 31-bit integer.
int get_mtrand_int31(void)
{
	return (int)(get_mtrand_int32()>>1);
}

// ----------------------------------------------------------------
// Generates a random number on [0,1)-real-interval.
double get_mtrand_float(void)
{
	return get_mtrand_int32()*(1.0/4294967296.0);
	// divided by 2^32
}

// ----------------------------------------------------------------
// Generates a random number on [0,1) with 53-bit resolution.
double get_mtrand_double(void)
{
	unsigned a = get_mtrand_int32() >> 5;
	unsigned b = get_mtrand_int32() >> 6;
	return (a*67108864.0+b) * (1.0/9007199254740992.0);
}
