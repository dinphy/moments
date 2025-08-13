<template>
  <div 
    :class="[
      'flex flex-col gap-2',
      'fixed inset-x-0 bottom-0 bg-white dark:bg-[#202020] border-t border-gray-200 dark:border-[#303030] shadow-lg z-50 max-h-[50vh] overflow-y-auto overscroll-contain sm:static sm:mt-2 sm:max-h-none sm:border-0 sm:shadow-none sm:bg-transparent',
      replyTo ? 'p-3 sm:p-1' : 'p-3'
    ]"
    v-if="currentCommentBox === pid"
    ref="commentBoxRef"
  >
    <div :class="[
      'relative flex items-start gap-2',
      'py-2 sm:bg-transparent'
    ]">
      <UTextarea 
        :rows="1" 
        :maxrows="3"
        autoresize
        autofocus 
        :placeholder="replyTo ? `回复${replyTo}:` : '说点什么...'" 
        v-model="state.content"
        class="flex-1 min-h-[32px] text-sm"
        :ui="{
          base: 'transition-all duration-200',
          rounded: 'rounded-lg',
          placeholder: 'placeholder:text-gray-400'
        }"
      />
      <div class="flex gap-1 items-center flex-shrink-0">
        <UIcon v-if="!global.userinfo.token" class="text-gray-400 w-7 h-7 cursor-pointer" name="i-carbon-user-avatar" @click="toggleUser"/>
        <UIcon class="text-gray-400 w-7 h-7 cursor-pointer select-none" name="i-carbon-face-activated" @click="toggleEmoji"/>
        <UButton 
          class="cursor-pointer text-sm" 
          :color="state.content.trim() ? 'primary' : 'white'" 
          :disabled="!state.content.trim()"
          @click="comment">
          发送
        </UButton>
      </div>
    </div>
    <Emoji v-if="emojiShow" @selected="emojiSelected"/>
    <div v-if="userShow" :class="[
      'flex gap-1 flex-col space-y-2 sm:flex-row sm:space-y-0 sm:space-x-2'
    ]">
      <template v-if="!global.userinfo.token">
        <UInput placeholder="姓名" v-model="state.username"/>
        <UInput placeholder="网站" v-model="state.website"/>
        <UInput placeholder="邮箱" v-model="state.email"/>
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