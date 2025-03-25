package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	list := []int{17, 45, 9, 33, 158, 1488, 7, 0}
	sorted := concurrentMergeSort(list)
	fmt.Println(sorted)

	ctx1, cancel1 := context.WithTimeout(context.WithValue(context.Background(), "name", "Арсен"), 2*time.Second)
	defer cancel1()

	ctx2, cancel2 := context.WithTimeout(context.WithValue(context.Background(), "quality", "Не жирный"), 5*time.Second)
	defer cancel2()

	mergedCtx := mergeContexts(ctx1, ctx2)

	fmt.Println("Пользователь:", mergedCtx.Value("name"))
	fmt.Println("Качество:", mergedCtx.Value("quality"))

	select {
	case <-mergedCtx.Done():
		fmt.Println("❌ Объединённый контекст завершён:", mergedCtx.Err())
	case <-time.After(6 * time.Second):
		fmt.Println("✅ Контекст всё ещё активен")
	}
}

func concurrentMergeSort(list []int) []int {
	n := len(list)
	if n < 2 {
		return list
	}

	mid := n / 2
	leftPart := list[:mid]
	rightPart := list[mid:]

	var leftResult, rightResult []int
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		leftResult = concurrentMergeSort(leftPart)
	}()

	go func() {
		defer wg.Done()
		rightResult = concurrentMergeSort(rightPart)
	}()

	wg.Wait()

	return merge(leftResult, rightResult)
}

func merge(left []int, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)

	return result
}

func mergeContexts(ctx1, ctx2 context.Context) context.Context {
	mergedCtx, cancel := context.WithCancel(context.Background())

	go func() {
		select {
		case <-ctx1.Done():
		case <-ctx2.Done():
		}
		cancel()
	}()

	return &mergedContext{mergedCtx, ctx1, ctx2}
}

type mergedContext struct {
	context.Context
	ctx1, ctx2 context.Context
}

func (m *mergedContext) Value(key any) any {
	if v := m.ctx1.Value(key); v != nil {
		return v
	}
	return m.ctx2.Value(key)
}
