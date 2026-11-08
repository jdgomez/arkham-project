package engine

// PlayerView is what the player sees of a game. It leaves out what they must
// not see, such as the seed.
type PlayerView struct {
	Turn  int   `json:"turn"`
	Phase Phase `json:"phase"`
}

// ViewOf composes the player's view of a game.
func ViewOf(game Game) PlayerView {
	return PlayerView{Turn: game.Turn, Phase: game.Phase}
}
