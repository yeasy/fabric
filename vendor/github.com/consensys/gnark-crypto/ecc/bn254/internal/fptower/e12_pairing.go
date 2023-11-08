package fptower

func (z *E12) nSquare(n int) {
	for i := 0; i < n; i++ {
		z.CyclotomicSquare(z)
	}
}

func (z *E12) nSquareCompressed(n int) {
	for i := 0; i < n; i++ {
		z.CyclotomicSquareCompressed(z)
	}
}

// Expt set z to xᵗ (mod q¹²) and return z (t is the generator of the curve)
func (z *E12) Expt(x *E12) *E12 {
	// Expt computation is derived from the addition chain:
	//
	//	_10     = 2*1
	//	_100    = 2*_10
	//	_1000   = 2*_100
	//	_10000  = 2*_1000
	//	_10001  = 1 + _10000
	//	_10011  = _10 + _10001
	//	_10100  = 1 + _10011
	//	_11001  = _1000 + _10001
	//	_100010 = 2*_10001
	//	_100111 = _10011 + _10100
	//	_101001 = _10 + _100111
	//	i27     = (_100010 << 6 + _100 + _11001) << 7 + _11001
	//	i44     = (i27 << 8 + _101001 + _10) << 6 + _10001
	//	i70     = ((i44 << 8 + _101001) << 6 + _101001) << 10
	//	return    (_100111 + i70) << 6 + _101001 + _1000
	//
	// Operations: 62 squares 17 multiplies
	//
	// Generated by github.com/mmcloughlin/addchain v0.4.0.

	// Allocate Temporaries.
	var result, t0, t1, t2, t3, t4, t5, t6 E12

	// Step 1: t3 = x^0x2
	t3.CyclotomicSquare(x)

	// Step 2: t5 = x^0x4
	t5.CyclotomicSquare(&t3)

	// Step 3: result = x^0x8
	result.CyclotomicSquare(&t5)

	// Step 4: t0 = x^0x10
	t0.CyclotomicSquare(&result)

	// Step 5: t2 = x^0x11
	t2.Mul(x, &t0)

	// Step 6: t0 = x^0x13
	t0.Mul(&t3, &t2)

	// Step 7: t1 = x^0x14
	t1.Mul(x, &t0)

	// Step 8: t4 = x^0x19
	t4.Mul(&result, &t2)

	// Step 9: t6 = x^0x22
	t6.CyclotomicSquare(&t2)

	// Step 10: t1 = x^0x27
	t1.Mul(&t0, &t1)

	// Step 11: t0 = x^0x29
	t0.Mul(&t3, &t1)

	// Step 17: t6 = x^0x880
	t6.nSquare(6)

	// Step 18: t5 = x^0x884
	t5.Mul(&t5, &t6)

	// Step 19: t5 = x^0x89d
	t5.Mul(&t4, &t5)

	// Step 26: t5 = x^0x44e80
	t5.nSquare(7)

	// Step 27: t4 = x^0x44e99
	t4.Mul(&t4, &t5)

	// Step 35: t4 = x^0x44e9900
	t4.nSquare(8)

	// Step 36: t4 = x^0x44e9929
	t4.Mul(&t0, &t4)

	// Step 37: t3 = x^0x44e992b
	t3.Mul(&t3, &t4)

	// Step 43: t3 = x^0x113a64ac0
	t3.nSquare(6)

	// Step 44: t2 = x^0x113a64ad1
	t2.Mul(&t2, &t3)

	// Step 52: t2 = x^0x113a64ad100
	t2.nSquare(8)

	// Step 53: t2 = x^0x113a64ad129
	t2.Mul(&t0, &t2)

	// Step 59: t2 = x^0x44e992b44a40
	t2.nSquare(6)

	// Step 60: t2 = x^0x44e992b44a69
	t2.Mul(&t0, &t2)

	// Step 70: t2 = x^0x113a64ad129a400
	t2.nSquare(10)

	// Step 71: t1 = x^0x113a64ad129a427
	t1.Mul(&t1, &t2)

	// Step 77: t1 = x^0x44e992b44a6909c0
	t1.nSquare(6)

	// Step 78: t0 = x^0x44e992b44a6909e9
	t0.Mul(&t0, &t1)

	// Step 79: result = x^0x44e992b44a6909f1
	z.Mul(&result, &t0)

	return z
}

// MulBy034 multiplication by sparse element (c0,0,0,c3,c4,0)
func (z *E12) MulBy034(c0, c3, c4 *E2) *E12 {

	var a, b, d E6

	a.MulByE2(&z.C0, c0)

	b.Set(&z.C1)
	b.MulBy01(c3, c4)

	c0.Add(c0, c3)
	d.Add(&z.C0, &z.C1)
	d.MulBy01(c0, c4)

	z.C1.Add(&a, &b).Neg(&z.C1).Add(&z.C1, &d)
	z.C0.MulByNonResidue(&b).Add(&z.C0, &a)

	return z
}

// Mul034By034 multiplication of sparse element (c0,0,0,c3,c4,0) by sparse element (d0,0,0,d3,d4,0)
func Mul034By034(d0, d3, d4, c0, c3, c4 *E2) [5]E2 {
	var z00, tmp, x0, x3, x4, x04, x03, x34 E2
	x0.Mul(c0, d0)
	x3.Mul(c3, d3)
	x4.Mul(c4, d4)
	tmp.Add(c0, c4)
	x04.Add(d0, d4).
		Mul(&x04, &tmp).
		Sub(&x04, &x0).
		Sub(&x04, &x4)
	tmp.Add(c0, c3)
	x03.Add(d0, d3).
		Mul(&x03, &tmp).
		Sub(&x03, &x0).
		Sub(&x03, &x3)
	tmp.Add(c3, c4)
	x34.Add(d3, d4).
		Mul(&x34, &tmp).
		Sub(&x34, &x3).
		Sub(&x34, &x4)

	z00.MulByNonResidue(&x4).
		Add(&z00, &x0)

	return [5]E2{z00, x3, x34, x03, x04}
}

// MulBy01234 multiplies z by an E12 sparse element of the form (x0, x1, x2, x3, x4, 0)
func (z *E12) MulBy01234(x *[5]E2) *E12 {
	var c1, a, b, c, z0, z1 E6
	c0 := &E6{B0: x[0], B1: x[1], B2: x[2]}
	c1.B0 = x[3]
	c1.B1 = x[4]
	a.Add(&z.C0, &z.C1)
	b.Add(c0, &c1)
	a.Mul(&a, &b)
	b.Mul(&z.C0, c0)
	c.Set(&z.C1).MulBy01(&x[3], &x[4])
	z1.Sub(&a, &b)
	z1.Sub(&z1, &c)
	z0.MulByNonResidue(&c)
	z0.Add(&z0, &b)

	z.C0 = z0
	z.C1 = z1

	return z
}
