package main

import (
	"math"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"zero size", 0},
		{"negative size", -5},
		{"one element", 1},
		{"hundred elements", 100},
		{"ten thousand elements", 10_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRandomElements(tt.size)

			if tt.size <= 0 {
				if len(got) != 0 {
					t.Fatalf("generateRandomElements(%d) len = %d, want 0", tt.size, len(got))
				}
				return
			}

			if len(got) != tt.size {
				t.Fatalf("generateRandomElements(%d) len = %d, want %d", tt.size, len(got), tt.size)
			}

			for i, v := range got {
				if v <= 0 || v > maxRandomValue {
					t.Errorf("generateRandomElements(%d)[%d] = %d, want in [1, %d]", tt.size, i, v, maxRandomValue)
				}
			}
		})
	}

	t.Run("different calls produce different data", func(t *testing.T) {
		a := generateRandomElements(100)
		b := generateRandomElements(100)

		if len(a) != len(b) {
			t.Fatalf("lengths differ: %d vs %d", len(a), len(b))
		}

		same := true
		for i := range a {
			if a[i] != b[i] {
				same = false
				break
			}
		}
		if same {
			t.Error("two consecutive calls returned identical slices, generator may be broken")
		}
	})

	t.Run("produces varied values", func(t *testing.T) {
		data := generateRandomElements(1000)
		unique := make(map[int]struct{}, len(data))
		for _, v := range data {
			unique[v] = struct{}{}
		}
		if len(unique) < 10 {
			t.Errorf("only %d unique values in 1000 elements, want at least 10", len(unique))
		}
	})
}

func TestMaximum(t *testing.T) {
	largeAsc := make([]int, 10_000)
	for i := range largeAsc {
		largeAsc[i] = i
	}

	tests := []struct {
		name string
		data []int
		want int
	}{
		{"nil", nil, 0},
		{"empty", []int{}, 0},
		{"single", []int{42}, 42},
		{"all equal", []int{5, 5, 5}, 5},
		{"max at start", []int{10, 1, 2, 3}, 10},
		{"max in middle", []int{1, 2, 10, 3}, 10},
		{"max at end", []int{1, 2, 3, 10}, 10},
		{"duplicate max", []int{7, 3, 7, 1}, 7},
		{"negative values", []int{-5, -1, -10}, -1},
		{"max int", []int{math.MinInt, math.MaxInt, 0}, math.MaxInt},
		{"min int and negatives", []int{math.MinInt, -1, -2}, -1},
		{"large ascending", largeAsc, 9_999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maximum(tt.data); got != tt.want {
				t.Errorf("maximum(%v) = %d, want %d", tt.data, got, tt.want)
			}
		})
	}
}

func TestMaxChunks(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"nil", 0},
		{"smaller than chunks", CHUNKS - 1},
		{"equal to chunks", CHUNKS},
		{"chunks plus one", CHUNKS + 1},
		{"multiple of chunks", CHUNKS * 2},
		{"with remainder", CHUNKS*2 + 1},
		{"large remainder", CHUNKS*3 + 5},
		{"large", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generateRandomElements(tt.size)
			gotChunks := maxChunks(data)
			gotSingle := maximum(data)
			if gotChunks != gotSingle {
				t.Errorf("maxChunks(%d elems) = %d, maximum = %d, want equal",
					tt.size, gotChunks, gotSingle)
			}
		})
	}

	t.Run("max in tail after division remainder", func(t *testing.T) {
		size := CHUNKS*2 + 3
		data := make([]int, size)
		for i := range data {
			data[i] = 1
		}
		data[size-1] = 999

		if got := maxChunks(data); got != 999 {
			t.Errorf("maxChunks = %d, want 999 (max in tail)", got)
		}
	})

	t.Run("all equal", func(t *testing.T) {
		data := make([]int, 1000)
		for i := range data {
			data[i] = 42
		}
		if got := maxChunks(data); got != 42 {
			t.Errorf("maxChunks = %d, want 42", got)
		}
	})

	t.Run("local maxima in each chunk", func(t *testing.T) {
		size := CHUNKS * 10
		data := make([]int, size)
		for i := range data {
			data[i] = 1
		}
		chunkSize := size / CHUNKS
		for c := 0; c < CHUNKS; c++ {
			data[c*chunkSize] = 100 + c
		}
		if got := maxChunks(data); got != 107 {
			t.Errorf("maxChunks = %d, want 107", got)
		}
	})
}

func TestMaxChunksStress(t *testing.T) {
	for i := 0; i < 50; i++ {
		data := generateRandomElements(1000)
		got := maxChunks(data)
		want := maximum(data)
		if got != want {
			t.Fatalf("iteration %d: maxChunks = %d, maximum = %d", i, got, want)
		}
	}
}

func BenchmarkMaximum(b *testing.B) {
	data := generateRandomElements(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maximum(data)
	}
}

func BenchmarkMaxChunks(b *testing.B) {
	data := generateRandomElements(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxChunks(data)
	}
}
