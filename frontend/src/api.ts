import type { PlayerView } from './generated/engine';

/** Fetches the player's view of a game. Rejects with a readable message when the server does not answer 200. */
export async function fetchGame(id: string): Promise<PlayerView> {
  const response = await fetch(`/api/games/${id}`);
  if (!response.ok) {
    throw new Error(`Could not load game "${id}" (HTTP ${response.status})`);
  }
  return response.json();
}
