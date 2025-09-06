<template>
  <meting-js
    v-if="id && server && type && api && server !== 'local'"
    :server="server"
    :type="type"
    :id="id"
    :api="api"
  />

  <div
    v-else-if="server === 'local' && id"
    class="music-player-container w-full mb-4"
  >
    <div
      class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 overflow-hidden transition-all duration-300 hover:shadow-lg"
    >
      <div class="p-3 space-y-3">
        <div class="flex items-center space-x-3">
          <div
            class="relative w-12 h-12 rounded-lg overflow-hidden bg-gradient-to-br from-blue-400 via-purple-500 to-pink-500 flex-shrink-0 group cursor-pointer touch-manipulation hover:-translate-y-px hover:shadow-lg dark:hover:shadow-xl dark:hover:shadow-black/30 active:scale-[0.98] motion-reduce:transform-none motion-reduce:transition-none transition-all duration-200 supports-[hover:none]:min-h-12 supports-[hover:none]:min-w-12 supports-[hover:none]:h-2 pointer-coarse:min-h-12 pointer-coarse:min-w-12"
            @click="togglePlay"
          >
            <div class="absolute inset-0 flex items-center justify-center">
              <UIcon
                name="i-carbon-music"
                class="w-6 h-6 text-white opacity-60 transition-opacity duration-200"
                :class="{ 'opacity-30': isPlaying }"
              />
            </div>

            <div
              class="absolute inset-0 flex items-center justify-center bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
            >
              <div
                class="w-8 h-8 rounded-full bg-white/90 dark:bg-gray-800/90 flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-200"
              >
                <UIcon
                  :name="isPlaying ? 'i-carbon-pause' : 'i-carbon-play'"
                  class="w-4 h-4 text-gray-700 dark:text-gray-200"
                />
              </div>
            </div>

            <div
              v-if="isPlaying"
              class="absolute inset-0 border-2 border-white/40 rounded-lg animate-pulse"
            ></div>
          </div>

          <div class="flex-1 min-w-0 overflow-hidden">
            <h3
              class="text-sm sm:text-base font-medium text-gray-900 dark:text-white truncate leading-tight"
              :title="displayTitle"
            >
              {{ displayTitle }}
            </h3>
            <p class="text-xs sm:text-sm text-gray-500 dark:text-gray-400 truncate mt-0.5" :title="displayArtist">
              {{ displayArtist }}
            </p>
          </div>

          <div class="hidden sm:flex items-center space-x-2">
            <UIcon
              name="i-carbon-volume-up"
              class="w-4 h-4 text-gray-600 dark:text-gray-300"
            />
            <input
              v-model="volume"
              @input="updateVolume"
              type="range"
              min="0"
              max="1"
              step="0.01"
              class="w-12 sm:w-16 h-1 sm:h-1.5 bg-gray-200 dark:bg-gray-600 rounded-lg appearance-none cursor-pointer volume-slider"
            />
          </div>
          <button
            class="sm:hidden p-1.5 rounded-full text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100 dark:hover:bg-gray-700 transition-all duration-200"
            @click="showMobileVolumeControl = !showMobileVolumeControl"
            title="音量控制"
          >
            <UIcon
              :name="
                volume > 0.5
                  ? 'i-carbon-volume-up'
                  : volume > 0
                  ? 'i-carbon-volume-down'
                  : 'i-carbon-volume-mute'
              "
              class="w-4 h-4 sm:w-5 sm:h-5"
            />
          </button>
        </div>
        <div
          v-if="showMobileVolumeControl"
          class="sm:hidden flex items-center justify-center space-x-3 p-2 sm:p-3 md:p-4 bg-gray-50 dark:bg-gray-700 rounded-lg"
        >
          <UIcon
            name="i-carbon-volume-mute"
            class="w-4 h-4 sm:w-5 sm:h-5 text-gray-600 dark:text-gray-300"
          />
          <input
            v-model="volume"
            @input="updateVolume"
            type="range"
            min="0"
            max="1"
            step="0.01"
            class="flex-1 h-1 sm:h-1.5 md:h-2 bg-gray-200 dark:bg-gray-600 rounded-lg appearance-none cursor-pointer volume-slider"
          />
          <UIcon
            name="i-carbon-volume-up"
            class="w-4 h-4 sm:w-5 sm:h-5 text-gray-600 dark:text-gray-300"
          />
        </div>

        <div class="space-y-2">
          <div class="flex items-center space-x-2">
            <span
              class="text-xs sm:text-sm font-mono text-gray-500 dark:text-gray-400 w-8 sm:w-10 md:w-12 text-right tabular-nums flex-shrink-0"
            >
              {{ formatTime(currentTime) }}
            </span>

            <div class="flex-1 relative">
              <div
                class="h-1 bg-gray-200 dark:bg-gray-600 rounded-full overflow-hidden"
              >
                <div
                  class="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full transition-all duration-100 ease-out"
                  :style="{ width: progressPercentage + '%' }"
                ></div>
              </div>
              <input
                v-model="seekValue"
                @input="seek"
                @mousedown="isDragging = true"
                @mouseup="isDragging = false"
                @touchstart="isDragging = true"
                @touchend="isDragging = false"
                type="range"
                min="0"
                :max="duration || 100"
                step="0.1"
                class="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
              />
            </div>

            <span
              class="text-xs sm:text-sm font-mono text-gray-500 dark:text-gray-400 w-8 sm:w-10 md:w-12 tabular-nums flex-shrink-0"
            >
              {{ formatTime(duration) }}
            </span>
          </div>
        </div>
      </div>

      <audio
        ref="audioRef"
        :src="id"
        @loadedmetadata="onLoadedMetadata"
        @timeupdate="onTimeUpdate"
        @ended="onEnded"
        @play="onPlay"
        @pause="onPause"
        @error="onError"
        class="hidden"
        preload="metadata"
      ></audio>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MusicDTO } from "@/types";
import { toast } from "vue-sonner";

const props = defineProps<MusicDTO>();

const audioRef = ref<HTMLAudioElement>();
const isPlaying = ref(false);
const currentTime = ref(0);
const duration = ref(0);
const volume = ref(0.8);
const seekValue = ref(0);
const isDragging = ref(false);
const showMobileVolumeControl = ref(false);

const progressPercentage = computed(() => {
  if (!duration.value) return 0;
  return (currentTime.value / duration.value) * 100;
});

const displayTitle = computed(() => {
  if (props.server === 'local' && props.title && props.title.trim()) {
    return props.title.length > 20 ? props.title.substring(0, 20) + '...' : props.title;
  }

  if (!props.id) return "未知名歌曲";
  const fileName = props.id.split("/").pop() || props.id;
  const nameWithoutExt = fileName.replace(/\.[^/.]+$/, "");

  if (nameWithoutExt.length > 20) {
    return nameWithoutExt.substring(0, 20) + "...";
  }

  return nameWithoutExt;
});

const displayArtist = computed(() => {
  if (props.server === 'local' && props.artist && props.artist.trim()) {
    return props.artist.length > 30 ? props.artist.substring(0, 30) + '...' : props.artist;
  }

  return props.server === 'local' ? "未知名歌手" : "本地音频文件";
});

const togglePlay = () => {
  if (!audioRef.value) return;

  if (isPlaying.value) {
    audioRef.value.pause();
  } else {
    audioRef.value.play().catch((error) => {
      console.error("播放失败:", error);
      toast.error("音频播放失败");
    });
  }
};

const seek = () => {
  if (!audioRef.value || !isDragging.value) return;
  audioRef.value.currentTime = seekValue.value;
};

const updateVolume = () => {
  if (!audioRef.value) return;
  audioRef.value.volume = volume.value;
};

const formatTime = (time: number): string => {
  if (!time || isNaN(time)) return "0:00";
  const minutes = Math.floor(time / 60);
  const seconds = Math.floor(time % 60);
  return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

const onLoadedMetadata = () => {
  if (!audioRef.value) return;
  duration.value = audioRef.value.duration;
  updateVolume();
};

const onTimeUpdate = () => {
  if (!audioRef.value || isDragging.value) return;
  currentTime.value = audioRef.value.currentTime;
  seekValue.value = currentTime.value;
};

const onEnded = () => {
  isPlaying.value = false;
  currentTime.value = 0;
  seekValue.value = 0;
};

const onPlay = () => {
  isPlaying.value = true;
};

const onPause = () => {
  isPlaying.value = false;
};

const onError = (error: Event) => {
  console.error("音频加载错误:", error);
  toast.error("音频文件加载失败");
};

onMounted(() => {
  if (audioRef.value) {
    audioRef.value.volume = volume.value;
  }
});

onUnmounted(() => {
  if (audioRef.value) {
    audioRef.value.pause();
  }
});
</script>

<style scoped></style>
