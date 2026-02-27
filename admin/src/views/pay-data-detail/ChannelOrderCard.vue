<template>
  <el-card class="inner-card">
    <template #header>
      <div class="card-header">
        <div class="title-with-toggle" :class="{ collapsed: isCollapsed }" @click="toggleCollapse">
          <el-icon class="toggle-icon"><ChevronDown v-if="!isCollapsed" /><ChevronRight v-else /></el-icon>
          <span>Channel Order #{{ channelIndex + 1 }}</span>
        </div>
        <el-button type="primary" size="small" @click="viewChannelOrderJson">查看详情</el-button>
      </div>
    </template>
    
    <div v-if="!isCollapsed">
      <div class="info-grid">
        <div class="info-item">
          <label>订单ID：</label>
          <span>{{ channelOrder.order_id }}</span>
        </div>
        <div class="info-item">
          <label>渠道订单ID：</label>
          <span>{{ channelOrder.channel_trade_id }}</span>
        </div>
        <div class="info-item">
          <label>渠道：</label>
          <span>{{ channelOrder.channel_id }}</span>
        </div>
        <div class="info-item">
          <label>状态：</label>
          <span :class="getStatusClass(channelOrder.status)">{{ channelOrder.status }}</span>
        </div>
        <div class="info-item">
          <label>支付金额：</label>
          <span>{{ formatAmount(channelOrder.amount, channelOrder.precision, channelOrder.currency) }}</span>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script>
import * as Icons from '@element-plus/icons-vue';

export default {
  name: 'ChannelOrderCard',
  components: {
    ChevronDown: Icons.ArrowDown,
    ChevronRight: Icons.ArrowRight
  },
  props: {
    channelOrder: {
      type: Object,
      required: true
    },
    paymentIndex: {
      type: Number,
      required: true
    },
    coreIndex: {
      type: Number,
      required: true
    },
    channelIndex: {
      type: Number,
      required: true
    },
    isCollapsed: {
      type: Boolean,
      default: false
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
      this.$emit('toggle-collapse', this.paymentIndex, this.coreIndex, this.channelIndex);
    },
    viewChannelOrderJson() {
      this.$emit('view-json', this.channelOrder);
    }
  }
};
</script>

<style scoped>
.inner-card {
  margin-left: 20px;
  margin-bottom: 10px;
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