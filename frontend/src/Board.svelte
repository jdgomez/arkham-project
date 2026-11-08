<script lang="ts">
  import type { Area, InvestigatorView } from './generated/engine';

  let { board, investigators }: { board: Area[]; investigators: InvestigatorView[] } = $props();

  const nameOf = (id: string) => board.find((area) => area.id === id)?.name ?? id;
</script>

<ul>
  {#each board as area (area.id)}
    <li>
      <section aria-labelledby="area-{area.id}">
        <h2 id="area-{area.id}">{area.name}</h2>
        <ul aria-label="Connections">
          {#each area.connections as connection (connection)}
            <li>{nameOf(connection)}</li>
          {/each}
        </ul>
        {#each investigators.filter((investigator) => investigator.location === area.id) as investigator (investigator.id)}
          <p>{investigator.name}</p>
        {/each}
      </section>
    </li>
  {/each}
</ul>
