<template>
  <Header v-if="memos.length > 0" v-bind:user="memos[0].user" />

  <!-- 按年份分组展示 -->
  <div v-for="(yearGroup, year) in groupedMemos" :key="year">
    <h2 class="text-lg font-semibold p-4">{{ year }}<span class="text-xs">年</span></h2>
    <div class="flex flex-col">
      <Memo v-bind:memo="m" v-for="m in yearGroup" :key="m.id" />
    </div>
  </div>

  <div
    ref="loadMoreEle"
    class="text-xs text-center text-gray-500 py-2"
    @click="loadMore"
    v-if="hasNext"
  >
    点击加载更多
  </div>
  <div class="text-xs text-center text-gray-500 py-2" @click="loadMore" v-else>
    已经到底啦
  </div>
</template>

<script setup lang="ts">
import type { MemoVO } from "~/types";
import Memo from "~/components/userMemo.vue";
import { memoChangedEvent, memoReloadEvent } from "~/event";
import { useElementVisibility } from "@vueuse/core";
import { ref, computed, reactive, onMounted, watch } from "vue";
import { useMyFetch } from "~/utils";
import { useRoute } from "vue-router";

const loadMoreEle = ref(null);
const targetIsVisible = useElementVisibility(loadMoreEle);
watch(targetIsVisible, async (visible) => {
  if (visible) {
    await loadMore();
  }
});

const hasNext = ref(false);
const route = useRoute();
const userId = route.params.id as any as string;
const state = reactive({
  page: 1,
  size: 10,
});

const memos = ref<Array<MemoVO>>([]);

const groupedMemos = computed(() => {
  const grouped: Record<string, MemoVO[]> = {};
  memos.value.forEach((memo) => {
    const year = new Date(memo.createdAt).getFullYear();
    if (!grouped[year]) {
      grouped[year] = [];
    }
    grouped[year].push(memo);
  });
  return grouped;
});

onMounted(async () => {
  await reload();
});

const reload = async () => {
  const res = await useMyFetch<{
    list: Array<MemoVO>;
    total: number;
    hasNext: boolean;
  }>("/memo/list", {
    ...state,
    userId: parseInt(userId),
  });
  memos.value = res.list;
  hasNext.value = res.hasNext;
};

const loadMore = async () => {
  state.page = state.page + 1;
  const res = await useMyFetch<{
    list: Array<MemoVO>;
    total: number;
    hasNext: boolean;
  }>("/memo/list", state);
  memos.value = [...memos.value, ...res.list];
  hasNext.value = res.hasNext;
};

memoReloadEvent.on(async () => {
  await reload();
});

memoChangedEvent.on(async (id: number) => {
  const res = await useMyFetch<MemoVO>("/memo/get?latest=1&id=" + id);
  const index = memos.value.findIndex((r) => r.id === id);
  if (index >= 0) {
    memos.value[index] = res;
  }
});
</script>

<style scoped></style>
