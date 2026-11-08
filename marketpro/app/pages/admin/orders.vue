<template>
  <div>
    <div class="page-header">
      <span class="page-title">订单管理</span>
      <div class="d-flex gap-8">
        <el-input
          v-model="keyword"
          placeholder="搜索订单号 / 用户名"
          clearable
          style="width:200px"
          @keydown.enter="fetchOrders"
          @clear="fetchOrders"
        >
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
        <el-select v-model="statusFilter" style="width:130px" @change="fetchOrders">
          <el-option :value="-1" label="全部状态" />
          <el-option :value="0" label="待支付" />
          <el-option :value="1" label="待发货" />
          <el-option :value="2" label="已完成" />
          <el-option :value="3" label="已取消" />
          <el-option :value="4" label="已发货" />
        </el-select>
      </div>
    </div>

    <el-card shadow="never">
      <el-table :data="orders" v-loading="loading" stripe>
        <el-table-column prop="order_no" label="订单号" min-width="170" show-overflow-tooltip />

        <!-- 用户列：头像 + 用户名，可点击跳转用户管理 -->
        <el-table-column label="用户" width="160">
          <template #default="{ row }">
            <div v-if="row.user" class="user-cell" @click="goUser(row.user.id)" title="查看用户">
              <el-avatar
                :src="row.user.avatar_url || ''"
                :size="28"
                style="flex-shrink:0"
              >{{ row.user.username?.charAt(0)?.toUpperCase() }}</el-avatar>
              <div class="user-cell-info">
                <div class="user-name">{{ row.user.username }}</div>
                <div v-if="row.user.name && row.user.name !== row.user.username" class="user-realname">{{ row.user.name }}</div>
              </div>
            </div>
            <span v-else class="text-gray-400">—</span>
          </template>
        </el-table-column>

        <el-table-column label="商品" show-overflow-tooltip>
          <template #default="{ row }">{{ row.product?.name || '—' }}</template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="60" align="center" />
        <el-table-column label="金额" width="100">
          <template #default="{ row }">
            <span style="color:#e6162d;font-weight:600">{{ row.total_amount }}</span>
          </template>
        </el-table-column>
        <el-table-column label="支付方式" width="90">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ payMethodText(row.pay_method) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tooltip v-if="row.refunded_at" :content="refundTip(row)" placement="top">
              <el-tag type="info" size="small">已退款</el-tag>
            </el-tooltip>
            <el-tooltip v-else-if="row.status === 3 && row.paid_at" content="已付款但订单已取消，需要退款" placement="top">
              <el-tag type="danger" effect="dark" size="small">待退款</el-tag>
            </el-tooltip>
            <el-tag v-else :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">{{ fmtDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-dropdown v-if="ACTIONS[row.status]?.length || canRefund(row) || canSync(row)" @command="(cmd: number | string) => onCommand(row, cmd)">
              <el-button link type="primary" size="small">操作 <el-icon><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="canSync(row)" command="sync">同步支付状态</el-dropdown-item>
                  <el-dropdown-item v-if="ACTIONS[row.status]?.includes(2)" :command="2">标记完成</el-dropdown-item>
                  <el-dropdown-item v-if="canRefund(row)" command="refund" style="color:#f56c6c">退款</el-dropdown-item>
                  <el-dropdown-item v-if="ACTIONS[row.status]?.includes(3)" :command="3" style="color:#f56c6c">取消订单</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <span v-else style="color:#c0c4cc">—</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchOrders"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ArrowDown, Search } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";

definePageMeta({ layout: "admin" });

const router = useRouter();
const { get, post, put } = useApi();
const loading = ref(true);
const orders = ref<any[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);
const statusFilter = ref(-1);
const keyword = ref("");

const STATUS_TEXT: Record<number, string> = { 0: '待支付', 1: '待发货', 2: '已完成', 3: '已取消', 4: '已发货' };
const STATUS_TYPE: Record<number, string> = { 0: 'warning', 1: 'primary', 2: 'success', 3: 'danger', 4: 'info' };
// 与后端 AdminStatusTransitions 保持一致：已完成/已取消为终态
// 待支付不能直接标记已支付：付款只认 NodeLoc，用「同步支付状态」走与回调相同的发卡/结算流程
const ACTIONS: Record<number, number[]> = { 0: [3], 1: [2, 3], 4: [2, 3] };
const PAY_METHOD: Record<string, string> = { nodeloc: '能量', balance: '余额', free: '免费' };

const REFUND_METHOD: Record<string, string> = { native: '原路退回', transfer: '转账退回', balance: '退回余额' };

// 已付款、未退款、非免费的订单可退款（含付款后被取消的订单）
const canRefund = (row: any) =>
  row.status !== 0 && !!row.paid_at && !row.refunded_at && row.pay_method !== 'free' && row.total_amount > 0;
// 能量订单尚未记录付款时，可向 NodeLoc 查询（补救丢失的支付回调，或发现取消后才付款的订单）
const canSync = (row: any) =>
  (row.status === 0 || row.status === 3) && row.pay_method === 'nodeloc' && !!row.transaction_id && !row.paid_at && !row.refunded_at;
const refundTip = (row: any) =>
  `${REFUND_METHOD[row.refund_method] ?? row.refund_method ?? ''} ${row.refund_amount}` +
  (row.refunded_at ? ` · ${fmtDate(row.refunded_at)}` : '');

const statusText = (s: number) => STATUS_TEXT[s] ?? '未知';
const statusType = (s: number) => (STATUS_TYPE[s] ?? 'info') as any;
const payMethodText = (m: string) => PAY_METHOD[m] ?? (m || '—');

async function fetchOrders() {
  loading.value = true;
  try {
    const params: Record<string, string> = {
      page: String(page.value),
      page_size: String(pageSize.value),
      status: String(statusFilter.value),
    };
    if (keyword.value.trim()) params.keyword = keyword.value.trim();
    const res = await get<{ orders: any[]; total: number }>(`/api/admin/orders?${new URLSearchParams(params)}`);
    orders.value = res.orders ?? [];
    total.value = res.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function onCommand(row: any, cmd: number | string) {
  if (cmd === 'refund') refundOrder(row);
  else if (cmd === 'sync') syncPayment(row);
  else updateStatus(row, cmd as number);
}

async function syncPayment(row: any) {
  try {
    const res = await post<{ message: string }>(`/api/admin/orders/${row.order_no}/sync-payment`);
    ElMessage({ message: res.message || '已同步', type: 'info', duration: 5000 });
  } catch (e: any) {
    ElMessage.error(e?.data?.error || e?.message || "同步失败");
  }
  await fetchOrders();
}

async function refundOrder(row: any) {
  const how = row.pay_method === 'balance'
    ? '金额将退回买家站内余额。'
    : '优先原路退回买家 NodeLoc 账户；原路退款不可用时从商户账户转账退回（商户承担转账手续费）。';
  try {
    await ElMessageBox.confirm(
      `订单 ${row.order_no} 全额退款 ${row.total_amount}。${how}已结算给店主的收入会一并扣回，退款后订单变为已取消。确定退款？`,
      '订单退款',
      { type: 'warning', confirmButtonText: '确定退款', cancelButtonText: '返回' },
    );
  } catch {
    return;
  }
  try {
    const res = await post<{ order: any }>(`/api/admin/orders/${row.order_no}/refund`);
    const o = res.order ?? {};
    ElMessage.success(`退款成功：${REFUND_METHOD[o.refund_method] ?? o.refund_method ?? ''} ${o.refund_amount ?? ''}`);
  } catch (e: any) {
    ElMessage.error(e?.data?.error || e?.message || "退款失败");
  }
  await fetchOrders();
}

async function updateStatus(row: any, status: number) {
  if (status === 3) {
    const paid = canRefund(row);
    try {
      await ElMessageBox.confirm(
        paid
          ? `订单 ${row.order_no} 已付款（${row.total_amount}）。仅取消不会退款，需要退款请用「退款」。确定只取消不退款？`
          : `确定取消订单 ${row.order_no}？`,
        '取消订单',
        { type: 'warning', confirmButtonText: '确定取消', cancelButtonText: '返回' },
      );
    } catch {
      return;
    }
  }
  if (status === 2) {
    const settle = row.shop_id > 0 && !row.shop_settled;
    try {
      await ElMessageBox.confirm(
        `确定将订单 ${row.order_no} 标记为完成？` +
          (settle ? `第三方店铺订单会把 ${row.total_amount} 结算到店主余额（官方店不结算）。` : ''),
        '标记完成',
        { type: 'warning', confirmButtonText: '确定完成', cancelButtonText: '返回' },
      );
    } catch {
      return;
    }
  }
  try {
    await put(`/api/admin/orders/${row.order_no}/status`, { status });
    ElMessage.success("更新成功");
  } catch (e: any) {
    ElMessage.error(e?.data?.error || e?.message || "更新失败");
  }
  await fetchOrders();
}

function goUser(userId: number) {
  router.push(`/admin/users?highlight=${userId}`);
}

function fmtDate(d: string) {
  if (!d) return '—';
  return new Date(d).toLocaleString('zh-CN', { hour12: false });
}

onMounted(fetchOrders);
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { font-size: 20px; font-weight: 600; color: #1d2129; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 6px;
  transition: background .15s;
}
.user-cell:hover { background: #f0f5ff; }
.user-cell-info { min-width: 0; }
.user-name { font-size: 13px; font-weight: 500; color: #1d2129; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.user-realname { font-size: 11px; color: #86909c; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
</style>
