<template>
  <UPopover :ui="{base:'min-w-[300px] min-h-[300px]'}" :popper="{ arrow: true }" mode="click">
    <UIcon name="i-carbon-music" class="cursor-pointer w-6 h-6"/>
    <template #panel="{close}">
      <div class="p-4 flex flex-col gap-2 max-h-[400px] overflow-auto">
        <UTabs :items="items" class="w-full">
          <template #musicID="{ item }">
            <UFormGroup label="选择平台" :ui="{label:{base:'font-bold'}}">
              <template #hint>
                <div class="text-xs text-gray-400">
                  <ULink class="underline" target="_blank" to="https://github.com/metowolf/MetingJS">MetingJS文档</ULink>
                </div>
              </template>
              <USelectMenu v-model="server" :options="servers" value-attribute="value" option-attribute="label"
                          placeholder="选择平台"></USelectMenu>
            </UFormGroup>

            <UFormGroup label="选择类型" :ui="{label:{base:'font-bold'}}">
              <USelectMenu v-model="type" :options="types" value-attribute="value" option-attribute="label"
                          placeholder="选择平台"></USelectMenu>
            </UFormGroup>

            <UFormGroup label="ID" :ui="{label:{base:'font-bold'}}">
              <UInput v-model="id" placeholder="输入歌曲ID/播放列表ID/专辑ID"/>
            </UFormGroup>
          </template>
          <template #musicAPI="{ item }">
            <UFormGroup label="API接口地址" :ui="{label:{base:'font-bold'}}">
              <UInput v-model="api"/>
            </UFormGroup>
          </template>
          <template #localMusic="{ item }">
            <UFormGroup label="选择音频" :ui="{label:{base:'font-bold'}}">
              <template #hint>
                <div class="text-xs text-gray-400">
                  支持 MP3、WAV、OGG 等音频格式
                </div>
              </template>
              <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-4 text-center hover:border-blue-400 dark:hover:border-blue-500 transition-colors duration-200">
                <UIcon name="i-carbon-cloud-upload" class="w-8 h-8 text-gray-400 dark:text-gray-500 mx-auto mb-2 cursor-pointer hover:text-blue-500 transition-colors duration-200" @click="openAudioFileDialog" />
                <div class="text-sm text-gray-600 dark:text-gray-300 mb-2">
                  拖拽音频到此处或
                  <label class="text-blue-500 dark:text-blue-400 cursor-pointer hover:underline">
                    点击上传
                    <input ref="audioInput" type="file" accept="audio/*" @change="uploadAudio" class="hidden" />
                  </label>
                </div>
              </div>
            </UFormGroup>
            <div v-if="isLocalAudio" class="space-y-3 mt-4">
              <UFormGroup label="歌曲名称" :ui="{label:{base:'font-bold'}}">
                <template #hint>
                  <div class="text-xs text-gray-400">
                    不填写将显示为"未知名歌曲"
                  </div>
                </template>
                <UInput v-model="musicTitle" placeholder="输入歌曲名称" />
              </UFormGroup>
              
              <UFormGroup label="歌手名称" :ui="{label:{base:'font-bold'}}">
                <template #hint>
                  <div class="text-xs text-gray-400">
                    不填写将显示为"未知名歌手"
                  </div>
                </template>
                <UInput v-model="musicArtist" placeholder="输入歌手名称" />
              </UFormGroup>
            </div>
            
            <div v-if="uploadingAudio" class="space-y-2 mt-4">
              <div class="flex justify-between text-sm">
                <span class="text-gray-600 dark:text-gray-300">{{ audioFilename }}</span>
                <span class="text-gray-500 dark:text-gray-400">{{ audioProgress }}%</span>
              </div>
              <UProgress :value="audioProgress" size="sm" />
            </div>
            
            <div v-if="uploadedAudioUrl || localMusicPath" class="mt-4 p-3 bg-green-50 dark:bg-green-900/20 rounded-lg">
              <div class="text-sm text-green-700 dark:text-green-300 font-medium">
                {{ uploadedAudioUrl ? '上传成功!' : '音频路径' }}
              </div>
              <div class="text-xs text-green-600 dark:text-green-400 truncate mt-1">
                {{ uploadedAudioUrl || localMusicPath }}
              </div>
            </div>
          </template>
        </UTabs>
        <MusicPreview v-if="previewing" 
                      :id="isLocalAudio ? (uploadedAudioUrl || localMusicPath) : id"
                      :server="isLocalAudio ? 'local' : server"
                      :type="isLocalAudio ? 'song' : type"
                      :api="isLocalAudio ? '' : api"
                      :title="musicTitle"
                      :artist="musicArtist" />

        <UButtonGroup class="shadow-none mt-1 gap-1">
          <UButton color="indigo" variant="solid" @click="preview(close)" :disabled="previewLoading || uploadingAudio"
                   :loading="previewLoading">预览
          </UButton>
          <UButton @click="confirm(close)" :disabled="uploadingAudio">确定</UButton>
          <UButton color="white" @click="reset(close)" :disabled="uploadingAudio">清空</UButton>
        </UButtonGroup>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
import type {MetingMusicServer, MetingMusicType, MusicDTO} from "@/types"
import {toast} from "vue-sonner";
import { useUpload } from "~/utils";

const props = withDefaults(defineProps<MusicDTO>(), {
  id: "",
  server: "netease" as MetingMusicServer,
  type: "song" as MetingMusicType,
  api: "https://api.i-meto.com/meting/api?server=:server&type=:type&id=:id&r=:r",
  title: "",
  artist: ""
})

const id = ref<string>(props.id)
const server = ref<MetingMusicServer>(props.server)
const type = ref<MetingMusicType>(props.type)
const api = ref<string>(props.api)
const musicTitle = ref<string>('')
const musicArtist = ref<string>('')
const localMusicPath = ref<string>('')
const uploadedAudioUrl = ref('')
const emit = defineEmits(['confirm'])
const items = [{
  slot: 'localMusic',
  label: '本地音乐'
}, {
  slot: 'musicID',
  label: '在线音乐'
}, {
  slot: 'musicAPI',
  label: 'API接口'
}]

watch(props, () => {
  if (props.server === 'local' && props.id) {
    localMusicPath.value = props.id
    server.value = 'local'
    type.value = 'song'
    api.value = ''
    id.value = ''
  } else {
    id.value = props.id || ''
    server.value = props.server || 'netease'
    type.value = props.type || 'song'
    localMusicPath.value = ''
    uploadedAudioUrl.value = ''
  }
  api.value = props.api || "https://api.i-meto.com/meting/api?server=:server&type=:type&id=:id&r=:r"
  musicTitle.value = props.title || ''
  musicArtist.value = props.artist || ''
}, { immediate: true })

const previewing = ref(false)
const previewLoading = ref(false)
const uploadingAudio = ref(false)
const audioProgress = ref(0)
const audioFilename = ref('')

const isLocalAudio = computed(() => !!uploadedAudioUrl.value || !!localMusicPath.value)

const preview = (close: Function) => {
  if (isLocalAudio.value) {
    previewing.value = true
    return
  }
  
  if (!server.value || !api.value || !id.value || !type.value) {
    toast.error("请完整填写所需信息")
    return
  }
  previewing.value = false
  previewLoading.value = true
  setTimeout(() => {
    previewing.value = true
    previewLoading.value = false
  }, 500)
}
const confirm = (close: Function) => {
  const isLocal = isLocalAudio.value || localMusicPath.value
  emit('confirm', {
    id: isLocal ? (uploadedAudioUrl.value || localMusicPath.value) : id.value,
    server: isLocal ? 'local' : server.value,
    type: isLocal ? 'song' : type.value,
    api: isLocal ? '' : api.value,
    title: musicTitle.value,
    artist: musicArtist.value
  })
  close()
}
const reset = (close: Function) => {
  previewing.value = false
  id.value = ""
  server.value = "netease"
  type.value = "song"
  api.value = "https://api.i-meto.com/meting/api?server=:server&type=:type&id=:id&r=:r"
  uploadedAudioUrl.value = ""
  localMusicPath.value = ""
  musicTitle.value = ""
  musicArtist.value = ""
  emit('confirm', { id: "", server: "netease", type: "song", api: api.value, title: "", artist: "" })
  close()
}
const servers = ref([{
  value: "netease",
  label: "网易云音乐",
}, {
  value: "tencent",
  label: "QQ音乐",
}, {
  value: "kugou",
  label: "酷狗音乐",
}, {
  value: "xiami",
  label: "虾米音乐",
}, {
  value: "baidu",
  label: "百度音乐",
},])

const types = ref([{
  value: "song",
  label: "歌曲",
}, {
  value: "playlist",
  label: "播放列表",
}, {
  value: "album",
  label: "专辑",
}, {
  value: "search",
  label: "搜索",
}, {
  value: "artist",
  label: "艺术家",
},])

const audioInput = ref<HTMLInputElement>()

const openAudioFileDialog = () => {
  if (audioInput.value) {
    audioInput.value.click()
  }
}

const uploadAudio = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (!files || files.length === 0) {
    return
  }

  const file = files[0]
  if (!file.type.startsWith('audio/')) {
    toast.error("只能上传音频文件")
    return
  }

  uploadingAudio.value = true
  audioFilename.value = file.name
  audioProgress.value = 0

  try {
    const result = await useUpload(files, (totalSize: number, index: number, name: string, progress: number) => {
      audioProgress.value = Math.round(progress * 100)
    })
    
    if (result && result.length > 0) {
      uploadedAudioUrl.value = result[0]
      toast.success("音频上传成功")
    } else {
      toast.error("音频上传失败")
    }
  } catch (error) {
    toast.error(`音频上传失败: ${error}`)
  } finally {
    uploadingAudio.value = false
    target.value = ''
  }
}
</script>

<style scoped>

</style>