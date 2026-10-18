<template>
  <el-config-provider>
    <div class="admin-app">
      <!-- 侧边栏 -->
      <div class="admin-sidebar" :class="{ collapsed }">
        <div class="sidebar-logo">
          <el-icon size="22"><Shop /></el-icon>
          <span v-show="!collapsed" class="logo-text">管理后台</span>
        </div>

        <el-menu
          :default-active="activeMenu"
          :collapse="collapsed"
          background-color="#001529"
          text-color="rgba(255,255,255,0.65)"
          active-text-color="#ffffff"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/admin">
            <el-icon><DataAnalysis /></el-icon>
            <template #title>仪表板</template>
          </el-menu-item>
          <el-menu-item index="/admin/categories">
            <el-icon><Grid /></el-icon>
            <template #title>分类管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/products">
            <el-icon><Goods /></el-icon>
            <template #title>商品管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/card-keys">
            <el-icon><Key /></el-icon>
            <template #title>卡密管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/orders">
            <el-icon><List /></el-icon>
            <template #title>订单管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/shops">
            <el-icon><OfficeBuilding /></el-icon>
            <template #title>店铺审核</template>
          </el-menu-item>
          <el-menu-item index="/admin/withdrawals">
            <el-icon><Wallet /></el-icon>
            <template #title>提现审核</template>
          </el-menu-item>
          <el-menu-item index="/admin/blog">
            <el-icon><Document /></el-icon>
            <template #title>博客管理</template>
          </el-menu-item>
          <el-menu-item index="/admin/newsletter">
            <el-icon><Message /></el-icon>
            <template #title>邮件订阅</template>
          </el-menu-item>
          <el-menu-item index="/admin/coupons">
            <el-icon><Ticket /></el-icon>
            <template #title>优惠券</template>
          </el-menu-item>
          <el-menu-item index="/admin/settings">
            <el-icon><Setting /></el-icon>
            <template #title>系统设置</template>
          </el-menu-item>
        </el-menu>
      </div>

      <!-- 主体 -->
      <div class="admin-main">
        <!-- 顶栏 -->
        <div class="admin-header">
          <el-icon class="collapse-btn" size="20" @click="collapsed = !collapsed">
            <Fold v-if="!collapsed" />
            <Expand v-else />
          </el-icon>
          <div class="header-right">
            <span class="username">{{ authStore.user?.username }}</span>
            <el-button size="small" @click="navigateTo('/')">返回前台</el-button>
            <el-button size="small" type="danger" plain @click="authStore.logout()">退出</el-button>
          </div>
        </div>

        <!-- 内容区 -->
        <div class="admin-content">
          <slot />
        </div>
      </div>
    </div>
  </el-config-provider>
</template>

<script setup lang="ts">
const route = useRoute();
const authStore = useAuthStore();
const collapsed = ref(false);

const { apply } = useNlTheme();
onMounted(() => apply(false));

const activeMenu = computed(() => {
  const p = route.path;
  if (p === "/admin" || p === "/admin/") return "/admin";
  const subs = [
    "/admin/categories", "/admin/products", "/admin/card-keys",
    "/admin/orders", "/admin/users", "/admin/shops",
    "/admin/withdrawals", "/admin/blog", "/admin/newsletter", "/admin/coupons", "/admin/settings",
  ];
  return subs.find(s => p.startsWith(s)) ?? "/admin";
});
</script>

<style>
/* 隔离 Bootstrap 对 Element Plus 的干扰 */
.admin-app .el-menu { border-right: none !important; }
.admin-app .el-menu-item.is-active { background-color: #1677ff !important; border-radius: 6px; }
.admin-app .el-menu-item:hover { background-color: rgba(255,255,255,0.08) !important; border-radius: 6px; }
.admin-app .el-menu--collapse { width: 64px !important; }
/* 防止 Bootstrap button reset 覆盖 el-button */
.admin-app .el-button { display: inline-flex !important; box-sizing: border-box !important; }
</style>

<style scoped>
.admin-app {
  display: flex;
  height: 100vh;
  overflow: hidden;
  background: #f0f2f5;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.admin-sidebar {
  width: 220px;
  background: #001529;
  display: flex;
  flex-direction: column;
  transition: width 0.25s;
  overflow: hidden;
  flex-shrink: 0;
}
.admin-sidebar.collapsed { width: 64px; }

.sidebar-logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #fff;
  font-size: 16px;
  font-weight: 700;
  background: rgba(255,255,255,0.05);
  padding: 0 16px;
  white-space: nowrap;
  overflow: hidden;
}
.logo-text { font-size: 15px; letter-spacing: 1px; }

.sidebar-menu {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px;
  box-sizing: border-box;
}

.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.admin-header {
  height: 64px;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  box-shadow: 0 1px 4px rgba(0,21,41,0.08);
  flex-shrink: 0;
  gap: 16px;
}

.collapse-btn {
  cursor: pointer;
  color: #595959;
  transition: color 0.2s;
}
.collapse-btn:hover { color: #1677ff; }

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}
.username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.admin-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}
</style>
