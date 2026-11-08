<script lang="ts">
  import ActionList from './ActionList.svelte';
  import { fetchGame, sendAction } from './api';
  import { DemoGameID, type Action, type PlayerView } from './generated/engine';

  let view = $state<PlayerView>();
  let error = $state<Error>();

  fetchGame(DemoGameID).then(show, fail);

  function show(next: PlayerView) {
    view = next;
    error = undefined;
  }

  function fail(cause: Error) {
    error = cause;
  }

  function act(action: Action) {
    sendAction(DemoGameID, action).then(show, fail);
  }
</script>

{#if error}
  <p role="alert">{error.message}</p>
{/if}
{#if view}
  <h1>Turn {view.turn} - {view.phase}</h1>
  <ActionList actions={view.availableActions} onaction={act} />
{:else if !error}
  <p>Loading...</p>
{/if}
