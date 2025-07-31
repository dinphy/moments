<template>
  <div class="min-h-screen bg-white-to-br from-blue-50 to-purple-50">
    
    <div class="max-w-4xl mx-auto px-4 py-16 text-center">
      <!-- 动态插图容器 -->
      <div class="relative h-64 mb-8">
        <div class="absolute inset-0 flex items-center justify-center">
          <div class="animate-pulse">
            <svg class="w-48 h-48 sm:w-64 sm:h-64 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="text-4xl sm:text-6xl font-bold text-blue-600 absolute inset-0 flex items-center justify-center">
            404
          </div>
        </div>
      </div>

      <h1 class="text-3xl sm:text-4xl font-bold text-gray-800 mb-4">
        哎呀，页面迷路了~
      </h1>
      
      <div class="space-y-4">
        <p class="text-md sm:text-lg text-gray-600">
          该页面可能已被移动、删除或暂时不可用
        </p>
        <div class="flex justify-center space-x-4">
          <NuxtLink to="/" 
            class="px-4 py-2 sm:px-6 sm:py-3 bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-lg shadow-lg
                   hover:from-blue-600 hover:to-purple-600 transition-all duration-300 transform hover:scale-105">
            返回首页
          </NuxtLink>
          <NuxtLink to="/user/1" 
            class="px-4 py-2 sm:px-6 sm:py-3 border-2 border-blue-500 text-blue-500 rounded-lg
                   hover:bg-blue-50 transition-colors duration-300">
            推荐页面
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
useHead({
  title: '页面未找到 - 404'
})
</script>

<style scoped>
.bg-gradient-to-br {
  background-image: linear-gradient(to bottom right, var(--tw-gradient-from), var(--tw-gradient-to));
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
</style>