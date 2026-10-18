<template>
  <article
    v-reveal
    v-spotlight
    class="nl-card"
    :class="{ 'nl-card--soldout': stock.out }"
    data-reveal
    :style="{ '--i': index % 8 }"
  >
    <div class="nl-card__thumb">
      <img v-if="product.image" :src="product.image" :alt="product.name" loading="lazy" />
      <i v-else :class="icon" aria-hidden="true"></i>
      <div class="nl-card__flags">
        <span v-for="f in flags" :key="f.label" class="nl-badge" :class="f.tone && `nl-badge--${f.tone}`">
          <i v-if="f.icon" :class="f.icon"></i>{{ f.label }}
        </span>
      </div>
      <button
        type="button"
        class="nl-card__fav"
        :aria-label="fav ? '取消收藏' : '收藏'"
        :aria-pressed="fav ? 'true' : 'false'"
        @click="toggleFav"
      >
        <i :class="fav ? 'ph-fill ph-heart' : 'ph ph-heart'"></i>
      </button>
    </div>
    <div class="nl-card__body">
      <div class="nl-card__seller"><i class="ph ph-storefront"></i>{{ product.shop?.name || "平台自营" }}</div>
      <h3 class="nl-card__title">
        <NuxtLink :to="`/product/${product.id}`">{{ product.name }}</NuxtLink>
      </h3>
      <div class="nl-card__meta">
        <template v-if="(product.review_count || 0) > 0 && product.avg_rating">
          <span class="star"><i class="ph-fill ph-star"></i></span>
          {{ product.avg_rating.toFixed(1) }} ·
        </template>
        已售 {{ product.sales_count || 0 }}
      </div>
      <NlStock :sales="product.sales_count" :stock="product.stock_count" />
      <div class="nl-card__foot">
        <NlPrice :now="price.now" :was="price.was" :promo="price.promo" />
        <button
          type="button"
          class="nl-card__add"
          :class="{ 'is-added': addedTick > 0 }"
          :disabled="stock.out"
          :aria-label="stock.out ? '已售罄' : '加入购物车'"
          @click="addToCart"
        >
          <i :key="addedTick" :class="addedTick > 0 ? 'ph ph-check' : 'ph ph-shopping-cart-simple'"></i>
        </button>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { useCartStore } from "~/stores/cart";
import { useWishlistStore } from "~/stores/wishlist";
import { categoryIcon, priceOf, productFlags, stockOf, type NlProduct } from "~/utils/nl";

const props = withDefaults(defineProps<{ product: NlProduct; index?: number }>(), { index: 0 });

const cart = useCartStore();
const wishlist = useWishlistStore();

const price = computed(() => priceOf(props.product));
const stock = computed(() => stockOf(props.product));
const flags = computed(() => productFlags(props.product));
const icon = computed(() => categoryIcon(props.product.category).cls);
const fav = computed(() => wishlist.has(props.product.id));

const addedTick = ref(0);
let resetTimer: ReturnType<typeof setTimeout> | null = null;

function addToCart() {
  if (stock.value.out) return;
  cart.add({ ...props.product, price: price.value.now });
  addedTick.value++;
  if (resetTimer) clearTimeout(resetTimer);
  resetTimer = setTimeout(() => (addedTick.value = 0), 1600);
}

function toggleFav() {
  wishlist.toggle({ ...props.product, price: price.value.now });
}

onBeforeUnmount(() => resetTimer && clearTimeout(resetTimer));
</script>
