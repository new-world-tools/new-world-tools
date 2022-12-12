package azcs

import (
	"encoding/json"
	"hash/crc32"
	"io"
	"os"
	"strings"
	"sync"
)

var DefaultHashRegistry *HashRegistry

func init() {
	DefaultHashRegistry = NewHashRegistry()

	for _, value := range hashBuffBucketsData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashGatherablesData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashLoreData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashNpcData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashPropertiesData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashVitalsCategoryData {
		DefaultHashRegistry.Add(value)
	}
	for _, value := range hashVitalsData {
		DefaultHashRegistry.Add(value)
	}
}

func NewHashRegistry() *HashRegistry {
	return &HashRegistry{
		hashes: map[uint32]string{},
	}
}

type HashRegistry struct {
	mu     sync.Mutex
	hashes map[uint32]string
}

func (registry *HashRegistry) Add(value string) error {
	lowerValue := strings.ToLower(value)
	hash := crc32.ChecksumIEEE([]byte(lowerValue))

	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.hashes[hash] = lowerValue

	return nil
}

func (registry *HashRegistry) Has(hash uint32) bool {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	_, ok := registry.hashes[hash]

	return ok
}

func (registry *HashRegistry) Get(hash uint32) string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	return registry.hashes[hash]
}

func (registry *HashRegistry) Remove(hash uint32) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	delete(registry.hashes, hash)
}

func (registry *HashRegistry) Hashes() map[uint32]string {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	return registry.hashes
}

func LoadHashes(filePath string, registry *HashRegistry) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	var names map[uint32]string

	dec := json.NewDecoder(f)
	err = dec.Decode(&names)
	if err == io.EOF {
		names = map[uint32]string{}
	}
	if err != nil && err != io.EOF {
		return err
	}

	for _, value := range names {
		registry.Add(value)
	}

	return nil
}

func StoreHashes(filePath string, registry *HashRegistry) error {
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
	err = enc.Encode(registry.Hashes())
	if err != nil {
		return err
	}

	return nil
}
