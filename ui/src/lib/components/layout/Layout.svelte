<script lang="ts">
  import { type Snippet } from "svelte";
  import Header from "./Header.svelte";
  import Footer from "./Footer.svelte";

  interface Props {
    children?: Snippet;
  }

  let { children }: Props = $props();
</script>

<div class="flex min-h-screen flex-col bg-background text-foreground">
  <Header />

  <main class="flex-1">
    {#if children}
      {@render children()}
    {/if}
  </main>

  <Footer />
</div>
