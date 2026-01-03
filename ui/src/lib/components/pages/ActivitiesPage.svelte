<script lang="ts">
  import Card from "$lib/components/ui/card/Card.svelte";
  import CardContent from "$lib/components/ui/card/CardContent.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";
  import { Calendar, Plus, Filter } from "lucide-svelte";
</script>

<div class="container mx-auto px-4 py-8">
  <!-- Page Header -->
  <div class="flex items-center justify-between mb-8">
    <div>
      <h1 class="text-3xl font-bold tracking-tight">Activities</h1>
      <p class="text-muted-foreground mt-2">
        Log and track your social interactions
      </p>
    </div>
    <div class="flex gap-2">
      <Button variant="outline">
        <Filter class="h-4 w-4 mr-2" />
        Filter
      </Button>
      <Button>
        <Plus class="h-4 w-4 mr-2" />
        Log Activity
      </Button>
    </div>
  </div>

  <!-- Empty State -->
  <Card>
    <CardContent class="py-16">
      <div class="flex flex-col items-center justify-center text-center">
        <div class="rounded-full bg-muted p-4 mb-4">
          <Calendar class="h-8 w-8 text-muted-foreground" />
        </div>
        <h3 class="text-lg font-semibold mb-2">No activities logged</h3>
        <p class="text-sm text-muted-foreground max-w-sm mb-4">
          Start tracking your social interactions. Log meetups, calls, and memorable moments with friends.
        </p>
        <Button>
          <Plus class="h-4 w-4 mr-2" />
          Log Your First Activity
        </Button>
      </div>
    </CardContent>
  </Card>
</div>
