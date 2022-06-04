package localization

import (
	"encoding/xml"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var ErrNotFound = errors.New("not found")

var reXml = regexp.MustCompile(`.xml`)

func Get(root string) (map[string]string, error) {
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

	keyValue := map[string]string{}
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
			val, ok := keyValue[resource.Key]
			if ok && val == resource.Key {
				log.Fatalf("multiple keys: %s", resource.Key)
			}
			if !ok {
				keyValue[resource.Key] = resource.Value
			}
		}
	}

	return keyValue, nil
}

type Resources struct {
	XMLName xml.Name  `xml:"resources"`
	Strings []*String `xml:"string"`
}

type String struct {
	XMLName xml.Name `xml:"string"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",chardata"`
}
