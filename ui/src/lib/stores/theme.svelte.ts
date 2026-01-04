export type Theme = "light" | "dark" | "system";

class ThemeState {
  value = $state<Theme>(getInitialTheme());

  constructor() {
    // Apply theme on initialization
    this.applyTheme();

    // Listen for system theme changes
    window.matchMedia("(prefers-color-scheme: dark)").addEventListener("change", () => {
      if (this.value === "system") {
        this.applyTheme();
      }
    });
  }

  setTheme(theme: Theme) {
    this.value = theme;
    localStorage.setItem("frens-theme", theme);
    this.applyTheme();
  }

  applyTheme() {
    const root = document.documentElement;
    const isDark =
      this.value === "dark" ||
      (this.value === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches);

    root.classList.toggle("dark", isDark);
  }

  get isDark(): boolean {
    return (
      this.value === "dark" ||
      (this.value === "system" && window.matchMedia("(prefers-color-scheme: dark)").matches)
    );
  }
}

function getInitialTheme(): Theme {
  if (typeof window === "undefined") return "system";

  const stored = localStorage.getItem("frens-theme") as Theme | null;
  if (stored && ["light", "dark", "system"].includes(stored)) {
    return stored;
  }
  return "system";
}

export const theme = new ThemeState();
