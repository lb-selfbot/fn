package fn

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestClamp(t *testing.T) {
	if Clamp(5, 0, 10) != 5 {
		t.Fatalf("expected 5")
	}
	if Clamp(-1, 0, 10) != 0 {
		t.Fatalf("below min not clamped")
	}
	if Clamp(11, 0, 10) != 10 {
		t.Fatalf("above max not clamped")
	}
}

func TestLimit(t *testing.T) {
	data := []int{1, 2, 3}
	if got := Limit(data, 2); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("limit mismatch: %v", got)
	}
	if got := Limit(data, 5); !reflect.DeepEqual(got, data) {
		t.Fatalf("limit beyond len mismatch: %v", got)
	}
	if got := Limit([]int{}, 3); len(got) != 0 {
		t.Fatalf("expected empty slice")
	}
}

func TestMap(t *testing.T) {
	data := []int{1, 2, 3}
	res := Map(data, func(v int) int { return v * v })
	if !reflect.DeepEqual(res, []int{1, 4, 9}) {
		t.Fatalf("map mismatch: %v", res)
	}
	// original unchanged
	if !reflect.DeepEqual(data, []int{1, 2, 3}) {
		t.Fatalf("original slice modified")
	}
}

func TestMapIndexed(t *testing.T) {
	data := []string{"a", "b", "c"}
	res := MapIndexed(data, func(i int, v string) string { return v + string(rune('0'+i)) })
	if !reflect.DeepEqual(res, []string{"a0", "b1", "c2"}) {
		t.Fatalf("mapindexed mismatch: %v", res)
	}
}

func TestFilter(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	res := Filter(data, func(v int) bool { return v%2 == 0 })
	if !reflect.DeepEqual(res, []int{2, 4}) {
		t.Fatalf("filter mismatch: %v", res)
	}
	// Ensure original unchanged
	if !reflect.DeepEqual(data, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("original slice modified")
	}
}

func TestFilterInPlace(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	filtered := FilterInPlace(data, func(v int) bool { return v > 2 })
	if !reflect.DeepEqual(filtered, []int{3, 4, 5}) {
		t.Fatalf("filter in place mismatch: %v", filtered)
	}
	// Capacity should be same, underlying slice reused
	if cap(filtered) != cap(data) {
		t.Fatalf("expected reuse of underlying array")
	}
}

func TestReduce(t *testing.T) {
	data := []int{1, 2, 3, 4}
	sum := Reduce(data, 0, func(acc, v int) int { return acc + v })
	if sum != 10 {
		t.Fatalf("expected sum 10 got %d", sum)
	}
	prod := Reduce(data, 1, func(acc, v int) int { return acc * v })
	if prod != 24 {
		t.Fatalf("expected product 24 got %d", prod)
	}
}

func TestAnyAll(t *testing.T) {
	data := []int{1, 3, 5}
	if Any(data, func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("expected no even numbers")
	}
	if !All(data, func(v int) bool { return v%2 == 1 }) {
		t.Fatalf("expected all odd numbers")
	}
	// Add even number for Any
	data = append(data, 4)
	if !Any(data, func(v int) bool { return v%2 == 0 }) {
		t.Fatalf("expected even number present")
	}
	if All(data, func(v int) bool { return v%2 == 1 }) {
		t.Fatalf("no longer all odd")
	}
}

func TestUnique(t *testing.T) {
	cases := [][]int{{}, {1}, {1, 1, 1}, {1, 2, 2, 3, 3, 3, 4}, {1, 2, 3, 4}}
	expects := [][]int{{}, {1}, {1}, {1, 2, 3, 4}, {1, 2, 3, 4}}
	for i, c := range cases {
		got := Unique(c)
		if !reflect.DeepEqual(got, expects[i]) {
			t.Fatalf("unique mismatch case %d: %v != %v", i, got, expects[i])
		}
	}
}

func TestIfElse(t *testing.T) {
	if IfElse(true, 1, 2) != 1 {
		t.Fatalf("expected 1")
	}
	if IfElse(false, 1, 2) != 2 {
		t.Fatalf("expected 2")
	}
}

func TestReverse(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	Reverse(data)
	if !reflect.DeepEqual(data, []int{5, 4, 3, 2, 1}) {
		t.Fatalf("reverse mismatch: %v", data)
	}
	// single element and empty
	one := []int{42}
	Reverse(one)
	if !reflect.DeepEqual(one, []int{42}) {
		t.Fatalf("reverse single mismatch")
	}
	var empty []int
	Reverse(empty)
}

func TestShuffle(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	copyData := append([]int(nil), data...)
	rand.Seed(1)
	Shuffle(data)
	if len(data) != len(copyData) {
		t.Fatalf("length changed after shuffle")
	}
	// after shuffling with deterministic seed, order should differ for this size, but allow rare equality
	identical := true
	for i := range data {
		if data[i] != copyData[i] {
			identical = false
			break
		}
	}
	if identical {
		t.Log("shuffle produced identical order (rare)")
	}
	// ensure all elements still present via sum & map
	m := map[int]int{}
	for _, v := range data {
		m[v]++
	}
	for _, v := range copyData {
		if m[v] != 1 {
			t.Fatalf("element count mismatch for %d", v)
		}
	}
}

func TestBatch(t *testing.T) {
	// typical
	data := []int{1, 2, 3, 4, 5, 6, 7}
	batches := Batch(data, 3)
	if len(batches) != 3 {
		t.Fatalf("expected 3 batches got %d", len(batches))
	}
	expectedLens := []int{3, 3, 1}
	for i, b := range batches {
		if len(b) != expectedLens[i] {
			t.Fatalf("batch %d size mismatch: %d", i, len(b))
		}
	}
	// batch size larger than slice
	batches = Batch(data, 10)
	if len(batches) != 1 || !reflect.DeepEqual(batches[0], data) {
		t.Fatalf("expected single batch containing all data")
	}

	// batch size 0 -> expect nil
	if res := Batch(data, 0); res != nil {
		// differentiate between empty and nil
		if len(res) == 0 {
			// acceptable? requirement states should return nil explicitly
			// enforce nil for clarity
			t.Fatalf("expected nil for batch size 0, got empty slice")
		} else {
			t.Fatalf("expected nil for batch size 0, got %v", res)
		}
	}
}

func TestFirst(t *testing.T) {
	data := []int{5, 7, 9, 10}
	if v, ok := First(data, func(x int) bool { return x%2 == 0 }); !ok || v != 10 {
		t.Fatalf("first even mismatch: %v %v", v, ok)
	}
	if v, ok := First(data, func(x int) bool { return x < 0 }); ok || v != 0 {
		t.Fatalf("expected no match and zero value")
	}
}

func TestDelete(t *testing.T) {
	data := []int{1, 2, 3, 2, 4}
	res := Delete(data, 2)
	// slices.DeleteFunc removes all matches
	if !reflect.DeepEqual(res, []int{1, 3, 4}) {
		t.Fatalf("delete mismatch: %v", res)
	}
	// ensure original underlying mutated appropriately (length shrink)
	if len(res) != 3 {
		t.Fatalf("unexpected length after delete")
	}
}

func TestToIfaceSlice(t *testing.T) {
	res := ToIfaceSlice(1, "a", true)
	if len(res) != 3 {
		t.Fatalf("length mismatch")
	}
	if _, ok := res[1].(string); !ok {
		t.Fatalf("expected string at index 1")
	}
	// empty call
	empty := ToIfaceSlice()
	if len(empty) != 0 {
		t.Fatalf("expected empty slice")
	}
}
