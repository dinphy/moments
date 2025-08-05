<template>
  <meting-js v-if="id && server && type && api" :server="server" :type="type" :id="id" :api="api"/>
</template>

<script setup lang="ts">
import type {MusicDTO} from "@/types";

const props = defineProps<MusicDTO>()
</script>

<style scoped>

</style>