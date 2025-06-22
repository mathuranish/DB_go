package internals

const BFactor = 4

type Node struct {
	// for normal non-leaf /internal node
	id int
	isLeaf bool
	keys []int
	children []int

	// for leaf nodes
	values [][]byte
	next int
	prev int
}

// to create internal node
func NewInternalNode(dm *DiskManager) *Node {
	id := dm.AllocatePage()
	node := &Node{
		id : id,
		isLeaf : false,
		keys : make([]int, 0, BFactor-1),
		children : make([]int, 0, BFactor),
		next : -1,
		prev : -1,
	}
	dm.WriteNode(id,node)
	return node
}

// to create leaf node
func NewLeafNode(dm *DiskManager) *Node {
	id := dm.AllocatePage()
	node := &Node{
		id: id,
		isLeaf : true,
		keys : make([]int, 0, BFactor-1),
		children : nil,
		values : make([][]byte, 0, BFactor-1),
		next : -1,
		prev : -1,
	}
	dm.WriteNode(id,node)
	return node
}

// to check if node is full, this needs to be methord
func (n *Node) IsNodeFull() bool {
	// leaf nodes are sctrictly checked with keys while internal nodes "fullness" depends on children
	if n.isLeaf{
		return len(n.keys)>=BFactor-1
	}
	return len(n.children)>=BFactor
}

// to find index of a key in a node
func (n *Node) FindAKey(keyToFind int) int {
	for ind, key := range n.keys {
		if keyToFind <= key {
			return ind
		}
	}
	return len(n.keys)
}

// insertion
func (n *Node) InsertIntoLeaf(key int, value []byte) {
	if !n.isLeaf {
		panic("Invalid method called, insertion in non leaf node called")
	}
	idx := n.FindAKey(key)
	// checking for existing key and updating it
	if idx < len(n.keys) && n.keys[idx] == key {
		n.values[idx] = value
		return
	}

	// need to create new key and shift array
	n.keys = append(n.keys, 0)
	copy(n.keys[idx+1:], n.keys[idx:])
	n.keys[idx] = key
	// same for values
	n.values = append(n.values, nil)
	copy(n.values[idx+1:], n.values[idx:])
	n.values[idx]=value
}


// splitting leaf and returns the new node and split key
func (n *Node) SplitLeaf(dm *DiskManager) (*Node, int) {
	if !n.isLeaf {
		panic("Invalid method called, spliting in leaf node called")
	}
	// finding middle element and spliting node
	mid := len(n.keys)/2
	newNode := NewLeafNode(dm)
	// ellipsis operator (...) to open the values if slice/list
	newNode.keys = append(newNode.keys, n.keys[mid:]...)
	newNode.values = append(newNode.values, n.values[mid:]...)
	n.keys = n.keys[:mid]
	n.values = n.values[:mid]

	// updateing ptrs
	newNode.next = n.next
	newNode.prev = n.id
	if n.next != -1 {
		nextNode := dm.ReadNode(n.next)
		nextNode.prev = newNode.id
		dm.WriteNode(nextNode.id, nextNode)
	}
	n.next = newNode.id
	dm.WriteNode(n.id, n)
	dm.WriteNode(newNode.id, newNode)
	return newNode, newNode.keys[0]
}

// spliting internal node and returns the new node and split key
// this is similiar to splitleaf, just with children updated instead of values
func (n *Node) SplitInternal(dm *DiskManager) (*Node, int) {
	if n.isLeaf {
		panic("Invalid method called, spliting in non leaf node called")
	}
	mid := len(n.keys) / 2
	splitKey := n.keys[mid]
	newNode := NewInternalNode(dm)
	newNode.keys = append(newNode.keys, n.keys[mid+1:]...)
	newNode.children = append(newNode.children, n.children[mid+1:]...)
	n.keys = n.keys[:mid]
	n.children = n.children[:mid+1]
	dm.WriteNode(n.id, n)
	dm.WriteNode(newNode.id, newNode)
	return newNode, splitKey
}

