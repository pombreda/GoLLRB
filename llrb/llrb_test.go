// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llrb

import (
	"math"
	"math/rand"
	"testing"
)

func TestCases(t *testing.T) {
	tree := New()
	tree.ReplaceOrInsert(Int(1))
	tree.ReplaceOrInsert(Int(1))
	if tree.Len() != 1 {
		t.Errorf("expecting len 1")
	}
	if !tree.Has(Int(1)) {
		t.Errorf("expecting to find key=1")
	}

	tree.Delete(Int(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(Int(1)) {
		t.Errorf("not expecting to find key=1")
	}

	tree.Delete(Int(1))
	if tree.Len() != 0 {
		t.Errorf("expecting len 0")
	}
	if tree.Has(Int(1)) {
		t.Errorf("not expecting to find key=1")
	}
}

func TestReverseInsertOrder(t *testing.T) {
	tree := New()
	n := 100
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(n - i))
	}
	i := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		i++
		if item.(Int) != Int(i) {
			t.Errorf("bad order: got %d, expect %d", item.(Int), i)
		}
		return true
	})
}

func TestRange(t *testing.T) {
	tree := New()
	order := []String{
		"ab", "aba", "abc", "a", "aa", "aaa", "b", "a-", "a!",
	}
	for _, i := range order {
		tree.ReplaceOrInsert(i)
	}
	k := 0
	tree.AscendRange(String("ab"), String("ac"), func(item Item) bool {
		if k > 3 {
			t.Fatalf("returned more items than expected")
		}
		i1 := order[k]
		i2 := item.(String)
		if i1 != i2 {
			t.Errorf("expecting %s, got %s", i1, i2)
		}
		k++
		return true
	})
}

func TestRandomInsertOrder(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		if item.(Int) != Int(j) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestRandomReplace(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	perm = rand.Perm(n)
	for i := 0; i < n; i++ {
		if replaced := tree.ReplaceOrInsert(Int(perm[i])); replaced == nil || replaced.(Int) != Int(perm[i]) {
			t.Errorf("error replacing")
		}
	}
}

func TestRandomInsertSequentialDelete(t *testing.T) {
	tree := New()
	n := 1000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	for i := 0; i < n; i++ {
		tree.Delete(Int(i))
	}
}

func TestRandomInsertDeleteNonExistent(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	if tree.Delete(Int(200)) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(Int(-2)) != nil {
		t.Errorf("deleted non-existent item")
	}
	for i := 0; i < n; i++ {
		if u := tree.Delete(Int(i)); u == nil || u.(Int) != Int(i) {
			t.Errorf("delete failed")
		}
	}
	if tree.Delete(Int(200)) != nil {
		t.Errorf("deleted non-existent item")
	}
	if tree.Delete(Int(-2)) != nil {
		t.Errorf("deleted non-existent item")
	}
}

func TestRandomInsertPartialDeleteOrder(t *testing.T) {
	tree := New()
	n := 100
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	for i := 1; i < n-1; i++ {
		tree.Delete(Int(i))
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		switch j {
		case 0:
			if item.(Int) != Int(0) {
				t.Errorf("expecting 0")
			}
		case 1:
			if item.(Int) != Int(n-1) {
				t.Errorf("expecting %d", n-1)
			}
		}
		j++
		return true
	})
}

func TestRandomInsertStats(t *testing.T) {
	tree := New()
	n := 100000
	perm := rand.Perm(n)
	for i := 0; i < n; i++ {
		tree.ReplaceOrInsert(Int(perm[i]))
	}
	avg, _ := tree.HeightStats()
	expAvg := math.Log2(float64(n)) - 1.5
	if math.Abs(avg-expAvg) >= 2.0 {
		t.Errorf("too much deviation from expected average height")
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(Int(i))
	}
}

func BenchmarkDeleteMin(b *testing.B) {
	b.StopTimer()
	tree := New()
	for i := 0; i < b.N; i++ {
		tree.ReplaceOrInsert(Int(b.N - i))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.DeleteMin()
	}
}

func TestInsertNoReplace(t *testing.T) {
	tree := New()
	n := 1000
	for q := 0; q < 2; q++ {
		perm := rand.Perm(n)
		for i := 0; i < n; i++ {
			tree.InsertNoReplace(Int(perm[i]))
		}
	}
	j := 0
	tree.AscendGreaterOrEqual(Int(0), func(item Item) bool {
		if item.(Int) != Int(j/2) {
			t.Fatalf("bad order")
		}
		j++
		return true
	})
}

func TestCoW(t *testing.T) {
	tree := NewCoW()
	tree.InsertNoReplace(Int(4))
	tree.InsertNoReplace(Int(6))
	tree.InsertNoReplace(Int(1))
	tree.InsertNoReplace(Int(3))

	tree2 := tree.Clone()
	tree.InsertNoReplace(Int(60))
	orig := tree.ReplaceOrInsert(Int(4))

	if tree2.Len() != 4 {
		t.Errorf("Expected Len of 4, but got %d", tree2.Len())
	}
	if tree2.Get(Int(60)) != nil {
		t.Errorf("Expected to get nil for 60")
	}
	if tree2.Get(Int(4)) != orig {
		t.Errorf("Expected to get original for 4")
	}

	tree2.InsertNoReplace(Int(20))
	if tree.Len() != 5 {
		t.Errorf("Expected Len of 5, but got %d", tree.Len())
	}
	if tree.Get(Int(20)) != nil {
		t.Errorf("Expected to get nil for 20")
	}

	orig = tree.Delete(Int(6))
	tree.Delete(Int(1))
	tree.DeleteMax()
	tree.DeleteMin()
	if tree2.Get(Int(6)) != orig {
		t.Errorf("Expected to get original for 6")
	}

	i := 0
	expect := []Int{4, 60}
	tree.AscendRange(Int(-1), Int(100), func(itm Item) bool {
		iv := itm.(Int)
		if iv != expect[i] {
			t.Errorf("expected %d got %d", expect[i], iv)
		}
		i++
		return true
	})

	i = 0
	expect = []Int{1, 3, 4, 6, 20}
	tree2.AscendRange(Int(-1), Int(100), func(itm Item) bool {
		iv := itm.(Int)
		if iv != expect[i] {
			t.Errorf("expected %d got %d", expect[i], iv)
		}
		i++
		return true
	})
}

type Intp int

func (x *Intp) Less(than Item) bool {
	return *x < *than.(*Intp)
}

func TestRandomCoW(t *testing.T) {
	tree := NewCoW()
	tree2 := tree.Clone()

	n := 1000
	nd := 50
	perm := rand.Perm(n)
	dels := make(map[*Intp]struct{}, nd)

	test := func(phase string, tree *LLRB) {
		prev := tree.Min().(*Intp)
		bottom := 0
		tree.AscendGreaterOrEqual((*Intp)(&bottom), func(item Item) bool {
			itm := item.(*Intp)
			if _, deleted := dels[itm]; !deleted {
				if *itm < *prev {
					t.Log(tree.Root())
					t.Fatalf("%s bad order: %v (expected %v)", phase, itm, prev)
				}
			}
			prev = itm
			return true
		})
	}

	for _, p := range perm {
		p := p
		itm := (*Intp)(&p)
		if rand.Int()%2 == 1 {
			tree.ReplaceOrInsert(itm)
			tree2.ReplaceOrInsert(itm)
		} else {
			tree2.ReplaceOrInsert(itm)
			tree.ReplaceOrInsert(itm)
		}
		test("insert(1)", tree)
		test("insert(2)", tree2)
	}
	test("insert-final(1)", tree)
	test("insert-final(2)", tree2)

	toDel := rand.Perm(n)[nd+2:]
	for _, p := range toDel[nd:] {
		p := p
		dels[(*Intp)(&p)] = struct{}{}
	}
	for d := range dels {
		if rand.Int()%2 == 1 {
			tree.Delete(d)
			tree2.Delete(d)
		} else {
			tree2.Delete(d)
			tree.Delete(d)
		}
		test("delete(1)", tree)
		test("delete(1)", tree2)
	}
	test("delete-final(1)", tree)
	test("delete-final(2)", tree2)

	td1 := toDel[nd]
	td2 := toDel[nd+1]
	if tree.Get((*Intp)(&td1)) != tree2.Get((*Intp)(&td1)) {
		t.Errorf("item %d not shared between trees", toDel[nd])
	}
	if tree.Get((*Intp)(&td2)) != tree2.Get((*Intp)(&td2)) {
		t.Errorf("item %d not shared between trees", toDel[nd+1])
	}
}
