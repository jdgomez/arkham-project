package content_test

import (
	"strings"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

// parts are the arguments of content.New, valid until a test case breaks them.
type parts struct {
	locations     []content.Area
	streets       []content.Area
	investigators []content.Investigator
	ancientOnes   []content.AncientOne
}

func validParts() parts {
	return parts{
		locations: []content.Area{
			{ID: "library", Name: "Library", Connections: []string{"streets"}},
			{ID: "docks", Name: "Docks", Connections: []string{"streets"}},
		},
		streets: []content.Area{
			{ID: "streets", Name: "Streets", Connections: []string{"library", "docks", "other-streets"}},
			{ID: "other-streets", Name: "Other Streets", Connections: nil},
		},
		investigators: []content.Investigator{{ID: "mandy", Name: "Mandy", StartingLocation: "library"}},
		ancientOnes:   []content.AncientOne{{ID: "azathoth", Name: "Azathoth"}},
	}
}

func (p parts) build() (content.Content, error) {
	return content.New(p.locations, p.streets, p.investigators, p.ancientOnes)
}

func TestNewAcceptsValidContent(t *testing.T) {
	if _, err := validParts().build(); err != nil {
		t.Fatalf("New returned error: %v", err)
	}
}

func TestNewDoesNotRequireStreetToStreetConnectionsToBeReciprocal(t *testing.T) {
	p := validParts()
	p.streets[0].Connections = append(p.streets[0].Connections, "other-streets")
	p.streets[1].Connections = nil

	if _, err := p.build(); err != nil {
		t.Fatalf("New returned error for a one-way street connection: %v", err)
	}
}

func TestNewFailsOnInvalidContent(t *testing.T) {
	tests := []struct {
		name   string
		break_ func(*parts)
		want   []string
	}{
		{"two locations with the same id", func(p *parts) { p.locations[1].ID = "library" }, []string{"duplicate", "library"}},
		{"two streets with the same id", func(p *parts) { p.streets[1].ID = "streets" }, []string{"duplicate", "streets"}},
		{"a location and a street with the same id", func(p *parts) { p.streets[0].ID = "docks" }, []string{"duplicate", "docks"}},
		{
			"two investigators with the same id",
			func(p *parts) { p.investigators = append(p.investigators, p.investigators[0]) },
			[]string{"duplicate", "mandy"},
		},
		{
			"two Ancient Ones with the same id",
			func(p *parts) { p.ancientOnes = append(p.ancientOnes, p.ancientOnes[0]) },
			[]string{"duplicate", "azathoth"},
		},
		{"a location without id", func(p *parts) { p.locations[0].ID = "" }, []string{"locations", "id"}},
		{"a location without name", func(p *parts) { p.locations[0].Name = "" }, []string{"library", "name"}},
		{"a street without id", func(p *parts) { p.streets[0].ID = "" }, []string{"streets", "id"}},
		{"a street without name", func(p *parts) { p.streets[0].Name = "" }, []string{"streets", "name"}},
		{"an investigator without id", func(p *parts) { p.investigators[0].ID = "" }, []string{"investigators", "id"}},
		{"an investigator without name", func(p *parts) { p.investigators[0].Name = "" }, []string{"mandy", "name"}},
		{"an Ancient One without id", func(p *parts) { p.ancientOnes[0].ID = "" }, []string{"ancient ones", "id"}},
		{"an Ancient One without name", func(p *parts) { p.ancientOnes[0].Name = "" }, []string{"azathoth", "name"}},
		{
			"a location that connects to its street while the street does not connect back",
			func(p *parts) { p.streets[0].Connections = []string{"docks", "other-streets"} },
			[]string{"library", "streets", "back"},
		},
		{
			"a street that connects to a location while the location does not connect back",
			func(p *parts) { p.locations[1].Connections = nil },
			[]string{"streets", "docks", "back"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := validParts()
			test.break_(&p)

			_, err := p.build()

			if err == nil {
				t.Fatal("New returned no error, want one")
			}
			for _, want := range test.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want it to mention %q", err, want)
				}
			}
		})
	}
}
