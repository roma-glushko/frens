<script lang="ts">
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import { MapPin, Plus, Search } from "lucide-svelte";
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Locations</h1>
      <p class="text-muted-foreground mt-2">
        Save meaningful places and venues
      </p>
    </div>
    <Button>
      <Plus class="h-4 w-4 mr-2" />
      Add Location
    </Button>
  </div>

  <!-- Search -->
  <div class="mb-6">
    <div class="relative max-w-md">
      <Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <input
        type="text"
        placeholder="Search locations..."
        class="w-full rounded-md border border-input bg-background px-10 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
    </div>
  </div>

  <!-- Empty State -->
  <Card>
    <CardContent class="py-16">
      <div class="flex flex-col items-center justify-center text-center">
        <div class="rounded-full bg-muted p-4 mb-4">
          <MapPin class="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 class="text-lg font-semibold mb-2">No locations saved</h3>
        <p class="text-sm text-muted-foreground max-w-sm mb-4">
          Save your favorite spots, restaurants, and venues. Link them to activities and friends for richer memories.
        </p>
        <Button>
          <Plus class="h-4 w-4 mr-2" />
          Add Your First Location
        </Button>
      </div>
    </CardContent>
  </Card>
</div>
