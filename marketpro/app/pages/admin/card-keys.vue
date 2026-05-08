<template>
  <div>
    <div class="page-header">
      <span class="page-title">卡密管理{{ productName ? ` — ${productName}` : "" }}</span>
      <el-button type="primary" :icon="Plus" @click="dialogVisible = true">批量添加卡密</el-button>
    </div>

    <el-card shadow="never" style="margin-bottom:16px">
      <el-form inline>
        <el-form-item label="商品ID">
          <el-input-number v-model="filterProductId" :min="0" placeholder="全部" style="width:160px" @change="fetchCardKeys" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchCardKeys">查询</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <el-table :data="cardKeys" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="product_id" label="商品ID" width="90" />
        <el-table-column prop="card_no" label="卡号" show-overflow-tooltip />
        <el-table-column prop="card_pwd" label="密码" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'danger' : 'success'" size="small">
              {{ row.status === 1 ? "已售出" : "未使用" }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="售出时间" width="160">
          <template #default="{ row }">
            <span v-if="row.sold_at" class="text-sm text-gray-500">
              {{ new Date(row.sold_at).toLocaleString("zh-CN") }}
            </span>
            <span v-else class="text-gray-300">—</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="90" fixed="right">
          <template #default="{ row }">
            <el-popconfirm title="确认删除？" @confirm="deleteItem(row.id)">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="fetchCardKeys" />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="批量添加卡密" width="520px">
      <el-form :model="addForm" label-width="90px">
        <el-form-item label="商品ID" required>
          <el-input-number v-model="addForm.product_id" style="width:100%" />
        </el-form-item>
        <el-form-item label="卡密内容" required>
          <el-input
            v-model="addForm.cards_text"
            type="textarea"
            :rows="10"
            placeholder="每行一条卡密，支持格式：&#10;卡号&#10;或&#10;卡号----密码"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="addCardKeys">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

definePageMeta({ layout: "admin" });

const route = useRoute();
const { get, post, del } = useApi();

const productId = computed(() => route.query.product_id ? Number(route.query.product_id) : null);
const productName = computed(() => route.query.product_name ? String(route.query.product_name) : "");

const loading = ref(true);
const saving = ref(false);
const dialogVisible = ref(false);
const cardKeys = ref<any[]>([]);
const page = ref(1);
const pageSize = ref(50);
const total = ref(0);
const filterProductId = ref<number | null>(productId.value);
const addForm = ref({ product_id: productId.value ?? (undefined as number | undefined), cards_text: "" });

async function fetchCardKeys() {
  loading.value = true;
  try {
    const pid = filterProductId.value ?? "";
    const res = await get<{ card_keys: any[]; total: number }>(`/api/admin/card-keys?product_id=${pid}&page=${page.value}&page_size=${pageSize.value}`);
    cardKeys.value = res.card_keys; total.value = res.total;
  } finally { loading.value = false; }
}

async function addCardKeys() {
  if (!addForm.value.product_id || !addForm.value.cards_text) return ElMessage.warning("请填写商品ID和卡密内容");
  saving.value = true;
  try {
    const res = await post<{ count: number }>("/api/admin/card-keys", addForm.value);
    ElMessage.success(`成功添加 ${res.count} 条卡密`);
    dialogVisible.value = false;
    addForm.value.cards_text = "";
    filterProductId.value = addForm.value.product_id ?? null;
    await fetchCardKeys();
  } catch { ElMessage.error("添加失败"); } finally { saving.value = false; }
}

async function deleteItem(id: number) {
  try { await del(`/api/admin/card-keys/${id}`); ElMessage.success("删除成功"); await fetchCardKeys(); }
  catch { ElMessage.error("删除失败"); }
}

onMounted(fetchCardKeys);
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { font-size: 20px; font-weight: 600; color: #1d2129; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
