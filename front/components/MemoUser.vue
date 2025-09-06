<template>
  <div
    class="flex flex-row text-sm w-full hover:bg-slate-200 hover:dark:bg-neutral-700"
    :class="{ 'bg-slate-100 dark:bg-neutral-800': props.memo.pinned }"
  >
    <div class="flex flex-col w-24 p-2 text-center">
      <template v-if="!isPinned">
        <div v-if="props.memo.displayDate" class="flex justify-center">
          <span v-if="$dayjs(props.memo.createdAt).isSame($dayjs(), 'day')" class="text-lg font-bold">今天</span>
          <span v-else-if="$dayjs(props.memo.createdAt).isSame($dayjs().subtract(1, 'day'), 'day')" class="text-lg font-bold">昨天</span>
          <template v-else>
            <span class="text-xl font-bold">{{ $dayjs(props.memo.createdAt).format("DD") }}</span>
            <span class="flex items-end text-xs">{{ $dayjs(props.memo.createdAt).format("MM") }}月</span>
          </template>
        </div>
        <div class="flex justify-center text-[#576b95] font-medium dark:text-white text-xs mt-2 select-none">
          {{ location }}
        </div>
      </template>
      <div v-else class="flex justify-center items-center">
        <span class="text-lg">置顶</span>
      </div>
    </div>
    <div class="flex w-full flex-col pr-4 py-2">
      <NuxtLink class="flex" :to="`/memo/${item.id}`">
        <div
          class="sm:w-24 sm:h-24 w-20 h-20 relative overflow-hidden flex-shrink-0"
          v-if="imageCount > 0 || hasMediaContent"
        >
          <div
            v-if="imageCount > 0"
            :class="getImageGridClass(imageCount)"
            class="h-full w-full gap-0.1 grid"
          >
            <div
              v-for="(img, index) in images.slice(0, Math.min(imageCount, 9))"
              :key="index"
              :class="getGridClass(index, imageCount)"
              class="border border-white dark:border-neutral-800 relative"
            >
              <img
                :src="img"
                class="absolute inset-0 w-full h-full object-cover"
              />
            </div>
          </div>

          <div v-else-if="hasMediaContent" class="w-full h-full bg-gradient-to-br from-blue-100 to-purple-100 dark:from-blue-900/30 dark:to-purple-900/30 rounded flex items-center justify-center border border-gray-200 dark:border-gray-700 relative overflow-hidden">
            <div class="absolute inset-0 bg-gradient-to-br from-white/20 to-transparent"></div>

            <div class="relative z-1">
              <div class="w-8 h-8 sm:w-10 sm:h-10 bg-white/90 dark:bg-gray-800/90 rounded-full flex items-center justify-center shadow-lg backdrop-blur-sm">
                <UIcon :name="getMediaIcon()" class="w-4 h-4 sm:w-5 sm:h-5 text-blue-600 dark:text-blue-400" />
              </div>
            </div>

            <div class="absolute top-2 right-2 w-2 h-2 bg-blue-400/60 rounded-full"></div>
            <div class="absolute bottom-2 left-2 w-1 h-1 bg-purple-400/60 rounded-full"></div>
          </div>
        </div>

        <div class="flex-1 flex flex-col justify-between">
          <div>
            <div
              class="markdown-content bg-neutral-100 dark:bg-neutral-800 p-2 sm:pb-2 pb-1 !leading-7 line-clamp-2 sm:line-clamp-3"
              v-if="(imageCount === 0 && !hasMediaContent) && content"
              v-html="content"
            ></div>
            <div
              class="markdown-content ml-1 !leading-5 line-clamp-2 sm:line-clamp-3"
              v-if="(imageCount > 0 || hasMediaContent) && content"
              v-html="content"
            ></div>
            <div
              class="text-gray-500 dark:text-gray-400 text-sm pl-2"
              v-if="!content"
            >
              无文字描述
            </div>
          </div>

          <div class="text-sm text-gray-500 ml-1" v-if="imageCount > 0 || hasMediaContent">
            <span v-if="imageCount > 0">有{{ imageCount }}图</span>
            <span v-if="imageCount > 0 && hasMediaContent"> · </span>
            <span v-if="hasMediaContent">{{ getMediaType() }}</span>
          </div>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ExtDTO, MemoVO } from "~/types";
import { md } from "~/utils";
import { computed } from "vue";
import { toast } from "vue-sonner";

const props = defineProps<{
  memo: MemoVO;
}>();

const item = computed(() => props.memo);
const isPinned = computed(() => item.value.pinned);
const location = computed(() =>
  (item.value.location || "").replaceAll(" ", " · ")
);

const extJSON = computed(() => {
  try {
    return JSON.parse(item.value.ext || "{}") as ExtDTO;
  } catch (error) {
    console.error("解析 ext 字段时出错:", error);
    return {} as ExtDTO;
  }
});

const content = computed(() => {
  if (item.value.content && item.value.content.length > 0) {
    try {
      return md.render(item.value.content);
    } catch (error) {
      console.error("内容渲染错误，请重新编辑:", error);
      toast.error("内容渲染错误，请重新编辑");
      return "内容渲染错误，请重新编辑";
    }
  }
  return "";
});

const imageCount = computed(() => {
  const imgs = item.value.imgs || "";
  return imgs.split(",").filter(Boolean).length;
});

const images = computed(() => {
  const imgs = item.value.imgs || "";
  return imgs.split(",").filter(Boolean);
});

const gridRules: Record<number, Record<number, string>> = {
  2: {0: 'col-span-1', 1: 'col-span-1'},
  3: {0: 'row-span-2', 1: 'row-span-1', 2: 'row-span-1'},
  5: {0: 'col-span-2 row-span-2',1: 'col-span-1 row-span-2'},
  6: {0: 'col-span-2 row-span-2'},
  7: {0: 'col-span-1 row-span-2',1: 'col-span-1 row-span-2'},
  8: {0: 'col-span-1 row-span-2'}
} as const;

const getImageGridClass = (count: number) => {
  if (count <= 1) return '';
  if (count === 2) return 'grid grid-cols-2';
  if (count === 3) return 'grid grid-cols-2';
  if (count <= 4) return 'grid grid-cols-2 grid-rows-2';
  return 'grid grid-cols-3 grid-rows-3';
};

const getGridClass = (index: number, count: keyof typeof gridRules): string => {
  return gridRules[count]?.[index] || '';
};

// 检查是否有媒体内容
const hasMediaContent = computed(() => {
  return (
    (item.value.externalFavicon && item.value.externalTitle && item.value.externalUrl) ||
    (extJSON.value.music && extJSON.value.music.id) ||
    (extJSON.value.doubanBook && extJSON.value.doubanBook.title) ||
    (extJSON.value.doubanMovie && extJSON.value.doubanMovie.title) ||
    (extJSON.value.video && 
      (['bilibili', 'youtube'].includes(extJSON.value.video.type) || extJSON.value.video.type === 'online') &&
      extJSON.value.video.value)
  );
});

// 获取媒体类型图标
const getMediaIcon = () => {
  if (item.value.externalFavicon && item.value.externalTitle && item.value.externalUrl) {
    return 'i-carbon-link';
  }
  if (extJSON.value.music && extJSON.value.music.id) {
    return 'i-carbon-play-filled';
  }
  if (extJSON.value.doubanBook && extJSON.value.doubanBook.title) {
    return 'i-carbon-book';
  }
  if (extJSON.value.doubanMovie && extJSON.value.doubanMovie.title) {
    return 'i-carbon-video';
  }
  if (extJSON.value.video && extJSON.value.video.value) {
    return 'i-carbon-play';
  }
  return 'i-carbon-media';
};

// 获取媒体类型文字
const getMediaType = () => {
  if (item.value.externalFavicon && item.value.externalTitle && item.value.externalUrl) {
    return '链接';
  }
  if (extJSON.value.music && extJSON.value.music.id) {
    return '音乐';
  }
  if (extJSON.value.doubanBook && extJSON.value.doubanBook.title) {
    return '豆瓣读书';
  }
  if (extJSON.value.doubanMovie && extJSON.value.doubanMovie.title) {
    return '豆瓣电影';
  }
  if (extJSON.value.video && extJSON.value.video.value) {
    return '视频';
  }
  return '媒体';
};
</script>

<style scoped></style>