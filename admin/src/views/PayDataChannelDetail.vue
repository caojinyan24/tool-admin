<template>
  <div class="pay-data-detail">
       <!-- 返回按钮 -->
       <el-button type="primary" icon="ArrowLeft" @click="goBack" style="margin-bottom: 10px;">返回</el-button>
   
       <!-- 加载状态 -->
        <div v-if="loading" class="loading-container">
          <el-loading-spinner />
          <p>加载中...</p>
        </div>
   
        <!-- 数据内容 -->
        <div v-else-if="detailData">
        <!-- Trade Order信息 -->
        <trade-order-card
          :trade-order="detailData.trade_order"
          @view-json="viewJson"
        />
        
        <!-- Payment Orders - 使用Tab页展示 -->
        <div class="payment-orders-section">
          <h3 class="section-title">Payment Orders</h3>
          <el-tabs v-model="activeTab" type="border-card" @tab-click="handleTabClick">
            <el-tab-pane
              v-for="(paymentOrder, index) in detailData.payment_orders"
              :key="paymentOrder.payment_order_id"
              :label="`Payment Order #${index + 1} (${paymentOrder.status})`"
              :name="`payment_${index}`"
            >
              <payment-order-tab
                :payment-order="paymentOrder"
                :payment-index="index"
                :core-orders="getCoreOrdersByPaymentOrderId(paymentOrder.transaction_id)"
                :channel-orders="detailData.channel_orders || []"
                :channel-logs="detailData.channel_logs || []"
                :notify-messages="detailData.notify_messages || []"
                :notify-histories="detailData.notify_histories || []"
                :expanded-logs="expandedLogs"
                :formatted-responses="formattedResponses"
                @view-json="viewJson"
                @toggle-log-expand="toggleLogExpand"
                @toggle-notify-message="toggleNotifyMessage"
                @toggle-notify-history="toggleNotifyHistory"
              />
            </el-tab-pane>
          </el-tabs>
        </div>
        </div>
        
        <!-- 无数据提示 -->
        <div v-else-if="!loading && !error" class="empty-container">
          <el-empty description="暂无详细数据" />
        </div>
        
        <!-- 错误提示 -->
        <div v-else-if="error" class="error-container">
          <el-alert
            title="加载失败"
            :description="error"
            type="error"
            show-icon
            :closable="false"
          />
          <el-button type="primary" @click="fetchData">重新加载</el-button>
        </div>
   
        <!-- JSON详情弹窗 -->
        <el-dialog v-model="jsonDialogVisible" title="记录详情（JSON格式）" width="80%">
          <pre class="json-pre">{{ jsonDataStr }}</pre>
          <template #footer>
            <el-button @click="jsonDialogVisible = false">关闭</el-button>
          </template>
        </el-dialog>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { payDataApi } from '@/api/payDataApi'
import { ArrowLeft } from '@element-plus/icons-vue'
import TradeOrderCard from './pay-data-detail/TradeOrderCard.vue'
import PaymentOrderTab from './pay-data-detail/PaymentOrderTab.vue'

export default {
  name: 'PayDataChannelDetail',
  components: {
    TradeOrderCard,
    PaymentOrderTab,
    ArrowLeft
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    
    // 路由参数
    const orderId = ref(route.query.order_id || '')
    const month = ref(route.query.month || '')
    
    // 组件挂载状态
    const isMounted = ref(true)
    
    // 数据状态
    const loading = ref(true)
    const error = ref('')
    const detailData = ref(null)
    
    // JSON弹窗状态
    const jsonDialogVisible = ref(false)
    const jsonData = ref({})
    
    // 标签页状态
    const activeTab = ref('')
    
    // 展开的日志状态
    const expandedLogs = ref([])
    
    // 格式化后的响应内容
    const formattedResponses = ref({})
    
    // 折叠状态
    const collapsedCoreOrders = ref([])
    const collapsedChannelOrders = ref([])
    const collapsedChannelLogs = ref([])
    const collapsedNotifyMessages = ref([])
    const collapsedNotifyHistories = ref([])
    const allNotifyMessagesCollapsed = ref(false)
    
    // 返回上一页
    const goBack = () => {
      // 检查是否有保存的渠道订单查询条件
      const savedChannelQueryForm = sessionStorage.getItem('payDataChannelQueryForm');
      if (savedChannelQueryForm) {
        // 如果有，返回到渠道订单查询tab
        router.push({
          path: '/pay-data',
          query: { activeTab: 'channel-query' }
        });
      } else {
        // 否则，返回到默认的交易查询tab
        router.push('/pay-data');
      }
    }
    
    // 数据获取
    const fetchData = async () => {
      loading.value = true
      error.value = ''

      try {
        // 使用实际的API调用获取后端数据
        detailData.value = await payDataApi.queryChannelOrderDetail({
          order_id: orderId.value,
          month: month.value
        })
        
        // 设置默认选中第一个tab
        if (detailData.value.payment_orders && detailData.value.payment_orders.length > 0) {
          activeTab.value = 'payment_0'
        }
      } catch (err) {
        console.error('Error fetching data:', err)
        if (!isMounted.value) return
        error.value = err.message || '网络错误，请稍后重试'
      } finally {
        if (!isMounted.value) return
        loading.value = false
      }
    }

    // 查看JSON详情
    const viewJson = (data) => {
      jsonData.value = data
      jsonDialogVisible.value = true
    }
    
    // 根据payment_order_id获取core orders
    const getCoreOrdersByPaymentOrderId = (paymentOrderId) => {
      if (!detailData.value || !detailData.value.core_orders) {
        return []
      }
      return detailData.value.core_orders
        .filter(order => order.trade_id === paymentOrderId)
        .sort((a, b) => new Date(b.create_time) - new Date(a.create_time))
    }
    
    // Tab点击处理
    const handleTabClick = (tab) => {
      console.log('Tab clicked:', tab)
    }
    
    // 切换日志展开状态
    const toggleLogExpand = (index, response) => {
      const logKey = `${index}-${response}`
      const expandedIndex = expandedLogs.value.indexOf(logKey)
      
      if (expandedIndex > -1) {
        expandedLogs.value.splice(expandedIndex, 1)
      } else {
        expandedLogs.value.push(logKey)
      }
    }
    
    // 切换通知消息折叠状态
    const toggleNotifyMessage = (paymentIndex, msgIndex) => {
      const key = `${paymentIndex}-${msgIndex}`
      const index = collapsedNotifyMessages.value.indexOf(key)
      
      if (index > -1) {
        collapsedNotifyMessages.value.splice(index, 1)
      } else {
        collapsedNotifyMessages.value.push(key)
      }
    }
    
    // 切换通知历史折叠状态
    const toggleNotifyHistory = (paymentIndex, msgIndex, historyIndex) => {
      const key = `${paymentIndex}-${msgIndex}-${historyIndex}`
      const index = collapsedNotifyHistories.value.indexOf(key)
      
      if (index > -1) {
        collapsedNotifyHistories.value.splice(index, 1)
      } else {
        collapsedNotifyHistories.value.push(key)
      }
    }
    
    // 计算属性：格式化JSON数据
    const jsonDataStr = computed(() => {
      return JSON.stringify(jsonData.value, null, 2)
    })
    
    // 页面加载时获取数据
    onMounted(() => {
      if (orderId.value && month.value) {
        fetchData()
      } else {
        error.value = '缺少必要的参数'
        loading.value = false
      }
    })

    // 组件卸载时设置挂载状态为false
    onUnmounted(() => {
      isMounted.value = false
    })

    return {
      orderId,
      month,
      loading,
      error,
      detailData,
      activeTab,
      jsonDialogVisible,
      jsonData,
      jsonDataStr,
      expandedLogs,
      formattedResponses,
      collapsedCoreOrders,
      collapsedChannelOrders,
      collapsedChannelLogs,
      collapsedNotifyMessages,
      collapsedNotifyHistories,
      allNotifyMessagesCollapsed,
      goBack,
      fetchData,
      viewJson,
      handleTabClick,
      getCoreOrdersByPaymentOrderId,
      toggleLogExpand,
      toggleNotifyMessage,
      toggleNotifyHistory
    }
  }
}
</script>

<style scoped>
.pay-data-detail {
  padding: 20px;
}

.loading-container,
.error-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
}

.loading-container p {
  margin-top: 20px;
  color: #606266;
}

.error-container {
  gap: 20px;
}

.detail-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 区域标题样式 */
.section-title {
  margin: 30px 0 15px 0;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  border-left: 4px solid #409eff;
  padding-left: 10px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  word-break: break-all;
}

.info-item label {
  width: 120px;
  font-weight: 500;
  color: #606266;
  margin-right: 8px;
}

.info-item span {
  flex: 1;
  color: #303133;
}

.payment-order {
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.payment-order:last-child {
  margin-bottom: 0;
  padding-bottom: 0;
  border-bottom: none;
}

.payment-order-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 16px;
}

.no-data {
  text-align: center;
  color: #909399;
  padding: 20px;
}

.message-section,
.log-section {
  margin-bottom: 20px;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.sub-section-title {
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 16px;
}

/* 状态高亮样式 */
:deep(.status-success) {
  color: #67c23a;
  background-color: #f0f9eb;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-fail),
:deep(.status-failed),
:deep(.status-error) {
  color: #f56c6c;
  background-color: #fef0f0;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-processing) {
  color: #e6a23c;
  background-color: #fdf6ec;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-pending) {
  color: #409eff;
  background-color: #ecf5ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-init) {
  color: #909399;
  background-color: #f0f2f5;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-closed),
:deep(.status-cancelled) {
  color: #909399;
  background-color: #f4f4f5;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-refunded) {
  color: #409eff;
  background-color: #ecf5ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-default) {
  color: #909399;
  background-color: #f0f2f5;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

/* JSON查看样式 */
.json-pre {
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 600px;
  overflow: auto;
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
}

/* 折叠相关样式 */
.title-with-toggle {
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-bottom: 15px;
  user-select: none;
}

.title-with-toggle:hover {
  color: #409eff;
}

.toggle-icon {
  margin-right: 8px;
  font-size: 14px;
  transition: transform 0.3s;
}

.subsection-title {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
  display: inline;
}
</style>