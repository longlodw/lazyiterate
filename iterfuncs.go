package lazyiterate

import (
	"fmt"
	"iter"
)

// Any returns true if any element in the sequence satisfies the predicate.
func All[T any](it iter.Seq[T], pred func(T) bool) bool {
	for v := range it {
		if !pred(v) {
			return false
		}
	}
	return true
}

// All2 returns true if all key-value pairs in the sequence satisfy the predicate.
func All2[K any, V any](it iter.Seq2[K, V], pred func(K, V) bool) bool {
	for k, v := range it {
		if !pred(k, v) {
			return false
		}
	}
	return true
}

// Any returns true if any element in the sequence satisfies the predicate.
func Any[T any](it iter.Seq[T], pred func(T) bool) bool {
	for v := range it {
		if pred(v) {
			return true
		}
	}
	return false
}

// Any2 returns true if any key-value pair in the sequence satisfies the predicate.
func Any2[K any, V any](it iter.Seq2[K, V], pred func(K, V) bool) bool {
	for k, v := range it {
		if pred(k, v) {
			return true
		}
	}
	return false
}

// Count returns the number of elements in the sequence.
func Count[T any](it iter.Seq[T]) int {
	count := 0
	for range it {
		count++
	}
	return count
}

// Count2 returns the number of key-value pairs in the sequence.
func Count2[K any, V any](it iter.Seq2[K, V]) int {
	count := 0
	for range it {
		count++
	}
	return count
}

// Find returns the first element in the sequence that satisfies the predicate.
func Find[T any](it iter.Seq[T], pred func(T) bool) (T, error) {
	for v := range it {
		if pred(v) {
			return v, nil
		}
	}
	var zero T
	return zero, fmt.Errorf("no element found")
}

// Find2 returns the first key-value pair in the sequence that satisfies the predicate.
func Find2[K any, V any](it iter.Seq2[K, V], pred func(K, V) bool) (K, V, error) {
	for k, v := range it {
		if pred(k, v) {
			return k, v, nil
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, fmt.Errorf("no element found")
}

// Filter returns a new sequence containing only the elements that satisfy the predicate.
func Filter[T any](it iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		it(func(v T) bool {
			if pred(v) {
				return yield(v)
			}
			return true // continue iteration
		})
	}
}

// Filter2 returns a new sequence containing only the key-value pairs that satisfy the predicate.
func Filter2[K any, V any](it iter.Seq2[K, V], pred func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		it(func(k K, v V) bool {
			if pred(k, v) {
				return yield(k, v)
			}
			return true // continue iteration
		})
	}
}

// Map returns a new sequence containing the results of applying the function to each element.
func Map[T, R any](it iter.Seq[T], fn func(T) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		it(func(v T) bool {
			return yield(fn(v))
		})
	}
}

// Map2 returns a new sequence containing the results of applying the function to each key-value pair.
func Map2[K, V, R any](it iter.Seq2[K, V], fn func(K, V) R) iter.Seq[R] {
	return func(yield func(R) bool) {
		it(func(k K, v V) bool {
			return yield(fn(k, v))
		})
	}
}

// Reduce applies a function cumulatively to the elements of the sequence, reducing it to a single value.
func Reduce[A any, T any](it iter.Seq[T], fn func(A, T) A, init A) A {
	acc := init
	for v := range it {
		acc = fn(acc, v)
	}
	return acc
}

// Reduce2 applies a function cumulatively to the key-value pairs of the sequence, reducing it to a single value.
func Reduce2[A any, K any, V any](it iter.Seq2[K, V], fn func(A, K, V) A, init A) A {
	acc := init
	for k, v := range it {
		acc = fn(acc, k, v)
	}
	return acc
}

// Reverse returns a new sequence with the elements in reverse order.
func Reverse[T any](it iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		var stack []T
		it(func(v T) bool {
			stack = append(stack, v)
			return true // continue iteration
		})
		for i := len(stack) - 1; i >= 0; i-- {
			if !yield(stack[i]) {
				break // stop if yield returns false
			}
		}
	}
}

// Reverse2 returns a new sequence with the key-value pairs in reverse order.
func Reverse2[K any, V any](it iter.Seq2[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var stack []struct {
			Key K
			Val V
		}
		it(func(k K, v V) bool {
			stack = append(stack, struct {
				Key K
				Val V
			}{k, v})
			return true // continue iteration
		})
		for i := len(stack) - 1; i >= 0; i-- {
			if !yield(stack[i].Key, stack[i].Val) {
				break // stop if yield returns false
			}
		}
	}
}

// Skip returns a new sequence that skips the first n elements.
func Skip[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		it(func(v T) bool {
			if count >= n {
				return yield(v)
			}
			count++
			return true // continue iteration
		})
	}
}

// Skip2 returns a new sequence that skips the first n key-value pairs.
func Skip2[K any, V any](it iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		it(func(k K, v V) bool {
			if count >= n {
				return yield(k, v)
			}
			count++
			return true // continue iteration
		})
	}
}

// Take returns a new sequence that takes the first n elements.
func Take[T any](it iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		it(func(v T) bool {
			if count < n {
				count++
				return yield(v)
			}
			return false // stop iteration after n elements
		})
	}
}

// Take2 returns a new sequence that takes the first n key-value pairs.
func Take2[K any, V any](it iter.Seq2[K, V], n int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		count := 0
		it(func(k K, v V) bool {
			if count < n {
				count++
				return yield(k, v)
			}
			return false // stop iteration after n pairs
		})
	}
}

// Zip returns a new sequence that combines elements from two sequences into pairs.
func Zip[T1, T2 any](it1 iter.Seq[T1], it2 iter.Seq[T2]) iter.Seq2[T1, T2] {
	return func(yield func(T1, T2) bool) {
		next1, stop1 := iter.Pull(it1)
		defer stop1()
		next2, stop2 := iter.Pull(it2)
		defer stop2()
		for {
			v1, ok1 := next1()
			v2, ok2 := next2()
			if !ok1 || !ok2 {
				break // stop if either sequence is exhausted
			}
			if !yield(v1, v2) {
				break // stop if yield returns false
			}
		}
	}
}
