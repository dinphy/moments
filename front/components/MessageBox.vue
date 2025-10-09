<template>
  <div class="absolute top-4 left-4 z-10">
    <UPopover v-model:open="showMessageBox" :popper="{ placement: 'bottom', offset: 8, strategy: 'fixed', arrow: true }">
      <div class="cursor-pointer relative">
        <UIcon name="i-carbon-notification" class="text-[#9fc84a] w-5 h-5" />
        <span
          v-if="unreadCount > 0"
          class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center"
        >
          {{ unreadCount }}
        </span>
      </div>

      <template #panel>
        <div class="w-[360px] max-w-[90vw] rounded-lg bg-white dark:bg-neutral-800 shadow-lg pb-2">
          <div class="flex justify-between items-center p-3 relative">
            <h3 class="flex flex-1 justify-center font-medium pb-2 border-b border-gray-200 dark:border-gray-700">消息</h3>
            <div
              v-if="messages.length > 0"
              @click="handleDeleteAllMessages()"
              class="absolute top-3 right-3 text-sm text-red-500 cursor-pointer"
            >
              <span>清空</span>
            </div>
          </div>

          <div
            v-if="!messages || messages.length === 0"
            class="py-10 text-center text-gray-500 dark:text-gray-400"
          >
            暂无消息
          </div>

          <div v-else class="max-h-96 overflow-y-auto scrollbar">
            <!-- 未读消息 -->
            <div v-for="message in messages" :key="message.id">
              <div v-if="!message.isRead">
                <div
                  class="p-3 hover:bg-gray-100 dark:hover:bg-neutral-700/80 cursor-pointer transition-colors duration-200"
                  @click="handleMessageClick(message)"
                >
                  <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 rounded bg-gray-200 dark:bg-neutral-600 flex items-center justify-center mr-3 relative">
                      <img
                        v-if="message.fromUserAvatar"
                        :src="message.fromUserAvatar"
                        alt="User avatar"
                        class="w-full h-full rounded-full object-cover"
                      />
                      <span v-else class="text-lg font-medium text-gray-600 dark:text-gray-300">
                        {{ message.fromName.charAt(0) }}
                      </span>
                      <div v-if="!message.isRead" class="absolute -top-1 -right-1 w-2 h-2 rounded-full bg-blue-500"></div>
                    </div>
                    
                    <div class="flex-grow mr-3">
                      <div class="flex flex-wrap items-center justify-between">
                        <span class="text-sm font-medium text-gray-800 dark:text-gray-200">
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
                        <template v-else-if="message.type === 'like'">
                          <div class="flex items-center">
                            <UIcon name="i-carbon-favorite" class="text-red-500 w-4 h-4 inline-block mr-1" />
                            <span>#{{ message.memoId }}</span>
                          </div>
                        </template>
                        <template v-else>
                          {{ message.content }}
                        </template>
                      </p>
                    </div>
                    <div class="flex-shrink-0 w-12 h-12 overflow-hidden bg-gray-100 dark:bg-neutral-700 flex items-center justify-center relative group/preview">
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
                      <button
                        @click.stop="deleteMessage(message.id)"
                        class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover/preview:opacity-100 transition-opacity duration-200"
                      >
                        <UIcon name="i-carbon-trash-can" class="w-5 h-5 text-white" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 已读消息分割线 -->
            <div v-if="messages.some(msg => msg.isRead) && unreadCount > 0" class="px-4 py-6 text-center text-xs text-gray-400 dark:text-gray-500">
              <span class="inline-block w-16 border-t border-gray-300 dark:border-gray-600 align-middle mx-2"></span>
              以下为已读消息
              <span class="inline-block w-16 border-t border-gray-300 dark:border-gray-600 align-middle mx-2"></span>
            </div>

            <!-- 已读消息 -->
            <div v-for="message in messages" :key="message.id">
              <div v-if="message.isRead">
                <div
                  class="p-3 hover:bg-gray-100 dark:hover:bg-neutral-700/80 cursor-pointer transition-colors duration-200"
                  @click="handleMessageClick(message)"
                >
                  <div class="flex items-start">
                    <div class="flex-shrink-0 w-12 h-12 rounded bg-gray-200 dark:bg-neutral-600 flex items-center justify-center mr-3 relative">
                      <img
                        v-if="message.fromUserAvatar"
                        :src="message.fromUserAvatar"
                        alt="User avatar"
                        class="w-full h-full rounded-full object-cover"
                      />
                      <span v-else class="text-lg font-medium text-gray-600 dark:text-gray-300">
                        {{ message.fromName.charAt(0) }}
                      </span>
                    </div>
                    <div class="flex-grow mr-3">
                      <div class="flex flex-wrap items-center justify-between">
                        <span class="text-sm font-medium text-gray-800 dark:text-gray-200">
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
                        <template v-else-if="message.type === 'like'">
                          <div class="flex items-center">
                            <UIcon name="i-carbon-favorite" class="text-red-500 w-4 h-4 inline-block mr-1" />
                            <span>了动态 #{{ message.memoId }}</span>
                          </div>
                        </template>
                        <template v-else>
                          {{ message.content }}
                        </template>
                      </p>
                    </div>
                    <div class="flex-shrink-0 w-12 h-12 overflow-hidden bg-gray-100 dark:bg-neutral-700 flex items-center justify-center relative group/preview">
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
                      <button
                        @click.stop="deleteMessage(message.id)"
                        class="absolute inset-0 flex items-center justify-center bg-black/50 opacity-0 group-hover/preview:opacity-100 transition-opacity duration-200"
                      >
                        <UIcon name="i-carbon-trash-can" class="w-5 h-5 text-white" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </UPopover>
  </div>

  <UModal v-model="showDeleteConfirm" :ui="{
    container: 'flex justify-center items-center backdrop-blur-sm',
  }">
    <div class="p-5">
      <h3 class="font-medium text-lg mb-4">清空提示</h3>
      <p class="text-gray-600 dark:text-gray-300 mb-6">
        确定要清空所有消息吗？此操作不可恢复。
      </p>
      <div class="flex justify-end space-x-3">
        <UButton
          color="gray"
          variant="solid"
          @click="showDeleteConfirm = false"
        >
          取消
        </UButton>
        <UButton
          color="primary"
          variant="solid"
          @click="confirmDeleteAllMessages"
        >
          确认
        </UButton>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";
import type { SysConfigVO } from "~/types";
import { md } from "~/utils";
import { messageChangedEvent } from "~/event";

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

// 监听showMessageBox的变化，当打开时刷新消息
watch(showMessageBox, (newValue) => {
  if (newValue) {
    fetchAllMessages();
  }
});

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
        if (messages.value && messages.value.length > 0) {
          const index = messages.value.findIndex((msg) => msg.id === messageId);
          if (index !== -1) {
            messages.value[index].isRead = true;
            unreadCount.value--;
            messageChangedEvent.emit(-1);
          }
        }
      }
    }
  } catch (error) {
    console.error("标记消息为已读失败:", error);
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

const fetchAllMessages = async () => {
  if (fetching.value || !global.value.userinfo.token) return;

  fetching.value = true;
  try {
    const response = await fetch("/api/message/all", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        messages.value = data.data.list || [];
        unreadCount.value = messages.value.filter((msg: any) => !msg.isRead).length;
        fetchMemoImages();
      } else {
        messages.value = [];
      }
    } else {
      messages.value = [];
    }
  } catch (error) {
    console.error("获取所有消息失败:", error);
    messages.value = [];
  } finally {
    fetching.value = false;
  }
};

const deleteMessage = async (messageId: number) => {
  if (!global.value.userinfo.token) return;

  try {
    const response = await fetch(`/api/message/delete?id=${messageId}`, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        messages.value = messages.value.filter(msg => msg.id !== messageId);
        unreadCount.value = messages.value.filter(msg => !msg.isRead).length;
        messageChangedEvent.emit(-1);
      }
    }
  } catch (error) {
    console.error("删除消息失败:", error);
  }
};

const showDeleteConfirm = ref(false);

const handleDeleteAllMessages = () => {
  if (messages.value.length === 0) return;
  showDeleteConfirm.value = true;
};

// 确认删除所有消息
const confirmDeleteAllMessages = async () => {
  showDeleteConfirm.value = false;
  if (!global.value.userinfo.token || messages.value.length === 0) return;

  try {
    const response = await fetch("/api/message/delete-all", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "x-api-token": global.value.userinfo.token,
      },
    });

    if (response.ok) {
      const data = await response.json();
      if (data.code === 0) {
        messages.value = [];
        unreadCount.value = 0;
        messageChangedEvent.emit(-1);
      }
    }
  } catch (error) {
    console.error("删除所有消息失败:", error);
  }
};

onMounted(() => {
  if (global.value.userinfo.token) {
    fetchAllMessages();
  }

  // 监听消息数量变化事件
  const unsubscribe = messageChangedEvent.on((change) => {
    // 如果是增加消息，直接增加未读计数
    if (change > 0) {
      unreadCount.value += change;
    } else {
      // 如果是减少消息或刷新，重新获取所有消息
      fetchAllMessages();
    }
  });

  // 存储取消订阅函数，以便在组件卸载时调用
  onUnmounted(() => {
    unsubscribe();
  });
});
</script>

<style scoped>
/* 自定义滚动条样式 */
.scrollbar::-webkit-scrollbar {
  width: 6px;
}

.scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(156, 163, 175, 0.5);
  border-radius: 3px;
}

.scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(156, 163, 175, 0.8);
}

/* 暗色模式下的滚动条样式 */
.dark .scrollbar::-webkit-scrollbar-thumb {
  background-color: rgba(75, 85, 99, 0.5);
}

.dark .scrollbar::-webkit-scrollbar-thumb:hover {
  background-color: rgba(75, 85, 99, 0.8);
}
</style>
