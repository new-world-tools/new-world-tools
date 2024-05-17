package structure

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/goccy/go-yaml"
	"sync"
)

type OrderedMap[K comparable, V any] struct {
	mu       sync.RWMutex
	keys     []K
	values   map[K]V
	position int
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		keys:     []K{},
		values:   map[K]V{},
		position: 0,
	}
}

func (orderedMap *OrderedMap[K, V]) Reset() {
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

	orderedMap.position = 0
}

func (orderedMap *OrderedMap[K, V]) Add(key K, value V) {
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

	orderedMap.keys = append(orderedMap.keys, key)
	orderedMap.values[key] = value
}

func (orderedMap *OrderedMap[K, V]) Get(key K) (V, bool) {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	value, ok := orderedMap.values[key]

	return value, ok
}

func (orderedMap *OrderedMap[K, V]) GetByPosition(position int) (key K, value V, ok bool) {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	if len(orderedMap.keys) > position {
		key := orderedMap.keys[position]
		value := orderedMap.values[key]

		return key, value, true
	}

	return
}

func (orderedMap *OrderedMap[K, V]) Has() bool {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	return len(orderedMap.keys) > orderedMap.position
}

func (orderedMap *OrderedMap[K, V]) Next() (K, V) {
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

	key := orderedMap.keys[orderedMap.position]
	value := orderedMap.values[key]
	orderedMap.position++

	return key, value
}

func (orderedMap *OrderedMap[K, V]) Size() int {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	return len(orderedMap.keys)
}

func (orderedMap *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	var data []byte
	var err error

	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, key := range orderedMap.keys {
		if i > 0 {
			buf.WriteRune(',')
		}

		data, err = json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
		buf.WriteRune(':')

		data, err = json.Marshal(orderedMap.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (orderedMap *OrderedMap[K, V]) MarshalYAML() (any, error) {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	mapSlice := yaml.MapSlice{}

	for _, key := range orderedMap.keys {
		value := orderedMap.values[key]
		mapSlice = append(mapSlice, yaml.MapItem{Key: key, Value: value})
	}

	return mapSlice, nil
}

func (orderedMap *OrderedMap[K, V]) ToMap() map[K]V {
	orderedMap.mu.RLock()
	defer orderedMap.mu.RUnlock()

	m := make(map[K]V)
	for _, key := range orderedMap.keys {
		m[key] = orderedMap.values[key]
	}
	return m
}
