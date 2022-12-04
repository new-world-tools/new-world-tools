package datasheet

import (
	"fmt"
	"github.com/new-world-tools/new-world-tools/internal"
	"sort"
)

func NewStore(dataTableDir string) (*Store, error) {
	files, err := FindAll(dataTableDir)
	if err != nil {
		return nil, err
	}

	store := &Store{
		store:     internal.NewSimpleStore[*DataSheetFile](),
		dataTypes: map[string]map[string]*DataSheetFile{},
	}

	keys := map[string]bool{}

	for _, file := range files {
		meta, err := ParseMeta(file)
		if err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s.%s", meta.DataType, meta.UniqueId)
		_, ok := keys[key]
		if ok {
			continue
		}

		keys[key] = true

		store.store.Add(fmt.Sprintf("%s.%s", meta.DataType, meta.UniqueId), file)
		_, ok = store.dataTypes[meta.DataType]
		if !ok {
			store.dataTypes[meta.DataType] = map[string]*DataSheetFile{}
		}
		store.dataTypes[meta.DataType][meta.UniqueId] = file
	}

	return store, nil
}

type Store struct {
	store     *internal.Store[*DataSheetFile]
	dataTypes map[string]map[string]*DataSheetFile
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

	dsMeta, err := ParseMeta(store.store.Get(key))
	if err != nil {
		return nil, fmt.Errorf("datasheet.ParseMeta: %s", err)
	}

	return dsMeta, nil
}

func (store *Store) GetDataTypes() []string {
	dataTypes := make([]string, len(store.dataTypes))

	var i int
	for dataType, _ := range store.dataTypes {
		dataTypes[i] = dataType
		i++
	}

	sort.Strings(dataTypes)

	return dataTypes
}

func (store *Store) GetKeys() []string {
	keys := []string{}
	for dataType, files := range store.dataTypes {
		for uniqueId, _ := range files {
			key := fmt.Sprintf("%s.%s", dataType, uniqueId)
			keys = append(keys, key)
		}
	}

	sort.Strings(keys)

	return keys
}

func (store *Store) GetDataSheetFiles(dataType string) []*DataSheetFile {
	files, ok := store.dataTypes[dataType]
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
		dataSheetFiles[i] = store.dataTypes[dataType][uniqueId]
	}

	return dataSheetFiles
}
