package engine

import (
	"slices"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

// Area is a place on the board with the areas it connects to.
type Area struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Connections []string `json:"connections"`
}

// InvestigatorView is an investigator as the player sees them.
type InvestigatorView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// AncientOneView is the Ancient One as the player sees it.
type AncientOneView struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PlayerView is what the player sees of a game. It leaves out what they must
// not see, such as the seed.
type PlayerView struct {
	Turn             int                `json:"turn"`
	Phase            Phase              `json:"phase"`
	Board            []Area             `json:"board"`
	Investigators    []InvestigatorView `json:"investigators"`
	AncientOne       AncientOneView     `json:"ancientOne"`
	AvailableActions []Action           `json:"availableActions"`
}

// ViewOf composes the player's view of a game: the content its ids refer to,
// the game state and the Available Actions.
func ViewOf(loaded content.Content, game Game) PlayerView {
	return PlayerView{
		Turn:             game.Turn,
		Phase:            game.Phase,
		Board:            boardOf(loaded),
		Investigators:    investigatorsOf(loaded, game),
		AncientOne:       ancientOneOf(loaded, game),
		AvailableActions: AvailableActions(game),
	}
}

func boardOf(loaded content.Content) []Area {
	board := []Area{}
	for _, area := range slices.Concat(loaded.Locations, loaded.Streets) {
		board = append(board, Area(area))
	}
	return board
}

func investigatorsOf(loaded content.Content, game Game) []InvestigatorView {
	investigators := []InvestigatorView{}
	for _, state := range game.Investigators {
		index := slices.IndexFunc(loaded.Investigators, func(i content.Investigator) bool { return i.ID == state.ID })
		if index < 0 {
			continue
		}
		investigators = append(investigators,
			InvestigatorView{ID: state.ID, Name: loaded.Investigators[index].Name, Location: state.Location})
	}
	return investigators
}

func ancientOneOf(loaded content.Content, game Game) AncientOneView {
	index := slices.IndexFunc(loaded.AncientOnes, func(a content.AncientOne) bool { return a.ID == game.AncientOne })
	if index < 0 {
		return AncientOneView{}
	}
	return AncientOneView{ID: game.AncientOne, Name: loaded.AncientOnes[index].Name}
}
