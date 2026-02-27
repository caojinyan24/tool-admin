<template>
  <el-card class="nested-card">
    <template #header>
      <div class="card-header">
        <div class="title-with-toggle" :class="{ collapsed: isCollapsed }" @click="toggleCollapse">
          <el-icon class="toggle-icon"><ChevronDown v-if="!isCollapsed" /><ChevronRight v-else /></el-icon>
          <span>Core Order #{{ index + 1 }}</span>
        </div>
        <el-button type="primary" size="small" @click="viewCoreOrderJson">查看详情</el-button>
      </div>
    </template>
    
    <div v-if="!isCollapsed">
      <div class="info-grid">
        <div class="info-item">
          <label>订单ID：</label>
          <span>{{ coreOrder.order_id }}</span>
        </div>
        <div class="info-item">
          <label>支付订单ID：</label>
          <span>{{ coreOrder.trade_id }}</span>
        </div>
        <div class="info-item">
          <label>支付方式：</label>
          <span>{{ coreOrder.payment_method }}</span>
        </div>
        <div class="info-item">
          <label>渠道：</label>
          <span>{{ coreOrder.channel_id }}</span>
        </div>
        <div class="info-item">
          <label>支付金额：</label>
          <span>{{ formatAmount(coreOrder.payment_amount, coreOrder.payment_currency_precision, coreOrder.payment_currency) }}</span>
        </div>
        <div class="info-item">
          <label>订单状态：</label>
          <span :class="getStatusClass(coreOrder.status)">{{ coreOrder.status }}</span>
        </div>
      </div>
      
      <!-- 渠道订单 -->
      <div v-if="channelOrders.length > 0">
        <div class="title-with-toggle" @click="toggleAllChannelOrders">
          <el-icon class="toggle-icon"><ChevronDown v-if="!allChannelOrdersCollapsed" /><ChevronRight v-else /></el-icon>
          <h5 class="subsection-title">Channel Orders ({{ channelOrders.length }})</h5>
        </div>
        <div v-if="!allChannelOrdersCollapsed">
          <channel-order-card
            v-for="(channelOrder, channelIndex) in channelOrders"
            :key="channelOrder.channel_order_id || channelIndex"
            :channel-order="channelOrder"
            :payment-index="paymentIndex"
            :core-index="index"
            :channel-index="channelIndex"
            :is-collapsed="isChannelOrderCollapsed(channelIndex)"
            @toggle-collapse="toggleChannelOrder"
            @view-json="viewChannelOrderJson"
          />
        </div>
      </div>
      
      <!-- 渠道日志 -->
      <div v-if="channelLogs.length > 0">
        <div class="title-with-toggle" @click="toggleAllChannelLogs">
          <el-icon class="toggle-icon"><ChevronDown v-if="!allChannelLogsCollapsed" /><ChevronRight v-else /></el-icon>
          <h5 class="subsection-title">Channel Logs ({{ channelLogs.length }})</h5>
        </div>
        <div v-if="!allChannelLogsCollapsed">
          <channel-log-card
            v-for="(channelLog, logIndex) in channelLogs"
            :key="channelLog.log_id || logIndex"
            :channel-log="channelLog"
            :payment-index="paymentIndex"
            :core-index="index"
            :log-index="logIndex"
              :is-collapsed="isChannelLogCollapsed(logIndex)"
              :expanded-logs="expandedLogs"
              :formatted-responses="formattedResponses"
              @toggle-collapse="toggleChannelLog"
              @view-json="viewChannelLogJson"
              @toggle-log-expand="handleToggleLogExpand"
          />
        </div>
      </div>
      
      <!-- 通知消息 -->
      <div v-if="getNotifyMessagesByCoreOrderId().length > 0">
        <div class="title-with-toggle" @click="toggleAllNotifyMessages">
          <el-icon class="toggle-icon"><ChevronDown v-if="!allNotifyMessagesCollapsed" /><ChevronRight v-else /></el-icon>
          <h5 class="subsection-title">通知消息 ({{ getNotifyMessagesByCoreOrderId().length }})</h5>
        </div>
        <div v-if="!allNotifyMessagesCollapsed">
          <notify-message-card
            v-for="(notifyMessage, msgIndex) in getNotifyMessagesByCoreOrderId()"
            :key="notifyMessage.id || msgIndex"
            :notify-message="notifyMessage"
            :payment-index="paymentIndex"
            :msg-index="msgIndex"
            :notify-histories="notifyMessage.notifyList || []"
            :is-collapsed="collapsedNotifyMessages.includes(msgIndex)"
            :expanded-logs="expandedLogs"
            :formatted-responses="formattedResponses"
            @toggle-collapse="toggleNotifyMessage"
            @toggle-history-collapse="toggleNotifyHistory"
            @view-json="viewJson"
            @toggle-log-expand="handleToggleLogExpand"
          />
        </div>
      </div>
    </div>
  </el-card>
</template>

<script>
import ChannelOrderCard from './ChannelOrderCard.vue';
import ChannelLogCard from './ChannelLogCard.vue';
import NotifyMessageCard from './NotifyMessageCard.vue';
import * as Icons from '@element-plus/icons-vue';

export default {
  name: 'CoreOrderCard',
  components: {
    ChannelOrderCard,
    ChannelLogCard,
    NotifyMessageCard,
    ChevronDown: Icons.ArrowDown,
    ChevronRight: Icons.ArrowRight
  },
  data() {
    return {
      collapsedChannelOrders: [],
      collapsedChannelLogs: [],
      collapsedNotifyMessages: [],
      collapsedNotifyHistories: []
    };
  },
  props: {
    coreOrder: {
      type: Object,
      required: true
    },
    paymentIndex: {
      type: Number,
      required: true
    },
    index: {
      type: Number,
      required: true
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
    collapsedState: {
      type: String,
      default: ''
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
  computed: {
    isCollapsed() {
      return this.collapsedState === 'collapsed';
    },
    allChannelOrdersCollapsed() {
      return this.channelOrders.length > 0 && this.channelOrders.every((_, channelIndex) => 
        this.isChannelOrderCollapsed(channelIndex)
      );
    },
    allChannelLogsCollapsed() {
      return this.channelLogs.length > 0 && this.channelLogs.every((_, logIndex) => 
        this.isChannelLogCollapsed(logIndex)
      );
    },
    allNotifyMessagesCollapsed() {
      const notifyMessages = this.getNotifyMessagesByCoreOrderId();
      return notifyMessages.length > 0 && notifyMessages.every((_, msgIndex) => 
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
    toggleCollapse() {
      this.$emit('toggle-collapse', this.paymentIndex, this.index);
    },
    viewCoreOrderJson() {
      this.$emit('view-json', this.coreOrder);
    },
    toggleAllChannelOrders() {
      const allCollapsed = this.allChannelOrdersCollapsed;
      
      // 清除所有折叠状态
      this.collapsedChannelOrders = [];
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed && this.channelOrders.length > 0) {
        this.channelOrders.forEach((_, index) => {
          this.collapsedChannelOrders.push(index);
        });
      }
    },
    toggleAllChannelLogs() {
      const allCollapsed = this.allChannelLogsCollapsed;
      
      // 清除所有折叠状态
      this.collapsedChannelLogs = [];
      
      // 如果不全折叠，则折叠所有；否则展开所有
      if (!allCollapsed && this.channelLogs.length > 0) {
        this.channelLogs.forEach((_, index) => {
          this.collapsedChannelLogs.push(index);
        });
      }
    },
    toggleChannelOrder(paymentIndex, coreIndex, channelIndex) {
      const idx = this.collapsedChannelOrders.indexOf(channelIndex);
      if (idx > -1) {
        this.collapsedChannelOrders.splice(idx, 1);
      } else {
        this.collapsedChannelOrders.push(channelIndex);
      }
      
      // 同时通知父组件
      this.$emit('toggle-channel-order', paymentIndex, coreIndex, channelIndex);
    },
    toggleChannelLog(paymentIndex, coreIndex, logIndex) {
      const idx = this.collapsedChannelLogs.indexOf(logIndex);
      if (idx > -1) {
        this.collapsedChannelLogs.splice(idx, 1);
      } else {
        this.collapsedChannelLogs.push(logIndex);
      }
      
      // 同时通知父组件
      this.$emit('toggle-channel-log', paymentIndex, coreIndex, logIndex);
    },
    isChannelOrderCollapsed(channelIndex) {
      return this.collapsedChannelOrders.includes(channelIndex);
    },
    isChannelLogCollapsed(logIndex) {
      return this.collapsedChannelLogs.includes(logIndex);
    },
    viewChannelOrderJson(channelOrder) {
      this.$emit('view-json', channelOrder);
    },
    viewChannelLogJson(channelLog) {
      this.$emit('view-json', channelLog);
    },
    handleToggleLogExpand(index, response) {
      // 转发事件到父组件
      this.$emit('toggle-log-expand', index, response);
    },
    // 获取与核心订单关联的通知消息，为每个消息添加notifyList字段
    // 严格按照文档中的关联关系：notify_msg.biz_id 和 core_orders.order_id 关联
    getNotifyMessagesByCoreOrderId() {
      console.log('获取核心订单通知消息 - orderId:', this.coreOrder.order_id);
      
      if (!this.notifyMessages || this.notifyMessages.length === 0) {
        console.log('没有通知消息数据');
        return [];
      }
      
      // 为每个通知消息添加notifyList字段，包含其对应的通知历史记录
      const messagesWithHistories = this.notifyMessages
        .filter(msg => msg.biz_id === this.coreOrder.order_id)
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
    getNotifyHistoriesByMsgId(msgId) {
      // 严格按照文档要求，使用notify_history.msg_id与notify_msg.id关联
      return this.notifyHistories
        .filter(history => history.msg_id === msgId)
        .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
    },
    // 切换所有通知消息的折叠状态
    toggleAllNotifyMessages() {
      const notifyMessages = this.getNotifyMessagesByCoreOrderId();
      if (notifyMessages.length === 0) return;
      
      if (this.allNotifyMessagesCollapsed) {
        // 如果全部折叠，则全部展开
        this.collapsedNotifyMessages = [];
      } else {
        // 如果不全部折叠，则全部折叠
        this.collapsedNotifyMessages = notifyMessages.map((_, index) => index);
      }
    },
    // 检查通知消息是否折叠
    isNotifyMessageCollapsed(msgIndex) {
      return this.collapsedNotifyMessages.includes(msgIndex);
    },
    // 切换通知消息的折叠状态
    toggleNotifyMessage(paymentIndex, msgIndex) {
      if (this.collapsedNotifyMessages.includes(msgIndex)) {
        this.collapsedNotifyMessages = this.collapsedNotifyMessages.filter(index => index !== msgIndex);
      } else {
        this.collapsedNotifyMessages.push(msgIndex);
      }
    },
    // 切换通知历史的折叠状态
    toggleNotifyHistory(paymentIndex, msgIndex, historyIndex) {
      const key = `${msgIndex}_${historyIndex}`;
      if (this.collapsedNotifyHistories.includes(key)) {
        this.collapsedNotifyHistories = this.collapsedNotifyHistories.filter(item => item !== key);
      } else {
        this.collapsedNotifyHistories.push(key);
      }
    }
  }
};
</script>

<style scoped>
.nested-card {
  margin-left: 20px;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.title-with-toggle.collapsed {
  font-weight: normal;
  color: #606266;
  background-color: #f4f4f5;
  border-left: 3px solid #dcdfe6;
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
</style>