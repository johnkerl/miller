#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include "lib/mlrutil.h"
#include "lib/mlrmath.h"
#include "lib/mlr_globals.h"

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

// ----------------------------------------------------------------
// Normal cumulative distribution function, expressed in terms of erfc library
// function (which is awkward, but exists).
double qnorm(double x) {
	return 0.5 * erfc(-x * M_SQRT1_2);
}

// ----------------------------------------------------------------
// This is a tangent-following method not unlike Newton-Raphson:
// * We can compute qnorm(y) = integral from -infinity to y of (1/sqrt(2pi)) exp(-t^2/2) dt.
// * We can compute derivative of qnorm(y) = (1/sqrt(2pi)) exp(-y^2/2).
// * We cannot explicitly compute invqnorm(y).
// * If dx/dy = (1/sqrt(2pi)) exp(-y^2/2) then dy/dx = sqrt(2pi) exp(y^2/2).
//
// This means we *can* compute the derivative of invqnorm even though we
// can't compute the function itself. So the essence of the method is to
// follow the tangent line to form successive approximations: we have known function input x
// and unknown function output y and initial guess y0.  At each step we find the intersection
// of the tangent line at y_n with the vertical line at x, to find y_{n+1}. Specificall:
//
// * Even though we can't compute y = q^-1(x) we can compute x = q(y).
// * Start with initial guess for y (y0 = 0.0 or y0 = x both are OK).
// * Find x = q(y). Since q (and therefore q^-1) are 1-1, we're done if qnorm(invqnorm(x)) is small.
// * Else iterate: using point-slope form, (y_{n+1} - y_n) / (x_{n+1} - x_n) = m = sqrt(2pi) exp(y_n^2/2).
//   Here x_2 = x (the input) and x_1 = q(y_1).
// * Solve for y_{n+1} and repeat.

#define INVQNORM_TOL 1e-9
#define INVQNORM_MAXITER 30

double invqnorm(double x) {
	// Initial approximation is linear. Starting with y0 = 0.0 works just as well.
	double y0 = x - 0.5;
	if (x <= 0.0)
		return 0.0;
	if (x >= 1.0)
		return 0.0;

	double y = y0;
	int niter = 0;
	while (1) {
		double backx = qnorm(y);
		double err = fabs(x - backx);
		if (err < INVQNORM_TOL)
			break;
		if (niter > INVQNORM_MAXITER) {
			fprintf(stderr, "%s: internal coding error: max iterations %d exceeded in invqnorm.\n",
				MLR_GLOBALS.argv0, INVQNORM_MAXITER);
			exit(1);
		}
		double m = sqrt(2*M_PI) * exp(y*y/2.0);
		double delta_y = m * (x - backx);

		y += delta_y;
		niter++;

	}
	return y;
}

// ================================================================
// Logisitic regression
//
// Real-valued x_0 .. x_{N-1}
// 0/1-valued  y_0 .. y_{N-1}
// Model p(x_i == 1)  as
//   p(x, m, b) = 1 / (1 + exp(-m*x-b)
// which is the same as
//   log(p/(1-p)) = m*x + b
// then
//   p(x, m, b) = 1 / (1 + exp(-m*x-b)
//              = exp(m*x+b) / (1 + exp(m*x+b)
// and
//   1-p        = exp(-m*x-b) / (1 + exp(-m*x-b)
//              = 1 / (1 + exp(m*x+b)
// Note for reference just below that
//   dp/dm      = -1 / [1 + exp(-m*x-b)]**2 * (-x) * exp(-m*x-b)
//              = [x exp(-m*x-b)) ] / [1 + exp(-m*x-b)]**2
//              = x * p * (1-p)
// and
//   dp/db      = -1 / [1 + exp(-m*x-b)]**2 * (-1) * exp(-m*x-b)
//              = [exp(-m*x-b)) ] / [1 + exp(-m*x-b)]**2
//              = p * (1-p)
// Write p_i for p(x_i, m, b)
//
// Maximum-likelihood equation:
//   L(m, b)    = prod_{i=0}^{N-1} [ p_i**y_i * (1-p_i)**(1-y_i) ]
//
// Log-likelihood equation:
//   ell(m, b)  = sum{i=0}^{N-1} [ y_i log(p_i) + (1-y_i) log(1-p_i) ]
//              = sum{i=0}^{N-1} [ log(1-p_i) + y_i log(p_i/(1-p_i)) ]
//              = sum{i=0}^{N-1} [ log(1-p_i) + y_i*(m*x_i+b) ]
// Differentiate with respect to parameters:
//
//   d ell/dm   = sum{i=0}^{N-1} [ -1/(1-p_i) dp_i/dm + x_i*y_i ]
//              = sum{i=0}^{N-1} [ -1/(1-p_i) x_i*p_i*(1-p_i) + x_i*y_i ]
//              = sum{i=0}^{N-1} [ x_i(y_i-p_i) ]
//
//   d ell/db   = sum{i=0}^{N-1} [ -1/(1-p_i) dp_i/db + y_i ]
//              = sum{i=0}^{N-1} [ -1/(1-p_i) p_i*(1-p_i) + y_i ]
//              = sum{i=0}^{N-1} [ y_i - p_i ]
//
//
//   d2ell/dm2  = sum{i=0}^{N-1} [ -x_i dp_i/dm ]
//              = sum{i=0}^{N-1} [ -x_i**2 * p_i * (1-p_i) ]
//
//   d2ell/dmdb = sum{i=0}^{N-1} [ -x_i dp_i/db ]
//              = sum{i=0}^{N-1} [ -x_i * p_i * (1-p_i) ]
//
//   d2ell/dbdm = sum{i=0}^{N-1} [ -dp_i/dm ]
//              = sum{i=0}^{N-1} [ -x_i * p_i * (1-p_i) ]
//
//   d2ell/db2  = sum{i=0}^{N-1} [ -dp_i/db ]
//              = sum{i=0}^{N-1} [ -p_i * (1-p_i) ]
//
// Newton-Raphson to minimize ell(m, b):
// * Pick m0, b0
// * [m_{j+1], b_{j+1}] = H^{-1} grad ell(m_j, b_j)
// * grad ell =
//   [ d ell/dm ]
//   [ d ell/db ]
// * H = Hessian of ell = Jacobian of grad ell =
//   [ d2ell/dm2  d2ell/dmdb ]
//   [ d2ell/dmdb d2ell/db2  ]

// p(x,m,b) for logistic regression:
static double lrp(double x, double m, double b) {
	return 1.0 / (1.0 + exp(-m*x-b));
}

// 1 - p(x,m,b) for logistic regression:
static double lrq(double x, double m, double b) {
	return 1.0 / (1.0 + exp(m*x+b));
}

// Supporting routine for mlr_logistic_regression():
static void mlr_logistic_regression_aux(double* xs, double* ys, int n, double* pm, double* pb,
	double m0, double b0, double tol, int maxits)
{
	int its = 0;
	int done = FALSE;
	double m = m0;
	double b = b0;

	while (!done) {
		// Compute derivatives
		double dldm    = 0.0;
		double dldb    = 0.0;
		double d2ldm2  = 0.0;
		double d2ldmdb = 0.0;
		double d2ldb2  = 0.0;
		double ell0    = 0.0;

		for (int i = 0; i < n; i++) {
			double xi = xs[i];
			double yi = ys[i];
			double pi = lrp(xi, m0, b0);
			double qi = lrq(xi, m0, b0);
			dldm += xi*(yi - pi);
			dldb += yi - pi;
			double piqi = pi * qi;
			double xipiqi = xi*piqi;
			double xi2piqi = xi*xipiqi;
			d2ldm2  -= xi2piqi;
			d2ldmdb -= xipiqi;
			d2ldb2  -= piqi;
			ell0 += log(qi) + yi * (m0 * xi + b0);
		}


		// Form the Hessian
		double ha = d2ldm2;
		double hb = d2ldmdb;
		double hc = d2ldmdb;
		double hd = d2ldb2;

		// Invert the Hessian
		double D = ha*hd - hb*hc;
		double Hinva =  hd/D;
		double Hinvb = -hb/D;
		double Hinvc = -hc/D;
		double Hinvd =  ha/D;

		// Compute H^-1 times grad ell
		double Hinvgradm = Hinva*dldm + Hinvb*dldb;
		double Hinvgradb = Hinvc*dldm + Hinvd*dldb;

		// Update [m,b]
		m = m0 - Hinvgradm;
		b = b0 - Hinvgradb;

		double ell = 0.0;
		for (int i = 0; i < n; i++) {
			double xi = xs[i];
			double yi = ys[i];
			double qi = lrq(xi, m, b);
			ell += log(qi) + yi * (m0 * xi + b0);
		}

		// Check for convergence
		double err = fabs(ell - ell0);

#if 0
		printf("its=%d,m=%e,b=%e,dm=%e,db=%e,ell=%e\n", its, m0, b0, -Hinvgradm, -Hinvgradb, ell);
#endif

		if (err < tol)
			done = TRUE;
		if (++its > maxits) {
			fprintf(stderr,
				"mlr_logistic_regression: Newton-Raphson convergence failed after %d iterations. m=%e, b=%e.\n",
					its, m, b);
			exit(1);
		}

		m0 = m;
		b0 = b;
	}

	*pm = m;
	*pb = b;
}

void mlr_logistic_regression(double* xs, double* ys, int n, double* pm, double* pb) {
	double m0     = -0.001;
	double b0     =  0.002;
	double tol    = 1e-9;
	int    maxits = 100;
	mlr_logistic_regression_aux(xs, ys, n, pm, pb, m0, b0, tol, maxits);
}
