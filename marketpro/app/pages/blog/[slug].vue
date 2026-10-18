<template>
  <div>
    <Breadcrumb :title="post?.title || '文章'" compact :trail="[{ label: '博客', to: '/blog' }]" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div v-if="loading" class="nl-article">
        <div class="nl-skel" style="min-height: 420px"></div>
      </div>
      <NlEmpty v-else-if="!post" icon="ph ph-article" title="文章不存在或已下线" text="它可能已被删除，或者链接有误。">
        <NuxtLink to="/blog" class="nl-btn nl-btn--secondary nl-btn--sm">返回博客</NuxtLink>
      </NlEmpty>

      <template v-else>
        <article class="nl-article">
          <header class="nl-article__head">
            <NuxtLink v-if="post.category" :to="{ path: '/blog', query: { category: post.category } }" class="nl-badge nl-badge--brand">{{ post.category }}</NuxtLink>
            <h1 class="nl-article__title">{{ post.title }}</h1>
            <div class="nl-post__meta nl-mono">
              <span><i class="ph ph-calendar-blank"></i>{{ fmtDate(post.published_at) }}</span>
              <span><i class="ph ph-eye"></i>{{ post.views || 0 }} 次阅读</span>
            </div>
          </header>
          <img v-if="post.cover_image" :src="post.cover_image" :alt="post.title" class="nl-article__cover" />
          <!-- eslint-disable-next-line vue/no-v-html -->
          <div class="nl-article__body" v-html="html"></div>
        </article>

        <nav v-if="prev || next" class="nl-article nl-postnav" aria-label="上下篇">
          <NuxtLink v-if="prev" :to="`/blog/${prev.slug}`" class="nl-postnav__item">
            <span><i class="ph ph-arrow-left"></i>上一篇</span>
            <b>{{ prev.title }}</b>
          </NuxtLink>
          <span v-else></span>
          <NuxtLink v-if="next" :to="`/blog/${next.slug}`" class="nl-postnav__item nl-postnav__item--next">
            <span>下一篇<i class="ph ph-arrow-right"></i></span>
            <b>{{ next.title }}</b>
          </NuxtLink>
        </nav>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });

interface Post {
  id: number;
  title: string;
  slug: string;
  category: string;
  cover_image: string;
  excerpt: string;
  content: string;
  published_at: string;
  views: number;
}
type PostLink = Pick<Post, "title" | "slug">;

const route = useRoute();
const site = useSiteStore();
const loading = ref(true);
const post = ref<Post | null>(null);
const prev = ref<PostLink | null>(null);
const next = ref<PostLink | null>(null);

const esc = (s: string) => s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");

// Content is authored by admins: if it already contains block-level HTML render it as-is,
// otherwise treat it as plain text and turn blank-line separated blocks into paragraphs.
const html = computed(() => {
  const c = post.value?.content || "";
  if (/<(p|h[1-6]|ul|ol|div|img|blockquote|pre|table)[\s>]/i.test(c)) return c;
  return c
    .split(/\n{2,}/)
    .map((p) => p.trim())
    .filter(Boolean)
    .map((p) => `<p>${esc(p).replace(/\n/g, "<br>")}</p>`)
    .join("");
});

async function load(slug: string) {
  loading.value = true;
  post.value = null;
  prev.value = next.value = null;
  try {
    post.value = await $fetch<Post>(`/api/blog/posts/${encodeURIComponent(slug)}`);
  } catch {
    post.value = null;
  } finally {
    loading.value = false;
  }
  if (!post.value) return;
  try {
    const res = await $fetch<{ posts: Post[] }>("/api/blog/posts", { query: { page: 1, page_size: 100 } });
    const all = res.posts ?? [];
    const idx = all.findIndex((p) => p.slug === slug);
    prev.value = idx > 0 ? all[idx - 1] : null;
    next.value = idx >= 0 && idx < all.length - 1 ? all[idx + 1] : null;
  } catch {
    /* neighbours are optional */
  }
}

watch(() => route.params.slug, (slug) => slug && load(String(slug)), { immediate: true });

watch(post, (p) => {
  if (!p) return;
  const siteName = (site.settings.site_name as string) || "NodeLoc";
  useHead({
    title: `${p.title} | ${siteName}`,
    meta: [
      { name: "description", content: p.excerpt || p.title },
      { property: "og:title", content: p.title },
      { property: "og:description", content: p.excerpt || p.title },
      ...(p.cover_image ? [{ property: "og:image", content: p.cover_image }] : []),
      { property: "og:type", content: "article" },
    ],
  });
});

function fmtDate(d: string) {
  return d ? new Date(d).toLocaleDateString("zh-CN", { year: "numeric", month: "long", day: "numeric" }) : "";
}
</script>
