# GoLLRB

GoLLRB is a Left-Leaning Red-Black (LLRB) implementation of 2-3 balanced binary
search trees in Go Language.

## Delta from upstream
(Forked from github.com/petar/GoLLRB)

Can now create a tree with NewCoW(), which returns a copy-on-write LLRB tree.
The interface is identical, but you can now Clone() the tree and obtain a
fully read-write duplicate of the tree.

There's a siginifcant performance penalty for writing to cow-enabled trees.
Every update to the tree requires O(Log N) node copy operations (with a constant
factor of about 3). It's also likely to make the garbage collector cry.

That said, it's a quick and dirty way to create forkable trees.

Read performance is unaffected, nor is the performance of New()-created LLRB
trees.

## Overview

As of this writing and to the best of the author's knowledge,
Go still does not have a balanced binary search tree (BBST) data structure.
These data structures are quite useful in a variety of cases. A BBST maintains
elements in sorted order under dynamic updates (inserts and deletes) and can
support various order-specific queries. Furthermore, in practice one often
implements other common data structures like Priority Queues, using BBST's.

2-3 trees (a type of BBST's), as well as the runtime-similar 2-3-4 trees, are
the de facto standard BBST algoritms found in implementations of Python, Java,
and other libraries. The LLRB method of implementing 2-3 trees is a recent
improvement over the traditional implementation. The LLRB approach was
discovered relatively recently (in 2008) by Robert Sedgewick of Princeton
University.

GoLLRB is a Go implementation of LLRB 2-3 trees.

## Maturity

GoLLRB has been used in some pretty heavy-weight machine learning tasks over many gigabytes of data.
I consider it to be in stable, perhaps even production, shape. There are no known bugs.

## Installation

With a healthy Go Language installed, simply run `go get github.com/petar/GoLLRB/llrb`

## Example

	package main

	import (
		"fmt"
		"github.com/petar/GoLLRB/llrb"
	)

	func lessInt(a, b interface{}) bool { return a.(int) < b.(int) }

	func main() {
		tree := llrb.New(lessInt)
		tree.ReplaceOrInsert(1)
		tree.ReplaceOrInsert(2)
		tree.ReplaceOrInsert(3)
		tree.ReplaceOrInsert(4)
		tree.DeleteMin()
		tree.Delete(4)
		c := tree.IterAscend()
		for {
			u := <-c
			if u == nil {
				break
			}
			fmt.Printf("%d\n", int(u.(int)))
		}
	}

## About

GoLLRB was written by [Petar Maymounkov](http://pdos.csail.mit.edu/~petar/).

Follow me on [Twitter @maymounkov](http://www.twitter.com/maymounkov)!
