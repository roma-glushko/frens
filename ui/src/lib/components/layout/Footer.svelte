<script lang="ts">
  import { Heart, Github, ExternalLink } from "lucide-svelte";

  const currentYear = new Date().getFullYear();

  const links = [
    { label: "Documentation", href: "https://github.com/roma-glushko/frens#readme" },
    { label: "GitHub", href: "https://github.com/roma-glushko/frens" },
    { label: "Issues", href: "https://github.com/roma-glushko/frens/issues" },
  ];
</script>

<footer class="border-t border-border bg-background">
  <div class="container mx-auto px-4">
    <!-- Main footer content -->
    <div class="py-8 md:py-12">
      <div class="grid gap-8 md:grid-cols-3">
        <!-- Brand -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
              <Heart class="h-4 w-4 text-primary-foreground" />
            </div>
            <span class="text-lg font-semibold">Frens</span>
          </div>
          <p class="text-sm text-muted-foreground max-w-xs">
            A friendship management and journaling app for introverts. Build relationships that last.
          </p>
        </div>

        <!-- Quick Links -->
        <div class="space-y-4">
          <h3 class="text-sm font-semibold">Resources</h3>
          <ul class="space-y-2">
            {#each links as link}
              <li>
                <a
                  href={link.href}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  {link.label}
                  <ExternalLink class="h-3 w-3" />
                </a>
              </li>
            {/each}
          </ul>
        </div>

        <!-- CLI Info -->
        <div class="space-y-4">
          <h3 class="text-sm font-semibold">CLI Commands</h3>
          <div class="space-y-2 text-sm text-muted-foreground">
            <p><code class="rounded bg-muted px-1.5 py-0.5">frens friend add</code> - Add a friend</p>
            <p><code class="rounded bg-muted px-1.5 py-0.5">frens activity add</code> - Log activity</p>
            <p><code class="rounded bg-muted px-1.5 py-0.5">frens journal stats</code> - View stats</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom bar -->
    <div class="border-t border-border py-4">
      <div class="flex flex-col items-center justify-between gap-4 md:flex-row">
        <p class="text-sm text-muted-foreground">
          &copy; {currentYear} Roma Hlushko. Apache-2.0 License.
        </p>
        <a
          href="https://github.com/roma-glushko/frens"
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
        >
          <Github class="h-4 w-4" />
          Star on GitHub
        </a>
      </div>
    </div>
  </div>
</footer>
