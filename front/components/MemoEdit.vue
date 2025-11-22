<template>
  <div class="px-4 space-y-2">
    <div class="flex justify-between items-center pt-4 text-gray-600">
      <div class="flex items-center" title="返回">
        <UIcon @click="goBack" name="i-carbon-chevron-left" class="w-5 h-5 cursor-pointer mr-4"/>
        <span v-if="props.inDrawer || $route.path==='/new'">添加内容</span>
        <span v-else>编辑内容</span>
      </div>
      <UButton @click="saveMemo">发表</UButton>
    </div>
    <div class="flex gap-2 text-lg text-gray-600 pt-4 ">
      <ExternalUrl v-model:favicon="state.externalFavicon" v-model:title="state.externalTitle"
                   v-model:url="state.externalUrl"/>

      <upload-image v-model:imgs="state.imgs"/>
      <music v-bind="state.music" @confirm="updateMusic"/>
      <upload-video @confirm="handleVideo" v-bind="state.video"/>
      <douban-edit v-model:type="doubanType" v-model:data="doubanData"/>
      <UPopover :popper="{ arrow: true }" mode="click">
        <UIcon name="i-carbon-calendar" class="w-6 h-6" title="自定义时间"/>
        <template #panel="{close}">
          <DatePicker
            v-model="state.createdAt"
            mode="datetime"
            is24hr
            :time-accuracy="2"
            :rules="{ seconds: 0 }"
            @close="close"
          />
        </template>
      </UPopover>
      <UIcon name="i-carbon-text-clear-format" @click="reset" class="w-6 h-6 cursor-pointer" title="清空"></UIcon>
    </div>

    <div class="w-full">
      <div class="relative">
        <UTextarea ref="contentRef" v-model="state.content" :rows="8" autoresize padded autofocus placeholder="这一刻的想法..."/>
        <UIcon class="text-[#9fc84a] w-7 h-7 animate-bounce absolute left-2 bottom-2 cursor-pointer select-none" :name="emojiShow ? 'weui-keyboard-outlined' : 'i-weui-sticker-outlined'" @click="toggleEmoji"/>
      </div>

      <Emoji v-if="emojiShow" @selected="emojiSelected" @close="emojiShow=false"/>

      <div class="bg-white dark:bg-gray-800 shadow-sm rounded-lg overflow-hidden mt-3">
        <div class="divide-y divide-gray-100 dark:divide-gray-700">
          <!-- 标签设置 -->
          <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="showTags = !showTags">
            <div class="flex items-center gap-2">
              <UIcon name="i-weui-tag-outlined" class="w-5 h-5 text-gray-500"/>
              <span class="text-gray-700 dark:text-gray-300">标签</span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-gray-500 dark:text-gray-400 text-sm">{{ selectedLabel.length ? `已选择${selectedLabel.length}个` : "未添加" }}</span>
              <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400 transition-transform" :class="{'rotate-90': showTags}"/>
            </div>
          </div>
          <div v-show="showTags" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
            <div class="mb-3">
              <UInput 
                v-model="newTag" 
                placeholder="按 Enter 确认" 
                size="md"
                @keyup.enter="addTag"
              >
              </UInput>
            </div>
            <div v-if="selectedLabel.length" class="mb-3">
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">已选择的标签</p>
              <div class="flex flex-wrap gap-2">
                <UBadge 
                  v-for="(tag,index) in selectedLabel" 
                  :key="index" 
                  size="sm" 
                  color="blue" 
                  variant="soft"
                  class="cursor-pointer hover:bg-blue-200 dark:hover:bg-blue-600 transition-colors flex items-center gap-1"
                >
                  {{ tag }}
                  <UIcon name="i-carbon-close" class="w-3 h-3" @click.stop="removeTag(index)"/>
                </UBadge>
              </div>
            </div>
            <div v-if="existTags.length">
              <p class="text-xs text-gray-500 dark:text-gray-400 mb-2">常用标签</p>
              <div class="flex flex-wrap gap-2">
                <UBadge 
                  v-for="(tag,index) in existTags" 
                  :key="index" 
                  size="sm" 
                  color="gray" 
                  variant="soft"
                  class="cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                  :class="{ 'opacity-50': selectedLabel.includes(tag) }"
                  @click="clickTag(tag)"
                >
                  {{ tag }}
                </UBadge>
              </div>
            </div>
          </div>
        </div>
        <div class="divide-y divide-gray-100 dark:divide-gray-700">
          <!-- 位置设置 -->
          <div class="px-4 py-3 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer" @click="toggleLocationPanel">
            <div class="flex items-center gap-2">
              <UIcon name="i-weui-location-outlined" class="w-5 h-5 text-[#576b95]"/>
              <span class="text-gray-700 dark:text-gray-300">所在位置</span>
            </div>
            <div class="flex items-center space-x-2">
              <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.location ? locationLabel : "未设置" }}</span>
              <UIcon name="i-heroicons-chevron-right" class="w-4 h-4 text-gray-400 transition-transform" :class="{'rotate-90': showLocationPanel}"/>
            </div>
          </div>
          <div v-show="showLocationPanel" class="px-4 py-3 bg-gray-50 dark:bg-gray-700/30 border-t border-gray-100 dark:border-gray-600">
            <UInput 
              v-model="state.location" 
              placeholder="例如：北京 朝阳区 三里屯"
              size="md"
              autofocus
            />
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-2">多个位置用空格分隔，显示时会用"·"连接</p>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white dark:bg-gray-800 shadow-sm rounded-lg overflow-hidden mt-3">
      <div class="divide-y divide-gray-100 dark:divide-gray-700">
        <!-- 可见性设置 -->
        <div class="px-4 py-3 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <UIcon name="i-weui-me-outlined" class="w-5 h-5 text-gray-500"/>
            <span class="text-gray-700 dark:text-gray-300">谁可以看</span>
          </div>
          <div class="flex items-center space-x-2">
            <span class="text-gray-500 dark:text-gray-400 text-sm">{{ state.showType ? '公开' : '私密' }}</span>
            <UToggle v-model="state.showType" color="primary"/>
          </div>
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-2">
      <external-url-preview :favicon="state.externalFavicon" :title="state.externalTitle" :url="state.externalUrl"/>
      <upload-image-preview :imgs="state.imgs" @remove-image="handleRemoveImage" @drag-image="handleDragImage"/>
      <music-preview v-if="state.music && state.music.id && state.music.type && state.music.server"
                     v-bind="state.music"/>
      <douban-book-preview :book="doubanData" v-if="doubanType === 'book' && doubanData&& doubanData.title"/>
      <douban-movie-preview :movie="doubanData" v-if="doubanType === 'movie' && doubanData&& doubanData.title"/>
      <video-preview-iframe v-if="['bilibili', 'youtube'].includes(state.video.type) && state.video.value" :url="state.video.value"/>
      <video-preview v-if="state.video.type === 'online' && state.video.value" :url="state.video.value"/>
    </div>
  </div>
</template>

<script setup lang="ts">
import {useMouse, useWindowScroll} from '@vueuse/core'
import type {
  DoubanBook,
  DoubanMovie,
  ExtDTO,
  MemoVO,
  MetingMusicServer,
  MetingMusicType,
  MusicDTO,
  Video,
  VideoType
} from "~/types";
import {toast} from "vue-sonner";
import UploadImage from "~/components/UploadImage.vue";
import Emoji from "~/components/Emoji.vue";
import dayjs from "dayjs";

const doubanType = ref<'book' | 'movie'>('book')
const doubanData = ref<DoubanBook | DoubanMovie>({})
const contentRef = ref(null)
const props = defineProps<{ id?: number, inDrawer?: boolean }>()
const emit = defineEmits(['success', 'close'])
const defaultState = {
  id: props.id || 0,
  createdAt: '' as string,
  content: "",
  ext: "",
  pinned: false,
  showType: true,
  location: "",
  externalFavicon: "",
  externalTitle: "",
  externalUrl: "",
  imgs: "",
  music: {
    id: '',
    api: 'https://api.i-meto.com/meting/api?server=:server&type=:type&id=:id&r=:r',
    server: 'netease' as MetingMusicServer,
    type: 'song' as MetingMusicType,
    title: '',
    artist: ''
  },
  video: {
    type: 'youtube' as VideoType,
    value: ""
  },
  doubanBook: {} as DoubanBook,
  doubanMovie: {} as DoubanMovie,
  tags: Array<string>(),
}
const selectedTags = ref<Array<string>>([])
const selectedLabel = computed({
  get:()=>selectedTags.value,
  set:(labels:Array<string>)=>{
    const tempLabels = Array<string>()
    labels.map(label=>{
      // @ts-ignore
      if(typeof  label !== 'string'){
        // @ts-ignore
        label = label.label
      }
      tempLabels.push(label)
      if(!existTags.value.includes(label)){
        existTags.value.push(label)
      }
    })
    selectedTags.value = [...tempLabels]
    console.log('selectedTags',selectedTags.value)
  }
})
const state = reactive({
  ...defaultState
})
const existTags = ref<string[]>([])
const reset = () => {
  Object.assign(state, defaultState)
}

const locationLabel = computed(() => {
  return state.location.split(" ").join(" · ")
})

const handleDragImage = (imgs: string[]) => {
  state.imgs = imgs.filter(Boolean).join(",")
}

const updateMusic = (music: MusicDTO) => {
  state.music.id = ""
  setTimeout(() => {
    Object.assign(state.music, music)
  }, 500)
}

const handleVideo = (video: Video) => {
  state.video = video
}

const handleRemoveImage = (img: string) => {
  state.imgs = state.imgs
    .split(",")
    .filter(item => item && item != img)
    .join(",")
}

const loadTags = async () => {
  const res = await useMyFetch<{
    tags: string[]
  }>("/tag/list")
  existTags.value = res.tags || []
}

const emojiShow = ref(false)
const showLocationPanel = ref(false)
const showTags = ref(false)
const newTag = ref('')

const toggleEmoji = () => {
  emojiShow.value = !emojiShow.value
}

const toggleLocationPanel = () => {
  showLocationPanel.value = !showLocationPanel.value
}

const addTag = () => {
  const tag = newTag.value.trim()
  if (tag && !selectedLabel.value.includes(tag)) {
    selectedLabel.value.push(tag)
    if (!existTags.value.includes(tag)) {
      existTags.value.push(tag)
    }
    newTag.value = ''
  }
}

const removeTag = (index: number) => {
  selectedLabel.value.splice(index, 1)
}
const emojiSelected = (emoji: string) => {
  state.content = state.content + emoji
}

const clickTag = (tag: string) => {
  if (!selectedLabel.value.includes(tag)){
    selectedLabel.value.push(tag)
  }
}

const router = useRouter();

const goBack = () => {
  if (props.inDrawer) {
    emit('close');
  } else {
    if (window.history.length > 1) {
      router.back();
    } else {
      navigateTo("/");
    }
  }
};

onMounted(async () => {
  if (state.id > 0) {
    const res = await useMyFetch<MemoVO>('/memo/get?id=' + state.id)
    Object.assign(state, res)
    state.showType = res.showType === 1
    const ext = JSON.parse(res.ext) as ExtDTO
    Object.assign(state.music, ext.music)
    Object.assign(state.video, ext.video)
    doubanType.value = ext.doubanBook && ext.doubanBook.title ? 'book' : 'movie'
    doubanData.value = doubanType.value === 'book' ? ext.doubanBook : ext.doubanMovie
    selectedLabel.value = res.tags ? res.tags.substring(0,res.tags.length-1).split(',') : []
    state.createdAt = dayjs(res.createdAt).format()
  }
  await loadTags()
})

// const keydown=(event:KeyboardEvent)=>{
//   if(event.key === '#'){
//     tagPopoverOpen.value = true
//   }
// }

const saveMemo = async () => {

  const doubanKey = doubanType.value === 'book' ? 'doubanBook' : 'doubanMovie'
  await useMyFetch('/memo/save', {
    id: state.id,
    content: state.content,
    ext: {
      music: state.music.id ? state.music : {},
      [doubanKey]: doubanData.value,
      video: state.video.value ? state.video : {},
    },
    pinned: state.pinned,
    showType: state.showType ? 1 : 0,
    externalFavicon: state.externalUrl ? state.externalFavicon : "",
    externalTitle: state.externalTitle,
    externalUrl: state.externalUrl,
    imgs: state.imgs.split(",").filter(Boolean),
    location: state.location,
    tags: selectedLabel.value,
    createdAt: state.createdAt || dayjs().format(),
  })
  toast.success("保存成功!")

  if (props.inDrawer) {
    emit('success')
  } else {
    await navigateTo('/')
  }
}

</script>

<style scoped>

</style>
