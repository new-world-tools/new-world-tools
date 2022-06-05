package localization

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var ErrNotFound = errors.New("not found")

var reXml = regexp.MustCompile(`.xml`)

type Localization struct {
	values map[string]string
}

func (loc *Localization) Has(key string) bool {
	_, ok := loc.values[strings.TrimPrefix(strings.ToLower(key), "@")]
	return ok
}

func (loc *Localization) Get(key string) string {
	val, _ := loc.values[strings.TrimPrefix(strings.ToLower(key), "@")]
	return val
}

func New(root string) (*Localization, error) {
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

	localizationData := &Localization{
		values: map[string]string{},
	}

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

			key := strings.ToLower(resource.Key)

			val, ok := localizationData.values[key]
			if ok && val != resource.Value {
				return nil, fmt.Errorf("multiple values: %s = %q and %q", key, val, resource.Value)
			}
			if !ok {
				localizationData.values[key] = resource.Value
			}
		}
	}

	return localizationData, nil
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
