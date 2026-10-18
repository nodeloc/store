<template>
  <div class="nl-container" style="padding-bottom: 96px">
    <NlEmpty v-if="!cart.items.length" icon="ph ph-shopping-cart-simple" title="购物车是空的" text="去挑几件好物吧。">
      <NuxtLink class="nl-btn nl-btn--volt" to="/shop"><i class="ph-fill ph-lightning"></i>去逛逛</NuxtLink>
    </NlEmpty>

    <div v-else class="nl-cartgrid">
      <div class="nl-stack">
        <div class="nl-cartlist">
          <div v-for="item in cart.items" :key="item.id" class="nl-cartrow">
            <NuxtLink :to="`/product/${item.id}`" class="nl-cartrow__thumb">
              <img v-if="item.image" :src="item.image" :alt="item.name" />
              <i v-else class="ph ph-package"></i>
            </NuxtLink>
            <div style="min-width: 0">
              <NuxtLink :to="`/product/${item.id}`" class="nl-cartrow__name">{{ item.name }}</NuxtLink>
              <div class="nl-cartrow__sub">单价 <NlPrice :now="item.price" /></div>
            </div>
            <div class="nl-qty nl-qty--sm nl-cartrow__qty" role="group" :aria-label="`${item.name} 数量`">
              <button type="button" :disabled="item.qty <= 1" aria-label="减少" @click="cart.updateQty(item.id, item.qty - 1)"><i class="ph ph-minus"></i></button>
              <input :value="item.qty" type="number" min="1" :max="item.stock" aria-label="数量" @change="cart.updateQty(item.id, Number(($event.target as HTMLInputElement).value) || 1)" />
              <button type="button" :disabled="item.qty >= (item.stock || 99)" aria-label="增加" @click="cart.updateQty(item.id, item.qty + 1)"><i class="ph ph-plus"></i></button>
            </div>
            <div class="nl-cartrow__total"><NlPrice :now="item.price * item.qty" /></div>
            <button type="button" class="nl-iconbtn nl-cartrow__del" :aria-label="`移除 ${item.name}`" style="color: var(--danger); font-size: 18px" @click="cart.remove(item.id)">
              <i class="ph ph-trash"></i>
            </button>
          </div>
        </div>
        <div class="nl-row">
          <NuxtLink class="nl-btn nl-btn--ghost nl-btn--sm" to="/shop"><i class="ph ph-arrow-left"></i>继续购物</NuxtLink>
          <button type="button" class="nl-btn nl-btn--ghost nl-btn--sm" style="color: var(--ink-muted)" @click="clearAll">
            <i class="ph ph-broom"></i>清空购物车
          </button>
        </div>
      </div>

      <aside class="nl-summary">
        <div class="nl-panel">
          <h2 class="nl-panel__title">订单汇总</h2>
          <div class="nl-sum">
            <div><span>商品（{{ cart.count }} 件）</span><b>{{ money(cart.total) }}</b></div>
            <div><span>发货</span><b style="color: var(--success)">自动 · 免运费</b></div>
            <div class="nl-sum__total"><span>合计</span><NlPrice :now="cart.total" lg /></div>
          </div>
          <div style="margin-top: 20px">
            <NuxtLink v-if="auth.isAuthenticated" class="nl-btn nl-btn--volt nl-btn--lg nl-btn--block" to="/checkout">
              前往结算 <i class="ph ph-arrow-right"></i>
            </NuxtLink>
            <a v-else class="nl-btn nl-btn--volt nl-btn--lg nl-btn--block" href="/auth/login?redirect=/checkout">
              <i class="ph ph-sign-in"></i>登录后结算
            </a>
          </div>
        </div>
        <div class="nl-note nl-note--brand"><i class="ph ph-shield-check"></i>确认收货后款项才会结算给卖家。</div>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useCartStore } from "~/stores/cart";
import { useAuthStore } from "~/stores/auth";
import { money } from "~/utils/nl";

const cart = useCartStore();
const auth = useAuthStore();

function clearAll() {
  if (confirm("确定清空购物车？")) cart.clear();
}
</script>
