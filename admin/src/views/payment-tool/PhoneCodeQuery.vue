<template>
  <div class="phone-code-query">
    <el-form ref="queryForm" :model="queryForm" :rules="rules" label-width="120px" class="query-form">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="请选择查询场景" prop="scene">
            <el-select v-model="queryForm.scene" placeholder="请选择场景" clearable>
              <el-option 
                v-for="option in sceneOptions" 
                :key="option.value" 
                :label="option.label" 
                :value="option.value"
              />
            </el-select>
            <div v-if="$refs.queryForm && $refs.queryForm.errors && $refs.queryForm.errors.scene" class="el-form-item__error">
              {{ $refs.queryForm.errors.scene[0] }}
            </div>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="请选择测试账号" prop="userName">
            <el-select v-model="queryForm.userName" placeholder="请选择账号" clearable>
              <el-option 
                v-for="account in accountOptions" 
                :key="account.userName" 
                :label="account.userName" 
                :value="account.userName"
              />
            </el-select>
            <div v-if="$refs.queryForm && $refs.queryForm.errors && $refs.queryForm.errors.userName" class="el-form-item__error">
              {{ $refs.queryForm.errors.userName[0] }}
            </div>
          </el-form-item>
        </el-col>
      </el-row>
      
      <el-form-item>
        <el-button type="primary" @click="handleQuery" :loading="loading">查询验证码</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 查询结果展示 -->
    <div v-if="resultVisible" class="result-section">
      <el-card class="result-card">
        <template #header>
          <div class="card-header">
            <span>查询结果</span>
          </div>
        </template>
        
        <div v-if="result.error" class="error-message">
          <el-alert
            title="查询失败"
            :description="result.error"
            type="error"
            show-icon
            :closable="false"
          />
        </div>
        
        <div v-else-if="result.code" class="success-message">
          <el-alert
            title="查询成功"
            type="success"
            show-icon
            :closable="false"
          />
          <div class="code-display">
            <h3>短信验证码：</h3>
            <div class="code-value">{{ result.code }}</div>
          </div>
          <div class="account-info">
            <p>账号信息：{{ selectedAccount?.userName }}</p>
            <p>手机号：{{ selectedAccount?.phone }}</p>
            <p>环境：{{ selectedAccount?.env }}</p>
            <p>地区：{{ selectedAccount?.region }}</p>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { getPhoneCode } from '@/api/payment-tool';

export default {
  name: 'PhoneCodeQuery',
  data() {
    return {
      loading: false,
      resultVisible: false,
      result: {},
      selectedAccount: null,
      queryForm: {
        scene: '',
        userName: ''
      },
      rules: {
        scene: [
          { required: true, message: '请选择查询场景', trigger: 'change' }
        ],
        userName: [
          { required: true, message: '请选择测试账号', trigger: 'change' }
        ]
      },
      sceneOptions: [
        { label: 'register', value: 'register' },
        { label: 'username_login', value: 'username_login' },
        { label: 'stc_bind_withdraw_phone', value: 'stc_bind_withdraw_phone' },
        { label: 'stcpay_sercurity', value: 'stcpay_sercurity' }
      ],
      accountOptions: [
        { userName: 'dxf001', phone: '13062263261', env: 'test', region: 'ID', comment: '印度尼西亚-不支持 Apple pay' },
        { userName: 'dxf002', phone: '19185812465', env: 'test', region: 'JP', comment: '日本-不支持 Apple pay' },
        { userName: 'dxf003', phone: '19085884168', env: 'test', region: 'TH', comment: '泰国-不支持 Apple pay' },
        { userName: 'dxf004', phone: '16017925281', env: 'test', region: '土耳其' },
        { userName: 'dxf005', phone: '17024632223', env: 'test', region: '巴西' },
        { userName: 'dxf006', phone: '18657083142', env: 'test', region: 'HK', comment: '香港-不支持 Apple pay' },
        { userName: 'dxf007', phone: '19205910626', env: 'test', region: 'US', comment: '美国-支持 Apple pay' },
        { userName: 'dxf009', phone: '15613986578', env: 'test', region: 'EG', comment: '埃及-钱包' },
        { userName: 'dxf010', phone: '13235632303', env: 'test', region: 'PK', comment: '巴基斯坦-钱包' },
        { userName: 'xftest05', phone: '17702881360', env: 'test', region: 'TR' },
        { userName: 'xftest06', phone: '13098148410', env: 'test', region: 'BR' },
        { userName: 'xftest07', phone: '966506678012', env: 'test', region: '巴基斯坦' },
        { userName: 'xftest09', phone: '201198765432', env: 'test', region: '科威特' }
      ]
    };
  },
  methods: {
      handleQuery() {
        this.$refs.queryForm.validate((valid) => {
          if (valid) {
            this.loading = true;
            this.selectedAccount = this.accountOptions.find(acc => acc.userName === this.queryForm.userName);
            
            // 调用API查询短信验证码
            getPhoneCode({
              scene: this.queryForm.scene,
              userName: this.queryForm.userName
            }).then(res => {
              this.loading = false;
              this.resultVisible = true;
              if (res.code === 200) {
                this.result = { code: res.data };
                this.$message.success('查询成功');
              } else {
                this.result = { error: res.message || '查询失败，请稍后重试' };
                this.$message.error(res.message || '查询失败');
              }
            }).catch(err => {
              this.loading = false;
              this.resultVisible = true;
              const errorMsg = err.message || '网络错误，请稍后重试';
              this.result = { error: errorMsg };
              this.$message.error(errorMsg);
            });
          }
        });
      },
    resetForm() {
      this.$refs.queryForm.resetFields();
      this.resultVisible = false;
      this.result = {};
      this.selectedAccount = null;
    }
  }
};
</script>

<style scoped>
.phone-code-query {
  padding: 20px;
}

.query-form {
  margin-bottom: 20px;
}

.result-section {
  margin-top: 30px;
}

.result-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.error-message,
.success-message {
  margin-top: 20px;
}

.code-display {
  margin-top: 20px;
  text-align: center;
}

.code-value {
  font-size: 36px;
  font-weight: bold;
  color: #1890ff;
  padding: 20px;
  background-color: #f0f9ff;
  border-radius: 8px;
  margin: 20px 0;
  letter-spacing: 4px;
}

.account-info {
  margin-top: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.account-info p {
  margin: 5px 0;
  color: #606266;
}
</style>