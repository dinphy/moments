<template>
  <div class="max-h-[200px] overflow-y-auto">
    <UTabs class="mt-2" :items="tabItems"
           :ui="{wrapper: 'space-y-0', list: {height: 'h-8', tab: {height: 'h-6', padding: 'px-1'}}}">
      <template #item="{ item: tabItem }">
        <div class="grid grid-cols-[repeat(auto-fill,minmax(50px,1fr))] gap-1 p-2">
          <div 
            v-for="emoji in tabItem.emojis" 
            :key="emoji.code"
            class="flex flex-col items-center p-1 rounded cursor-pointer transition-colors duration-200 hover:bg-gray-100 dark:hover:bg-gray-700"
            @click="selectEmoji(emoji.code)"
            :title="emoji.name"
          >
            <img 
              :src="emoji.path" 
              :alt="emoji.name"
              class="w-6 h-6 object-contain"
            />
            <span class="text-[10px] text-gray-600 dark:text-gray-400 mt-0.5 text-center leading-tight">{{ emoji.name }}</span>
          </div>
        </div>
      </template>
    </UTabs>
  </div>
</template>

<script setup lang="ts">
import { getAllEmojis } from '~/utils/emoji'

const emit = defineEmits(['selected'])

const selectEmoji = (emojiCode: string) => {
  emit('selected', emojiCode)
}

// 获取所有本地表情包
const allEmojis = getAllEmojis()

// 将表情包按类别分组（只保留常用和所有表情两组）
const tabItems = [
  {
    key: 'common',
    label: '常用',
    emojis: [
      // 前30个基础表情
      ...allEmojis.slice(0, 30),
      // 添加常用的手势表情
      ...allEmojis.filter(emoji => [
        '强', '弱', '握手', '胜利', 'OK', '抱拳', '勾引', '拳头', '合十'
      ].includes(emoji.name))
    ]
  },
  {
    key: 'all',
    label: '所有表情',
    emojis: allEmojis
  }
]
</script>

<style scoped>
/* 移除自定义CSS，使用Tailwind类替代 */
</style>