package azcs

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

var DefaultTypeRegistry *TypeRegistry

func init() {
	DefaultTypeRegistry = NewTypeRegistry()

	for id, value := range defaultTypes {
		DefaultTypeRegistry.Add(id, value)
	}
}

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		types: map[string]string{},
	}
}

type TypeRegistry struct {
	mu    sync.Mutex
	types map[string]string
}

func (registry *TypeRegistry) Add(id string, value string) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.types[id] = value

	return nil
}

func (registry *TypeRegistry) Has(id string) bool {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	_, ok := registry.types[id]

	return ok
}

func (registry *TypeRegistry) Get(id string) string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	return registry.types[id]
}

func (registry *TypeRegistry) Remove(id string) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	delete(registry.types, id)
}

func (registry *TypeRegistry) Types() map[string]string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	return registry.types
}

func LoadTypes(filePath string, registry *TypeRegistry) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var types map[string]string

	dec := json.NewDecoder(f)
	err = dec.Decode(&types)
	if err == io.EOF {
		types = map[string]string{}
	}
	if err != nil && err != io.EOF {
		return err
	}

	for id, value := range types {
		registry.Add(id, value)
	}

	return nil
}

func StoreTypes(filePath string, registry *TypeRegistry) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	err = enc.Encode(registry.Types())
	if err != nil {
		return err
	}

	return nil
}
