<script lang="ts">
  import { cn } from "$lib/utils";
  import { Sun, Moon, Monitor } from "lucide-svelte";
  import { theme, type Theme } from "$lib/stores/theme.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";

  let isOpen = $state(false);

  const themes: { value: Theme; label: string; icon: typeof Sun }[] = [
    { value: "light", label: "Light", icon: Sun },
    { value: "dark", label: "Dark", icon: Moon },
    { value: "system", label: "System", icon: Monitor },
  ];

  function toggleDropdown() {
    isOpen = !isOpen;
  }

  function selectTheme(newTheme: Theme) {
    theme.setTheme(newTheme);
    isOpen = false;
  }

  function handleClickOutside(e: MouseEvent) {
    const target = e.target as HTMLElement;
    if (!target.closest(".theme-toggle")) {
      isOpen = false;
    }
  }

  $effect(() => {
    if (isOpen) {
      document.addEventListener("click", handleClickOutside);
      return () => document.removeEventListener("click", handleClickOutside);
    }
  });
</script>

<div class="theme-toggle relative">
  <Button
    variant="ghost"
    size="icon"
    onclick={toggleDropdown}
    class="relative"
  >
    <Sun class={cn(
      "h-5 w-5 transition-all",
      theme.isDark ? "rotate-90 scale-0" : "rotate-0 scale-100"
    )} />
    <Moon class={cn(
      "absolute h-5 w-5 transition-all",
      theme.isDark ? "rotate-0 scale-100" : "-rotate-90 scale-0"
    )} />
    <span class="sr-only">Toggle theme</span>
  </Button>

  {#if isOpen}
    <div class="absolute right-0 top-full mt-2 w-36 rounded-md border border-border bg-popover p-1 shadow-lg">
      {#each themes as t}
        {@const Icon = t.icon}
        <button
          onclick={() => selectTheme(t.value)}
          class={cn(
            "flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm",
            "hover:bg-accent hover:text-accent-foreground",
            "transition-colors cursor-pointer",
            theme.value === t.value && "bg-accent text-accent-foreground"
          )}
        >
          <Icon class="h-4 w-4" />
          {t.label}
        </button>
      {/each}
    </div>
  {/if}
</div>
