package lazyiterate_test

import (
	"errors"
	"maps"
	"slices"
	"testing"

	"github.com/longlodw/lazyiterate"
)

func TestLazyIterChain(t *testing.T) {
	got := slices.Collect(lazyiterate.From(slices.Values([]int{1, 2, 3, 4, 5})).
		Filter(func(v int) bool { return v%2 == 1 }).
		Map(func(v int) string { return string(rune('a' + v - 1)) }).
		Skip(1).Take(2).Reverse().Seq())
	want := []string{"e", "c"}
	if !slices.Equal(got, want) {
		t.Errorf("chain result = %v, want %v", got, want)
	}
}

func TestLazyIterTerminalOperations(t *testing.T) {
	it := lazyiterate.From(slices.Values([]int{2, 4, 6}))
	if !it.All(func(v int) bool { return v%2 == 0 }) || it.Any(func(v int) bool { return v%2 == 1 }) {
		t.Fatal("unexpected predicate result")
	}
	if got := it.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}
	if got := it.Reduce(func(acc, v int) int { return acc + v }, 0); got != 12 {
		t.Errorf("Reduce() = %d, want 12", got)
	}
	v, err := it.Find(func(v int) bool { return v == 4 })
	if err != nil || v != 4 {
		t.Errorf("Find() = %d, %v; want 4, nil", v, err)
	}
	_, err = it.Find(func(v int) bool { return v == 7 })
	if !errors.Is(err, lazyiterate.ErrNotFound) {
		t.Errorf("Find() error = %v, want ErrNotFound", err)
	}
}

func TestLazyIter2ChainAndTerminals(t *testing.T) {
	it := lazyiterate.From2(maps.All(map[int]string{1: "a", 2: "bb", 3: "ccc"}))
	if !it.All(func(k int, v string) bool { return k == len(v) }) {
		t.Fatal("All() = false, want true")
	}
	if !it.Any(func(k int, v string) bool { return k == 2 && v == "bb" }) {
		t.Fatal("Any() = false, want true")
	}
	if got := it.Count(); got != 3 {
		t.Errorf("Count() = %d, want 3", got)
	}
	key, value, err := it.Find(func(k int, _ string) bool { return k == 2 })
	if err != nil || key != 2 || value != "bb" {
		t.Errorf("Find() = %d, %q, %v", key, value, err)
	}
	got := slices.Collect(it.Filter(func(k int, _ string) bool { return k > 1 }).Map(func(_ int, v string) string { return v }).Seq())
	if !slices.Equal(slices.Sorted(slices.Values(got)), []string{"bb", "ccc"}) {
		t.Errorf("Map() = %v, want [bb ccc]", got)
	}

	gotPairs := make(map[string]int)
	for key, value := range it.Filter(func(k int, _ string) bool { return k > 1 }).Map2(func(k int, v string) (string, int) { return v, k * 10 }) {
		gotPairs[key] = value
	}
	want := map[string]int{"bb": 20, "ccc": 30}
	if !maps.Equal(gotPairs, want) {
		t.Errorf("Map2() = %v, want %v", gotPairs, want)
	}
	if got := it.Reduce(func(acc, k int, _ string) int { return acc + k }, 0); got != 6 {
		t.Errorf("Reduce() = %d, want 6", got)
	}

	pairsAfterReverse := make([][2]int, 0)
	for key, value := range it.Skip(1).Take(1).Reverse() {
		pairsAfterReverse = append(pairsAfterReverse, [2]int{key, len(value)})
	}
	if len(pairsAfterReverse) != 1 {
		t.Errorf("Skip().Take().Reverse() = %v, want one pair", pairsAfterReverse)
	}
}

func TestZipStopsAtShortestSequence(t *testing.T) {
	got := make([][2]string, 0)
	for number, letter := range lazyiterate.From(slices.Values([]int{1, 2, 3})).Zip(lazyiterate.From(slices.Values([]string{"a", "b"}))) {
		got = append(got, [2]string{string(rune('0' + number)), letter})
	}
	want := [][2]string{{"1", "a"}, {"2", "b"}}
	if !slices.Equal(got, want) {
		t.Errorf("Zip() = %v, want %v", got, want)
	}
}

func TestTakeWithNonPositiveCountDoesNotReadInput(t *testing.T) {
	called := false
	it := lazyiterate.From(func(yield func(int) bool) { called = true; yield(1) })
	if got := slices.Collect(it.Take(0).Seq()); len(got) != 0 || called {
		t.Errorf("Take(0) = %v, input called = %t; want empty without input", got, called)
	}
}

func TestTakeDoesNotReadPastLimit(t *testing.T) {
	produced := 0
	it := lazyiterate.From(func(yield func(int) bool) {
		for _, v := range []int{1, 2, 3} {
			produced++
			if !yield(v) {
				return
			}
		}
	})
	got := slices.Collect(it.Take(2).Seq())
	if !slices.Equal(got, []int{1, 2}) || produced != 2 {
		t.Errorf("Take(2) = %v after producing %d values; want [1 2] after 2", got, produced)
	}
}
