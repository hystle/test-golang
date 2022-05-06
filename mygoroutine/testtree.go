package testGoroutinePkg

import (
	"fmt"
	"math/rand"
)

type Tree struct {
    Left  *Tree
    Value int
    Right *Tree
}

func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	s:= ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}

func New(k int) *Tree {
	var t *Tree
	for _, v := range rand.Perm(k) {
		t = insert(t, v+1)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{
			Left: nil,
			Value: v,
			Right: nil,	
		}
	}
	if t.Value > v {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func Walk(t *Tree, ch chan int) {
    defer close(ch)
	var w func(t *Tree)
	w = func(t *Tree) {
		if t != nil {
			w(t.Left)
			ch <- t.Value
			w(t.Right)
		}
	}
	w(t)
}

func Same(t1, t2 *Tree) bool {
	ch1 := make(chan int)
	go Walk(t1, ch1)

	ch2 := make(chan int)
	go Walk(t2, ch2)

	// // approach1
	// for v := range ch1 {
	// 	if v != <-ch2 {
	// 		return false
	// 	}
	// }
	// return true

	// approach2
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2

		if v1 != v2 || ok1 != ok2 {
			return false
		}
		if !ok1 {
			break
		}
	}
	return true
}