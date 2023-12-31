package bkdeltree

// BKDelTree is a BK-Tree with support for deletion of words.

import (
	"fmt"
	"strings"

	"github.com/creasty/go-levenshtein"
)

// Distance calculates hamming distance.
func distance(x, y string) int {
	return levenshtein.Distance(x, y)
}

// BKDelTree is a BK-Tree with deletion support.
type BKDelTree struct {
	root *BKNode
	stats
}

// stats is a struct to hold BKDelTree stats.
type stats struct {
	numNodes int
}

// NewBKDelTree creates a new empty BKDelTree.
func NewBKDelTree() *BKDelTree {
	return &BKDelTree{
		root:  nil,
		stats: stats{},
	}
}

// BKNode is a node in a BK-Tree containing the word and a map of children.
type BKNode struct {
	word          string
	childrenCount int
	child         map[int]*BKNode
}

// newBKNode creates a new BKNode with the given word.
func newBKNode(word string) *BKNode {
	return &BKNode{
		word:  word,
		child: make(map[int]*BKNode),
	}
}

// Insert inserts a word into the BKDelTree starting from the root node and goes down the tree.
func (t *BKDelTree) Insert(word string) error {
	if t.root == nil {
		t.root = newBKNode(word)
		t.numNodes++
		return nil
	}
	_, err := t.root.insert(word)
	if err != nil {
		return err
	}
	t.numNodes++
	return nil
}

// insert inserts a word into the BKNode.
func (n *BKNode) insert(word string) (int, error) {
	if n.word == word {
		return 0, fmt.Errorf("word already exists")
	}
	dist := distance(n.word, word)
	if _, ok := n.child[dist]; !ok {
		n.child[dist] = newBKNode(word)
		n.childrenCount++
		return 1, nil
	}
	insertNum, err := n.child[dist].insert(word)
	if err != nil {
		return 0, err
	}
	return insertNum + 1, nil
}

// Search searches the BKDelTree for a word with a maximum distance of maxDist.
// It returns a slice of BKNodes that match the search criteria.
func (t *BKDelTree) Search(word string, maxDist int) []BKNode {
	if t.root == nil {
		return nil
	}
	return t.root.search(word, maxDist, make([]BKNode, 0))
}

// search searches the BKNode for a word with a maximum distance of maxDist.
// It returns a slice of BKNodes that match the search criteria.
func (n *BKNode) search(word string, maxDist int, results []BKNode) []BKNode {
	dist := distance(n.word, word)
	if dist <= maxDist {
		results = append(results, *n)
	}

	for i := dist - maxDist; i <= dist+maxDist; i++ {
		if _, ok := n.child[i]; ok {
			results = n.child[i].search(word, maxDist, results)
		}
	}
	return results
}

// GetParent returns the parent of the BKNode with the given word.
func (t *BKDelTree) GetParent(word string) (*BKNode, error) {
	if t.root == nil {
		return nil, fmt.Errorf("root is nil")
	}
	if t.root.word == word {
		return nil, nil
	}
	return t.root.getParent(word)
}

// getParent returns the parent of the BKNode with the given word.
func (n *BKNode) getParent(word string) (*BKNode, error) {
	dist := distance(n.word, word)
	if childNode, ok := n.child[dist]; ok && childNode.word == word {
		return n, nil
	}
	for i := dist; i <= dist; i++ {
		if _, ok := n.child[i]; ok {
			return n.child[i].getParent(word)
		}
	}
	return nil, fmt.Errorf("word not found")
}

// Delete deletes the BKNode with the given word from the BKDelTree.
func (t *BKDelTree) Delete(word string) error {
	if t.root == nil {
		// Tree is empty
		return nil
	}
	parentNode, err := t.GetParent(word)
	if err != nil {
		// Word not found
		return fmt.Errorf("word not found: %s", word)
	}
	if parentNode == nil {
		// Root Node - Get a random child from root and promote it as root
		// TODO choose child with minimum distance and promote it as root
		newTree := NewBKDelTree()
		allSubNodes := t.root.collectFamily()
		for _, subNode := range allSubNodes[1:] {
			newTree.Insert(subNode.word)
		}
		t.root = newTree.root
		t.numNodes = newTree.numNodes
	} else {
		// Non-root Node
		parentNode.deleteChild(word, t)
	}
	return nil
}

// collectFamily collects all the nodes in sub-tree of the BKNode and returns a slice of BKNodes.
func (n *BKNode) collectFamily() []*BKNode {
	allSubNodes := make([]*BKNode, 0)
	allSubNodes = append(allSubNodes, n)
	for _, child := range n.child {
		allSubNodes = append(allSubNodes, child.collectFamily()...)
	}
	return allSubNodes
}

// deleteChild deletes the child BKNode with the given word from the BKNode.
// It also inserts all the children of the deleted BKNode back into the BKDelTree.
func (n *BKNode) deleteChild(word string, t *BKDelTree) {
	dist := distance(n.word, word)
	if childNode, ok := n.child[dist]; ok && childNode.word == word {
		allSubNodes := childNode.collectFamily()

		delete(n.child, dist)
		t.numNodes--
		n.childrenCount--
		for _, subNode := range allSubNodes[1:] {
			t.numNodes--
			t.Insert(subNode.word)
		}
	}
}

// SPPrintf returns a pretty printed BKDelTree.
// root node is at depth 0
// each child node is indented by 1 period (.), it's children count in brackets and its distance in parentheses
func (t *BKDelTree) SPPrintf(indentChar string) string {
	if indentChar == "" {
		indentChar = ". "
	}
	if t.root == nil {
		// tree is empty
		return ""
	}
	sb := &strings.Builder{}
	t.root.pprint(sb, 0, indentChar, 0)
	return sb.String()
}

// pprint pretty prints the BKNode into a string Builder.
// takes 4 arguments:
//
//	sb: the string builder to write to
//	indentLevel: the indent level at which to start printing
//	indentChar: the character to use for indentation
//	distance: the distance of the BKNode in relation to it's parent
func (n *BKNode) pprint(sb *strings.Builder, indentLevel int, indentChar string, distance int) {
	sb.WriteString(fmt.Sprintf("%s%s:cc[%d] (%d)\n", strings.Repeat(indentChar, indentLevel), n.word, n.childrenCount, distance))
	for dist, child := range n.child {
		child.pprint(sb, indentLevel+1, indentChar, dist)
	}

}
