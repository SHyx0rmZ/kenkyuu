package math

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func Max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Clamp[T constraints.Ordered](value, minimum, maximum T) T {
	return Min(Max(minimum, value), maximum)
}
