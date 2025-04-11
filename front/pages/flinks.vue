<template>
  <Header v-bind:user="currentUser" @add-flinks="showAddLinkModal = true" />
  <div class="grid sm:grid-cols-2 grid-cols gap-4 p-4">
    <div
      v-for="(link, index) in friendLinkList"
      :key="link.url"
      class="bg-neutral-100 dark:bg-neutral-800 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 duration-300 relative"
      @mouseenter="hoverIndex = index"
      @mouseleave="hoverIndex = -1"
    >
      <UIcon
        v-if="global.userinfo.id === 1 && hoverIndex === index"
        name="i-carbon-trash-can"
        class="w-5 h-5 cursor-pointer absolute top-2 right-2 text-gray-500 hover:text-red-500"
        @click="
          showDeleteConfirm = true;
          deleteIndex = index;
        "
      />
      <a :href="link.url" target="_blank" class="block p-4">
        <div class="flex items-center gap-2 mb-2">
          <img
            :src="link.icon"
            alt="Friend Avatar"
            class="w-8 h-8 rounded-full"
          />
          <span class="text font-semibold">{{ link.name }}</span>
        </div>
        <p class="text-gray-600 dark:text-gray-300 text-sm">
          {{ link.desc || "暂无描述" }}
        </p>
      </a>
    </div>
  </div>
  <div class="flex justify-center item-center mb-4 text-sm text-gray-400">
    共 {{ friendLinkList.length }} 个朋友
  </div>

  <UModal
    v-model="showAddLinkModal"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center',
    }"
  >
    <div class="p-4">
      <h3 class="text-lg font-bold mb-2">添加友情链接</h3>
      <UFormGroup
        label="名称"
        name="name"
        :ui="{ label: { base: 'font-bold' } }"
      >
        <UInput v-model="newLink.name" class="mb-2" placeholder="(*此项必填)" />
      </UFormGroup>
      <UFormGroup
        label="网址"
        name="url"
        :ui="{ label: { base: 'font-bold' } }"
      >
        <UInput
          v-model="newLink.url"
          class="mb-2"
          placeholder="(*此项必填，须以 http(s):// 开头)"
        />
      </UFormGroup>
      <UFormGroup
        label="图标"
        name="icon"
        :ui="{ label: { base: 'font-bold' } }"
      >
        <UInput
          v-model="newLink.icon"
          class="mb-2"
          placeholder="(*此项必填，请输入图标链接)"
        />
      </UFormGroup>
      <UFormGroup
        label="描述"
        name="desc"
        :ui="{ label: { base: 'font-bold' } }"
      >
        <UInput v-model="newLink.desc" class="mb-2" placeholder="此项可为空" />
      </UFormGroup>
      <div class="flex justify-end gap-2 mt-4">
        <UButton color="white" @click="showAddLinkModal = false">取消</UButton>
        <UButton @click="addFriendLink">确认添加</UButton>
      </div>
    </div>
  </UModal>

  <UModal
    v-model="showDeleteConfirm"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center',
    }"
  >
    <div class="p-4">
      <h3 class="font-bold mb-2">温情提示</h3>
      <p class="text-gray-600 dark:text-gray-300 py-4">
        你确定要删除这条链接吗？
      </p>
      <div class="flex justify-end gap-2 mt-4">
        <UButton color="white" @click="showDeleteConfirm = false">取消</UButton>
        <UButton @click="confirmDelete">确认删除</UButton>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useState } from "#app";
import type { SysConfigVO, UserVO } from "~/types";
import { toast } from "vue-sonner";
import { useMyFetch } from "~/utils";
import { useGlobalState } from "~/store";

const currentUser = useState<UserVO>("userinfo");
const sysConfig = useState<SysConfigVO>("sysConfig");
const global = useGlobalState();

const friendLinkList = computed(() => {
  if (!sysConfig.value.friendLinks) {
    return [];
  }
  const lines = sysConfig.value.friendLinks.split("\n");
  return lines.map((line) => {
    const [name, url, icon, desc = ""] = line.split("|");
    return {
      name,
      url,
      icon,
      desc,
    };
  });
});

const showAddLinkModal = ref(false);
const newLink = reactive({
  name: "",
  url: "",
  icon: "",
  desc: "",
});

const hoverIndex = ref(-1);
const showDeleteConfirm = ref(false);
const deleteIndex = ref(-1);

const addFriendLink = async () => {
  if (!newLink.name || !newLink.url || !newLink.icon) {
    toast.error("请根据提示填写相关信息");
    return;
  }
  if (!newLink.url.startsWith("http")) {
    toast.error("网址必须以 http 开头");
    return;
  }
  const newLinkStr = `${newLink.name}|${newLink.url}|${newLink.icon}|${newLink.desc}`;
  const newFriendLinks = sysConfig.value.friendLinks
    ? `${sysConfig.value.friendLinks}\n${newLinkStr}`
    : newLinkStr;
  await useMyFetch("/sysConfig/save", {
    ...sysConfig.value,
    friendLinks: newFriendLinks,
  });
  toast.success("友情链接添加成功");
  showAddLinkModal.value = false;
  const res = await useMyFetch<SysConfigVO>("/sysConfig/getFull");
  if (res) {
    sysConfig.value = res;
  }
};

const confirmDelete = async () => {
  const lines = sysConfig.value.friendLinks.split("\n");
  lines.splice(deleteIndex.value, 1);
  const newFriendLinks = lines.join("\n");
  await useMyFetch("/sysConfig/save", {
    ...sysConfig.value,
    friendLinks: newFriendLinks,
  });
  toast.success("友情链接删除成功");
  showDeleteConfirm.value = false;
  const res = await useMyFetch<SysConfigVO>("/sysConfig/getFull");
  if (res) {
    sysConfig.value = res;
  }
};
</script>

<style scoped></style>
