package localization

import (
	"encoding/xml"
	"errors"
	"github.com/new-world-tools/new-world-tools/internal"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrNotFound = errors.New("not found")

var reXml = regexp.MustCompile(`.xml`)

func New(root string) (*internal.Store[string], error) {
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

	localizationStore := internal.NewStore[string](func(key string) string {
		return strings.TrimPrefix(strings.ToLower(key), "@")
	})

	for _, xmlPath := range files {
		xmlFile, err := os.Open(xmlPath)
		if err != nil {
			return nil, err
		}

		var resources Resources
		err = xml.NewDecoder(xmlFile).Decode(&resources)
		if err != nil {
			return nil, err
		}

		for _, resource := range resources.Strings {
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
	XMLName xml.Name  `xml:"resources"`
	Strings []*String `xml:"string"`
}

type String struct {
	XMLName xml.Name `xml:"string"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",chardata"`
	Nil     bool     `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
}
