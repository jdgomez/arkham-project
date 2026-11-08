<script lang="ts">
  import type { Action } from './generated/engine';

  let { actions, onaction }: { actions: Action[]; onaction: (action: Action) => void } = $props();
</script>

{#each actions as action (action.type)}
  <button type="button" onclick={() => onaction(action)}>{action.type}</button>
{/each}
