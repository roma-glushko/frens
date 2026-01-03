<script lang="ts">
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import CardHeader from "$lib/components/ui/card/CardHeader.svelte";
  import CardTitle from "$lib/components/ui/card/CardTitle.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import {
    ArrowLeft,
    Hash,
    MapPin,
    Calendar,
    Activity,
    FileText,
    Loader2,
    Globe,
    Users,
  } from "lucide-svelte";
  import { currentPath } from "$lib/stores/router.svelte";
  import { api, type Location, type Event } from "$lib/api";

  // Extract location ID from path
  const locationId = $derived(() => {
    const parts = currentPath.value.split("/").filter((p) => p);
    return parts[1]; // /locations/:id -> parts = ["locations", "id"]
  });

  let location = $state<Location | null>(null);
  let activities = $state<Event[]>([]);
  let notes = $state<Event[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  $effect(() => {
    const id = locationId();
    if (id) {
      loadLocation(id);
    }
  });

  async function loadLocation(id: string) {
    loading = true;
    error = null;
    try {
      const [locationData, activitiesData, notesData] = await Promise.all([
        api.locations.get(id),
        api.locations.activities(id),
        api.locations.notes(id),
      ]);
      location = locationData;
      activities = activitiesData;
      notes = notesData;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load location";
    } finally {
      loading = false;
    }
  }

  function goBack() {
    currentPath.navigate("/locations");
  }

  function formatLastActivity(dateStr?: string): string {
    if (!dateStr) return "No activity yet";
    const date = new Date(dateStr);
    if (isNaN(date.getTime()) || date.getFullYear() < 1970) return "No activity yet";
    const days = Math.floor(
      (Date.now() - date.getTime()) / (1000 * 60 * 60 * 24)
    );
    if (days === 0) return "Today";
    if (days === 1) return "Yesterday";
    if (days < 7) return `${days} days ago`;
    if (days < 30) return `${Math.floor(days / 7)} weeks ago`;
    return `${Math.floor(days / 30)} months ago`;
  }

  function formatEventDate(dateStr: string): string {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "";
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  }

  function formatRelativeDate(dateStr: string): string {
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "";
    const days = Math.floor(
      (Date.now() - date.getTime()) / (1000 * 60 * 60 * 24)
    );
    if (days === 0) return "Today";
    if (days === 1) return "Yesterday";
    if (days < 7) return `${days} days ago`;
    return formatEventDate(dateStr);
  }
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Back Button -->
  <div class="mb-6">
    <Button variant="ghost" onclick={goBack} class="-ml-2">
      <ArrowLeft class="h-4 w-4 mr-2" />
      Back to Locations
    </Button>
  </div>

  {#if loading}
    <!-- Loading State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading location...</p>
        </div>
      </CardContent>
    </Card>
  {:else if error}
    <!-- Error State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-destructive/10 p-4 mb-4">
            <MapPin class="h-8 w-8 text-destructive" />
          </div>
          <h3 class="text-lg font-semibold mb-2">Location not found</h3>
          <p class="text-sm text-muted-foreground max-w-sm mb-4">{error}</p>
          <Button onclick={goBack}>Go Back</Button>
        </div>
      </CardContent>
    </Card>
  {:else if location}
    <!-- Profile Header -->
    <div class="flex flex-col md:flex-row gap-6 mb-8">
      <!-- Icon -->
      <div
        class="flex h-24 w-24 shrink-0 items-center justify-center rounded-2xl bg-primary/10"
      >
        <MapPin class="h-12 w-12 text-primary" />
      </div>

      <!-- Basic Info -->
      <div class="flex-1">
        <div class="flex flex-col sm:flex-row sm:items-center gap-2 mb-2">
          <h1 class="text-3xl font-bold tracking-tight">{location.name}</h1>
          {#if location.country}
            <span class="inline-flex items-center gap-1 text-lg text-muted-foreground">
              <Globe class="h-4 w-4" />
              {location.country}
            </span>
          {/if}
        </div>

        {#if location.aliases && location.aliases.length > 0}
          <p class="text-muted-foreground mb-2">
            a.k.a. {location.aliases.join(", ")}
          </p>
        {/if}

        {#if location.description}
          <p class="text-muted-foreground mb-4">{location.description}</p>
        {/if}

        <!-- Tags -->
        {#if location.tags && location.tags.length > 0}
          <div class="flex items-center gap-2 flex-wrap">
            {#each location.tags as tag}
              <span
                class="inline-flex items-center gap-1 rounded-full bg-secondary px-3 py-1 text-sm text-secondary-foreground"
              >
                <Hash class="h-3 w-3" />
                {tag}
              </span>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <!-- Stats Row -->
    <div class="grid grid-cols-2 md:grid-cols-3 gap-4 mb-8">
      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="rounded-full bg-primary/10 p-2">
            <Activity class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{location.activitiesCount}</p>
            <p class="text-xs text-muted-foreground">Activities</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="rounded-full bg-primary/10 p-2">
            <FileText class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{location.notesCount}</p>
            <p class="text-xs text-muted-foreground">Notes</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="rounded-full bg-primary/10 p-2">
            <Calendar class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-sm font-medium">
              {formatLastActivity(location.lastActivity)}
            </p>
            <p class="text-xs text-muted-foreground">Last Activity</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Activities Section -->
    {#if activities.length > 0}
      <Card class="mb-6">
        <CardHeader>
          <CardTitle class="text-lg flex items-center gap-2">
            <Activity class="h-4 w-4" />
            Recent Activities
          </CardTitle>
        </CardHeader>
        <CardContent class="pt-0">
          <div class="space-y-4">
            {#each activities.slice(0, 10) as activity}
              <div class="flex gap-4 pb-4 border-b border-border last:border-0 last:pb-0">
                <div class="flex-shrink-0 w-20 text-xs text-muted-foreground">
                  {formatRelativeDate(activity.date)}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm">{activity.description}</p>
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
                  {#if activity.friendIds && activity.friendIds.length > 0}
                    <div class="flex items-center gap-1.5 mt-1.5 text-xs text-muted-foreground">
                      <Users class="h-3 w-3" />
                      {activity.friendIds.join(", ")}
                    </div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
          {#if activities.length > 10}
            <p class="text-xs text-muted-foreground mt-4 text-center">
              Showing 10 of {activities.length} activities
            </p>
          {/if}
        </CardContent>
      </Card>
    {/if}

    <!-- Notes Section -->
    {#if notes.length > 0}
      <Card class="mb-6">
        <CardHeader>
          <CardTitle class="text-lg flex items-center gap-2">
            <FileText class="h-4 w-4" />
            Notes
          </CardTitle>
        </CardHeader>
        <CardContent class="pt-0">
          <div class="space-y-4">
            {#each notes.slice(0, 10) as note}
              <div class="flex gap-4 pb-4 border-b border-border last:border-0 last:pb-0">
                <div class="flex-shrink-0 w-20 text-xs text-muted-foreground">
                  {formatRelativeDate(note.date)}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm">{note.description}</p>
                  {#if note.tags && note.tags.length > 0}
                    <div class="flex items-center gap-1.5 mt-2 flex-wrap">
                      {#each note.tags as tag}
                        <span class="inline-flex items-center gap-1 rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
                          <Hash class="h-2.5 w-2.5" />
                          {tag}
                        </span>
                      {/each}
                    </div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
          {#if notes.length > 10}
            <p class="text-xs text-muted-foreground mt-4 text-center">
              Showing 10 of {notes.length} notes
            </p>
          {/if}
        </CardContent>
      </Card>
    {/if}

    <!-- Empty State for no activities/notes -->
    {#if activities.length === 0 && notes.length === 0}
      <Card>
        <CardContent class="py-12">
          <div class="flex flex-col items-center justify-center text-center">
            <div class="rounded-full bg-muted p-4 mb-4">
              <MapPin class="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold mb-2">No activities yet</h3>
            <p class="text-sm text-muted-foreground max-w-sm">
              Log activities or notes at {location.name} to see them here.
            </p>
          </div>
        </CardContent>
      </Card>
    {/if}
  {/if}
</div>
