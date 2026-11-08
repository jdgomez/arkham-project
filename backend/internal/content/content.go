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
		investigators: map[string]Investigator{},
		ancientOnes:   map[string]AncientOne{},
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

// identity is the id and name that every kind of content must have.
type identity struct{ id, name string }

func identities[T any](items []T, of func(T) identity) []identity {
	result := make([]identity, len(items))
	for i, item := range items {
		result[i] = of(item)
	}
	return result
}

// checkIdentities fails on a missing id or name, or on an id already in seen.
func checkIdentities(kind string, entries []identity, seen map[string]bool) error {
	for i, entry := range entries {
		switch {
		case entry.id == "":
			return fmt.Errorf("%s: entry %d has no id", kind, i+1)
		case entry.name == "":
			return fmt.Errorf("%s: %q has no name", kind, entry.id)
		case seen[entry.id]:
			return fmt.Errorf("%s: duplicate id %q", kind, entry.id)
		}
		seen[entry.id] = true
	}
	return nil
}

func areaIdentity(area Area) identity { return identity{area.ID, area.Name} }

// validate checks ids and names, and that every reference points to content
// that exists. It also builds the indexes.
func (c *Content) validate() error {
	areaIDs := map[string]bool{}
	if err := checkIdentities("locations", identities(c.Locations, areaIdentity), areaIDs); err != nil {
		return err
	}
	if err := checkIdentities("streets", identities(c.Streets, areaIdentity), areaIDs); err != nil {
		return err
	}
	investigatorIdentity := func(i Investigator) identity { return identity{i.ID, i.Name} }
	if err := checkIdentities("investigators", identities(c.Investigators, investigatorIdentity), map[string]bool{}); err != nil {
		return err
	}
	ancientOneIdentity := func(a AncientOne) identity { return identity{a.ID, a.Name} }
	if err := checkIdentities("ancient ones", identities(c.AncientOnes, ancientOneIdentity), map[string]bool{}); err != nil {
		return err
	}
	if err := c.validateConnections(areaIDs); err != nil {
		return err
	}
	for _, investigator := range c.Investigators {
		if !areaIDs[investigator.StartingLocation] {
			return fmt.Errorf("%s: starting location %q does not exist", investigator.ID, investigator.StartingLocation)
		}
		c.investigators[investigator.ID] = investigator
	}
	for _, ancientOne := range c.AncientOnes {
		c.ancientOnes[ancientOne.ID] = ancientOne
	}
	return nil
}

// validateConnections checks that every connection leads to an existing area,
// and that a location and its street connect to each other. Other kinds of
// connection may be one-way, so nothing else is required to be reciprocal.
func (c *Content) validateConnections(areaIDs map[string]bool) error {
	for _, area := range c.areas {
		for _, connection := range area.Connections {
			if !areaIDs[connection] {
				return fmt.Errorf("%s: connection to %q, which does not exist", area.ID, connection)
			}
		}
	}
	if err := connectBack("location", c.Locations, "street", c.Streets); err != nil {
		return err
	}
	return connectBack("street", c.Streets, "location", c.Locations)
}

// connectBack checks that every area of the given kind that connects to an
// area of the other kind is connected back by it.
func connectBack(kind string, areas []Area, otherKind string, others []Area) error {
	otherByID := indexAreas(others)
	for _, area := range areas {
		for _, connection := range area.Connections {
			other, ok := otherByID[connection]
			if ok && !slices.Contains(other.Connections, area.ID) {
				return fmt.Errorf("%s %q connects to %s %q, which does not connect back", kind, area.ID, otherKind, other.ID)
			}
		}
	}
	return nil
}

func indexAreas(areas []Area) map[string]Area {
	index := map[string]Area{}
	for _, area := range areas {
		index[area.ID] = area
	}
	return index
}
