<template>
  <div
    class="relative flex gap-4 text-sm dark:bg-neutral-800 p-4"
    :class="[item.pinned ? 'bg-slate-100 dark:bg-neutral-700' : '']"
  >
    <div class="avatar">
      <NuxtLink
        :to="isDetailPage ? '' : `/memo/${item.id}`"
      >
        <UAvatar :src="item.user.avatarUrl" alt="Avatar" />
      </NuxtLink>
    </div>
    <div class="flex flex-col gap-1 flex-1">
      <div
        class="username text-[#576b95] mb-1 dark:text-white flex justify-between"
      >
        <NuxtLink class="cursor-pointer" :to="`/user/${item.user.id}`">
          {{ item.user.nickname }}
        </NuxtLink>
        <div>
          <UIcon v-if="item.pinned" name="i-carbon-pin" />
          <UIcon
            v-if="item.showType === 0"
            name="i-carbon-locked"
            class="text-red-500 ml-2 dark:text-white"
          />
        </div>
      </div>
      <div class="mb-2">
        <div :style="getMemoMaxHeightStyle()" class="overflow-hidden">
          <div
            class="markdown-content"
            ref="contentRef"
            v-html="content"
          ></div>
        </div>
        <div
          v-if="showMore"
          class="text-[#576b95] text-sm my-1 cursor-pointer"
          @click="doShowMore"
        >
          {{ getMemoMaxHeightStyle() === "" ? "收起" : "全文" }}
        </div>
        <div v-if="tags.length > 0" class="flex flex-wrap gap-2 mt-3">
          <span v-for="(tag, index) in tags" :key="`tag-${index}`">
            <NuxtLink :to="`/tags/${item.user.username}/${tag}`">
              <UBadge 
                size="sm" 
                color="gray" 
                variant="soft"
                class="px-3 py-1.5 text-xs transition-all duration-300 hover:scale-105 hover:shadow-lg cursor-pointer bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-800/50 dark:to-gray-700/50 hover:from-gray-100 hover:to-gray-200 dark:hover:from-gray-700/70 dark:hover:to-gray-600/70 border border-gray-200/50 dark:border-gray-600/30"
              >
                <UIcon name="i-carbon-hashtag" />
                <span class="text-gray-700 dark:text-gray-200">{{ tag }}</span>
              </UBadge>
            </NuxtLink>
          </span>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <external-url-preview
          v-if="
            item.externalFavicon && item.externalTitle && item.externalUrl
          "
          :favicon="item.externalFavicon"
          :title="item.externalTitle"
          :url="item.externalUrl"
        />
        <upload-image-preview
          :imgs="item.imgs"
          :imgConfigs="item.imgConfigs"
          :memo-id="item.id"
        />

        <music-preview
          v-if="extJSON.music && extJSON.music.id"
          v-bind="extJSON.music"
        />
        <douban-book-preview
          v-if="extJSON.doubanBook && extJSON.doubanBook.title"
          :book="extJSON.doubanBook"
        />
        <douban-movie-preview
          v-if="extJSON.doubanMovie && extJSON.doubanMovie.title"
          :movie="extJSON.doubanMovie"
        />
        <video-preview-iframe
          v-if="
            extJSON.video &&
            ['bilibili', 'youtube'].includes(extJSON.video.type) &&
            extJSON.video.value
          "
          :url="extJSON.video.value"
        />
        <video-preview
          v-if="
            extJSON.video &&
            extJSON.video.type === 'online' &&
            extJSON.video.value
          "
          :url="extJSON.video.value"
        />
      </div>

      <div
        v-if="location"
        class="text-[#576b95] font-medium dark:text-white text-xs mt-2 mb-1 select-none flex items-center gap-0.5"
      >
        <UIcon name="i-carbon-location" />
        <span>{{ location }}</span>
      </div>

      <div class="flex justify-between items-center relative">
        <div class="flex text-xs text-[#9DA4B0]">
          {{
            sysConfig.timeFormat === "timeAgo"
              ? $dayjs(item.createdAt).fromNow()
              : $dayjs(item.createdAt).format("YYYY-MM-DD HH:mm")
          }}
          {{
            $dayjs(item.createdAt).isAfter($dayjs()) ? '，未到发布时间，仅自己可见' : ''
          }}
        </div>
        <div
          @click="showToolbar = !showToolbar"
          class="toolbar-icon px-2 py-1 bg-[#f7f7f7] dark:bg-slate-700 hover:bg-[#dedede] cursor-pointer rounded flex items-center justify-center"
        >
          <img
            class="w-3 h-3"
            src="data:image/svg+xml,%3csvg%20t='1709204592505'%20class='icon'%20viewBox='0%200%201024%201024'%20version='1.1'%20xmlns='http://www.w3.org/2000/svg'%20p-id='16237'%20width='16'%20height='16'%3e%3cpath%20d='M229.2%20512m-140%200a140%20140%200%201%200%20280%200%20140%20140%200%201%200-280%200Z'%20p-id='16238'%20fill='%238a8a8a'%3e%3c/path%3e%3cpath%20d='M794.8%20512m-140%200a140%20140%200%201%200%20280%200%20140%20140%200%201%200-280%200Z'%20p-id='16239'%20fill='%238a8a8a'%3e%3c/path%3e%3c/svg%3e"
          />
        </div>

        <div
          v-if="showToolbar"
          ref="toolbarRef"
          class="absolute top-[-8px] right-[32px] bg-[#4c4c4c] rounded text-white p-2"
        >
          <div class="flex flex-row gap-2">
            <div
              class="flex flex-row gap-1 cursor-pointer items-center px-4"
              @click="liked ? unlikeMemo(item.id) : likeMemo(item.id)"
            >
              <UIcon
                v-if="liked"
                name="i-weui-like-filled"
                class="w-5 h-5 text-red-400"
              />
              <UIcon v-else name="i-weui-like-outlined" class="w-5 h-5" />
              <div>{{ liked ? "取消" : "赞" }}</div>
            </div>
            <template v-if="sysConfig.enableComment">
              <span class="bg-[#6b7280] h-[20px] w-[1px]"></span>
              <div
                class="flex flex-row gap-1 cursor-pointer items-center px-4"
                @click="doComment"
              >
                <UIcon
                  name="i-weui-comment-outlined"
                  class="w-5 h-5 relative"
                />
                <div>评论</div>
              </div>
            </template>
          </div>
        </div>
      </div>

      <div
        class="rounded bottom-shadow bg-[#f7f7f7] dark:bg-[#202020] flex flex-col gap-1"
      >
        <div
          v-if="likeInfo && likeInfo.length > 0"
          class="flex flex-row py-2 px-3 gap-2 items-center text-sm"
          :class="[
            item.comments && item.comments.length > 0
              ? 'border-b-[1px] border-neutral-[100] dark:border-neutral-800'
              : '',
          ]"
        >
          <div class="text-[#576b95] gap-1">
            <UIcon name="i-carbon-favorite" class="mr-1 relative top-[1px]" />
            {{ visibleLikeInfo.map(info => info.name).join(', ') }}
            <span v-if="likeNum > visibleLikeInfo.length">
              {{ visibleLikeInfo.length > 0 ? ', ' : '' }}{{ guestLikeNum }}位访客
            </span>
          </div>
        </div>
        <div class="flex flex-col gap-1" v-if="sysConfig.enableComment">
          <CommentBox :comment-id="0" :memo-id="item.id" :memo-user-id="item.user.id" />
          <div
            class="space-y-1"
            :class="[item.comments && item.comments.length > 0 ? 'py-2' : '']"
          >
            <div
              v-if="item.comments && item.comments.length > 0"
              v-for="c in item.comments"
              :key="c.id"
              class="px-3 relative flex-col text-sm"
            >
              <Comment
                :comment="c"
                :memo-id="item.id"
                :memo-user-id="item.user.id"
                :is-detail-page="isDetailPage"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ExtDTO, MemoVO, SysConfigVO } from "~/types";
import { toast } from "vue-sonner";
import { memoChangedEvent, memoReloadEvent, messageChangedEvent } from "~/event";
import Comment from "~/components/Comment.vue";
import { useGlobalState } from "~/store";
import { md, getGuestId } from "~/utils";
import {useStorage} from '@vueuse/core'

const showMore = ref(false);
const showMoreClicked = ref(false);
const isDetailPage = computed(() => {
  return route.path.startsWith("/memo/");
});
const contentRef = ref<HTMLDivElement | null>(null);
const sysConfig = useState<SysConfigVO>("sysConfig");
const route = useRoute();
const router = useRouter();
const { y } = useWindowScroll();

const getMemoMaxHeightStyle = () => {
  if (isDetailPage.value || showMoreClicked.value) {
    return "";
  }
  if (sysConfig.value.memoMaxHeight) {
    return `max-height:${sysConfig.value.memoMaxHeight}px`;
  }
  return "";
};

const currentCommentBox = useState("currentCommentBox");
const props = defineProps<{
  memo: MemoVO;
}>();
const extJSON = computed(() => {
  return JSON.parse(props.memo.ext || "{}") as ExtDTO;
});
const item = computed(() => {
  return props.memo;
});

const global = useGlobalState();
const showToolbar = ref(false);
const toolbarRef = ref(null);

onClickOutside(toolbarRef, () =>
  setTimeout(() => {
    showToolbar.value = false;
  }, 10)
);

const location = computed(() => {
  return (item.value.location || "").replaceAll(" ", " · ");
});

const tags = computed(() => {
  const tagsStr = item.value.tags;
  if (!tagsStr) {
    return [];
  }
  const len = tagsStr.length;
  if (tagsStr[len - 1] === ",") {
    return tagsStr.substring(0, len - 1).split(",");
  }
  return tagsStr.split(",");
});

const doComment = () => {
  const value = item.value.id + "#0";
  if (currentCommentBox.value === value) {
    currentCommentBox.value = "";
  } else {
    currentCommentBox.value = value;
  }
  showToolbar.value = false;
};

const doShowMore = () => {
  showMoreClicked.value = !showMoreClicked.value;
};

const liked = ref(false);
const likeInfo = ref<{ id: number | string; name: string }[] | null>(null);
const likeNum = ref(0);
const likeShowAll = ref(false);
const isLoading = ref(false);
const localCommentUserinfo = useStorage('localCommentUserinfo', {
  username: "",
  website: "",
  email: "",
})

const visibleLikeInfo = computed(() => {
  if (!likeInfo.value) return [];
  return likeInfo.value.filter(info => info.name && info.name !== info.id);
});

const guestLikeNum = computed(() => {
  if (!likeInfo.value) return 0;
  return likeInfo.value.filter(info => typeof info.id === 'string').length;
});

const doLike = async (params: string) => {
  showToolbar.value = false;
  try {
    await useMyFetch(`/like/add?${params}`);
    toast.success("点赞成功!");
    liked.value = true;
  } catch (error) {
    console.error("点赞失败:", error);
    toast.warning("点赞失败，请稍后重试！");
  }
};

const likeMemo = async (id: number) => {
  if (isLoading.value) return;
  isLoading.value = true;

  try {
    const guestId = await getGuestId();
    if (!guestId) return;
    let params = `id=${id}&guestId=${guestId}`;
    const guestName = localCommentUserinfo.value.username;
    if (guestName) {
      params += `&guestName=${guestName}`;
    }

    if (sysConfig.value.enableGoogleRecaptcha) {
      await new Promise<void>((resolve) => {
        grecaptcha.ready(() => {
          grecaptcha
            .execute(sysConfig.value.googleSiteKey, { action: "newComment" })
            .then(async (token) => {
              params += `&token=${token}`;
              await doLike(params);
              await getLike(id);
              memoChangedEvent.emit(id);
              if (global.value.userinfo.id !== item.value.userId) {
                messageChangedEvent.emit(1);
              }
              resolve();
            });
        });
      });
    } else {
      await doLike(params);
      await getLike(id);
      memoChangedEvent.emit(id);
      if (global.value.userinfo.id !== item.value.userId) {
        messageChangedEvent.emit(1);
      }
    }
  } finally {
    isLoading.value = false;
  }
};

const doUnlike = async (params: string) => {
  showToolbar.value = false;
  if (!global.value.userinfo.token) {
    toast.warning("访客不允许取消点赞！");
    return false;
  }
  try {
    await useMyFetch(`/like/remove?${params}`);
    toast.success("取消点赞成功!");
    liked.value = false;
    return true;
  } catch (error) {
    console.error("取消点赞失败:", error);
    toast.warning("取消点赞失败，请稍后重试！");
    return false;
  }
};

const unlikeMemo = async (id: number) => {
  if (isLoading.value) return;
  isLoading.value = true;

  try {
    const guestId = await getGuestId();
    if (!guestId) return;
    let params = `id=${id}&guestId=${guestId}`;
    const guestName = localCommentUserinfo.value.username;
    if (guestName) {
      params += `&guestName=${guestName}`;
    }

    if (sysConfig.value.enableGoogleRecaptcha) {
      await new Promise<void>((resolve) => {
        grecaptcha.ready(() => {
          grecaptcha
            .execute(sysConfig.value.googleSiteKey, { action: "newComment" })
            .then(async (token) => {
              params += `&token=${token}`;
              const success = await doUnlike(params);
              if (success) {
                await getLike(id);
                memoChangedEvent.emit(id);
                if (global.value.userinfo.id !== item.value.userId) {
                  messageChangedEvent.emit(-1);
                }
              }
              resolve();
            });
        });
      });
    } else {
      const success = await doUnlike(params);
      if (success) {
        await getLike(id);
        memoChangedEvent.emit(id);
        if (global.value.userinfo.id !== item.value.userId) {
          messageChangedEvent.emit(-1);
        }
      }
    }
  } finally {
    isLoading.value = false;
  }
};

const getLike = async (id: number) => {
  const guestId = await getGuestId();
  if (!guestId) return;
  let params = `id=${id}&guestId=${guestId}`;

  try {
    const response = await useMyFetch<{
      likes: { id: number | string; name: string }[];
      total: number;
    }>(`/like/get?${params}`);
    likeInfo.value = (response.likes || []).sort((a, b) => {
      return Number(typeof b.id === 'number') - Number(typeof a.id === 'number');
    });
    likeNum.value = response.total;
    if (global.value.userinfo.token) {
      const userId = global.value.userinfo.id;
      liked.value = likeInfo.value?.some((info) => info.id === userId) || false;
    } else {
      liked.value = likeInfo.value?.some((info) => info.id === guestId) || false;
    }
    return true;
  } catch (error) {
    console.error("获取点赞信息失败:", error);
    toast.error("获取点赞信息失败，请稍后重试！");
    return false;
  }
};

onMounted(async () => {
  await getLike(item.value.id);

  if (!isDetailPage.value) {
    setTimeout(() => {
      const { height } = useElementSize(contentRef.value);
      if (height.value > sysConfig.value.memoMaxHeight) {
        showMore.value = true;
      }
    }, 20);
  }
});

const content = computed(() => {
  if (item.value.content && item.value.content.length > 0) {
    try {
      return md.render(item.value.content);
    } catch (e) {
      console.log("内容渲染错误,请重新编辑", e);
      return "内容渲染错误,请重新编辑";
    }
  }
  return "";
});
const goBack = () => {
  if (window.history.length > 1) {
    router.back();
  } else {
    navigateTo("/");
  }
};
</script>

<style lang="scss" scoped></style>
