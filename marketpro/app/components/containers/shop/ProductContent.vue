<template>
  <div>
    <div class="nl-tabs" role="tablist">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        role="tab"
        class="nl-tab"
        :aria-selected="activeTab === tab.key ? 'true' : 'false'"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
        <span v-if="tab.key === 'reviews' && stats?.total" class="nl-badge nl-badge--brand">{{ stats.total }}</span>
      </button>
    </div>

    <div v-if="activeTab === 'description'" class="nl-panel">
      <div v-if="description" class="nl-prose">{{ description }}</div>
      <NlEmpty v-else icon="ph ph-article" title="暂无商品介绍" />
    </div>

    <div v-else class="nl-panel">
      <div v-if="stats && stats.total > 0" class="nl-rating">
        <div class="nl-rating__big">
          <b>{{ stats.avg_rating.toFixed(1) }}</b>
          <div class="nl-stars" style="margin: 6px 0">
            <i v-for="n in 5" :key="n" class="ph-fill ph-star" :class="{ 'is-off': n > Math.round(stats.avg_rating) }"></i>
          </div>
          <span>{{ stats.total }} 条评价</span>
        </div>
        <div class="nl-rating__bars">
          <div v-for="n in [5, 4, 3, 2, 1]" :key="n" class="nl-rating__bar">
            <span>{{ n }}★</span>
            <div class="nl-stock__track">
              <div class="nl-stock__fill" :style="{ width: `${((stats.rating_distribution?.[n] || 0) / stats.total) * 100}%` }"></div>
            </div>
            <span>{{ stats.rating_distribution?.[n] || 0 }}</span>
          </div>
        </div>
      </div>

      <div v-if="auth.isAuthenticated && !hasReviewed" class="nl-stack" style="gap: 16px; margin-bottom: 24px">
        <h3 class="nl-panel__title" style="margin-bottom: 0 !important">发表评价</h3>
        <div class="nl-starpick" role="radiogroup" aria-label="评分">
          <button
            v-for="n in 5"
            :key="n"
            type="button"
            role="radio"
            :aria-checked="form.rating === n ? 'true' : 'false'"
            :aria-label="`${n} 星`"
            :class="{ 'is-on': n <= form.rating }"
            @click="form.rating = n"
          >
            <i class="ph-fill ph-star"></i>
          </button>
        </div>
        <label class="nl-field">
          <span>标题（可选）</span>
          <input v-model="form.title" class="nl-input" type="text" maxlength="100" />
        </label>
        <label class="nl-field">
          <span>评价内容</span>
          <textarea v-model="form.content" class="nl-input" rows="3" placeholder="分享你的使用体验…"></textarea>
        </label>
        <div class="nl-row">
          <button type="button" class="nl-btn nl-btn--primary" :disabled="submitting" @click="submitReview">
            <i class="ph ph-paper-plane-tilt"></i>{{ submitting ? "提交中…" : "提交评价" }}
          </button>
          <span v-if="submitError" class="nl-note nl-note--danger" style="padding: 8px 12px"><i class="ph ph-warning"></i>{{ submitError }}</span>
        </div>
      </div>
      <div v-else-if="auth.isAuthenticated" class="nl-note nl-note--success" style="margin-bottom: 16px">
        <i class="ph ph-check-circle"></i>你已评价过该商品
      </div>
      <a v-else class="nl-btn nl-btn--secondary nl-btn--sm" style="margin-bottom: 16px" :href="`/auth/login?redirect=${encodeURIComponent(route.fullPath)}`">
        <i class="ph ph-sign-in"></i>登录后评价
      </a>

      <div v-if="reviewsLoading" class="nl-skel" style="min-height: 120px"></div>
      <NlEmpty v-else-if="!reviews.length" icon="ph ph-star" title="暂无评价" text="购买后来抢先评价吧。" />
      <div v-else>
        <div v-for="r in reviews" :key="r.id" class="nl-review">
          <span class="nl-avatar">
            <img v-if="r.user?.avatar_url" :src="r.user.avatar_url" alt="" />
            <i v-else class="ph ph-user"></i>
          </span>
          <div style="min-width: 0">
            <div class="nl-review__head">
              <b>{{ r.user?.username || "匿名用户" }}</b>
              <span class="nl-stars" style="font-size: 12px">
                <i v-for="n in 5" :key="n" class="ph-fill ph-star" :class="{ 'is-off': n > r.rating }"></i>
              </span>
              <span class="nl-mono">{{ fmtDateTime(r.created_at) }}</span>
            </div>
            <p v-if="r.title" style="font-weight: 600">{{ r.title }}</p>
            <p>{{ r.content }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from "~/stores/auth";
import { fmtDateTime } from "~/utils/nl";

const props = defineProps<{ productId: number; description?: string }>();

const auth = useAuthStore();
const route = useRoute();

const tabs = [
  { key: "description", label: "商品介绍" },
  { key: "reviews", label: "用户评价" },
];
const activeTab = ref("description");

interface Review {
  id: number;
  rating: number;
  title?: string;
  content?: string;
  created_at: string;
  user?: { id?: number; username?: string; avatar_url?: string };
}
interface ReviewStats {
  total: number;
  avg_rating: number;
  rating_distribution: Record<number, number>;
}

const reviews = ref<Review[]>([]);
const stats = ref<ReviewStats | null>(null);
const reviewsLoading = ref(false);
const hasReviewed = ref(false);
const form = ref({ rating: 5, title: "", content: "" });
const submitting = ref(false);
const submitError = ref("");

async function loadReviews() {
  if (!props.productId) return;
  reviewsLoading.value = true;
  try {
    const data = await $fetch<{ reviews: Review[]; stats: ReviewStats }>(`/api/products/${props.productId}/reviews`, {
      credentials: "include",
    });
    reviews.value = data.reviews || [];
    stats.value = data.stats || null;
    hasReviewed.value = !!auth.user && reviews.value.some(r => r.user?.id === auth.user?.id);
  } catch {
    /* keep previous */
  } finally {
    reviewsLoading.value = false;
  }
}

async function submitReview() {
  if (!form.value.content.trim()) {
    submitError.value = "请填写评价内容";
    return;
  }
  submitting.value = true;
  submitError.value = "";
  try {
    await $fetch(`/api/products/${props.productId}/reviews`, {
      method: "POST",
      credentials: "include",
      body: { ...form.value },
    });
    hasReviewed.value = true;
    form.value = { rating: 5, title: "", content: "" };
    await loadReviews();
  } catch (e: unknown) {
    const err = e as { data?: { error?: string }; message?: string };
    submitError.value = err?.data?.error || err?.message || "提交失败，请稍后重试";
  } finally {
    submitting.value = false;
  }
}

watch(() => props.productId, id => id && loadReviews(), { immediate: true });
</script>
