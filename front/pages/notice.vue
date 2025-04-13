<template>
  <Header v-bind:user="currentUser" @add-notice="showAddModal = true" />
  <UModal
    v-model="showAddModal"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center',
    }"
  >
    <div class="p-4">
      <p class="text-center text-lg font-bold mb-2">添加公告</p>
      <UForm class="space-y-4" size="sm" :state="state">
        <UFormGroup
          label="公告标题"
          name="title"
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput
            v-model="state.title"
            class="mb-2"
            placeholder="(*此项必填)"
          />
        </UFormGroup>
        <UFormGroup
          label="公告内容"
          name="content"
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput
            v-model="state.content"
            class="mb-2"
            placeholder="(*此项必填)"
          />
        </UFormGroup>
        <UFormGroup
          label="公告链接"
          name="noticeUrl"
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput
            v-model="state.noticeUrl"
            class="mb-2"
            placeholder="(*此项必填，须以 http(s):// 开头)"
          />
        </UFormGroup>
        <UFormGroup
          label="公告描述"
          name="description"
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput
            v-model="state.description"
            class="mb-2"
            placeholder="(此项非必填)"
          />
        </UFormGroup>
        <div class="flex justify-end gap-2 mt-4">
          <UButton color="white" @click="showAddModal = false">取消</UButton>
          <UButton @click="addNotice">确认添加</UButton>
        </div>
      </UForm>
    </div>
  </UModal>
  <div class="grid sm:grid-cols-2 grid-cols gap-4 p-4">
    <div
      v-for="notice in noticeList"
      :key="notice.id"
      class="bg-neutral-100 dark:bg-neutral-800 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 duration-300"
    >
      <a :href="notice.noticeUrl" target="_blank" class="block p-4">
        <div class="flex items-center gap-2 mb-2">
          <UIcon name="i-carbon-notebook" class="w-8 h-8" />
          <span class="text font-semibold">{{ notice.title }}</span>
        </div>
        <p class="text-gray-600 dark:text-gray-300 text-sm">
          {{ notice.description || "暂无描述" }}
        </p>
      </a>
    </div>
  </div>
  <div class="flex justify-center items-center mb-4 text-sm text-gray-400">
    共 {{ noticeList.length }} 条公告
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useMyFetch } from "~/utils";
import { useState } from "#app";
import type { UserVO, Notice } from "~/types";
import { toast } from "vue-sonner";

const state = reactive({
  title: "",
  content: "",
  noticeUrl: "",
  description: "",
});

const currentUser = useState<UserVO>("userinfo");
const noticeList = ref<Notice[]>([]);
const showAddModal = ref(false);

const addNotice = async () => {
  if (state.title.length === 0) {
    toast.warning("公告标题不能为空");
    return;
  }
  if (state.content.length === 0) {
    toast.warning("公告内容不能为空");
    return;
  }
  if (state.noticeUrl.length === 0) {
    toast.warning("公告链接不能为空");
    return;
  }

  if (!/^https?:\/\//.test(state.noticeUrl)) {
    toast.warning("公告链接必须以 http 或 https 开头");
    return;
  }

  if (state.description.length === 0) {
    state.description = "暂无描述";
  }

  try {
    const response = await useMyFetch("/notice/add", state);
    console.log("公告添加成功", response);
    toast.success("公告添加成功");
  } catch (error) {
    console.error("公告添加失败", error);
    toast.error("公告添加失败，请稍后重试");
  }
};
const getNoticeList = async () => {
  try {
    const response = await useMyFetch("/notice/list");
    noticeList.value = response as Notice[];
  } catch (error) {
    console.error("获取公告列表失败", error);
    toast.error("获取公告列表失败，请稍后重试");
  }
};

onMounted(() => {
  getNoticeList();
});
</script>

<style scoped></style>
