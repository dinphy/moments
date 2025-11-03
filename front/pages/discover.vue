<template>
  <Header :user="currentUser" />
  <div class="bg-gray-100 dark:bg-gray-900 min-h-screen rounded-b-lg p-2">
    <div class="grid grid-cols-2 gap-2 mb-2">
      <!-- 发动态卡片 -->
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
          <UIcon
                v-if="mode.preference === 'light'"
                name="i-carbon-sun"
                class="w-10 h-10"
              />

              <UIcon
                v-else-if="mode.preference === 'dark'"
                name="i-carbon-moon"
                class="w-10 h-10"
              />

              <UIcon
                v-else
                name="i-carbon-laptop"
                class="w-10 h-10"
              />
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
