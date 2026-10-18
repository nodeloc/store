<template>
  <div>
    <Breadcrumb title="商家入驻" compact />
    <div class="nl-container" style="padding-bottom: 96px">
      <section class="nl-pitch">
        <div class="nl-pitch__copy">
          <span class="nl-hero__eyebrow"><span class="nl-live"></span>商家入驻 · OPEN A SHOP</span>
          <h1 class="nl-pitch__title">在 NodeLoc 开一家<br /><em>自己的小店</em></h1>
          <p class="nl-pitch__sub">用社区账号开店，买家用积分付款，卡密自动发货。VPS、AI 账号、点卡、兑换码都能卖。</p>
          <div class="nl-hero__cta">
            <NuxtLink to="/account/shop" class="nl-btn nl-btn--volt nl-btn--lg">申请开店<i class="ph ph-arrow-right"></i></NuxtLink>
            <NuxtLink to="/vendor" class="nl-btn nl-btn--outline-ink nl-btn--lg">看看其他店铺</NuxtLink>
          </div>
        </div>
        <div v-if="shopTotal > 0" class="nl-pitch__proof">
          <div class="nl-pitch__faces">
            <span v-for="s in faces" :key="s.id" :title="s.name">
              <img v-if="s.logo" :src="s.logo" :alt="s.name" />
              <i v-else class="ph ph-storefront"></i>
            </span>
          </div>
          <b class="nl-mono">{{ shopTotal }}</b>
          <span>家店铺已入驻</span>
        </div>
      </section>

      <section class="nl-section">
        <NlSectionHead eyebrow="WHY HERE" title="为什么在这里开店" />
        <div class="nl-bento nl-perks">
          <div v-for="p in perks" :key="p.title" class="nl-bento__cell">
            <i :class="p.icon"></i>
            <b>{{ p.title }}</b>
            <span>{{ p.text }}</span>
          </div>
        </div>
      </section>

      <section class="nl-section">
        <NlSectionHead eyebrow="HOW IT WORKS" title="四步开店" />
        <ol class="nl-flow">
          <li v-for="(s, i) in flow" :key="s.title" v-reveal data-reveal :style="{ '--i': i }">
            <span class="nl-flow__no">{{ String(i + 1).padStart(2, "0") }}</span>
            <b>{{ s.title }}</b>
            <span>{{ s.text }}</span>
          </li>
        </ol>
      </section>

      <section class="nl-section">
        <NlSectionHead eyebrow="FAQ" title="常见问题" />
        <div class="nl-faq">
          <details v-for="(f, i) in faq" :key="f.q" :open="i === 0">
            <summary>{{ f.q }}<i class="ph ph-plus"></i></summary>
            <p>{{ f.a }}</p>
          </details>
        </div>
      </section>

      <section class="nl-section">
        <div class="nl-cta">
          <div>
            <b>准备好了吗？</b>
            <span>提交申请只需要一分钟，审核结果会显示在「我的店铺」页面。</span>
          </div>
          <NuxtLink to="/account/shop" class="nl-btn nl-btn--lg nl-cta__btn">立即申请<i class="ph ph-arrow-right"></i></NuxtLink>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });
useSeo({ title: "商家入驻", description: "用 NodeLoc 账号开店，积分收款，卡密自动发货。" });

interface ShopFace { id: number; name: string; logo: string }

const faces = ref<ShopFace[]>([]);
const shopTotal = ref(0);

const perks = [
  { icon: "ph ph-user-circle-check", title: "社区账号直接开店", text: "使用 NodeLoc 账号登录即可申请，不用另外注册。" },
  { icon: "ph ph-lightning", title: "积分收款", text: "买家用 NodeLoc 积分付款，订单完成后计入店铺余额。" },
  { icon: "ph ph-robot", title: "自动发货", text: "提前上传卡密库存，买家付款后系统自动发放；也支持手动发货。" },
  { icon: "ph ph-bank", title: "余额提现", text: "店铺余额可随时申请提现，进度在账户中心可查。" },
];

const flow = [
  { title: "登录账号", text: "使用 NodeLoc 社区账号登录商店。" },
  { title: "提交申请", text: "在「我的店铺」填写店铺名称、简介和联系方式。" },
  { title: "等待审核", text: "管理员审核通过后，店铺即可上线。" },
  { title: "上架开卖", text: "添加商品与库存，订单和收入都在店铺后台管理。" },
];

const faq = [
  { q: "一个账号可以开几家店？", a: "每个 NodeLoc 账号最多开一家店铺。" },
  { q: "可以卖哪些商品？", a: "以虚拟商品为主，例如 VPS、AI 账号、点卡和兑换码。请遵守 NodeLoc 社区规则，违规店铺会被下线。" },
  { q: "审核需要多久？", a: "管理员会尽快处理。审核结果和拒绝原因（如有）都会显示在「我的店铺」页面。" },
  { q: "货款怎么结算？", a: "订单完成后金额计入店铺余额，可在「余额与提现」中申请提现，手续费以申请页面显示为准。" },
];

onMounted(async () => {
  try {
    const res = await $fetch<{ data: ShopFace[]; total: number }>("/api/shops", { query: { page: 1, page_size: 8 } });
    faces.value = (res.data ?? []).slice(0, 6);
    shopTotal.value = res.total ?? 0;
  } catch {
    /* social proof is optional */
  }
});
</script>
