package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	A int
	B int
}

func TestCopyMap(t *testing.T) {
	t.Run("not a map", func(t *testing.T) {
		_, err := CopyMap(1)
		assert.EqualError(t, err, "input is not a map")
	})

	t.Run("copy map and change", func(t *testing.T) {
		a := map[string]map[int]testStruct{"a": {1: {A: 1, B: 2}}}

		b, err := CopyMap(a)
		assert.NoError(t, err)
		assert.Equal(t, a, b)

		v, ok := b.(map[string]map[int]testStruct)
		assert.True(t, ok)
		s := v["a"][1]
		s.A = 2
		v["a"][1] = s
		assert.NotEqual(t, a, v)
	})
}

// Benchmarks

var benchmarkMap1 map[string]int
var benchmarkMap2 map[string]int
var benchmarkMap3 map[string]map[int]testStruct

func initBenchmarkMaps() {
	benchmarkMap1 = map[string]int{}
	benchmarkMap2 = map[string]int{"abc": 1, "adsa": 2, "sad": 3, "gfjiadsjf": 4, " ": 5}
	benchmarkMap3 = map[string]map[int]testStruct{"a": {1: {A: 1, B: 2}, 2: {A: 1, B: 2}, 3: {A: 1, B: 2}, 4: {A: 1, B: 2}, 5: {A: 1, B: 2}}, "2": {1: {A: 1, B: 2}, 2: {A: 1, B: 2}, 3: {A: 1, B: 2}, 4: {A: 1, B: 2}, 5: {A: 1, B: 2}}}

}

func BenchmarkCopyMap1(b *testing.B) {
	initBenchmarkMaps()

	for n := 0; n < b.N; n++ {
		CopyMap(benchmarkMap1)
	}
}

func BenchmarkCopyMap1NonGeneric(b *testing.B) {
	initBenchmarkMaps()

	f := func(m map[string]int) map[string]int {
		out := make(map[string]int)
		for k, v := range m {
			out[k] = v
		}
		return out
	}

	for n := 0; n < b.N; n++ {
		f(benchmarkMap1)
	}
}

func BenchmarkCopyMap2(b *testing.B) {
	initBenchmarkMaps()

	for n := 0; n < b.N; n++ {
		CopyMap(benchmarkMap2)
	}
}

func BenchmarkCopyMap2NonGeneric(b *testing.B) {
	initBenchmarkMaps()

	f := func(m map[string]int) map[string]int {
		out := make(map[string]int)
		for k, v := range m {
			out[k] = v
		}
		return out
	}

	for n := 0; n < b.N; n++ {
		f(benchmarkMap2)
	}
}

func BenchmarkCopyMap3(b *testing.B) {
	initBenchmarkMaps()

	for n := 0; n < b.N; n++ {
		CopyMap(benchmarkMap3)
	}
}

func BenchmarkCopyMap3NonGeneric(b *testing.B) {
	initBenchmarkMaps()

	f := func(m map[string]map[int]testStruct) map[string]map[int]testStruct {
		out := make(map[string]map[int]testStruct)
		for k, v := range m {
			out[k] = make(map[int]testStruct)
			for k2, v2 := range v {
				out[k][k2] = testStruct{A: v2.A, B: v2.B}
			}
		}
		return out
	}

	for n := 0; n < b.N; n++ {
		f(benchmarkMap3)
	}
}
