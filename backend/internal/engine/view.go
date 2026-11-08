package engine

import (
	"fmt"

	"github.com/jdgomez/arkham-project/backend/internal/content"
)

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
// the game state and the Available Actions. It fails if the game refers to
// content that does not exist.
func ViewOf(loaded content.Content, game Game) (PlayerView, error) {
	investigators, err := investigatorsOf(loaded, game)
	if err != nil {
		return PlayerView{}, err
	}
	ancientOne, ok := loaded.AncientOne(game.AncientOne)
	if !ok {
		return PlayerView{}, fmt.Errorf("composing view: Ancient One %q does not exist", game.AncientOne)
	}
	return PlayerView{
		Turn:             game.Turn,
		Phase:            game.Phase,
		Board:            loaded.Areas(),
		Investigators:    investigators,
		AncientOne:       ancientOne,
		AvailableActions: AvailableActions(game),
	}, nil
}

func investigatorsOf(loaded content.Content, game Game) ([]InvestigatorView, error) {
	investigators := []InvestigatorView{}
	for _, state := range game.Investigators {
		investigator, ok := loaded.Investigator(state.ID)
		if !ok {
			return nil, fmt.Errorf("composing view: investigator %q does not exist", state.ID)
		}
		investigators = append(investigators,
			InvestigatorView{ID: state.ID, Name: investigator.Name, Location: state.Location})
	}
	return investigators, nil
}
