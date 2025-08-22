package fn

import (
	"cmp"
	"math/rand"
	"slices"
)

// Clamp constrains a value to be within the specified minimum and maximum bounds.
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Limit returns a new slice containing at most n elements from the input slice.
// If n is greater than the length of the slice, returns the entire slice.
func Limit[T any](slice []T, n int) []T {
	return slice[:min(len(slice), n)]
}

// Map performs an in-place transformation of the slice if the input and output types are the same
func Map[T, R any](slice []T, f func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

// MapIndexed performs an in-place transformation if possible
func MapIndexed[T, R any](slice []T, f func(int, T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = f(i, v)
	}
	return result
}

// Filter returns a new slice containing only the elements that satisfy the predicate
func Filter[T any](slice []T, pred func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, x := range slice {
		if pred(x) {
			result = append(result, x)
		}
	}
	return result
}

// FilterInPlace filters a slice in-place without allocation, modifying and returning the original slice
func FilterInPlace[T any](a []T, f func(T) bool) []T {
	b := a[:0]
	for _, x := range a {
		if f(x) {
			b = append(b, x)
		}
	}
	// Clear remaining elements for garbage collection
	var zero T
	for i := len(b); i < len(a); i++ {
		a[i] = zero
	}
	return b
}

// Reduce uses a more efficient single-pass reduction
func Reduce[T, R any](slice []T, initial R, f func(R, T) R) R {
	result := initial
	for i := range slice {
		result = f(result, slice[i])
	}
	return result
}

// Any returns true if any element satisfies the predicate
func Any[T any](slice []T, pred func(T) bool) bool {
	return slices.ContainsFunc(slice, pred)
}

// All returns true if all elements satisfy the predicate
func All[T any](slice []T, pred func(T) bool) bool {
	for _, v := range slice {
		if !pred(v) {
			return false
		}
	}
	return true
}

// Unique modifies the slice in-place to remove duplicates while preserving order.
func Unique[T comparable](slice []T) []T {
	if len(slice) <= 1 {
		return slice
	}
	seen := make(map[T]struct{}, len(slice))
	n := 0
	for _, v := range slice {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		slice[n] = v
		n++
	}
	// Clear remaining elements to help GC of references
	var zero T
	for i := n; i < len(slice); i++ {
		slice[i] = zero
	}
	return slice[:n]
}

// IfElse provides a conditional operator that returns trueVal if condition is true, falseVal otherwise
func IfElse[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

// Reverse reverses the elements of a slice in-place
func Reverse[T any](a []T) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

// Shuffle randomly reorders the elements in a slice using Fisher-Yates algorithm
func Shuffle[T any](a []T) {
	for i := len(a) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// Batch splits a slice into batches of specified size with minimal allocation
func Batch[T any](slice []T, batchSize int) [][]T {
	if batchSize <= 0 {
		return nil
	}

	var batches [][]T
	for batchSize < len(slice) {
		slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
	}
	if len(slice) > 0 {
		batches = append(batches, slice)
	}
	return batches
}

// First returns the first element that satisfies the predicate
func First[T any](slice []T, pred func(T) bool) (T, bool) {
	for _, x := range slice {
		if pred(x) {
			return x, true
		}
	}
	var zero T
	return zero, false
}

// Delete removes all occurrences of an element from a slice
// Warning! You must reassign the slice to the result of this function:
// slice = fn.Delete(slice, value)
func Delete[T comparable](slice []T, value T) []T {
	return slices.DeleteFunc(slice, func(el T) bool { return el == value })
}

// ToIfaceSlice converts a slice of an arbitrary type to a []any
func ToIfaceSlice(arr ...any) []any {
	return arr
}
