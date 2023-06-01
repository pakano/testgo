package main

import (
	"container/heap"
	"fmt"
	"testing"
)

type Queue []int

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i, j int) bool {
	if i <= 0 || j <= 0 {
		return false
	}
	return q[i] > q[j]
}

func (q Queue) Swap(i, j int) {
	if i <= 0 || j <= 0 {
		return
	}
	q[i], q[j] = q[j], q[i]
}

func (q *Queue) Push(x interface{}) {
	*q = append(*q, x.(int))
}

func (q *Queue) Pop() interface{} {
	l := len(*q)
	if l <= 1 {
		return nil
	}
	x := (*q)[l-1]
	*q = (*q)[0 : l-1]
	return x
}

func TestHeap(t *testing.T) {
	q := &Queue{1}
	heap.Init(q)
	fmt.Println(q)
	heap.Push(q, 9)
	fmt.Println(q)
	heap.Pop(q)
	fmt.Println(q)
	heap.Pop(q)
	fmt.Println(q)
	heap.Pop(q)
	fmt.Println(q)
	heap.Pop(q)
	fmt.Println(q)
	heap.Fix()
}
