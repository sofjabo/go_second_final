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
	sizeTests := []struct {
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

	for _, tt := range sizeTests {
		t.Run("differential/"+tt.name, func(t *testing.T) {
			data := generateRandomElements(tt.size)
			gotChunks := maxChunks(data)
			gotSingle := maximum(data)
			if gotChunks != gotSingle {
				t.Errorf("maxChunks(%d elems) = %d, maximum = %d, want equal",
					tt.size, gotChunks, gotSingle)
			}
		})
	}

	deterministic := []struct {
		name string
		data []int
		want int
	}{
		{
			name: "all negative",
			data: []int{-10, -5, -3, -20, -1, -8, -7, -6},
			want: -1,
		},
		{
			name: "all min int",
			data: []int{
				math.MinInt, math.MinInt, math.MinInt, math.MinInt,
				math.MinInt, math.MinInt, math.MinInt, math.MinInt,
			},
			want: math.MinInt,
		},
		{
			name: "max int in last chunk",
			data: []int{
				-10, -20, -30, -40,
				-50, -60, -70, math.MaxInt,
			},
			want: math.MaxInt,
		},
		{
			name: "max in tail after division remainder",
			data: func() []int {
				d := make([]int, CHUNKS*2+3)
				for i := range d {
					d[i] = 1
				}
				d[len(d)-1] = 999
				return d
			}(),
			want: 999,
		},
		{
			name: "all equal",
			data: func() []int {
				d := make([]int, 1000)
				for i := range d {
					d[i] = 42
				}
				return d
			}(),
			want: 42,
		},
		{
			name: "local maxima at chunk starts",
			data: func() []int {
				size := CHUNKS * 10
				d := make([]int, size)
				for i := range d {
					d[i] = 1
				}
				chunkSize := size / CHUNKS
				for c := 0; c < CHUNKS; c++ {
					d[c*chunkSize] = 100 + c
				}
				return d
			}(),
			want: 107,
		},
		{
			name: "max at chunk boundaries",
			data: func() []int {
				size := CHUNKS * 4
				d := make([]int, size)
				for i := range d {
					d[i] = 1
				}
				chunkSize := size / CHUNKS
				d[chunkSize] = 777
				return d
			}(),
			want: 777,
		},
	}

	for _, tt := range deterministic {
		t.Run("deterministic/"+tt.name, func(t *testing.T) {
			if got := maxChunks(tt.data); got != tt.want {
				t.Errorf("maxChunks() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSearchFunctionsDoNotModifyInput(t *testing.T) {
	original := []int{3, 10, -5, 7, 1, 42, -100, 8, 9, 9, 9}
	cases := []struct {
		name string
		fn   func([]int) int
	}{
		{"maximum", maximum},
		{"maxChunks", maxChunks},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data := append([]int(nil), original...)
			c.fn(data)
			for i := range data {
				if data[i] != original[i] {
					t.Fatalf("%s modified input at index %d: got %d, want %d",
						c.name, i, data[i], original[i])
				}
			}
		})
	}
}

func TestMaxChunksRepeated(t *testing.T) {
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
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}
	want := len(data) - 1

	b.ResetTimer()
	var result int
	for i := 0; i < b.N; i++ {
		result = maximum(data)
	}
	b.StopTimer()

	if result != want {
		b.Fatalf("maximum() = %d, want %d", result, want)
	}
}

func BenchmarkMaxChunks(b *testing.B) {
	data := make([]int, 1_000_000)
	for i := range data {
		data[i] = i
	}
	want := len(data) - 1

	b.ResetTimer()
	var result int
	for i := 0; i < b.N; i++ {
		result = maxChunks(data)
	}
	b.StopTimer()

	if result != want {
		b.Fatalf("maxChunks() = %d, want %d", result, want)
	}
}
