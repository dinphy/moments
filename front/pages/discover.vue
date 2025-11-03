<template>
  <Header :user="currentUser" />
  <div class="bg-gray-100 dark:bg-gray-900 min-h-screen rounded-b-lg p-2">
    <div class="grid grid-cols-2 gap-2 mb-2">
      <!-- 发表动态卡片 -->
      <div 
        v-if="global.userinfo.token"
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 flex flex-col items-center justify-center hover:shadow-lg transition-all duration-200 cursor-pointer group"
        @click="navigateTo('/new')"
      >
        <div class="mb-2">
          <UIcon name="i-carbon-camera" class="w-10 h-10 text-blue-500 transition-transform duration-200 group-hover:scale-110"/>
        </div>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">发表动态</span>
      </div>
      <!-- 日历检索卡片 -->
      <div 
        v-if="global.userinfo.token"
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 flex flex-col items-center justify-center hover:shadow-lg transition-all duration-200 cursor-pointer group"
        @click="navigateTo('/user/calendar')"
      >
        <div class="mb-2">
          <UIcon name="i-jam-search-folder" class="w-10 h-10 text-purple-500 transition-transform duration-200 group-hover:scale-110"/>
        </div>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">日历检索</span>
      </div>
      <!-- 友情链接卡片 -->
      <div 
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 flex flex-col items-center justify-center hover:shadow-lg transition-all duration-200 cursor-pointer group"
        @click="navigateTo('/friend')"
      >
        <div class="mb-2">
          <UIcon name="i-carbon-friendship" class="w-10 h-10 text-green-500 transition-transform duration-200 group-hover:scale-110"/>
        </div>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">友情链接</span>
      </div>
      <!-- 主题切换卡片 -->
      <div 
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md p-4 flex flex-col items-center justify-center hover:shadow-lg transition-all duration-200 cursor-pointer group"
        @click="toggleMode"
      >
        <div class="relative mb-2">
          <svg
            v-if="mode.value === 'light'"
            class="w-10 h-10 text-yellow-500 transition-transform duration-200 group-hover:scale-110"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9"></path>
            <path d="M20 3v4"></path>
            <path d="M22 5h-4"></path>
          </svg>
          <svg
            v-else
            class="w-10 h-10 text-yellow-500 transition-transform duration-200 group-hover:scale-110"
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <circle cx="12" cy="12" r="4"></circle>
            <path d="M12 2v2"></path>
            <path d="M12 20v2"></path>
            <path d="m4.93 4.93 1.41 1.41"></path>
            <path d="m17.66 17.66 1.41 1.41"></path>
            <path d="M2 12h2"></path>
            <path d="M20 12h2"></path>
            <path d="m6.34 17.66-1.41 1.41"></path>
            <path d="m19.07 4.93-1.41 1.41"></path>
          </svg>
        </div>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ modeText }}</span>
      </div>
    </div>
    <!-- 探索更多 -->
    <div v-if="!global.userinfo.token" class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="px-4 py-6 text-center">
        <UIcon name="i-carbon-locked" class="w-12 h-12 text-gray-400 mx-auto mb-3" />
        <p class="text-gray-600 dark:text-gray-400 text-sm">登录后，发现更多精彩内容</p>
        <button 
          @click="loginReg = true" 
          class="mt-3 px-4 py-2 bg-green-500 text-white rounded-md text-sm font-medium hover:bg-green-600 transition-colors"
        >
          立即登录
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toast } from "vue-sonner";
import type { UserVO } from "~/types";
import { useGlobalState } from "~/store";

const global = useGlobalState();
const currentUser = useState<UserVO>('userinfo');
const mode = useColorMode();
const loginReg = useState<boolean>("loginReg", () => false);

const modeText = computed(() => {
  if (mode.preference === 'light') {
    return '夜间模式';
  } else if (mode.preference === 'dark') {
    return '日间模式';
  } else {
    return '跟随系统';
  }
});

const toggleMode = () => {
  if (mode.preference === "system") {
    mode.preference = "dark";
  } else if (mode.preference === "dark") {
    mode.preference = "light";
  } else {
    mode.preference = "system";
    toast.success("显示模式将跟随系统设置");
  }
};
</script>

<style scoped>
</style>
