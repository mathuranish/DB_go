package internals

// actual tree operations

// rootpage is
type BPTree struct {
	dm *DiskManager
	rootPage int
}

func NewTree(dm *DiskManager) *BPTree {
	return &BPTree{
		dm: dm,
		rootPage: -1,
	}
}

// rec function for data insertion
func (t *BPTree) insertIntoNode(node *Node, key int, value []byte) (*Node, int, bool) {
	// adding to leaf node
	if node.isLeaf {
		node.InsertIntoLeaf(key, value)
		if len(node.keys) > BFactor-1 {
			newNode, splitKey := node.SplitLeaf(t.dm)
			return newNode, splitKey, false
		}
		return nil, 0, false
	}

	// adding to internal node
	idx := node.FindAKey(key)
	childId := node.children[idx]
	child := t.dm.ReadNode(childId)
	newChild, splitKey, didSplit := t.insertIntoNode(child, key, value)
	if didSplit {
		node.keys = append(node.keys, 0)
		copy(node.keys[idx+1:], node.keys[idx:])
		node.keys[idx] = splitKey
		node.children = append(node.children, 0)
		copy(node.children[idx+2:], node.children[idx+1:])
		node.children[idx+1] = newChild.id
		if len(node.keys) > BFactor-1 {
			newNode, splitKey := node.SplitInternal(t.dm)
			return newNode, splitKey, true
		}
	}
	return nil, 0, false
}

// insert data to node
func (t *BPTree) InsertNode(key int, value []byte) {
	// when no root is created
	if t.rootPage == -1 {
		root := NewLeafNode(t.dm)
		t.rootPage = root.id
	}
	root := t.dm.ReadNode(t.rootPage)
	newNode, splitKey, didSplit := t.insertIntoNode(root, key, value)
	if didSplit {
		newRoot := NewInternalNode(t.dm)
		newRoot.keys = []int{splitKey}
		newRoot.children = []int{root.id, newNode.id}
		t.rootPage = newRoot.id
	}

}


// seaching functional
func (t *BPTree) SeachInTree(key int) ([]byte, bool) {
	if t.rootPage == -1 {
		return nil, false
	}
	node := t.dm.ReadNode(t.rootPage)
	for !node.isLeaf {
		idx := node.FindAKey(key)
		if idx == len(node.keys) {
			node = t.dm.ReadNode(node.children[idx])
		} else {
			node = t.dm.ReadNode(node.children[idx])
		}
	}
	// feaf node searching
	for i, k := range node.keys {
		if k == key {
			return node.values[i], true
		}
	}
	return nil, false
}
