// +build amd64
// errorcheck -0 -d=ssa/prove/debug=1

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "math"

func f0(a []int) int {
	a[0] = 1
	a[0] = 1 // ERROR "Proved IsInBounds$"
	a[6] = 1
	a[6] = 1 // ERROR "Proved IsInBounds$"
	a[5] = 1 // ERROR "Proved IsInBounds$"
	a[5] = 1 // ERROR "Proved IsInBounds$"
	return 13
}

func f1(a []int) int {
	if len(a) <= 5 {
		return 18
	}
	a[0] = 1 // ERROR "Proved IsInBounds$"
	a[0] = 1 // ERROR "Proved IsInBounds$"
	a[6] = 1
	a[6] = 1 // ERROR "Proved IsInBounds$"
	a[5] = 1 // ERROR "Proved IsInBounds$"
	a[5] = 1 // ERROR "Proved IsInBounds$"
	return 26
}

func f1b(a []int, i int, j uint) int {
	if i >= 0 && i < len(a) {
		return a[i] // ERROR "Proved IsInBounds$"
	}
	if i >= 10 && i < len(a) {
		return a[i] // ERROR "Proved IsInBounds$"
	}
	if i >= 10 && i < len(a) {
		return a[i] // ERROR "Proved IsInBounds$"
	}
	if i >= 10 && i < len(a) { // todo: handle this case
		return a[i-10]
	}
	if j < uint(len(a)) {
		return a[j] // ERROR "Proved IsInBounds$"
	}
	return 0
}

func f1c(a []int, i int64) int {
	c := uint64(math.MaxInt64 + 10) // overflows int
	d := int64(c)
	if i >= d && i < int64(len(a)) {
		// d overflows, should not be handled.
		return a[i]
	}
	return 0
}

func f2(a []int) int {
	for i := range a { // ERROR "Induction variable: limits \[0,\?\), increment 1"
		a[i+1] = i
		a[i+1] = i // ERROR "Proved IsInBounds$"
	}
	return 34
}

func f3(a []uint) int {
	for i := uint(0); i < uint(len(a)); i++ {
		a[i] = i // ERROR "Proved IsInBounds$"
	}
	return 41
}

func f4a(a, b, c int) int {
	if a < b {
		if a == b { // ERROR "Disproved Eq64$"
			return 47
		}
		if a > b { // ERROR "Disproved Greater64$"
			return 50
		}
		if a < b { // ERROR "Proved Less64$"
			return 53
		}
		// We can't get to this point and prove knows that, so
		// there's no message for the next (obvious) branch.
		if a != a {
			return 56
		}
		return 61
	}
	return 63
}

func f4b(a, b, c int) int {
	if a <= b {
		if a >= b {
			if a == b { // ERROR "Proved Eq64$"
				return 70
			}
			return 75
		}
		return 77
	}
	return 79
}

func f4c(a, b, c int) int {
	if a <= b {
		if a >= b {
			if a != b { // ERROR "Disproved Neq64$"
				return 73
			}
			return 75
		}
		return 77
	}
	return 79
}

func f4d(a, b, c int) int {
	if a < b {
		if a < c {
			if a < b { // ERROR "Proved Less64$"
				if a < c { // ERROR "Proved Less64$"
					return 87
				}
				return 89
			}
			return 91
		}
		return 93
	}
	return 95
}

func f4e(a, b, c int) int {
	if a < b {
		if b > a { // ERROR "Proved Greater64$"
			return 101
		}
		return 103
	}
	return 105
}

func f4f(a, b, c int) int {
	if a <= b {
		if b > a {
			if b == a { // ERROR "Disproved Eq64$"
				return 112
			}
			return 114
		}
		if b >= a { // ERROR "Proved Geq64$"
			if b == a { // ERROR "Proved Eq64$"
				return 118
			}
			return 120
		}
		return 122
	}
	return 124
}

func f5(a, b uint) int {
	if a == b {
		if a <= b { // ERROR "Proved Leq64U$"
			return 130
		}
		return 132
	}
	return 134
}

// These comparisons are compile time constants.
func f6a(a uint8) int {
	if a < a { // ERROR "Disproved Less8U$"
		return 140
	}
	return 151
}

func f6b(a uint8) int {
	if a < a { // ERROR "Disproved Less8U$"
		return 140
	}
	return 151
}

func f6x(a uint8) int {
	if a > a { // ERROR "Disproved Greater8U$"
		return 143
	}
	return 151
}

func f6d(a uint8) int {
	if a <= a { // ERROR "Proved Leq8U$"
		return 146
	}
	return 151
}

func f6e(a uint8) int {
	if a >= a { // ERROR "Proved Geq8U$"
		return 149
	}
	return 151
}

func f7(a []int, b int) int {
	if b < len(a) {
		a[b] = 3
		if b < len(a) { // ERROR "Proved Less64$"
			a[b] = 5 // ERROR "Proved IsInBounds$"
		}
	}
	return 161
}

func f8(a, b uint) int {
	if a == b {
		return 166
	}
	if a > b {
		return 169
	}
	if a < b { // ERROR "Proved Less64U$"
		return 172
	}
	return 174
}

func f9(a, b bool) int {
	if a {
		return 1
	}
	if a || b { // ERROR "Disproved Arg$"
		return 2
	}
	return 3
}

func f10(a string) int {
	n := len(a)
	// We optimize comparisons with small constant strings (see cmd/compile/internal/gc/walk.go),
	// so this string literal must be long.
	if a[:n>>1] == "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" {
		return 0
	}
	return 1
}

func f11a(a []int, i int) {
	useInt(a[i])
	useInt(a[i]) // ERROR "Proved IsInBounds$"
}

func f11b(a []int, i int) {
	useSlice(a[i:])
	useSlice(a[i:]) // ERROR "Proved IsSliceInBounds$"
}

func f11c(a []int, i int) {
	useSlice(a[:i])
	useSlice(a[:i]) // ERROR "Proved IsSliceInBounds$"
}

func f11d(a []int, i int) {
	useInt(a[2*i+7])
	useInt(a[2*i+7]) // ERROR "Proved IsInBounds$"
}

func f12(a []int, b int) {
	useSlice(a[:b])
}

func f13a(a, b, c int, x bool) int {
	if a > 12 {
		if x {
			if a < 12 { // ERROR "Disproved Less64$"
				return 1
			}
		}
		if x {
			if a <= 12 { // ERROR "Disproved Leq64$"
				return 2
			}
		}
		if x {
			if a == 12 { // ERROR "Disproved Eq64$"
				return 3
			}
		}
		if x {
			if a >= 12 { // ERROR "Proved Geq64$"
				return 4
			}
		}
		if x {
			if a > 12 { // ERROR "Proved Greater64$"
				return 5
			}
		}
		return 6
	}
	return 0
}

func f13b(a int, x bool) int {
	if a == -9 {
		if x {
			if a < -9 { // ERROR "Disproved Less64$"
				return 7
			}
		}
		if x {
			if a <= -9 { // ERROR "Proved Leq64$"
				return 8
			}
		}
		if x {
			if a == -9 { // ERROR "Proved Eq64$"
				return 9
			}
		}
		if x {
			if a >= -9 { // ERROR "Proved Geq64$"
				return 10
			}
		}
		if x {
			if a > -9 { // ERROR "Disproved Greater64$"
				return 11
			}
		}
		return 12
	}
	return 0
}

func f13c(a int, x bool) int {
	if a < 90 {
		if x {
			if a < 90 { // ERROR "Proved Less64$"
				return 13
			}
		}
		if x {
			if a <= 90 { // ERROR "Proved Leq64$"
				return 14
			}
		}
		if x {
			if a == 90 { // ERROR "Disproved Eq64$"
				return 15
			}
		}
		if x {
			if a >= 90 { // ERROR "Disproved Geq64$"
				return 16
			}
		}
		if x {
			if a > 90 { // ERROR "Disproved Greater64$"
				return 17
			}
		}
		return 18
	}
	return 0
}

func f13d(a int) int {
	if a < 5 {
		if a < 9 { // ERROR "Proved Less64$"
			return 1
		}
	}
	return 0
}

func f13e(a int) int {
	if a > 9 {
		if a > 5 { // ERROR "Proved Greater64$"
			return 1
		}
	}
	return 0
}

func f13f(a int64) int64 {
	if a > math.MaxInt64 {
		if a == 0 { // ERROR "Disproved Eq64$"
			return 1
		}
	}
	return 0
}

func f13g(a int) int {
	if a < 3 {
		return 5
	}
	if a > 3 {
		return 6
	}
	if a == 3 { // ERROR "Proved Eq64$"
		return 7
	}
	return 8
}

func f13h(a int) int {
	if a < 3 {
		if a > 1 {
			if a == 2 { // ERROR "Proved Eq64$"
				return 5
			}
		}
	}
	return 0
}

func f13i(a uint) int {
	if a == 0 {
		return 1
	}
	if a > 0 { // ERROR "Proved Greater64U$"
		return 2
	}
	return 3
}

func f14(p, q *int, a []int) {
	// This crazy ordering usually gives i1 the lowest value ID,
	// j the middle value ID, and i2 the highest value ID.
	// That used to confuse CSE because it ordered the args
	// of the two + ops below differently.
	// That in turn foiled bounds check elimination.
	i1 := *p
	j := *q
	i2 := *p
	useInt(a[i1+j])
	useInt(a[i2+j]) // ERROR "Proved IsInBounds$"
}

func f15(s []int, x int) {
	useSlice(s[x:])
	useSlice(s[:x]) // ERROR "Proved IsSliceInBounds$"
}

func f16(s []int) []int {
	if len(s) >= 10 {
		return s[:10] // ERROR "Proved IsSliceInBounds$"
	}
	return nil
}

func f17(b []int) {
	for i := 0; i < len(b); i++ { // ERROR "Induction variable: limits \[0,\?\), increment 1"
		// This tests for i <= cap, which we can only prove
		// using the derived relation between len and cap.
		// This depends on finding the contradiction, since we
		// don't query this condition directly.
		useSlice(b[:i]) // ERROR "Proved IsSliceInBounds$"
	}
}

func f18(b []int, x int, y uint) {
	_ = b[x]
	_ = b[y]

	if x > len(b) { // ERROR "Disproved Greater64$"
		return
	}
	if y > uint(len(b)) { // ERROR "Disproved Greater64U$"
		return
	}
	if int(y) > len(b) { // ERROR "Disproved Greater64$"
		return
	}
}

func sm1(b []int, x int) {
	// Test constant argument to slicemask.
	useSlice(b[2:8]) // ERROR "Proved slicemask not needed$"
	// Test non-constant argument with known limits.
	if cap(b) > 10 {
		useSlice(b[2:]) // ERROR "Proved slicemask not needed$"
	}
}

func lim1(x, y, z int) {
	// Test relations between signed and unsigned limits.
	if x > 5 {
		if uint(x) > 5 { // ERROR "Proved Greater64U$"
			return
		}
	}
	if y >= 0 && y < 4 {
		if uint(y) > 4 { // ERROR "Disproved Greater64U$"
			return
		}
		if uint(y) < 5 { // ERROR "Proved Less64U$"
			return
		}
	}
	if z < 4 {
		if uint(z) > 4 { // Not provable without disjunctions.
			return
		}
	}
}

// fence1–4 correspond to the four fence-post implications.

func fence1(b []int, x, y int) {
	// Test proofs that rely on fence-post implications.
	if x+1 > y {
		if x < y { // ERROR "Disproved Less64$"
			return
		}
	}
	if len(b) < cap(b) {
		// This eliminates the growslice path.
		b = append(b, 1) // ERROR "Disproved Greater64$"
	}
}

func fence2(x, y int) {
	if x-1 < y {
		if x > y { // ERROR "Disproved Greater64$"
			return
		}
	}
}

func fence3(b []int, x, y int64) {
	if x-1 >= y {
		if x <= y { // Can't prove because x may have wrapped.
			return
		}
	}

	if x != math.MinInt64 && x-1 >= y {
		if x <= y { // ERROR "Disproved Leq64$"
			return
		}
	}

	if n := len(b); n > 0 {
		b[n-1] = 0 // ERROR "Proved IsInBounds$"
	}
}

func fence4(x, y int64) {
	if x >= y+1 {
		if x <= y {
			return
		}
	}
	if y != math.MaxInt64 && x >= y+1 {
		if x <= y { // ERROR "Disproved Leq64$"
			return
		}
	}
}

// Check transitive relations
func trans1(x, y int64) {
	if x > 5 {
		if y > x {
			if y > 2 { // ERROR "Proved Greater64"
				return
			}
		} else if y == x {
			if y > 5 { // ERROR "Proved Greater64"
				return
			}
		}
	}
	if x >= 10 {
		if y > x {
			if y > 10 { // ERROR "Proved Greater64"
				return
			}
		}
	}
}

func trans2(a, b []int, i int) {
	if len(a) != len(b) {
		return
	}

	_ = a[i]
	_ = b[i] // ERROR "Proved IsInBounds$"
}

func trans3(a, b []int, i int) {
	if len(a) > len(b) {
		return
	}

	_ = a[i]
	_ = b[i] // ERROR "Proved IsInBounds$"
}

//go:noinline
func useInt(a int) {
}

//go:noinline
func useSlice(a []int) {
}

func main() {
}
