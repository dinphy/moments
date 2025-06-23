<template>
  <div>
    <Header :user="currentUser"/>

    <div class="p-4 space-y-4">
      <div class="flex items-center justify-between gap-2">
          <UButtonGroup class="flex-1 flex">
          <UInput
              v-model="state.contentContains"
              placeholder="请输入关键词..."
              size="xl"
              name="contentContains"
              class="flex-1"
              :ui="{ base: 'rounded-r-none pr-3 z-[1]' }"
          />
          <USelectMenu 
            v-model="state.showType"
            :options="[
              { label: '所有', value: -1 }, 
              { label: '公开', value: 1 }, 
              { label: '私密', value: 0 }
            ]"
            option-attribute="label"
            value-attribute="value"
            size="xl"
            :ui="{ base: 'rounded-l-none' }"
          >
            <template #label>
              <span v-if="state.showType === -1">所有</span>
              <span v-else-if="state.showType === 1">公开</span>
              <span v-else>私密</span>
            </template>
          </USelectMenu>
          </UButtonGroup>
          <UButton @click="reload" class="ml-2 px-4">搜索</UButton>
      </div>
      <div
          class="flex items-center text-sm text-gray-500 gap-1 border-b pb-4"
          :class="[(state.isSearching || state.total !== undefined) ? 'justify-between' : 'justify-end']"
      >
          <span v-if="state.isSearching">
            正在检索<UBadge class="text-neutral mx-1" variant="outline">{{ state.contentContains }}</UBadge>中..
          </span>
          <span v-else-if="state.total !== undefined">
            <UBadge variant="solid">共 {{ state.total }} 条内容</UBadge>
          </span>
          <span v-else class="flex-1">暂无相关内容</span>
          <span class="flex items-center">
          高级：<UToggle v-model="openSwitch" />
          </span>
      </div>
      <div v-if="openSwitch" class="flex flex-wrap gap-4">
          <UFormGroup
          label="按时间"
          name="contentContains"
          class="flex-1"
          :ui="{ label: { base: 'font-bold' } }"
          >
          <UPopover :popper="{ placement: 'bottom-start' }">
              <UButton
              icon="i-heroicons-calendar-days-20-solid"
              color="white"
              variant="solid"
              class="w-full"
              >
              {{ isRangeSelected(null) ? '日期不限' : ranges.find(r => isRangeSelected(r.duration))?.label }}
              </UButton>

              <template #panel="{ close }">
                <div class="flex flex-col py-4 space-y-2">
                    <UButton
                        v-for="(range, index) in ranges"
                        :key="index"
                        :label="range.label"
                        color="gray"
                        variant="ghost"
                        class="px-6"
                        :class="[
                        isRangeSelected(range.duration)
                            ? 'bg-gray-100 dark:bg-gray-800'
                            : 'hover:bg-gray-50 dark:hover:bg-gray-800/50',
                        ]"
                        truncate
                        @click="selectRange(range.duration); close()"
                    />
                </div>
              </template>
          </UPopover>
          </UFormGroup>

          <UFormGroup
          label="按标签"
          name="tagContains"
          class="flex-1"
          :ui="{ label: { base: 'font-bold' } }"
          >
          <USelectMenu multiple v-model="state.tags" searchable :options="tags">
              <template #label>
              <span v-if="state.tags.length" class="truncate">{{
                  state.tags.join(", ")
              }}</span>
              <span v-else>选择标签</span>
              </template>
          </USelectMenu>
          </UFormGroup>
      </div>
    </div>

    <div class="flex flex-col divide-y divide-[#C0BEBF]/20 ">
      <Memo v-bind:memo="m" v-for="m in memos" :key="m.id"/>
    </div>
    <div ref="loadMoreEle" class="text-xs text-center text-gray-500 py-2 cursor-pointer" @click="loadMore" v-if="hasNext">
      点击加载更多
    </div>
    <div class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-else>
      已经到底啦
    </div>
  </div>
</template>

<script setup lang="ts">
import type {MemoVO, UserVO, SysConfigVO} from "~/types";

import {isSameDay, sub} from "date-fns";
import Memo from "~/components/Memo.vue";
import {memoChangedEvent, memoReloadEvent} from "~/event";
import {useElementVisibility} from '@vueuse/core'

const ranges = [
  {label: '日期不限', duration: null as null | Duration},
  {label: '最近一周', duration: {days: 7}},
  {label: '一个月内', duration: {days: 31}},
  {label: '三个月内', duration: {days: 90}},
  {label: '一年以内', duration: {months: 12}},
  {label: '最近三年', duration: {years: 3}},
]
const tags = ref<string[]>([])
const currentUser = useState<UserVO>('userinfo')
const state = reactive({
  page: 1,
  size: 10,
  contentContains: "",
  tags: [],
  showType: -1,
  total: 0,
  isSearching: false,
  range: {
    start: new Date(1970, 0, 1),
    end: new Date(2100, 0, 1)
  }
})
const openSwitch = ref(false)

function isRangeSelected(duration: Duration | null) {
  if (!duration) {
    return isSameDay(state.range.start, new Date(1970, 0, 1)) && 
           isSameDay(state.range.end, new Date(2100, 0, 1))
  }
  return isSameDay(state.range.start, sub(new Date(), duration)) && 
         isSameDay(state.range.end, new Date())
}

function selectRange(duration: Duration | null) {
  if (!duration) {
    state.range = {start: new Date(1970, 0, 1), end: new Date(2100, 0, 1)}
  } else {
    state.range = {start: sub(new Date(), duration), end: new Date()}
  }
}

const loadTags = async () => {
  const res = await useMyFetch<{
    tags: string[]
  }>("/tag/list")
  tags.value = res.tags
}

const loadMoreEle = ref(null)
const targetIsVisible = useElementVisibility(loadMoreEle)
const sysConfig = useState<SysConfigVO>("sysConfig")

watch(targetIsVisible, async (visible) => {
  if (visible && sysConfig.value?.enableAutoLoadNextPage) {
    await loadMore()
  }
});
const hasNext = ref(false)

const memos = ref<Array<MemoVO>>([])
onMounted(async () => {
  await loadTags()
  await reload()
})

const reload = async () => {
  state.page = 1
  state.isSearching = true
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', {
    size: state.size,
    start: state.range.start,
    end: state.range.end,
    contentContains: state.contentContains,
    showType: state.showType,
    tag: state.tags.join(','),
  })
  memos.value = res.list
  hasNext.value = res.hasNext
  state.total = res.total
  state.isSearching = false
}

const loadMore = async () => {
  state.page = state.page + 1
  const res = await useMyFetch<{
    list: Array<MemoVO>,
    total: number,
    hasNext: boolean
  }>('/memo/list', {
    page: state.page,
    size: state.size,
    start: state.range.start,
    end: state.range.end,
    contentContains: state.contentContains,
    showType: state.showType,
    tag: state.tags.join(','),
  })
  memos.value = [...memos.value, ...res.list]
  hasNext.value = res.hasNext
  state.total = res.total
}

memoReloadEvent.on(async () => {
  await reload()
})

memoChangedEvent.on(async (id: number) => {
  const res = await useMyFetch<MemoVO>('/memo/get?latest=1&id=' + id)
  const index = memos.value.findIndex(r => r.id === id)
  if (index >= 0) {
    memos.value[index] = res
  }
})

</script>

<style scoped>

</style>