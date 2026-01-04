<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import type { Location } from "$lib/api";
  import { currentPath } from "$lib/stores/router.svelte";
  import L from "leaflet";
  import "leaflet/dist/leaflet.css";

  // Custom marker icon to avoid Leaflet's bundler issues
  const customIcon = L.divIcon({
    className: "custom-marker",
    html: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" width="32" height="32">
      <path fill-rule="evenodd" d="M11.54 22.351l.07.04.028.016a.76.76 0 00.723 0l.028-.015.071-.041a16.975 16.975 0 001.144-.742 19.58 19.58 0 002.683-2.282c1.944-1.99 3.963-4.98 3.963-8.827a8.25 8.25 0 00-16.5 0c0 3.846 2.02 6.837 3.963 8.827a19.58 19.58 0 002.682 2.282 16.975 16.975 0 001.145.742zM12 13.5a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd" />
    </svg>`,
    iconSize: [32, 32],
    iconAnchor: [16, 32],
    popupAnchor: [0, -32],
  });

  interface Props {
    locations: Location[];
  }

  let { locations }: Props = $props();

  let mapContainer: HTMLDivElement;
  let map: L.Map | null = null;
  let markers: L.Marker[] = [];

  // Filter locations that have coordinates
  const locationsWithCoords = $derived(
    locations.filter((l) => l.lat !== undefined && l.lng !== undefined)
  );

  function viewLocation(id: string) {
    currentPath.navigate(`/locations/${id}`);
  }

  function initMap() {
    if (!mapContainer || map) return;

    // Default center (world view)
    let center: L.LatLngExpression = [20, 0];
    let zoom = 2;

    // If we have locations with coords, fit bounds to them
    if (locationsWithCoords.length > 0) {
      const lats = locationsWithCoords.map((l) => l.lat!);
      const lngs = locationsWithCoords.map((l) => l.lng!);
      center = [
        (Math.min(...lats) + Math.max(...lats)) / 2,
        (Math.min(...lngs) + Math.max(...lngs)) / 2,
      ];
    }

    map = L.map(mapContainer).setView(center, zoom);

    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution:
        '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
      maxZoom: 19,
    }).addTo(map);

    updateMarkers();

    // Fit bounds if we have locations
    if (locationsWithCoords.length > 0) {
      const bounds = L.latLngBounds(
        locationsWithCoords.map((l) => [l.lat!, l.lng!] as L.LatLngTuple)
      );
      map.fitBounds(bounds, { padding: [50, 50] });
    }
  }

  function updateMarkers() {
    if (!map) return;

    // Clear existing markers
    markers.forEach((m) => m.remove());
    markers = [];

    // Add markers for locations with coordinates
    locationsWithCoords.forEach((location) => {
      const marker = L.marker([location.lat!, location.lng!], { icon: customIcon }).addTo(map!);

      // Create popup content
      const popupContent = document.createElement("div");
      popupContent.className = "location-popup";
      popupContent.innerHTML = `
        <div style="min-width: 150px;">
          <h3 style="font-weight: 600; margin: 0 0 4px 0;">${location.name}</h3>
          ${location.country ? `<p style="color: #666; margin: 0 0 4px 0; font-size: 0.875rem;">${location.country}</p>` : ""}
          ${location.description ? `<p style="color: #888; margin: 0 0 8px 0; font-size: 0.75rem;">${location.description}</p>` : ""}
          <p style="color: #888; margin: 0; font-size: 0.75rem;">${location.activitiesCount} activities</p>
        </div>
      `;

      const viewBtn = document.createElement("button");
      viewBtn.textContent = "View Details";
      viewBtn.style.cssText =
        "margin-top: 8px; padding: 4px 8px; font-size: 0.75rem; background: #0ea5e9; color: white; border: none; border-radius: 4px; cursor: pointer; width: 100%;";
      viewBtn.onclick = () => viewLocation(location.id);
      popupContent.appendChild(viewBtn);

      marker.bindPopup(popupContent);
      markers.push(marker);
    });
  }

  onMount(() => {
    // Small delay to ensure container is rendered
    setTimeout(initMap, 0);
  });

  onDestroy(() => {
    if (map) {
      map.remove();
      map = null;
    }
  });

  // Track location IDs for reactivity
  const locationIds = $derived(locationsWithCoords.map((l) => l.id).join(","));

  // Update markers when locations change
  $effect(() => {
    // Access locationIds to track changes
    const _ids = locationIds;
    if (map) {
      updateMarkers();

      // Fit bounds to new markers if we have any
      if (locationsWithCoords.length > 0) {
        const bounds = L.latLngBounds(
          locationsWithCoords.map((l) => [l.lat!, l.lng!] as L.LatLngTuple)
        );
        map.fitBounds(bounds, { padding: [50, 50] });
      }
    }
  });
</script>

<div class="w-full h-full min-h-[400px] rounded-lg overflow-hidden border border-border">
  <div bind:this={mapContainer} class="w-full h-full min-h-[400px]"></div>
</div>

{#if locationsWithCoords.length === 0}
  <div class="absolute inset-0 flex items-center justify-center bg-background/80 rounded-lg">
    <div class="text-center p-4">
      <p class="text-muted-foreground">No locations have coordinates yet.</p>
      <p class="text-sm text-muted-foreground mt-1">
        Add lat/lng to your location files to see them on the map.
      </p>
    </div>
  </div>
{/if}

<style>
  :global(.leaflet-container) {
    font-family: inherit;
  }

  :global(.custom-marker) {
    background: none;
    border: none;
  }

  :global(.custom-marker svg) {
    color: #dc2626;
    filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
  }
</style>
