// Package content holds the static game content (the map, the investigators,
// the Ancient Ones) as plain data. It imports no file, OS or YAML package, so
// the engine can depend on it; contentyaml reads the YAML files into it.
package content

import (
	"fmt"
	"slices"
)

// Area is a place on the board that other areas connect to.
type Area struct {
	ID          string   `json:"id" yaml:"id"`
	Name        string   `json:"name" yaml:"name"`
	Connections []string `json:"connections" yaml:"connections"`
}

// Investigator is a playable character.
type Investigator struct {
	ID               string `yaml:"id"`
	Name             string `yaml:"name"`
	StartingLocation string `yaml:"starting_location"`
}

// AncientOne is the great old one the investigators face.
type AncientOne struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

// Content is everything the game loads once and never changes. Build it with
// New, which validates it and indexes it by id.
type Content struct {
	Locations     []Area
	Streets       []Area
	Investigators []Investigator
	AncientOnes   []AncientOne

	areas         []Area
	investigators map[string]Investigator
	ancientOnes   map[string]AncientOne
}

// New validates the content and indexes it by id.
func New(locations, streets []Area, investigators []Investigator, ancientOnes []AncientOne) (Content, error) {
	areas := slices.Concat(locations, streets)
	content := Content{
		Locations:     locations,
		Streets:       streets,
		Investigators: investigators,
		AncientOnes:   ancientOnes,
		areas:         areas,
	}
	if err := content.validate(); err != nil {
		return Content{}, err
	}
	return content, nil
}

// Areas lists the locations and then the streets: the whole board.
func (c Content) Areas() []Area {
	return c.areas
}

// Investigator returns the investigator with the id, and whether it exists.
func (c Content) Investigator(id string) (Investigator, bool) {
	investigator, ok := c.investigators[id]
	return investigator, ok
}

// AncientOne returns the Ancient One with the id, and whether it exists.
func (c Content) AncientOne(id string) (AncientOne, bool) {
	ancientOne, ok := c.ancientOnes[id]
	return ancientOne, ok
}

// validate checks that every reference points to content that exists.
func (c *Content) validate() error {
	areas := map[string]bool{}
	for _, area := range c.areas {
		areas[area.ID] = true
	}
	for _, area := range c.areas {
		for _, connection := range area.Connections {
			if !areas[connection] {
				return fmt.Errorf("%s: connection to %q, which does not exist", area.ID, connection)
			}
		}
	}
	c.investigators = map[string]Investigator{}
	for _, investigator := range c.Investigators {
		if !areas[investigator.StartingLocation] {
			return fmt.Errorf("%s: starting location %q does not exist", investigator.ID, investigator.StartingLocation)
		}
		c.investigators[investigator.ID] = investigator
	}
	c.ancientOnes = map[string]AncientOne{}
	for _, ancientOne := range c.AncientOnes {
		c.ancientOnes[ancientOne.ID] = ancientOne
	}
	return nil
}
