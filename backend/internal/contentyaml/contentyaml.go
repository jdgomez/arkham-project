// Package contentyaml reads the game content from the YAML files of the
// content directory (D13), one file per content type.
package contentyaml

import (
	"bytes"
	"fmt"
	"io/fs"

	"gopkg.in/yaml.v3"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

// Load reads and validates all the content files in the file system. Unknown
// fields are errors.
func Load(files fs.FS) (content.Content, error) {
	var (
		locations, streets []content.Area
		investigators      []content.Investigator
		ancientOnes        []content.AncientOne
	)
	for _, file := range []struct {
		name   string
		target any
	}{
		{"locations.yaml", &locations},
		{"streets.yaml", &streets},
		{"investigators.yaml", &investigators},
		{"ancient-ones.yaml", &ancientOnes},
	} {
		if err := decodeFile(files, file.name, file.target); err != nil {
			return content.Content{}, err
		}
	}
	return content.New(locations, streets, investigators, ancientOnes)
}

func decodeFile(files fs.FS, name string, target any) error {
	data, err := fs.ReadFile(files, name)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}
