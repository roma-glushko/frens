<script lang="ts">
  import { Layout } from "$lib/components/layout";
  import {
    DashboardPage,
    FriendsPage,
    ActivitiesPage,
    NotesPage,
    LocationsPage,
  } from "$lib/components/pages";
  import { currentPath } from "$lib/stores/router.svelte";
  import { theme } from "$lib/stores/theme.svelte";

  // Initialize theme on mount
  $effect(() => {
    theme.applyTheme();
  });
</script>

<Layout>
  {#if currentPath.value === "/" || currentPath.value === ""}
    <DashboardPage />
  {:else if currentPath.value.startsWith("/friends")}
    <FriendsPage />
  {:else if currentPath.value.startsWith("/activities")}
    <ActivitiesPage />
  {:else if currentPath.value.startsWith("/notes")}
    <NotesPage />
  {:else if currentPath.value.startsWith("/locations")}
    <LocationsPage />
  {:else}
    <DashboardPage />
  {/if}
</Layout>
