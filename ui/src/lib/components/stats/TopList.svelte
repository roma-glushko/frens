<script lang="ts">
  import type { RankedItem } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";

  interface Props {
    title: string;
    items: RankedItem[];
    type: "friends" | "locations" | "tags";
    emptyMessage?: string;
  }

  let { title, items, type, emptyMessage = "No data yet" }: Props = $props();

  function handleClick(item: RankedItem) {
    if (type === "friends") {
      currentPath.navigate(`/friends/${item.id}`);
    } else if (type === "locations") {
      currentPath.navigate(`/locations/${item.id}`);
    }
  }

  function isClickable(): boolean {
    return type === "friends" || type === "locations";
  }
</script>

<div>
  <h3 class="text-sm font-semibold mb-3">{title}</h3>
  {#if items.length === 0}
    <p class="text-sm text-muted-foreground">{emptyMessage}</p>
  {:else}
    <div class="space-y-2">
      {#each items as item, i}
        {#if isClickable()}
          <button
            type="button"
            class="flex items-center justify-between w-full text-left p-2 rounded-md hover:bg-muted/50 transition-colors"
            onclick={() => handleClick(item)}
          >
            <div class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground w-4">{i + 1}.</span>
              <span class="text-sm">{item.name}</span>
            </div>
            <span class="text-xs text-muted-foreground">{item.count}</span>
          </button>
        {:else}
          <div class="flex items-center justify-between p-2">
            <div class="flex items-center gap-2">
              <span class="text-xs text-muted-foreground w-4">{i + 1}.</span>
              <span class="text-sm">#{item.name}</span>
            </div>
            <span class="text-xs text-muted-foreground">{item.count}</span>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>
