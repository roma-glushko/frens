<script lang="ts">
  import { onMount } from "svelte";
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import { Calendar, Plus, Search, Hash, MapPin, Users, Loader2 } from "lucide-svelte";
  import { api, type Event } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";

  let activities = $state<Event[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      activities = await api.activities.list();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load activities";
    } finally {
      loading = false;
    }
  });

  let searchQuery = $state("");

  const filteredActivities = $derived(
    activities
      .filter(a =>
        a.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
        a.tags?.some(t => t.toLowerCase().includes(searchQuery.toLowerCase()))
      )
      .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
  );

  // Group activities by date
  const groupedActivities = $derived(() => {
    const groups: { label: string; activities: Event[] }[] = [];
    let currentLabel = "";

    for (const activity of filteredActivities) {
      const label = getDateLabel(activity.date);
      if (label !== currentLabel) {
        groups.push({ label, activities: [activity] });
        currentLabel = label;
      } else {
        groups[groups.length - 1].activities.push(activity);
      }
    }

    return groups;
  });

  function viewFriend(id: string) {
    currentPath.navigate(`/friends/${id}`);
  }

  function viewLocation(id: string) {
    currentPath.navigate(`/locations/${id}`);
  }

  function getDateLabel(dateStr: string): string {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "Unknown";

    const today = new Date();
    const activityDate = new Date(date.getFullYear(), date.getMonth(), date.getDate());
    const todayDate = new Date(today.getFullYear(), today.getMonth(), today.getDate());
    const days = Math.floor((todayDate.getTime() - activityDate.getTime()) / (1000 * 60 * 60 * 24));

    if (days === 0) return "Today";
    if (days === 1) return "Yesterday";
    return date.toLocaleDateString("en-US", {
      weekday: "long",
      month: "long",
      day: "numeric",
      year: date.getFullYear() !== today.getFullYear() ? "numeric" : undefined,
    });
  }
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Activities</h1>
      <p class="text-muted-foreground mt-1">
        {#if loading}
          Loading...
        {:else}
          {activities.length} {activities.length === 1 ? "activity" : "activities"}
        {/if}
      </p>
    </div>
    <Button>
      <Plus class="h-4 w-4 mr-2" />
      Log Activity
    </Button>
  </div>

  <!-- Loading State -->
  {#if loading}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading activities...</p>
        </div>
      </CardContent>
    </Card>
  {:else if error}
    <!-- Error State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-destructive/10 p-4 mb-4">
            <Calendar class="h-8 w-8 text-destructive" />
          </div>
          <h3 class="text-lg font-semibold mb-2">Failed to load activities</h3>
          <p class="text-sm text-muted-foreground max-w-sm">{error}</p>
        </div>
      </CardContent>
    </Card>
  {:else}
    <!-- Search -->
    <div class="mb-6">
      <div class="relative max-w-md">
        <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          type="text"
          bind:value={searchQuery}
          placeholder="Search activities..."
          class="w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>
    </div>

    <!-- Activities List -->
    {#if filteredActivities.length === 0}
      <Card>
        <CardContent class="py-16">
          <div class="flex flex-col items-center justify-center text-center">
            <div class="rounded-full bg-muted p-4 mb-4">
              <Calendar class="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold mb-2">
              {searchQuery ? "No activities found" : "No activities logged"}
            </h3>
            <p class="text-sm text-muted-foreground max-w-sm mb-4">
              {searchQuery
                ? "Try a different search term"
                : "Start tracking your social interactions. Log meetups, calls, and memorable moments."}
            </p>
            {#if !searchQuery}
              <Button>
                <Plus class="h-4 w-4 mr-2" />
                Log Your First Activity
              </Button>
            {/if}
          </div>
        </CardContent>
      </Card>
    {:else}
      <!-- Timeline grouped by date -->
      <div class="space-y-8">
        {#each groupedActivities() as group}
          <div>
            <!-- Date Header -->
            <h3 class="text-sm font-semibold text-foreground mb-4">{group.label}</h3>

            <!-- Activities for this date -->
            <div class="relative">
              <!-- Timeline line -->
              <div class="absolute left-[7px] top-2 bottom-2 w-px bg-border"></div>

              <div class="space-y-4">
                {#each group.activities as activity}
                  <div class="relative flex gap-4 pl-6">
                    <!-- Timeline dot -->
                    <div class="absolute left-0 top-1 h-[15px] w-[15px] rounded-full border-2 border-primary bg-background"></div>

                    <!-- Content -->
                    <div class="flex-1 min-w-0">
                      <!-- Activity content -->
                      <p class="text-sm">{activity.description}</p>

                      <!-- Tags -->
                      {#if activity.tags && activity.tags.length > 0}
                        <div class="flex items-center gap-1.5 mt-2 flex-wrap">
                          {#each activity.tags as tag}
                            <span class="inline-flex items-center gap-1 rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
                              <Hash class="h-2.5 w-2.5" />
                              {tag}
                            </span>
                          {/each}
                        </div>
                      {/if}

                      <!-- Friends & Locations -->
                      <div class="flex items-center gap-4 mt-2 text-xs text-muted-foreground">
                        {#if activity.friendIds && activity.friendIds.length > 0}
                          <div class="flex items-center gap-1">
                            <Users class="h-3 w-3" />
                            {#each activity.friendIds as friendId, i}
                              <button
                                type="button"
                                class="text-primary hover:underline"
                                onclick={() => viewFriend(friendId)}
                              >
                                {friendId}
                              </button>{#if i < activity.friendIds.length - 1}<span>,</span>{/if}
                            {/each}
                          </div>
                        {/if}
                        {#if activity.locationIds && activity.locationIds.length > 0}
                          <div class="flex items-center gap-1">
                            <MapPin class="h-3 w-3" />
                            {#each activity.locationIds as locationId, i}
                              <button
                                type="button"
                                class="text-primary hover:underline"
                                onclick={() => viewLocation(locationId)}
                              >
                                {locationId}
                              </button>{#if i < activity.locationIds.length - 1}<span>,</span>{/if}
                            {/each}
                          </div>
                        {/if}
                      </div>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  {/if}
</div>
