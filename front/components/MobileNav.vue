<template>
  <div class="md:hidden fixed bottom-0 left-0 right-0 bg-[#F5F5F5] dark:bg-[#202020] backdrop-blur-md z-10">
    <div class="flex justify-around items-center">
      <!-- 动态 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-[#07C160]': $route.path === '/' }"
        @click="handleHomeClick"
      >
        <UIcon name="i-carbon-activity" class="w-6 h-6" />
        <span class="text-xs mt-0.5">动态</span>
      </div>

      <!-- 发现 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-[#07C160]': $route.path === '/discover' }"
        @click="navigate('/discover')"
      >
        <UIcon :name="$route.path === '/discover' ? 'i-weui-discover-filled' : 'i-weui-discover-outlined'" class="w-6 h-6" />
        <span class="text-xs mt-0.5">发现</span>
      </div>

      <!-- 我 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-[#07C160]': $route.path === '/user/settings' }"
        @click="handleUserClick"
      >
        <UIcon :name="$route.path === '/user/settings' ? 'i-weui-me-filled' : 'i-weui-me-outlined'" class="w-6 h-6" />
        <span class="text-xs mt-0.5">我</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";

const global = useGlobalState();
const loginReg = useState<boolean>("loginReg", () => false);

const navigate = async (url: string) => {
  await navigateTo(url);
};

const handleUserClick = async () => {
  if (global.value.userinfo.token) {
    await navigateTo("/user/settings");
  } else {
    loginReg.value = true;
  }
};

const handleHomeClick = async () => {
  const route = useRoute();
  if (route.path === '/') {
    // 如果当前在主页，触发内容刷新事件
    const { memoReloadEvent } = await import('~/event');
    memoReloadEvent.emit('refresh');
  } else {
    // 否则导航到主页
    await navigateTo('/');
  }
};
</script>

<style scoped></style>
