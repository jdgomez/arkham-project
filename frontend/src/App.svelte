<script lang="ts">
  import { fetchGame } from './api';
  import { DemoGameID } from './generated/engine';

  const game = fetchGame(DemoGameID);
</script>

{#await game}
  <p>Loading...</p>
{:then view}
  <h1>Turn {view.turn} - {view.phase}</h1>
{:catch error}
  <p role="alert">{error.message}</p>
{/await}
