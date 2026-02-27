<template>
  <div class="pay-data-list">
    <!-- <h1>支付数据查询</h1> -->
    
    <!-- Tab页切换 -->
    <el-tabs v-model:model-value="activeTab" type="border-card" class="tab-container">
      <!-- 原有的交易查询tab -->
      <el-tab-pane label="交易查询" name="trade-query">
        <!-- 查询条件表单 -->
        <div class="query-form">
          <el-form :inline="true" :model="queryForm" class="demo-form-inline">
            <el-form-item label="交易ID">
              <el-input v-model="queryForm.trade_id" placeholder="请输入交易ID" />
            </el-form-item>
            <el-form-item label="外部交易ID">
              <el-input v-model="queryForm.out_trade_id" placeholder="请输入外部交易ID" />
            </el-form-item>
            <el-form-item label="用户标识">
              <el-input v-model="queryForm.user_identity" placeholder="请输入用户标识" />
            </el-form-item>
            <el-form-item label="交易状态">
              <!-- 简化交易状态下拉框 -->
              <el-select v-model="queryForm.trade_status" placeholder="请选择交易状态" style="width: 200px;">
                <el-option label="全部" value="" />
                <el-option label="未知" value="UNKNOWN" />
                <el-option label="初始" value="INIT" />
                <el-option label="处理中" value="PAYING" />
                <el-option label="成功" value="SUCCESS" />
                <el-option label="失败" value="FAILED" />
                <el-option label="已取消" value="CANCEL" />
                <el-option label="已过期" value="EXPIRED" />
              </el-select>
            </el-form-item>
            <el-form-item label="国家">
              <el-input v-model="queryForm.region" placeholder="请输入国家" />
            </el-form-item>
            <el-form-item label="支付方式">
              <el-select v-model="queryForm.payment_method" placeholder="请选择支付方式" style="width: 200px;">
                <el-option label="全部" value="" />
                <el-option label="STCPay" value="STCPay" />
                <el-option label="BankCard" value="BankCard" />
                <el-option label="Wallet" value="Wallet" />
                <el-option label="ApplePay" value="ApplePay" />
              </el-select>
            </el-form-item>
            <el-form-item label="支付方式ID">
              <el-input v-model="queryForm.payment_method_id" placeholder="支付方式ID" />
            </el-form-item>
       
            <el-form-item label="查询月份">
              <!-- 简化月份下拉框 -->
              <el-select v-model="queryForm.month" placeholder="请选择查询月份" style="width: 200px;">
                <el-option v-for="month in monthOptions" :key="month" :label="month" :value="month" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleQuery">查询</el-button>
              <el-button @click="resetQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <!-- 结果表格 -->
        <div class="result-table" style="padding: 0;">
          <el-table
            :data="tableData"
            :loading="loading"
            style="width: 100% !important"
            border
            stripe
            fit
          >
            <el-table-column prop="trade_id" label="交易ID" width="300" show-overflow-tooltip>
              <template #default="scope">
                <el-button type="text" @click="viewDetail(scope.row)">{{ scope.row.trade_id }}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="out_trade_id" label="外部交易ID" show-overflow-tooltip />
            <el-table-column prop="user_identity" label="用户标识" width="150" show-overflow-tooltip />
            <el-table-column prop="trade_status" label="交易状态" width="150" >
              <template #default="scope">
                <span :class="getStatusClass(scope.row.trade_status)">{{ scope.row.trade_status }}</span>
              </template>
            </el-table-column>
            <el-table-column label="支付单号" show-overflow-tooltip>
              <template #default="scope">
                {{ getPayTransactionId(scope.row.extra) }}
              </template>
            </el-table-column>
            <el-table-column prop="amount" width="120" label="金额" >
              <template #default="scope">
                {{ formatAmount(scope.row.amount, scope.row.trade_precision, scope.row.currency) }}
              </template>
            </el-table-column>
            <el-table-column prop="region" label="国家" width="80" show-overflow-tooltip  />

            <el-table-column prop="create_time_utc" label="创建时间" />
          </el-table>
          
          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="queryForm.page"
              v-model:page-size="queryForm.page_size"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="total"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </div>
      </el-tab-pane>
      
      <!-- 新增的渠道订单查询tab -->
      <el-tab-pane label="渠道订单查询" name="channel-query">
        <!-- 渠道订单查询条件表单 -->
        <div class="query-form">
          <el-form :inline="true" :model="channelQueryForm" class="demo-form-inline">
            <el-form-item label="订单ID">
              <el-input v-model="channelQueryForm.order_id" placeholder="支付订单" />
            </el-form-item>
             <el-form-item label="渠道标识">
              <el-input v-model="channelQueryForm.channel_id" placeholder="可模糊查询" />
            </el-form-item>
            <el-form-item label="用户标识">
              <el-input v-model="channelQueryForm.user_identity" placeholder="用户标识" />
            </el-form-item>
            <el-form-item label="支付方式">
              <el-select v-model="channelQueryForm.payment_method" placeholder="支付方式" style="width: 200px;">
                <el-option label="全部" value="" />
                <el-option label="STCPay" value="STCPay" />
                <el-option label="BankCard" value="BankCard" />
                <el-option label="Wallet" value="Wallet" />
                <el-option label="ApplePay" value="ApplePay" />
              </el-select>
            </el-form-item>
            <el-form-item label="支付方式ID">
              <el-input v-model="channelQueryForm.payment_method_id" placeholder="支付方式ID" />
            </el-form-item>
            
            <el-form-item label="订单状态">
              <el-select v-model="channelQueryForm.status" placeholder="订单状态" style="width: 200px;">
                <el-option label="全部" value="" />
                <el-option label="未知" value="UNKNOWN" />
                <el-option label="初始" value="INIT" />
                <el-option label="处理中" value="PAYING" />
                <el-option label="成功" value="SUCCESS" />
                <el-option label="失败" value="FAILED" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="查询月份">
              <!-- 简化月份下拉框 -->
              <el-select v-model="channelQueryForm.month" placeholder="请选择查询月份" style="width: 200px;">
                <el-option v-for="month in monthOptions" :key="month" :label="month" :value="month" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="handleChannelQuery">查询</el-button>
              <el-button @click="resetChannelQuery">重置</el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <!-- 渠道订单结果表格 -->
        <div class="result-table">
          <el-table
            :data="channelTableData"
            :loading="channelLoading"
            style="width: 100%"
            border
            stripe
            fit
          >
            <el-table-column prop="order_id" label="订单ID" width="280" show-overflow-tooltip>
              <template #default="scope">
                <el-button type="text" @click="viewChannelDetail(scope.row)">{{ scope.row.order_id }}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="channel_id" label="渠道标识" width="120" show-overflow-tooltip />
            <el-table-column prop="user_identity" label="用户标识" width="120" show-overflow-tooltip />
            <el-table-column prop="payment_method" label="支付方式" width="100" />
            <el-table-column prop="payment_method_id" label="支付方式ID" width="150" show-overflow-tooltip />
                                    <el-table-column prop="region" label="国家"  width="80"/>

            <el-table-column prop="amount" label="金额" width="100">
              <template #default="scope">
                {{ formatAmount(scope.row.amount, scope.row.trade_precision, scope.row.currency) }}
              </template>
            </el-table-column>

            <el-table-column prop="instrument_id" label="支付工具ID" width="250" />

            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <span :class="getStatusClass(scope.row.status)">{{ scope.row.status }}</span>
              </template>
            </el-table-column>
                        <el-table-column prop="channel_code" label="渠道错误码"  />
                        <el-table-column prop="channel_message" label="渠道错误信息" show-overflow-tooltip />
            <el-table-column prop="create_time" label="创建时间" width="150" />
          </el-table>
          
          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="channelQueryForm.page"
              v-model:page-size="channelQueryForm.page_size"
              :page-sizes="[10, 20, 50, 100]"
              layout="total, sizes, prev, pager, next, jumper"
              :total="channelTotal"
              @size-change="handleChannelSizeChange"
              @current-change="handleChannelCurrentChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
    
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

export default {
  name: 'PayDataList',
  data() {
    return {
      // Tab页状态
      activeTab: 'trade-query',
      
      // 原有的交易查询表单
      queryForm: {
        trade_id: '',
        out_trade_id: '',
        user_identity: '',
        trade_status: '',
        region: '',
        payment_method: '',
        payment_method_id: '',
        month: '',
        page: 1,
        page_size: 20
      },
      tableData: [],
      total: 0,
      loading: false,
      
      // 新增渠道订单查询表单
      channelQueryForm: {
        channel_id: '',
        order_id: '',
        user_identity: '',
        payment_method: '',
        payment_method_id: '',
        status: '',
        month: '',
        page: 1,
        page_size: 20
      },
      channelTableData: [],
      channelTotal: 0,
      channelLoading: false,
      
      jsonDialogVisible: false,
      jsonData: {},
      monthOptions: []
    };
  },
  mounted() {
    console.log('PayDataList组件已挂载，开始初始化');
    
    // 检查URL参数中是否有指定的activeTab
    if (this.$route.query.activeTab) {
      this.activeTab = this.$route.query.activeTab;
    }
    
    // 生成月份选项字符串数组
    const currentMonth = new Date();
    const year = currentMonth.getFullYear();
    const month = String(currentMonth.getMonth() + 1).padStart(2, '0');
    const currentMonthStr = `${year}${month}`;
    
    console.log('当前月份:', currentMonthStr);
    
    // 先设置选项数组
    this.monthOptions = [currentMonthStr]; // 先添加当前月份
    
    // 添加剩余5个月
    for (let i = 1; i < 6; i++) {
      const date = new Date(year, currentMonth.getMonth() - i, 1);
      const mYear = date.getFullYear();
      const mMonth = String(date.getMonth() + 1).padStart(2, '0');
      const monthStr = `${mYear}${mMonth}`;
      if (!this.monthOptions.includes(monthStr)) {
        this.monthOptions.push(monthStr);
      }
    }
    
    console.log('生成的月份选项:', this.monthOptions);
    
    // 设置默认月份
    this.queryForm.month = currentMonthStr;
    this.channelQueryForm.month = currentMonthStr;
    
    console.log('默认月份设置完成:', currentMonthStr);
    
    // 从sessionStorage恢复查询条件
    const savedQueryForm = sessionStorage.getItem('payDataQueryForm');
    if (savedQueryForm) {
      console.log('从sessionStorage发现保存的查询条件，尝试恢复');
      try {
        const parsedForm = JSON.parse(savedQueryForm);
        // 确保恢复有效的月份值
        if (this.monthOptions.includes(parsedForm.month)) {
          this.queryForm = parsedForm;
          console.log('查询条件恢复成功:', parsedForm);
          // 恢复后自动查询数据
          this.handleQuery();
        } else {
          console.warn('恢复的月份不在有效选项中，跳过自动查询');
        }
      } catch (e) {
        console.error('恢复查询条件失败:', e);
      }
    }
    
    // 从sessionStorage恢复渠道订单查询条件
    const savedChannelQueryForm = sessionStorage.getItem('payDataChannelQueryForm');
    if (savedChannelQueryForm) {
      console.log('从sessionStorage发现保存的渠道订单查询条件，尝试恢复');
      try {
        const parsedForm = JSON.parse(savedChannelQueryForm);
        // 确保恢复有效的月份值
        if (this.monthOptions.includes(parsedForm.month)) {
          this.channelQueryForm = parsedForm;
          console.log('渠道订单查询条件恢复成功:', parsedForm);
          // 恢复后自动查询数据
          this.handleChannelQuery();
        } else {
          console.warn('恢复的渠道查询月份不在有效选项中，跳过自动查询');
        }
      } catch (e) {
        console.error('恢复渠道订单查询条件失败:', e);
      }
    }
    console.log('组件初始化完成');
  },
  computed: {
    jsonDataStr() {
      return JSON.stringify(this.jsonData, null, 2);
    }
  },
  
  watch: {
    // 监控标签页切换
    activeTab(newVal, oldVal) {
      console.log(`标签页切换: ${oldVal} -> ${newVal}`);
    }
  },

  methods: {

    
    // 原有交易查询方法
      handleQuery() {
        console.log('开始交易查询，参数:', this.queryForm);
        
        if (!this.queryForm.month) {
          this.$message.error('请输入查询月份');
          console.warn('交易查询失败：未选择查询月份');
          return;
        }
        
        this.loading = true;
        payDataApi.queryTradeList(this.queryForm).then(res => {
          console.log('交易查询API返回结果:', res);
          // 保留当前页码，不再重置为1
          // 直接使用当前的queryForm值进行查询，不依赖sessionStorage缓存
          const result = res;
          this.tableData = result.items || [];
          this.total = result.total || 0;
          console.log('交易查询成功，数据量:', (result.items || []).length, '总条数:', (result.total || 0));
          // 只有在查询成功后才保存查询条件
          sessionStorage.setItem('payDataQueryForm', JSON.stringify(this.queryForm));
          console.log('查询条件已保存到sessionStorage');
        }).catch(error => {
          console.error('交易查询失败:', error);
          this.$message.error('查询失败：' + (error.message || '未知错误'));
        }).finally(() => {
          this.loading = false;
          console.log('交易查询请求完成');
        });
      },
    
    // 原有交易查询重置
    resetQuery() {
      console.log('执行交易查询重置操作');
      this.queryForm = {
        trade_id: '',
        out_trade_id: '',
        user_identity: '',
        trade_status: '',
        region: '',
        payment_method: '',
        payment_method_id: '',
        month: new Date().toISOString().slice(0, 7).replace('-', ''),
        page: 1,
        page_size: 20
      };
      // 清空保存的查询条件
      sessionStorage.removeItem('payDataQueryForm');
      console.log('交易查询条件已重置并清空sessionStorage');
    },
    
    // 原有交易查询详情
    viewDetail(row) {
      console.log('查看交易详情，交易ID:', row.trade_id);
      // 保存当前查询状态，确保返回时能恢复
      sessionStorage.setItem('payDataQueryForm', JSON.stringify(this.queryForm));
      console.log('详情跳转前保存查询条件');
      
      this.$router.push({
        path: '/pay-data/detail',
        query: { 
          trade_id: row.trade_id,
          month: this.queryForm.month 
        }
      });
      console.log('跳转到详情页:', '/pay-data/detail', { trade_id: row.trade_id, month: this.queryForm.month });
    },
    
    // 分页大小变化
    handleSizeChange(size) {
      console.log('交易查询分页大小变化:', this.queryForm.page_size, '->', size);
      this.queryForm.page_size = size;
      this.handleQuery();
    },
    
    // 当前页码变化
    handleCurrentChange(current) {
      console.log('交易查询页码变化:', this.queryForm.page, '->', current);
      this.queryForm.page = current;
      this.handleQuery();
    },
    
    // 格式化金额
    formatAmount(amount, precision, currency) {
      try {
        if (!amount && amount !== 0) return '-';
        const precisionNum = precision || 2;
        const formattedAmount = (amount / Math.pow(10, precisionNum)).toFixed(precisionNum);
        const result = `${currency || ''} ${formattedAmount}`.trim();
        return result;
      } catch (error) {
        console.error('金额格式化失败:', error, { amount, precision, currency });
        return '-';
      }
    },
    
    // 解析extra字段获取支付交易ID
    getPayTransactionId(extra) {
      if (!extra) return '-';
      try {
        const parsedExtra = typeof extra === 'string' ? JSON.parse(extra) : extra;
        return parsedExtra.pay_transaction_id || '-';
      } catch (error) {
        console.error('解析extra字段失败:', error, { extra });
        return '-';
      }
    },
    
    // 查看详情
    viewJson(row) {
      console.log('查看JSON详情:', row);
      this.jsonData = row;
      this.jsonDialogVisible = true;
      console.log('详情对话框已打开');
    },
    
    // 新增渠道订单查询方法
      handleChannelQuery() {
        console.log('开始渠道订单查询，参数:', this.channelQueryForm);
        
        if (!this.channelQueryForm.month) {
          this.$message.error('请输入查询月份');
          console.warn('渠道订单查询失败：未选择查询月份');
          return;
        }
        
        this.channelLoading = true;
        payDataApi.queryChannelOrderList(this.channelQueryForm).then(res => {
          console.log('渠道订单查询API返回结果:', res);
          const result = res;
          this.channelTableData = result.items || [];
          this.channelTotal = result.total || 0;
          console.log('渠道订单查询成功，数据量:', (result.items || []).length, '总条数:', (result.total || 0));
          // 保存查询条件
          sessionStorage.setItem('payDataChannelQueryForm', JSON.stringify(this.channelQueryForm));
          console.log('渠道查询条件已保存到sessionStorage');
        }).catch(error => {
          console.error('渠道订单查询失败:', error);
          this.$message.error('渠道订单查询失败：' + (error.message || '未知错误'));
        }).finally(() => {
          this.channelLoading = false;
          console.log('渠道订单查询请求完成');
        });
      },
    
    // 渠道订单查询重置
    resetChannelQuery() {
      console.log('执行渠道订单查询重置操作');
      this.channelQueryForm = {
        channel_id: '',
        order_id: '',
        user_identity: '',
        payment_method: '',
        payment_method_id: '',
        status: '',
        month: new Date().toISOString().slice(0, 7).replace('-', ''),
        page: 1,
        page_size: 20
      };
      // 清空保存的查询条件
      sessionStorage.removeItem('payDataChannelQueryForm');
      console.log('渠道订单查询条件已重置并清空sessionStorage');
    },
    
    // 渠道订单查询详情
    viewChannelDetail(row) {
      console.log('查看渠道订单详情，订单ID:', row.order_id);
      // 保存当前查询状态，确保返回时能恢复
      sessionStorage.setItem('payDataChannelQueryForm', JSON.stringify(this.channelQueryForm));
      console.log('详情跳转前保存渠道查询条件');
      
      this.$router.push({
        path: '/pay-data/channel-detail',
        query: { 
          order_id: row.order_id,
          month: this.channelQueryForm.month 
        }
      });
      console.log('跳转到渠道详情页:', '/pay-data/channel-detail', { order_id: row.order_id, month: this.channelQueryForm.month });
    },
    
    // 渠道订单分页大小变化
    handleChannelSizeChange(size) {
      console.log('渠道订单分页大小变化:', this.channelQueryForm.page_size, '->', size);
      this.channelQueryForm.page_size = size;
      this.handleChannelQuery();
    },
    
    // 渠道订单当前页码变化
    handleChannelCurrentChange(current) {
      console.log('渠道订单页码变化:', this.channelQueryForm.page, '->', current);
      this.channelQueryForm.page = current;
      this.handleChannelQuery();
    },
    getStatusClass(status) {
      if (!status) return '';
      
      const statusMap = {
        'SUCCESS': 'status-success',
        'FAILED': 'status-fail',
        'PAYING': 'status-processing',
        'PENDING': 'status-pending',
        'CANCEL': 'status-cancelled',
        'EXPIRED': 'status-cancelled',
        'INIT': 'status-init'
      };
      
      return statusMap[status.toUpperCase()] || '';
    },
  }
};
</script>

<style scoped>
.pay-data-list {
  padding: 20px;
}

.tab-container {
  margin-bottom: 20px;
}

.query-form {
  margin-bottom: 20px;
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
}

.result-table {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  width: 100%;
  box-sizing: border-box;
}

.result-table :deep(.el-table) {
  width: 100% !important;
}

.result-table :deep(.el-table__header-wrapper),
.result-table :deep(.el-table__body-wrapper) {
  width: 100% !important;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.json-pre {
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 500px;
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

:deep(.status-init) {
  color: #909399;
  background-color: #f0f2f5;
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
</style>