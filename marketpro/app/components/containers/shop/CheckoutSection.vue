<template>
  <div class="nl-container" style="padding-bottom: 96px">
    <NlEmpty v-if="!cart.items.length" icon="ph ph-shopping-cart-simple" title="没有待结算的商品" text="先把商品加入购物车。">
      <NuxtLink class="nl-btn nl-btn--volt" to="/shop"><i class="ph-fill ph-lightning"></i>去逛逛</NuxtLink>
    </NlEmpty>

    <div v-else class="nl-cartgrid">
      <div class="nl-stack" style="gap: 16px">
        <div class="nl-panel">
          <h2 class="nl-panel__title">下单账号</h2>
          <div v-if="auth.isAuthenticated" class="nl-userline">
            <span class="nl-avatar">
              <img v-if="auth.user?.avatar_url" :src="auth.user.avatar_url" alt="" />
              <i v-else class="ph ph-user"></i>
            </span>
            <div>
              <b>{{ auth.user?.username }}</b>
              <span v-if="auth.user?.email">{{ auth.user.email }}</span>
            </div>
          </div>
          <div v-else class="nl-row">
            <span class="nl-note nl-note--warning" style="flex: 1"><i class="ph ph-warning"></i>请先登录 NodeLoc 账号再下单。</span>
            <a class="nl-btn nl-btn--primary nl-btn--sm" href="/auth/login?redirect=/checkout"><i class="ph ph-sign-in"></i>登录</a>
          </div>
        </div>

        <div class="nl-panel nl-stack" style="gap: 16px">
          <label class="nl-field">
            <span>备注（选填）</span>
            <textarea v-model="remark" class="nl-input" rows="2" placeholder="如有特殊说明可填写"></textarea>
          </label>
          <div class="nl-field">
            <span>优惠券</span>
            <div class="nl-inline-form">
              <input v-model="couponCode" class="nl-input" type="text" placeholder="输入优惠券码" :disabled="!!appliedCoupon" aria-label="优惠券码" />
              <button v-if="!appliedCoupon" type="button" class="nl-btn nl-btn--secondary" :disabled="!couponCode.trim() || couponValidating" @click="validateCoupon">
                {{ couponValidating ? "验证中…" : "使用" }}
              </button>
              <button v-else type="button" class="nl-btn nl-btn--ghost" @click="removeCoupon">取消</button>
            </div>
            <div v-if="couponError" class="nl-note nl-note--danger"><i class="ph ph-warning"></i>{{ couponError }}</div>
            <div v-if="appliedCoupon" class="nl-note nl-note--success">
              <i class="ph-fill ph-ticket"></i>已优惠 {{ money(appliedCoupon.discount) }}
              <template v-if="appliedCoupon.description">（{{ appliedCoupon.description }}）</template>
            </div>
          </div>
        </div>

        <div class="nl-panel">
          <h2 class="nl-panel__title">支付方式</h2>
          <div class="nl-paypick" role="radiogroup" aria-label="支付方式">
            <label v-for="m in paymentMethods" :key="m.id" :class="{ 'is-on': selectedPayment === m.id }">
              <input v-model="selectedPayment" type="radio" name="payment" :value="m.id" />
              <div><b>{{ m.label }}</b><span>{{ m.desc }}</span></div>
              <i :class="m.icon"></i>
            </label>
          </div>
        </div>
      </div>

      <aside class="nl-summary">
        <div class="nl-panel">
          <h2 class="nl-panel__title">订单汇总 <small>{{ cart.count }} 件</small></h2>
          <div class="nl-sum">
            <div v-for="item in cart.items" :key="item.id">
              <span style="min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{ item.name }} × {{ item.qty }}</span>
              <b>{{ money(item.price * item.qty) }}</b>
            </div>
            <div v-if="appliedCoupon"><span>优惠</span><b style="color: var(--success)">-{{ money(appliedCoupon.discount) }}</b></div>
            <div class="nl-sum__total"><span>应付</span><NlPrice :now="finalTotal" lg /></div>
          </div>
          <div v-if="errorMsg" class="nl-note nl-note--danger" style="margin-top: 16px"><i class="ph ph-warning"></i>{{ errorMsg }}</div>
          <button
            type="button"
            class="nl-btn nl-btn--volt nl-btn--lg nl-btn--block"
            style="margin-top: 20px"
            :disabled="submitting || !auth.isAuthenticated"
            @click="placeOrder"
          >
            <template v-if="submitting">处理中…</template>
            <template v-else><i class="ph-fill ph-lightning"></i>确认下单</template>
          </button>
        </div>
        <div class="nl-note nl-note--brand"><i class="ph ph-key"></i>付款成功后卡密自动发放到订单页。</div>
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
const router = useRouter();

const remark = ref("");
const selectedPayment = ref("balance");
const submitting = ref(false);
const errorMsg = ref("");

const couponCode = ref("");
const couponValidating = ref(false);
const couponError = ref("");
const appliedCoupon = ref<{ coupon_id: number; code: string; discount: number; final_amount: number; description: string } | null>(null);

const finalTotal = computed(() => (appliedCoupon.value ? appliedCoupon.value.final_amount : cart.total));

const paymentMethods = [
  { id: "balance", label: "余额支付", desc: "使用账户余额，即时完成", icon: "ph-fill ph-wallet" },
  { id: "nodeloc", label: "NodeLoc 能量支付", desc: "跳转 NodeLoc 完成支付", icon: "ph-fill ph-lightning" },
];

async function validateCoupon() {
  couponError.value = "";
  if (!couponCode.value.trim()) return;
  couponValidating.value = true;
  try {
    const shopId = (cart.items[0] as { shop_id?: number } | undefined)?.shop_id ?? 0;
    appliedCoupon.value = await $fetch("/api/coupons/validate", {
      method: "POST",
      credentials: "include",
      body: { code: couponCode.value.trim().toUpperCase(), shop_id: shopId, amount: cart.total },
    });
  } catch (e: unknown) {
    couponError.value = (e as { data?: { error?: string } })?.data?.error || "优惠券无效";
    appliedCoupon.value = null;
  } finally {
    couponValidating.value = false;
  }
}

function removeCoupon() {
  appliedCoupon.value = null;
  couponCode.value = "";
  couponError.value = "";
}

async function placeOrder() {
  if (!cart.items.length) {
    errorMsg.value = "购物车为空";
    return;
  }
  const contactVal = auth.user?.email || auth.user?.username || "";
  submitting.value = true;
  errorMsg.value = "";
  try {
    const orderNos: string[] = [];
    for (const item of cart.items) {
      const body: Record<string, unknown> = {
        product_id: item.id,
        quantity: item.qty,
        contact: contactVal,
        remark: remark.value,
        pay_method: selectedPayment.value,
      };
      // Coupons do not stack: only the first order line gets it.
      if (appliedCoupon.value && orderNos.length === 0) body.coupon_id = appliedCoupon.value.coupon_id;
      const res = await $fetch<{ order_no?: string; payment_url?: string }>("/api/orders/create", {
        method: "POST",
        credentials: "include",
        signal: AbortSignal.timeout(15000),
        body,
      });
      if (res.order_no) orderNos.push(res.order_no);
      if (selectedPayment.value !== "balance" && res.payment_url) {
        cart.clear();
        window.location.href = res.payment_url;
        return;
      }
    }
    cart.clear();
    router.push(orderNos.length ? `/order/${orderNos[0]}` : "/account/orders");
  } catch (e: unknown) {
    const err = e as { name?: string; cause?: { name?: string }; data?: { error?: string } };
    errorMsg.value =
      err?.name === "TimeoutError" || err?.cause?.name === "TimeoutError" ? "请求超时，请稍后重试" : err?.data?.error || "下单失败，请重试";
  } finally {
    submitting.value = false;
  }
}
</script>
