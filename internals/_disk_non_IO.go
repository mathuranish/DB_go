package internals

// not doing actual i/o here, we'll do that later
// data persistence will also not happen

// we're using a dict/map to mock our actual i/o
type DiskManager struct {
	// nextpage is for managing new id allocation
	nextPage int
	pages map[int]*Node
}

// creating new diskmanager, we'll use this disk manager object throught codebase
func NewDiskManager() *DiskManager {
	return &DiskManager{
		nextPage: 0,
		pages: make(map[int]*Node),
	}
}

// updating and managing id/pagesize counter
func (dm *DiskManager) AllocatePage() int {
	id := dm.nextPage
	dm.nextPage++
	return id
}

// writing to page
func (dm *DiskManager) WriteNode(id int, node *Node) {
	dm.pages[id] = node
}

// reading from pages
func (dm *DiskManager) ReadNode(id int) *Node {
	return dm.pages[id]
}

// deleting from pages
func (dm *DiskManager) Deletenode(id int) {
	delete(dm.pages, id)
}

