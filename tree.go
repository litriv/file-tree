package tree

import (
	"code.litriv.com/comparison"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type WalkFn func(*Node) error

// node
type nodeType int

const (
	containerType nodeType = iota
	leafType
)

type Node struct {
	tree     *Tree
	parent   *Node
	children []*Node
	nodeType nodeType
	Level    int
	Context  interface{}
}

func (n *Node) Index() int {
	if n.parent != nil {
		for i, child := range n.parent.children {
			if child == n {
				return i
			}
		}
	}
	return 0
}

func (n *Node) IsRoot() bool {
	return n.parent == nil
}

func (n *Node) IsLeaf() bool {
	return n.nodeType == leafType
}

func (c *Node) hasSpace() bool {
	return len(c.children) < c.tree.width
}

func (c *Node) findChildWithSpace() *Node {
	for _, child := range c.children {
		if child.hasSpace() {
			return child
		}
	}
	r := &Node{c.tree, c, make([]*Node, 0), containerType, c.Level + 1, nil}
	c.children = append(c.children, r)
	return r
}

func (c *Node) add(l *Node) error {
	if c.hasSpace() {
		if c.Level == c.tree.depth { // for last level node
			c.children = append(c.children, l)
			l.parent = c
			if c.tree.Listener != nil {
				err := c.tree.Listener.LeafAdded(l)
				if err != nil {
					return err
				}
			}
		} else { // for nodes higher up
			child := c.findChildWithSpace()
			c.tree.lastTouchedContainer = child
			return child.add(l)
		}
	} else {
		if c.parent != nil {
			return c.parent.add(l)
		} else {
			r := &Node{
				tree:     c.tree,
				children: make([]*Node, 1),
				nodeType: containerType}
			c.tree.root = r
			r.children[0] = c
			c.parent = r
			l.Level++
			r.incChildLevels()
			c.tree.depth += 1
			if c.tree.Listener != nil {
				err := c.tree.Listener.NewRootInserted(r)
				if err != nil {
					return err
				}
			}
			return r.add(l)
		}
	}
	return nil
}

func (n *Node) incChildLevels() {
	for _, c := range n.children {
		c.Level += 1
		c.incChildLevels()
	}
}

// Listener
type Listener interface {
	LeafAdded(l *Node) error
	NewRootInserted(r *Node) error
}

// Tree

// Contexts must implement this interface when wanting to reject duplicates (Tree.rejectDuplicates is true)
type DuplicateTracker interface {
	comparison.Eq
	Key() interface{}
}

type Tree struct {
	width, depth         int
	lastTouchedContainer *Node
	root                 *Node
	pathTracker          *pathTracker
	rejectDuplicates     bool
	contexts             map[interface{}]bool
	Added                int
	Rejected             int
	Listener             Listener
}

func NewTree(width int, rejectDuplicates bool) *Tree {
	t := &Tree{width: width, rejectDuplicates: rejectDuplicates}
	if rejectDuplicates {
		t.contexts = make(map[interface{}]bool, 0)
	}
	c := &Node{t, nil, make([]*Node, 0), containerType, 0, nil}
	t.lastTouchedContainer = c
	t.root = c
	return t
}

func (t *Tree) AddLeaf(Context interface{}) error {
	if t.rejectDuplicates {
		dupTrackerCtx := Context.(DuplicateTracker)
		if _, ok := t.contexts[dupTrackerCtx.Key()]; ok {
			t.Rejected++
			return nil
		} else {
			t.contexts[dupTrackerCtx.Key()] = true
		}
	}
	l := &Node{t, nil, make([]*Node, 0), leafType, t.depth + 1, Context}
	l.tree = t
	err := t.lastTouchedContainer.add(l)
	if err != nil {
		return err
	}
	t.Added++
	return nil
}

func (t *Tree) Stats() string {
	return fmt.Sprintf("Added: %d\nRejected: %d\n", t.Added, t.Rejected)
}

func stringWalker(n *Node) error {
	fmt.Printf("%s", n.Path())
	if n.Context != nil {
		fmt.Printf("/%d", n.Context)
	}
	fmt.Print("\n")
	return nil
}

func (t *Tree) String() string {
	r := ""
	t.Walk(func(n *Node) error {
		r = fmt.Sprintf("%s%s", r, n.Path())
		if n.Context != nil {
			r = fmt.Sprintf("%s/%d", r, n.Context)
		}
		r = fmt.Sprintf("%s\n", r)
		return nil
	})
	r = fmt.Sprintf("%s%s", r, t.Stats())
	return r
}

// Eq compares to another tree, param o, and returns true if they contain the same elements in the same structure, with the same contexts. Returns a bool for equality and a string that contains details about the inequality if appropriate, or nil otherwise. The context has to implement comparison.Eq
func (t *Tree) Eq(o *Tree) (bool, error) {
	tItems := [][2]interface{}{}
	t.Walk(func(n *Node) error {
		tItems = append(tItems, [2]interface{}{n.Path(), n.Context})
		return nil
	})

	oItems := [][2]interface{}{}
	o.Walk(func(n *Node) error {
		oItems = append(oItems, [2]interface{}{n.Path(), n.Context})
		return nil
	})

	if len(tItems) != len(oItems) {
		return false, errors.New("trees have different sizes or structures")
	}

	for i, item := range tItems {
		switch {
		case item[0] != oItems[i][0]:
			return false, errors.New("paths don't match")
		case item[1] == nil && oItems[i][1] == nil:
			continue
		case item[1] == nil && oItems[i][1] != nil || item[1] != nil && oItems[i][1] == nil || !item[1].(comparison.Eq).Eq(oItems[i][1].(comparison.Eq)):
			return false, errors.New("contexts don't match")
		}
	}
	return true, nil
}

// Walker
type pathTracker struct {
	path  string
	level int
}

func (p *pathTracker) update(n *Node) {
	index := strconv.FormatInt(int64(n.Index()), 10)
	if n.IsRoot() {
		p.path = "0"
	} else if n.Level > p.level {
		p.path = filepath.Join(p.path, index)
		p.level = n.Level
	} else if n.Level < p.level {
		p.path = cutPathAndAppend(p.path, n.Level, index)
		p.level = n.Level
	} else {
		splitPath := strings.Split(p.path, string(filepath.Separator))
		p.path = cutPathAndAppend(p.path, len(splitPath)-1, index)
	}
}

func cutPathAndAppend(path string, whereToCut int, toAppend string) string {
	splitPath := append(strings.Split(path, string(filepath.Separator))[0:whereToCut], toAppend)
	return filepath.Join(splitPath...)
}

func (n *Node) Path() string {
	return n.tree.pathTracker.path
}

func (n *Node) walk(fn WalkFn) error {
	n.tree.pathTracker.update(n)
	err := fn(n)
	if err != nil {
		return err
	}
	for _, child := range n.children {
		err := child.walk(fn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tree) Walk(fn WalkFn) error {
	if t.root != nil {
		t.pathTracker = &pathTracker{"", -1}
		err := t.root.walk(fn)
		if err != nil {
			return err
		}
	}
	return nil
}
