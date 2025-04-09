<template>
  <Header v-bind:user="currentUser" @add-link-clicked="showModal = true" />
  <div class="grid sm:grid-cols-2 grid-cols gap-4 p-4">
    <div
      v-for="(link, index) in friendLinkList"
      :key="link.url"
      class="bg-neutral-100 dark:bg-neutral-800 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 duration-300 relative"
      @mouseenter="showDeleteIcon[index] = true"
      @mouseleave="showDeleteIcon[index] = false"
    >
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
      <div v-if="showDeleteIcon[index]" class="absolute top-2 right-2">
        <UIcon
          name="i-carbon-delete"
          class="text-red-500 w-4 h-4 cursor-pointer"
          @click="
            showDeleteConfirm = true;
            selectedIndex = index;
          "
        />
      </div>
    </div>
  </div>
  <div class="flex justify-center item-center mb-4 text-sm text-gray-400">
    共 {{ friendLinkList.length }} 个朋友
  </div>
  <UModal
    v-model="showModal"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center',
    }"
  >
    <div class="bg-white dark:bg-neutral-800 p-4 rounded-lg shadow-md">
      <h2 class="text-lg font-bold mb-2">添加友情链接</h2>
      <UTextarea
        v-model="newFriendLink"
        :rows="5"
        placeholder="每行示例：名称 | 网址(须以 http(s):// 开头) | 图标链接"
      />
      <div class="flex justify-end gap-2 mt-4">
        <UButton color="white" @click="showModal = false">取消</UButton>
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
    <div class="bg-white dark:bg-neutral-800 p-4 rounded-lg shadow-md">
      <h2 class="text-lg font-bold mb-2">温情提示</h2>
      <p>确定要删除这个友情链接吗？</p>
      <div class="flex justify-end gap-2 mt-4">
        <UButton color="white" @click="showDeleteConfirm = false">取消</UButton>
        <UButton @click="deleteFriendLink">确认删除</UButton>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { useState } from "#app";
import type { SysConfigVO, UserVO } from "~/types";
import { toast } from "vue-sonner";

const currentUser = useState<UserVO>("userinfo");
const sysConfig = useState<SysConfigVO>("sysConfig");

const showModal = ref(false);
const newFriendLink = ref("");
const showDeleteIcon = ref<boolean[]>([]);
const showDeleteConfirm = ref(false);
const selectedIndex = ref(-1);

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

const findInvalidFriendLink = (links: string): string | undefined => {
  const invalidLink = links
    .split("\n")
    .filter(Boolean)
    .find((line) => {
      const [name, url, icon] = line.split("|");
      if (!name || !url || !icon || !url.startsWith("http")) {
        return true;
      }
    });

  if (!invalidLink) {
    return undefined;
  }
  return invalidLink;
};

const addFriendLink = async () => {
  if (!newFriendLink.value.trim()) {
    toast.error("输入框不能为空");
    return;
  }
  const invalidLink = findInvalidFriendLink(newFriendLink.value);
  if (invalidLink) {
    toast.error(`网址格式不正确：${invalidLink}`);
    return;
  }
  const currentLinks = sysConfig.value.friendLinks || "";
  const newLinks = currentLinks
    ? `${currentLinks}\n${newFriendLink.value}`
    : newFriendLink.value;
  await useMyFetch("/sysConfig/save", { friendLinks: newLinks });
  sysConfig.value.friendLinks = newLinks;
  showModal.value = false;
  newFriendLink.value = "";
  toast.success("友情链接添加成功");
};

const deleteFriendLink = async () => {
  if (selectedIndex.value === -1) return;
  const linkList = sysConfig.value.friendLinks?.split("\n") || [];
  linkList.splice(selectedIndex.value, 1);
  const updatedLinks = linkList.join("\n");
  await useMyFetch("/sysConfig/save", { friendLinks: updatedLinks });
  sysConfig.value.friendLinks = updatedLinks;
  showDeleteConfirm.value = false;
  selectedIndex.value = -1;
  toast.success("友情链接删除成功");
};
</script>

<style scoped></style>
