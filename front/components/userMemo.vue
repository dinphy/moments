<template>
  <div
    class="flex flex-row sm:gap-4 text-sm sm:py-2 sm:px-4 w-full"
    :class="{ 'bg-slate-100 dark:bg-neutral-800': props.memo.pinned }"
  >
    <div class="flex flex-col w-2/5 sm:w-1/5 p-2">
      <template v-if="!isPinned">
        <div>
          <span class="text-lg">{{ formattedDate.day }}</span>
          <span>{{ formattedDate.month }}月</span>
        </div>
        <div
          class="text-[#576b95] font-medium dark:text-white text-xs mt-2 mb-1 select-none"
        >
          {{ location }}
        </div>
      </template>
      <div v-else class="flex items-center">
        <span class="text-lg">置顶</span>
      </div>
    </div>
    <div class="flex w-full flex-col">
      <div class="flex">
        <div class="w-32 h-32" v-if="imageCount > 0">
          <upload-image-preview
            :imgs="item.imgs"
            :imgConfigs="item.imgConfigs"
            :memo-id="item.id"
          />
        </div>
        <div class="flex-1 flex flex-col justify-between">
          <div
            class="markdown-content bg-neutral-100 dark:bg-neutral-800 p-2 pb-1"
            ref="contentRef"
            v-if="imageCount === 0"
            v-html="content"
          ></div>
          <div
            class="markdown-content ml-2"
            ref="contentRef"
            v-if="imageCount > 0"
            v-html="content"
          ></div>
          <div
            v-if="imageCount > 0"
            class="image-count text-sm text-gray-500 mt-1 ml-2"
          >
            有{{ imageCount }}图
          </div>
        </div>
      </div>
      <div class="flex flex-col gap-2 mt-2">
        <external-url-preview
          v-if="hasExternalUrl"
          :favicon="item.externalFavicon"
          :title="item.externalTitle"
          :url="item.externalUrl"
        />
        <music-preview v-if="hasMusic" v-bind="extJSON.music" />
        <douban-book-preview v-if="hasDoubanBook" :book="extJSON.doubanBook" />
        <douban-movie-preview
          v-if="hasDoubanMovie"
          :movie="extJSON.doubanMovie"
        />
        <video-preview-iframe
          v-if="hasVideoIframe"
          :url="extJSON.video.value"
        />
        <video-preview v-if="hasVideo" :url="extJSON.video.value" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ExtDTO, MemoVO, SysConfigVO } from "~/types";
import { md } from "~/utils";
import { useGlobalState } from "~/store";
import { ref, computed, onMounted } from "vue";
import { toast } from "vue-sonner";

const router = useRouter();
const sysConfig = useState<SysConfigVO>("sysConfig");
const props = defineProps<{
  memo: MemoVO;
}>();

const item = computed(() => props.memo);
const isPinned = computed(() => item.value.pinned);
const formattedDate = computed(() => {
  const date = new Date(item.value.createdAt);
  return {
    day: date.getDate().toString().padStart(2, "0"),
    month: (date.getMonth() + 1).toString().padStart(2, "0"),
  };
});
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

const hasExternalUrl = computed(
  () =>
    item.value.externalFavicon &&
    item.value.externalTitle &&
    item.value.externalUrl
);
const hasMusic = computed(() => extJSON.value.music && extJSON.value.music.id);
const hasDoubanBook = computed(
  () => extJSON.value.doubanBook && extJSON.value.doubanBook.title
);
const hasDoubanMovie = computed(
  () => extJSON.value.doubanMovie && extJSON.value.doubanMovie.title
);
const hasVideoIframe = computed(
  () =>
    extJSON.value.video &&
    ["bilibili", "youtube"].includes(extJSON.value.video.type) &&
    extJSON.value.video.value
);
const hasVideo = computed(
  () =>
    extJSON.value.video &&
    extJSON.value.video.type === "online" &&
    extJSON.value.video.value
);

const contentRef = ref<HTMLElement | null>(null);

onMounted(() => {
  if (contentRef.value) {
    contentRef.value.classList.add("line-clamp-3");
  }
});
</script>

<style scoped>
.upload-image-preview img {
  width: 100%;
  height: auto;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.image-count {
  margin-top: 4px;
}
</style>
