<template>
  <Header :user="currentUser" />
  <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 min-h-screen rounded-b-lg p-4 overflow-hidden relative">
    <div class="relative grid grid-cols-2 gap-4 mb-6"> 
      <div 
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-4 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
        @click="navigateTo('/friend')"
      >
        <div class="absolute top-0 right-0 w-32 h-32 bg-slate-300 dark:bg-slate-600 opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
        <div class="relative flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-4 items-center">
          <div class="w-14 h-14 bg-green-500 bg-opacity-10 dark:bg-green-400 dark:bg-opacity-10 rounded-xl flex items-center justify-center group-hover:bg-opacity-20 transition-all duration-300">
            <UIcon name="i-weui-contacts-outlined" class="w-8 h-8 text-green-600 dark:text-green-400"/>
          </div>
          <div class="text-center sm:text-left">
            <h3 class="font-bold text-md text-slate-800 dark:text-slate-200">友情链接</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">发现更多精彩博客</p>
          </div>
        </div>
      </div>

      <div 
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-4 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
        @click="toggleMode"
      >
        <div class="absolute top-0 right-0 w-32 h-32 bg-slate-300 dark:bg-slate-600 opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
        <div class="relative flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-4 items-center">
          <div class="w-14 h-14 bg-indigo-500 bg-opacity-10 dark:bg-indigo-400 dark:bg-opacity-10 rounded-xl flex items-center justify-center group-hover:bg-opacity-20 transition-all duration-300">
            <UIcon
              v-if="mode.preference === 'light'"
              name="i-carbon-sun"
              class="w-8 h-8 text-indigo-600 dark:text-indigo-400"
            />
            <UIcon
              v-else-if="mode.preference === 'dark'"
              name="i-carbon-moon"
              class="w-8 h-8 text-indigo-600 dark:text-indigo-400"
            />
            <UIcon
              v-else
              name="i-weui-display-outlined"
              class="w-8 h-8 text-indigo-600 dark:text-indigo-400"
            />
          </div>
          <div class="text-center sm:text-left">
            <h3 class="font-bold text-md text-slate-800 dark:text-slate-200">{{ modeText }}</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">主题切换</p>
          </div>
        </div>
      </div>
    </div>

    <div v-if="global.userinfo.token && global.userinfo.id === 1" class="relative bg-white dark:bg-gray-800 bg-opacity-90 backdrop-blur-md shadow-lg rounded-2xl p-4 overflow-hidden border border-slate-200 dark:border-slate-700 mb-6">
      <div class="absolute top-0 right-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -mr-20 -mt-20"></div>
      <div class="absolute bottom-0 left-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -ml-20 -mb-20"></div>
      <div class="relative">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <NuxtLink to="/user/manage" class="block bg-slate-100 dark:bg-gray-700/50 rounded-lg p-3 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center">
                  <UIcon name="i-weui-group-detail-outlined" class="w-5 h-5 text-amber-600 dark:text-amber-400"/>
                </div>
                <div>
                  <h3 class="font-medium text-gray-800 dark:text-gray-200">用户管理</h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">管理系统用户</p>
                </div>
              </div>
              <UIcon name="i-weui-arrow-outlined" class="w-4 h-4 text-gray-400" />
            </div>
          </NuxtLink>
          <NuxtLink to="/sys/settings" class="block bg-slate-100 dark:bg-gray-700/50 rounded-lg p-3 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-3">
                <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                  <UIcon name="i-weui-setting-outlined" class="w-5 h-5 text-blue-600 dark:text-blue-400"/>
                </div>
                <div>
                  <h3 class="font-medium text-gray-800 dark:text-gray-200">系统设置</h3>
                  <p class="text-sm text-gray-600 dark:text-gray-400">配置系统参数</p>
                </div>
              </div>
              <UIcon name="i-weui-arrow-outlined" class="w-4 h-4 text-gray-400" />
            </div>
          </NuxtLink>
        </div>
      </div>
    </div>

    <div v-if="!global.userinfo.token" class="relative bg-white dark:bg-gray-800 bg-opacity-90 backdrop-blur-md shadow-lg rounded-2xl p-8 overflow-hidden border border-slate-200 dark:border-slate-700">
      <div class="absolute top-0 right-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -mr-20 -mt-20"></div>
      <div class="absolute bottom-0 left-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -ml-20 -mb-20"></div>
      <div class="relative text-center">
        <div class="hidden sm:inline-flex items-center justify-center w-16 h-16 bg-slate-100 dark:bg-slate-700 rounded-full mb-4 border border-slate-200 dark:border-slate-600">
          <UIcon name="i-carbon-locked" class="w-8 h-8 text-slate-600 dark:text-slate-400" />
        </div>
        <h3 class="text-xl font-bold text-slate-800 dark:text-slate-200 mb-2">解锁更多功能</h3>
        <p class="text-slate-600 dark:text-slate-400 mb-6 max-w-md mx-auto hidden sm:block">登录后，发现更多精彩内容，记录你的生活点滴</p>
        <button 
          @click="loginReg = true" 
          class="inline-flex items-center justify-center px-6 py-3 bg-slate-600 dark:bg-slate-500 text-white rounded-xl font-medium hover:bg-slate-700 dark:hover:bg-slate-400 transition-all duration-300 transform hover:scale-105 shadow-md"
        >
          <UIcon name="i-carbon-login" class="w-5 h-5 mr-2" />
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
