package engine

import "github.com/jdgomez/arkham-project/backend/internal/content"

// InvestigatorView is an investigator as the player sees them.
type InvestigatorView struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// PlayerView is what the player sees of a game. It leaves out what they must
// not see, such as the seed.
type PlayerView struct {
	Turn             int                `json:"turn"`
	Phase            Phase              `json:"phase"`
	Board            []content.Area     `json:"board"`
	Investigators    []InvestigatorView `json:"investigators"`
	AncientOne       content.AncientOne `json:"ancientOne"`
	AvailableActions []Action           `json:"availableActions"`
}

// ViewOf composes the player's view of a game: the content its ids refer to,
// the game state and the Available Actions.
func ViewOf(loaded content.Content, game Game) PlayerView {
	return PlayerView{
		Turn:             game.Turn,
		Phase:            game.Phase,
		Board:            loaded.Areas(),
		Investigators:    investigatorsOf(loaded, game),
		AncientOne:       ancientOneOf(loaded, game),
		AvailableActions: AvailableActions(game),
	}
}

func investigatorsOf(loaded content.Content, game Game) []InvestigatorView {
	investigators := []InvestigatorView{}
	for _, state := range game.Investigators {
		investigator, ok := loaded.Investigator(state.ID)
		if !ok {
			continue
		}
		investigators = append(investigators,
			InvestigatorView{ID: state.ID, Name: investigator.Name, Location: state.Location})
	}
	return investigators
}

func ancientOneOf(loaded content.Content, game Game) content.AncientOne {
	ancientOne, _ := loaded.AncientOne(game.AncientOne)
	return ancientOne
}
