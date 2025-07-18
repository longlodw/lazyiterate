package lazyiterate_test

import (
	"maps"
	"reflect"
	"slices"
	"testing"

	"github.com/longlodw/lazyiterate.git"
)

func TestAll(t *testing.T) {
	seq := slices.Values([]int{2, 4, 6})
	if !lazyiterate.All(seq, func(v int) bool { return v%2 == 0 }) {
		t.Error("All failed for all even")
	}
	if lazyiterate.All(seq, func(v int) bool { return v > 2 }) {
		t.Error("All failed for not all > 2")
	}
}

func TestAll2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	seq := maps.All(m)
	if !lazyiterate.All2(seq, func(k int, v string) bool { return len(v) == 1 }) {
		t.Error("All2 failed for all len==1")
	}
	if lazyiterate.All2(seq, func(k int, v string) bool { return k > 1 }) {
		t.Error("All2 failed for not all k>1")
	}
}

func TestAny(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	if !lazyiterate.Any(seq, func(v int) bool { return v%2 == 0 }) {
		t.Error("Any failed for one even")
	}
	if lazyiterate.Any(seq, func(v int) bool { return v > 3 }) {
		t.Error("Any failed for none > 3")
	}
}

func TestAny2(t *testing.T) {
	m := map[int]string{1: "a", 2: "bb"}
	seq := maps.All(m)
	if !lazyiterate.Any2(seq, func(k int, v string) bool { return len(v) == 2 }) {
		t.Error("Any2 failed for one len==2")
	}
	if lazyiterate.Any2(seq, func(k int, v string) bool { return k > 2 }) {
		t.Error("Any2 failed for none k>2")
	}
}

func TestCount(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3, 4})
	if lazyiterate.Count(seq) != 4 {
		t.Error("Count failed")
	}
}

func TestCount2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	seq := maps.All(m)
	if lazyiterate.Count2(seq) != 2 {
		t.Error("Count2 failed")
	}
}

func TestFind(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	v, err := lazyiterate.Find(seq, func(x int) bool { return x == 2 })
	if v != 2 || err != nil {
		t.Error("Find failed to find 2")
	}
	_, err = lazyiterate.Find(seq, func(x int) bool { return x == 5 })
	if err == nil {
		t.Error("Find should not find 5")
	}
}

func TestFind2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	seq := maps.All(m)
	k, v, err := lazyiterate.Find2(seq, func(k int, v string) bool { return v == "b" })
	if v != "b" || k != 2 || err != nil {
		t.Error("Find2 failed to find b")
	}
	_, _, err = lazyiterate.Find2(seq, func(k int, v string) bool { return v == "c" })
	if err == nil {
		t.Error("Find2 should not find c")
	}
}

func TestFilter(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3, 4})
	var got []int
	lazyiterate.Filter(seq, func(v int) bool { return v%2 == 0 })(func(v int) bool {
		got = append(got, v)
		return true
	})
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter got %v, want %v", got, want)
	}
}

func TestFilter2(t *testing.T) {
	m := map[int]string{1: "a", 2: "bb"}
	seq := maps.All(m)
	var got []int
	lazyiterate.Filter2(seq, func(k int, v string) bool { return len(v) == 2 })(func(k int, v string) bool {
		got = append(got, k)
		return true
	})
	want := []int{2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter2 got %v, want %v", got, want)
	}
}

func TestMap(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	var got []string
	lazyiterate.Map(seq, func(v int) string { return string(rune('a' + v)) })(func(s string) bool {
		got = append(got, s)
		return true
	})
	want := []string{"b", "c", "d"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map got %v, want %v", got, want)
	}
}

func TestMap2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	seq := maps.All(m)
	var got []string
	lazyiterate.Map2(seq, func(k int, v string) string { return v + string(rune('0'+k)) })(func(s string) bool {
		got = append(got, s)
		return true
	})
	want := []string{"a1", "b2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Map2 got %v, want %v", got, want)
	}
}

func TestReduce(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	sum := lazyiterate.Reduce(seq, func(acc, v int) int { return acc + v }, 0)
	if sum != 6 {
		t.Errorf("Reduce got %v, want 6", sum)
	}
}

func TestReduce2(t *testing.T) {
	m := map[int]int{1: 2, 2: 3}
	seq := maps.All(m)
	sum := lazyiterate.Reduce2(seq, func(acc, k, v int) int { return acc + k + v }, 0)
	if sum != 1+2+2+3 {
		t.Errorf("Reduce2 got %v, want %v", sum, 1+2+2+3)
	}
}

func TestReverse(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3})
	var got []int
	lazyiterate.Reverse(seq)(func(v int) bool {
		got = append(got, v)
		return true
	})
	want := []int{3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse got %v, want %v", got, want)
	}
}

func TestReverse2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b"}
	seq := maps.All(m)
	var got []int
	lazyiterate.Reverse2(seq)(func(k int, v string) bool {
		got = append(got, k)
		return true
	})
	if len(got) != 2 || (got[0] != 2 && got[0] != 1) {
		t.Errorf("Reverse2 got %v, want [2,1] or [1,2]", got)
	}
}

func TestSkip(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3, 4})
	var got []int
	lazyiterate.Skip(seq, 2)(func(v int) bool {
		got = append(got, v)
		return true
	})
	want := []int{3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Skip got %v, want %v", got, want)
	}
}

func TestSkip2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	seq := maps.All(m)
	var got []int
	lazyiterate.Skip2(seq, 2)(func(k int, v string) bool {
		got = append(got, k)
		return true
	})
	if len(got) != 1 {
		t.Errorf("Skip2 got %v, want 1 element", got)
	}
}

func TestTake(t *testing.T) {
	seq := slices.Values([]int{1, 2, 3, 4})
	var got []int
	lazyiterate.Take(seq, 2)(func(v int) bool {
		got = append(got, v)
		return true
	})
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Take got %v, want %v", got, want)
	}
}

func TestTake2(t *testing.T) {
	m := map[int]string{1: "a", 2: "b", 3: "c"}
	seq := maps.All(m)
	var got []int
	lazyiterate.Take2(seq, 2)(func(k int, v string) bool {
		got = append(got, k)
		return true
	})
	if len(got) != 2 {
		t.Errorf("Take2 got %v, want 2 elements", got)
	}
}

func TestZip(t *testing.T) {
	seq1 := slices.Values([]int{1, 2, 3, 4})
	seq2 := slices.Values([]string{"a", "b", "c"})
	var got []struct {
		Int    int
		String string
	}
	lazyiterate.Zip(seq1, seq2)(func(i int, s string) bool {
		got = append(got, struct {
			Int    int
			String string
		}{i, s})
		return true
	})
	want := []struct {
		Int    int
		String string
	}{{1, "a"}, {2, "b"}, {3, "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Zip got %v, want %v", got, want)
	}
}
