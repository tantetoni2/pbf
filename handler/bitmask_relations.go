package handler

import (
	"sync"

	"github.com/missinglink/pbf/lib"

	"github.com/missinglink/gosmparse"
)

// BitmaskBoundaries - Load all elements in to memory
type BitmaskRelations struct {
	Pass            int
	Mutex           *sync.Mutex
	Masks           *lib.BitmaskMap
	RelationMembers map[int64][]gosmparse.RelationMember
}

// ReadNode - called once per node
func (b *BitmaskRelations) ReadNode(item gosmparse.Node) { /* noop */ }

// ReadWay - called once per way
func (b *BitmaskRelations) ReadWay(item gosmparse.Way) {

	// only run on second pass
	if b.Pass != 1 {
		return
	}

	// must be in bitmask
	if !b.Masks.Ways.Has(item.ID) {
		return
	}

	// insert dependents in their masks
	for _, ref := range item.NodeIDs {
		b.Masks.Nodes.Insert(ref)
	}
}

// ReadRelation - called once per relation
func (b *BitmaskRelations) ReadRelation(item gosmparse.Relation) {

	// only run on first pass
	if b.Pass != 0 {
		return
	}

	// store ALL relation members in memory
	b.Mutex.Lock()
	b.RelationMembers[item.ID] = item.Members
	b.Mutex.Unlock()

	// insert item in the relations mask
	b.Masks.Relations.Insert(item.ID)
}
