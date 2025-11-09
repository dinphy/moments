<template>
  <USlideover
    v-model="showSearchDrawer"
    :ui="{
      width: 'sm:max-w-md md:max-w-lg lg:max-w-xl w-screen',
      overlay: {
        base: 'fixed inset-0 bg-gray-900/50 backdrop-blur-sm',
      },
      background: 'bg-white dark:bg-gray-900',
      ring: '',
      rounded: '',
      shadow: 'shadow-xl',
      padding: 'p-0',
      margin: '',
      height: 'h-screen',
    }"
  >
    <div class="h-[calc(100vh-72px)] overflow-y-auto p-4 sm:p-6">
      <div class="flex items-center mb-4">
        <UIcon
          @click="showSearchDrawer = false"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-4"
        />
        <h3 class="text-lg font-bold">搜索内容</h3>
      </div>
      <!-- 搜索框 -->
      <div class="mb-4">
        <UInput
          v-model="searchQuery"
          placeholder="输入关键词搜索..."
          size="xl"
          icon="i-heroicons-magnifying-glass"
          @keyup.enter="performSearch"
          class="w-full"
        />
      </div>
      <!-- 搜索结果 -->
      <div v-if="searchResults.length > 0" class="space-y-3">
        <p class="text-sm text-gray-500 mb-3">找到 {{ searchResults.length }} 条结果</p>
        <div class="space-y-3">
          <div v-for="result in searchResults" :key="result.id" 
              class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all duration-200 cursor-pointer"
              @click="navigateToMemo(result.id)">
            <div class="flex items-start space-x-3">
              <img v-if="result.user && result.user.avatarUrl" 
                  :src="result.user.avatarUrl" 
                  :alt="result.user.nickname || '用户头像'"
                  class="w-10 h-10 rounded-full flex-shrink-0">
              <div v-else 
                  class="w-10 h-10 rounded-full bg-gray-300 dark:bg-gray-600 flex items-center justify-center flex-shrink-0">
                <UIcon name="i-heroicons-user" class="w-5 h-5 text-gray-500" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-2 mb-2">
                  <span class="font-medium text-sm">{{ result.user ? result.user.nickname : '匿名用户' }}</span>
                  <span class="text-xs text-gray-500">{{ formatDate(result.createdAt) }}</span>
                </div>
                <p class="text-gray-700 dark:text-gray-300 text-sm leading-relaxed line-clamp-3">
                  {{ result.content }}
                </p>
                <div v-if="result.images && result.images.length > 0" class="grid grid-cols-3 gap-2 mt-3">
                  <img v-for="(img, index) in result.images.slice(0, 3)" 
                      :key="index"
                      :src="img" 
                      :alt="`图片 ${index + 1}`"
                      class="w-full h-20 object-cover rounded-md">
                </div>
                <div v-if="result.tags && result.tags.length > 0" class="flex flex-wrap gap-1 mt-2">
                  <span v-for="(tag, index) in result.tags.slice(0, 3)" 
                        :key="index"
                        class="inline-block px-2 py-1 text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 rounded-full">
                    #{{ typeof tag === 'string' ? tag : tag.name }}
                  </span>
                  <span v-if="result.tags.length > 3" 
                        class="inline-block px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-400 rounded-full">
                    +{{ result.tags.length - 3 }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- 加载更多 -->
        <div v-if="hasMore" class="text-center mt-6">
          <UButton @click="loadMore" :loading="loading" size="sm" variant="outline">加载更多</UButton>
        </div>
      </div>
      <!-- 无结果提示 -->
      <div v-else-if="hasSearched && !loading" class="text-center py-8">
        <UIcon name="i-heroicons-magnifying-glass" class="w-10 h-10 mx-auto text-gray-400 mb-3" />
        <p class="text-gray-500 text-sm">没有找到相关内容</p>
        <p class="text-xs text-gray-400 mt-1">尝试使用不同的关键词</p>
      </div>
      <!-- 初始状态提示 -->
      <div v-else-if="!hasSearched" class="text-center py-8">
        <UIcon name="i-heroicons-magnifying-glass" class="w-10 h-10 mx-auto text-gray-400 mb-3" />
        <p class="text-gray-500 text-sm">输入关键词开始搜索</p>
      </div>
    </div>
  </USlideover>
</template>

<script setup lang="ts">
import type {MemoVO, UserVO} from "~/types";

// 使用全局状态来控制搜索抽屉
const showSearchDrawer = useState<boolean>('showSearchDrawer', () => false)

// 搜索相关状态
const searchQuery = ref('')
const searchResults = ref<MemoVO[]>([])
const hasSearched = ref(false)
const loading = ref(false)
const hasMore = ref(false)
const page = ref(1)
const pageSize = ref(10)

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 执行搜索
const performSearch = async (resetPage = true) => {
  if (!searchQuery.value.trim()) return

  loading.value = true

  if (resetPage) {
    page.value = 1
    searchResults.value = []
  }

  try {
    const response = await useMyFetch<{
      list: MemoVO[],
      total: number,
      hasNext: boolean
    }>('/memo/list', {
      page: page.value,
      size: pageSize.value,
      contentContains: searchQuery.value,
      showType: 1, // 只搜索公开内容
    })

    if (resetPage) {
      searchResults.value = response.list
    } else {
      searchResults.value = [...searchResults.value, ...response.list]
    }

    hasMore.value = response.hasNext
    hasSearched.value = true
  } catch (error) {
    console.error('搜索失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载更多结果
const loadMore = () => {
  page.value++
  performSearch(false)
}

// 导航到详情页
const navigateToMemo = (id: number) => {
  showSearchDrawer.value = false
  navigateTo(`/memo/${id}`)
}
</script>
