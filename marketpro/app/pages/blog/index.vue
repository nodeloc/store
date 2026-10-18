<template>
  <div>
    <Breadcrumb :title="heading" sub="平台公告、使用指南和商家故事。" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-toolbar" style="margin-top: 0">
        <div class="nl-blogcats" role="list">
          <NuxtLink
            v-for="c in cats"
            :key="c.key"
            role="listitem"
            class="nl-blogcat"
            :to="{ path: '/blog', query: c.key ? { category: c.key } : {} }"
            :aria-current="c.key === category ? 'true' : undefined"
          >{{ c.label }}<span v-if="c.count !== undefined" class="nl-mono">{{ c.count }}</span></NuxtLink>
        </div>
        <form class="nl-toolbar__search" role="search" @submit.prevent="search">
          <label for="blog-q" class="nl-sr">搜索文章</label>
          <input id="blog-q" v-model="keyword" class="nl-input" type="search" placeholder="搜索文章" />
          <button type="submit" aria-label="搜索"><i class="ph ph-magnifying-glass"></i></button>
        </form>
      </div>

      <div v-if="loading" class="nl-postgrid">
        <div v-for="n in 3" :key="n" class="nl-skel"></div>
      </div>
      <NlEmpty
        v-else-if="!posts.length"
        icon="ph ph-article"
        :title="filtered ? '没有找到相关文章' : '还没有发布文章'"
        :text="filtered ? '换个关键词或分类试试。' : '公告和指南发布后会出现在这里。'"
      >
        <NuxtLink v-if="filtered" to="/blog" class="nl-btn nl-btn--secondary nl-btn--sm" @click="keyword = ''">查看全部文章</NuxtLink>
        <NuxtLink v-else to="/shop" class="nl-btn nl-btn--secondary nl-btn--sm">先去逛逛商品</NuxtLink>
      </NlEmpty>
      <div v-else class="nl-postgrid">
        <NuxtLink
          v-for="(p, i) in posts"
          :key="p.id"
          v-reveal
          data-reveal
          :style="{ '--i': i }"
          :to="`/blog/${p.slug}`"
          class="nl-post"
          :class="{ 'nl-post--lead': i === 0 && page === 1 && !filtered }"
        >
          <div class="nl-post__cover" :class="{ 'nl-post__cover--blank': !p.cover_image }">
            <img v-if="p.cover_image" :src="p.cover_image" :alt="p.title" loading="lazy" />
            <template v-else>
              <span class="nl-post__mark" aria-hidden="true">{{ p.category || "NodeLoc" }}</span>
              <i class="ph ph-article" aria-hidden="true"></i>
            </template>
          </div>
          <div class="nl-post__body">
            <span v-if="p.category" class="nl-badge nl-badge--brand">{{ p.category }}</span>
            <h3 class="nl-post__title">{{ p.title }}</h3>
            <p v-if="p.excerpt" class="nl-post__excerpt">{{ p.excerpt }}</p>
            <div class="nl-post__meta nl-mono">
              <span>{{ fmtDate(p.published_at) }}</span>
              <span><i class="ph ph-eye"></i>{{ p.views || 0 }}</span>
            </div>
          </div>
        </NuxtLink>
      </div>

      <NlPager :page="page" :total="totalPages" @go="goPage" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });
useSeo({ title: "博客", description: "NodeLoc 商店的平台公告、使用指南和商家故事。" });

interface Post {
  id: number;
  title: string;
  slug: string;
  category: string;
  cover_image: string;
  excerpt: string;
  published_at: string;
  views: number;
}

const route = useRoute();
const router = useRouter();
const pageSize = 9;

const posts = ref<Post[]>([]);
const total = ref(0);
const loading = ref(true);
const categories = ref<{ name: string; count: number }[]>([]);

const category = computed(() => (typeof route.query.category === "string" ? route.query.category : ""));
const q = computed(() => (typeof route.query.q === "string" ? route.query.q.trim() : ""));
const page = computed(() => Math.max(1, Number(route.query.page) || 1));
const totalPages = computed(() => Math.ceil(total.value / pageSize));
const filtered = computed(() => !!(category.value || q.value));
const keyword = ref(q.value);

const heading = computed(() => (q.value ? `搜索：${q.value}` : category.value || "博客"));
const cats = computed(() => [
  { key: "", label: "全部", count: undefined as number | undefined },
  ...categories.value.map((c) => ({ key: c.name, label: c.name, count: c.count })),
]);

async function load() {
  loading.value = true;
  try {
    const res = await $fetch<{ posts: Post[]; total: number }>("/api/blog/posts", {
      query: { page: page.value, page_size: pageSize, category: category.value || undefined, q: q.value || undefined },
    });
    posts.value = res.posts ?? [];
    total.value = res.total ?? 0;
  } catch {
    posts.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function search() {
  router.push({ query: { category: category.value || undefined, q: keyword.value.trim() || undefined } });
}

function goPage(p: number) {
  router.push({ query: { ...route.query, page: p > 1 ? p : undefined } });
  window.scrollTo({ top: 0, behavior: "smooth" });
}

function fmtDate(d: string) {
  return d ? new Date(d).toLocaleDateString("zh-CN", { year: "numeric", month: "2-digit", day: "2-digit" }) : "";
}

watch(() => [route.query.category, route.query.q, route.query.page], load, { immediate: true });

onMounted(async () => {
  try {
    const res = await $fetch<{ name: string; count: number }[]>("/api/blog/categories");
    categories.value = Array.isArray(res) ? res : [];
  } catch {
    categories.value = [];
  }
});
</script>
