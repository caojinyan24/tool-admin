<template>
  <div class="payment-order-tab">
    <el-card class="detail-card">
      <template #header>
        <div class="card-header">
          <h3 class="section-title">Payment Order</h3>
          <el-button type="primary" size="small" @click="viewPaymentOrderJson">查看详情</el-button>
        </div>
      </template>
      
      <div class="info-grid">
        <div class="info-item">
          <label>交易ID：</label>
          <span>{{ paymentOrder.transaction_id }}</span>
        </div>
        <div class="info-item">
          <label>支付方式：</label>
          <span>{{ paymentOrder.payment_method }}</span>
        </div>
        <div class="info-item">
          <label>创建时间：</label>
          <span>{{ paymentOrder.create_time_utc }}</span>
        </div>
        <div class="info-item">
          <label>支付状态：</label>
          <span :class="getStatusClass(paymentOrder.status)">{{ paymentOrder.status }}</span>
        </div>
        <div class="info-item">
          <label>支付金额：</label>
          <span>{{ formatAmount(paymentOrder.payment_amount, paymentOrder.payment_precision, paymentOrder.payment_currency) }}</span>
        </div>
      </div>
      
      <!-- 子Tab结构 -->
      <el-tabs v-model="activeSubTab" class="sub-tabs">
        <!-- 第一个Tab：支付订单信息(包括coreorder和channelorders) -->
        <el-tab-pane label="支付订单信息" name="payment-order-info">
          <!-- Core Orders -->
          <div v-if="coreOrders.length > 0">
            <div class="title-with-toggle" @click="toggleAllCoreOrders">
              <el-icon class="toggle-icon">
                <svg v-if="!allCoreOrdersCollapsed" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M8 12a.5.5 0 0 0 .5-.5V5.707l2.146 2.147a.5.5 0 0 0 .708-.708l-3-3a.5.5 0 0 0-.708 0l-3 3a.5.5 0 1 0 .708.708L7.5 5.707V11.5a.5.5 0 0 0 .5.5z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M4.5 11.5a.5.5 0 0 0 .5.5h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L13.293 11H5a.5.5 0 0 0-.5.5z"/>
                </svg>
              </el-icon>
              <h5 class="subsection-title">Core Orders ({{ coreOrders.length }})</h5>
            </div>
            <div v-if="!allCoreOrdersCollapsed">
              <core-order-card
                v-for="(coreOrder, coreIndex) in coreOrders"
                :key="coreOrder.order_id || coreIndex"
                :core-order="coreOrder"
                :payment-index="paymentIndex"
                :index="coreIndex"
                :channel-orders="getChannelOrdersByOrderId(coreOrder.order_id)"
                :collapsed-state="isCoreOrderCollapsed(coreIndex) ? 'collapsed' : 'expanded'"
                :expanded-logs="expandedLogs"
                :formatted-responses="formattedResponses"
                @toggle-collapse="toggleCoreOrder"
                @view-json="viewJson"
                @toggle-log-expand="toggleLogExpand"
              />
            </div>
          </div>
        </el-tab-pane>
        
        <!-- 第二个Tab：通知消息 -->
        <el-tab-pane label="通知消息" name="notify-messages">
          <div v-if="getFilteredNotifyMessages().length > 0">
            <div class="title-with-toggle" @click="toggleAllNotifyMessages">
              <el-icon class="toggle-icon">
                <svg v-if="!allNotifyMessagesCollapsed" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M8 12a.5.5 0 0 0 .5-.5V5.707l2.146 2.147a.5.5 0 0 0 .708-.708l-3-3a.5.5 0 0 0-.708 0l-3 3a.5.5 0 1 0 .708.708L7.5 5.707V11.5a.5.5 0 0 0 .5.5z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M4.5 11.5a.5.5 0 0 0 .5.5h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L13.293 11H5a.5.5 0 0 0-.5.5z"/>
                </svg>
              </el-icon>
              <h5 class="subsection-title">通知消息 ({{ getFilteredNotifyMessages().length }})</h5>
            </div>
            <div v-if="!allNotifyMessagesCollapsed">
              <notify-message-card
                v-for="(notifyMessage, msgIndex) in getFilteredNotifyMessages()"
                :key="notifyMessage.id || msgIndex"
                :notify-message="notifyMessage"
                :payment-index="paymentIndex"
                :msg-index="msgIndex"
                :notify-histories="notifyMessage.notifyList || []"
                :is-collapsed="isNotifyMessageCollapsed(msgIndex)"
                :expanded-logs="expandedLogs"
                :formatted-responses="formattedResponses"
                @toggle-collapse="$emit('toggle-notify-message', paymentIndex, msgIndex)"
                @toggle-history-collapse="$emit('toggle-notify-history', paymentIndex, msgIndex, $event)"
                @view-json="viewJson"
                @toggle-log-expand="toggleLogExpand"
              />
            </div>
          </div>
          <div v-else class="empty-message">
            <el-empty description="暂无通知消息" />
          </div>
        </el-tab-pane>
        
        <!-- 第三个Tab：channel logs - 只显示与当前transaction_id关联的日志 -->
        <el-tab-pane label="外部渠道交互日志" name="channel-logs">
          <div v-if="getChannelLogsByTransactionId().length > 0">
            <div class="title-with-toggle" @click="toggleAllChannelLogs">
              <el-icon class="toggle-icon">
                <svg v-if="!allChannelLogsCollapsed" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M8 12a.5.5 0 0 0 .5-.5V5.707l2.146 2.147a.5.5 0 0 0 .708-.708l-3-3a.5.5 0 0 0-.708 0l-3 3a.5.5 0 1 0 .708.708L7.5 5.707V11.5a.5.5 0 0 0 .5.5z"/>
                </svg>
                <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
                  <path fill="currentColor" d="M4.5 11.5a.5.5 0 0 0 .5.5h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L13.293 11H5a.5.5 0 0 0-.5.5z"/>
                </svg>
              </el-icon>
              <h5 class="subsection-title">Channel Logs ({{ getChannelLogsByTransactionId().length }})</h5>
            </div>
            <div v-if="!allChannelLogsCollapsed">
              <channel-log-card
                v-for="(log, logIndex) in getChannelLogsByTransactionId()"
                :key="log.id || logIndex"
                :channel-log="log"
                :payment-index="paymentIndex"
                :log-index="logIndex"
                :is-collapsed="isChannelLogCollapsed(logIndex)"
                :expanded-logs="expandedLogs"
                :formatted-responses="formattedResponses"
                @toggle-collapse="toggleChannelLog(logIndex)"
                @view-json="viewJson"
                @toggle-log-expand="toggleLogExpand"
              />
            </div>
          </div>
          <div v-else class="empty-message">
            <el-empty description="暂无Channel Logs" />
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script>
import CoreOrderCard from './CoreOrderCard.vue';
import NotifyMessageCard from './NotifyMessageCard.vue';
import ChannelLogCard from './ChannelLogCard.vue';
import { ElIcon } from 'element-plus';

export default {
  name: 'PaymentOrderTab',
  components: {
    CoreOrderCard,
    NotifyMessageCard,
    ChannelLogCard,
    ElIcon
  },
  props: {
    paymentOrder: {
      type: Object,
      required: true
    },
    paymentIndex: {
      type: Number,
      required: true
    },
    coreOrders: {
      type: Array,
      default: () => []
    },
    channelOrders: {
      type: Array,
      default: () => []
    },
    channelLogs: {
      type: Array,
      default: () => []
    },
    notifyMessages: {
      type: Array,
      default: () => []
    },
    notifyHistories: {
      type: Array,
      default: () => []
    },
    expandedLogs: {
      type: Array,
      default: () => []
    },
    formattedResponses: {
      type: Object,
      default: () => ({})
    }
  },
  data() {
    return {
      activeSubTab: 'payment-order-info',
      collapsedCoreOrders: [],
      // channel logs折叠状态
      collapsedChannelLogs: [],
      // 所有channel logs的折叠状态
      allChannelLogsCollapsed: false
    };
  },
  mounted() {
    // 组件挂载时的逻辑
    console.log('PaymentOrderTab组件挂载 - paymentIndex:', this.paymentIndex);
    console.log('PaymentOrderTab组件挂载 - 通知消息数量:', this.notifyMessages.length);
    console.log('PaymentOrderTab组件挂载 - 通知历史数量:', this.notifyHistories.length);
    if (this.notifyMessages.length > 0) {
      console.log('PaymentOrderTab - 通知消息详情:', this.notifyMessages);
    }
  },
  computed: {
    allCoreOrdersCollapsed() {
      return this.coreOrders.length > 0 && this.coreOrders.every((_, coreIndex) => 
        this.isCoreOrderCollapsed(coreIndex)
      );
    },
    allNotifyMessagesCollapsed() {
      const filteredMessages = this.getFilteredNotifyMessages();
      return filteredMessages.length > 0 && filteredMessages.every((_, msgIndex) => 
        this.isNotifyMessageCollapsed(msgIndex)
      );
    }
  },
  methods: {
    getStatusClass(status) {
      if (!status) return '';
      
      const statusMap = {
        'SUCCESS': 'status-success',
        'FAILED': 'status-fail',
        'PAYING': 'status-processing',
        'PENDING': 'status-pending',
        'CANCEL': 'status-cancelled',
        'EXPIRED': 'status-cancelled',
        'INIT': 'status-init',
      };
      
      return statusMap[status.toUpperCase()] || '';
    },
    formatAmount(amount, precision, currency) {
      if (!amount && amount !== 0) return '-';
      const precisionNum = precision || 0;
      const formattedAmount = (amount / Math.pow(10, precisionNum)).toString();
      return `${currency || ''} ${formattedAmount}`.trim();
    },
    viewPaymentOrderJson() {
      this.$emit('view-json', this.paymentOrder);
    },
    viewJson(data) {
      this.$emit('view-json', data);
    },
    getChannelOrdersByOrderId(orderId) {
      return this.channelOrders.filter(order => order.order_id === orderId);
    },
    // 获取与当前transaction_id关联的channel logs
    getChannelLogsByTransactionId() {
      // 首先获取当前paymentOrder的transaction_id
      const transactionId = this.paymentOrder.transaction_id;
      
      // 然后获取关联的coreOrders
      const coreOrders = this.coreOrders;
      
      // 从coreOrders中提取所有的order_id
      const orderIds = coreOrders.map(order => order.order_id).filter(Boolean);
      
      // 最后获取与这些order_id关联的channel logs
      return this.channelLogs
        .filter(log => orderIds.includes(log.order_id))
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    },
    // 获取所有channel logs，按时间倒序
    getAllChannelLogs() {
      return this.channelLogs
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    },
    // 获取所有相关的通知消息，为每个消息添加notifyList字段
    // 包含对应的通知历史记录
    getFilteredNotifyMessages() {
      // 获取当前payment order的transaction_id和payment_id
      const transactionId = this.paymentOrder.transaction_id;
      const paymentId = this.paymentOrder.payment_id;
      
      // 为每个通知消息添加notifyList字段，包含其对应的通知历史记录
      const messagesWithHistories = this.notifyMessages
        .filter(message => {
          // 只保留与当前payment order的transaction_id或payment_id相关的通知消息
          return message.biz_id === transactionId || message.biz_id === paymentId;
        })
        .map(message => {
          // 获取该消息对应的通知历史记录
          const notifyList = this.getNotifyHistoriesByMsgId(message.id);
          
          // 返回包含notifyList字段的新消息对象
          return {
            ...message,
            notifyList: notifyList
          };
        })
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      
      return messagesWithHistories;
    },
    // 获取与通知消息关联的历史记录
    // 严格按照文档中的关联关系：notify_history.msg_id 和 notify_msg.id 关联
    getNotifyHistoriesByMsgId(msgId) {
      // 严格按照文档要求，使用notify_history.msg_id与notify_msg.id关联
      const filteredHistories = this.notifyHistories
        .filter(history => history.msg_id === msgId)
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
      
      return filteredHistories;
    },
    toggleAllCoreOrders() {
      const allCollapsed = this.allCoreOrdersCollapsed;
      
      // 清除所有折叠状态
      this.collapsedCoreOrders = [];
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed && this.coreOrders.length > 0) {
        this.coreOrders.forEach((_, index) => {
          this.collapsedCoreOrders.push(index);
        });
      }
    },
    toggleCoreOrder(paymentIndex, coreIndex) {
      const idx = this.collapsedCoreOrders.indexOf(coreIndex);
      if (idx > -1) {
        this.collapsedCoreOrders.splice(idx, 1);
      } else {
        this.collapsedCoreOrders.push(coreIndex);
      }
    },
    toggleAllNotifyMessages() {
      // 通知父组件切换所有通知消息的折叠状态
      console.log('toggleAllNotifyMessages clicked');
      // 模拟实现折叠/展开所有通知消息的功能
      const allCollapsed = this.allNotifyMessagesCollapsed;
      
      // 遍历所有通知消息并切换折叠状态
      this.notifyMessages.forEach((_, msgIndex) => {
        const currentState = this.isNotifyMessageCollapsed(msgIndex);
        // 如果当前是展开状态且需要折叠所有，或者当前是折叠状态且需要展开所有，则触发切换
        if ((!allCollapsed && !currentState) || (allCollapsed && currentState)) {
          this.$emit('toggle-notify-message', this.paymentIndex, msgIndex);
        }
      });
    },
    // 切换所有channel logs的折叠状态
    toggleAllChannelLogs() {
      // 切换状态
      this.allChannelLogsCollapsed = !this.allChannelLogsCollapsed;
      
      // 清除所有折叠状态
      this.collapsedChannelLogs = [];
      
      // 如果需要折叠所有，则添加所有索引
      if (this.allChannelLogsCollapsed) {
        const logs = this.getChannelLogsByTransactionId();
        logs.forEach((_, index) => {
          this.collapsedChannelLogs.push(index);
        });
      }
    },
    // 检查单个channel log是否折叠
    isChannelLogCollapsed(logIndex) {
      return this.collapsedChannelLogs.includes(logIndex);
    },
    // 切换单个channel log的折叠状态
    toggleChannelLog(logIndex) {
      const idx = this.collapsedChannelLogs.indexOf(logIndex);
      if (idx > -1) {
        this.collapsedChannelLogs.splice(idx, 1);
      } else {
        this.collapsedChannelLogs.push(logIndex);
      }
    },
    isCoreOrderCollapsed(coreIndex) {
      return this.collapsedCoreOrders.includes(coreIndex);
    },
    // 检查通知消息是否折叠
    isNotifyMessageCollapsed(msgIndex) {
      // 使用与PayDataDetail.vue中一致的索引格式
      if (this.$parent && typeof this.$parent.isNotifyMessageCollapsed === 'function') {
        return this.$parent.isNotifyMessageCollapsed(this.paymentIndex, msgIndex);
      }
      // 默认展开
      return false;
    },
    toggleLogExpand(index, response) {
      // 确保正确地将事件传递给父组件PayDataDetail.vue
      this.$emit('toggle-log-expand', index, response);
    }
  }
};
</script>

<style scoped>
.payment-order-tab {
  padding: 10px;
}

.detail-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.section-title {
  margin: 0;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
  border-left: 4px solid #409eff;
  padding-left: 10px;
}

.sub-tabs {
  margin-top: 20px;
}

.title-with-toggle {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.title-with-toggle:hover {
  color: #409EFF;
  background-color: #ecf5ff;
}

.toggle-icon {
  margin-right: 5px;
  font-size: 12px;
  transition: transform 0.2s;
}

.subsection-title {
  margin: 20px 0 10px 0;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
  background-color: #f0f9ff;
  border-left: 4px solid #409EFF;
}

.subsection-title:hover {
  color: #409EFF;
  background-color: #ecf5ff;
  border-left-color: #66b1ff;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 15px;
  margin-top: 10px;
}

.info-item {
  display: flex;
  align-items: flex-start;
}

.info-item label {
  font-weight: normal;
  margin-right: 10px;
  min-width: 100px;
  color: #909399;
  font-size: 14px;
}

.info-item span {
  color: #303133;
  word-break: break-all;
  font-size: 14px;
}

.empty-message {
  padding: 40px 0;
  text-align: center;
}
</style>