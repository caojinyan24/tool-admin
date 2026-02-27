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
            :key="paymentOrder.transaction_id"
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
      <el-button type="primary" @click="fetchDetail">重新加载</el-button>
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
import { payDataApi } from '@/api/payDataApi';
import TradeOrderCard from './pay-data-detail/TradeOrderCard.vue';
import PaymentOrderTab from './pay-data-detail/PaymentOrderTab.vue';
import * as Icons from '@element-plus/icons-vue';

export default {
  name: 'PayDataDetail',
  components: {
    TradeOrderCard,
    PaymentOrderTab,
    ArrowLeft: Icons.ArrowLeft
  },
  data() {
    return {
      loading: false,
      detailData: null,
      error: null,
      jsonDialogVisible: false,
      jsonData: {},
      activeTab: '',
      expandedLogs: [],
      formattedResponses: {}, // 存储格式化后的响应内容
      // 折叠状态，默认全部展开
      collapsedCoreOrders: [],
      collapsedChannelOrders: [],
      collapsedChannelLogs: [],
      // 新增：折叠状态 - 通知消息
      collapsedNotifyMessages: [],
      collapsedNotifyHistories: [],
      // 所有通知消息折叠状态
      allNotifyMessagesCollapsed: false
    };
  },
  computed: {
    jsonDataStr() {
      return JSON.stringify(this.jsonData, null, 2);
    }
  },
  mounted() {
    this.fetchDetail();
  },
  
  created() {
  },
  methods: {
    // 获取详情数据
    async fetchDetail() {
      const { trade_id, month } = this.$route.query;
      
      if (!trade_id || !month) {
        this.error = '缺少必要的查询参数';
        this.loading = false;
        return;
      }
      
      this.loading = true;
      this.error = null;
      
      try {
        // 使用实际的API调用获取后端数据
        this.detailData = await payDataApi.queryDetail({ trade_id, month });
        
        // 设置默认选中第一个tab
        if (this.detailData.payment_orders && this.detailData.payment_orders.length > 0) {
          this.activeTab = 'payment_0';
        }
      } catch (error) {
        this.error = error.message || '加载详情失败';
      } finally {
        this.loading = false;
      }
    },
    
    // 返回上一页
    goBack() {
      this.$router.push('/pay-data');
    },
    
    // 查看JSON详情
    viewJson(data) {
      this.jsonData = data;
      this.jsonDialogVisible = true;
    },
    
    // Tab点击处理
    handleTabClick(tab) {
      console.log('Tab clicked:', tab);
    },
    
    // 根据payment_order_id获取core orders
    getCoreOrdersByPaymentOrderId(paymentOrderId) {
      if (!this.detailData || !this.detailData.core_orders) {
        return [];
      }
      return this.detailData.core_orders.filter(order => order.trade_id === paymentOrderId).sort((a, b) => new Date(b.create_time) - new Date(a.create_time));
    },
    
    // 根据payment_order的transaction_id获取关联的通知消息
    // 严格按照文档中的关联关系：notify_msg.biz_id 和 payment_orders.transaction_id 关联
    getNotifyMessagesByPaymentOrderId(transactionId) {
      if (!this.detailData || !this.detailData.notify_messages) {
        return [];
      }
      
      // 严格按照文档要求，使用notify_msg.biz_id与payment_orders.transaction_id关联
      return this.detailData.notify_messages
        .filter(msg => msg.biz_id === transactionId)
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    },
    
    // 根据core_order的order_id获取关联的通知消息
    // 严格按照文档中的关联关系：notify_msg.biz_id 和 core_orders.order_id 关联
    getNotifyMessagesByCoreOrderId(orderId) {
      if (!this.detailData || !this.detailData.notify_messages) {
        return [];
      }
      
      // 严格按照文档要求，使用notify_msg.biz_id与core_orders.order_id关联
      return this.detailData.notify_messages
        .filter(msg => msg.biz_id === orderId)
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    },
    
    // 折叠功能相关方法
    toggleCoreOrder(paymentIndex, coreIndex) {
      const key = `${paymentIndex}_${coreIndex}`;
      const idx = this.collapsedCoreOrders.indexOf(key);
      if (idx > -1) {
        this.collapsedCoreOrders.splice(idx, 1);
      } else {
        this.collapsedCoreOrders.push(key);
      }
    },
    
    isCoreOrderCollapsed(paymentIndex, coreIndex) {
      const key = `${paymentIndex}_${coreIndex}`;
      return this.collapsedCoreOrders.includes(key);
    },
    
    toggleAllCoreOrders(paymentIndex) {
      const coreOrders = this.getCoreOrdersByPaymentOrderId(this.detailData.payment_orders[paymentIndex].transaction_id);
      const allCollapsed = coreOrders.every((_, index) => 
        this.isCoreOrderCollapsed(paymentIndex, index)
      );
      
      // 清除该paymentOrder下所有coreOrder的折叠状态
      this.collapsedCoreOrders = this.collapsedCoreOrders.filter(key => 
        !key.startsWith(`${paymentIndex}_`)
      );
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed) {
        coreOrders.forEach((_, index) => {
          this.collapsedCoreOrders.push(`${paymentIndex}_${index}`);
        });
      }
    },
    
    toggleChannelOrder(paymentIndex, coreIndex, channelIndex) {
      const key = `${paymentIndex}_${coreIndex}_${channelIndex}`;
      const idx = this.collapsedChannelOrders.indexOf(key);
      if (idx > -1) {
        this.collapsedChannelOrders.splice(idx, 1);
      } else {
        this.collapsedChannelOrders.push(key);
      }
    },
    
    isChannelOrderCollapsed(paymentIndex, coreIndex, channelIndex) {
      const key = `${paymentIndex}_${coreIndex}_${channelIndex}`;
      return this.collapsedChannelOrders.includes(key);
    },
    
    toggleAllChannelOrders(paymentIndex, coreIndex) {
      const coreOrder = this.getCoreOrdersByPaymentOrderId(
        this.detailData.payment_orders[paymentIndex].transaction_id
      )[coreIndex];
      const channelOrders = this.detailData.channel_orders.filter(order => order.order_id === coreOrder.order_id);
      
      const allCollapsed = channelOrders.every((_, index) => 
        this.isChannelOrderCollapsed(paymentIndex, coreIndex, index)
      );
      
      // 清除该coreOrder下所有channelOrder的折叠状态
      this.collapsedChannelOrders = this.collapsedChannelOrders.filter(key => 
        !key.startsWith(`${paymentIndex}_${coreIndex}_`)
      );
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed) {
        channelOrders.forEach((_, index) => {
          this.collapsedChannelOrders.push(`${paymentIndex}_${coreIndex}_${index}`);
        });
      }
    },
    
    toggleChannelLog(paymentIndex, coreIndex, logIndex) {
      const key = `${paymentIndex}_${coreIndex}_${logIndex}`;
      const idx = this.collapsedChannelLogs.indexOf(key);
      if (idx > -1) {
        this.collapsedChannelLogs.splice(idx, 1);
      } else {
        this.collapsedChannelLogs.push(key);
      }
    },
    
    isChannelLogCollapsed(paymentIndex, coreIndex, logIndex) {
      const key = `${paymentIndex}_${coreIndex}_${logIndex}`;
      return this.collapsedChannelLogs.includes(key);
    },
    
    toggleAllChannelLogs(paymentIndex, coreIndex) {
      const coreOrder = this.getCoreOrdersByPaymentOrderId(
        this.detailData.payment_orders[paymentIndex].transaction_id
      )[coreIndex];
      const channelLogs = this.detailData.channel_logs
        .filter(log => log.order_id === coreOrder.order_id)
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      
      const allCollapsed = channelLogs.every((_, index) => 
        this.isChannelLogCollapsed(paymentIndex, coreIndex, index)
      );
      
      // 清除该coreOrder下所有channelLog的折叠状态
      this.collapsedChannelLogs = this.collapsedChannelLogs.filter(key => 
        !key.startsWith(`${paymentIndex}_${coreIndex}_`)
      );
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed) {
        channelLogs.forEach((_, index) => {
          this.collapsedChannelLogs.push(`${paymentIndex}_${coreIndex}_${index}`);
        });
      }
    },
    
    // 新增：切换通知消息折叠状态
    toggleNotifyMessage(paymentIndex, msgIndex) {
      const key = `${paymentIndex}_msg_${msgIndex}`;
      const idx = this.collapsedNotifyMessages.indexOf(key);
      if (idx > -1) {
        this.collapsedNotifyMessages.splice(idx, 1);
      } else {
        this.collapsedNotifyMessages.push(key);
      }
    },
    
    // 新增：检查通知消息折叠状态
    isNotifyMessageCollapsed(paymentIndex, msgIndex) {
      const key = `${paymentIndex}_msg_${msgIndex}`;
      return this.collapsedNotifyMessages.includes(key);
    },
    
    // 新增：切换通知历史折叠状态
    toggleNotifyHistory(paymentIndex, msgIndex, historyIndex) {
      const key = `${paymentIndex}_msg_${msgIndex}_history_${historyIndex}`;
      const idx = this.collapsedNotifyHistories.indexOf(key);
      if (idx > -1) {
        this.collapsedNotifyHistories.splice(idx, 1);
      } else {
        this.collapsedNotifyHistories.push(key);
      }
    },
    
    // 新增：检查通知历史折叠状态
    isNotifyHistoryCollapsed(paymentIndex, msgIndex, historyIndex) {
      const key = `${paymentIndex}_msg_${msgIndex}_history_${historyIndex}`;
      return this.collapsedNotifyHistories.includes(key);
    },
    
    // 切换所有通知消息折叠状态
    toggleAllNotifyMessages(paymentIndex) {
      console.log('切换所有通知消息折叠状态 - paymentIndex:', paymentIndex);
      const paymentOrders = this.detailData?.payment_orders || [];
      if (paymentIndex >= 0 && paymentIndex < paymentOrders.length) {
        const paymentOrder = paymentOrders[paymentIndex];
        if (paymentOrder && paymentOrder.transaction_id) {
          // 获取该支付订单的通知消息
          const notifyMessages = this.getNotifyMessagesByPaymentOrderId(paymentOrder.transaction_id) || [];
          console.log('切换所有通知消息 - 消息数量:', notifyMessages.length);
          
          // 反转折叠状态
          const newCollapsedState = !this.allNotifyMessagesCollapsed;
          console.log('切换到新状态:', newCollapsedState);
          this.allNotifyMessagesCollapsed = newCollapsedState;
          
          // 更新所有消息的折叠状态
          for (let msgIndex = 0; msgIndex < notifyMessages.length; msgIndex++) {
            this.toggleNotifyMessage(paymentIndex, msgIndex);
          }
        }
      }
    },
    
    // 获取所有通知消息折叠状态
    getAllNotifyMessagesCollapsed() {
      return this.allNotifyMessagesCollapsed;
    },
    
    // 切换调试数据显示
    toggleDebugData() {
      this.debugDataVisible = !this.debugDataVisible;
    },
    
    // 切换日志内容展开/收起状态并格式化JSON
    toggleLogExpand(index, response) {
      // 确保formattedResponses对象已初始化
      if (!this.formattedResponses) {
        this.formattedResponses = {};
      }
      
      const logIndex = this.expandedLogs.indexOf(index);
      if (logIndex > -1) {
        // 收起时从展开列表中移除
        this.expandedLogs.splice(logIndex, 1);
      } else {
        // 展开时添加到展开列表并格式化内容
        this.expandedLogs.push(index);
        // 格式化JSON响应内容
        try {
          // 尝试解析为JSON并格式化
          if (response && typeof response === 'string') {
            this.formattedResponses[index] = JSON.stringify(JSON.parse(response), null, 2);
          } else {
            this.formattedResponses[index] = JSON.stringify(response, null, 2);
          }
        } catch (e) {
          // 如果不是有效的JSON，则使用原始内容
          this.formattedResponses[index] = response || '';
        }
      }
      
      // 强制更新，确保UI正确响应状态变化
      this.$forceUpdate();
    }
  },
  watch: {
    // 监听路由参数变化，重新加载数据
    '$route.query': {
      handler() {
        this.fetchDetail();
      },
      deep: true
    }
  }
};
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

.payment-orders-section {
  margin-top: 30px;
}

.section-title {
  margin: 30px 0 15px 0;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  border-left: 4px solid #409eff;
  padding-left: 10px;
}

.json-pre {
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 600px;
  overflow: auto;
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
}

/* 状态高亮样式 */
:deep(.status-success) {
  color: #67c23a;
  background-color: #f0f9eb;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

:deep(.status-fail) {
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

:deep(.status-cancelled) {
  color: #909399;
  background-color: #f4f4f5;
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
</style>