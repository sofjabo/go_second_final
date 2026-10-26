package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE           = 100_000_000
	CHUNKS         = 8
	maxRandomValue = 1_000_000
)

func generateRandomElements(size int) []int {
	if size <= 0 {
		return nil
	}

	numbers := make([]int, size)
	for i := range numbers {
		numbers[i] = rand.Intn(maxRandomValue) + 1
	}

	return numbers
}

func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	max := data[0]
	for _, v := range data[1:] {
		if v > max {
			max = v
		}
	}

	return max
}

func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}

	if len(data) < CHUNKS {
		return maximum(data)
	}

	chunkSize := len(data) / CHUNKS
	maxes := make([]int, CHUNKS)

	var wg sync.WaitGroup
	wg.Add(CHUNKS)

	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == CHUNKS-1 {
			end = len(data)
		}
		go func(idx int, chunk []int) {
			defer wg.Done()
			maxes[idx] = maximum(chunk)
		}(i, data[start:end])
	}

	wg.Wait()
	return maximum(maxes)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	maxSingle := maximum(data)
	elapsedSingle := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d us\n", maxSingle, elapsedSingle)

	fmt.Printf("Ищем максимальное значение в %d горутин\n", CHUNKS)
	start = time.Now()
	maxParallel := maxChunks(data)
	elapsedParallel := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d us\n", maxParallel, elapsedParallel)

	if maxSingle != maxParallel {
		fmt.Printf("ОШИБКА: результаты не совпадают: single=%d, parallel=%d\n", maxSingle, maxParallel)
	}
}
