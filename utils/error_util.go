package utils

import "fmt"

func NewSliceBoundsOutOfRangeError(i, cap int) error {
	return fmt.Errorf(" runtime error: slice bounds out of range [:%d] with capacity %d", i, cap)
}
