package store

import (
	"cmp"
	"sort"
	"strings"
)

type Store[K cmp.Ordered, V any] struct {
	values        map[K]V
	keyNormalizer func(key K) K
}

func (store *Store[K, V]) Add(key K, item V) {
	store.values[store.keyNormalizer(key)] = item
}

func (store *Store[K, V]) Has(key K) bool {
	_, ok := store.values[store.keyNormalizer(key)]
	return ok
}

func (store *Store[K, V]) Get(key K) V {
	val, _ := store.values[store.keyNormalizer(key)]
	return val
}

func (store *Store[K, V]) GetKeys() []K {
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

func NewSimpleStore[K cmp.Ordered, V any]() *Store[K, V] {
	return &Store[K, V]{
		values: map[K]V{},
		keyNormalizer: func(key K) K {
			return key
		},
	}
}

func NewCaseInsensitiveStore[V any]() *Store[string, V] {
	return &Store[string, V]{
		values: map[string]V{},
		keyNormalizer: func(key string) string {
			return strings.ToLower(key)
		},
	}
}
func NewStore[K cmp.Ordered, V any](keyNormalizer func(key K) K) *Store[K, V] {
	return &Store[K, V]{
		values:        map[K]V{},
		keyNormalizer: keyNormalizer,
	}
}
