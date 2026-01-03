<script lang="ts">
  import { onMount } from "svelte";
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardHeader from "$lib/components/ui/card/CardHeader.svelte";
  import CardTitle from "$lib/components/ui/card/CardTitle.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import { Users, Calendar, StickyNote, MapPin, TrendingUp, Loader2 } from "lucide-svelte";
  import { api, type ComprehensiveStats } from "$lib/api";
  import { ActivityChart, TopList, InsightCard } from "$lib/components/stats";

  let stats = $state<ComprehensiveStats | null>(null);
  let loading = $state(true);
  let error = $state<string | null>(null);

  onMount(async () => {
    try {
      stats = await api.stats.getComprehensive();
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load stats";
    } finally {
      loading = false;
    }
  });
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="mb-8">
    <h1 class="text-3xl font-bold tracking-tight">Statistics</h1>
    <p class="text-muted-foreground mt-1">
      Your friendship insights and trends
    </p>
  </div>

  {#if loading}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading statistics...</p>
        </div>
      </CardContent>
    </Card>
  {:else if error}
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-destructive/10 p-4 mb-4">
            <TrendingUp class="h-8 w-8 text-destructive" />
          </div>
          <h3 class="text-lg font-semibold mb-2">Failed to load statistics</h3>
          <p class="text-sm text-muted-foreground max-w-sm">{error}</p>
        </div>
      </CardContent>
    </Card>
  {:else if stats}
    <!-- Numbers Section -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
            <Users class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{stats.counts.friends}</p>
            <p class="text-xs text-muted-foreground">Friends</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
            <Calendar class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{stats.counts.activities}</p>
            <p class="text-xs text-muted-foreground">Activities</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
            <StickyNote class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{stats.counts.notes}</p>
            <p class="text-xs text-muted-foreground">Notes</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
            <MapPin class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{stats.counts.locations}</p>
            <p class="text-xs text-muted-foreground">Locations</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Activity Chart -->
    <Card class="mb-8">
      <CardHeader>
        <CardTitle class="text-lg">Activity Over Time</CardTitle>
      </CardHeader>
      <CardContent>
        <ActivityChart data={stats.activityTimeline} />
      </CardContent>
    </Card>

    <!-- Tops Section -->
    <div class="grid md:grid-cols-3 gap-6 mb-8">
      <Card>
        <CardContent class="p-4">
          <TopList
            title="Top Friends"
            items={stats.topFriends}
            type="friends"
            emptyMessage="No friends with activities yet"
          />
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4">
          <TopList
            title="Top Locations"
            items={stats.topLocations}
            type="locations"
            emptyMessage="No locations with activities yet"
          />
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4">
          <TopList
            title="Top Tags"
            items={stats.topTags}
            type="tags"
            emptyMessage="No tags used yet"
          />
        </CardContent>
      </Card>
    </div>

    <!-- Insights Section -->
    <Card>
      <CardContent class="p-4">
        <InsightCard insights={stats.insights} />
      </CardContent>
    </Card>
  {/if}
</div>
