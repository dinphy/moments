<template>
  <div 
    :class="[
      'flex flex-col',
      'fixed inset-x-0 bottom-0 bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 shadow-lg z-50 md:static md:border-0 md:shadow-none',
      replyTo ? 'p-3' : 'p-3'
    ]"
    v-if="currentCommentBox === pid"
    ref="commentBoxRef"
  >
    <!-- 回复提示 -->
    <div v-if="replyTo" class="text-xs text-gray-500 mb-2 px-1">
      回复 <span class="text-blue-500">{{ replyTo }}</span>
    </div>
    
    <!-- 输入区域 -->
    <div class="flex items-end gap-2">
      <div class="flex-1 relative">
        <UTextarea 
          :rows="1" 
          :maxrows="4"
          autoresize
          autofocus 
          :placeholder="'评论'" 
          v-model="state.content"
          class="w-full text-sm"
          :ui="{
            base: 'transition-all duration-200 px-4 py-2',
            rounded: 'rounded-full',
            placeholder: 'placeholder:text-gray-400',
            background: 'bg-gray-100 dark:bg-gray-700'
          }"
        />
      </div>
      
      <!-- 功能按钮和发送按钮 -->
      <div class="flex items-center gap-2">
        <UIcon v-if="!global.userinfo.token" class="text-gray-500 w-6 h-6 cursor-pointer" :name="userShow ? 'weui-keyboard-outlined' : 'i-ep-user'" @click="toggleUser"/>
        <UIcon class="text-gray-500 w-6 h-6 cursor-pointer select-none" :name="emojiShow ? 'weui-keyboard-outlined' : 'i-weui-sticker-outlined'" @click="toggleEmoji"/>
        <UButton 
          class="cursor-pointer text-sm px-3" 
          :color="state.content.trim() ? 'primary' : 'gray'" 
          :variant="state.content.trim() ? 'solid' : 'ghost'"
          :disabled="!state.content.trim()"
          size="sm"
          @click="comment">
          发送
        </UButton>
      </div>
    </div>
    
    <!-- 表情选择器 -->
    <div v-if="emojiShow" class="mt-2 border-t border-gray-100 dark:border-gray-700 pt-2">
      <Emoji @selected="emojiSelected"/>
    </div>
    
    <!-- 用户信息输入 -->
    <div v-if="userShow" class="mt-3 grid grid-cols-1 gap-2 md:grid-cols-3">
      <template v-if="!global.userinfo.token">
        <UInput placeholder="姓名" v-model="state.username" size="sm"/>
        <UInput placeholder="网站" v-model="state.website" size="sm"/>
        <UInput placeholder="邮箱" v-model="state.email" size="sm"/>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import {toast} from "vue-sonner";
import {memoChangedEvent, messageChangedEvent} from "~/event";
import Emoji from "~/components/Emoji.vue";
import {useGlobalState} from "~/store";
import {useStorage} from '@vueuse/core'
import type {SysConfigVO} from "~/types";
import { getGuestId } from "~/utils";

const props = defineProps<{
  commentId: number
  memoId: number
  memoUserId: number
  replyTo?: string
  replyEmail?: string
}>()
const pid = computed(() => {
  return `${props.memoId}#${props.commentId}`
})
const global = useGlobalState()
const localCommentUserinfo = useStorage('localCommentUserinfo', {
  username: "",
  website: "",
  email: "",
})
const userShow = ref(false)
const emojiShow = ref(false)
const currentCommentBox = useState('currentCommentBox')
const sysConfig = useState<SysConfigVO>('sysConfig')
const commentBoxRef = ref<HTMLElement>()

onClickOutside(commentBoxRef, () => {
  if (currentCommentBox.value === pid.value) {
    currentCommentBox.value = ''
  }
})

const state = reactive({
  content: "",
  memoId: props.memoId,
  replyTo: props.replyTo,
  replyEmail: props.replyEmail,
  username: localCommentUserinfo.value.username,
  website: localCommentUserinfo.value.website,
  email: localCommentUserinfo.value.email,
})

const comment = async () => {
  if (!state.content.trim()) {
    toast.warning("发送失败，内容不能为空")
    return
  }

  if (sysConfig.value.enableGoogleRecaptcha) {
    grecaptcha.ready(() => {
      grecaptcha.execute(sysConfig.value.googleSiteKey, {action: 'newComment'}).then(async (token) => {
        await doComment(token)
      })
    })
  } else {
    await doComment()
  }
}

const doComment = async (token?: string) => {
  if (!global.value.userinfo.token) {
    localCommentUserinfo.value = {
      username: state.username,
      website: state.website,
      email: state.email,
    }
  }

  if (state.content.length > sysConfig.value.maxCommentLength) {
    toast.error("评论字数超过限制长度:" + sysConfig.value.maxCommentLength)
    return
  }

  const guestId = await getGuestId()
  if (!guestId) return
  await useMyFetch(`/comment/add`, {...state, token, guestId: guestId})
  toast.success("评论成功!")
  currentCommentBox.value = ''
  state.content = ''
  memoChangedEvent.emit(props.memoId)
  // 只有当评论用户不是动态发布者时才更新消息计数
  if (global.value.userinfo.id !== props.memoUserId) {
    messageChangedEvent.emit(1)
  }
}

const toggleUser = () => {
  userShow.value = !userShow.value
}
const toggleEmoji = () => {
  emojiShow.value = !emojiShow.value
}
const emojiSelected = (emoji: string) => {
  state.content = state.content + emoji
}
</script>

<style scoped>

</style>