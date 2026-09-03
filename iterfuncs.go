// Package lazyiterate provides lazy, chainable operations for iterators.
package lazyiterate

import (
	"errors"
	"iter"
)

// ErrNotFound is returned by Find when no element satisfies its predicate.
var ErrNotFound = errors.New("lazyiterate: no matching element found")

// LazyIter is a chainable wrapper around an iter.Seq.
//
// Use From to wrap an existing sequence. Operations that produce a sequence
// return LazyIter, so they can be chained; terminal operations return a value.
type LazyIter[T any] iter.Seq[T]

// LazyIter2 is a chainable wrapper around an iter.Seq2.
type LazyIter2[K, V any] iter.Seq2[K, V]

// From wraps seq in a LazyIter.
func From[T any](seq iter.Seq[T]) LazyIter[T] { return LazyIter[T](seq) }

// From2 wraps seq in a LazyIter2.
func From2[K, V any](seq iter.Seq2[K, V]) LazyIter2[K, V] { return LazyIter2[K, V](seq) }

// Seq returns the underlying iter.Seq.
func (it LazyIter[T]) Seq() iter.Seq[T] { return iter.Seq[T](it) }

// Seq returns the underlying iter.Seq2.
func (it LazyIter2[K, V]) Seq() iter.Seq2[K, V] { return iter.Seq2[K, V](it) }

// All reports whether every element satisfies pred.
func (it LazyIter[T]) All(pred func(T) bool) bool {
	for v := range it {
		if !pred(v) {
			return false
		}
	}
	return true
}

// All reports whether every key-value pair satisfies pred.
func (it LazyIter2[K, V]) All(pred func(K, V) bool) bool {
	for k, v := range it {
		if !pred(k, v) {
			return false
		}
	}
	return true
}

// Any reports whether any element satisfies pred.
func (it LazyIter[T]) Any(pred func(T) bool) bool {
	for v := range it {
		if pred(v) {
			return true
		}
	}
	return false
}

// Any reports whether any key-value pair satisfies pred.
func (it LazyIter2[K, V]) Any(pred func(K, V) bool) bool {
	for k, v := range it {
		if pred(k, v) {
			return true
		}
	}
	return false
}

// Count returns the number of elements.
func (it LazyIter[T]) Count() int {
	count := 0
	for range it {
		count++
	}
	return count
}

// Count returns the number of key-value pairs.
func (it LazyIter2[K, V]) Count() int {
	count := 0
	for range it {
		count++
	}
	return count
}

// Find returns the first element that satisfies pred, or ErrNotFound.
func (it LazyIter[T]) Find(pred func(T) bool) (T, error) {
	for v := range it {
		if pred(v) {
			return v, nil
		}
	}
	var zero T
	return zero, ErrNotFound
}

// Find returns the first key-value pair that satisfies pred, or ErrNotFound.
func (it LazyIter2[K, V]) Find(pred func(K, V) bool) (K, V, error) {
	for k, v := range it {
		if pred(k, v) {
			return k, v, nil
		}
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, ErrNotFound
}

// Filter returns the elements that satisfy pred.
func (it LazyIter[T]) Filter(pred func(T) bool) LazyIter[T] {
	return func(yield func(T) bool) {
		it(func(v T) bool { return !pred(v) || yield(v) })
	}
}

// Filter returns the key-value pairs that satisfy pred.
func (it LazyIter2[K, V]) Filter(pred func(K, V) bool) LazyIter2[K, V] {
	return func(yield func(K, V) bool) {
		it(func(k K, v V) bool { return !pred(k, v) || yield(k, v) })
	}
}

// Map applies fn to each element.
func (it LazyIter[T]) Map[R any](fn func(T) R) LazyIter[R] {
	return func(yield func(R) bool) {
		it(func(v T) bool { return yield(fn(v)) })
	}
}

// Map applies fn to each key-value pair and returns a single-value sequence.
func (it LazyIter2[K, V]) Map[R any](fn func(K, V) R) LazyIter[R] {
	return func(yield func(R) bool) {
		it(func(k K, v V) bool { return yield(fn(k, v)) })
	}
}

// Map2 applies fn to each key-value pair while preserving a two-value
// sequence. fn may change both the key and value types.
func (it LazyIter2[K, V]) Map2[K2, V2 any](fn func(K, V) (K2, V2)) LazyIter2[K2, V2] {
	return func(yield func(K2, V2) bool) {
		it(func(k K, v V) bool { return yield(fn(k, v)) })
	}
}

// Reduce combines the elements into one value, starting with init.
func (it LazyIter[T]) Reduce[A any](fn func(A, T) A, init A) A {
	acc := init
	for v := range it {
		acc = fn(acc, v)
	}
	return acc
}

// Reduce combines the key-value pairs into one value, starting with init.
func (it LazyIter2[K, V]) Reduce[A any](fn func(A, K, V) A, init A) A {
	acc := init
	for k, v := range it {
		acc = fn(acc, k, v)
	}
	return acc
}

// Reverse returns the elements in reverse order. It buffers the input before
// yielding, so it is not streaming.
func (it LazyIter[T]) Reverse() LazyIter[T] {
	return func(yield func(T) bool) {
		var values []T
		for v := range it {
			values = append(values, v)
		}
		for i := len(values) - 1; i >= 0; i-- {
			if !yield(values[i]) {
				return
			}
		}
	}
}

// Reverse returns the key-value pairs in reverse order. It buffers the input
// before yielding, so it is not streaming.
func (it LazyIter2[K, V]) Reverse() LazyIter2[K, V] {
	type pair struct {
		key   K
		value V
	}
	return func(yield func(K, V) bool) {
		var values []pair
		for k, v := range it {
			values = append(values, pair{k, v})
		}
		for i := len(values) - 1; i >= 0; i-- {
			if !yield(values[i].key, values[i].value) {
				return
			}
		}
	}
}

// Skip returns the sequence after its first n elements. A non-positive n does
// not skip any elements.
func (it LazyIter[T]) Skip(n int) LazyIter[T] {
	return func(yield func(T) bool) {
		skipped := 0
		it(func(v T) bool {
			if skipped < n {
				skipped++
				return true
			}
			return yield(v)
		})
	}
}

// Skip returns the sequence after its first n key-value pairs. A non-positive
// n does not skip any pairs.
func (it LazyIter2[K, V]) Skip(n int) LazyIter2[K, V] {
	return func(yield func(K, V) bool) {
		skipped := 0
		it(func(k K, v V) bool {
			if skipped < n {
				skipped++
				return true
			}
			return yield(k, v)
		})
	}
}

// Take returns at most the first n elements. A non-positive n returns an empty
// sequence.
func (it LazyIter[T]) Take(n int) LazyIter[T] {
	return func(yield func(T) bool) {
		if n <= 0 {
			return
		}
		taken := 0
		it(func(v T) bool {
			if !yield(v) {
				return false
			}
			taken++
			return taken < n
		})
	}
}

// Take returns at most the first n key-value pairs. A non-positive n returns
// an empty sequence.
func (it LazyIter2[K, V]) Take(n int) LazyIter2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		taken := 0
		it(func(k K, v V) bool {
			if !yield(k, v) {
				return false
			}
			taken++
			return taken < n
		})
	}
}

// Zip pairs this sequence with other, stopping when either sequence ends.
func (it LazyIter[T]) Zip[U any](other LazyIter[U]) LazyIter2[T, U] {
	return func(yield func(T, U) bool) {
		next, stop := iter.Pull(it.Seq())
		defer stop()
		otherNext, otherStop := iter.Pull(other.Seq())
		defer otherStop()
		for {
			v, ok := next()
			u, otherOK := otherNext()
			if !ok || !otherOK || !yield(v, u) {
				return
			}
		}
	}
}
