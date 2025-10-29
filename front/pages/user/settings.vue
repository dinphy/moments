<template>
  <Header :user="currentUser" v-if="!isAdminMode"/>

  <div class="bg-gray-100 dark:bg-gray-900 min-h-screen p-2 rounded-b-lg">
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showAvatar = !showAvatar">
          <span class="text-gray-700 dark:text-gray-300">头像</span>
          <div class="flex items-center space-x-2">
            <UAvatar :src="state.avatarUrl" size="sm"/>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showAvatar" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600 space-y-3">
          <UInput v-model="state.avatarUrl" placeholder="输入头像地址或上传" size="md"/>
          <label class="cursor-pointer inline-block">
            <UInput
              type="file"
              @change="uploadAvatarUrl"
              accept="image/*"
              class="hidden"
            />
            <div class="flex items-center justify-center p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors text-sm">
              <UIcon name="i-heroicons-arrow-up-tray" class="w-4 h-4 mr-1"/>
              上传头像
            </div>
          </label>
        </div>
        
        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showNickname = !showNickname">
          <span class="text-gray-700 dark:text-gray-300">昵称</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.nickname || "未设置" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showNickname" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="state.nickname" placeholder="请输入昵称" size="md"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showSlogan = !showSlogan">
          <span class="text-gray-700 dark:text-gray-300">个性签名</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.slogan || "这个人很懒，什么都没留下" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showSlogan" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
          <UInput v-model="state.slogan" placeholder="请输入个性签名" size="md"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showCover = !showCover">
          <span class="text-gray-700 dark:text-gray-300">封面</span>
          <div class="flex items-center space-x-2">
            <span v-if="state.coverUrl" class="text-gray-500 dark:text-gray-400 text-sm">已设置</span>
            <span v-else class="text-gray-500 dark:text-gray-400 text-sm">未设置</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showCover" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600 space-y-3">
          <UInput v-model="state.coverUrl" placeholder="输入封面地址或上传" size="md"/>
          <label class="cursor-pointer inline-block">
            <UInput
              type="file"
              @change="uploadCoverUrl"
              accept="image/*"
              class="hidden"
            />
            <div class="flex items-center justify-center p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 transition-colors text-sm">
              <UIcon name="i-heroicons-arrow-up-tray" class="w-4 h-4 mr-1"/>
              上传封面
            </div>
          </label>
          <div v-if="state.coverUrl" class="rounded overflow-hidden mt-3">
            <img :src="state.coverUrl" class="w-full h-32 object-cover" alt="封面预览"/>
          </div>
        </div>
      </div>
    </div>

    <!-- 账号设置 -->
    <div class="bg-white dark:bg-gray-800 shadow-sm mb-2 rounded-lg overflow-hidden">
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <div class="px-4 py-3 flex items-center justify-between">
          <span class="text-gray-700 dark:text-gray-300">登录名</span>
          <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.username }}</span>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showPassword = !showPassword">
          <span class="text-gray-700 dark:text-gray-300">修改密码</span>
          <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
        </div>
        <div v-show="showPassword" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600 space-y-3">
          <UInput v-model="state.password" type="password" placeholder="请输入新密码（留空则不修改）" size="md"/>
        </div>

        <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showEmail = !showEmail">
          <span class="text-gray-700 dark:text-gray-300">邮箱</span>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.email || "未设置" }}</span>
            <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400"/>
          </div>
        </div>
        <div v-show="showEmail" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600 space-y-3">
          <UInput v-model="state.email" type="email" placeholder="请输入邮箱地址" size="md"/>
          <p class="text-xs text-gray-500 dark:text-gray-400">若管理员启用了邮件通知，将在收到评论时发送邮件通知</p>
        </div>
      </div>
    </div>

    <!-- 保存按钮 -->
    <div class="py-3">
      <UButton class="w-full justify-center bg-blue-500 hover:bg-blue-600" @click="save" size="md">保存设置</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {UserVO} from "~/types";
import {toast} from "vue-sonner";
import {useUpload} from "~/utils";
import {useGlobalState} from "~/store";

const props = defineProps<{
  targetUser?: UserVO,
  isAdminMode?: boolean,
  onSave?: () => void
}>()

const global = useGlobalState()
const currentUser = useState<UserVO>('userinfo')
const state = reactive({
  password: "",
  username: "",
  nickname: "",
  slogan: "",
  avatarUrl: "",
  coverUrl: "",
  email: "",
  css: "",
  js: "",
})

const showAvatar = ref(false)
const showCover = ref(false)
const showNickname = ref(false)
const showSlogan = ref(false)
const showPassword = ref(false)
const showEmail = ref(false)

const logout = async () => {
  global.value.userinfo = {}
  await navigateTo('/')
}
const reload = async () => {
  const res = await useMyFetch<UserVO>('/user/profile')
  if (res) {
    Object.assign(state, res)
    currentUser.value = res
  }
}

const save = async () => {
  try {
    if (state.password && state.password.length < 6) {
      toast.warning("密码长度至少6位")
      return
    }

    if (props.isAdminMode) {
      // 管理员模式：更新其他用户信息
      await useMyFetch('/user/update', {
        id: props.targetUser!.id,
        username: state.username,
        nickname: state.nickname,
        email: state.email,
        slogan: state.slogan,
        avatarUrl: state.avatarUrl,
        coverUrl: state.coverUrl,
        ...(state.password && { password: state.password })
      })
      toast.success("用户更新成功")
      if (props.onSave) {
        props.onSave()
      }
    } else {
      // 普通用户模式：更新自己的信息
      await useMyFetch('/user/saveProfile', state)
      toast.success("保存成功")
      location.reload()
    }
  } catch (error) {
    toast.error(props.isAdminMode ? "用户更新失败" : "保存失败")
  }
}

const uploadAvatarUrl = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0){
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.avatarUrl = result[0]
  }
}

const uploadCoverUrl = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0){
      toast.error("只能上传图片");
      return
    }
  }
  const result = await useUpload(files)
  if (result.length) {
    toast.success("上传成功")
    state.coverUrl = result[0]
  }
}

onMounted(async () => {
  if (props.isAdminMode && props.targetUser) {
    Object.assign(state, props.targetUser)
    state.password = ""
  } else {
    Object.assign(state, currentUser.value)
  }
})

</script>
<style scoped>

</style>
