<template>
  <meting-js v-if="id && server && type && api && server !== 'local'" :server="server" :type="type" :id="id" :api="api"/>
  <audio v-else-if="server === 'local' && id" :src="id" controls class="w-full mb-2"></audio>
</template>

<script setup lang="ts">
import type {MusicDTO} from "@/types";

const props = defineProps<MusicDTO>()
</script>

<style scoped>

</style>