<template>
  <Header :user="currentUser" />
  <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 min-h-screen rounded-b-lg p-4 overflow-hidden relative">
    
    <!-- 功能卡片区域 -->
    <div class="relative grid grid-cols-2 gap-4 mb-6">
      <!-- 发动态卡片 -->
      <div 
        v-if="global.userinfo.token"
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-6 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
        @click="navigateTo('/new')"
      >
        <div class="absolute top-0 right-0 w-32 h-32 bg-slate-300 dark:bg-slate-600 opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
        <div class="relative flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-4 items-center">
          <div class="w-14 h-14 bg-blue-500 bg-opacity-10 dark:bg-blue-400 dark:bg-opacity-10 rounded-xl flex items-center justify-center group-hover:bg-opacity-20 transition-all duration-300">
            <UIcon name="i-carbon-camera" class="w-8 h-8 text-blue-600 dark:text-blue-400"/>
          </div>
          <div class="text-center sm:text-left">
            <h3 class="font-bold text-lg text-slate-800 dark:text-slate-200">发表动态</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">分享你的精彩瞬间</p>
          </div>
        </div>
      </div>
      
      <!-- 日历检索卡片 -->
      <div 
        v-if="global.userinfo.token"
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-6 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
        @click="navigateTo('/user/calendar')"
      >
        <div class="absolute top-0 right-0 w-32 h-32 bg-slate-300 dark:bg-slate-600 opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
        <div class="relative flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-4 items-center">
          <div class="w-14 h-14 bg-purple-500 bg-opacity-10 dark:bg-purple-400 dark:bg-opacity-10 rounded-xl flex items-center justify-center group-hover:bg-opacity-20 transition-all duration-300">
            <UIcon name="i-jam-search-folder" class="w-8 h-8 text-purple-600 dark:text-purple-400"/>
          </div>
          <div class="text-center sm:text-left">
            <h3 class="font-bold text-lg text-slate-800 dark:text-slate-200">日历检索</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">那往昔的美好时光</p>
          </div>
        </div>
      </div>
      
      <!-- 友情链接卡片 -->
      <div 
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-6 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
        @click="navigateTo('/friend')"
      >
        <div class="absolute top-0 right-0 w-32 h-32 bg-slate-300 dark:bg-slate-600 opacity-10 rounded-full -mr-16 -mt-16 group-hover:scale-150 transition-transform duration-500"></div>
        <div class="relative flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-4 items-center">
          <div class="w-14 h-14 bg-green-500 bg-opacity-10 dark:bg-green-400 dark:bg-opacity-10 rounded-xl flex items-center justify-center group-hover:bg-opacity-20 transition-all duration-300">
            <UIcon name="i-weui-contacts-outlined" class="w-8 h-8 text-green-600 dark:text-green-400"/>
          </div>
          <div class="text-center sm:text-left">
            <h3 class="font-bold text-lg text-slate-800 dark:text-slate-200">友情链接</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">发现更多精彩博客</p>
          </div>
        </div>
      </div>
      
      <!-- 主题切换卡片 -->
      <div 
        class="group relative overflow-hidden bg-gradient-to-br from-slate-100 to-slate-200 dark:from-slate-700 dark:to-slate-800 rounded-2xl p-6 cursor-pointer transform transition-all duration-300 hover:scale-105 hover:shadow-xl border border-slate-200 dark:border-slate-600"
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
            <h3 class="font-bold text-lg text-slate-800 dark:text-slate-200">{{ modeText }}</h3>
            <p class="text-sm text-slate-600 dark:text-slate-400 opacity-90 hidden sm:block">主题切换</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 管理员功能区域 -->
    <div v-if="global.userinfo.token && global.userinfo.id === 1" class="relative bg-white dark:bg-gray-800 bg-opacity-90 backdrop-blur-md shadow-lg rounded-2xl px-4 py-8 overflow-hidden border border-slate-200 dark:border-slate-700 mb-6">
      <div class="absolute top-0 right-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -mr-20 -mt-20"></div>
      <div class="absolute bottom-0 left-0 w-40 h-40 bg-slate-200 dark:bg-slate-700 opacity-20 rounded-full -ml-20 -mb-20"></div>
      <div class="relative">
        <h3 class="text-lg font-semibold text-slate-800 dark:text-slate-200 mb-4 flex items-center">
          <UIcon name="i-heroicons-shield-check" class="w-5 h-5 mr-2 text-amber-500" />
          系统与安全
        </h3>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <NuxtLink to="/user/manage" class="group bg-gradient-to-br from-amber-50 to-orange-50 dark:from-amber-900/20 dark:to-orange-900/20 rounded-xl p-4 flex items-center justify-between hover:shadow-md transition-all duration-300 border border-amber-200 dark:border-amber-800/30">
            <div class="flex items-center">
              <div class="w-10 h-10 bg-amber-100 dark:bg-amber-800/30 rounded-lg flex items-center justify-center mr-3 group-hover:scale-110 transition-transform duration-300">
                <UIcon name="i-weui-group-detail-outlined" class="w-5 h-5 text-amber-600 dark:text-amber-400"/>
              </div>
              <span class="font-medium text-gray-700 dark:text-gray-300">用户管理</span>
            </div>
            <UIcon name="i-weui-arrow-outlined" class="w-4 h-4 text-amber-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300"/>
          </NuxtLink>
          <NuxtLink to="/sys/settings" class="group bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-xl p-4 flex items-center justify-between hover:shadow-md transition-all duration-300 border border-blue-200 dark:border-blue-800/30">
            <div class="flex items-center">
              <div class="w-10 h-10 bg-blue-100 dark:bg-blue-800/30 rounded-lg flex items-center justify-center mr-3 group-hover:scale-110 transition-transform duration-300">
                <UIcon name="i-weui-setting-outlined" class="w-5 h-5 text-blue-600 dark:text-blue-400"/>
              </div>
              <span class="font-medium text-gray-700 dark:text-gray-300">系统设置</span>
            </div>
            <UIcon name="i-weui-arrow-outlined" class="w-4 h-4 text-blue-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300"/>
          </NuxtLink>
        </div>
      </div>
    </div>
    
    <!-- 未登录提示区域 -->
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
