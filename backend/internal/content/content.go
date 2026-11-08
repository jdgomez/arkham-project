// Package content loads the static game content (the map, the investigators,
// the Ancient Ones) from the YAML files of the content directory.
package content

import (
	"bytes"
	"fmt"
	"io/fs"
	"slices"

	"gopkg.in/yaml.v3"
)

// Area is a place on the board that other areas connect to.
type Area struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Connections []string `yaml:"connections"`
}

// Investigator is a playable character.
type Investigator struct {
	ID               string `yaml:"id"`
	Name             string `yaml:"name"`
	StartingLocation string `yaml:"starting_location"`
}

// AncientOne is the great old one the investigators face.
type AncientOne struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

// Content is everything the game loads once and never changes.
type Content struct {
	Locations     []Area
	Streets       []Area
	Investigators []Investigator
	AncientOnes   []AncientOne
}

// Load reads all the content files in the file system, one per content type.
// Unknown fields are errors.
func Load(files fs.FS) (Content, error) {
	var content Content
	for _, file := range []struct {
		name   string
		target any
	}{
		{"locations.yaml", &content.Locations},
		{"streets.yaml", &content.Streets},
		{"investigators.yaml", &content.Investigators},
		{"ancient-ones.yaml", &content.AncientOnes},
	} {
		if err := decodeFile(files, file.name, file.target); err != nil {
			return Content{}, err
		}
	}
	if err := content.validate(); err != nil {
		return Content{}, err
	}
	return content, nil
}

// validate checks that every reference points to content that exists.
func (c Content) validate() error {
	areas := map[string]bool{}
	for _, area := range slices.Concat(c.Locations, c.Streets) {
		areas[area.ID] = true
	}
	for _, area := range slices.Concat(c.Locations, c.Streets) {
		for _, connection := range area.Connections {
			if !areas[connection] {
				return fmt.Errorf("%s: connection to %q, which does not exist", area.ID, connection)
			}
		}
	}
	for _, investigator := range c.Investigators {
		if !areas[investigator.StartingLocation] {
			return fmt.Errorf("%s: starting location %q does not exist", investigator.ID, investigator.StartingLocation)
		}
	}
	return nil
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
