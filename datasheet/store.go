package datasheet

import (
	"github.com/new-world-tools/new-world-tools/store"
	"sort"
)

func NewStore(dataTableDir string) (*Store, error) {
	dataSheetFiles, err := FindAll(dataTableDir)
	if err != nil {
		return nil, err
	}

	store := &Store{
		store: store.NewCaseInsensitiveStore[*DataSheetFile](),
	}

	for _, dataSheetFile := range dataSheetFiles {
		store.store.Add(dataSheetFile.GetPath(), dataSheetFile)
	}

	return store, nil
}

type Store struct {
	store *store.Store[string, *DataSheetFile]
}

func (store *Store) GetByUniqueId(typ string, uniqueId string) (*DataSheetFile, error) {
	keys := store.store.GetKeys()

	for _, key := range keys {
		dsFile := store.store.Get(key)

		meta, err := dsFile.GetMeta()
		if err != nil {
			return nil, err
		}

		if meta.Type == typ && meta.UniqueId == uniqueId {
			return dsFile, nil
		}
	}

	return nil, nil
}

func (store *Store) GetByType(typ string) ([]*DataSheetFile, error) {
	keys := store.store.GetKeys()

	dsFiles := []*DataSheetFile{}
	for _, key := range keys {
		dsFile := store.store.Get(key)

		meta, err := dsFile.GetMeta()
		if err != nil {
			return nil, err
		}

		if meta.Type == typ {
			dsFiles = append(dsFiles, dsFile)
		}
	}

	return dsFiles, nil
}

func (store *Store) GetAll() []*DataSheetFile {
	keys := store.store.GetKeys()

	dsFiles := make([]*DataSheetFile, len(keys))
	for i, key := range keys {
		dsFile := store.store.Get(key)
		dsFiles[i] = dsFile
	}

	return dsFiles
}

func (store *Store) GetTypes() ([]string, error) {
	keys := store.store.GetKeys()

	typesMap := map[string]bool{}
	for _, key := range keys {
		dsFile := store.store.Get(key)

		meta, err := dsFile.GetMeta()
		if err != nil {
			return nil, err
		}

		typesMap[meta.Type] = true
	}

	types := make([]string, len(typesMap))
	var i int
	for typ, _ := range typesMap {
		types[i] = typ
		i++
	}

	sort.Strings(types)

	return types, nil
}
