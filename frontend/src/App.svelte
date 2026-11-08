<script lang="ts">
  import Board from './Board.svelte';
  import ActionList from './ActionList.svelte';
  import { fetchGame, sendAction } from './api';
  import { DemoGameID, type Action, type PlayerView } from './generated/engine';

  let view = $state<PlayerView>();
  let error = $state<Error>();
  // One request at a time: the Actions stay disabled until the answer arrives, so a late response never overwrites a newer one.
  let pending = $state(false);

  fetchGame(DemoGameID).then(show, fail);

  function show(next: PlayerView) {
    view = next;
    error = undefined;
  }

  function fail(cause: Error) {
    error = cause;
  }

  function act(action: Action) {
    pending = true;
    sendAction(DemoGameID, action)
      .then(show, fail)
      .finally(() => (pending = false));
  }
</script>

{#if error}
  <p role="alert">{error.message}</p>
{/if}
{#if view}
  <h1>Turn {view.turn} - {view.phase}</h1>
  <p>Ancient One: {view.ancientOne.name}</p>
  <Board board={view.board} investigators={view.investigators} />
  <ActionList actions={view.availableActions} onaction={act} disabled={pending} />
{:else if !error}
  <p>Loading...</p>
{/if}
