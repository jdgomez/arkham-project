package contentyaml_test

import (
	"os"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/jdgomez/arkham-project/backend/internal/content"
	"github.com/jdgomez/arkham-project/backend/internal/contentyaml"
)

// loadAll loads the real content directory of the repository: the validation
// test of D13 runs over all of content/.
func loadAll(t *testing.T) content.Content {
	t.Helper()
	loaded, err := contentyaml.Load(os.DirFS("../../../content"))
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	return loaded
}

func ids(areas []content.Area) []string {
	result := make([]string, len(areas))
	for i, area := range areas {
		result[i] = area.ID
	}
	return result
}

func connections(t *testing.T, areas []content.Area, id string) []string {
	t.Helper()
	for _, area := range areas {
		if area.ID == id {
			return area.Connections
		}
	}
	t.Fatalf("no area %q in %v", id, ids(areas))
	return nil
}

func TestTheMapHasSixLocationsAndTwoStreets(t *testing.T) {
	loaded := loadAll(t)

	wantLocations := []string{
		"unvisited-isle", "river-docks", "the-unnamable",
		"administration", "library", "science-building",
	}
	wantStreets := []string{"merchant-district-streets", "miskatonic-u-streets"}
	if got := ids(loaded.Locations); !slices.Equal(got, wantLocations) {
		t.Errorf("locations = %v, want %v", got, wantLocations)
	}
	if got := ids(loaded.Streets); !slices.Equal(got, wantStreets) {
		t.Errorf("streets = %v, want %v", got, wantStreets)
	}
}

func TestEachLocationConnectsToItsDistrictsStreet(t *testing.T) {
	loaded := loadAll(t)
	tests := []struct {
		location string
		street   string
	}{
		{"unvisited-isle", "merchant-district-streets"},
		{"river-docks", "merchant-district-streets"},
		{"the-unnamable", "merchant-district-streets"},
		{"administration", "miskatonic-u-streets"},
		{"library", "miskatonic-u-streets"},
		{"science-building", "miskatonic-u-streets"},
	}
	for _, test := range tests {
		t.Run(test.location, func(t *testing.T) {
			if got := connections(t, loaded.Locations, test.location); !slices.Equal(got, []string{test.street}) {
				t.Errorf("%s connects to %v, want only %s", test.location, got, test.street)
			}
			if got := connections(t, loaded.Streets, test.street); !slices.Contains(got, test.location) {
				t.Errorf("%s connects to %v, want it to include %s", test.street, got, test.location)
			}
		})
	}
}

func TestTheTwoStreetsConnectToEachOther(t *testing.T) {
	loaded := loadAll(t)

	if got := connections(t, loaded.Streets, "merchant-district-streets"); !slices.Contains(got, "miskatonic-u-streets") {
		t.Errorf("Merchant District Streets connects to %v, want Miskatonic U. Streets", got)
	}
	if got := connections(t, loaded.Streets, "miskatonic-u-streets"); !slices.Contains(got, "merchant-district-streets") {
		t.Errorf("Miskatonic U. Streets connects to %v, want Merchant District Streets", got)
	}
}

// validFiles is a small valid content directory that the error cases break in one place.
func validFiles() fstest.MapFS {
	return fstest.MapFS{
		"locations.yaml":     {Data: []byte("- id: library\n  name: Library\n  connections: [streets]\n")},
		"streets.yaml":       {Data: []byte("- id: streets\n  name: Streets\n  connections: [library]\n")},
		"investigators.yaml": {Data: []byte("- id: mandy\n  name: Mandy\n  starting_location: library\n")},
		"ancient-ones.yaml":  {Data: []byte("- id: azathoth\n  name: Azathoth\n")},
	}
}

func TestLoadAcceptsValidContent(t *testing.T) {
	if _, err := contentyaml.Load(validFiles()); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
}

func TestLoadFailsOnAnUnknownFieldNamingTheFileAndTheField(t *testing.T) {
	files := validFiles()
	files["locations.yaml"] = &fstest.MapFile{Data: []byte("- id: library\n  name: Library\n  conections: [streets]\n")}

	_, err := contentyaml.Load(files)

	if err == nil {
		t.Fatal("Load returned no error, want one about the unknown field")
	}
	for _, want := range []string{"locations.yaml", "conections"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want it to name %q", err, want)
		}
	}
}

func TestLoadFailsOnBrokenCrossReferences(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		content string
		want    []string
	}{
		{
			"a location connected to an area that does not exist",
			"locations.yaml",
			"- id: library\n  name: Library\n  connections: [nowhere]\n",
			[]string{"library", "nowhere"},
		},
		{
			"a street connected to an area that does not exist",
			"streets.yaml",
			"- id: streets\n  name: Streets\n  connections: [nowhere]\n",
			[]string{"streets", "nowhere"},
		},
		{
			"an investigator starting at a location that does not exist",
			"investigators.yaml",
			"- id: mandy\n  name: Mandy\n  starting_location: nowhere\n",
			[]string{"mandy", "nowhere"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files := validFiles()
			files[test.file] = &fstest.MapFile{Data: []byte(test.content)}

			_, err := contentyaml.Load(files)

			if err == nil {
				t.Fatal("Load returned no error, want one about the broken reference")
			}
			for _, want := range test.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want it to name %q", err, want)
				}
			}
		})
	}
}

var contentFiles = []string{"locations.yaml", "streets.yaml", "investigators.yaml", "ancient-ones.yaml"}

func TestLoadFailsOnAFileWithNoEntries(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"an empty file", ""},
		{"only a comment", "# nothing here yet\n"},
		{"an empty list", "[]\n"},
	}
	for _, file := range contentFiles {
		for _, test := range tests {
			t.Run(file+" with "+test.name, func(t *testing.T) {
				files := validFiles()
				files[file] = &fstest.MapFile{Data: []byte(test.data)}

				_, err := contentyaml.Load(files)

				if err == nil {
					t.Fatal("Load returned no error, want one about the missing entries")
				}
				for _, want := range []string{file, "no entries"} {
					if !strings.Contains(err.Error(), want) {
						t.Errorf("error = %q, want it to name %q", err, want)
					}
				}
			})
		}
	}
}

func TestLoadFailsOnMoreThanOneDocument(t *testing.T) {
	for _, file := range contentFiles {
		t.Run(file, func(t *testing.T) {
			files := validFiles()
			files[file] = &fstest.MapFile{Data: slices.Concat(files[file].Data, []byte("---\n"), files[file].Data)}

			_, err := contentyaml.Load(files)

			if err == nil {
				t.Fatal("Load returned no error, want one about the extra document")
			}
			for _, want := range []string{file, "more than one document"} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want it to name %q", err, want)
				}
			}
		})
	}
}
