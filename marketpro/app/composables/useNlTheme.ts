export type NlTheme = "dark" | "light";

const KEY = "nl-theme";

export function useNlTheme() {
  const theme = useState<NlTheme>("nl-theme", () => {
    if (import.meta.client) {
      const saved = localStorage.getItem(KEY);
      if (saved === "light" || saved === "dark") return saved;
    }
    return "dark";
  });

  // Element Plus reads html.dark for its dark palette; the admin layout must not inherit it.
  function apply(active = true) {
    if (!import.meta.client) return;
    const root = document.documentElement;
    root.dataset.theme = theme.value;
    root.classList.toggle("nl-on", active);
    root.classList.toggle("dark", active && theme.value === "dark");
  }

  function toggle() {
    theme.value = theme.value === "dark" ? "light" : "dark";
    try {
      localStorage.setItem(KEY, theme.value);
    } catch {
      /* storage unavailable */
    }
    apply();
  }

  return { theme, toggle, apply };
}
