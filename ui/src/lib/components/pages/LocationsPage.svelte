<script lang="ts">
  import { onMount } from "svelte";
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import { MapPin, Plus, Search, Hash, Activity, Calendar, Loader2, LayoutGrid, Map } from "lucide-svelte";
  import { api, type Location } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";
  import LocationMap from "$lib/components/locations/LocationMap.svelte";

  type ViewMode = "cards" | "map";
  let viewMode = $state<ViewMode>("cards");

  let locations = $state<Location[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      locations = await api.locations.list();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load locations";
    } finally {
      loading = false;
    }
  });

  let searchQuery = $state("");

  const filteredLocations = $derived(
    locations.filter(l =>
      l.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.country?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      l.tags?.some(t => t.toLowerCase().includes(searchQuery.toLowerCase()))
    )
  );

  function formatLastActivity(dateStr?: string): string {
    if (!dateStr) return "No activity";
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "No activity";
    const days = Math.floor((Date.now() - date.getTime()) / (1000 * 60 * 60 * 24));
    if (days === 0) return "Today";
    if (days === 1) return "Yesterday";
    if (days < 7) return `${days} days ago`;
    if (days < 30) return `${Math.floor(days / 7)} weeks ago`;
    return `${Math.floor(days / 30)} months ago`;
  }

  function viewLocation(id: string) {
    currentPath.navigate(`/locations/${id}`);
  }
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Locations</h1>
      <p class="text-muted-foreground mt-1">
        {#if loading}
          Loading...
        {:else}
          {locations.length} saved {locations.length === 1 ? "place" : "places"}
        {/if}
      </p>
    </div>
    <div class="flex items-center gap-2">
      <!-- View Toggle -->
      <div class="flex items-center rounded-md border border-input bg-background p-1">
        <button
          type="button"
          class="inline-flex items-center justify-center rounded px-2.5 py-1.5 text-sm font-medium transition-colors {viewMode === 'cards' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}"
          onclick={() => viewMode = 'cards'}
        >
          <LayoutGrid class="h-4 w-4 mr-1.5" />
          Cards
        </button>
        <button
          type="button"
          class="inline-flex items-center justify-center rounded px-2.5 py-1.5 text-sm font-medium transition-colors {viewMode === 'map' ? 'bg-primary text-primary-foreground' : 'text-muted-foreground hover:text-foreground'}"
          onclick={() => viewMode = 'map'}
        >
          <Map class="h-4 w-4 mr-1.5" />
          Map
        </button>
      </div>
      <Button>
        <Plus class="h-4 w-4 mr-2" />
        Add Location
      </Button>
    </div>
  </div>

  <!-- Loading State -->
  {#if loading}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading locations...</p>
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
          <h3 class="text-lg font-semibold mb-2">Failed to load locations</h3>
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
          placeholder="Search by name, country, or tag..."
          class="w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>
    </div>

    <!-- Locations Content -->
    {#if viewMode === "map"}
      <!-- Map View -->
      <div class="relative h-[calc(100vh-280px)] min-h-[400px]">
        <LocationMap locations={filteredLocations} />
      </div>
    {:else}
      <!-- Cards View -->
      {#if filteredLocations.length === 0}
      <Card>
        <CardContent class="py-16">
          <div class="flex flex-col items-center justify-center text-center">
            <div class="rounded-full bg-muted p-4 mb-4">
              <MapPin class="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold mb-2">
              {searchQuery ? "No locations found" : "No locations yet"}
            </h3>
            <p class="text-sm text-muted-foreground max-w-sm mb-4">
              {searchQuery
                ? "Try a different search term"
                : "Save your favorite spots, restaurants, and venues."}
            </p>
            {#if !searchQuery}
              <Button>
                <Plus class="h-4 w-4 mr-2" />
                Add Your First Location
              </Button>
            {/if}
          </div>
        </CardContent>
      </Card>
    {:else}
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {#each filteredLocations as location}
          <button type="button" class="text-left w-full" onclick={() => viewLocation(location.id)}>
            <Card class="hover:shadow-md transition-shadow cursor-pointer h-full">
              <CardContent class="p-4">
              <div class="flex items-start gap-3">
                <!-- Icon -->
                <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                  <MapPin class="h-5 w-5 text-primary" />
                </div>

                <!-- Info -->
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <h3 class="font-semibold truncate">{location.name}</h3>
                    {#if location.country}
                      <span class="text-xs text-muted-foreground">
                        {location.country}
                      </span>
                    {/if}
                  </div>

                  {#if location.aliases && location.aliases.length > 0}
                    <p class="text-xs text-muted-foreground">
                      a.k.a. {location.aliases.join(", ")}
                    </p>
                  {/if}

                  {#if location.description}
                    <p class="text-sm text-muted-foreground truncate mt-1">
                      {location.description}
                    </p>
                  {/if}

                  <!-- Tags -->
                  {#if location.tags && location.tags.length > 0}
                    <div class="flex items-center gap-1.5 mt-2 flex-wrap">
                      {#each location.tags as tag}
                        <span class="inline-flex items-center gap-1 rounded-full bg-secondary px-2 py-0.5 text-xs text-secondary-foreground">
                          <Hash class="h-2.5 w-2.5" />
                          {tag}
                        </span>
                      {/each}
                    </div>
                  {/if}

                  <!-- Stats -->
                  <div class="flex items-center gap-4 mt-3 text-xs text-muted-foreground">
                    <div class="flex items-center gap-1">
                      <Activity class="h-3 w-3" />
                      {location.activitiesCount}
                    </div>
                    <div class="flex items-center gap-1">
                      <Calendar class="h-3 w-3" />
                      {formatLastActivity(location.lastActivity)}
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
            </Card>
          </button>
        {/each}
      </div>
      {/if}
    {/if}
  {/if}
</div>
