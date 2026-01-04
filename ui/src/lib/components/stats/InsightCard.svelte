<script lang="ts">
  import type { Insight } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";
  import { UserX } from "lucide-svelte";

  interface Props {
    insights: Insight[];
  }

  let { insights }: Props = $props();

  function handleClick(insight: Insight) {
    if (insight.friendId) {
      currentPath.navigate(`/friends/${insight.friendId}`);
    }
  }
</script>

<div>
  <h3 class="text-sm font-semibold mb-3">Friends to Reconnect With</h3>
  {#if insights.length === 0}
    <p class="text-sm text-muted-foreground">
      You're keeping up with everyone! No friends need reconnecting.
    </p>
  {:else}
    <div class="space-y-2">
      {#each insights as insight}
        <button
          type="button"
          class="flex items-center gap-3 w-full text-left p-2 rounded-md hover:bg-muted/50 transition-colors"
          onclick={() => handleClick(insight)}
        >
          <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-orange-500/10">
            <UserX class="h-4 w-4 text-orange-500" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate">{insight.title}</p>
            <p class="text-xs text-muted-foreground">{insight.description}</p>
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>
