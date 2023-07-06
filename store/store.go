package store

import (
	"sort"
	"strings"
)

type Store[T any] struct {
	values        map[string]T
	keyNormalizer func(key string) string
}

func (store *Store[T]) Add(key string, item T) {
	store.values[store.keyNormalizer(key)] = item
}

func (store *Store[T]) Has(key string) bool {
	_, ok := store.values[store.keyNormalizer(key)]
	return ok
}

func (store *Store[T]) Get(key string) T {
	val, _ := store.values[store.keyNormalizer(key)]
	return val
}

func (store *Store[T]) GetKeys() []string {
	keys := make([]string, len(store.values))

	var i int
	for key, _ := range store.values {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	return keys
}

func NewSimpleStore[T any]() *Store[T] {
	return &Store[T]{
		values: map[string]T{},
		keyNormalizer: func(key string) string {
			return key
		},
	}
}

func NewCaseInsensitiveStore[T any]() *Store[T] {
	return &Store[T]{
		values: map[string]T{},
		keyNormalizer: func(key string) string {
			return strings.ToLower(key)
		},
	}
}
func NewStore[T any](keyNormalizer func(key string) string) *Store[T] {
	return &Store[T]{
		values:        map[string]T{},
		keyNormalizer: keyNormalizer,
	}
}
