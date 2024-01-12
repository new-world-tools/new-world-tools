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

func New(root string) (*store.Store[string, string], error) {
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

	localizationStore := store.NewStore[string, string](func(key string) string {
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
			if resource.Nil || resource.Key == "" {
				continue
			}

			if localizationStore.Has(resource.Key) && localizationStore.Get(resource.Key) != resource.Value {
				//return nil, fmt.Errorf("multiple values: %s = %q and %q", key, val, resource.Value)
			}
			//if !ok {w
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
	Key                     string `xml:"key,attr"`
	RelVersion              string `xml:"rel_version,attr"`
	RelVerstion             string `xml:"rel_verstion,attr"`
	Comment                 string `xml:"comment,attr"`
	VO                      string `xml:"VO,attr"`
	VOStatus                string `xml:"VO_Status,attr"`
	VOType                  string `xml:"VO_Type,attr"`
	CameraEnterBlendTime    string `xml:"cameraEnterBlendTime,attr"`
	CameraExitBlendTime     string `xml:"cameraExitBlendTime,attr"`
	CameraState             string `xml:"cameraState,attr"`
	CameraStateLookAt       string `xml:"cameraStateLookAt,attr"`
	CameraStateOrigin       string `xml:"cameraStateOrigin,attr"`
	DialogueNext            string `xml:"dialogue-next,attr"`
	DialoguePrompt          string `xml:"dialogue-prompt,attr"`
	End                     string `xml:"end,attr"`
	Gender                  string `xml:"gender,attr"`
	HideNearbyPlayerAvatars string `xml:"hideNearbyPlayerAvatars,attr"`
	HidePlayerAvatar        string `xml:"hidePlayerAvatar,attr"`
	LineId                  string `xml:"line_id,attr"`
	Location                string `xml:"location,attr"`
	Name                    string `xml:"name,attr"`
	OriginEnterBlendTime    string `xml:"originEnterBlendTime,attr"`
	QuestId                 string `xml:"quest_id,attr"`
	Speaker                 string `xml:"speaker,attr"`
	Start                   string `xml:"start,attr"`

	Value string `xml:",innerxml"`

	Nil bool `xml:"http://www.w3.org/2001/XMLSchema-instance nil,attr"`
}
