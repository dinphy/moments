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
        <UInput
          v-model="searchQuery"
          placeholder="搜索"
          size="md"
          icon="i-heroicons-magnifying-glass"
          :ui="{ icon: { trailing: { pointer: '' } } }"
          @keyup.enter="performSearch"
          class="w-full"
        >
          <template #trailing>
            <UButton
              v-show="searchQuery !== ''"
              color="gray"
              variant="link"
              icon="i-heroicons-x-mark-20-solid"
              :padded="false"
              @click="clearSearch"
            />
          </template>
        </UInput>
      </div>
      <!-- 搜索结果 -->
      <div v-if="hasSearched && searchResults && searchResults.length > 0" class="space-y-3">
        <p class="text-sm text-gray-500 mb-3">找到 {{ totalResults }} 条结果</p>
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
                  <span class="font-medium text-sm">{{ result.user.nickname }}</span>
                  <span class="text-xs text-gray-500">{{ formatDate(result.createdAt) }}</span>
                </div>
                <div 
                  class="markdown-content text-gray-700 dark:text-gray-300 text-sm leading-relaxed line-clamp-3" 
                  v-html="renderMarkdown(result.content)"
                ></div>
                <div v-if="result.imgs" class="grid grid-cols-3 gap-2 mt-3">
                  <img v-for="(img, index) in getImages(result.imgs).slice(0, 3)" 
                      :key="index"
                      :src="img" 
                      :alt="`图片 ${index + 1}`"
                      class="w-full h-20 object-cover rounded-md">
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- 加载更多 -->
        <div ref="loadMoreEle" class="text-xs text-center text-gray-500 py-2 cursor-pointer mt-6" @click="loadMore" v-if="hasMore">
          点击加载更多
        </div>
        <div class="text-xs text-center text-gray-500 py-2 mt-6" v-else-if="hasSearched">
          已经到底啦
        </div>
      </div>
      <!-- 无结果提示 -->
      <div v-else-if="hasSearched && !loading" class="text-center py-12">
        <div class="relative inline-block mb-4">
          <div class="absolute inset-0 bg-gray-200 dark:bg-gray-700 rounded-full blur-2xl opacity-30"></div>
          <UIcon name="i-heroicons-magnifying-glass" class="relative w-12 h-12 mx-auto text-gray-400" />
        </div>
        <h3 class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2">未找到相关内容</h3>
        <p class="text-gray-500 dark:text-gray-400 text-sm max-w-xs mx-auto mb-2">试试这些关键词，可能会有结果</p>
        <div class="flex flex-wrap justify-center gap-2">
          <UBadge
            v-for="suggestion in ['学习', '生活', '旅行', '工作', '音乐']"
            :key="suggestion"
            size="sm"
            color="gray"
            variant="soft"
            class="cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-300"
            @click="searchQuery = suggestion; performSearch()"
          >
            {{ suggestion }}
          </UBadge>
        </div>
      </div>
      <!-- 初始状态提示 -->
      <div v-else-if="!hasSearched" class="text-center py-12">
        <div class="relative inline-block mb-4">
          <div class="absolute inset-0 bg-primary-200 dark:bg-primary-800 rounded-full blur-2xl opacity-30"></div>
          <UIcon name="i-heroicons-magnifying-glass" class="relative w-12 h-12 mx-auto text-primary-500 dark:text-primary-400" />
        </div>
        <h3 class="text-xl font-medium text-gray-700 dark:text-gray-300 mb-2">探索精彩内容</h3>
        <p class="text-gray-500 dark:text-gray-400 text-sm max-w-xs mx-auto">输入关键词，发现感兴趣的内容、照片和回忆</p>
        <div class="flex flex-wrap justify-center gap-2 mt-6" v-if="tags.length > 0">
          <UBadge
            v-for="tag in tags.slice(0, 8)"
            :key="tag"
            size="md"
            color="primary"
            variant="soft"
            class="cursor-pointer hover:bg-primary-200 dark:hover:bg-primary-600 transition-colors duration-300"
            @click="navigateToTag(tag)"
          >
            #{{ tag }}
          </UBadge>
        </div>
        <div v-if="!global.userinfo.token" class="flex flex-wrap justify-center gap-2 mt-2">
          <UBadge
            v-for="suggestion in ['学习', '生活', '旅行', '工作', '音乐']"
            :key="suggestion"
            size="sm"
            color="gray"
            variant="soft"
            class="cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-300"
            @click="searchQuery = suggestion; performSearch()"
          >
            {{ suggestion }}
          </UBadge>
        </div>
      </div>
    </div>
  </USlideover>
</template>

<script setup lang="ts">
import type {MemoVO, UserVO} from "~/types";
import {useElementVisibility} from '@vueuse/core'
import { useGlobalState } from "~/store";
import { md } from "~/utils";

const global = useGlobalState();
const showSearchDrawer = useState<boolean>('showSearchDrawer', () => false)

// 搜索相关状态
const searchQuery = ref('')
const searchResults = ref<MemoVO[]>([])
const hasSearched = ref(false)
const loading = ref(false)
const hasMore = ref(false)
const page = ref(1)
const pageSize = ref(10)
const totalResults = ref(0) // 存储总结果数
const tags = ref<string[]>([]) // 热门标签

// 加载更多元素引用
const loadMoreEle = ref(null)
const targetIsVisible = useElementVisibility(loadMoreEle)

// 格式化日期
const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

// 处理图片字符串，转换为数组
const getImages = (imgs: string) => {
  if (!imgs) return []
  try {
    return JSON.parse(imgs)
  } catch (e) {
    return imgs.split(',').filter(url => url.trim())
  }
}

// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  hasSearched.value = false
  searchResults.value = []
  page.value = 1
  hasMore.value = false
}

// 导航到标签页面并关闭抽屉
const navigateToTag = (tag: string) => {
  navigateTo(`/tags/${global.value.userinfo.username}/${tag}`)
  showSearchDrawer.value = false
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
      totalResults.value = response.total // 保存总结果数
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

// 监听加载更多元素是否可见
watch(targetIsVisible, async (visible) => {
  if (visible && hasMore.value && !loading.value) {
    await loadMore()
  }
})

// 获取标签列表
const fetchTags = async () => {
  if (global.value.userinfo.token) {
    try {
      const response = await useMyFetch<{tags: string[]}>('/tag/list');
      tags.value = response.tags || [];
    } catch (error) {
      console.error('获取标签列表失败:', error);
    }
  }
};

// 页面加载时获取标签
onMounted(() => {
  fetchTags();
})

// 渲染Markdown内容
const renderMarkdown = (content: string) => {
  if (!content) return "";
  try {
    return md.render(content);
  } catch (e) {
    console.error("Markdown渲染错误:", e);
    return content;
  }
};

// 导航到详情页
const navigateToMemo = (id: number) => {
  showSearchDrawer.value = false
  navigateTo(`/memo/${id}`)
}
</script>
