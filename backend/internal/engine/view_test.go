package engine_test

import (
	"reflect"
	"testing"

	"github.com/jdgomez/arkham-project/backend/internal/content"

	"github.com/jdgomez/arkham-project/backend/internal/engine"
)

func TestDemoViewShowsTheBoardMandyAndAzathoth(t *testing.T) {
	view := viewOf(t, demoGame(t))

	wantBoard := []content.Area{
		{ID: "unvisited-isle", Name: "Unvisited Isle", Connections: []string{"merchant-district-streets"}},
		{ID: "river-docks", Name: "River Docks", Connections: []string{"merchant-district-streets"}},
		{ID: "the-unnamable", Name: "The Unnamable", Connections: []string{"merchant-district-streets"}},
		{ID: "administration", Name: "Administration", Connections: []string{"miskatonic-u-streets"}},
		{ID: "library", Name: "Library", Connections: []string{"miskatonic-u-streets"}},
		{ID: "science-building", Name: "Science Building", Connections: []string{"miskatonic-u-streets"}},
		{
			ID: "merchant-district-streets", Name: "Merchant District Streets",
			Connections: []string{"unvisited-isle", "river-docks", "the-unnamable", "miskatonic-u-streets"},
		},
		{
			ID: "miskatonic-u-streets", Name: "Miskatonic U. Streets",
			Connections: []string{"administration", "library", "science-building", "merchant-district-streets"},
		},
	}
	if !reflect.DeepEqual(view.Board, wantBoard) {
		t.Errorf("board = %+v, want %+v", view.Board, wantBoard)
	}
	wantInvestigators := []engine.InvestigatorView{{ID: "mandy-thompson", Name: "Mandy Thompson", Location: "library"}}
	if !reflect.DeepEqual(view.Investigators, wantInvestigators) {
		t.Errorf("investigators = %+v, want %+v", view.Investigators, wantInvestigators)
	}
	wantAncientOne := content.AncientOne{ID: "azathoth", Name: "Azathoth"}
	if view.AncientOne != wantAncientOne {
		t.Errorf("ancient one = %+v, want %+v", view.AncientOne, wantAncientOne)
	}
}
