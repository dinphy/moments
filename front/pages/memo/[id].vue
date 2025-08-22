<template>
  <Header v-if="memo && memo.user" v-bind:user="memo.user" v-bind:memo-item="memo"/>

  <Memo v-if="memo" v-bind:memo="memo"/>

</template>

<script setup lang="ts">
import type {MemoVO} from "~/types";
import {memoChangedEvent} from "~/event";

const route = useRoute()
const id = route.params.id as any as number
const memo = ref<MemoVO>()
const reload = async () => {
  try {
    const res = await useMyFetch<MemoVO>('/memo/get?id=' + id)
    if (!res) {
      navigateTo('/404');
      return;
    }
    memo.value = res;
  } catch (error) {
    navigateTo('/404');
  }
}

memoChangedEvent.on(async () => {
  await reload()
})
onMounted(async () => {
  await reload()
})


</script>

<style scoped>

</style>