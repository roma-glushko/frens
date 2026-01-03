// Simple client-side router state
class RouterState {
  value = $state(window.location.pathname);

  constructor() {
    // Listen for browser back/forward navigation
    window.addEventListener("popstate", () => {
      this.value = window.location.pathname;
    });
  }

  navigate(path: string) {
    if (this.value !== path) {
      this.value = path;
      window.history.pushState({}, "", path);
    }
  }
}

export const currentPath = new RouterState();

// Update the path and push to history when value changes
$effect.root(() => {
  $effect(() => {
    const path = currentPath.value;
    if (window.location.pathname !== path) {
      window.history.pushState({}, "", path);
    }
  });
});
