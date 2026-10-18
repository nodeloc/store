// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  ssr: false,
  devtools: { enabled: false },
  css: [
    "@/assets/icons/icons.css",
    "bootstrap/dist/css/bootstrap.min.css",
    "@/assets/scss/main.scss",
    "@/assets/nl/tokens.css",
    "@/assets/nl/components.css",
    "@/assets/nl/compat.css",
    "@/assets/nl/site.css",
  ],
  modules: ["@nuxt/eslint", "@nuxt/image", "@nuxt/scripts", "@pinia/nuxt"],

  image: {
    // Images under /uploads/ are served by the Go backend (proxied via Nitro route rules).
    // IPX's filesystem storage can't find them; using 'none' passes URLs through as-is.
    provider: "none",
  },

  nitro: {
    routeRules: {
      "/api/**": { proxy: `${process.env.BACKEND_URL || "http://backend:8080"}/api/**` },
      "/uploads/**": { proxy: `${process.env.BACKEND_URL || "http://backend:8080"}/uploads/**` },
    },
  },
  runtimeConfig: {
    public: {
      seoDefaults: {
        description: "NodeLoc 社区发卡平台：VPS、AI 账号、点卡与兑换码，积分付款，卡密自动发放。",
        keywords: ["NodeLoc", "发卡", "卡密", "VPS", "AI 账号", "点卡", "兑换码"],
        baseUrl: "https://example.com",
        defaultImage: {
          url: "/og-image.png",
          width: 1200,
          height: 630,
          alt: "Open Graph Image",
        },
        siteName: "NodeLoc 商店",
        twitterCard: "summary_large_image",
        type: "website",
      },
    },
  },

  app: {
    head: {
      htmlAttrs: { lang: "zh-CN" },
      meta: [{ name: "theme-color", content: "#08090b" }],
      script: [
        {
          // Runs before the SPA boots so the stored theme is applied without a light flash.
          innerHTML:
            "try{var t=localStorage.getItem('nl-theme');var r=document.documentElement;r.dataset.theme=t==='light'?'light':'dark';if(!location.pathname.startsWith('/admin')){r.classList.add('nl-on');if(r.dataset.theme==='dark')r.classList.add('dark')}}catch(e){}",
          tagPosition: "head",
        },
      ],
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        { rel: "preconnect", href: "https://fonts.gstatic.com", crossorigin: "" },
        {
          rel: "stylesheet",
          href: "https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@12..96,500..800&family=IBM+Plex+Sans:wght@400;500;600&family=JetBrains+Mono:wght@500;700&display=swap",
        },
        {
          rel: "icon",
          type: "image/x-icon",
          href: "/favicon.ico",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "32x32",
          href: "/favicon-32x32.png",
        },
        {
          rel: "icon",
          type: "image/png",
          sizes: "16x16",
          href: "/favicon-16x16.png",
        },
        {
          rel: "apple-touch-icon",
          sizes: "180x180",
          href: "/apple-touch-icon.png",
        },
      ],
    },
  },
});
