package datasheet

import (
	"fmt"
	"github.com/new-world-tools/new-world-tools/store"
	"os"
	"sort"
)

func NewStore(dataTableDir string) (*Store, error) {
	files, err := FindAll(dataTableDir)
	if err != nil {
		return nil, err
	}

	store := &Store{
		store: store.NewSimpleStore[*DataSheetFile](),
		types: map[string]map[string]*DataSheetFile{},
	}

	keys := map[string]bool{}

	for _, file := range files {
		f, err := os.Open(file.GetPath())
		if err != nil {
			return nil, err
		}

		meta, err := ParseMeta(f)
		if err != nil {
			return nil, err
		}

		f.Close()

		key := fmt.Sprintf("%s.%s", meta.Type, meta.UniqueId)
		_, ok := keys[key]
		if ok {
			continue
		}

		keys[key] = true

		store.store.Add(fmt.Sprintf("%s.%s", meta.Type, meta.UniqueId), file)
		_, ok = store.types[meta.Type]
		if !ok {
			store.types[meta.Type] = map[string]*DataSheetFile{}
		}
		store.types[meta.Type][meta.UniqueId] = file
	}

	return store, nil
}

type Store struct {
	store *store.Store[*DataSheetFile]
	types map[string]map[string]*DataSheetFile
}

func (store *Store) GetDataSheet(key string) (*DataSheet, error) {
	if !store.store.Has(key) {
		return nil, fmt.Errorf("%q is not exists", key)
	}

	ds, err := Parse(store.store.Get(key))
	if err != nil {
		return nil, fmt.Errorf("datasheet.Parse: %s", err)
	}

	return ds, nil
}

func (store *Store) GetDataSheets(key string) (*DataSheet, error) {
	if !store.store.Has(key) {
		return nil, fmt.Errorf("%q is not exists", key)
	}

	ds, err := Parse(store.store.Get(key))
	if err != nil {
		return nil, fmt.Errorf("datasheet.Parse: %s", err)
	}

	return ds, nil
}

func (store *Store) GetDataSheetMeta(key string) (*Meta, error) {
	if !store.store.Has(key) {
		return nil, fmt.Errorf("%q is not exists", key)
	}

	f, err := os.Open(store.store.Get(key).GetPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dsMeta, err := ParseMeta(f)
	if err != nil {
		return nil, fmt.Errorf("datasheet.ParseMeta: %s", err)
	}

	return dsMeta, nil
}

func (store *Store) GetTypes() []string {
	types := make([]string, len(store.types))

	var i int
	for typ, _ := range store.types {
		types[i] = typ
		i++
	}

	sort.Strings(types)

	return types
}

func (store *Store) GetKeys() []string {
	keys := []string{}
	for typ, files := range store.types {
		for uniqueId, _ := range files {
			key := fmt.Sprintf("%s.%s", typ, uniqueId)
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys
}

func (store *Store) GetDataSheetFiles(typ string) []*DataSheetFile {
	files, ok := store.types[typ]
	if !ok {
		return []*DataSheetFile{}
	}

	dataSheetFiles := make([]*DataSheetFile, len(files))

	uniqueIds := make([]string, len(files))
	var i int
	for uniqueId, _ := range files {
		uniqueIds[i] = uniqueId
		i++
	}
	sort.Strings(uniqueIds)

	for i, uniqueId := range uniqueIds {
		dataSheetFiles[i] = store.types[typ][uniqueId]
	}

	return dataSheetFiles
}
