#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "mlrmath.h"

#define JACOBI_TOLERANCE 1e-12
#define JACOBI_MAXITER   20

static void  matmul2(double C[2][2], double A[2][2], double B[2][2]);
static void matmul2t(double C[2][2], double A[2][2], double B[2][2]);

// ----------------------------------------------------------------
// Jacobi real-symmetric eigensolver. Loosely adapted from Numerical Recipes.
//
// Note: this is coded for n=2 (to implement PCA linear regression on 2
// variables) but the algorithm is quite general. Changing from 2 to n is a
// matter of updating the top and bottom of the function: function signature to
// take double** matrix, double* eigenvector_1, double* eigenvector_2, and n;
// create copy-matrix and make-identity matrix functions; free temp matrices at
// the end; etc.

void mlr_get_real_symmetric_eigensystem(
	double matrix[2][2],      // Input
	double *peigenvalue_1,    // Output: dominant eigenvalue
	double *peigenvalue_2,    // Output: less-dominant eigenvalue
	double eigenvector_1[2],  // Output: corresponding to dominant eigenvalue
	double eigenvector_2[2])  // Output: corresponding to less-dominant eigenvalue
{
	double L[2][2] = {
		{ matrix[0][0], matrix[0][1] },
		{ matrix[1][0], matrix[1][1] }
	};
	double V[2][2] = {
		{ 1.0, 0.0 },
		{ 0.0, 1.0 },
	};
	double P[2][2], PT_A[2][2];
	int n = 2;

	int found = 0;
	for (int iter = 0; iter < JACOBI_MAXITER; iter++) {
		double sum = 0.0;
		for (int i = 1; i < n; i++)
			for (int j = 0; j < i; j++)
				sum += fabs(L[i][j]);
		if (fabs(sum*sum) < JACOBI_TOLERANCE) {
			found = 1;
			break;
		}

		for (int p = 0; p < n; p++) {
			for (int q = p+1; q < n; q++) {
				double numer = L[p][p] - L[q][q];
				double denom = L[p][q] + L[q][p];
				if (fabs(denom) < JACOBI_TOLERANCE)
					continue;
				double theta = numer / denom;
				int sign_theta = (theta < 0) ? -1 : 1;
				double t = sign_theta / (fabs(theta) + sqrt(theta*theta + 1));
				double c = 1.0 / sqrt(t*t + 1);
				double s = t * c;

				for (int pi = 0; pi < n; pi++)
					for (int pj = 0; pj < n; pj++)
						P[pi][pj] = (pi == pj) ? 1.0 : 0.0;
				P[p][p] =  c;
				P[p][q] = -s;
				P[q][p] =  s;
				P[q][q] =  c;

				// L = P.transpose() * L * P
				// V = V * P
				matmul2t(PT_A, P, L);
				matmul2(L, PT_A, P);
				matmul2(V, V, P);
			}
		}
	}

	if (!found) {
		fprintf(stderr,
			"Jacobi eigensolver: max iterations (%d) exceeded.  Non-symmetric input?\n", JACOBI_MAXITER);
			exit(1);
	}

	double eigenvalue_1 = L[0][0];
	double eigenvalue_2 = L[1][1];
	double abs1 = fabs(eigenvalue_1);
	double abs2 = fabs(eigenvalue_2);
	if (abs1 > abs2) {
		*peigenvalue_1 = eigenvalue_1;
		*peigenvalue_2 = eigenvalue_2;
		eigenvector_1[0] = V[0][0]; // Column 0 of V
		eigenvector_1[1] = V[1][0];
		eigenvector_2[0] = V[0][1]; // Column 1 of V
		eigenvector_2[1] = V[1][1];
	} else {
		*peigenvalue_1 = eigenvalue_2;
		*peigenvalue_2 = eigenvalue_1;
		eigenvector_1[0] = V[0][1];
		eigenvector_1[1] = V[1][1];
		eigenvector_2[0] = V[0][0];
		eigenvector_2[1] = V[1][0];
	}
}

// xxx cmt mem-mgmt
static void matmul2(
	double C[2][2], // Output
	double A[2][2], // Input
	double B[2][2]) // Input
{
	double T[2][2];
	int n = 2;
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < n; j++) {
			double sum = 0.0;
			for (int k = 0; k < n; k++)
				sum += A[i][k] * B[k][j];
			T[i][j] = sum;
		}
	}
	for (int i = 0; i < n; i++)
		for (int j = 0; j < n; j++)
			C[i][j] = T[i][j];
}

static void matmul2t(
	double C[2][2], // Output
	double A[2][2], // Input
	double B[2][2]) // Input
{
	double T[2][2];
	int n = 2;
	for (int i = 0; i < n; i++) {
		for (int j = 0; j < n; j++) {
			double sum = 0.0;
			for (int k = 0; k < n; k++)
				sum += A[k][i] * B[k][j];
			T[i][j] = sum;
		}
	}
	for (int i = 0; i < n; i++)
		for (int j = 0; j < n; j++)
			C[i][j] = T[i][j];
}
