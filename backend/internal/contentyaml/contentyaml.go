// Package contentyaml reads the game content from the YAML files of the
// content directory (D13), one file per content type.
package contentyaml

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"

	"gopkg.in/yaml.v3"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

// Load reads and validates all the content files in the file system. Unknown
// fields, empty files and files with more than one YAML document are errors.
func Load(files fs.FS) (content.Content, error) {
	locations, err := decodeFile[content.Area](files, "locations.yaml")
	if err != nil {
		return content.Content{}, err
	}
	streets, err := decodeFile[content.Area](files, "streets.yaml")
	if err != nil {
		return content.Content{}, err
	}
	investigators, err := decodeFile[content.Investigator](files, "investigators.yaml")
	if err != nil {
		return content.Content{}, err
	}
	ancientOnes, err := decodeFile[content.AncientOne](files, "ancient-ones.yaml")
	if err != nil {
		return content.Content{}, err
	}
	return content.New(locations, streets, investigators, ancientOnes)
}

// decodeFile reads the one YAML document of the file: a non-empty list.
func decodeFile[T any](files fs.FS, name string) ([]T, error) {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		return nil, fmt.Errorf("reading content: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var entries []T
	if err := decoder.Decode(&entries); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("%s: no entries", name)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("%s: more than one document", name)
	}
	return entries, nil
}
