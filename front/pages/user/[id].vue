<template>
  <Header v-if="memos.length > 0" v-bind:user="memos[0].user" />

  <div class="flex border-b border-gray-200 dark:border-gray-700 mb-4 bg-gray-50 dark:bg-gray-800 rounded-t-lg">
    <button
      @click="activeTab = 'memos'"
      class="flex-1 px-4 py-3 text-base font-medium transition-all duration-200 relative"
      :class="{
        'text-blue-600 dark:text-blue-400 bg-white dark:bg-gray-700 shadow-sm': activeTab === 'memos',
        'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50': activeTab !== 'memos'
      }"
    >
      <span class="flex items-center justify-center gap-2">
        <UIcon name="i-carbon-activity" class="text-lg" />
        动态
      </span>
      <div
        v-if="activeTab === 'memos'"
        class="absolute bottom-0 left-0 w-full h-0.5 bg-blue-600 dark:bg-blue-400 rounded-t-lg"
      ></div>
    </button>
    <button
      @click="activeTab = 'album'"
      class="flex-1 px-4 py-3 text-base font-medium transition-all duration-200 relative"
      :class="{
        'text-blue-600 dark:text-blue-400 bg-white dark:bg-gray-700 shadow-sm': activeTab === 'album',
        'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700/50': activeTab !== 'album'
      }"
    >
      <span class="flex items-center justify-center gap-2">
        <UIcon name="i-carbon-image" class="text-lg" />
        相册
      </span>
      <div
        v-if="activeTab === 'album'"
        class="absolute bottom-0 left-0 w-full h-0.5 bg-blue-600 dark:bg-blue-400 rounded-t-lg"
      ></div>
    </button>
  </div>

  <!-- 动态内容 -->
  <div v-if="activeTab === 'memos'" class="flex flex-col">
    <div v-if="sysConfig.enableNewMemo">
      <div v-for="(memo, index) in pinnedMemos" :key="index">
        <MemoItem v-bind:memo="memo" />
      </div>
      <div v-for="(memo, index) in nonPinnedMemoList" :key="index">
        <div v-if="memo.displayYear" class="pl-4 py-4">
          <span class="text-xl">{{ memo.displayYear }}</span>
          <span class="text-sm">年</span>
        </div>
        <MemoItem v-bind:memo="memo" />
      </div>
    </div>
    <div v-else class="flex flex-col divide-y divide-[#C0BEBF]/20">
      <Memo v-bind:memo="memo" v-for="memo in memos" :key="memo.id" />
    </div>
    <div
      v-if="hasNext"
      ref="loadMoreEle"
      class="text-xs text-center text-gray-500 py-2 cursor-pointer"
      @click="loadMore"
    >
      点击加载更多
    </div>
    <div class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-else>
      已经到底啦
    </div>
  </div>

  <!-- 相册内容 -->
  <div v-if="activeTab === 'album'" class="flex flex-col">
    <MemoAlbum :memos="memos" :loading="loading" />
    <div
      v-if="hasNext"
      ref="loadMoreEle"
      class="text-xs text-center text-gray-500 py-2 cursor-pointer"
      @click="loadMore"
    >
      点击加载更多
    </div>
    <div class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-else>
      已经到底啦
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MemoVO, SysConfigVO } from "~/types";
import Memo from "~/components/Memo.vue";
import MemoItem from "~/components/MemoItem.vue";
import MemoAlbum from "~/components/MemoAlbum.vue";
import { memoChangedEvent, memoReloadEvent } from "~/event";
import { useElementVisibility } from "@vueuse/core";
import dayjs from "dayjs";

const route = useRoute();
const userId = route.params.id as any as string;
const sysConfig = useState<SysConfigVO>("sysConfig");

const loading = ref(true);
const loadingMore = ref(false);
const hasNext = ref(false);
const loadMoreEle = ref(null);
const targetIsVisible = useElementVisibility(loadMoreEle);

const state = reactive({
  page: 1,
  size: 10,
});

const memos = ref<Array<MemoVO>>([]);

watch(targetIsVisible, async (visible) => {
  if (visible && sysConfig.value?.enableAutoLoadNextPage && !loadingMore.value) {
    await loadMore();
  }
});

const reload = async () => {
  state.page = 1;
  try {
    loading.value = true;
    const res = await useMyFetch<{
      list: Array<MemoVO>;
      total: number;
      hasNext: boolean;
    }>("/memo/list", {
      ...state,
      userId: parseInt(userId),
    });

    if (!res || typeof res.total === 'undefined') {
      navigateTo('/404');
      return;
    }
    
    memos.value = res.list;
    hasNext.value = res.hasNext;
  } catch (error) {
    navigateTo('/404');
  } finally {
    loading.value = false;
  }
};

const loadMore = async () => {
  if (loadingMore.value || !hasNext.value) return;

  loadingMore.value = true;
  try {
    state.page = state.page + 1;
    const res = await useMyFetch<{
      list: Array<MemoVO>;
      total: number;
      hasNext: boolean;
    }>("/memo/list", {
      ...state,
      userId: parseInt(userId),
    });
    memos.value = [...memos.value, ...res.list];
    hasNext.value = res.hasNext;
  } catch (error) {
    console.error('加载更多数据失败:', error);
  } finally {
    loadingMore.value = false;
  }
};
const activeTab = ref('memos');

onMounted(async () => {
  await reload();
});



memoReloadEvent.on(async () => {
  await reload();
});

memoChangedEvent.on(async (id: number) => {
  const res = await useMyFetch<MemoVO>("/memo/get?latest=1&id=" + id);
  if (memos.value && memos.value.length > 0) {
    const index = memos.value.findIndex((r) => r.id === id);
    if (index >= 0) {
      const updatedMemos = [...memos.value];
      updatedMemos[index] = res;
      memos.value = updatedMemos;
    }
  }
});

const pinnedMemos = computed(() => memos.value.filter((memo) => memo.pinned));
const nonPinnedMemos = computed(() =>
  memos.value.filter((memo) => !memo.pinned)
);

const nonPinnedMemoList = computed(() => {
  if (!nonPinnedMemos.value.length) return [];
  let lastYear: string | null = null;
  let lastDate: string | null = null;
  return nonPinnedMemos.value.map((memo) => {
    const currentYear = dayjs(memo.createdAt).locale("zh-cn").format("YYYY");
    const currentDate = dayjs(memo.createdAt).locale("zh-cn").format("YYYY-MM-DD");
    const displayYear = currentYear !== lastYear ? currentYear : null;
    const displayDate = currentDate !== lastDate ? currentDate : null;

    if (currentYear !== lastYear) lastYear = currentYear;
    if (currentDate !== lastDate) lastDate = currentDate;

    return { ...memo, displayYear, displayDate };
  });
});
</script>

<style scoped></style>