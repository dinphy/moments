<template>
  <UPopover :popper="{ arrow: true }" mode="click">
    <div class="relative group">
      <UIcon name="i-carbon-image" class="cursor-pointer w-6 h-6 text-gray-600 transition-colors duration-200" />
      <div v-if="imgList.length > 0" class="absolute -top-1 -right-1 bg-blue-500 text-white text-xs rounded-full w-4 h-4 flex items-center justify-center">
        {{ imgList.length }}
      </div>
    </div>
    
    <template #panel="{ close }">
      <div class="w-96 max-w-[90vw] bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-100 dark:border-gray-700">
        <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100">本地图片</h3>
        </div>
        
        <div class="p-4 space-y-4">
          <div class="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-6 text-center hover:border-blue-400 dark:hover:border-blue-500 transition-colors duration-200">
            <UIcon name="i-carbon-cloud-upload" class="w-12 h-12 text-gray-400 dark:text-gray-500 mx-auto mb-2 cursor-pointer hover:text-blue-500 transition-colors duration-200" @click="openFileDialog" />
            <div class="text-sm text-gray-600 dark:text-gray-300 mb-2">
              拖拽图片到此处或
              <label class="text-blue-500 dark:text-blue-400 cursor-pointer hover:underline">
                点击上传
                <input ref="fileInput" type="file" accept="image/*" multiple @change="upload" class="hidden" />
              </label>
            </div>
            <p class="text-xs text-gray-400 dark:text-gray-500">支持 JPG、PNG、WEBP 等格式</p>
          </div>

          <div v-if="isUploading" class="space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ filename }}</span>
              <span class="text-gray-500 dark:text-gray-400">{{ current }}/{{ total }}</span>
            </div>
            <div class="relative">
              <UProgress :value="progress" size="sm" />
              <span class="absolute right-0 top-0 text-xs text-gray-500 dark:text-gray-400 -mt-1">{{ progress }}%</span>
            </div>
          </div>

          <div v-if="imgList.length > 0" class="space-y-3 max-h-56 overflow-y-auto">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">已添加图片</h4>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-2">
              <div v-for="(img, index) in imgList" :key="index" class="relative group">
                <img :src="img" class="w-full h-20 object-cover rounded-md border border-gray-200 dark:border-gray-600" />
                <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-opacity duration-200 rounded-md flex items-center justify-center">
                  <UIcon name="i-carbon-close" class="w-4 h-4 text-white opacity-0 group-hover:opacity-100 cursor-pointer" @click="removeImg(index)" />
                </div>
              </div>
            </div>
          </div>

          <div class="space-y-2">
            <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300">网络图片</h4>
            <div class="flex gap-2">
              <UInput v-model="imgUrlToAdd" placeholder="输入图片URL" class="flex-1" size="sm" />
              <UButton @click="addImg" size="sm" color="primary" variant="ghost">
                <UIcon name="i-carbon-add" class="w-4 h-4" />
              </UButton>
            </div>
          </div>

          <div class="flex justify-end gap-2 pt-2">
            <UButton @click="clear(close)" size="sm" color="gray" variant="ghost">
              清空
            </UButton>
            <UButton @click="close" size="sm" color="primary">
              完成
            </UButton>
          </div>
        </div>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
import { useUpload } from "~/utils";
import { toast } from "vue-sonner";

const imgs = defineModel<string>('imgs')
const progress = ref(0)
const filename = ref('')
const total = ref(0)
const current = ref(0)
const imgUrlToAdd = ref("")
const isUploading = ref(false)

const imgList = computed(() => {
  return (imgs.value || '').split(',').filter(Boolean)
})

const upload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  
  if (!files || files.length === 0) {
    return
  }

  const fileArray = Array.from(files)
  const containsOtherFile = fileArray.some(file => !file.type.startsWith('image/'))
  
  if (containsOtherFile) {
    toast.error("只能上传图片")
    return
  }

  isUploading.value = true
  
  try {
    const result = await useUpload(files, (totalSize: number, index: number, name: string, p: number) => {
      progress.value = Math.round(p * 100)
      filename.value = name
      total.value = totalSize
      current.value = index
    })
    
    if (result && result.length) {
      toast.success("上传成功")
      imgs.value = [...imgList.value, ...result].filter(Boolean).join(',')
    }
  } finally {
    isUploading.value = false
    progress.value = 0
    // Reset file input
    target.value = ''
  }
}

const addImg = () => {
  if (!imgUrlToAdd.value) {
    toast.error("请输入图片URL");
    return
  }

  const imgsArr = imgList.value
  if (imgsArr.includes(imgUrlToAdd.value)) {
    toast.error("不能使用重复的图片地址")
    return
  }

  imgs.value = [...imgsArr, imgUrlToAdd.value].join(',')
  imgUrlToAdd.value = ''
}

const removeImg = (index: number) => {
  const imgsArr = imgList.value
  imgsArr.splice(index, 1)
  imgs.value = imgsArr.join(',')
}

const clear = (close: Function) => {
  imgs.value = ''
  close()
}

const openFileDialog = () => {
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
  if (fileInput) {
    fileInput.click()
  }
}
</script>

<style scoped>
/* Custom styles for better visual appearance */
</style>
