import type { Action, PlayerView } from './generated/engine';

/** Fetches the player's view of a game. Rejects with a readable message when the server does not answer 200. */
export function fetchGame(id: string): Promise<PlayerView> {
  return requestView(gameURL(id), undefined, `load game "${id}"`);
}

/** Sends an Action and returns the new player's view. Rejects with a readable message when the server does not answer 200. */
export function sendAction(id: string, action: Action): Promise<PlayerView> {
  return requestView(
    `${gameURL(id)}/actions`,
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(action),
    },
    `do "${action.type}"`,
  );
}

async function requestView(
  url: string,
  init: RequestInit | undefined,
  what: string,
): Promise<PlayerView> {
  const response = await fetch(url, init);
  if (!response.ok) {
    throw new Error(`Could not ${what} (HTTP ${response.status})`);
  }
  return response.json();
}

function gameURL(id: string): string {
  return `/api/games/${encodeURIComponent(id)}`;
}
