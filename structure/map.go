package structure

import (
	"bytes"
	"encoding/json"
	"sync"
)

type OrderedMap[K comparable, V any] struct {
	mu       sync.Mutex
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
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

	value, ok := orderedMap.values[key]

	return value, ok
}

func (orderedMap *OrderedMap[K, V]) GetByPosition(position int) (key K, value V, ok bool) {
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

	if len(orderedMap.keys) > position {
		key := orderedMap.keys[position]
		value := orderedMap.values[key]

		return key, value, true
	}

	return
}

func (orderedMap *OrderedMap[K, V]) Has() bool {
	orderedMap.mu.Lock()
	defer orderedMap.mu.Unlock()

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

func (orderedMap *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	var data []byte
	var err error
	var buf bytes.Buffer

	buf.WriteRune('{')
	orderedMap.Reset()
	var i int
	for orderedMap.Has() {
		key, value := orderedMap.Next()

		if i > 0 {
			buf.WriteRune(',')
		}

		data, err = json.Marshal(key)
		if err != nil {
			return nil, err
		}

		buf.Write(data)
		buf.WriteRune(':')

		data, err = json.Marshal(value)
		if err != nil {
			return nil, err
		}

		buf.Write(data)

		i++
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (orderedMap *OrderedMap[K, V]) ToMap() map[K]V {
	m := make(map[K]V)
	orderedMap.Reset()
	for orderedMap.Has() {
		key, value := orderedMap.Next()
		m[key] = value
	}
	return m
}

func (orderedMap *OrderedMap[K, V]) MarshalYAML() (any, error) {
	return orderedMap.ToMap(), nil
}
