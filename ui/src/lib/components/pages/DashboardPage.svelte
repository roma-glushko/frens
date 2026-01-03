<script lang="ts">
  import { onMount } from "svelte";
  import { JournalInput, Timeline, type TimelineEntry } from "$lib/components/journal";
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardHeader from "$lib/components/ui/card/CardHeader.svelte";
  import CardTitle from "$lib/components/ui/card/CardTitle.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import { Users, Calendar, StickyNote, MapPin } from "lucide-svelte";
  import { api, type Stats, type Event } from "$lib/api";

  let stats = $state<Stats>({
    friends: 0,
    activities: 0,
    notes: 0,
    locations: 0,
  });

  let recentEntries = $state<TimelineEntry[]>([]);

  function eventToTimelineEntry(event: Event): TimelineEntry {
    return {
      id: event.id,
      type: event.type,
      content: event.description,
      date: new Date(event.date),
      tags: event.tags,
      location: event.locationIds?.[0],
    };
  }

  onMount(async () => {
    try {
      const [statsData, activities, notes] = await Promise.all([
        api.stats.get(),
        api.activities.list(),
        api.notes.list(),
      ]);

      stats = statsData;

      // Combine and sort by date, take most recent 10
      const allEvents = [...activities, ...notes];
      allEvents.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
      recentEntries = allEvents.slice(0, 10).map(eventToTimelineEntry);
    } catch (e) {
      console.error("Failed to load dashboard data:", e);
    }
  });
</script>

<div class="container mx-auto px-4 py-6">
  <!-- Quick Stats Bar -->
  <div class="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
    <div class="flex items-center gap-3 rounded-lg border border-border bg-card p-3">
      <div class="flex h-9 w-9 items-center justify-center rounded-md bg-primary/10">
        <Users class="h-4 w-4 text-primary" />
      </div>
      <div>
        <p class="text-2xl font-bold">{stats.friends}</p>
        <p class="text-xs text-muted-foreground">Friends</p>
      </div>
    </div>

    <div class="flex items-center gap-3 rounded-lg border border-border bg-card p-3">
      <div class="flex h-9 w-9 items-center justify-center rounded-md bg-primary/10">
        <Calendar class="h-4 w-4 text-primary" />
      </div>
      <div>
        <p class="text-2xl font-bold">{stats.activities}</p>
        <p class="text-xs text-muted-foreground">Activities</p>
      </div>
    </div>

    <div class="flex items-center gap-3 rounded-lg border border-border bg-card p-3">
      <div class="flex h-9 w-9 items-center justify-center rounded-md bg-primary/10">
        <StickyNote class="h-4 w-4 text-primary" />
      </div>
      <div>
        <p class="text-2xl font-bold">{stats.notes}</p>
        <p class="text-xs text-muted-foreground">Notes</p>
      </div>
    </div>

    <div class="flex items-center gap-3 rounded-lg border border-border bg-card p-3">
      <div class="flex h-9 w-9 items-center justify-center rounded-md bg-primary/10">
        <MapPin class="h-4 w-4 text-primary" />
      </div>
      <div>
        <p class="text-2xl font-bold">{stats.locations}</p>
        <p class="text-xs text-muted-foreground">Locations</p>
      </div>
    </div>
  </div>

  <!-- Main content area -->
  <div class="grid gap-6 lg:grid-cols-3">
    <!-- Journal section (takes 2 columns on large screens) -->
    <div class="lg:col-span-2 space-y-6">
      <!-- Quick input -->
      <div>
        <h2 class="text-lg font-semibold mb-3">Quick Journal</h2>
        <JournalInput />
      </div>

      <!-- Timeline -->
      <div>
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold">Recent Activity</h2>
          <a href="/activities" class="text-sm text-muted-foreground hover:text-foreground transition-colors">
            View all
          </a>
        </div>
        <Timeline entries={recentEntries} />
      </div>
    </div>

    <!-- Sidebar -->
    <div class="space-y-6">
      <!-- Tips card -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">Frentxt Syntax</CardTitle>
        </CardHeader>
        <CardContent class="space-y-2.5 text-sm text-muted-foreground">
          <p>
            <strong class="text-foreground">Dates:</strong> <code class="rounded bg-muted px-1">yesterday ::</code> or <code class="rounded bg-muted px-1">2 days ago ::</code>
          </p>
          <p>
            <strong class="text-foreground">Tags:</strong> <code class="rounded bg-muted px-1">#work</code> <code class="rounded bg-muted px-1">#family</code>
          </p>
          <p>
            <strong class="text-foreground">Locations:</strong> <code class="rounded bg-muted px-1">@home</code> <code class="rounded bg-muted px-1">@office</code>
          </p>
          <div class="pt-2 border-t border-border mt-2">
            <p class="text-xs">
              <span class="text-foreground">Example:</span><br/>
              <code class="text-[11px]">yesterday :: Coffee with Sarah #catchup @cafe</code>
            </p>
          </div>
        </CardContent>
      </Card>

      <!-- Upcoming dates -->
      <Card>
        <CardHeader>
          <CardTitle class="text-base">Upcoming</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="space-y-3 text-sm">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-primary"></div>
                <span>Mom's Birthday</span>
              </div>
              <span class="text-muted-foreground">in 3 weeks</span>
            </div>
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2">
                <div class="h-2 w-2 rounded-full bg-muted-foreground"></div>
                <span>Sarah's Anniversary</span>
              </div>
              <span class="text-muted-foreground">in 2 months</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- CLI hint -->
      <Card class="bg-muted/50">
        <CardContent class="pt-6">
          <p class="text-sm text-muted-foreground mb-2">
            Prefer the terminal?
          </p>
          <code class="block rounded bg-background border border-border p-2 text-xs">
            frens note add "your thought"
          </code>
        </CardContent>
      </Card>
    </div>
  </div>
</div>
