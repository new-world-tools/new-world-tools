package localization

import (
	"encoding/xml"
	"errors"
	"github.com/new-world-tools/new-world-tools/store"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrNotFound = errors.New("not found")

var reXml = regexp.MustCompile(`.xml`)

func New(root string) (*store.Store[string], error) {
	files := []string{}

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil || !reXml.MatchString(info.Name()) {
			return err
		}

		path, err = filepath.Abs(filepath.Clean(path))
		if err != nil {
			return err
		}

		files = append(files, path)

		return nil
	})
	if err != nil {
		return nil, err
	}

	localizationStore := store.NewStore[string](func(key string) string {
		return strings.TrimPrefix(strings.ToLower(key), "@")
	})

	for _, xmlPath := range files {
		xmlFile, err := os.Open(xmlPath)
		if err != nil {
			return nil, err
		}

		var isLegacyFile bool

		var resources Resources
		err = xml.NewDecoder(xmlFile).Decode(&resources)
		if err != nil {
			if strings.Contains(err.Error(), "expected element type") {
				isLegacyFile = true
			}
			if !isLegacyFile {
				return nil, err
			}
		}
		xmlFile.Close()

		if isLegacyFile {
			continue
		}

		for _, resource := range resources.Resources {
			if resource.Nil {
				continue
			}

			if localizationStore.Has(resource.Key) && localizationStore.Get(resource.Key) != resource.Value {
				//return nil, fmt.Errorf("multiple values: %s = %q and %q", key, val, resource.Value)
			}
			//if !ok {
			localizationStore.Add(resource.Key, resource.Value)
			//}
		}
	}

	return localizationStore, nil
}

type Resources struct {
	XMLName   xml.Name   `xml:"resources"`
	Resources []Resource `xml:"string"`
}

type Resource struct {
	Key        string `xml:"key,attr"`
	RelVersion string `xml:"rel_version,attr"`
	Comment    string `xml:"comment,attr"`
	Value      string `xml:",innerxml"`

	Nil     bool     `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
}
