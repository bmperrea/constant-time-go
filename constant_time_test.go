package constant_time_test

import (
	"math/rand"
	"testing"

	"github.com/bmperrea/constant-time-go"
)

func TestLessThanUint32(t *testing.T) {
	tests := []struct {
		x uint32
		y uint32
		a uint32
	}{
		{0,1, 1},
		{2, 2, 0},
		{1 << 31, 1 << 31, 0},
		{17, 1 << 31, 1},
		{2^32, 0, 0},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		answer := constant_time.LessThanUint32(test.x, test.y)
		if test.a != answer {
			t.Errorf("LessThanUint32 #%d wrong result\ngot: %v\n"+
				"want: %v", i, answer, test.a)
			continue
		}
	}
}

func TestLessOrEqUint32(t *testing.T) {
	tests := []struct {
		x uint32
		y uint32
		a uint32
	}{
		{0, 1, 1},
		{2, 2, 1},
		{1 << 31, 1 << 31, 1},
		{17, 1 << 31, 1},
		{2 ^ 32, 0, 0},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		answer := constant_time.LessOrEqUint32(test.x, test.y)
		if test.a != answer {
			t.Errorf("LessOrEqUint32 #%d wrong result\ngot: %v\n"+
				"want: %v", i, answer, test.a)
			continue
		}
	}
}

func TestEqUint32(t *testing.T) {
	tests := []struct {
		x uint32
		y uint32
		a uint32
	}{
		{0, 1, 0},
		{2, 2, 1},
		{1 << 31, 1 << 31, 1},
		{17, 1 << 31, 0},
		{2 ^ 32, 0, 0},
	}

	t.Logf("Running %d tests", len(tests))
	for i, test := range tests {
		answer := constant_time.EqUint32(test.x, test.y)
		if test.a != answer {
			t.Errorf("LessOrEqUint32 #%d wrong result\ngot: %v\n"+
				"want: %v", i, answer, test.a)
			continue
		}
	}
}

// Notes on decent benchamrks

// 1. always record the result of your function to prevent
//    the compiler eliminating the function call. print
//    something related to all results to further avoid compiler optimization

// 2. prefer to use a pre-allocated array of input values, preferrably random,
//    to avoid artificial branch prediction speedups

// 3. stop and start the timer manually to avoid mixing in setup
//    and overhead time

// 4. Compare against a benchmark with your function removed
//    to measure how much time is spent in the remaining overhead+setup

var num_rand = 100

// generates an slice of somewhat random numbers (actually only num_rand to keep it quick)
func randUint32s(n int) []uint32 {
	arr := make([]uint32, n)
	for i := range arr {
		arr[i] = rand.Uint32()
	}
	return arr
}

// The control benchmark
func BenchmarkNothingUint32(b *testing.B) {
	x := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= x[i]
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The implementation based on the one in crypto/subtle using 64bit ops internally
func BenchmarkLessThanUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.LessThanUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// A branching implementation
func BenchmarkBranchingLessThanUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.BranchingLessThanUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The implementation based on the one in crypto/subtle
func BenchmarkLessOrEqUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.LessOrEqUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// A branching implementation
func BenchmarkBranchingLessOrEqUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.BranchingLessOrEqUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The implementation based on the one in crypto/subtle
func BenchmarkEqUint32Alternate(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.EqUint32Alternate(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The implementation based on int64
func BenchmarkEqUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.EqUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// A branching implementation
func BenchmarkBranchingEqUint32(b *testing.B) {
	x := randUint32s(b.N)
	y := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.BranchingEqUint32(x[i], y[i])
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The implementation based on the one in crypto/subtle
func BenchmarkSelectUint32(b *testing.B) {
	x := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.SelectUint32(x[i], 2, 1)
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// A branching implementation
func BenchmarkBranchingSelectUint32(b *testing.B) {
	x := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= constant_time.BranchingSelectUint32(x[i], 2, 1)
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}

// The control benchmark again
func BenchmarkNothingUint32Again(b *testing.B) {
	x := randUint32s(b.N)
	result := uint32(0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result &= x[i]
	}
	b.StopTimer()
	if (result == 73) {
		print("whatever")
	}
}
