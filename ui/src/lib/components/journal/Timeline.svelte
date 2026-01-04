<script lang="ts">
  import { cn } from "$lib/utils";
  import { Calendar, MapPin, Hash, MessageSquare, Activity, UserPlus, MapPinPlus } from "lucide-svelte";
  import { currentPath } from "$lib/stores/router.svelte";

  export interface TimelineEntry {
    id: string;
    type: "note" | "activity" | "friend_added" | "location_added";
    content: string;
    date: Date;
    tags?: string[];
    location?: string;
    entityId?: string;
    entityName?: string;
  }

  interface Props {
    entries?: TimelineEntry[];
  }

  let { entries = [] }: Props = $props();

  function viewLocation(id: string) {
    currentPath.navigate(`/locations/${id}`);
  }

  function viewFriend(id: string) {
    currentPath.navigate(`/friends/${id}`);
  }

  function formatTime(date: Date): string {
    return date.toLocaleTimeString("en-US", {
      hour: "numeric",
      minute: "2-digit",
    });
  }

  // Group entries by date
  function groupByDate(entries: TimelineEntry[]): Map<string, TimelineEntry[]> {
    const groups = new Map<string, TimelineEntry[]>();

    for (const entry of entries) {
      const dateKey = entry.date.toDateString();
      if (!groups.has(dateKey)) {
        groups.set(dateKey, []);
      }
      groups.get(dateKey)!.push(entry);
    }

    return groups;
  }

  function getDateLabel(dateStr: string): string {
    const date = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) return "Today";
    if (days === 1) return "Yesterday";
    if (days < 7) return date.toLocaleDateString("en-US", { weekday: "long" });

    return date.toLocaleDateString("en-US", {
      weekday: "long",
      month: "long",
      day: "numeric",
    });
  }

  const groupedEntries = $derived(groupByDate(entries));
</script>

{#if entries.length === 0}
  <div class="flex flex-col items-center justify-center py-16 text-center">
    <div class="rounded-full bg-muted p-4 mb-4">
      <MessageSquare class="h-8 w-8 text-muted-foreground" />
    </div>
    <h3 class="text-lg font-semibold mb-2">No entries yet</h3>
    <p class="text-sm text-muted-foreground max-w-sm">
      Start journaling by adding your first thought, note, or activity above.
    </p>
  </div>
{:else}
  <div class="space-y-8">
    {#each [...groupedEntries.entries()] as [dateStr, dayEntries]}
      <div>
        <!-- Date header -->
        <div class="flex items-center gap-3 mb-4">
          <div class="flex items-center gap-2 text-sm font-medium text-foreground">
            <Calendar class="h-4 w-4 text-muted-foreground" />
            {getDateLabel(dateStr)}
          </div>
          <div class="flex-1 h-px bg-border"></div>
          <span class="text-xs text-muted-foreground">
            {dayEntries.length} {dayEntries.length === 1 ? "entry" : "entries"}
          </span>
        </div>

        <!-- Entries for this date -->
        <div class="space-y-3 pl-2">
          {#each dayEntries as entry}
            <div class="group relative flex gap-4">
              <!-- Timeline line and dot -->
              <div class="flex flex-col items-center">
                <div
                  class={cn(
                    "flex h-8 w-8 shrink-0 items-center justify-center rounded-full border-2",
                    entry.type === "activity"
                      ? "border-primary bg-primary/10"
                      : entry.type === "friend_added"
                        ? "border-green-500 bg-green-500/10"
                        : entry.type === "location_added"
                          ? "border-blue-500 bg-blue-500/10"
                          : "border-muted-foreground/30 bg-muted"
                  )}
                >
                  {#if entry.type === "activity"}
                    <Activity class="h-4 w-4 text-primary" />
                  {:else if entry.type === "friend_added"}
                    <UserPlus class="h-4 w-4 text-green-500" />
                  {:else if entry.type === "location_added"}
                    <MapPinPlus class="h-4 w-4 text-blue-500" />
                  {:else}
                    <MessageSquare class="h-4 w-4 text-muted-foreground" />
                  {/if}
                </div>
                <div class="w-px flex-1 bg-border"></div>
              </div>

              <!-- Content -->
              <div class="flex-1 pb-6">
                <div class="rounded-lg border border-border bg-card p-4 shadow-sm transition-shadow hover:shadow-md">
                  <!-- Entry content -->
                  {#if entry.type === "friend_added" && entry.entityId}
                    <p class="text-sm text-foreground">
                      Added <button
                        type="button"
                        onclick={() => viewFriend(entry.entityId!)}
                        class="font-medium text-green-600 hover:text-green-700 dark:text-green-400 dark:hover:text-green-300 hover:underline"
                      >{entry.entityName}</button> as a friend
                    </p>
                  {:else if entry.type === "location_added" && entry.entityId}
                    <p class="text-sm text-foreground">
                      Added <button
                        type="button"
                        onclick={() => viewLocation(entry.entityId!)}
                        class="font-medium text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 hover:underline"
                      >{entry.entityName}</button> as a location
                    </p>
                  {:else}
                    <p class="text-sm text-foreground whitespace-pre-wrap">{entry.content}</p>
                  {/if}

                  <!-- Metadata -->
                  <div class="mt-3 flex flex-wrap items-center gap-x-3 gap-y-1.5 text-xs text-muted-foreground">
                    <span>{formatTime(entry.date)}</span>

                    {#if entry.location}
                      <button
                        type="button"
                        onclick={() => viewLocation(entry.location!)}
                        class="flex items-center gap-1 hover:text-foreground transition-colors"
                      >
                        <MapPin class="h-3 w-3" />
                        <span>@{entry.location}</span>
                      </button>
                    {/if}

                    {#if entry.tags && entry.tags.length > 0}
                      <div class="flex items-center gap-1.5">
                        <Hash class="h-3 w-3" />
                        <span>{entry.tags.map(t => `#${t}`).join(" ")}</span>
                      </div>
                    {/if}
                  </div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/each}
  </div>
{/if}
