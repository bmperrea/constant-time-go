# constant-time-go
A repository for testing constant time approaches to programming in go

This was made for analyzing constant time implementations of comparisons in go
particularly for use in [github.com/btcsuite/btcd/btcec](https://github.com/btcsuite/btcd/tree/master/btcec)

## Intro

Constant time implementations are useful in cryptography where attacks based on runtime have been shown to
completely break the security of otherwise secure algorithms. To achieve this the most accepted approach
is to write logic directly in terms of bitwise operations, avoiding different code branches that take
significantly longer time than others. To aid in this endeavor, golang provides the library [crypto/subtle](https://golang.org/src/crypto/subtle/constant_time.go),
which gives reference implementations of pure bitwise comparison functions and other basic tools.

for example
```
func LessThanUint32(x, y uint32) uint32 {
	xs := int64(x)
	ys := int64(y)
	// the msb keeps the sign
	return uint32(((xs - ys) >> 63) & 1)
}
```

An even more subtle timing issue that is solved by crypto/subtle is that of branch prediction. By
returning integers instead of boolean data, cyrpto/subtle sets one up for further bitwise operations
avoiding using branching statements at all like if, else.

Nonetheless, libraries like btcd often use another apporach to constant time by using branch statments
that take the same amount of compute time. This can be much faster than computing every branch every
time, but has also been shown to be susceptible to timing attacks due to the existence of branch prediction
at the processor level, which can make unusual branches take more compute time.

for example
```
func BranchingLessThanUint32(x, y uint32) uint32 {
   	result := uint32(1)
   	if x < y {
   		result &= 1
   	} else {
   		result &= 0
   	}
   	return result
   }
}
```

Using such a same-time branch approach is often preferrable because it can make for much faster compute,
without sacrificing too much in terms of non-constant run time. This repo analyzes these different
approaches to pin down how much compute time is saved by using same-time branches in place of the
constant time comparators in crypto/subtle.

Note - these problems could be solved in a simple way if there was a native golang way to cast a bool
to an integer so that one could use internal comparators (<, ==) with only one additional operation
to continue with bitwise logic. Unfortunately, in the effort to keep go to a minimal feature set, the
creators of go have agreed not to support this. It would be interesting to try to acheive a more
reliably constant time bool->int function using golang assembly

## Results

```
$ go test --bench mark --benchtime=80ms
goos: darwin
goarch: amd64
pkg: github.com/bmperrea/constant-time-go
BenchmarkNothingUint32-8                   	200000000	         0.63 ns/op
BenchmarkConstantTimeLessThanUint32-8      	100000000	         1.00 ns/op
BenchmarkBranchingLessThanUint32-8         	100000000	         0.98 ns/op
BenchmarkConstantTimeLessOrEqUint32-8      	100000000	         1.06 ns/op
BenchmarkBranchingLessOrEqUint32-8         	100000000	         1.07 ns/op
BenchmarkConstantTimeEqUint32Alternate-8   	50000000	         2.14 ns/op
BenchmarkConstantTimeEqUint32-8            	100000000	         1.09 ns/op
BenchmarkBranchingEqUint32-8               	100000000	         0.98 ns/op
BenchmarkConstantTimeSelectUint32-8        	100000000	         1.00 ns/op
BenchmarkBranchingSelectUint32-8           	100000000	         0.92 ns/op
BenchmarkNothingUint32Again-8                   	200000000	         0.63 ns/op
PASS
ok  	github.com/bmperrea/constant-time-go	53.353s
```

So - the overhead of the loop is pretty high, but even when we take that away, the branching case gives us for 
the less than and lessOrEqual functions, as sell as selection, but results in a sizable speedup for equality checks. 
In particular, after subtracting the control benchmark we get a slow down factor of (2.15 - .63) / (.98 - .63) ~ 4.
However, an implementation based on int64 leads to a significant speedup so that the slow down is only
(1.09 - .63) / (.98 - .63) ~ 1.3 which is significant only if the operation represents a large portion of a calculation.


Conclusion: I recommend using the `crypto/subtle` functions instead of using branching for most situations
    since the additional computation cost is most often immeasurable, and one avoids the possibility of
    timing attacks based on branch prediction. The only possible exception is with equality checks - one might want to 
    use a branching statement to convert (a==b) to an integer before doing more bitwise operations on the result. However,
    a faster implementation based on int64 here improves the situation making these branch-free
    comparisons quite practical.
