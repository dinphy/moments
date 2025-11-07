<template>
  <div v-if="showNav" class="md:hidden fixed bottom-0 left-0 right-0 bg-[#F9F9F9]/95 dark:bg-[#202020]/95 backdrop-blur-md z-10 transition-all duration-300">
    <div class="flex justify-around items-center">
      <div
        class="flex flex-col items-center py-1"
        :class="{ 'text-[#07C160]': $route.path === '/' }"
        @click="handleHomeClick"
      >
        <UIcon name="i-carbon-activity" class="w-6 h-6" />
        <span class="text-[10px] mt-0.5">动态</span>
      </div>

      <div
        class="flex flex-col items-center py-1"
        :class="{ 'text-[#07C160]': $route.path === '/discover' }"
        @click="navigate('/discover')"
      >
        <UIcon :name="$route.path === '/discover' ? 'i-weui-discover-filled' : 'i-weui-discover-outlined'" class="w-6 h-6" />
        <span class="text-[10px] mt-0.5">发现</span>
      </div>

      <div
        class="flex flex-col items-center py-1"
        :class="{ 'text-[#07C160]': $route.path === '/user/settings' }"
        @click="handleUserClick"
      >
        <UIcon :name="$route.path === '/user/settings' ? 'i-weui-me-filled' : 'i-weui-me-outlined'" class="w-6 h-6" />
        <span class="text-[10px] mt-0.5">我的</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";

const global = useGlobalState();
const loginReg = useState<boolean>("loginReg", () => false);
const route = useRoute();

// 添加计算属性判断是否显示导航栏
const showNav = computed(() => {
  return ['/', '/discover', '/user/settings'].includes(route.path);
});

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
    const { memoReloadEvent } = await import('~/event');
    memoReloadEvent.emit('refresh');
  } else {
    await navigateTo('/');
  }
};
</script>

<style scoped></style>
