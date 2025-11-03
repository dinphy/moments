<template>
  <div class="md:hidden fixed bottom-0 left-0 right-0 bg-white/95 dark:bg-[#202020]/95 backdrop-blur-md z-10 border-t border-gray-100 dark:border-gray-700">
    <div class="flex justify-around items-center">
      <!-- 动态 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-green-500': $route.path === '/' }"
        @click="handleHomeClick"
      >
        <UIcon name="i-carbon-activity" class="w-6 h-6" />
        <span class="text-xs mt-0.5">动态</span>
      </div>

      <!-- 发现 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-green-500': $route.path === '/discover' }"
        @click="navigate('/discover')"
      >
        <UIcon name="i-system-uicons-compass" class="w-6 h-6" />
        <span class="text-xs mt-0.5">发现</span>
      </div>

      <!-- 我 -->
      <div
        class="flex flex-col items-center py-1.5"
        :class="{ 'text-green-500': $route.path === '/user/settings' }"
        @click="handleUserClick"
      >
        <UIcon name="i-system-uicons-user-male" class="w-6 h-6" />
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
