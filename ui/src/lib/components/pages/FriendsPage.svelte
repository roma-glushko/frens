<script lang="ts">
  import { onMount } from "svelte";
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import { Users, Plus, Search, MapPin, Calendar, Activity, Hash, Loader2 } from "lucide-svelte";
  import { api, type Friend } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";

  let friends = $state<Friend[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      friends = await api.friends.list();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load friends";
    } finally {
      loading = false;
    }
  });

  let searchQuery = $state("");

  const filteredFriends = $derived(
    friends.filter(f =>
      f.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      f.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      f.tags?.some(t => t.toLowerCase().includes(searchQuery.toLowerCase()))
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

  function viewFriend(id: string) {
    currentPath.navigate(`/friends/${id}`);
  }
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Friends</h1>
      <p class="text-muted-foreground mt-1">
        {#if loading}
          Loading...
        {:else}
          {friends.length} {friends.length === 1 ? "friend" : "friends"} in your circle
        {/if}
      </p>
    </div>
    <Button>
      <Plus class="h-4 w-4 mr-2" />
      Add Friend
    </Button>
  </div>

  <!-- Loading State -->
  {#if loading}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading friends...</p>
        </div>
      </CardContent>
    </Card>
  {:else if error}
    <!-- Error State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-destructive/10 p-4 mb-4">
            <Users class="h-8 w-8 text-destructive" />
          </div>
          <h3 class="text-lg font-semibold mb-2">Failed to load friends</h3>
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
          placeholder="Search by name, description, or tag..."
          class="w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </div>
    </div>

    <!-- Friends List -->
    {#if filteredFriends.length === 0}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-muted p-4 mb-4">
            <Users class="h-8 w-8 text-muted-foreground" />
          </div>
          <h3 class="text-lg font-semibold mb-2">
            {searchQuery ? "No friends found" : "No friends yet"}
          </h3>
          <p class="text-sm text-muted-foreground max-w-sm mb-4">
            {searchQuery
              ? "Try a different search term"
              : "Start building your network by adding your first friend."}
          </p>
          {#if !searchQuery}
            <Button>
              <Plus class="h-4 w-4 mr-2" />
              Add Your First Friend
            </Button>
          {/if}
        </div>
      </CardContent>
    </Card>
  {:else}
    <div class="grid gap-4 md:grid-cols-2">
      {#each filteredFriends as friend}
        <button
          type="button"
          class="text-left w-full"
          onclick={() => viewFriend(friend.id)}
        >
          <Card class="hover:shadow-md transition-shadow cursor-pointer h-full">
            <CardContent class="p-4">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary font-semibold text-lg">
                {friend.name.charAt(0)}
              </div>

              <!-- Info -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="font-semibold truncate">{friend.name}</h3>
                  {#if friend.nicknames && friend.nicknames.length > 0}
                    <span class="text-xs text-muted-foreground">
                      ({friend.nicknames.join(", ")})
                    </span>
                  {/if}
                </div>

                {#if friend.description}
                  <p class="text-sm text-muted-foreground truncate mt-0.5">
                    {friend.description}
                  </p>
                {/if}

                <!-- Tags -->
                {#if friend.tags && friend.tags.length > 0}
                  <div class="flex items-center gap-1.5 mt-2 flex-wrap">
                    {#each friend.tags as tag}
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
                    {friend.activitiesCount} activities
                  </div>
                  <div class="flex items-center gap-1">
                    <Calendar class="h-3 w-3" />
                    {formatLastActivity(friend.lastActivity)}
                  </div>
                  {#if friend.locations && friend.locations.length > 0}
                    <div class="flex items-center gap-1">
                      <MapPin class="h-3 w-3" />
                      {friend.locations.length}
                    </div>
                  {/if}
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
</div>
