<script lang="ts">
  import { Layout } from "$lib/components/layout";
  import {
    DashboardPage,
    FriendsPage,
    FriendProfilePage,
    ActivitiesPage,
    NotesPage,
    LocationsPage,
    LocationProfilePage,
    StatsPage,
  } from "$lib/components/pages";
  import { currentPath } from "$lib/stores/router.svelte";
  import { theme } from "$lib/stores/theme.svelte";

  // Initialize theme on mount
  $effect(() => {
    theme.applyTheme();
  });

  // Check if we're on a friend profile page (/friends/:id)
  const isFriendProfile = $derived(() => {
    const path = currentPath.value;
    if (!path.startsWith("/friends")) return false;
    const parts = path.split("/").filter((p) => p);
    return parts.length > 1; // ["friends", "id"]
  });

  // Check if we're on a location profile page (/locations/:id)
  const isLocationProfile = $derived(() => {
    const path = currentPath.value;
    if (!path.startsWith("/locations")) return false;
    const parts = path.split("/").filter((p) => p);
    return parts.length > 1; // ["locations", "id"]
  });
</script>

<Layout>
  {#if currentPath.value === "/" || currentPath.value === ""}
    <DashboardPage />
  {:else if currentPath.value.startsWith("/friends")}
    {#if isFriendProfile()}
      <FriendProfilePage />
    {:else}
      <FriendsPage />
    {/if}
  {:else if currentPath.value.startsWith("/activities")}
    <ActivitiesPage />
  {:else if currentPath.value.startsWith("/notes")}
    <NotesPage />
  {:else if currentPath.value.startsWith("/locations")}
    {#if isLocationProfile()}
      <LocationProfilePage />
    {:else}
      <LocationsPage />
    {/if}
  {:else if currentPath.value.startsWith("/stats")}
    <StatsPage />
  {:else}
    <DashboardPage />
  {/if}
</Layout>
