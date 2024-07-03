package store

import (
	"cmp"
	"sort"
	"strings"
	"sync"
)

type ConcurrentStore[K cmp.Ordered, V any] struct {
	values        map[K]V
	keyNormalizer func(key K) K
	mu            sync.Mutex
}

func (store *ConcurrentStore[K, V]) Add(key K, item V) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.values[store.keyNormalizer(key)] = item
}

func (store *ConcurrentStore[K, V]) Has(key K) bool {
	store.mu.Lock()
	defer store.mu.Unlock()

	_, ok := store.values[store.keyNormalizer(key)]
	return ok
}

func (store *ConcurrentStore[K, V]) Get(key K) V {
	store.mu.Lock()
	defer store.mu.Unlock()

	val, _ := store.values[store.keyNormalizer(key)]
	return val
}

func (store *ConcurrentStore[K, V]) GetKeys() []K {
	store.mu.Lock()
	defer store.mu.Unlock()

	keys := make([]K, len(store.values))

	var i int
	for key, _ := range store.values {
		keys[i] = key
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func NewSimpleConcurrentStore[K cmp.Ordered, V any]() *ConcurrentStore[K, V] {
	return &ConcurrentStore[K, V]{
		values: map[K]V{},
		keyNormalizer: func(key K) K {
			return key
		},
	}
}

func NewCaseInsensitiveConcurrentStore[V any]() *ConcurrentStore[string, V] {
	return &ConcurrentStore[string, V]{
		values: map[string]V{},
		keyNormalizer: func(key string) string {
			return strings.ToLower(key)
		},
	}
}
func NewConcurrentStore[K cmp.Ordered, V any](keyNormalizer func(key K) K) *ConcurrentStore[K, V] {
	return &ConcurrentStore[K, V]{
		values:        map[K]V{},
		keyNormalizer: keyNormalizer,
	}
}
