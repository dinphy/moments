<template>
  <div>
    <Header v-bind:user="currentUser" />

    <div v-if="loading" class="flex justify-center items-center py-10">
      <UIcon name="i-carbon-circle-dash" class="animate-spin text-3xl text-gray-500" />
    </div>

    <div v-else-if="groupedImages.length > 0" class="px-2 sm:px-4 pb-8">
      <div v-for="(yearGroup, yearIndex) in groupedImages" :key="yearGroup.year" class="mb-8">
        <div class="sticky top-14 z-20 bg-white dark:bg-gray-900 py-3 mb-4 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-xl font-bold text-gray-800 dark:text-white">{{ yearGroup.year }}年</h2>
        </div>

        <div v-for="(monthGroup, monthIndex) in yearGroup.months" :key="monthGroup.month" class="mb-6">
          <div class="flex items-center mb-3">
            <h3 class="text-lg font-medium text-gray-700 dark:text-gray-300">{{ monthGroup.month }}月</h3>
            <span class="ml-2 text-sm text-gray-500 dark:text-gray-400">({{ monthGroup.images.length }}张)</span>
          </div>

          <div class="grid grid-cols-3 gap-1 sm:gap-2">
            <div 
              v-for="(image, imgIndex) in monthGroup.images" 
              :key="image.id"
              class="relative w-full pb-[100%] overflow-hidden rounded-lg cursor-pointer group bg-gray-100 dark:bg-gray-800"
              @click="openImagePreview(yearIndex, monthIndex, imgIndex)"
            >
              <img 
                :src="image.url" 
                :alt="image.memoContent"
                class="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                loading="lazy"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <div class="absolute bottom-0 left-0 right-0 p-2 text-white text-xs truncate opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                {{ image.memoContent || '无描述' }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div
        v-if="hasNext"
        ref="loadMoreEle"
        class="text-xs text-center text-gray-500 py-2 cursor-pointer"
        @click="loadMore"
      >
        点击加载更多
      </div>
      <div class="text-xs text-center text-gray-500 py-2" v-else>
        已经到底啦
      </div>
    </div>

    <div v-else class="flex flex-col items-center justify-center py-16">
      <UIcon name="i-carbon-image" class="text-6xl text-gray-300 dark:text-gray-600 mb-4" />
      <p class="text-gray-500 dark:text-gray-400">暂无图片</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MemoVO, SysConfigVO, UserVO } from "~/types";
import dayjs from 'dayjs';
import { useElementVisibility } from "@vueuse/core";

const currentUser = useState<UserVO>('userinfo');
const loading = ref(true);
const loadingMore = ref(false);
const hasNext = ref(false);
const loadMoreEle = ref(null);
const targetIsVisible = useElementVisibility(loadMoreEle);
const sysConfig = useState<SysConfigVO>("sysConfig");

watch(targetIsVisible, async (visible) => {
  if (visible && sysConfig.value?.enableAutoLoadNextPage && !loadingMore.value) {
    await loadMore();
  }
});

const state = reactive({
  page: 1,
  size: 50,
});

const allMemos = ref<Array<MemoVO>>([]);

const groupedImages = computed(() => {
  const images: Array<{
    id: number;
    url: string;
    memoId: number;
    memoContent: string;
    createdAt: string;
    displayDate: string;
  }> = [];

  allMemos.value.forEach(memo => {
    if (memo.imgs) {
      const imgUrls = memo.imgs.split(',').filter(Boolean);
      imgUrls.forEach(url => {
        images.push({
          id: memo.id,
          url,
          memoId: memo.id,
          memoContent: memo.content,
          createdAt: memo.createdAt,
          displayDate: dayjs(memo.createdAt).format('YYYY-MM-DD'),
        });
      });
    }
  });

  images.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime());

  const result: any[] = [];
  let currentYear: number | null = null;
  let currentYearGroup: { months: any; year?: number; } | null = null;
  let currentMonth: number | null = null;
  let currentMonthGroup: { images: any; month?: number; } | null = null;

  images.forEach(image => {
    const date = dayjs(image.createdAt);
    const year = date.year();
    const month = date.month() + 1;

    if (currentYear !== year) {
      currentYear = year;
      currentYearGroup = {
        year,
        months: [],
      };
      result.push(currentYearGroup);
      currentMonth = null;
      currentMonthGroup = null;
    }

    if (currentMonth !== month) {
      currentMonth = month;
      currentMonthGroup = {
        month,
        images: [],
      };
      currentYearGroup?.months.push(currentMonthGroup);
    }

    currentMonthGroup?.images.push(image);
  });

  return result;
});

const getMonthImagesForPreview = (yearIndex: number, monthIndex: number) => {
  const monthGroup = groupedImages.value[yearIndex].months[monthIndex];
  return monthGroup.images.map(image => ({
    src: image.url,
    caption: `
      <div class="absolute bottom-0 left-0 right-0 p-3 bg-black/75 text-white backdrop-blur-sm">
        <div class="line-clamp-2 mb-2 text-sm leading-6">${image.memoContent || '暂无描述'}</div>
        <div class="flex justify-between items-center text-xs">
          <span class="text-white/80">${image.displayDate}</span>
          <a href="/memo/${image.memoId}" class="flex items-center gap-1 text-white opacity-80 hover:opacity-100 transition-opacity">
            详情
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </a>
        </div>
      </div>
    `,
    thumb: image.url,
  }));
};

const openImagePreview = (yearIndex: number, monthIndex: number, imgIndex: number) => {
  const currentMonthImages = getMonthImagesForPreview(yearIndex, monthIndex);
  
  import('@fancyapps/ui').then(({ Fancybox }) => {
    Fancybox.show(currentMonthImages, {
      startIndex: imgIndex,
      Thumbs: {
        type: "modern",
      },
      Toolbar: {
        display: {
          left: ["zoom"],
          middle: [],
          right: ["slideshow", "fullscreen", "download", "thumbs", "close"]
        }
      },
      caption: (fancybox, slide) => {
        return slide.caption || false;
      },
      on: {
        init: () => {
          const style = document.createElement('style');
          style.innerHTML = `
            .fancybox__caption {
              z-index: 9999 !important;
            }
          `;
          document.head.appendChild(style);
        }
      }
    });
  });
};

const loadInitialData = async () => {
  loading.value = true;
  try {
    state.page = 1;
    const res = await useMyFetch<{
      list: Array<MemoVO>,
      total: number,
      hasNext: boolean
    }>('/memo/list', state);
    allMemos.value = res.list;
    hasNext.value = res.hasNext;
  } catch (error) {
    console.error('加载相册数据失败:', error);
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
      list: Array<MemoVO>,
      total: number,
      hasNext: boolean
    }>('/memo/list', state);
    allMemos.value = [...allMemos.value, ...res.list];
    hasNext.value = res.hasNext;
  } catch (error) {
    console.error('加载更多相册数据失败:', error);
  } finally {
    loadingMore.value = false;
  }
};

onMounted(async () => {
  await loadInitialData();
});
</script>
