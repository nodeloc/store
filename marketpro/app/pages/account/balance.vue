<template>
  <div>
    <Breadcrumb title="余额管理" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-account">
        <AccountSidebar />
        <div style="min-width: 0">

            <!-- 余额卡片 -->
            <div class="border border-gray-100 rounded-8 p-32 mb-24">
              <h5 class="mb-20">账户余额</h5>
              <div v-if="loading" class="text-gray-400">加载中...</div>
              <div v-else class="d-flex align-items-center gap-16">
                <i class="ph-fill ph-wallet text-main-600" style="font-size:2.5rem"></i>
                <div>
                  <div class="text-gray-500 text-sm mb-4">可用余额（能量）</div>
                  <div class="fw-bold text-main-600" style="font-size:2rem">
                    {{ balance?.balance?.toFixed(0) ?? "0" }}
                  </div>
                </div>
              </div>
            </div>

            <!-- 交易记录 -->
            <div class="border border-gray-100 rounded-8 p-24">
              <div class="d-flex justify-content-between align-items-center mb-20">
                <h6 class="mb-0">余额变动记录</h6>
                <span class="text-sm text-gray-400">共 {{ total }} 条</span>
              </div>

              <div v-if="txLoading" class="text-center py-32 text-gray-400">加载中...</div>
              <div v-else-if="txs.length === 0" class="text-center py-32 text-gray-400">
                <i class="ph ph-receipt d-block mb-8" style="font-size:2.5rem"></i>
                暂无记录
              </div>
              <template v-else>
                <div class="border border-gray-100 rounded-8 overflow-hidden">
                  <table class="table table-hover mb-0">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="py-12 px-16 text-sm fw-medium">时间</th>
                        <th class="py-12 px-16 text-sm fw-medium">类型</th>
                        <th class="py-12 px-16 text-sm fw-medium">金额</th>
                        <th class="py-12 px-16 text-sm fw-medium">变动后余额</th>
                        <th class="py-12 px-16 text-sm fw-medium">说明</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="tx in txs" :key="tx.id">
                        <td class="py-12 px-16 text-sm text-gray-500">{{ formatDate(tx.created_at) }}</td>
                        <td class="py-12 px-16">
                          <span class="badge rounded-pill" :class="txTypeBadge(tx.type)">
                            {{ txTypeLabel(tx.type) }}
                          </span>
                        </td>
                        <td class="py-12 px-16 fw-semibold"
                          :class="tx.amount >= 0 ? 'text-success-600' : 'text-danger-600'">
                          {{ tx.amount >= 0 ? "+" : "" }}{{ tx.amount?.toFixed(0) }}
                        </td>
                        <td class="py-12 px-16 text-sm text-gray-600">
                          {{ tx.balance_after?.toFixed(0) }}
                        </td>
                        <td class="py-12 px-16 text-sm text-gray-500">{{ tx.description || "-" }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <!-- 分页 -->
                <div v-if="total > pageSize" class="d-flex justify-content-center mt-20">
                  <nav>
                    <ul class="pagination pagination-sm mb-0">
                      <li class="page-item" :class="{ disabled: page <= 1 }">
                        <button class="page-link" @click="fetchTxs(page - 1)">«</button>
                      </li>
                      <li v-for="n in pageCount" :key="n" class="page-item" :class="{ active: n === page }">
                        <button class="page-link" @click="fetchTxs(n)">{{ n }}</button>
                      </li>
                      <li class="page-item" :class="{ disabled: page >= pageCount }">
                        <button class="page-link" @click="fetchTxs(page + 1)">»</button>
                      </li>
                    </ul>
                  </nav>
                </div>
              </template>
            </div>

        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import AccountSidebar from "~/components/containers/account/AccountSidebar.vue";

definePageMeta({ layout: "layout-three" });

const balance = ref<any>(null);
const txs = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 20;
const pageCount = computed(() => Math.ceil(total.value / pageSize));
const loading = ref(true);
const txLoading = ref(false);

onMounted(async () => {
  try {
    balance.value = await $fetch<any>("/api/balance", { credentials: "include" });
  } catch {
    balance.value = null;
  } finally {
    loading.value = false;
  }
  await fetchTxs(1);
});

async function fetchTxs(p: number) {
  txLoading.value = true;
  page.value = p;
  try {
    const res = await $fetch<any>("/api/balance/txs", {
      credentials: "include",
      params: { page: p, page_size: pageSize },
    });
    txs.value = res.data || [];
    total.value = res.total || 0;
  } catch {
    txs.value = [];
    total.value = 0;
  } finally {
    txLoading.value = false;
  }
}

function formatDate(d: string) {
  if (!d) return "-";
  return new Date(d).toLocaleString("zh-CN", { hour12: false });
}

const TX_LABELS: Record<string, string> = {
  sale_income:       "销售收入",
  withdrawal:        "提现",
  withdrawal_refund: "提现退款",
  purchase:          "余额支付",
  admin_adjust:      "后台调整",
};

const TX_BADGE: Record<string, string> = {
  sale_income:       "bg-success-100 text-success-600",
  withdrawal:        "bg-warning-100 text-warning-600",
  withdrawal_refund: "bg-main-100 text-main-600",
  purchase:          "bg-danger-100 text-danger-600",
  admin_adjust:      "bg-gray-100 text-gray-600",
};

function txTypeLabel(type: string) {
  return TX_LABELS[type] ?? type;
}
function txTypeBadge(type: string) {
  return TX_BADGE[type] ?? "bg-gray-100 text-gray-600";
}
</script>
