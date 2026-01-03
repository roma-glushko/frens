<script lang="ts">
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardHeader from "$lib/components/ui/card/CardHeader.svelte";
  import CardTitle from "$lib/components/ui/card/CardTitle.svelte";
  import CardDescription from "$lib/components/ui/card/CardDescription.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import { Users, Calendar, MapPin, StickyNote, TrendingUp, Clock } from "lucide-svelte";
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="mb-8">
    <h1 class="text-3xl font-bold tracking-tight">Dashboard</h1>
    <p class="text-muted-foreground mt-2">
      Overview of your friendships and activities
    </p>
  </div>

  <!-- Stats Grid -->
  <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mb-8">
    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Total Friends</CardTitle>
        <Users class="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">--</div>
        <p class="text-xs text-muted-foreground">Connect your journal</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Activities</CardTitle>
        <Calendar class="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">--</div>
        <p class="text-xs text-muted-foreground">This month</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Notes</CardTitle>
        <StickyNote class="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">--</div>
        <p class="text-xs text-muted-foreground">Total notes</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle class="text-sm font-medium">Locations</CardTitle>
        <MapPin class="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div class="text-2xl font-bold">--</div>
        <p class="text-xs text-muted-foreground">Saved places</p>
      </CardContent>
    </Card>
  </div>

  <!-- Content Grid -->
  <div class="grid gap-6 md:grid-cols-2">
    <!-- Recent Activity -->
    <Card>
      <CardHeader>
        <div class="flex items-center gap-2">
          <Clock class="h-5 w-5 text-muted-foreground" />
          <CardTitle>Recent Activity</CardTitle>
        </div>
        <CardDescription>Your latest interactions</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-col items-center justify-center py-8 text-center">
          <Calendar class="h-12 w-12 text-muted-foreground/50 mb-4" />
          <p class="text-sm text-muted-foreground">No recent activities</p>
          <p class="text-xs text-muted-foreground mt-1">
            Activities will appear here once you start logging
          </p>
        </div>
      </CardContent>
    </Card>

    <!-- Friends Overview -->
    <Card>
      <CardHeader>
        <div class="flex items-center gap-2">
          <TrendingUp class="h-5 w-5 text-muted-foreground" />
          <CardTitle>Connection Insights</CardTitle>
        </div>
        <CardDescription>How you're staying connected</CardDescription>
      </CardHeader>
      <CardContent>
        <div class="flex flex-col items-center justify-center py-8 text-center">
          <Users class="h-12 w-12 text-muted-foreground/50 mb-4" />
          <p class="text-sm text-muted-foreground">No data available</p>
          <p class="text-xs text-muted-foreground mt-1">
            Add friends and activities to see insights
          </p>
        </div>
      </CardContent>
    </Card>
  </div>
</div>
