<template>
  <div>
    <div class="page-header">
      <span class="page-title">商品管理</span>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增商品</el-button>
    </div>

    <el-card shadow="never">
      <el-table :data="products" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="图片" width="80">
          <template #default="{ row }">
            <el-image v-if="row.image" :src="row.image" style="width:48px;height:48px;border-radius:4px;object-fit:cover" fit="cover" />
            <span v-else style="color:#ccc;font-size:12px">无图</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" show-overflow-tooltip />
        <el-table-column label="价格" width="140">
          <template #default="{ row }">
            <div>
              <span v-if="row.is_promo_active" style="color:#e6162d;font-weight:600"><i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ row.promo_price }}</span>
              <span v-else style="color:#e6162d;font-weight:600"><i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ row.price }}</span>
              <span v-if="row.orig_price || row.is_promo_active" style="color:#86909c;text-decoration:line-through;font-size:12px;margin-left:4px">
                <i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ row.is_promo_active ? row.price : row.orig_price }}
              </span>
            </div>
            <el-tag v-if="row.is_promo_active" type="danger" size="small" style="margin-top:2px">限时促销</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="库存/已售" width="110">
          <template #default="{ row }">
            <span>{{ row.stock_count }}</span>
            <span style="color:#86909c;font-size:12px"> / {{ row.sales_count }}</span>
          </template>
        </el-table-column>
        <el-table-column label="发货方式" width="110">
          <template #default="{ row }">
            <el-tag :type="row.delivery_type === 1 ? 'warning' : row.delivery_type === 2 ? 'success' : 'primary'" size="small">
              {{ row.delivery_type === 1 ? '手动发货' : row.delivery_type === 2 ? '固定内容' : '卡密自动' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort" label="排序" width="70" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'info'" size="small">{{ row.is_active ? "上架" : "下架" }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="openEdit(row)">编辑</el-button>
            <el-button v-if="(row.delivery_type ?? 0) === 0" link type="warning" size="small" @click="goCardKeys(row)">卡密</el-button>
            <el-popconfirm title="确认删除？" @confirm="deleteItem(row.id)">
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchProducts"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑商品' : '新增商品'" width="820px">
      <el-form :model="form" label-width="80px">
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="商品名称" required>
              <el-input v-model="form.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属分类" required>
              <el-select v-model="form.category_id" style="width:100%" placeholder="选择分类">
                <el-option v-for="c in allCategories" :key="c.id" :value="c.id" :label="c.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>

        <!-- 价格行：两列，每列更宽 -->
        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item label="售价" required>
              <el-input-number v-model="form.price" :min="0" :precision="2" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="原价">
              <el-input-number v-model="form.orig_price" :min="0" :precision="2" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 库存 / 已售 / 排序 -->
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="库存数量">
              <el-input-number v-model="form.stock_count" :min="0" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="已售">
              <el-input :model-value="form.sales_count" disabled style="width:100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" style="width:100%" controls-position="right" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider content-position="left" style="margin:8px 0">限时促销</el-divider>
        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="促销价">
              <el-input-number v-model="form.promo_price" :min="0" :precision="2" style="width:100%" controls-position="right" placeholder="0=不启用" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="开始时间">
              <el-date-picker v-model="form.promo_start" type="datetime" style="width:100%" placeholder="留空立即生效" value-format="YYYY-MM-DDTHH:mm:ss" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="结束时间">
              <el-date-picker v-model="form.promo_end" type="datetime" style="width:100%" placeholder="留空永久有效" value-format="YYYY-MM-DDTHH:mm:ss" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 多图上传 -->
        <el-form-item label="商品图片">
          <div class="images-manager">
            <div class="images-grid">
              <div
                v-for="(url, idx) in imageList"
                :key="idx"
                class="image-item"
                :class="{ 'is-cover': url === form.image }"
              >
                <el-image :src="url" fit="cover" class="image-thumb" />
                <div class="image-actions">
                  <el-tooltip content="设为封面" placement="top">
                    <el-icon class="action-icon cover-icon" @click="setCover(url)"><StarFilled v-if="url === form.image" /><Star v-else /></el-icon>
                  </el-tooltip>
                  <el-tooltip content="删除" placement="top">
                    <el-icon class="action-icon delete-icon" @click="removeImage(idx)"><Delete /></el-icon>
                  </el-tooltip>
                </div>
                <div v-if="url === form.image" class="cover-badge">封面</div>
              </div>

              <!-- 上传按钮 -->
              <el-upload
                class="image-upload-btn"
                action="/api/admin/upload/image"
                :show-file-list="false"
                :on-success="handleUploadSuccess"
                :on-error="handleUploadError"
                :before-upload="beforeUpload"
                accept="image/jpeg,image/png,image/gif,image/webp"
                with-credentials
                multiple
              >
                <div class="upload-placeholder">
                  <el-icon><Plus /></el-icon>
                  <span>上传图片</span>
                </div>
              </el-upload>
            </div>
            <div class="images-tip">点击 ★ 设为封面图，支持多张图片</div>
          </div>
        </el-form-item>

        <el-form-item label="发货方式">
          <el-radio-group v-model="form.delivery_type">
            <el-radio :value="0">
              <span><i class="ph ph-key" style="margin-right:4px"></i>卡密自动发货</span>
            </el-radio>
            <el-radio :value="2">
              <span><i class="ph ph-file-text" style="margin-right:4px"></i>固定内容自动发货（源码 / 教程 / 链接）</span>
            </el-radio>
            <el-radio :value="1">
              <span><i class="ph ph-hand-coins" style="margin-right:4px"></i>手动发货（TG 号 / 实物 / 其他）</span>
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="form.delivery_type === 2" label="发货内容" required>
          <el-input
            v-model="form.fixed_content"
            type="textarea"
            :rows="8"
            class="mono-input"
            placeholder="买家付款后在订单页看到的完整内容（所有买家相同，付款前不公开）"
          />
        </el-form-item>

        <el-form-item label="状态">
          <el-switch v-model="form.is_active" active-text="上架" inactive-text="下架" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Plus, Star, StarFilled, Delete } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";

definePageMeta({ layout: "admin" });

const router = useRouter();
const { get, post, put, del } = useApi();
const loading = ref(true);
const saving = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const products = ref<any[]>([]);
const allCategories = ref<any[]>([]);
const page = ref(1);
const pageSize = ref(20);
const total = ref(0);

// imageList 是当前编辑/创建时的图片数组
const imageList = ref<string[]>([]);

const defaultForm = () => ({
  name: "",
  category_id: undefined as number | undefined,
  description: "",
  price: 0,
  orig_price: 0,
  image: "",
  images: "[]",
  stock_count: 0,
  sales_count: 0,
  sort: 0,
  is_active: true,
  delivery_type: 0,
  fixed_content: "",
  promo_price: 0,
  promo_start: null as string | null,
  promo_end: null as string | null,
});
const form = ref(defaultForm());

// 同步 imageList → form.images / form.image
watch(imageList, (list) => {
  form.value.images = JSON.stringify(list);
  // 若封面不在列表中则清空
  if (form.value.image && !list.includes(form.value.image)) {
    form.value.image = list[0] ?? "";
  }
}, { deep: true });

async function fetchProducts() {
  loading.value = true;
  try {
    const res = await get<{ products: any[]; total: number }>(`/api/admin/products?page=${page.value}&page_size=${pageSize.value}`);
    products.value = res.products; total.value = res.total;
  } finally { loading.value = false; }
}

async function fetchCategories() {
  const res = await get<{ categories: any[] }>("/api/admin/categories");
  allCategories.value = res.categories;
}

function openCreate() {
  editingId.value = null;
  form.value = defaultForm();
  imageList.value = [];
  dialogVisible.value = true;
}

function openEdit(row: any) {
  editingId.value = row.id;
  form.value = {
    name: row.name,
    category_id: row.category_id,
    description: row.description ?? "",
    price: row.price,
    orig_price: row.orig_price ?? 0,
    image: row.image ?? "",
    images: row.images ?? "[]",
    stock_count: row.stock_count ?? 0,
    sales_count: row.sales_count ?? 0,
    sort: row.sort,
    is_active: row.is_active,
    delivery_type: row.delivery_type ?? 0,
    fixed_content: row.fixed_content ?? "",
    promo_price: row.promo_price ?? 0,
    promo_start: row.promo_start ? row.promo_start.slice(0, 19) : null,
    promo_end: row.promo_end ? row.promo_end.slice(0, 19) : null,
  };
  try {
    imageList.value = JSON.parse(row.images || "[]");
  } catch {
    imageList.value = [];
  }
  // 若有封面但不在列表中，加入列表
  if (form.value.image && !imageList.value.includes(form.value.image)) {
    imageList.value.unshift(form.value.image);
  }
  dialogVisible.value = true;
}

function goCardKeys(row: any) {
  router.push(`/admin/card-keys?product_id=${row.id}&product_name=${encodeURIComponent(row.name)}`);
}

function setCover(url: string) {
  form.value.image = url;
}

function removeImage(idx: number) {
  imageList.value.splice(idx, 1);
}

function beforeUpload(file: File) {
  const isImage = ["image/jpeg", "image/png", "image/gif", "image/webp"].includes(file.type);
  const isLt5M = file.size / 1024 / 1024 < 5;
  if (!isImage) { ElMessage.error("只支持 JPG、PNG、GIF、WEBP 格式"); return false; }
  if (!isLt5M) { ElMessage.error("图片大小不能超过 5MB"); return false; }
  return true;
}

function handleUploadSuccess(res: any) {
  if (res.success) {
    imageList.value.push(res.url);
    // 若没有封面则自动设第一张
    if (!form.value.image) form.value.image = res.url;
    ElMessage.success("上传成功");
  } else {
    ElMessage.error("上传失败");
  }
}

function handleUploadError() {
  ElMessage.error("图片上传失败");
}

async function save() {
  if (!form.value.name || !form.value.category_id) return ElMessage.warning("请填写商品名称和分类");
  if (form.value.delivery_type === 2 && !form.value.fixed_content.trim()) return ElMessage.warning("固定内容发货需要填写发货内容");
  saving.value = true;
  try {
    const toRFC3339 = (s: string | null) => s ? new Date(s).toISOString() : null;
    const payload = {
      ...form.value,
      images: JSON.stringify(imageList.value),
      promo_start: toRFC3339(form.value.promo_start),
      promo_end: toRFC3339(form.value.promo_end),
    };
    editingId.value
      ? await put(`/api/admin/products/${editingId.value}`, payload)
      : await post("/api/admin/products", payload);
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    await fetchProducts();
  } catch (e: any) { ElMessage.error(e?.data?.error || "保存失败"); } finally { saving.value = false; }
}

async function deleteItem(id: number) {
  try { await del(`/api/admin/products/${id}`); ElMessage.success("删除成功"); await fetchProducts(); }
  catch { ElMessage.error("删除失败"); }
}

onMounted(() => { fetchProducts(); fetchCategories(); });
</script>

<style scoped>
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.page-title { font-size: 20px; font-weight: 600; color: #1d2129; }
.pagination { display: flex; justify-content: flex-end; margin-top: 16px; }

.images-manager { width: 100%; }
.images-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: flex-start;
}

.image-item {
  position: relative;
  width: 90px;
  height: 90px;
  border-radius: 6px;
  overflow: hidden;
  border: 2px solid #e4e7ed;
  cursor: default;
}
.image-item.is-cover {
  border-color: #409eff;
}
.image-thumb {
  width: 100%;
  height: 100%;
  display: block;
}
.image-actions {
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  opacity: 0;
  transition: opacity 0.2s;
}
.image-item:hover .image-actions { opacity: 1; }
.action-icon {
  font-size: 18px;
  cursor: pointer;
  color: #fff;
}
.cover-icon:hover { color: #ffd700; }
.delete-icon:hover { color: #f56c6c; }

.cover-badge {
  position: absolute;
  bottom: 0; left: 0; right: 0;
  background: rgba(64,158,255,0.85);
  color: #fff;
  font-size: 11px;
  text-align: center;
  padding: 2px 0;
}

.image-upload-btn {
  width: 90px;
  height: 90px;
}
.image-upload-btn :deep(.el-upload) {
  width: 90px;
  height: 90px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.image-upload-btn :deep(.el-upload:hover) { border-color: #409eff; color: #409eff; }
.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #8c939d;
}
.upload-placeholder .el-icon { font-size: 22px; }

.images-tip { font-size: 12px; color: #909399; margin-top: 6px; }
.mono-input :deep(textarea) { font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: 13px; }
</style>
