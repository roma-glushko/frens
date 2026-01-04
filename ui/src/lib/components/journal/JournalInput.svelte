<script lang="ts">
  import { cn } from "$lib/utils";
  import Button from "$lib/components/ui/button/Button.svelte";
  import {
    Send,
    Lightbulb,
    MessageSquare,
    Activity,
    User,
    MapPin,
    Calendar,
    Gift,
  } from "lucide-svelte";

  type EntryType = "note" | "activity" | "friend" | "location" | "date" | "wishlist";

  const entryTypeConfig: Record<
    EntryType,
    { label: string; icon: typeof MessageSquare; buttonLabel: string }
  > = {
    note: { label: "Note", icon: MessageSquare, buttonLabel: "Note" },
    activity: { label: "Activity", icon: Activity, buttonLabel: "Activity" },
    friend: { label: "Friend", icon: User, buttonLabel: "Friend" },
    location: { label: "Location", icon: MapPin, buttonLabel: "Location" },
    date: { label: "Date", icon: Calendar, buttonLabel: "Date" },
    wishlist: { label: "Wishlist", icon: Gift, buttonLabel: "Gift Idea" },
  };

  let content = $state("");
  let entryType = $state<EntryType>("note");
  let isFocused = $state(false);
  let showHelp = $state(false);

  function handleSubmit() {
    if (!content.trim()) return;
    // TODO: Submit to API
    console.log("Submitting:", { type: entryType, content });
    content = "";
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) {
      e.preventDefault();
      handleSubmit();
    }
  }

  const placeholders: Record<EntryType, string[]> = {
    note: [
      "Remember to call mom about weekend plans #family",
      "Book recommendation from Sarah: The Midnight Library #books",
      "Project idea: build a habit tracker app #ideas",
    ],
    activity: [
      "Had coffee with Sarah at the new place #catchup @blue-bottle",
      "yesterday :: Team lunch, great discussions #work @office",
      "2 days ago :: Movie night with Alex and Jordan #friends @home",
    ],
    friend: [
      "Sarah Miller :: Met at yoga class, works in marketing #yoga #marketing",
      "John Smith :: College roommate, lives in NYC #college @new-york",
      "Emma :: Friend of Alex, loves hiking #hiking",
    ],
    location: [
      "Blue Bottle Coffee :: Great pour-over, quiet for work @hayes-valley #coffee",
      "Central Park :: Perfect for morning runs @new-york #outdoors",
      "The Office :: Coworking space downtown #work @downtown",
    ],
    date: [
      "Sarah :: birthday :: March 15",
      "John :: anniversary :: June 20, 2020",
      "Mom :: birthday :: December 3",
    ],
    wishlist: [
      "Sarah :: Kindle Paperwhite :: She mentioned wanting to read more #books",
      "John :: AirPods Pro :: His old ones broke #tech",
      "Emma :: Hiking boots :: Planning a trip together #hiking #gift",
    ],
  };

  const syntaxHelp: Record<
    EntryType,
    {
      sections: { title: string; description: string; examples: string[] }[];
      example: string;
      quickHints: string[];
    }
  > = {
    note: {
      sections: [
        {
          title: "Date (optional)",
          description: "Start with a date followed by ::",
          examples: ["yesterday :: description", "2 days ago :: description", "Jan 15 :: description"],
        },
        {
          title: "Tags",
          description: "Add anywhere in your text:",
          examples: ["#family", "#work", "#ideas"],
        },
      ],
      example: "yesterday :: Remember to call mom #family",
      quickHints: ["#tag", "date ::"],
    },
    activity: {
      sections: [
        {
          title: "Date (optional)",
          description: "Start with a date followed by ::",
          examples: ["yesterday :: description", "2 days ago :: description", "last week :: description"],
        },
        {
          title: "Tags & Locations",
          description: "Add anywhere in your text:",
          examples: ["#coffee #work - tags", "@cafe-nero @home - locations"],
        },
      ],
      example: "yesterday :: Had coffee with Sarah #catchup @blue-bottle",
      quickHints: ["#tag", "@location", "date ::"],
    },
    friend: {
      sections: [
        {
          title: "Name & Description",
          description: "Use :: to separate name from description:",
          examples: ["Name :: description", "First Last :: how you met, context"],
        },
        {
          title: "Tags & Location",
          description: "Add context about the friend:",
          examples: ["#yoga #marketing - interests/context", "@new-york - where they live"],
        },
      ],
      example: "Sarah Miller :: Met at yoga class, works in marketing #yoga @san-francisco",
      quickHints: ["name :: desc", "#tag", "@location"],
    },
    location: {
      sections: [
        {
          title: "Name & Description",
          description: "Use :: to separate name from description:",
          examples: ["Place Name :: what it's good for", "Venue :: description of the place"],
        },
        {
          title: "Area & Tags",
          description: "Add context about the location:",
          examples: ["@neighborhood - area/region", "#coffee #quiet - categories"],
        },
      ],
      example: "Blue Bottle Coffee :: Great pour-over, quiet for work @hayes-valley #coffee",
      quickHints: ["name :: desc", "@area", "#tag"],
    },
    date: {
      sections: [
        {
          title: "Format",
          description: "Friend :: type :: date:",
          examples: ["Name :: birthday :: March 15", "Name :: anniversary :: June 20, 2020"],
        },
        {
          title: "Date Types",
          description: "Common date types:",
          examples: ["birthday", "anniversary", "first-met"],
        },
      ],
      example: "Sarah :: birthday :: March 15",
      quickHints: ["name :: type :: date"],
    },
    wishlist: {
      sections: [
        {
          title: "Format",
          description: "Friend :: item :: notes:",
          examples: ["Name :: Gift idea :: why they'd like it", "Name :: Item :: context or notes"],
        },
        {
          title: "Tags",
          description: "Categorize gift ideas:",
          examples: ["#books #tech #hiking - categories", "#birthday #christmas - occasions"],
        },
      ],
      example: "Sarah :: Kindle Paperwhite :: She mentioned wanting to read more #books",
      quickHints: ["name :: item :: note", "#tag"],
    },
  };

  let currentExample = $state(0);

  // Cycle placeholder when type changes
  $effect(() => {
    entryType; // dependency
    currentExample = 0;
  });

  const placeholder = $derived(placeholders[entryType][currentExample % placeholders[entryType].length]);
</script>

<div
  class={cn(
    "rounded-lg border bg-card transition-all duration-200",
    isFocused ? "border-ring shadow-md" : "border-border"
  )}
>
  <!-- Type toggle header -->
  <div class="flex items-center gap-1 px-4 pt-3 pb-1 flex-wrap">
    {#each Object.entries(entryTypeConfig) as [type, config]}
      {@const Icon = config.icon}
      <button
        type="button"
        onclick={() => (entryType = type as EntryType)}
        class={cn(
          "flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
          entryType === type
            ? "bg-primary text-primary-foreground"
            : "text-muted-foreground hover:text-foreground hover:bg-muted"
        )}
      >
        <Icon class="h-3.5 w-3.5" />
        {config.label}
      </button>
    {/each}
  </div>

  <!-- Textarea -->
  <div class="px-4 pb-4 pt-2">
    <textarea
      bind:value={content}
      onfocus={() => (isFocused = true)}
      onblur={() => (isFocused = false)}
      onkeydown={handleKeydown}
      placeholder={placeholder}
      rows={3}
      class="w-full resize-none bg-transparent text-foreground placeholder:text-muted-foreground/60 focus:outline-none text-[15px]"
    ></textarea>
  </div>

  <!-- Syntax help panel -->
  {#if showHelp}
    {@const help = syntaxHelp[entryType]}
    <div class="border-t border-border bg-muted/30 px-4 py-3">
      <div class="grid gap-3 sm:grid-cols-2 text-sm">
        {#each help.sections as section}
          <div>
            <p class="font-medium text-foreground mb-1.5">{section.title}</p>
            <p class="text-muted-foreground text-xs">{section.description}</p>
            <div class="mt-1 space-y-0.5 text-xs text-muted-foreground">
              {#each section.examples as ex}
                <p><code class="text-foreground">{ex}</code></p>
              {/each}
            </div>
          </div>
        {/each}
      </div>
      <div class="mt-3 pt-3 border-t border-border">
        <p class="text-xs text-muted-foreground">
          <span class="font-medium text-foreground">Example:</span>
          <code class="ml-1">{help.example}</code>
        </p>
      </div>
    </div>
  {/if}

  <!-- Bottom bar -->
  <div class="flex items-center justify-between border-t border-border px-4 py-2.5">
    <!-- Left side -->
    <div class="flex items-center gap-3">
      <button
        type="button"
        onclick={() => (showHelp = !showHelp)}
        class={cn(
          "flex items-center gap-1.5 text-xs transition-colors rounded px-2 py-1 -ml-2",
          showHelp
            ? "text-foreground bg-muted"
            : "text-muted-foreground hover:text-foreground hover:bg-muted/50"
        )}
      >
        <Lightbulb class="h-3.5 w-3.5" />
        <span>Syntax</span>
      </button>

      <div class="hidden sm:flex items-center gap-2 text-xs text-muted-foreground">
        {#each syntaxHelp[entryType].quickHints as hint}
          <code class="rounded bg-muted px-1.5 py-0.5">{hint}</code>
        {/each}
      </div>
    </div>

    <!-- Right side actions -->
    <div class="flex items-center gap-2">
      {#if content.trim()}
        <span class="text-xs text-muted-foreground hidden sm:inline">
          <kbd class="rounded border border-border bg-muted px-1 py-0.5 text-[10px]">⌘</kbd>
          <kbd class="rounded border border-border bg-muted px-1 py-0.5 text-[10px]">↵</kbd>
        </span>
      {/if}
      <Button
        onclick={handleSubmit}
        disabled={!content.trim()}
        size="sm"
      >
        <Send class="h-4 w-4 sm:mr-2" />
        <span class="hidden sm:inline">Save {entryTypeConfig[entryType].buttonLabel}</span>
      </Button>
    </div>
  </div>
</div>
