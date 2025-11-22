<template>
  <div>
    <div v-if="loading" class="flex justify-center items-center py-10">
      <UIcon name="i-carbon-circle-dash" class="animate-spin text-3xl text-gray-500" />
    </div>

    <div v-else-if="groupedImages.length > 0" class="px-2 sm:px-4 pb-8">
      <div v-for="(yearGroup, yearIndex) in groupedImages" :key="yearGroup.year" class="mb-8">
        <h2 class="py-4">
          <span class="text-xl">{{ yearGroup.year }}</span>
          <span class="text-sm">年</span>
        </h2>

        <div v-for="(monthGroup, monthIndex) in yearGroup.months" :key="monthGroup.month" class="mb-6">
          <div class="flex flex-row">
            <div class="flex flex-col w-24 p-2 text-center">
              <div class="flex">
                <span class="text-xl font-bold">{{ dayjs().month(monthGroup.month - 1).format("MM") }}</span>
                <span class="flex items-end text-xs">月</span>
              </div>
              <div class="flex justify-center text-[#576b95] font-medium dark:text-white text-xs mt-2 select-none">
                {{ monthGroup.images.length }}图
              </div>
            </div>

            <div class="flex w-full flex-col pr-4 py-2">
              <div class="grid grid-cols-3 gap-1 sm:gap-2">
                <div 
                  v-for="(image, imgIndex) in monthGroup.images" 
                  :key="image.id"
                  class="relative w-full pb-[100%] overflow-hidden rounded-lg cursor-pointer group bg-gray-100 dark:bg-gray-800"
                  :data-fancybox="`gallery-${yearIndex}-${monthIndex}`"
                  :data-src="image.url"
                  :data-caption="image.memoContentHtml || '暂无描述'"
                  :data-date="image.displayDate"
                  :data-memo-id="image.memoId"
                >
                  <img 
                    :src="image.url" 
                    :alt="image.memoContent"
                    class="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                    loading="lazy"
                  />
                  <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                  <div class="absolute bottom-0 left-0 right-0 p-2 text-white text-xs truncate opacity-0 group-hover:opacity-100 transition-opacity duration-300" v-html="image.memoContentHtml || '暂无描述'"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
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
import { Fancybox } from '@fancyapps/ui';
import { md } from '~/utils';

const props = defineProps({
  memos: {
    type: Array as () => Array<MemoVO>,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const currentUser = useState<UserVO>('userinfo');
const allMemos = computed(() => props.memos || []);
const loading = computed(() => props.loading);

// 初始化 Fancybox
onMounted(() => {
  Fancybox.bind('[data-fancybox]', {
    Thumbs: false,
    Toolbar: {
      display: {
        left: ["zoom", "infobar"],
        middle: [],
        right: ["slideshow", "fullscreen", "close"]
      }
    },
    caption: function (fancybox, slide) {
      const caption = slide.caption || '';
      const date = slide.triggerEl?.dataset.date || '';
      const memoId = slide.triggerEl?.dataset.memoId || '';
      
      return `
        <div class="p-4 text-white fixed bottom-0 left-0 right-0 z-50 bg-black/50 backdrop-blur-sm text-xs">
          <div class="line-clamp-2 mb-2 text-sm leading-6">${caption}</div>
          <div class="flex justify-between items-center text-xs">
            <span class="text-gray-300">${date}</span>
            <a href="/memo/${memoId}" class="flex items-center gap-1 text-white opacity-80 hover:opacity-100 transition-opacity">
              详情
              <svg class="w-4 h-4 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </a>
          </div>
        </div>
      `;
    }
  });
});

// 组件卸载时解绑
onUnmounted(() => {
  Fancybox.unbind('[data-fancybox]');
  Fancybox.close();
});

// 处理内容，渲染emoji
const processContent = (content: string) => {
  if (content && content.length > 0) {
    try {
      return md.render(content);
    } catch (e) {
      console.log("内容渲染错误", e);
      return content;
    }
  }
  return content;
};

const groupedImages = computed(() => {
  const images: Array<{
    id: number;
    url: string;
    memoId: number;
    memoContent: string;
    memoContentHtml: string;
    createdAt: string;
    displayDate: string;
  }> = [];

  allMemos.value.forEach(memo => {
      if (memo.imgs) {
        const imgUrls = memo.imgs.split(',').filter(Boolean);
        const processedContent = processContent(memo.content || '');
        imgUrls.forEach(url => {
          images.push({
            id: memo.id,
            url,
            memoId: memo.id,
            memoContent: memo.content,
            memoContentHtml: processedContent,
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
</script>
