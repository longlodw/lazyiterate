# lazyiterate

A Go package providing generic, lazy if possible iterator utility functions for functional-style programming.
For full documentation, see the [GoDoc](https://pkg.go.dev/github.com/longlodw/lazyiterate).

## Features

- **All / Any**: Test if all or any elements (or key-value pairs) satisfy a predicate.
- **Count**: Count elements or key-value pairs in a sequence.
- **Find**: Find the first element or key-value pair matching a predicate.
- **Filter**: Lazily filter elements or key-value pairs.
- **Map**: Lazily transform elements or key-value pairs.
- **Reduce**: Accumulate elements or key-value pairs into a single value.
- **Reverse**: Iterate elements or key-value pairs in reverse order.
- **Skip / Take**: Skip or take a fixed number of elements or key-value pairs.
- **Zip**: Combine two sequences into pairs.

All functions are generic and work with the [iter](https://pkg.go.dev/iter) package's `Seq` and `Seq2` types.

## Example Usage

```go
import (
    "slices"

    "github.com/longlodw/lazyiterate"
)

int main() {
    seq := slices.Values([]int{1, 2, 3, 4, 5})
    even := lazyiterate.Filter(seq, func(x int) bool { return x%2 == 0 })
    doubled := lazyiterate.Map(even, func(x int) int { return x * 2 })
    count := lazyiterate.Count(doubled)
}
```

## Requirements

- Go 1.23+ (for iterator)
- [iter](https://pkg.go.dev/iter) package

## License

MIT
