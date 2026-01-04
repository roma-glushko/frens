<script lang="ts">
  import { cn } from "$lib/utils";
  import { Heart, Menu, X } from "lucide-svelte";
  import Navigation from "./Navigation.svelte";
  import ThemeToggle from "./ThemeToggle.svelte";
  import Button from "$lib/components/ui/button/Button.svelte";

  let mobileMenuOpen = $state(false);

  function toggleMobileMenu() {
    mobileMenuOpen = !mobileMenuOpen;
  }

  function closeMobileMenu() {
    mobileMenuOpen = false;
  }
</script>

<header class="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
  <div class="container mx-auto px-4">
    <div class="flex h-14 items-center justify-between">
      <!-- Logo -->
      <a href="/" class="flex items-center gap-2 font-semibold">
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
          <Heart class="h-4 w-4 text-primary-foreground" />
        </div>
        <span class="text-lg">Frens</span>
      </a>

      <!-- Desktop Navigation -->
      <div class="hidden md:flex md:items-center md:gap-6">
        <Navigation />
      </div>

      <!-- Right side actions -->
      <div class="flex items-center gap-2">
        <ThemeToggle />

        <!-- Mobile menu button -->
        <Button
          variant="ghost"
          size="icon"
          class="md:hidden"
          onclick={toggleMobileMenu}
        >
          {#if mobileMenuOpen}
            <X class="h-5 w-5" />
          {:else}
            <Menu class="h-5 w-5" />
          {/if}
          <span class="sr-only">Toggle menu</span>
        </Button>
      </div>
    </div>

    <!-- Mobile Navigation -->
    {#if mobileMenuOpen}
      <div class="border-t border-border md:hidden">
        <nav class="flex flex-col py-4">
          <Navigation mobile={true} onNavigate={closeMobileMenu} />
        </nav>
      </div>
    {/if}
  </div>
</header>
