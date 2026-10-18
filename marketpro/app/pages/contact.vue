<template>
  <div>
    <Breadcrumb title="联系我们" sub="订单、支付、开店或提现遇到问题，可以通过下面的方式找到我们。" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-contact">
        <a
          v-for="c in channels"
          :key="c.label"
          v-reveal
          data-reveal
          :href="c.href"
          :target="c.external ? '_blank' : undefined"
          :rel="c.external ? 'noopener noreferrer' : undefined"
          class="nl-contact__card"
        >
          <i :class="c.icon"></i>
          <span class="nl-contact__label">{{ c.label }}</span>
          <b>{{ c.value }}</b>
          <span class="nl-contact__hint">{{ c.hint }}</span>
        </a>
      </div>

      <div class="nl-contact__grid">
        <section class="nl-panel">
          <h2 class="nl-panel__title"><i class="ph ph-envelope-simple"></i>写邮件给我们</h2>
          <form v-if="email" class="nl-stack" style="gap: var(--space-4)" @submit.prevent="send">
            <div class="nl-contact__row">
              <label class="nl-field">
                <span>问题类型</span>
                <div class="nl-select" style="width: 100%">
                  <select v-model="topic" style="width: 100%">
                    <option v-for="t in topics" :key="t">{{ t }}</option>
                  </select>
                  <i class="ph ph-caret-down"></i>
                </div>
              </label>
              <label class="nl-field">
                <span>订单号（选填）</span>
                <input v-model.trim="orderNo" class="nl-input nl-mono" placeholder="例如 NL2026…" />
              </label>
            </div>
            <label class="nl-field">
              <span>问题描述</span>
              <textarea v-model.trim="message" class="nl-input" rows="5" placeholder="尽量写清楚发生了什么、什么时候发生的。" required></textarea>
            </label>
            <div class="nl-row" style="justify-content: space-between">
              <span class="nl-contact__note">点击后会打开你的邮件客户端，发送到 {{ email }}</span>
              <button type="submit" class="nl-btn nl-btn--volt" :disabled="!message"><i class="ph ph-paper-plane-tilt"></i>打开邮件</button>
            </div>
          </form>
          <NlEmpty v-else icon="ph ph-envelope-simple" title="暂未设置客服邮箱" text="请通过右侧的社区渠道联系我们。" />
        </section>

        <aside class="nl-panel">
          <h2 class="nl-panel__title"><i class="ph ph-lifebuoy"></i>先看看这些</h2>
          <ul class="nl-helplist">
            <li v-for="h in help" :key="h.to">
              <NuxtLink :to="h.to">
                <i :class="h.icon"></i>
                <div><b>{{ h.title }}</b><span>{{ h.text }}</span></div>
                <i class="ph ph-arrow-right nl-helplist__go"></i>
              </NuxtLink>
            </li>
          </ul>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useSeo } from "~/composables/useSeo";
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";

definePageMeta({ layout: "layout-three" });
useSeo({ title: "联系我们" });

const site = useSiteStore();
const str = (k: string) => (typeof site.settings[k] === "string" ? (site.settings[k] as string).trim() : "");

const email = computed(() => str("contact_email"));

const channels = computed(() => {
  const out: { icon: string; label: string; value: string; hint: string; href: string; external?: boolean }[] = [];
  const tel = str("contact_tel");
  const tg = tel.match(/(?:TG|Telegram)\s*[:：]?\s*@?([A-Za-z0-9_]{4,})/i);
  if (tg) {
    out.push({ icon: "ph ph-telegram-logo", label: "Telegram", value: `@${tg[1]}`, hint: str("contact_tel_label") || "响应最快", href: `https://t.me/${tg[1]}`, external: true });
  } else if (tel) {
    out.push({ icon: "ph ph-phone", label: "电话", value: tel, hint: str("contact_tel_label") || "工作时间接听", href: `tel:${tel.replace(/\s+/g, "")}` });
  }
  if (email.value) {
    out.push({ icon: "ph ph-envelope-simple", label: "邮箱", value: email.value, hint: "适合附带截图和订单号", href: `mailto:${email.value}` });
  }
  const qq = str("contact_qq");
  if (qq) out.push({ icon: "ph ph-chat-circle-dots", label: "QQ", value: qq, hint: "添加时请备注来意", href: `https://wpa.qq.com/msgrd?v=3&uin=${encodeURIComponent(qq)}`, external: true });
  const addr = str("contact_address") || "NodeLoc.com";
  const domain = addr.match(/[a-z0-9-]+\.[a-z.]{2,}/i)?.[0];
  out.push({
    icon: "ph ph-users-three",
    label: "NodeLoc 社区",
    value: addr,
    hint: "发帖或私信管理员",
    href: domain ? `https://www.${domain.replace(/^www\./i, "").toLowerCase()}` : "https://www.nodeloc.com",
    external: true,
  });
  return out;
});

const topics = ["订单问题", "支付与积分", "开店与提现", "账号问题", "其他"];
const topic = ref(topics[0]);
const orderNo = ref("");
const message = ref("");

function send() {
  if (!email.value || !message.value) return;
  const subject = `[${topic.value}]${orderNo.value ? ` 订单 ${orderNo.value}` : ""}`;
  const body = `${message.value}${orderNo.value ? `\n\n订单号：${orderNo.value}` : ""}`;
  window.location.href = `mailto:${email.value}?subject=${encodeURIComponent(subject)}&body=${encodeURIComponent(body)}`;
}

const help = [
  { to: "/account/orders", icon: "ph ph-receipt", title: "订单没收到货？", text: "在订单详情查看卡密，或联系店主。" },
  { to: "/account/withdrawals", icon: "ph ph-bank", title: "提现进度", text: "查看每笔提现的审核状态。" },
  { to: "/become-seller", icon: "ph ph-storefront", title: "想开店？", text: "了解入驻流程和常见问题。" },
  { to: "/blog", icon: "ph ph-article", title: "公告与博客", text: "平台规则和更新会发在这里。" },
];
</script>
