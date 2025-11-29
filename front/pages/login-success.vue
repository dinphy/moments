<template>
  <div class="login-success-container">
    <div class="loading-spinner"></div>
    <p>登录成功，正在处理...</p>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useGlobalState } from '~/store';

const router = useRouter();
const global = useGlobalState();

onMounted(() => {
  // 从URL参数中提取token、userId和username
  const urlParams = new URLSearchParams(window.location.search);
  const token = urlParams.get('token');
  const userId = urlParams.get('userId');
  const username = urlParams.get('username');
  
  if (token) {
    // 将用户信息设置到global.value.userinfo中
    global.value.userinfo = {
      token: token,
      id: userId ? parseInt(userId) : 0,
      username: username || ''
    };
    console.log('用户信息已设置到global.value.userinfo中');
    
    // 延迟一小段时间再跳转，确保用户可以看到过渡效果
    setTimeout(() => {
      // 跳转到首页并刷新页面
      window.location.href = '/';
    }, 1000);
  } else {
    console.error('URL中未找到token参数');
    // 如果没有token，也跳转到首页并刷新
    setTimeout(() => {
      window.location.href = '/';
    }, 1000);
  }
});
</script>

<style scoped>
.login-success-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: #f5f5f5;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e0e0e0;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

p {
  color: #333;
  font-size: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>