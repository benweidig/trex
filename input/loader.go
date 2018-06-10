package input

import (
	"encoding/json"
	"errors"

	yaml "gopkg.in/yaml.v2"
)

// Load loads/unmarshals the bytes into the correct map according to the filetype.
// If no filltype is set it tries to sniff the actual content.
func Load(fileType FileType, data []byte) (interface{}, error) {

	switch fileType {
	case FileTypeJSON:
		return loadFromJSON(data)

	case FileTypeYAML:
		return loadFromYAML(data)
	}

	// we couldn't detect the filetype so we have to sniff it
	first := string(data[0])
	if first == "[" || first == "{" {
		return loadFromJSON(data)
	}

	// TODO: sniff content
	return nil, errors.New("Failed to load")
}

func loadFromJSON(data []byte) (interface{}, error) {
	var rawJSON interface{}
	err := json.Unmarshal(data, &rawJSON)
	return rawJSON, err
}

func loadFromYAML(data []byte) (interface{}, error) {
	var raw interface{}
	err := yaml.Unmarshal([]byte(data), &raw)
	return raw, err
}
