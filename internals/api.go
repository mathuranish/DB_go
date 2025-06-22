package internals

type Database struct {
	tree *BPTree
	dm *DiskManager
}

func NewDatabase() (*Database, error) {
	dm, err := NewDiskManager("database_test.db")
	if err != nil {
		return nil, err
	}
	tree := NewTree(dm)
	return &Database{
		tree: tree,
		dm:   dm,
	}, nil
}

func (db *Database) Put(key int, value []byte) {
	db.tree.InsertNode(key, value)
}

func (db *Database) Get(key int) ([]byte, bool) {
	return db.tree.SeachInTree(key)
}

func (db *Database) Close() error {
	return db.dm.Close()
}
