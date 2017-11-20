package constant_time_go

var (
	notMSB = ^(uint32(1) << 31)
)

// constant_time.go provides constant time implementations of useful mathematical operations. In addition, these
// functions return integers, using 0 or 1 to represent false or true respectively,
// which is useful for writing logic in terms of bitwise operators

// References
// These functions are based on the sample implementation in golang.org/src/crypto/subtle/constant_time.go
// Here we have converted these functions, and written new ones, for direct use in uint32 arithmetic

// ConstantTimeSelect returns x if v is 1 and y if v is 0.
// Its behavior is undefined if v takes any other value.
func ConstantTimeSelectUint32(v, x, y uint32) uint32 { return ^(v-1)&x | (v-1)&y }

func ConstantTimeLessThanUint32(x, y uint32) uint32 {
	diff := int64(x) - int64(y)
	return uint32((diff >> 63) & 1)
}

func ConstantTimeLessOrEqUint32(x, y uint32) uint32 {
	diff := int64(x) - int64(y)
	return uint32(((diff - 1) >> 63) & 1)
}

// ConstantTimeEq returns 1 if x == y and 0 otherwise.
func ConstantTimeEqUint32Alternate(x, y uint32) uint32 {
	z := ^(x ^ y)
	z &= z >> 16
	z &= z >> 8
	z &= z >> 4
	z &= z >> 2
	z &= z >> 1

	return z & 1
}

// ConstantTimeEq utilizing the same strategy as ConstantTimeLessOrEqUint32 - the sign bit in int64
func ConstantTimeEqUint32(x, y uint32) uint32 {
	diff := int64(x) - int64(y)
	return uint32((((diff - 1) ^ diff) >> 63) & 1)
}

func BranchingSelectUint32(v, x, y uint32) uint32 {
	result := uint32(0)
	if v == 1 {
		result |= x
	} else {
		result |= y
	}
	return result
}

func BranchingLessThanUint32(x, y uint32) uint32 {
	result := uint32(3)
	if x < y {
		result &= 1
	} else {
		result &= 4
	}
	return result
}

func BranchingLessOrEqUint32(x, y uint32) uint32 {
	result := uint32(3)
	if x <= y {
		result &= 1
	} else {
		result &= 4
	}
	return result
}

func BranchingEqUint32(x, y uint32) uint32 {
	result := uint32(3)
	if x == y {
		result &= 1
	} else {
		result &= 4
	}
	return result
}
