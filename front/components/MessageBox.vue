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
          'text-[#9fc84a] cursor-pointer': unreadCount > 0,
          'text-gray-400 cursor-not-allowed': unreadCount === 0,
        }"
        class="flex items-center text-sm"
      >
        <UIcon name="i-carbon-checkmark-outline" class="w-4 h-4 mr-1" />全部已读
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
            class="flex-shrink-0 w-8 h-8 rounded-full bg-gray-200 dark:bg-neutral-600 flex items-center justify-center mr-3"
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
          <div class="flex-grow">
            <div class="flex justify-between items-center">
              <span
                class="text-sm font-medium text-gray-800 dark:text-gray-200"
              >
                {{ message.fromName }}
              </span>
              <span class="text-xs text-gray-500 dark:text-gray-400">
                {{
                  sysConfig.timeFormat === "timeAgo"
                    ? $dayjs(message.createdAt).fromNow()
                    : $dayjs(message.createdAt).format("YYYY-MM-DD HH:mm")
                }}
              </span>
            </div>
            <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">
              {{
                message.type === "comment" && message.content.length > 20
                  ? message.content.substring(0, 20) + "..."
                  : message.content
              }}
            </p>
          </div>
          <div
            v-if="!message.isRead"
            class="flex-shrink-0 w-2 h-2 rounded-full bg-blue-500 ml-2"
          ></div>
        </div>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useGlobalState } from "~/store";
import type { SysConfigVO } from "~/types";

const global = useGlobalState();
const sysConfig = useState<SysConfigVO>("sysConfig");
const showMessageBox = ref(false);
const messages = ref<any[]>([]);
const unreadCount = ref(0);
const fetching = ref(false);

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
      }
    }
  } catch (error) {
    console.error("获取未读消息失败:", error);
  } finally {
    fetching.value = false;
  }
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
          unreadCount.value = 0;
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
  // 可以在这里初始化获取一次未读消息数量
  if (global.value.userinfo.token) {
    fetchUnreadMessages();
  }
});

onUnmounted(() => {
  // 可以在这里清理定时器等
});
</script>

<style scoped></style>
