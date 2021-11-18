package hash

import (
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

func (registry *Registry) Hashes() []*Hash {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	return registry.hashes
}
