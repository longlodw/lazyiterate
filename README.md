# lazyiterate

A Go 1.27 package for lazy, functional-style iterator pipelines. Wrap an
[`iter.Seq`](https://pkg.go.dev/iter#Seq) or
[`iter.Seq2`](https://pkg.go.dev/iter#Seq2) in a `LazyIter` object, then chain
operations directly on it.

For full API documentation, see the [Go package documentation](https://pkg.go.dev/github.com/longlodw/lazyiterate).

## Features

- **Mapping shapes**: `LazyIter2.Map` projects pairs to `LazyIter`; `Map2` preserves pairs as `LazyIter2`.
- **Terminal operations**: `All`, `Any`, `Count`, `Find`, and `Reduce` consume a sequence and return a result.
- **Single- and two-value sequences**: `LazyIter[T]` wraps `iter.Seq[T]`; `LazyIter2[K, V]` wraps `iter.Seq2[K, V]`.

## Example

```go
import (
	"slices"

	"github.com/longlodw/lazyiterate"
)

func main() {
	count := lazyiterate.From(slices.Values([]int{1, 2, 3, 4, 5})).
		Filter(func(x int) bool { return x%2 == 0 }).
		Map(func(x int) int { return x * 2 }).
		Count()
}
```

`From` and `From2` infer their type arguments from an existing `iter.Seq` or
`iter.Seq2`. You can also convert explicitly, for example
`lazyiterate.LazyIter[int](seq)`.

`Map` on `LazyIter2` projects each pair to a single value. Use `Map2` when a
transformation needs to preserve both values:

```go
renamed := lazyiterate.From2(maps.All(users)).Map2(func(id int, name string) (string, int) {
	return name, id
})
```

`Reverse` buffers its full input, so unlike the other sequence-producing
operations it is lazy to consume but not streaming to produce.

## Requirements

- Go 1.27+ (generic methods and the standard-library `iter` package)

## License

MIT
