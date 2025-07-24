<template>
  <div
    v-if="unreadCount > 0"
    class="absolute top-4 left-4 cursor-pointer z-10"
    @click="toggleMessageBox"
  >
    <UIcon name="i-carbon-notification" class="text-[#9fc84a] w-5 h-5" />
    <span
      class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center"
    >
      {{ unreadCount }}
    </span>
  </div>

  <UModal
    v-model="showMessageBox"
    :ui="{
      container: 'flex justify-center items-center backdrop-blur',
    }"
  >
    <div class="flex justify-between items-center p-5">
      <h3 class="font-medium">消息盒子</h3>
      <div
        @click="unreadCount > 0 ? markAllAsRead() : null"
        :class="{
          'text-blue-500 cursor-pointer': unreadCount > 0,
          'text-gray-400 cursor-not-allowed': unreadCount === 0,
        }"
        class="flex items-center text-sm"
      >
        <UIcon name="i-carbon-notification-new" class="w-4 h-4 mr-1" />
        <span v-if="unreadCount > 0">标为已读({{ unreadCount }})</span>
        <span v-else>全部已读</span>
      </div>
    </div>

    <div
      v-if="!messages || messages.length === 0"
      class="p-4 text-center text-gray-500 dark:text-gray-400"
    >
      暂无消息
    </div>

    <div v-else class="max-h-96 overflow-y-auto mb-5">
      <div
        v-for="message in messages"
        :key="message.id"
        :class="{ 'bg-gray-50 dark:bg-neutral-700/50': !message.isRead }"
        class="p-3 hover:bg-gray-100 dark:hover:bg-neutral-700/80 cursor-pointer transition-colors duration-200"
        @click="handleMessageClick(message)"
      >
        <div class="flex items-start">
          <div
            class="flex-shrink-0 w-12 h-12 rounded bg-gray-200 dark:bg-neutral-600 flex items-center justify-center mr-3"
          >
            <img
              v-if="message.fromUserAvatar"
              :src="message.fromUserAvatar"
              alt="User avatar"
              class="w-full h-full rounded-full object-cover"
            />
            <span
              v-else
              class="text-lg font-medium text-gray-600 dark:text-gray-300"
            >
              {{ message.fromName.charAt(0) }}
            </span>
          </div>
          <div class="flex-grow mr-3">
            <div class="flex flex-wrap items-center justify-between">
              <span
                class="text-sm font-medium text-gray-800 dark:text-gray-200"
              >
                {{ message.fromName }}
              </span>
              <span class="text-xs text-gray-500 dark:text-gray-400">
                {{ sysConfig.timeFormat === "timeAgo" ? $dayjs(message.createdAt).fromNow() : $dayjs(message.createdAt).format("YYYY-MM-DD HH:mm") }}
              </span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-300 mt-1 line-clamp-1">
              <template v-if="message.replyTo">
                <span class="mr-1">回复</span>
                <span class="text-[#576b95] text-nowrap">{{ message.replyTo }}</span>
                <span class="mr-1">:</span>
              </template>
              {{ message.content }}
            </p>
          </div>
          <div class="flex-shrink-0 w-12 h-12 overflow-hidden bg-gray-100 dark:bg-neutral-700 flex items-center justify-center">
            <img
              v-if="memoImages[message.memoId] && memoImages[message.memoId].length > 0"
              :src="memoImages[message.memoId][0]"
              alt="Message image"
              class="w-full h-full object-cover"
            />
            <span v-else-if="memoContents[message.memoId]" class="text-[10px] text-gray-500 dark:text-gray-400 line-clamp-3 px-1" v-html="renderMarkdown(memoContents[message.memoId])">
            </span>
            <span v-else class="text-xs text-gray-500 dark:text-gray-400 p-1 text-center">
              无内容
            </span>
          </div>
          <div
            v-if="!message.isRead"
            class="flex-shrink-0 w-2 h-2 rounded-full bg-blue-500 ml-2 mt-2"
          ></div>
        </div>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";
import type { SysConfigVO } from "~/types";
import { messageChangedEvent } from "~/event";
import { md } from "~/utils";

const global = useGlobalState();
const sysConfig = useState<SysConfigVO>("sysConfig");
const showMessageBox = ref(false);
const messages = ref<any[]>([]);
const unreadCount = ref(0);
const fetching = ref(false);
const memoImages = ref<Record<number, string[]>>({});
const memoContents = ref<Record<number, string>>({});
const fetchingMemoImages = ref<Record<number, boolean>>({});
const showMoreClicked = ref(false);

const toggleMessageBox = () => {
  showMessageBox.value = !showMessageBox.value;
  if (showMessageBox.value) {
    fetchUnreadMessages();
  }
};

// 获取未读消息
const fetchUnreadMessages = async () => {
  if (fetching.value || !global.value.userinfo.token) return;

  fetching.value = true;
  try {
    const response = await fetch("/api/message/unread", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        messages.value = data.data.list;
        unreadCount.value = data.data.total;
        // 获取相关动态的图片
        fetchMemoImages();
      }
    }
  } catch (error) {
    console.error("获取未读消息失败:", error);
  } finally {
    fetching.value = false;
  }
};

// 渲染markdown内容
const renderMarkdown = (content: string) => {
  if (content.length > 20 && !showMoreClicked.value) {
    const truncated = content.substring(0, 20) + '...';
    return md.render(truncated);
  }
  return md.render(content);
}

// 获取动态图片
const fetchMemoImages = () => {
  if (!messages.value || messages.value.length === 0) return;

  messages.value.forEach(async (message) => {
    if (message.memoId && !fetchingMemoImages.value[message.memoId] && !memoImages.value[message.memoId]) {
      fetchingMemoImages.value[message.memoId] = true;
      try {
        // 获取单个动态数据，将id作为查询参数传递
        const memo = await useMyFetch<any>(`/memo/get?id=${message.memoId}`);

        if (memo) {
          // 存储memo内容
          memoContents.value[message.memoId] = memo.content || '';
          
          // 解析图片字符串为数组
          if (memo.imgs) {
            const imgs = memo.imgs.split(',').filter((img: string) => img.trim() !== '');
            memoImages.value[message.memoId] = imgs;
          }
        }
      } catch (error) {
        console.error(`获取动态 ${message.memoId} 图片失败:`, error);
      } finally {
        fetchingMemoImages.value[message.memoId] = false;
      }
    }
  });
};

// 标记消息为已读
const markAsRead = async (messageId: number) => {
  if (!global.value.userinfo.token) return;

  try {
    const response = await fetch(`/api/message/read?id=${messageId}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        if (messages.value) {
          const index = messages.value.findIndex((msg) => msg.id === messageId);
          if (index !== -1) {
            messages.value[index].isRead = true;
            unreadCount.value--;
            // 触发消息数量变化事件
            messageChangedEvent.emit(-1);
          }
        }
      }
    }
  } catch (error) {
    console.error("标记消息为已读失败:", error);
  }
};

// 标记所有消息为已读
const markAllAsRead = async () => {
  if (!global.value.userinfo.token || unreadCount.value === 0) return;

  try {
    const response = await fetch("/api/message/read-all", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        // 更新本地消息状态
        if (messages.value) {
          messages.value.forEach((msg) => {
            msg.isRead = true;
          });
          const change = -unreadCount.value;
          unreadCount.value = 0;
          // 触发消息数量变化事件
          messageChangedEvent.emit(change);
        }
      }
    }
  } catch (error) {
    console.error("标记所有消息为已读失败:", error);
  }
};

// 处理消息点击
const handleMessageClick = (message: any) => {
  if (!message.isRead) {
    markAsRead(message.id);
  }

  // 根据消息类型跳转到相应页面
  if (message.type === "comment" || message.type === "like") {
    navigateTo(`/memo/${message.memoId}`);
  }
};

onMounted(() => {
  // 初始化获取一次未读消息数量
  if (global.value.userinfo.token) {
    fetchUnreadMessages();
  }

  // 监听消息数量变化事件
  const unsubscribe = messageChangedEvent.on((change) => {
    // 根据变化量更新未读数量
    unreadCount.value = Math.max(0, unreadCount.value + change);
  });

  // 存储取消订阅函数，以便在组件卸载时调用
  onUnmounted(() => {
    unsubscribe();
  });
});

onUnmounted(() => {
  // 组件卸载时的清理工作
});
</script>

<style scoped></style>
