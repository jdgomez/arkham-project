<script lang="ts">
  import type { Action } from './generated/engine';

  let {
    actions,
    onaction,
    disabled = false,
  }: { actions: Action[]; onaction: (action: Action) => void; disabled?: boolean } = $props();
</script>

{#each actions as action (action.type)}
  <button type="button" {disabled} onclick={() => onaction(action)}>{action.type}</button>
{/each}
