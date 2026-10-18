<template>
  <div>
    <Breadcrumb title="订单详情" :sub="order ? `订单号 ${order.order_no}` : ''" :trail="[{ label: '我的订单', to: '/account/orders' }]" />

    <div class="nl-container" style="padding-bottom: 96px">
      <div v-if="loading" class="nl-order">
        <div class="nl-skel"></div>
        <div class="nl-skel"></div>
      </div>
      <NlEmpty v-else-if="!order" icon="ph ph-receipt-x" title="订单不存在" text="订单号有误，或不属于当前账号。">
        <NuxtLink class="nl-btn nl-btn--secondary nl-btn--sm" to="/account/orders">返回订单列表</NuxtLink>
      </NlEmpty>

      <div v-else class="nl-order">
        <div class="nl-order__main">
          <div class="nl-panel">
            <h2 class="nl-panel__title">
              订单状态
              <span class="nl-badge" :class="status.tone && `nl-badge--${status.tone}`">
                <i :class="status.icon"></i>{{ status.label }}
              </span>
            </h2>
            <div v-if="order.status === 3" class="nl-note">
              <i class="ph ph-x-circle"></i>订单已取消，库存已释放。
            </div>
            <NlSteps
              v-else
              :status="order.status"
              :created-at="order.created_at"
              :paid-at="order.paid_at"
              :updated-at="order.updated_at"
              :auto-delivered="!!order.card_keys?.length || isFixed"
            />
          </div>

          <div v-if="order.status === 0" class="nl-panel">
            <div class="nl-note nl-note--warning" style="margin-bottom: 16px">
              <i class="ph ph-clock"></i>{{ order.product?.delivery_type === 1 ? "订单待支付，付款后等待卖家发货，发货内容会显示在本页面。" : "订单待支付，完成付款后商品会自动发放到本页面。" }}
            </div>
            <a v-if="order.payment_url" class="nl-btn nl-btn--volt nl-btn--lg" :href="order.payment_url">
              <i class="ph-fill ph-lightning"></i>前往支付 <i class="ph ph-arrow-right"></i>
            </a>
          </div>

          <div v-if="order.card_keys?.length" class="nl-panel">
            <h2 class="nl-panel__title">
              <i class="ph-fill ph-key" style="color: var(--accent)"></i>卡密信息
              <small>共 {{ order.card_keys.length }} 张</small>
            </h2>
            <div class="nl-tickets">
              <NlTicket v-for="(k, i) in order.card_keys" :key="i" :index="i" :rows="ticketRows(k)" />
            </div>
          </div>

          <div v-if="order.deliver_content && isFixed && order.status === 2" class="nl-panel">
            <h2 class="nl-panel__title">
              <i class="ph-fill ph-file-text"></i>商品内容
              <button type="button" class="nl-btn nl-btn--ghost nl-btn--sm" style="margin-left: auto" @click="copyDeliver">
                <i :class="copied ? 'ph ph-check' : 'ph ph-copy'"></i>{{ copied ? "已复制" : "复制全部" }}
              </button>
            </h2>
            <pre class="nl-deliver nl-deliver--doc">{{ order.deliver_content }}</pre>
          </div>
          <div v-else-if="order.deliver_content && (order.status === 4 || order.status === 2)" class="nl-panel">
            <h2 class="nl-panel__title"><i class="ph-fill ph-package" style="color: var(--accent)"></i>卖家已发货</h2>
            <div class="nl-deliver">{{ order.deliver_content }}</div>
            <p v-if="order.deliver_note" style="margin-top: 8px; font-size: 13px; color: var(--ink-muted)">备注：{{ order.deliver_note }}</p>
            <div v-if="order.status === 4" style="margin-top: 16px" class="nl-row">
              <button type="button" class="nl-btn nl-btn--primary" :disabled="confirming" @click="confirmReceipt">
                <i class="ph ph-check-circle"></i>{{ confirming ? "确认中…" : "确认收货" }}
              </button>
              <span style="font-size: 13px; color: var(--ink-muted)">确认后交易完成，款项结算给卖家。</span>
            </div>
          </div>
          <div v-else-if="order.status === 1" class="nl-panel">
            <div class="nl-note nl-note--brand"><i class="ph ph-package"></i>已付款，等待卖家发货。发货后内容会显示在这里。</div>
          </div>
        </div>

        <aside class="nl-order__aside">
          <div class="nl-panel">
            <div class="nl-order__product">
              <span class="nl-order__thumb">
                <img v-if="order.product?.image" :src="order.product.image" :alt="order.product.name" />
                <i v-else :class="categoryIcon(order.product?.category).cls"></i>
              </span>
              <div style="min-width: 0">
                <NuxtLink v-if="order.product?.id" :to="`/product/${order.product.id}`">{{ order.product.name }}</NuxtLink>
                <span v-else>{{ order.product?.name || "-" }}</span>
                <div style="font-size: 12px; color: var(--ink-muted)">× {{ order.quantity }}</div>
              </div>
            </div>
            <div class="nl-sum" style="margin-top: 16px">
              <div><span>下单时间</span><b>{{ fmtDateTime(order.created_at) }}</b></div>
              <div v-if="order.paid_at"><span>支付时间</span><b>{{ fmtDateTime(order.paid_at) }}</b></div>
              <div v-if="order.contact"><span>联系方式</span><b>{{ order.contact }}</b></div>
              <div v-if="order.remark"><span>备注</span><b style="font-family: var(--font-sans); font-weight: 500">{{ order.remark }}</b></div>
              <div class="nl-sum__total"><span>实付</span><NlPrice :now="order.total_amount || 0" /></div>
            </div>
          </div>

          <div v-if="order.status === 2 && order.product?.id" class="nl-panel">
            <button v-if="!reviewed" type="button" class="nl-btn nl-btn--secondary nl-btn--block" @click="openReview">
              <i class="ph ph-star"></i>评价商品
            </button>
            <div v-else class="nl-note nl-note--success"><i class="ph-fill ph-check-circle"></i>已评价，感谢反馈</div>
          </div>
          <NuxtLink class="nl-btn nl-btn--ghost nl-btn--block" to="/account/orders"><i class="ph ph-arrow-left"></i>返回订单列表</NuxtLink>
        </aside>
      </div>
    </div>

    <div v-if="reviewModal" class="nl-modal" role="dialog" aria-modal="true" aria-label="评价商品" @keydown.esc="reviewModal = false">
      <div class="nl-modal__scrim" @click="reviewModal = false"></div>
      <div class="nl-modal__box">
        <div class="nl-modal__head">
          <h3>评价商品</h3>
          <button type="button" class="nl-iconbtn" aria-label="关闭" @click="reviewModal = false"><i class="ph ph-x"></i></button>
        </div>
        <div class="nl-order__product">
          <span class="nl-order__thumb">
            <img v-if="order?.product?.image" :src="order.product.image" alt="" />
            <i v-else class="ph ph-package"></i>
          </span>
          <div>
            <b style="display: block">{{ order?.product?.name }}</b>
            <span class="nl-mono" style="font-size: 12px; color: var(--ink-subtle)">{{ order?.order_no }}</span>
          </div>
        </div>
        <div class="nl-field">
          <span>评分 · {{ ratingLabel }}</span>
          <div class="nl-starpick" role="radiogroup" aria-label="评分">
            <button
              v-for="n in 5"
              :key="n"
              type="button"
              role="radio"
              :aria-checked="reviewForm.rating === n ? 'true' : 'false'"
              :aria-label="`${n} 星`"
              :class="{ 'is-on': n <= reviewForm.rating }"
              @click="reviewForm.rating = n"
            >
              <i class="ph-fill ph-star"></i>
            </button>
          </div>
        </div>
        <label class="nl-field">
          <span>评价标题（选填）</span>
          <input v-model="reviewForm.title" class="nl-input" type="text" maxlength="100" placeholder="一句话概括你的体验" />
        </label>
        <label class="nl-field">
          <span>评价内容（选填）</span>
          <textarea v-model="reviewForm.content" class="nl-input" rows="4" maxlength="500" placeholder="分享你的使用体验"></textarea>
        </label>
        <div v-if="reviewError" class="nl-note nl-note--danger"><i class="ph ph-warning"></i>{{ reviewError }}</div>
        <div class="nl-modal__foot">
          <button type="button" class="nl-btn nl-btn--ghost" @click="reviewModal = false">取消</button>
          <button type="button" class="nl-btn nl-btn--primary" :disabled="!reviewForm.rating || submitting" @click="submitReview">
            {{ submitting ? "提交中…" : "提交评价" }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import { ORDER_STATUS, categoryIcon, fmtDateTime, type NlProduct } from "~/utils/nl";

definePageMeta({ layout: "layout-three" });
useHead({ title: "订单详情" });

interface CardKey {
  card_no: string;
  card_pwd?: string;
}
interface Order {
  order_no: string;
  status: number;
  quantity: number;
  total_amount: number;
  created_at: string;
  updated_at?: string;
  paid_at?: string | null;
  contact?: string;
  remark?: string;
  payment_url?: string;
  deliver_content?: string;
  deliver_note?: string;
  card_keys?: CardKey[];
  product?: NlProduct;
}

const route = useRoute();
const order = ref<Order | null>(null);
const loading = ref(true);
const confirming = ref(false);
const reviewed = ref(false);
const reviewModal = ref(false);
const submitting = ref(false);
const reviewError = ref("");
const reviewForm = ref({ rating: 5, title: "", content: "" });

const status = computed(() => ORDER_STATUS[order.value?.status ?? -1] || { label: "未知", tone: "", icon: "ph ph-question" });
const isFixed = computed(() => order.value?.product?.delivery_type === 2);
const copied = ref(false);

async function copyDeliver() {
  const text = order.value?.deliver_content || "";
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    const ta = document.createElement("textarea");
    ta.value = text;
    ta.style.position = "fixed";
    ta.style.opacity = "0";
    document.body.appendChild(ta);
    ta.select();
    document.execCommand("copy");
    ta.remove();
  }
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
}
const ratingLabel = computed(() => ["点击星星评分", "非常差", "较差", "一般", "满意", "非常满意"][reviewForm.value.rating] || "");

function ticketRows(k: CardKey) {
  const rows = [{ label: "卡号", value: k.card_no }];
  if (k.card_pwd) rows.push({ label: "卡密", value: k.card_pwd });
  return rows;
}

onMounted(async () => {
  try {
    order.value = await $fetch<Order>(`/api/orders/${route.params.orderNo}`, { credentials: "include" });
  } catch {
    order.value = null;
  } finally {
    loading.value = false;
  }
  syncPayment();
});

// 付款后如果支付平台的跳转没回到本站（关掉了页面等），订单会一直停在待支付。
// 进入订单页时向后端要一次支付状态，后端会去 NodeLoc 查询并补发货。
async function syncPayment() {
  if (order.value?.status !== 0 || !order.value.payment_url) return;
  try {
    const res = await $fetch<{ paid: boolean }>(`/api/order/${route.params.orderNo}/status`, { credentials: "include" });
    if (res.paid) {
      order.value = await $fetch<Order>(`/api/orders/${route.params.orderNo}`, { credentials: "include" });
    }
  } catch {
    // 查询失败不影响页面展示
  }
}

function openReview() {
  reviewForm.value = { rating: 5, title: "", content: "" };
  reviewError.value = "";
  reviewModal.value = true;
}

async function submitReview() {
  if (!reviewForm.value.rating || !order.value?.product?.id) return;
  submitting.value = true;
  reviewError.value = "";
  try {
    await $fetch(`/api/products/${order.value.product.id}/reviews`, {
      method: "POST",
      credentials: "include",
      body: { ...reviewForm.value },
    });
    reviewed.value = true;
    reviewModal.value = false;
  } catch (e: unknown) {
    reviewError.value = (e as { data?: { error?: string } })?.data?.error || "提交失败，请稍后重试";
  } finally {
    submitting.value = false;
  }
}

async function confirmReceipt() {
  if (!confirm("确认已收到商品？确认后将完成交易，款项将结算给卖家。")) return;
  confirming.value = true;
  try {
    order.value = await $fetch<Order>(`/api/orders/${route.params.orderNo}/confirm`, {
      method: "POST",
      credentials: "include",
    });
  } catch (e: unknown) {
    alert((e as { data?: { error?: string } })?.data?.error || "操作失败");
  } finally {
    confirming.value = false;
  }
}
</script>
