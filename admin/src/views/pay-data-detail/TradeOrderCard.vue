<template>
  <el-card class="detail-card">
    <template #header>
      <div class="card-header">
        <h3 class="section-title">Trade Order</h3>
        <el-button type="primary" size="small" @click="viewTradeOrderJson">查看详情</el-button>
      </div>
    </template>
    
    <div class="info-grid">
      <div class="info-item">
        <label>交易ID：</label>
        <span>{{ tradeOrder.trade_id }}</span>
      </div>
      <div class="info-item">
        <label>创建时间：</label>
        <span>{{ tradeOrder.create_time_utc }}</span>
      </div>
      <div class="info-item">
        <label>交易金额：</label>
        <span>{{ formatAmount(tradeOrder.amount, tradeOrder.trade_precision, tradeOrder.currency) }}</span>
      </div>
      <div class="info-item">
        <label>交易状态：</label>
        <span :class="getStatusClass(tradeOrder.status)">{{ tradeOrder.trade_status }}</span>
      </div>
    </div>
  </el-card>
</template>

<script>
export default {
  name: 'TradeOrderCard',
  props: {
    tradeOrder: {
      type: Object,
      required: true
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
    viewTradeOrderJson() {
      this.$emit('view-json', this.tradeOrder);
    },
     formatAmount(amount, precision, currency) {
      if (!amount && amount !== 0) return '-';
      const precisionNum = precision || 0;
      const formattedAmount = (amount / Math.pow(10, precisionNum)).toString();
      return `${currency || ''} ${formattedAmount}`.trim();
    },
  }
};
</script>

<style scoped>
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