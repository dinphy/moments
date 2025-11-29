<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-100" style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;">
    <div class="w-10 h-10 border-4 border-gray-300 border-t-blue-500 rounded-full animate-spin mb-5"></div>
    <p class="text-gray-800 text-base">登录成功，正在处理...000</p>
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

<style scoped></style>
