<template>
  <div class="payment-tool">
    <!-- <h1>支付工具</h1> -->
    
    <el-tabs v-model="activeTab" type="border-card" class="test-tabs">
      <el-tab-pane label="查询短信验证码" name="phone-code">
        <PhoneCodeQuery />
      </el-tab-pane>
      <el-tab-pane label="其他功能" name="other">
        <div class="empty-tab">
          <p>其他功能开发中...</p>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import PhoneCodeQuery from './payment-tool/PhoneCodeQuery.vue';

export default {
  name: 'PaymentTest',
  components: {
    PhoneCodeQuery
  },
  data() {
    return {
      activeTab: 'phone-code'
    };
  }
};
</script>

<style scoped>
.payment-tool {
  padding: 20px;
}

h1 {
  margin-bottom: 20px;
  color: #303133;
}

.test-tabs {
  margin-top: 20px;
}

.empty-tab {
  padding: 60px 0;
  text-align: center;
  color: #909399;
}
</style>