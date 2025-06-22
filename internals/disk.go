package internals

import (
	"encoding/binary"
	"os"
	"sync"
)

const (
	PageSize    = 4096 // fixed page size in bytes
	MaxKeys     = BFactor - 1
	MaxChildren = BFactor
	MaxValueLen = 256 // max length of a single value
)

type DiskManager struct {
	// nextpage is for managing new id allocation
	file *os.File
	nextPage int
	// pages map[int]*Node 
	mutex sync.Mutex
}

// creating new diskmanager, we'll use this disk manager object throught codebase
func NewDiskManager(filename string) (*DiskManager, error ){
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	
	// checking file size
	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}
	nextPage := int(stat.Size() / PageSize)
	return &DiskManager{
		file:     file,
		nextPage: nextPage,
	}, nil
}

func (dm *DiskManager) Close() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	return dm.file.Close()
}

// updating and managing id/pagesize counter
func (dm *DiskManager) AllocatePage() int {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	page := dm.nextPage
	dm.nextPage++
	// Extend file size if necessary
	if err := dm.file.Truncate(int64(dm.nextPage * PageSize)); err != nil {
		panic(err)
	}
	return page
}

// writing to page
func (dm *DiskManager) WriteNode(page int, node *Node) {

	// lock for thread safety
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	buf := make([]byte, PageSize)
	offset := 0

	if node.isLeaf {
		buf[offset] = 1
	} else {
		buf[offset] = 0
	}
	offset++
	binary.BigEndian.PutUint32(buf[offset:offset+4], uint32(node.id))
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:offset+4], uint32(node.next))
	offset += 4
	binary.BigEndian.PutUint32(buf[offset:offset+4], uint32(node.prev))
	offset += 4
	binary.BigEndian.PutUint16(buf[offset:offset+2], uint16(len(node.keys)))
	offset += 2

	// writing keys
	for _, key := range node.keys {
		binary.BigEndian.PutUint32(buf[offset:offset+4], uint32(key))
		offset += 4
	}
	// padding unused keys for consistency
	for i := len(node.keys); i < MaxKeys; i++ {
		binary.BigEndian.PutUint32(buf[offset:offset+4], 0)
		offset += 4
	}

	if node.isLeaf {
		// writing values
		for _, val := range node.values {
			length := len(val)
			if length > MaxValueLen {
				length = MaxValueLen
			}
			binary.BigEndian.PutUint16(buf[offset:offset+2], uint16(length))
			offset += 2
			copy(buf[offset:offset+length], val[:length])
			offset += length
			// pad remaining space for this value
			for j := length; j < MaxValueLen; j++ {
				buf[offset] = 0
				offset++
			}
		}
		// padding extra
		for i := len(node.values); i < MaxKeys; i++ {
			binary.BigEndian.PutUint16(buf[offset:offset+2], 0)
			offset += 2
			for j := 0; j < MaxValueLen; j++ {
				buf[offset] = 0
				offset++
			}
		}
	} else {
		// writing to children
		for _, child := range node.children {
			binary.BigEndian.PutUint32(buf[offset:offset+4], uint32(child))
			offset += 4
		}
		for i := len(node.children); i < MaxChildren; i++ {
			binary.BigEndian.PutUint32(buf[offset:offset+4], 0)
			offset += 4
		}
	}

	// file write final
	_, err := dm.file.WriteAt(buf, int64(page*PageSize))
	if err != nil {
		panic(err)
	}
}

// reading from pages
func (dm *DiskManager) ReadNode(page int) *Node {
	// thread safety
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	buf := make([]byte, PageSize)
	_, err := dm.file.ReadAt(buf, int64(page*PageSize))
	if err != nil {
		panic(err)
	}

	node := &Node{}
	offset := 0

	node.isLeaf = buf[offset] == 1
	offset++
	node.id = int(binary.BigEndian.Uint32(buf[offset : offset+4]))
	offset += 4
	node.next = int(binary.BigEndian.Uint32(buf[offset : offset+4]))
	offset += 4
	node.prev = int(binary.BigEndian.Uint32(buf[offset : offset+4]))
	offset += 4
	keyCount := int(binary.BigEndian.Uint16(buf[offset : offset+2]))
	offset += 2

	node.keys = make([]int, 0, MaxKeys)
	for i := 0; i < keyCount; i++ {
		key := int(binary.BigEndian.Uint32(buf[offset : offset+4]))
		node.keys = append(node.keys, key)
		offset += 4
	}
	// ignoring padding
	offset += 4 * (MaxKeys - keyCount) 

	if node.isLeaf {
		node.values = make([][]byte, 0, MaxKeys)
		for i := 0; i < keyCount; i++ {
			length := int(binary.BigEndian.Uint16(buf[offset : offset+2]))
			offset += 2
			val := make([]byte, length)
			copy(val, buf[offset:offset+length])
			node.values = append(node.values, val)
			offset += length
			// ignoring padding
			offset += MaxValueLen - length
		}
		// ignoring padding
		offset += (2 + MaxValueLen) * (MaxKeys - keyCount)
	} else {
		node.children = make([]int, 0, MaxChildren)
		for i := 0; i < keyCount+1; i++ {
			if i < len(node.keys) || (i == len(node.keys) && len(node.children) < MaxChildren) {
				child := int(binary.BigEndian.Uint32(buf[offset : offset+4]))
				node.children = append(node.children, child)
			}
			offset += 4
		}
	}

	return node
}

// deleting from pages
func (dm *DiskManager) DeletePage(page int) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()
	// marking page as reusable by resetting it
	buf := make([]byte, PageSize)
	_, err := dm.file.WriteAt(buf, int64(page*PageSize))
	if err != nil {
		panic(err)
	}
}
