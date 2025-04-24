package internals

type Database struct {
	tree *BPTree
}

func NewDatabase() *Database {
	dm := NewDiskManager()
	tree := NewTree(dm)
	return &Database{tree: tree}
}

func (db *Database) Put(key int, value []byte) {
	db.tree.InsertNode(key, value)
}

func (db *Database) Get(key int) ([]byte, bool) {
	return db.tree.SeachInTree(key)
}