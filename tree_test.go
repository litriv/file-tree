package tree

import (
	"fmt"
)

// Accept
func ExampleAcceptDuplicates() {
	tree := NewTree(2, false)
	for i := 1; i < 7; i++ {
		tree.AddLeaf(i)
	}
	for i := 1; i < 7; i++ {
		tree.AddLeaf(i)
	}
	fmt.Println(tree)
	// Output:
	// 0
	// 0/0
	// 0/0/0
	// 0/0/0/0
	// 0/0/0/0/0/1
	// 0/0/0/0/1/2
	// 0/0/0/1
	// 0/0/0/1/0/3
	// 0/0/0/1/1/4
	// 0/0/1
	// 0/0/1/0
	// 0/0/1/0/0/5
	// 0/0/1/0/1/6
	// 0/0/1/1
	// 0/0/1/1/0/1
	// 0/0/1/1/1/2
	// 0/1
	// 0/1/0
	// 0/1/0/0
	// 0/1/0/0/0/3
	// 0/1/0/0/1/4
	// 0/1/0/1
	// 0/1/0/1/0/5
	// 0/1/0/1/1/6
	// Added: 12
	// Rejected: 0
}

// Reject
type Test_ctx int

func (c Test_ctx) Eq(other interface{}) bool {
	return c == other
}
func (c Test_ctx) Key() interface{} {
	return c
}
func ExampleRejectDuplicates() {
	tree := NewTree(2, true)
	var i Test_ctx
	for i = 1; i < 7; i++ {
		tree.AddLeaf(i)
	}
	for i = 1; i < 7; i++ {
		tree.AddLeaf(i)
	}
	fmt.Println(tree)
	// Output:
	// 0
	// 0/0
	// 0/0/0
	// 0/0/0/0/1
	// 0/0/0/1/2
	// 0/0/1
	// 0/0/1/0/3
	// 0/0/1/1/4
	// 0/1
	// 0/1/0
	// 0/1/0/0/5
	// 0/1/0/1/6
	// Added: 6
	// Rejected: 6
}

func ExampleEquals() {
	tree := NewTree(2, true)
	var i Test_ctx
	for i = 1; i < 100; i++ {
		tree.AddLeaf(i)
	}
	tree2 := NewTree(2, true)
	for i = 1; i < 100; i++ {
		tree2.AddLeaf(i)
	}
	fmt.Println(tree.Eq(tree2))
	// Output:
	// true <nil>
}

func ExampleNotEqualsStructuresDiffer() {
	tree := NewTree(4, true)
	var i Test_ctx
	for i = 1; i < 100; i++ {
		tree.AddLeaf(i)
	}
	tree2 := NewTree(2, true)
	for i = 1; i < 100; i++ {
		tree2.AddLeaf(i)
	}
	fmt.Println(tree.Eq(tree2))
	// Output:
	// false trees have different sizes or structures
}

func ExampleNotEqualsSizessDiffer() {
	tree := NewTree(2, true)
	var i Test_ctx
	for i = 1; i < 100; i++ {
		tree.AddLeaf(i)
	}
	tree2 := NewTree(2, true)
	for i = 1; i < 101; i++ {
		tree2.AddLeaf(i)
	}
	fmt.Println(tree.Eq(tree2))
	// Output:
	// false trees have different sizes or structures
}

func ExampleNotEqualContextsDiffer() {
	tree := NewTree(2, true)
	var i Test_ctx
	for i = 1; i < 100; i++ {
		tree.AddLeaf(i)
	}
	tree2 := NewTree(2, true)
	for i = 1; i < 100; i++ {
		tree2.AddLeaf(i + 1)
	}
	fmt.Println(tree.Eq(tree2))
	// Output:
	// false contexts don't match
}
func ExampleSize() {
	tree := NewTree(4, true)
	var i Test_ctx
	for i = 0; i < 100; i++ {
		tree.AddLeaf(i)
	}
	fmt.Println(tree.Added)
	// Output:
	// 100
}