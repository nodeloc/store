import type { Directive } from "vue";
import { prefersReducedMotion } from "~/utils/nl";

type RevealEl = HTMLElement & { _nlIn?: boolean; _nlIo?: IntersectionObserver };

// Vue rewrites `class` on re-render, so the revealed state is kept on the element and re-applied.
const reveal: Directive<RevealEl> = {
  mounted(el) {
    if (prefersReducedMotion() || !("IntersectionObserver" in window)) {
      el._nlIn = true;
      el.classList.add("is-in");
      return;
    }
    el._nlIo = new IntersectionObserver(
      entries => {
        for (const e of entries) {
          if (e.isIntersecting) {
            el._nlIn = true;
            el.classList.add("is-in");
            el._nlIo?.disconnect();
          }
        }
      },
      { threshold: 0.12, rootMargin: "0px 0px -40px 0px" }
    );
    el._nlIo.observe(el);
  },
  updated(el) {
    if (el._nlIn) el.classList.add("is-in");
  },
  unmounted(el) {
    el._nlIo?.disconnect();
  },
};

const spotlight: Directive<HTMLElement> = {
  mounted(el) {
    el.addEventListener("pointermove", e => {
      const r = el.getBoundingClientRect();
      el.style.setProperty("--mx", `${e.clientX - r.left}px`);
      el.style.setProperty("--my", `${e.clientY - r.top}px`);
    });
  },
};

export default defineNuxtPlugin(nuxtApp => {
  nuxtApp.vueApp.directive("reveal", reveal);
  nuxtApp.vueApp.directive("spotlight", spotlight);
});
