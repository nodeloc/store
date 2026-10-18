<template>
  <div>
    <Breadcrumb title="提现管理" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-account">
        <AccountSidebar />
        <div style="min-width: 0">

            <!-- 余额 + 申请表单 -->
            <div class="border border-gray-100 rounded-8 p-32 mb-24">
              <div class="d-flex align-items-center justify-content-between mb-24 flex-wrap gap-12">
                <h5 class="mb-0">申请提现</h5>
                <div class="d-flex align-items-center gap-8 bg-main-50 rounded-8 px-20 py-10">
                  <i class="ph-fill ph-wallet text-main-600 fs-5"></i>
                  <span class="text-sm text-gray-500">可用余额</span>
                  <span class="fw-bold text-main-600 fs-5">{{ balance.toFixed(0) }}</span>
                  <span class="text-sm text-gray-400">能量</span>
                </div>
              </div>

              <!-- 未开放 -->
              <div v-if="!withdrawalEnabled" class="alert alert-warning rounded-8 d-flex align-items-center gap-12">
                <i class="ph ph-warning fs-5"></i>
                <span>提现功能暂未开放，请稍后再试。</span>
              </div>

              <!-- 申请表单 -->
              <template v-else>
                <div class="row gy-3">
                  <div class="col-sm-8 col-md-6">
                    <label class="form-label fw-medium">提现能量数量</label>
                    <div class="input-group">
                      <input
                        v-model.number="amount"
                        type="number"
                        :min="minAmount"
                        step="1"
                        class="form-control py-12 px-16 border border-gray-200 rounded-start-8"
                        placeholder="请输入能量数量"
                      />
                      <span class="input-group-text bg-gray-50 border border-gray-200 rounded-end-8 text-gray-500">能量</span>
                    </div>
                    <div class="text-sm text-gray-400 mt-6 d-flex gap-16">
                      <span>最低 {{ minAmount }} 能量</span>
                      <span v-if="feeRate > 0">手续费 {{ feeRate }}%</span>
                      <span v-if="amount > 0 && feeRate > 0" class="text-main-600">
                        实际到账 {{ actualAmount }} 能量
                      </span>
                    </div>
                  </div>
                  <div class="col-12">
                    <label class="form-label fw-medium">备注（可选）</label>
                    <input
                      v-model="remark"
                      type="text"
                      class="form-control py-12 px-16 border border-gray-200 rounded-8"
                      placeholder="如有特殊说明可填写"
                      maxlength="100"
                    />
                  </div>
                  <div class="col-12">
                    <div class="bg-gray-50 rounded-8 p-16 text-sm text-gray-500 mb-16">
                      <i class="ph ph-info me-6 text-main-600"></i>
                      提现将转账至您的 NodeLoc 账户能量，通常在管理员审核后 1-3 个工作日内到账。
                    </div>
                    <div v-if="applyMsg" :class="applyOk ? 'text-success-600' : 'text-danger-600'" class="mb-12 text-sm">
                      {{ applyMsg }}
                    </div>
                    <button type="button" class="btn btn-main px-32 rounded-8" :disabled="applying || !canApply" @click="applyWithdrawal">
                      {{ applying ? "提交中..." : "提交申请" }}
                    </button>
                    <button type="button" class="btn btn-outline-secondary px-20 rounded-8 ms-12" @click="amount = balance; remark = ''">
                      全部提现
                    </button>
                  </div>
                </div>
              </template>
            </div>

            <!-- 提现记录 -->
            <div class="border border-gray-100 rounded-8 p-24">
              <div class="d-flex justify-content-between align-items-center mb-20">
                <h6 class="mb-0">提现记录</h6>
                <span class="text-sm text-gray-400">共 {{ total }} 条</span>
              </div>

              <div v-if="loading" class="text-center text-gray-400 py-32">加载中...</div>
              <div v-else-if="withdrawals.length === 0" class="text-center text-gray-400 py-32">
                <i class="ph ph-arrow-circle-up d-block mb-8" style="font-size:2.5rem"></i>
                暂无提现记录
              </div>
              <template v-else>
                <div class="border border-gray-100 rounded-8 overflow-hidden">
                  <table class="table table-hover mb-0">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="py-12 px-16 text-sm fw-medium">申请时间</th>
                        <th class="py-12 px-16 text-sm fw-medium">申请能量</th>
                        <th class="py-12 px-16 text-sm fw-medium">手续费</th>
                        <th class="py-12 px-16 text-sm fw-medium">实际到账</th>
                        <th class="py-12 px-16 text-sm fw-medium">状态</th>
                        <th class="py-12 px-16 text-sm fw-medium">备注</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="w in withdrawals" :key="w.id">
                        <td class="py-12 px-16 text-sm text-gray-500">{{ formatDate(w.created_at) }}</td>
                        <td class="py-12 px-16 fw-semibold">{{ w.amount?.toFixed(0) }}</td>
                        <td class="py-12 px-16 text-sm text-gray-500">{{ w.fee?.toFixed(0) || "0" }}</td>
                        <td class="py-12 px-16 fw-semibold text-success-600">{{ w.actual_amount?.toFixed(0) }}</td>
                        <td class="py-12 px-16">
                          <span class="badge rounded-pill" :class="wStatusClass(w.status)">
                            {{ wStatusLabel(w.status) }}
                          </span>
                        </td>
                        <td class="py-12 px-16 text-sm text-gray-500">{{ w.remark || "-" }}</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <!-- 分页 -->
                <div v-if="total > pageSize" class="d-flex justify-content-center mt-20">
                  <nav>
                    <ul class="pagination pagination-sm mb-0">
                      <li class="page-item" :class="{ disabled: page <= 1 }">
                        <button class="page-link" @click="loadWithdrawals(page - 1)">«</button>
                      </li>
                      <li v-for="n in pageCount" :key="n" class="page-item" :class="{ active: n === page }">
                        <button class="page-link" @click="loadWithdrawals(n)">{{ n }}</button>
                      </li>
                      <li class="page-item" :class="{ disabled: page >= pageCount }">
                        <button class="page-link" @click="loadWithdrawals(page + 1)">»</button>
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

const balance = ref(0);
const withdrawalEnabled = ref(true);
const feeRate = ref(0);
const minAmount = ref(1);

const withdrawals = ref<any[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 15;
const pageCount = computed(() => Math.ceil(total.value / pageSize));
const loading = ref(true);

const applying = ref(false);
const amount = ref<number>(0);
const remark = ref("");
const applyMsg = ref("");
const applyOk = ref(false);

const actualAmount = computed(() => {
  if (!amount.value || amount.value <= 0) return 0;
  const fee = Math.round(amount.value * feeRate.value / 100);
  return amount.value - fee;
});

const canApply = computed(() =>
  withdrawalEnabled.value &&
  amount.value > 0 &&
  amount.value >= minAmount.value &&
  amount.value <= balance.value
);

onMounted(async () => {
  try {
    const res = await $fetch<any>("/api/balance", { credentials: "include" });
    balance.value = res.balance ?? 0;
    withdrawalEnabled.value = res.withdrawal_enabled !== false;
    feeRate.value = res.fee_rate ?? 0;
    minAmount.value = res.min_amount ?? 1;
  } catch { /* keep defaults */ }
  await loadWithdrawals(1);
});

async function loadWithdrawals(p = 1) {
  loading.value = true;
  page.value = p;
  try {
    const res = await $fetch<any>("/api/withdrawals", {
      credentials: "include",
      params: { page: p, page_size: pageSize },
    });
    withdrawals.value = res.data || [];
    total.value = res.total || 0;
  } catch {
    withdrawals.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

async function applyWithdrawal() {
  applyMsg.value = "";
  if (!amount.value || amount.value <= 0) {
    applyMsg.value = "请输入提现能量数量";
    applyOk.value = false;
    return;
  }
  if (amount.value < minAmount.value) {
    applyMsg.value = `最低提现 ${minAmount.value} 能量`;
    applyOk.value = false;
    return;
  }
  if (amount.value > balance.value) {
    applyMsg.value = "提现能量不能超过可用余额";
    applyOk.value = false;
    return;
  }
  applying.value = true;
  try {
    await $fetch("/api/withdrawals/apply", {
      method: "POST",
      credentials: "include",
      body: { amount: amount.value, remark: remark.value },
    });
    applyMsg.value = "申请成功，等待管理员审核后将转账至您的 NodeLoc 账户";
    applyOk.value = true;
    // 更新余额显示
    balance.value = Math.max(0, balance.value - amount.value);
    amount.value = 0;
    remark.value = "";
    await loadWithdrawals(1);
  } catch (e: any) {
    applyMsg.value = e?.data?.error || "申请失败，请稍后重试";
    applyOk.value = false;
  } finally {
    applying.value = false;
  }
}

function formatDate(d: string) {
  if (!d) return "-";
  return new Date(d).toLocaleString("zh-CN", { hour12: false });
}

const STATUS_LABEL: Record<number, string> = {
  0: "待审核",
  1: "已完成",
  2: "已拒绝",
  3: "处理中",
};
const STATUS_CLASS: Record<number, string> = {
  0: "bg-warning-100 text-warning-600",
  1: "bg-success-100 text-success-600",
  2: "bg-danger-100 text-danger-600",
  3: "bg-main-100 text-main-600",
};

function wStatusLabel(s: number) { return STATUS_LABEL[s] ?? "未知"; }
function wStatusClass(s: number) { return STATUS_CLASS[s] ?? "bg-gray-100 text-gray-500"; }
</script>
