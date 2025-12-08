<template>
  <div class="fixed inset-0 bg-gradient-to-br from-gray-900/80 to-black/90 flex items-center justify-center z-50" @click="closeDialog">
    <div class="bg-white dark:bg-gray-900 rounded-lg w-full max-w-lg overflow-hidden flex flex-col shadow-2xl mx-2" @click.stop>
      <!-- 标题栏 -->
      <div class="p-4 border-b border-gray-200 dark:border-gray-800 flex items-center justify-between bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-gray-800 dark:to-gray-900">
        <div class="flex items-center space-x-3">
          <div class="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex items-center justify-center">
            <UIcon name="i-carbon-bot" class="w-6 h-6 text-white" />
          </div>
          <h3 class="text-md text-gray-900 dark:text-white">AI 润色助手</h3>
        </div>
        <div class="flex items-center space-x-5">
          <div class="group relative">
            <UIcon name="i-weui-delete-outlined" class="w-5 h-5 cursor-pointer text-gray-600 dark:text-gray-300 hover:text-red-500 dark:hover:text-red-400 transition-all duration-200 transform hover:scale-110 px-3" @click="clearHistory" />
            <span class="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs rounded px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">清空历史</span>
          </div>
          <div class="group relative">
            <UIcon name="i-weui-done-filled" 
              class="w-5 h-5 cursor-pointer transition-all duration-200 transform hover:scale-110 px-3" 
              :class="!lastAIResponse ? 'text-gray-400 dark:text-gray-600 cursor-not-allowed' : 'text-green-600 dark:text-green-400 hover:text-green-700 dark:hover:text-green-300'" 
              @click="applyToContent" 
              :disabled="!lastAIResponse" />
            <span class="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs rounded px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">应用内容</span>
          </div>
          <div class="group relative">
            <UIcon name="i-weui-close-outlined" 
              class="w-5 h-5 cursor-pointer text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-all duration-200 transform hover:scale-110 hover:rotate-90 px-3" 
              @click="closeDialog" />
            <span class="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs rounded px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">关闭</span>
          </div>
        </div>
      </div>

      <!-- 对话内容区 -->
      <div ref="chatContainer" class="flex-1 min-h-[25vh] max-h-[50vh] overflow-y-auto p-4 space-y-6 bg-gray-50 dark:bg-gray-900/50 scrollbar">
        <div v-for="(message, index) in messages" :key="index" class="flex" :class="message.role === 'user' ? 'justify-end' : 'justify-start'">
          <div v-if="message.role === 'assistant'" class="flex items-start space-x-3 max-w-[85%]">
            <div class="w-5 h-5 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex-shrink-0 flex items-center justify-center">
              <UIcon name="i-carbon-bot" class="w-3 h-3 text-white" />
            </div>
            <div class="bg-white dark:bg-gray-800 rounded-lg rounded-tl-none px-4 py-2 shadow-sm border border-gray-200 dark:border-gray-700">
              <p class="text-sm text-gray-800 dark:text-gray-200 whitespace-pre-wrap leading-relaxed">{{ message.content }}</p>
            </div>
          </div>
          
          <!-- 用户消息 -->
          <div v-else class="flex items-start space-x-3 max-w-[85%] justify-end">
            <div class="bg-gradient-to-r from-blue-500 to-indigo-600 rounded-lg rounded-tr-none px-4 py-2 shadow-sm text-white">
              <p class="text-sm whitespace-pre-wrap leading-relaxed">{{ message.content }}</p>
            </div>
            <div class="w-5 h-5 rounded-full bg-gradient-to-br from-gray-400 to-gray-600 flex-shrink-0 flex items-center justify-center">
              <UIcon name="i-carbon-user" class="w-3 h-3 text-white" />
            </div>
          </div>
        </div>

        <div v-if="isLoading" class="flex items-start space-x-3 max-w-[85%]">
          <div class="w-6 h-6 rounded-full bg-gradient-to-br from-blue-500 to-indigo-600 flex-shrink-0 flex items-center justify-center">
            <UIcon name="i-carbon-bot" class="w-5 h-5 text-white" />
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg rounded-tl-none px-4 py-3 shadow-sm border border-gray-200 dark:border-gray-700">
            <div class="flex space-x-2">
              <div class="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 0ms"></div>
              <div class="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 150ms"></div>
              <div class="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style="animation-delay: 300ms"></div>
            </div>
          </div>
        </div>
      </div>

      <div class="border-t border-gray-200 dark:border-gray-800 p-4 bg-white dark:bg-gray-900">
        <!-- 预设指令 -->
        <div v-if="props.content" class="flex items-center flex-wrap gap-2 mb-4">
          <UBadge 
            v-for="(prompt, index) in presetPrompts" 
            :key="index"
            @click="applyPresetPromptAndSend(prompt)"
            class="cursor-pointer bg-gray-100 hover:bg-blue-100 dark:bg-gray-800 dark:hover:bg-blue-900/50 text-gray-700 dark:text-gray-300 transition-colors px-3 py-1 text-sm"
          >
            {{ prompt.title }}
          </UBadge>
        </div>

        <div class="flex flex-col space-y-3">
          <div class="flex space-x-3 relative">
            <UTextarea 
              v-model="userInput" 
              placeholder="输入您的指令或内容..." 
              :rows="3"
              autofocus
              @keydown.enter.prevent="sendMessage"
              :disabled="isLoading"
              class="flex-1 rounded-xl border-gray-300 dark:border-gray-700 focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
            />
            <div class="flex space-x-2 absolute right-2 bottom-2">
              <UButton 
                @click="sendMessage" 
                :disabled="!userInput.trim() || isLoading" 
                class="p-2 rounded-full text-white transition-all"
                size="sm"
              >
                <UIcon name="i-carbon-location-current" class="w-4 h-4" />
              </UButton>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toast } from "vue-sonner";

const props = defineProps<{
  content: string
}>();

const emit = defineEmits(['close', 'applyContent']);

// 对话消息
const messages = ref<Array<{role: string, content: string}>>([]);
const userInput = ref('');
const isLoading = ref(false);
const lastAIResponse = ref('');
const chatContainer = ref<HTMLElement | null>(null);

// 滚动到最新消息
const scrollToBottom = () => {
  nextTick(() => {
    if (chatContainer.value) {
      chatContainer.value.scrollTop = chatContainer.value.scrollHeight;
    }
  });
};

// 预设指令
const presetPrompts = ref([
  { title: "润色", prompt: "请将内容斧正、润色，保持原意不变。" },
  { title: "简化", prompt: "请将内容提炼，用更简洁的语言表达。" },
  { title: "扩展", prompt: "请对内容扩展，丰富细节和背景信息。" },
  { title: "续写", prompt: "请基于内容进行逻辑连贯的续写，保持风格一致。" },
]);

// 初始化对话，添加原始内容
onMounted(() => {
  if (props.content) {
    messages.value.push({
      role: 'user',
      content: `${props.content}`
    });
    // 自动发送第一条消息
    nextTick(() => {
      sendMessage();
    });
  } else {
    // 添加欢迎消息
    messages.value.push({
      role: 'assistant',
      content: '嗨，您好！有需要尽管吩咐哦~'
    });

    // 滚动到最新消息
    scrollToBottom();
  }
});

// 发送消息
const sendMessage = async () => {
  const content = userInput.value.trim();
  if (!content || isLoading.value) return;

  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: content
  });

  userInput.value = '';
  isLoading.value = true;

  try {
    // 构建请求
    const requestMessages = messages.value.map(msg => ({
      role: msg.role,
      content: msg.content
    }));

    // 调用AI对话API
    const response = await useMyFetch<any>('/ai/chat', {
      messages: requestMessages
    });

    if (response && response.message && response.message.content) {
      // 添加AI回复
      messages.value.push({
        role: 'assistant',
        content: response.message.content
      });

      // 保存最后一条AI回复，用于应用到内容
      lastAIResponse.value = response.message.content;

      // 滚动到最新消息
      scrollToBottom();
    } else {
      toast.error("AI回复失败");
    }
  } catch (error) {
    console.error("AI对话错误:", error);
    toast.error("AI对话失败，请稍后重试");
  } finally {
    isLoading.value = false;
  }
};

// 应用预设指令
const applyPresetPrompt = (preset: any) => {
  userInput.value = preset.prompt;
};

// 应用预设指令并发送
const applyPresetPromptAndSend = (preset: any) => {
  userInput.value = preset.prompt;
  sendMessage();
};

// 清空对话历史
const clearHistory = () => {
  messages.value = [];
  lastAIResponse.value = '';

  // 添加欢迎消息
  messages.value.push({
    role: 'assistant',
    content: '已清空，开始新的对话吧~'
  });

  // 滚动到最新消息
  scrollToBottom();
};

// 应用AI回复到内容
const applyToContent = () => {
  if (lastAIResponse.value) {
    emit('applyContent', lastAIResponse.value);
    closeDialog();
  }
};

// 关闭对话框
const closeDialog = () => {
  emit('close');
};
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