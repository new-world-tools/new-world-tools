package hash

import (
	"sort"
	"sync"
)

type Hash struct {
	FileName string
	Hash     []byte
}

type Registry struct {
	hashes []*Hash
	mu     sync.Mutex
}

func NewRegistry() *Registry {
	return &Registry{
		hashes: []*Hash{},
	}
}

func (registry *Registry) Add(fileName string, hash []byte) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.hashes = append(registry.hashes, &Hash{
		FileName: fileName,
		Hash:     hash,
	})
}

func (registry *Registry) Remove(fileName string) bool {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	for i, h := range registry.hashes {
		if h.FileName == fileName {
			last := len(registry.hashes) - 1
			registry.hashes[i] = registry.hashes[last]
			registry.hashes[last] = nil
			registry.hashes = registry.hashes[:last]

			return true
		}
	}

	return false
}

func (registry *Registry) Hashes() []*Hash {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	sort.Slice(registry.hashes, func(i, j int) bool {
		return registry.hashes[i].FileName < registry.hashes[j].FileName
	})

	return registry.hashes
}
