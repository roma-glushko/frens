<script lang="ts">
  import { cn } from "$lib/utils";
  import { Users, Calendar, MapPin, StickyNote, BarChart3 } from "lucide-svelte";
  import { currentPath } from "$lib/stores/router.svelte";

  interface Props {
    mobile?: boolean;
    onNavigate?: () => void;
  }

  let { mobile = false, onNavigate }: Props = $props();

  interface NavItem {
    label: string;
    href: string;
    icon: typeof Users;
  }

  const navItems: NavItem[] = [
    { label: "Dashboard", href: "/", icon: BarChart3 },
    { label: "Friends", href: "/friends", icon: Users },
    { label: "Activities", href: "/activities", icon: Calendar },
    { label: "Notes", href: "/notes", icon: StickyNote },
    { label: "Locations", href: "/locations", icon: MapPin },
  ];

  function isActive(href: string): boolean {
    if (href === "/") {
      return currentPath.value === "/";
    }
    return currentPath.value.startsWith(href);
  }

  function handleClick(e: MouseEvent, href: string) {
    e.preventDefault();
    currentPath.value = href;
    onNavigate?.();
  }
</script>

{#if mobile}
  <!-- Mobile Navigation -->
  {#each navItems as item}
    {@const Icon = item.icon}
    <a
      href={item.href}
      onclick={(e) => handleClick(e, item.href)}
      class={cn(
        "flex items-center gap-3 px-4 py-3 text-sm font-medium transition-colors",
        "hover:bg-accent hover:text-accent-foreground",
        "rounded-md mx-2",
        isActive(item.href)
          ? "bg-accent text-accent-foreground"
          : "text-muted-foreground"
      )}
    >
      <Icon class="h-4 w-4" />
      {item.label}
    </a>
  {/each}
{:else}
  <!-- Desktop Navigation -->
  <nav class="flex items-center gap-1">
    {#each navItems as item}
      {@const Icon = item.icon}
      <a
        href={item.href}
        onclick={(e) => handleClick(e, item.href)}
        class={cn(
          "flex items-center gap-2 px-3 py-2 text-sm font-medium transition-colors rounded-md",
          "hover:bg-accent hover:text-accent-foreground",
          isActive(item.href)
            ? "bg-accent text-accent-foreground"
            : "text-muted-foreground"
        )}
      >
        <Icon class="h-4 w-4" />
        {item.label}
      </a>
    {/each}
  </nav>
{/if}
