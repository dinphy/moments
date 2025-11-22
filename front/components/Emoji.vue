<template>
  <div class="max-h-[230px] overflow-y-auto">
    <!-- 最近使用 -->
    <div v-if="recentEmojis.length > 0" class="border-b border-gray-200 dark:border-gray-700">
      <div class="px-3 py-2 text-sm text-gray-600 dark:text-gray-400 font-medium">
        最近使用
      </div>
      <div class="grid grid-cols-[repeat(auto-fill,minmax(30px,1fr))] gap-2 p-2">
        <div 
          v-for="emoji in recentEmojis" 
          :key="emoji.code"
          class="flex flex-col items-center rounded cursor-pointer transition-colors duration-200 hover:bg-gray-100 dark:hover:bg-gray-700 flex-shrink-0"
          @click="selectEmoji(emoji.code)"
          :title="emoji.name"
        >
          <img 
            :src="emoji.path" 
            :alt="emoji.name"
            class="w-7 h-7 object-contain"
          />
        </div>
      </div>
    </div>
    
    <!-- 所有表情 -->
    <div>
      <div class="px-3 py-2 text-sm text-gray-600 dark:text-gray-400 font-medium">
        所有表情
      </div>
      <div class="grid grid-cols-[repeat(auto-fill,minmax(30px,1fr))] gap-2 p-2">
        <div 
          v-for="emoji in allEmojis" 
          :key="emoji.code"
          class="flex flex-col items-center rounded cursor-pointer transition-colors duration-200 hover:bg-gray-100 dark:hover:bg-gray-700"
          @click="selectEmoji(emoji.code)"
          :title="emoji.name"
        >
          <img 
            :src="emoji.path" 
            :alt="emoji.name"
            class="w-7 h-7 object-contain"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getAllEmojis } from '~/utils/emoji'

const emit = defineEmits(['selected'])
const RECENT_EMOJIS_KEY = 'recent_emojis'
const MAX_RECENT_EMOJIS = 8
const allEmojis = getAllEmojis()
const recentEmojis = ref<any[]>([])

const selectEmoji = (emojiCode: string) => {
  addToRecent(emojiCode)
  emit('selected', emojiCode)
}

const addToRecent = (emojiCode: string) => {
  try {
    const recentCodes = getRecentCodes()
    const filteredCodes = recentCodes.filter(code => code !== emojiCode)
    const newRecentCodes = [emojiCode, ...filteredCodes].slice(0, MAX_RECENT_EMOJIS)
    localStorage.setItem(RECENT_EMOJIS_KEY, JSON.stringify(newRecentCodes))

    updateRecentEmojis(newRecentCodes)
  } catch (error) {
    console.error('保存最近使用的表情失败:', error)
  }
}

const getRecentCodes = (): string[] => {
  try {
    const stored = localStorage.getItem(RECENT_EMOJIS_KEY)
    return stored ? JSON.parse(stored) : []
  } catch {
    return []
  }
}

const updateRecentEmojis = (codes: string[]) => {
  recentEmojis.value = codes.map(code => {
    return allEmojis.find(emoji => emoji.code === code)
  }).filter(Boolean)
}

onMounted(() => {
  const recentCodes = getRecentCodes()
  updateRecentEmojis(recentCodes)
})
</script>

<style scoped>
</style>