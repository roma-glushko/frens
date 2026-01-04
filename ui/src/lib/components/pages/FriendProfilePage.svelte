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
    Mail,
    Phone,
    MessageCircle,
    Loader2,
    User,
    ExternalLink,
  } from "lucide-svelte";
  import { currentPath } from "$lib/stores/router.svelte";
  import { api, type Friend, type Event } from "$lib/api";

  // Extract friend ID from path
  const friendId = $derived(() => {
    const parts = currentPath.value.split("/").filter((p) => p);
    return parts[1]; // /friends/:id -> parts = ["friends", "id"]
  });

  let friend = $state<Friend | null>(null);
  let activities = $state<Event[]>([]);
  let notes = $state<Event[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  $effect(() => {
    const id = friendId();
    if (id) {
      loadFriend(id);
    }
  });

  async function loadFriend(id: string) {
    loading = true;
    error = null;
    try {
      const [friendData, activitiesData, notesData] = await Promise.all([
        api.friends.get(id),
        api.friends.activities(id),
        api.friends.notes(id),
      ]);
      friend = friendData;
      activities = activitiesData;
      notes = notesData;
    } catch (e) {
      error = e instanceof Error ? e.message : "Failed to load friend";
    } finally {
      loading = false;
    }
  }

  function goBack() {
    currentPath.navigate("/friends");
  }

  function viewLocation(id: string) {
    currentPath.navigate(`/locations/${id}`);
  }

  function formatLastActivity(dateStr?: string): string {
    if (!dateStr) return "No activity yet";
    const date = new Date(dateStr);
    if (isNaN(date.getTime())) return "No activity yet";
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

  function getContactIcon(type: string) {
    switch (type.toLowerCase()) {
      case "email":
        return Mail;
      case "phone":
        return Phone;
      case "telegram":
      case "whatsapp":
      case "signal":
        return MessageCircle;
      default:
        return ExternalLink;
    }
  }

  function getContactHref(type: string, value: string): string | null {
    switch (type.toLowerCase()) {
      case "email":
        return `mailto:${value}`;
      case "phone":
        return `tel:${value}`;
      case "telegram":
        return `https://t.me/${value.replace("@", "")}`;
      case "whatsapp":
        return `https://wa.me/${value.replace(/[^0-9]/g, "")}`;
      default:
        return null;
    }
  }
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Back Button -->
  <div class="mb-6">
    <Button variant="ghost" onclick={goBack} class="-ml-2">
      <ArrowLeft class="h-4 w-4 mr-2" />
      Back to Friends
    </Button>
  </div>

  {#if loading}
    <!-- Loading State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <Loader2 class="h-8 w-8 animate-spin text-muted-foreground mb-4" />
          <p class="text-sm text-muted-foreground">Loading friend profile...</p>
        </div>
      </CardContent>
    </Card>
  {:else if error}
    <!-- Error State -->
    <Card>
      <CardContent class="py-16">
        <div class="flex flex-col items-center justify-center text-center">
          <div class="rounded-full bg-destructive/10 p-4 mb-4">
            <User class="h-8 w-8 text-destructive" />
          </div>
          <h3 class="text-lg font-semibold mb-2">Friend not found</h3>
          <p class="text-sm text-muted-foreground max-w-sm mb-4">{error}</p>
          <Button onclick={goBack}>Go Back</Button>
        </div>
      </CardContent>
    </Card>
  {:else if friend}
    <!-- Profile Header -->
    <div class="flex flex-col md:flex-row gap-6 mb-8">
      <!-- Avatar -->
      <div
        class="flex h-24 w-24 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary font-bold text-4xl"
      >
        {friend.name.charAt(0)}
      </div>

      <!-- Basic Info -->
      <div class="flex-1">
        <div class="flex flex-col sm:flex-row sm:items-center gap-2 mb-2">
          <h1 class="text-3xl font-bold tracking-tight">{friend.name}</h1>
          {#if friend.nicknames && friend.nicknames.length > 0}
            <span class="text-lg text-muted-foreground">
              (a.k.a. {friend.nicknames.join(", ")})
            </span>
          {/if}
        </div>

        {#if friend.description}
          <p class="text-muted-foreground mb-4">{friend.description}</p>
        {/if}

        <!-- Tags -->
        {#if friend.tags && friend.tags.length > 0}
          <div class="flex items-center gap-2 flex-wrap">
            {#each friend.tags as tag}
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
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="rounded-full bg-primary/10 p-2">
            <Activity class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{friend.activitiesCount}</p>
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
            <p class="text-2xl font-bold">{friend.notesCount}</p>
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
              {formatLastActivity(friend.lastActivity)}
            </p>
            <p class="text-xs text-muted-foreground">Last Activity</p>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent class="p-4 flex items-center gap-3">
          <div class="rounded-full bg-primary/10 p-2">
            <MapPin class="h-5 w-5 text-primary" />
          </div>
          <div>
            <p class="text-2xl font-bold">{friend.locations?.length ?? 0}</p>
            <p class="text-xs text-muted-foreground">Locations</p>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Details Grid -->
    <div class="grid md:grid-cols-2 gap-6">
      <!-- Contacts -->
      {#if friend.contacts && friend.contacts.length > 0}
        <Card>
          <CardHeader>
            <CardTitle class="text-lg flex items-center gap-2">
              <Mail class="h-4 w-4" />
              Contacts
            </CardTitle>
          </CardHeader>
          <CardContent class="pt-0">
            <div class="space-y-3">
              {#each friend.contacts as contact}
                {@const Icon = getContactIcon(contact.type)}
                {@const href = getContactHref(contact.type, contact.value)}
                <div class="flex items-center gap-3">
                  <div class="rounded-full bg-muted p-2">
                    <Icon class="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-xs text-muted-foreground capitalize">
                      {contact.type}
                    </p>
                    {#if href}
                      <a
                        {href}
                        class="text-sm text-primary hover:underline truncate block"
                      >
                        {contact.value}
                      </a>
                    {:else}
                      <p class="text-sm truncate">{contact.value}</p>
                    {/if}
                  </div>
                </div>
              {/each}
            </div>
          </CardContent>
        </Card>
      {/if}

      <!-- Locations -->
      {#if friend.locations && friend.locations.length > 0}
        <Card>
          <CardHeader>
            <CardTitle class="text-lg flex items-center gap-2">
              <MapPin class="h-4 w-4" />
              Associated Locations
            </CardTitle>
          </CardHeader>
          <CardContent class="pt-0">
            <div class="flex flex-wrap gap-2">
              {#each friend.locations as location}
                <button
                  type="button"
                  onclick={() => viewLocation(location)}
                  class="inline-flex items-center gap-1.5 rounded-md bg-muted px-3 py-1.5 text-sm hover:bg-muted/80 transition-colors"
                >
                  <MapPin class="h-3.5 w-3.5 text-muted-foreground" />
                  {location}
                </button>
              {/each}
            </div>
          </CardContent>
        </Card>
      {/if}
    </div>

    <!-- Activities Section -->
    {#if activities.length > 0}
      <Card class="mt-6">
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
                  {#if activity.locationIds && activity.locationIds.length > 0}
                    <div class="flex items-center gap-1.5 mt-1.5 text-xs text-muted-foreground flex-wrap">
                      <MapPin class="h-3 w-3" />
                      {#each activity.locationIds as locationId, i}
                        <button
                          type="button"
                          onclick={() => viewLocation(locationId)}
                          class="hover:text-foreground transition-colors"
                        >
                          {locationId}{i < activity.locationIds.length - 1 ? "," : ""}
                        </button>
                      {/each}
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
      <Card class="mt-6">
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

    <!-- Empty State for no additional info -->
    {#if (!friend.contacts || friend.contacts.length === 0) && (!friend.locations || friend.locations.length === 0) && activities.length === 0 && notes.length === 0}
      <Card class="mt-6">
        <CardContent class="py-12">
          <div class="flex flex-col items-center justify-center text-center">
            <div class="rounded-full bg-muted p-4 mb-4">
              <User class="h-8 w-8 text-muted-foreground" />
            </div>
            <h3 class="text-lg font-semibold mb-2">No additional details</h3>
            <p class="text-sm text-muted-foreground max-w-sm">
              Add contacts, locations, activities, or notes to build a richer profile for
              {friend.name}.
            </p>
          </div>
        </CardContent>
      </Card>
    {/if}
  {/if}
</div>
