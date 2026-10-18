<template>
  <div>
    <Breadcrumb title="我的店铺" />
    <div class="nl-container" style="padding-bottom: 96px">
      <div class="nl-account">
        <AccountSidebar />
        <div style="min-width: 0">

            <!-- 加载中 -->
            <div v-if="loading" class="text-center py-48 text-gray-400">加载中...</div>

            <!-- 无店铺 -->
            <template v-else-if="!shop">
              <!-- 已提交申请 -->
              <div v-if="applySuccess" class="border border-gray-100 rounded-8 p-48 text-center">
                <i class="ph ph-clock text-warning-600" style="font-size:4rem"></i>
                <h6 class="mt-16 mb-8 text-gray-600">申请已提交，等待审核</h6>
                <p class="text-gray-400 text-sm">管理员审核通过后，您即可开始运营店铺。</p>
              </div>
              <!-- 申请表单 -->
              <div v-else class="border border-gray-100 rounded-8 p-32">
                <h5 class="mb-24 fw-semibold">申请开店</h5>
                <form @submit.prevent="applyShop">
                  <div class="row gy-3">
                    <div class="col-12">
                      <label class="form-label fw-medium">店铺名称 <span class="text-danger-600">*</span></label>
                      <input v-model="applyForm.name" type="text" class="form-control py-12 px-16 border border-gray-200 rounded-8" placeholder="请输入店铺名称" required />
                    </div>
                    <div class="col-12">
                      <label class="form-label fw-medium">联系方式</label>
                      <input v-model="applyForm.contact" type="text" class="form-control py-12 px-16 border border-gray-200 rounded-8" placeholder="QQ / 微信 / 邮箱等" />
                    </div>
                    <div class="col-12">
                      <label class="form-label fw-medium">店铺简介</label>
                      <textarea v-model="applyForm.description" class="form-control py-12 px-16 border border-gray-200 rounded-8" rows="3" placeholder="简单介绍一下您的店铺"></textarea>
                    </div>
                    <div class="col-12">
                      <div v-if="applyError" class="mb-8 text-sm text-danger-600">{{ applyError }}</div>
                      <button type="submit" class="btn btn-main px-40 rounded-8" :disabled="applying">
                        {{ applying ? "提交中..." : "提交申请" }}
                      </button>
                    </div>
                  </div>
                </form>
              </div>
            </template>

            <!-- 有店铺 -->
            <template v-else>
              <!-- 状态横幅 -->
              <div v-if="shop.status === 0" class="alert alert-warning d-flex align-items-center gap-12 mb-24 rounded-8">
                <i class="ph ph-clock fs-5"></i>
                <span>您的店铺正在审核中，请耐心等待。</span>
              </div>
              <div v-else-if="shop.status === 2" class="alert alert-danger d-flex align-items-center gap-12 mb-24 rounded-8">
                <i class="ph ph-x-circle fs-5"></i>
                <div>
                  <div class="fw-semibold">店铺申请被拒绝</div>
                  <div v-if="shop.reject_reason" class="text-sm mt-4">原因：{{ shop.reject_reason }}</div>
                </div>
              </div>
              <div v-else-if="shop.status === 3" class="alert alert-danger d-flex align-items-center gap-12 mb-24 rounded-8">
                <i class="ph ph-prohibit fs-5"></i>
                <span>您的店铺已被封禁。如有疑问请联系管理员。</span>
              </div>

              <!-- Tab 导航 -->
              <ul class="nav nav-tabs mb-24">
                <li v-if="shop.status === 1" class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'dashboard' }" @click="switchTab('dashboard')">Dashboard</button>
                </li>
                <li class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'info' }" @click="tab = 'info'">店铺信息</button>
                </li>
                <li v-if="shop.status === 1" class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'products' }" @click="switchTab('products')">商品管理</button>
                </li>
                <li v-if="shop.status === 1" class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'orders' }" @click="switchTab('orders')">店铺订单</button>
                </li>
                <li v-if="shop.status === 1" class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'reviews' }" @click="switchTab('reviews')">评价管理</button>
                </li>
                <li v-if="shop.status === 1" class="nav-item">
                  <button class="nav-link" :class="{ active: tab === 'coupons' }" @click="switchTab('coupons')">优惠券</button>
                </li>
              </ul>

              <!-- Tab：Dashboard -->
              <div v-if="tab === 'dashboard'">
                <div v-if="dashLoading" class="text-center py-48 text-gray-400">加载中...</div>
                <template v-else>
                  <!-- 核心指标卡片 -->
                  <div class="row g-3 mb-24">
                    <div class="col-sm-6 col-xl-3">
                      <div class="border border-gray-100 rounded-8 p-20 text-center">
                        <div class="text-gray-400 text-sm mb-8">今日订单</div>
                        <div class="fs-3 fw-bold text-main-600">{{ dashData.order_stats?.today_orders ?? 0 }}</div>
                        <div class="text-sm text-gray-400 mt-4">收入 <span class="text-danger-600 fw-semibold">{{ dashData.order_stats?.today_income?.toFixed(0) ?? 0 }}</span> 能量</div>
                      </div>
                    </div>
                    <div class="col-sm-6 col-xl-3">
                      <div class="border border-gray-100 rounded-8 p-20 text-center">
                        <div class="text-gray-400 text-sm mb-8">本周订单</div>
                        <div class="fs-3 fw-bold text-success-600">{{ dashData.order_stats?.week_orders ?? 0 }}</div>
                        <div class="text-sm text-gray-400 mt-4">收入 <span class="text-danger-600 fw-semibold">{{ dashData.order_stats?.week_income?.toFixed(0) ?? 0 }}</span> 能量</div>
                      </div>
                    </div>
                    <div class="col-sm-6 col-xl-3">
                      <div class="border border-gray-100 rounded-8 p-20 text-center">
                        <div class="text-gray-400 text-sm mb-8">本月订单</div>
                        <div class="fs-3 fw-bold text-warning-600">{{ dashData.order_stats?.month_orders ?? 0 }}</div>
                        <div class="text-sm text-gray-400 mt-4">收入 <span class="text-danger-600 fw-semibold">{{ dashData.order_stats?.month_income?.toFixed(0) ?? 0 }}</span> 能量</div>
                      </div>
                    </div>
                    <div class="col-sm-6 col-xl-3">
                      <div class="border border-gray-100 rounded-8 p-20 text-center">
                        <div class="text-gray-400 text-sm mb-8">总收入</div>
                        <div class="fs-3 fw-bold text-danger-600">{{ dashData.order_stats?.total_income?.toFixed(0) ?? 0 }}</div>
                        <div class="text-sm text-gray-400 mt-4">能量</div>
                      </div>
                    </div>
                  </div>

                  <!-- 订单概览 + 评价概览 -->
                  <div class="row g-3 mb-24">
                    <div class="col-lg-6">
                      <div class="border border-gray-100 rounded-8 p-24">
                        <h6 class="fw-semibold mb-16">订单概览</h6>
                        <div class="d-flex flex-column gap-12">
                          <div class="d-flex justify-content-between align-items-center">
                            <span class="text-gray-500 text-sm">总订单数</span>
                            <span class="fw-semibold">{{ dashData.order_stats?.total_orders ?? 0 }}</span>
                          </div>
                          <div class="d-flex justify-content-between align-items-center">
                            <span class="text-gray-500 text-sm">已完成</span>
                            <span class="fw-semibold text-success-600">{{ dashData.order_stats?.completed_orders ?? 0 }}</span>
                          </div>
                          <div class="d-flex justify-content-between align-items-center">
                            <span class="text-gray-500 text-sm">待支付</span>
                            <span class="fw-semibold text-warning-600">{{ dashData.order_stats?.pending_orders ?? 0 }}</span>
                          </div>
                          <div class="d-flex justify-content-between align-items-center">
                            <span class="text-gray-500 text-sm">已取消</span>
                            <span class="fw-semibold text-gray-400">{{ dashData.order_stats?.cancelled_orders ?? 0 }}</span>
                          </div>
                          <div class="d-flex justify-content-between align-items-center">
                            <span class="text-gray-500 text-sm">商品数量</span>
                            <span class="fw-semibold">{{ dashData.product_count ?? 0 }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                    <div class="col-lg-6">
                      <div class="border border-gray-100 rounded-8 p-24">
                        <h6 class="fw-semibold mb-16">评价概览</h6>
                        <div v-if="(dashData.review_stats?.total ?? 0) === 0" class="text-center text-gray-400 py-20 text-sm">暂无评价</div>
                        <template v-else>
                          <div class="d-flex align-items-center gap-16 mb-16">
                            <div class="text-center">
                              <div class="fs-2 fw-bold text-warning-600">{{ dashData.review_stats?.avg_rating?.toFixed(1) }}</div>
                              <div class="text-sm text-gray-400">平均评分</div>
                            </div>
                            <div class="text-center">
                              <div class="fs-2 fw-bold">{{ dashData.review_stats?.total }}</div>
                              <div class="text-sm text-gray-400">总评价数</div>
                            </div>
                          </div>
                          <div class="d-flex flex-column gap-8">
                            <div v-for="star in [5,4,3,2,1]" :key="star" class="d-flex align-items-center gap-8">
                              <span class="text-sm text-warning-600" style="width:40px">{{ star }} 星</span>
                              <div class="flex-grow-1 bg-gray-100 rounded-pill" style="height:8px">
                                <div class="bg-warning-600 rounded-pill h-100" :style="{ width: ratingPercent(star) + '%' }"></div>
                              </div>
                              <span class="text-sm text-gray-400" style="width:30px;text-align:right">{{ dashData.review_stats?.rating_distribution?.[star] ?? 0 }}</span>
                            </div>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>

                  <!-- 最近订单 -->
                  <div class="border border-gray-100 rounded-8 p-24">
                    <div class="d-flex justify-content-between align-items-center mb-16">
                      <h6 class="fw-semibold mb-0">最近订单</h6>
                      <button class="btn btn-sm btn-outline-main rounded-8" @click="switchTab('orders')">查看全部</button>
                    </div>
                    <div v-if="!dashData.order_stats?.recent_orders?.length" class="text-center text-gray-400 py-16 text-sm">暂无订单</div>
                    <div v-else class="overflow-hidden">
                      <table class="table table-hover mb-0">
                        <thead class="bg-gray-50">
                          <tr>
                            <th class="py-10 px-12 text-sm fw-medium">订单号</th>
                            <th class="py-10 px-12 text-sm fw-medium">商品</th>
                            <th class="py-10 px-12 text-sm fw-medium">金额</th>
                            <th class="py-10 px-12 text-sm fw-medium">状态</th>
                            <th class="py-10 px-12 text-sm fw-medium">时间</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="o in dashData.order_stats?.recent_orders" :key="o.id">
                            <td class="py-10 px-12 text-sm font-monospace">{{ o.order_no }}</td>
                            <td class="py-10 px-12 text-sm" style="max-width:140px">{{ o.product?.name }}</td>
                            <td class="py-10 px-12 fw-semibold text-danger-600">{{ o.total_amount }}</td>
                            <td class="py-10 px-12"><span class="badge" :class="orderStatusClass(o.status)">{{ orderStatusText(o.status) }}</span></td>
                            <td class="py-10 px-12 text-sm text-gray-400">{{ fmtDate(o.created_at) }}</td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  </div>
                </template>
              </div>

              <!-- Tab：店铺信息 -->
              <div v-if="tab === 'info'" class="border border-gray-100 rounded-8 p-32">
                <div class="d-flex align-items-center gap-20 mb-28">
                  <div class="position-relative">
                    <img v-if="infoForm.logo" :src="infoForm.logo" class="rounded-circle" style="width:72px;height:72px;object-fit:cover" alt="店铺头像" />
                    <div v-else class="rounded-circle bg-gray-100 d-flex align-items-center justify-content-center" style="width:72px;height:72px">
                      <i class="ph ph-storefront text-gray-400 fs-3"></i>
                    </div>
                    <label class="position-absolute bottom-0 end-0 bg-white border border-gray-200 rounded-circle d-flex align-items-center justify-content-center" style="width:24px;height:24px;cursor:pointer">
                      <i class="ph ph-camera" style="font-size:12px"></i>
                      <input type="file" accept="image/*" class="d-none" @change="uploadLogo" />
                    </label>
                  </div>
                  <div>
                    <h5 class="mb-2">{{ shop.name }}</h5>
                    <span class="badge" :class="statusBadge">{{ statusText }}</span>
                  </div>
                </div>
                <form @submit.prevent="saveShop">
                  <div class="row gy-3">
                    <div class="col-12">
                      <label class="form-label fw-medium">店铺名称</label>
                      <input v-model="infoForm.name" type="text" class="form-control py-12 px-16 border border-gray-200 rounded-8" />
                    </div>
                    <div class="col-12">
                      <label class="form-label fw-medium">联系方式</label>
                      <input v-model="infoForm.contact" type="text" class="form-control py-12 px-16 border border-gray-200 rounded-8" />
                    </div>
                    <div class="col-12">
                      <label class="form-label fw-medium">店铺简介</label>
                      <textarea v-model="infoForm.description" class="form-control py-12 px-16 border border-gray-200 rounded-8" rows="3"></textarea>
                    </div>

                    <!-- 特色服务 -->
                    <div class="col-12">
                      <div class="d-flex justify-content-between align-items-center mb-4">
                        <div>
                          <div class="fw-medium mb-2">特色服务</div>
                          <div class="text-sm text-gray-400">显示在商品详情页右侧，最多 6 条</div>
                        </div>
                        <button type="button" class="btn btn-main btn-sm rounded-8 px-16"
                          style="white-space:nowrap"
                          :disabled="infoForm.features.length >= 6"
                          @click="addFeature">
                          <i class="ph ph-plus me-4"></i>添加条目
                        </button>
                      </div>
                      <div v-if="infoForm.features.length === 0" class="border border-dashed border-gray-200 rounded-8 p-20 text-center text-gray-400 text-sm mt-12">
                        暂未设置，将显示系统默认内容
                      </div>
                      <div v-else class="d-flex flex-column gap-12 mt-12">
                        <div v-for="(f, fi) in infoForm.features" :key="fi"
                          class="border border-gray-100 rounded-8 p-16 d-flex gap-12 align-items-start">
                          <!-- 图标选择 -->
                          <div class="flex-shrink-0" style="width:44px">
                            <div class="w-44 h-44 bg-main-50 text-main-600 rounded-circle d-flex align-items-center justify-content-center fs-5 mb-6">
                              <i :class="f.icon || 'ph-fill ph-star'"></i>
                            </div>
                            <select v-model="f.icon" class="form-select form-select-sm p-0 border-0 text-xs text-gray-400" style="font-size:10px">
                              <option v-for="ic in featureIcons" :key="ic.value" :value="ic.value">{{ ic.label }}</option>
                            </select>
                          </div>
                          <!-- 标题 + 描述 -->
                          <div class="flex-grow-1 d-flex flex-column gap-8">
                            <input v-model="f.title" type="text" class="form-control form-control-sm py-8 px-12 border border-gray-200 rounded-6"
                              placeholder="标题，如：正版授权" maxlength="20" />
                            <input v-model="f.text" type="text" class="form-control form-control-sm py-8 px-12 border border-gray-200 rounded-6"
                              placeholder="描述，如：所有商品均为正版授权" maxlength="60" />
                          </div>
                          <!-- 删除 -->
                          <button type="button" class="btn btn-danger btn-sm flex-shrink-0"
                            style="padding:6px 10px"
                            @click="infoForm.features.splice(fi, 1)">
                            <i class="ph ph-trash"></i>
                          </button>
                        </div>
                      </div>
                    </div>

                    <div class="col-12">
                      <div v-if="saveMsg" :class="saveOk ? 'text-success-600' : 'text-danger-600'" class="mb-8 text-sm">{{ saveMsg }}</div>
                      <button type="submit" class="btn btn-main px-32 rounded-8" :disabled="saving">
                        {{ saving ? "保存中..." : "保存修改" }}
                      </button>
                      <span v-if="shop.status === 1" class="ms-12 text-sm text-gray-400">修改后需重新审核</span>
                    </div>
                  </div>
                </form>
              </div>

              <!-- Tab：商品管理 -->
              <div v-if="tab === 'products'">
                <div class="d-flex justify-content-between align-items-center mb-16">
                  <h6 class="mb-0">商品列表（{{ products.length }}）</h6>
                  <button class="btn btn-main btn-sm px-20 rounded-8" @click="openProductCreate">
                    <i class="ph ph-plus me-4"></i>新增商品
                  </button>
                </div>

                <div v-if="products.length === 0" class="border border-gray-100 rounded-8 p-32 text-center text-gray-400">暂无商品</div>
                <div v-else class="border border-gray-100 rounded-8 overflow-hidden">
                  <table class="table table-hover mb-0">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="py-12 px-16 text-sm fw-medium">商品</th>
                        <th class="py-12 px-16 text-sm fw-medium">价格</th>
                        <th class="py-12 px-16 text-sm fw-medium">库存/已售</th>
                        <th class="py-12 px-16 text-sm fw-medium">状态</th>
                        <th class="py-12 px-16 text-sm fw-medium">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="p in products" :key="p.id">
                        <td class="py-12 px-16">
                          <div class="d-flex align-items-center gap-12">
                            <img v-if="p.image" :src="p.image" class="rounded-4" style="width:40px;height:40px;object-fit:cover" :alt="p.name" />
                            <div class="fw-medium text-sm" style="max-width:180px">{{ p.name }}</div>
                          </div>
                        </td>
                        <td class="py-12 px-16">
                          <div class="text-danger-600 fw-semibold"><i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ p.is_promo_active ? p.promo_price : p.price }}</div>
                          <div v-if="p.is_promo_active" class="text-gray-400 text-xs text-decoration-line-through"><i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ p.price }}</div>
                        </td>
                        <td class="py-12 px-16 text-sm">{{ p.stock_count }} / {{ p.sales_count }}</td>
                        <td class="py-12 px-16">
                          <span class="badge" :class="p.is_active ? 'bg-success-100 text-success-600' : 'bg-gray-100 text-gray-500'">
                            {{ p.is_active ? "上架" : "下架" }}
                          </span>
                        </td>
                        <td class="py-12 px-16">
                          <button class="btn btn-sm btn-primary rounded-6 me-6" @click="openProductEdit(p)">编辑</button>
                          <button class="btn btn-sm btn-warning rounded-6 me-6" @click="openProductStock(p)">库存</button>
                          <button class="btn btn-sm btn-danger rounded-6" @click="deleteProduct(p.id)">删除</button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              <!-- Tab：店铺订单 -->
              <div v-if="tab === 'orders'">
                <div class="d-flex justify-content-between align-items-center mb-16">
                  <h6 class="mb-0">订单列表（共 {{ orderTotal }} 条）</h6>
                  <div class="d-flex gap-8">
                    <select v-model="orderStatus" class="form-select form-select-sm rounded-8" style="width:120px" @change="fetchOrders(1)">
                      <option value="-1">全部状态</option>
                      <option value="0">待支付</option>
                      <option value="1">已支付</option>
                      <option value="2">已完成</option>
                      <option value="3">已取消</option>
                    </select>
                  </div>
                </div>
                <div v-if="orders.length === 0" class="border border-gray-100 rounded-8 p-32 text-center text-gray-400">暂无订单</div>
                <div v-else class="border border-gray-100 rounded-8 overflow-hidden">
                  <table class="table table-hover mb-0">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="py-12 px-16 text-sm fw-medium">订单号</th>
                        <th class="py-12 px-16 text-sm fw-medium">商品</th>
                        <th class="py-12 px-16 text-sm fw-medium">买家联系</th>
                        <th class="py-12 px-16 text-sm fw-medium">金额</th>
                        <th class="py-12 px-16 text-sm fw-medium">状态</th>
                        <th class="py-12 px-16 text-sm fw-medium">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="o in orders" :key="o.id">
                        <td class="py-12 px-16 text-sm font-monospace">{{ o.order_no }}</td>
                        <td class="py-12 px-16 text-sm" style="max-width:140px">{{ o.product?.name }}</td>
                        <td class="py-12 px-16 text-sm text-gray-500">{{ o.contact || o.user?.username || "-" }}</td>
                        <td class="py-12 px-16 fw-semibold text-danger-600"><i class="ph-fill ph-lightning-a" style="font-size:0.85em;vertical-align:middle;margin-right:1px"></i>{{ o.total_amount }}</td>
                        <td class="py-12 px-16">
                          <span class="badge" :class="orderStatusClass(o.status)">{{ orderStatusText(o.status) }}</span>
                        </td>
                        <td class="py-12 px-16">
                          <!-- 待发货：显示发货按钮 -->
                          <button
                            v-if="o.status === 1 && o.product?.delivery_type === 1"
                            class="btn btn-sm btn-warning rounded-6"
                            @click="openShipDialog(o)"
                          >发货</button>
                          <!-- 已发货：显示已发货内容 -->
                          <span v-else-if="o.status === 4" class="text-sm text-gray-400">
                            已发货，等待确认
                          </span>
                          <span v-else class="text-sm text-gray-300">-</span>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
                <!-- 分页 -->
                <div v-if="orderTotal > orderPageSize" class="d-flex justify-content-center mt-20">
                  <nav>
                    <ul class="pagination pagination-sm mb-0">
                      <li class="page-item" :class="{ disabled: orderPage <= 1 }">
                        <button class="page-link" @click="fetchOrders(orderPage - 1)">«</button>
                      </li>
                      <li v-for="n in orderPageCount" :key="n" class="page-item" :class="{ active: n === orderPage }">
                        <button class="page-link" @click="fetchOrders(n)">{{ n }}</button>
                      </li>
                      <li class="page-item" :class="{ disabled: orderPage >= orderPageCount }">
                        <button class="page-link" @click="fetchOrders(orderPage + 1)">»</button>
                      </li>
                    </ul>
                  </nav>
                </div>
              </div>

              <!-- Tab：评价管理 -->
              <div v-if="tab === 'reviews'">
                <div v-if="reviewLoading" class="text-center py-48 text-gray-400">加载中...</div>
                <template v-else>
                  <!-- 评价统计 -->
                  <div class="border border-gray-100 rounded-8 p-24 mb-24">
                    <div class="d-flex align-items-center gap-24 flex-wrap">
                      <div class="d-flex align-items-center gap-12">
                        <div class="text-center">
                          <div class="fs-2 fw-bold text-warning-600">{{ reviewStats?.avg_rating?.toFixed(1) ?? '0.0' }}</div>
                          <div class="text-sm text-gray-400">平均评分</div>
                        </div>
                        <div class="d-flex gap-2">
                          <i v-for="s in 5" :key="s" class="ph-fill ph-star" :class="s <= Math.round(reviewStats?.avg_rating ?? 0) ? 'text-warning-600' : 'text-gray-200'" style="font-size:18px"></i>
                        </div>
                      </div>
                      <div class="vr d-none d-sm-block"></div>
                      <div class="text-center">
                        <div class="fs-4 fw-bold">{{ reviewStats?.total ?? 0 }}</div>
                        <div class="text-sm text-gray-400">总评价数</div>
                      </div>
                      <div class="vr d-none d-sm-block"></div>
                      <div class="d-flex gap-8 flex-wrap">
                        <button class="btn btn-sm rounded-pill px-16" :class="reviewRating === 0 ? 'btn-main' : 'btn-outline-main'" @click="filterReviews(0)">全部</button>
                        <button v-for="star in [5,4,3,2,1]" :key="star" class="btn btn-sm rounded-pill px-16"
                          :class="reviewRating === star ? 'btn-main' : 'btn-outline-main'"
                          @click="filterReviews(star)">
                          {{ star }}星 ({{ reviewStats?.rating_distribution?.[star] ?? 0 }})
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- 评价列表 -->
                  <div v-if="reviews.length === 0" class="border border-gray-100 rounded-8 p-32 text-center text-gray-400">暂无评价</div>
                  <div v-else class="d-flex flex-column gap-16">
                    <div v-for="r in reviews" :key="r.id" class="border border-gray-100 rounded-8 p-20">
                      <div class="d-flex justify-content-between align-items-start mb-12">
                        <div class="d-flex align-items-center gap-12">
                          <img v-if="r.user?.avatar_url" :src="r.user.avatar_url" class="rounded-circle" style="width:36px;height:36px;object-fit:cover" :alt="r.user?.username" />
                          <div v-else class="rounded-circle bg-gray-100 d-flex align-items-center justify-content-center" style="width:36px;height:36px">
                            <i class="ph ph-user text-gray-400"></i>
                          </div>
                          <div>
                            <div class="fw-medium text-sm">{{ r.user?.username || '匿名用户' }}</div>
                            <div class="d-flex gap-2 mt-2">
                              <i v-for="s in 5" :key="s" class="ph-fill ph-star" :class="s <= r.rating ? 'text-warning-600' : 'text-gray-200'" style="font-size:13px"></i>
                            </div>
                          </div>
                        </div>
                        <span class="text-sm text-gray-400">{{ fmtDate(r.created_at) }}</span>
                      </div>
                      <div class="bg-gray-50 rounded-6 px-12 py-8 mb-10 text-sm text-gray-500">
                        <i class="ph ph-package me-4"></i>{{ r.product?.name || '未知商品' }}
                      </div>
                      <div v-if="r.title" class="fw-semibold text-sm mb-4">{{ r.title }}</div>
                      <div class="text-sm text-gray-600">{{ r.content || '用户未填写评价内容' }}</div>
                    </div>
                  </div>

                  <!-- 分页 -->
                  <div v-if="reviewTotal > reviewPageSize" class="d-flex justify-content-center mt-20">
                    <nav>
                      <ul class="pagination pagination-sm mb-0">
                        <li class="page-item" :class="{ disabled: reviewPage <= 1 }">
                          <button class="page-link" @click="fetchReviews(reviewPage - 1)">«</button>
                        </li>
                        <li v-for="n in reviewPageCount" :key="n" class="page-item" :class="{ active: n === reviewPage }">
                          <button class="page-link" @click="fetchReviews(n)">{{ n }}</button>
                        </li>
                        <li class="page-item" :class="{ disabled: reviewPage >= reviewPageCount }">
                          <button class="page-link" @click="fetchReviews(reviewPage + 1)">»</button>
                        </li>
                      </ul>
                    </nav>
                  </div>
                </template>
              </div>

              <!-- Tab：优惠券 -->
              <div v-if="tab === 'coupons'">
                <div class="d-flex justify-content-between align-items-center mb-16">
                  <h6 class="mb-0">店铺优惠券</h6>
                  <button class="btn btn-main btn-sm px-20 rounded-8" @click="openCouponCreate">
                    <i class="ph ph-plus me-4"></i>新建优惠券
                  </button>
                </div>
                <div v-if="couponLoading" class="text-center py-32 text-gray-400">加载中...</div>
                <div v-else-if="coupons.length === 0" class="border border-dashed border-gray-200 rounded-8 p-32 text-center text-gray-400">
                  暂无优惠券
                </div>
                <div v-else class="border border-gray-100 rounded-8 overflow-hidden">
                  <table class="table table-hover mb-0">
                    <thead class="bg-gray-50">
                      <tr>
                        <th class="py-12 px-16 text-sm fw-medium">券码</th>
                        <th class="py-12 px-16 text-sm fw-medium">折扣</th>
                        <th class="py-12 px-16 text-sm fw-medium">最低金额</th>
                        <th class="py-12 px-16 text-sm fw-medium">使用次数</th>
                        <th class="py-12 px-16 text-sm fw-medium">过期时间</th>
                        <th class="py-12 px-16 text-sm fw-medium">状态</th>
                        <th class="py-12 px-16 text-sm fw-medium">操作</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr v-for="cp in coupons" :key="cp.id">
                        <td class="py-12 px-16 font-monospace fw-bold">{{ cp.code }}</td>
                        <td class="py-12 px-16 text-danger-600 fw-semibold">
                          {{ cp.discount_type === 'percent' ? `减${cp.discount_value}%` : `减${cp.discount_value}能量` }}
                        </td>
                        <td class="py-12 px-16 text-sm">{{ cp.min_amount > 0 ? cp.min_amount + ' 能量' : '不限' }}</td>
                        <td class="py-12 px-16 text-sm">{{ cp.used_count }} / {{ cp.max_uses > 0 ? cp.max_uses : '不限' }}</td>
                        <td class="py-12 px-16 text-sm text-gray-500">
                          {{ cp.expires_at ? new Date(cp.expires_at).toLocaleDateString('zh-CN') : '永不过期' }}
                        </td>
                        <td class="py-12 px-16">
                          <span class="badge" :class="cp.is_active ? 'bg-success-100 text-success-600' : 'bg-gray-100 text-gray-500'">
                            {{ cp.is_active ? '启用' : '停用' }}
                          </span>
                        </td>
                        <td class="py-12 px-16 d-flex gap-6">
                          <button class="btn btn-sm btn-primary rounded-6" @click="toggleCoupon(cp)">
                            {{ cp.is_active ? '停用' : '启用' }}
                          </button>
                          <button class="btn btn-sm btn-danger rounded-6" @click="deleteCoupon(cp.id)">删除</button>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>
            </template>

        </div>
      </div>
    </div>

    <!-- 商品编辑弹窗 -->
    <div v-if="productDialog" class="modal-backdrop-custom">
      <div class="modal-box-custom">
        <div class="d-flex justify-content-between align-items-center mb-16">
          <h6 class="mb-0">{{ editingProductId ? "编辑商品" : "新增商品" }}</h6>
          <button class="btn-close" @click="productDialog = false"></button>
        </div>

        <!-- Tab 导航（编辑时才显示库存 tab） -->
        <ul class="nav nav-tabs mb-20">
          <li class="nav-item">
            <button class="nav-link" :class="{ active: pTab === 'info' }" @click="pTab = 'info'">基本信息</button>
          </li>
          <li v-if="editingProductId" class="nav-item">
            <button class="nav-link" :class="{ active: pTab === 'stock' }" @click="pTab = 'stock'">
              库存管理
              <span class="badge bg-main-100 text-main-600 ms-6" style="font-size:11px">{{ ckTotal }}</span>
            </button>
          </li>
        </ul>

        <!-- Tab：基本信息 -->
        <div v-show="pTab === 'info'">
          <div class="row gy-3">
            <div class="col-12">
              <label class="form-label fw-medium">商品名称 <span class="text-danger-600">*</span></label>
              <input v-model="pForm.name" type="text" class="form-control" placeholder="请输入商品名称" />
            </div>
            <div class="col-12">
              <label class="form-label fw-medium">分类 <span class="text-danger-600">*</span></label>
              <select v-model="pForm.category_id" class="form-select">
                <option value="" disabled>请选择分类</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <div class="col-12">
              <label class="form-label fw-medium">发货方式 <span class="text-danger-600">*</span></label>
              <div class="d-flex gap-24 mt-6">
                <label class="d-flex align-items-center gap-8" style="cursor:pointer">
                  <input type="radio" v-model.number="pForm.delivery_type" :value="0" class="form-check-input mt-0" />
                  <span>卡密自动发货</span>
                </label>
                <label class="d-flex align-items-center gap-8" style="cursor:pointer">
                  <input type="radio" v-model.number="pForm.delivery_type" :value="1" class="form-check-input mt-0" />
                  <span>手动发货（TG 号 / 实物 / 其他）</span>
                </label>
              </div>
            </div>
            <div class="col-sm-6">
              <label class="form-label fw-medium">价格（能量）<span class="text-danger-600">*</span></label>
              <input v-model.number="pForm.price" type="number" min="0" step="1" class="form-control" />
            </div>
            <div class="col-sm-6">
              <label class="form-label fw-medium">划线价（能量）</label>
              <input v-model.number="pForm.orig_price" type="number" min="0" step="1" class="form-control" />
            </div>
            <!-- 手动发货：库存字段 -->
            <div v-if="pForm.delivery_type === 1" class="col-sm-6">
              <label class="form-label fw-medium">库存数量 <span class="text-danger-600">*</span></label>
              <input v-model.number="pForm.stock_count" type="number" min="0" step="1" class="form-control" placeholder="剩余可售数量" />
            </div>
            <div class="col-sm-4">
              <label class="form-label fw-medium">促销价（能量）</label>
              <input v-model.number="pForm.promo_price" type="number" min="0" step="1" class="form-control" placeholder="0=不启用" />
            </div>
            <div class="col-sm-4">
              <label class="form-label fw-medium">促销开始</label>
              <input v-model="pForm.promo_start" type="datetime-local" class="form-control" />
            </div>
            <div class="col-sm-4">
              <label class="form-label fw-medium">促销结束</label>
              <input v-model="pForm.promo_end" type="datetime-local" class="form-control" />
            </div>
            <div class="col-12">
              <label class="form-label fw-medium">商品描述</label>
              <textarea v-model="pForm.description" class="form-control" rows="3"></textarea>
            </div>
            <!-- 多图上传 -->
            <div class="col-12">
              <label class="form-label fw-medium d-block mb-8">商品图片</label>
              <div class="product-images-grid">
                <div
                  v-for="(url, idx) in imageList"
                  :key="idx"
                  class="pimg-item"
                  :class="{ 'is-cover': url === pForm.image }"
                >
                  <img :src="url" class="pimg-thumb" />
                  <div class="pimg-actions">
                    <button
                      class="pimg-btn"
                      :class="url === pForm.image ? 'pimg-btn-star-active' : 'pimg-btn-star'"
                      :title="url === pForm.image ? '当前封面' : '设为封面'"
                      @click="pForm.image = url"
                    >★</button>
                    <button class="pimg-btn pimg-btn-del" title="删除" @click="removeImage(idx)">✕</button>
                  </div>
                  <div v-if="url === pForm.image" class="pimg-cover-badge">封面</div>
                </div>
                <label class="pimg-upload-btn" title="上传图片">
                  <span class="pimg-upload-icon">+</span>
                  <span class="pimg-upload-text">上传</span>
                  <input type="file" accept="image/*" multiple class="d-none" @change="uploadProductImages" />
                </label>
              </div>
              <div class="text-sm text-gray-400 mt-6">点击 ★ 设为封面，支持多张图片</div>
            </div>
            <div class="col-12">
              <label class="form-label fw-medium">上架状态</label>
              <div class="form-check form-switch">
                <input v-model="pForm.is_active" class="form-check-input" type="checkbox" />
                <label class="form-check-label">{{ pForm.is_active ? "上架" : "下架" }}</label>
              </div>
            </div>
          </div>
          <div class="d-flex justify-content-end gap-12 mt-24">
            <button class="btn btn-outline-main rounded-8 px-24" @click="productDialog = false">取消</button>
            <button class="btn btn-main rounded-8 px-24" :disabled="pSaving" @click="saveProduct">
              {{ pSaving ? "保存中..." : "保存" }}
            </button>
          </div>
        </div>

        <!-- Tab：库存管理 -->
        <div v-if="pTab === 'stock' && editingProductId">
          <!-- 添加卡密 -->
          <div class="mb-20">
            <label class="form-label fw-medium">批量添加卡密</label>
            <textarea v-model="cardsText" class="form-control font-monospace" rows="5"
              placeholder="每行一条，支持两种格式：&#10;卡号&#10;卡号----密码"></textarea>
            <div class="mt-8 d-flex justify-content-between align-items-center">
              <span class="text-sm text-gray-400">共 {{ cardsText.split('\n').filter(l => l.trim()).length }} 条</span>
              <button class="btn btn-main btn-sm rounded-8 px-20" :disabled="ckAdding" @click="addCardKeys">
                {{ ckAdding ? "添加中..." : "批量添加" }}
              </button>
            </div>
          </div>
          <!-- 已有卡密 -->
          <div class="border-top pt-16">
            <div class="d-flex justify-content-between align-items-center mb-12">
              <span class="fw-medium text-sm">已有卡密（{{ ckTotal }}）</span>
              <button class="btn btn-outline-main btn-sm rounded-8" @click="fetchCardKeys">刷新</button>
            </div>
            <div v-if="cardKeys.length === 0" class="text-center text-gray-400 py-16 text-sm">暂无卡密，请在上方批量添加</div>
            <div v-else class="overflow-auto" style="max-height:320px">
              <table class="table table-sm mb-0">
                <thead class="table-light">
                  <tr>
                    <th>卡号</th>
                    <th>密码</th>
                    <th>状态</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="ck in cardKeys" :key="ck.id">
                    <td class="font-monospace text-sm">{{ ck.card_no }}</td>
                    <td class="font-monospace text-sm">{{ ck.card_pwd || "-" }}</td>
                    <td>
                      <span class="badge" :class="ck.status === 0 ? 'bg-success-100 text-success-600' : 'bg-gray-100 text-gray-500'">
                        {{ ck.status === 0 ? "未售出" : ck.status === 1 ? "已售出" : "已锁定" }}
                      </span>
                    </td>
                    <td>
                      <button v-if="ck.status === 0" class="btn btn-sm btn-outline-danger rounded-6" @click="deleteCardKey(ck.id)">删除</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div class="d-flex justify-content-end mt-20">
            <button class="btn btn-outline-main rounded-8 px-24" @click="productDialog = false">关闭</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 发货弹窗 -->
    <div v-if="shipDialog" class="modal-backdrop-custom">
      <div class="modal-box-custom" style="max-width:500px">
        <div class="d-flex justify-content-between align-items-center mb-20">
          <h6 class="mb-0">发货 - {{ shipOrder?.product?.name }}</h6>
          <button class="btn-close" @click="shipDialog = false"></button>
        </div>
        <div v-if="shipOrder" class="mb-16 p-12 bg-gray-50 rounded-8 text-sm">
          <div class="text-gray-500 mb-4">买家联系方式</div>
          <div class="fw-semibold">{{ shipOrder.contact || "未填写" }}</div>
          <div class="text-gray-500 mt-8 mb-4">订单号</div>
          <div class="font-monospace">{{ shipOrder.order_no }}</div>
        </div>
        <div class="mb-16">
          <label class="form-label fw-medium">发货内容 <span class="text-danger-600">*</span></label>
          <textarea v-model="shipContent" class="form-control" rows="4"
            placeholder="TG 号、快递单号、账号密码等发货信息..."></textarea>
        </div>
        <div class="mb-20">
          <label class="form-label fw-medium">备注（可选）</label>
          <input v-model="shipNote" type="text" class="form-control" placeholder="给买家的附加说明" />
        </div>
        <div class="d-flex justify-content-end gap-12">
          <button class="btn btn-outline-main rounded-8 px-24" @click="shipDialog = false">取消</button>
          <button class="btn btn-main rounded-8 px-24" :disabled="shipping" @click="doShip">
            {{ shipping ? "发货中..." : "确认发货" }}
          </button>
        </div>
      </div>
    </div>

    <!-- 优惠券创建弹窗 -->
    <div v-if="couponDialog" class="modal-backdrop-custom">
      <div class="modal-box-custom" style="max-width:480px">
        <div class="d-flex justify-content-between align-items-center mb-20">
          <h6 class="mb-0">新建优惠券</h6>
          <button class="btn-close" @click="couponDialog = false"></button>
        </div>
        <div class="row gy-3">
          <div class="col-12">
            <label class="form-label fw-medium">券码 <span class="text-danger-600">*</span></label>
            <input v-model="couponForm.code" type="text" class="form-control" placeholder="如：SAVE20" style="text-transform:uppercase" />
          </div>
          <div class="col-12">
            <label class="form-label fw-medium">折扣类型</label>
            <div class="d-flex gap-16">
              <label class="d-flex align-items-center gap-8" style="cursor:pointer">
                <input type="radio" v-model="couponForm.discount_type" value="percent" class="form-check-input mt-0" />
                <span>百分比（如10=减10%）</span>
              </label>
              <label class="d-flex align-items-center gap-8" style="cursor:pointer">
                <input type="radio" v-model="couponForm.discount_type" value="fixed" class="form-check-input mt-0" />
                <span>固定能量</span>
              </label>
            </div>
          </div>
          <div class="col-sm-6">
            <label class="form-label fw-medium">折扣值</label>
            <input v-model.number="couponForm.discount_value" type="number" min="0.01" step="0.01" class="form-control" />
          </div>
          <div class="col-sm-6">
            <label class="form-label fw-medium">最低金额（能量，0=不限）</label>
            <input v-model.number="couponForm.min_amount" type="number" min="0" step="1" class="form-control" />
          </div>
          <div class="col-sm-6">
            <label class="form-label fw-medium">使用上限（0=不限）</label>
            <input v-model.number="couponForm.max_uses" type="number" min="0" step="1" class="form-control" />
          </div>
          <div class="col-sm-6">
            <label class="form-label fw-medium">过期时间（留空=永不过期）</label>
            <input v-model="couponForm.expires_at" type="datetime-local" class="form-control" />
          </div>
          <div class="col-12">
            <label class="form-label fw-medium">说明（可选）</label>
            <input v-model="couponForm.description" type="text" class="form-control" placeholder="如：新用户专享" />
          </div>
        </div>
        <div class="d-flex justify-content-end gap-12 mt-24">
          <button class="btn btn-outline-main rounded-8 px-24" @click="couponDialog = false">取消</button>
          <button class="btn btn-main rounded-8 px-24" :disabled="couponSaving" @click="saveCoupon">
            {{ couponSaving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "~/components/layout/banner/Breadcrumb.vue";
import AccountSidebar from "~/components/containers/account/AccountSidebar.vue";

definePageMeta({ layout: "layout-three" });

// ==================== 状态 ====================
const loading = ref(true);
const shop = ref<any>(null);
const tab = ref<"dashboard" | "info" | "products" | "orders" | "reviews" | "coupons">("info");

// 申请表单
const applyForm = ref({ name: "", contact: "", description: "" });
const applying = ref(false);
const applyError = ref("");
const applySuccess = ref(false);

// 店铺信息编辑
const infoForm = ref({ name: "", contact: "", description: "", logo: "", features: [] as { icon: string; title: string; text: string }[] });
const saving = ref(false);
const saveMsg = ref("");
const saveOk = ref(false);

// 特色服务可选图标
const featureIcons = [
  { value: "ph-fill ph-seal-check",  label: "✅ 认证" },
  { value: "ph-fill ph-key",         label: "🔑 钥匙" },
  { value: "ph-fill ph-lock-key",    label: "🔒 安全" },
  { value: "ph-fill ph-headset",     label: "🎧 客服" },
  { value: "ph-fill ph-truck",       label: "🚚 发货" },
  { value: "ph-fill ph-shield-check","label": "🛡️ 保障" },
  { value: "ph-fill ph-star",        label: "⭐ 星级" },
  { value: "ph-fill ph-gift",        label: "🎁 礼品" },
  { value: "ph-fill ph-clock",       label: "⏰ 时效" },
  { value: "ph-fill ph-currency-circle-dollar", label: "💰 价格" },
];

function addFeature() {
  if (infoForm.value.features.length >= 6) return;
  infoForm.value.features.push({ icon: "ph-fill ph-seal-check", title: "", text: "" });
}

// 商品管理
const products = ref<any[]>([]);
const categories = ref<any[]>([]);
const productDialog = ref(false);
const pTab = ref<"info" | "stock">("info");
const editingProductId = ref<number | null>(null);
const pSaving = ref(false);
const imageList = ref<string[]>([]);
const defaultPForm = () => ({
  name: "", category_id: "" as any, delivery_type: 0, stock_count: 0,
  price: 0, orig_price: 0, promo_price: 0, promo_start: "", promo_end: "",
  description: "", image: "", is_active: true,
});
const pForm = ref(defaultPForm());

watch(imageList, (list) => {
  if (pForm.value.image && !list.includes(pForm.value.image)) {
    pForm.value.image = list[0] ?? "";
  }
}, { deep: true });

// 卡密管理（集成在商品弹窗内）
const cardKeys = ref<any[]>([]);
const ckTotal = ref(0);
const cardsText = ref("");
const ckAdding = ref(false);

// 订单
const orders = ref<any[]>([]);
const orderTotal = ref(0);
const orderPage = ref(1);
const orderPageSize = 15;
const orderStatus = ref("-1");
const orderPageCount = computed(() => Math.ceil(orderTotal.value / orderPageSize));

// 发货
const shipDialog = ref(false);
const shipOrder = ref<any>(null);
const shipContent = ref("");
const shipNote = ref("");
const shipping = ref(false);

// Dashboard
const dashLoading = ref(false);
const dashData = ref<any>({});

// 评价管理
const reviewLoading = ref(false);
const reviews = ref<any[]>([]);
const reviewTotal = ref(0);
const reviewPage = ref(1);
const reviewPageSize = 15;
const reviewRating = ref(0);
const reviewStats = ref<any>(null);
const reviewPageCount = computed(() => Math.ceil(reviewTotal.value / reviewPageSize));

// ==================== 计算属性 ====================
const statusText = computed(() => {
  if (!shop.value) return "";
  const map: Record<number, string> = { 0: "审核中", 1: "运营中", 2: "已拒绝", 3: "已封禁" };
  return map[shop.value.status] ?? "未知";
});
const statusBadge = computed(() => {
  if (!shop.value) return "";
  const map: Record<number, string> = {
    0: "bg-warning-100 text-warning-600",
    1: "bg-success-100 text-success-600",
    2: "bg-danger-100 text-danger-600",
    3: "bg-gray-100 text-gray-500",
  };
  return map[shop.value.status] ?? "bg-gray-100 text-gray-500";
});

// ==================== 初始化 ====================
onMounted(async () => {
  try {
    shop.value = await $fetch<any>("/api/shop/me", { credentials: "include" });
    infoForm.value = {
      name: shop.value.name || "",
      contact: shop.value.contact || "",
      description: shop.value.description || "",
      logo: shop.value.logo || "",
      features: (() => { try { return JSON.parse(shop.value.features || "[]"); } catch { return []; } })(),
    };
    // 已审核通过的店铺默认显示 Dashboard
    if (shop.value.status === 1) {
      tab.value = "dashboard";
      await fetchDashboard();
    }
  } catch {
    shop.value = null;
  } finally {
    loading.value = false;
  }
});

async function switchTab(t: "dashboard" | "products" | "orders" | "reviews" | "coupons") {
  tab.value = t;
  if (t === "dashboard" && !dashData.value.order_stats) {
    await fetchDashboard();
  }
  if (t === "products" && products.value.length === 0) {
    await Promise.all([fetchProducts(), fetchCategories()]);
  }
  if (t === "orders" && orders.value.length === 0) {
    await fetchOrders(1);
  }
  if (t === "reviews" && reviews.value.length === 0) {
    await fetchReviews(1);
  }
  if (t === "coupons" && coupons.value.length === 0) {
    await fetchCoupons();
  }
}

// ==================== 申请开店 ====================
async function applyShop() {
  applying.value = true;
  applyError.value = "";
  try {
    await $fetch("/api/shop/apply", {
      method: "POST",
      credentials: "include",
      body: applyForm.value,
    });
    applySuccess.value = true;
  } catch (e: any) {
    applyError.value = e?.data?.error || "提交失败，请稍后重试";
  } finally {
    applying.value = false;
  }
}

// ==================== 店铺信息 ====================
async function uploadLogo(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0];
  if (!file) return;
  const fd = new FormData();
  fd.append("file", file);
  try {
    const res = await $fetch<any>("/api/upload/image", { method: "POST", credentials: "include", body: fd });
    if (res.success) infoForm.value.logo = res.url;
  } catch { /* ignore */ }
}

async function saveShop() {
  saving.value = true;
  saveMsg.value = "";
  try {
    const updated = await $fetch<any>("/api/shop/me", {
      method: "PUT",
      credentials: "include",
      body: {
        ...infoForm.value,
        features: JSON.stringify(infoForm.value.features),
      },
    });
    shop.value = updated;
    infoForm.value = {
      name: updated.name || "",
      contact: updated.contact || "",
      description: updated.description || "",
      logo: updated.logo || "",
      features: (() => { try { return JSON.parse(updated.features || "[]"); } catch { return []; } })(),
    };
    saveMsg.value = updated.status === 0 ? "保存成功，已提交重新审核，审核通过前店铺暂停营业" : "保存成功";
    saveOk.value = true;
  } catch (e: any) {
    saveMsg.value = e?.data?.error || "保存失败";
    saveOk.value = false;
  } finally {
    saving.value = false;
  }
}

// ==================== 商品管理 ====================
async function fetchProducts() {
  const data = await $fetch<any[]>("/api/shop/me/products", { credentials: "include" });
  products.value = data || [];
}

async function fetchCategories() {
  const data = await $fetch<any[]>("/api/categories/with-products", { credentials: "include" }).catch(() => []);
  categories.value = (data as any[]) || [];
}

function openProductCreate() {
  editingProductId.value = null;
  pForm.value = defaultPForm();
  imageList.value = [];
  cardKeys.value = [];
  ckTotal.value = 0;
  cardsText.value = "";
  pTab.value = "info";
  productDialog.value = true;
}

async function openProductEdit(p: any) {
  editingProductId.value = p.id;
  pForm.value = {
    name: p.name,
    category_id: p.category_id,
    delivery_type: p.delivery_type ?? 0,
    stock_count: p.stock_count ?? 0,
    price: p.price,
    orig_price: p.orig_price || 0,
    promo_price: p.promo_price || 0,
    promo_start: p.promo_start ? p.promo_start.slice(0, 16) : "",
    promo_end: p.promo_end ? p.promo_end.slice(0, 16) : "",
    description: p.description || "",
    image: p.image || "",
    is_active: p.is_active,
  };
  try {
    imageList.value = JSON.parse(p.images || "[]");
  } catch {
    imageList.value = [];
  }
  if (pForm.value.image && !imageList.value.includes(pForm.value.image)) {
    imageList.value.unshift(pForm.value.image);
  }
  cardsText.value = "";
  pTab.value = "info";
  productDialog.value = true;
  await fetchCardKeys();
}

async function openProductStock(p: any) {
  await openProductEdit(p);
  pTab.value = "stock";
}

function removeImage(idx: number) {
  imageList.value.splice(idx, 1);
}

async function uploadProductImages(e: Event) {
  const files = (e.target as HTMLInputElement).files;
  if (!files || files.length === 0) return;
  for (const file of Array.from(files)) {
    const fd = new FormData();
    fd.append("file", file);
    try {
      const res = await $fetch<any>("/api/upload/image", { method: "POST", credentials: "include", body: fd });
      if (res.success) {
        imageList.value.push(res.url);
        if (!pForm.value.image) pForm.value.image = res.url;
      }
    } catch { /* ignore */ }
  }
  // reset input so same file can be re-selected
  (e.target as HTMLInputElement).value = "";
}

async function saveProduct() {
  if (!pForm.value.name || !pForm.value.category_id) {
    alert("请填写商品名称和分类");
    return;
  }
  pSaving.value = true;
  const toRFC3339 = (s: string) => s ? new Date(s).toISOString() : null;
  const body = {
    ...pForm.value,
    images: JSON.stringify(imageList.value),
    promo_start: toRFC3339(pForm.value.promo_start),
    promo_end: toRFC3339(pForm.value.promo_end),
  };
  try {
    if (editingProductId.value) {
      await $fetch(`/api/shop/me/products/${editingProductId.value}`, { method: "PUT", credentials: "include", body });
    } else {
      await $fetch("/api/shop/me/products", { method: "POST", credentials: "include", body });
    }
    productDialog.value = false;
    await fetchProducts();
  } catch (e: any) {
    alert(e?.data?.error || "保存失败");
  } finally {
    pSaving.value = false;
  }
}

async function deleteProduct(id: number) {
  if (!confirm("确认删除该商品？")) return;
  try {
    await $fetch(`/api/shop/me/products/${id}`, { method: "DELETE", credentials: "include" });
    await fetchProducts();
  } catch (e: any) {
    alert(e?.data?.error || "删除失败");
  }
}

// ==================== 卡密管理（集成在商品弹窗） ====================
async function fetchCardKeys() {
  if (!editingProductId.value) return;
  const data = await $fetch<any>(`/api/shop/me/products/${editingProductId.value}/cardkeys`, { credentials: "include" });
  cardKeys.value = (data as any).card_keys || [];
  ckTotal.value = (data as any).total || 0;
}

async function addCardKeys() {
  if (!cardsText.value.trim() || !editingProductId.value) return;
  ckAdding.value = true;
  try {
    await $fetch(`/api/shop/me/products/${editingProductId.value}/cardkeys`, {
      method: "POST",
      credentials: "include",
      body: { cards_text: cardsText.value },
    });
    cardsText.value = "";
    await fetchCardKeys();
    // 同步更新列表中的库存数
    const idx = products.value.findIndex(p => p.id === editingProductId.value);
    if (idx !== -1) products.value[idx].stock_count = ckTotal.value;
  } catch (e: any) {
    alert(e?.data?.error || "添加失败");
  } finally {
    ckAdding.value = false;
  }
}

async function deleteCardKey(id: number) {
  if (!confirm("确认删除该卡密？")) return;
  try {
    await $fetch(`/api/shop/me/cardkeys/${id}`, { method: "DELETE", credentials: "include" });
    await fetchCardKeys();
    const idx = products.value.findIndex(p => p.id === editingProductId.value);
    if (idx !== -1) products.value[idx].stock_count = ckTotal.value;
  } catch (e: any) {
    alert(e?.data?.error || "删除失败");
  }
}

// ==================== 订单 ====================
async function fetchOrders(page: number) {
  orderPage.value = page;
  const data = await $fetch<any>("/api/shop/me/orders", {
    credentials: "include",
    params: { page, page_size: orderPageSize, status: orderStatus.value },
  });
  orders.value = (data as any).data || [];
  orderTotal.value = (data as any).total || 0;
}

function orderStatusText(s: number) {
  const map: Record<number,string> = { 0:"待支付", 1:"待发货", 2:"已完成", 3:"已取消", 4:"已发货/待确认" };
  return map[s] ?? "未知";
}
function orderStatusClass(s: number) {
  const map: Record<number,string> = {
    0: "bg-warning-100 text-warning-600",
    1: "bg-main-100 text-main-600",
    2: "bg-success-100 text-success-600",
    3: "bg-gray-100 text-gray-500",
    4: "bg-purple-100 text-purple-600",
  };
  return map[s] ?? "bg-gray-100 text-gray-500";
}

function openShipDialog(o: any) {
  shipOrder.value = o;
  shipContent.value = "";
  shipNote.value = "";
  shipDialog.value = true;
}

async function doShip() {
  if (!shipContent.value.trim()) { alert("请填写发货内容"); return; }
  shipping.value = true;
  try {
    await $fetch(`/api/shop/me/orders/${shipOrder.value.id}/ship`, {
      method: "POST",
      credentials: "include",
      body: { content: shipContent.value, note: shipNote.value },
    });
    shipDialog.value = false;
    // 更新本地订单状态
    const idx = orders.value.findIndex(o => o.id === shipOrder.value.id);
    if (idx !== -1) orders.value[idx].status = 4;
  } catch (e: any) {
    alert(e?.data?.error || "发货失败");
  } finally {
    shipping.value = false;
  }
}

function fmtDate(d: string) {
  if (!d) return "-";
  return new Date(d).toLocaleString("zh-CN", { hour12: false });
}

// ==================== 优惠券 ====================
const coupons = ref<any[]>([]);
const couponLoading = ref(false);
const couponDialog = ref(false);
const couponSaving = ref(false);
const couponForm = ref({
  code: "", discount_type: "percent", discount_value: 10,
  min_amount: 0, max_uses: 0, expires_at: "", description: "",
});

async function fetchCoupons() {
  couponLoading.value = true;
  try {
    const res = await $fetch<any>("/api/shop/me/coupons", { credentials: "include" });
    coupons.value = res.data || [];
  } catch { coupons.value = []; } finally { couponLoading.value = false; }
}

function openCouponCreate() {
  couponForm.value = { code: "", discount_type: "percent", discount_value: 10, min_amount: 0, max_uses: 0, expires_at: "", description: "" };
  couponDialog.value = true;
}

async function saveCoupon() {
  if (!couponForm.value.code.trim()) { alert("请填写券码"); return; }
  couponSaving.value = true;
  try {
    const body: any = { ...couponForm.value, code: couponForm.value.code.toUpperCase() };
    if (!body.expires_at) delete body.expires_at;
    else body.expires_at = new Date(body.expires_at).toISOString();
    await $fetch("/api/shop/me/coupons", { method: "POST", credentials: "include", body });
    couponDialog.value = false;
    await fetchCoupons();
  } catch (e: any) { alert(e?.data?.error || "创建失败"); } finally { couponSaving.value = false; }
}

async function toggleCoupon(cp: any) {
  try {
    await $fetch(`/api/shop/me/coupons/${cp.id}`, {
      method: "PUT", credentials: "include",
      body: { is_active: !cp.is_active },
    });
    await fetchCoupons();
  } catch { alert("操作失败"); }
}

async function deleteCoupon(id: number) {
  if (!confirm("确认删除该优惠券？")) return;
  try {
    await $fetch(`/api/shop/me/coupons/${id}`, { method: "DELETE", credentials: "include" });
    await fetchCoupons();
  } catch { alert("删除失败"); }
}

// ==================== Dashboard ====================
async function fetchDashboard() {
  dashLoading.value = true;
  try {
    dashData.value = await $fetch<any>("/api/shop/me/dashboard", { credentials: "include" });
  } catch {
    dashData.value = {};
  } finally {
    dashLoading.value = false;
  }
}

function ratingPercent(star: number): number {
  const total = dashData.value.review_stats?.total ?? 0;
  if (total === 0) return 0;
  return ((dashData.value.review_stats?.rating_distribution?.[star] ?? 0) / total) * 100;
}

// ==================== 评价管理 ====================
async function fetchReviews(page: number) {
  reviewLoading.value = true;
  reviewPage.value = page;
  try {
    const data = await $fetch<any>("/api/shop/me/reviews", {
      credentials: "include",
      params: { page, page_size: reviewPageSize, rating: reviewRating.value },
    });
    reviews.value = data.data || [];
    reviewTotal.value = data.total || 0;
    reviewStats.value = data.stats || null;
  } catch {
    reviews.value = [];
    reviewTotal.value = 0;
  } finally {
    reviewLoading.value = false;
  }
}

function filterReviews(rating: number) {
  reviewRating.value = rating;
  fetchReviews(1);
}
</script>

<style scoped>
.modal-backdrop-custom {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.45);
  z-index: 1050;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.modal-box-custom {
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  width: 100%;
  max-width: 680px;
  max-height: 90vh;
  overflow-y: auto;
}

/* 多图上传网格 */
.product-images-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.pimg-item {
  position: relative;
  width: 88px;
  height: 88px;
  border-radius: 8px;
  overflow: hidden;
  border: 2px solid #e5e7eb;
  flex-shrink: 0;
}
.pimg-item.is-cover {
  border-color: #e6162d;
}
.pimg-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}
.pimg-actions {
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,.42);
  display: none;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.pimg-item:hover .pimg-actions {
  display: flex;
}
.pimg-btn {
  border: none;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}
.pimg-btn-star       { background: rgba(255,255,255,.85); color: #888; }
.pimg-btn-star-active{ background: #fff; color: #f59e0b; }
.pimg-btn-del        { background: rgba(255,255,255,.85); color: #e6162d; }
.pimg-cover-badge {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: #e6162d;
  color: #fff;
  font-size: 11px;
  text-align: center;
  padding: 2px 0;
}
.pimg-upload-btn {
  width: 88px;
  height: 88px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #9ca3af;
  flex-shrink: 0;
  transition: border-color .2s, color .2s;
}
.pimg-upload-btn:hover {
  border-color: #e6162d;
  color: #e6162d;
}
.pimg-upload-icon { font-size: 24px; line-height: 1; }
.pimg-upload-text { font-size: 12px; margin-top: 4px; }
</style>
