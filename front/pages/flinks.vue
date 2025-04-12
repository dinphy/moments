<template>
  <Header v-bind:user="currentUser" />
  <div class="grid sm:grid-cols-2 grid-cols gap-4 p-4">
    <div
      v-for="link in friendLinkList"
      :key="link.url"
      class="bg-neutral-100 dark:bg-neutral-800 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 duration-300"
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
    </div>
  </div>
  <div class="flex justify-center item-center mb-4 text-sm text-gray-400">
    共 {{ friendLinkList.length }} 个朋友
  </div>
</template>

<script setup lang="ts">
import { useState } from "#app";
import type { SysConfigVO, UserVO } from "~/types";

const currentUser = useState<UserVO>("userinfo");
const sysConfig = useState<SysConfigVO>("sysConfig");

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
</script>

<style scoped></style>
