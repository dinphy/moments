<template>
  <UModal
    v-model="open"
    :ui="{
      base: 'p-6',
      container:
        'flex justify-center items-end sm:items-end md:items-end backdrop-blur',
    }"
  >
    <div class="flex justify-between items-center py-4 sm:py-6 mb-6 w-full">
      <div class="flex items-center gap-2">
        <img
          :src="global.userinfo.token ? currentUser.avatarUrl : '/avatar.webp'"
          class="avatar w-8 h-8 rounded-md"
          :alt="global.userinfo.token ? '用户头像' : '访客头像'"
        />
        <div class="flex flex-col">
          <span class="text-sm font-bold">{{
            global.userinfo.token ? currentUser.nickname : guestId
          }}</span>
          <span class="text-xs">{{
            global.userinfo.token ? currentUser.slogan : guestMood
          }}</span>
        </div>
      </div>

      <div
        class="flex items-center gap-3 text-gray-500 dark:text-white/80 bg-gray-100 dark:bg-gray-800/75 rounded-lg px-3 py-2"
      >
        <div class="flex flex-col items-center" title="切换主题">
          <span class="flex items-center">
            <svg
              v-if="mode.value === 'light'"
              class="lucide lucide-moon-star-icon cursor-pointer"
              @click="toggleMode"
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M12 3a6 6 0 0 0 9 9 9 9 0 1 1-9-9"></path>
              <path d="M20 3v4"></path>
              <path d="M22 5h-4"></path>
            </svg>

            <svg
              v-else
              class="lucide lucide-sun-icon cursor-pointer"
              @click="toggleMode"
              xmlns="http://www.w3.org/2000/svg"
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="4"></circle>
              <path d="M12 2v2"></path>
              <path d="M12 20v2"></path>
              <path d="m4.93 4.93 1.41 1.41"></path>
              <path d="m17.66 17.66 1.41 1.41"></path>
              <path d="M2 12h2"></path>
              <path d="M20 12h2"></path>
              <path d="m6.34 17.66-1.41 1.41"></path>
              <path d="m19.07 4.93-1.41 1.41"></path>
            </svg>
          </span>
        </div>

        <div
          class="flex flex-col items-center"
          v-if="!global.userinfo.token"
          @click="mobileloginReg = true"
          title="登录/注册"
        >
          <span class="flex items-center">
            <UIcon name="i-carbon-login" class="w-5 h-5 cursor-pointer" />
          </span>
        </div>
        <div class="flex flex-col items-center" v-else @click="logout" title="退出登录">

          <span class="flex items-center">
            <UIcon name="i-carbon-logout" class="w-5 h-5 cursor-pointer" />
          </span>
        </div>
      </div>
    </div>

    <div
      class="grid grid-cols-5 gap-3 py-4 sm:py-6 text-gray-500 dark:text-white/80"
    >
      <div
        v-if="global.userinfo.token"
        class="flex flex-col items-center"
        @click="navigate('/new')"
        title="发表"
      >
        <span class="flex items-center">
          <UIcon
            name="i-carbon-camera"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">发表</span>
      </div>
      <div
        v-if="$route.path !== '/user/calendar' && global.userinfo.token"
        class="flex flex-col items-center"
        @click="navigate('/user/calendar')"
        title="日历检索"
      >
        <span class="flex items-center">
          <UIcon
            name="i-jam-search-folder"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">检索</span>
      </div>
      <div
        v-if="$route.path == '/'"
        class="flex flex-col items-center"
        @click="navigate('/friend')"
        title="友情链接"
      >
        <span class="flex items-center">
          <UIcon
            name="i-carbon-friendship"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">友链</span>
      </div>
      <div
        v-if="$route.path === '/user/settings' && global.userinfo.id === 1"
        class="flex flex-col items-center"
        @click="navigate('/user/manage')"
        title="用户管理"
      >
        <span class="flex items-center">
          <UIcon
            name="i-carbon-user-multiple"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">管理</span>
      </div>
      <div
        v-if="$route.path !== '/sys/settings' && global.userinfo.id === 1"
        class="flex flex-col items-center"
        @click="navigate('/sys/settings')"
        title="系统设置"
      >
        <span class="flex items-center">
          <UIcon
            name="i-carbon-settings"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">系统</span>
      </div>
      <div
        v-if="$route.path !== '/user/settings' && global.userinfo.token"
        class="flex flex-col items-center"
        @click="navigate('/user/settings')"
        title="用户中心"
      >
        <span class="flex items-center">
          <UIcon
            name="i-carbon-user-avatar"
            class="w-6 h-6 sm:w-7 sm:h-7 cursor-pointer"
          />
        </span>
        <span class="text-xs mt-1">用户</span>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import { toast } from "vue-sonner";
import { useGlobalState } from "~/store";
import { getGuestId } from "~/utils";
import type { UserVO } from "~/types";

const global = useGlobalState();
const mode = useColorMode();
const open = useState<boolean>("sidebarOpen", () => false);
const mobileloginReg = useState<boolean>("mobileloginReg", () => false);
const currentUser = useState<UserVO>("userinfo");
const guestId = ref<string>("");

onMounted(async () => {
  if (!global.value.userinfo.token) {
    const id = await getGuestId();
    if (id) {
      guestId.value = id;
    }
  }
});
const toggleMode = () => {
  if (mode.preference === "system") {
    mode.preference = "dark";
  } else if (mode.preference === "dark") {
    mode.preference = "light";
  } else {
    mode.preference = "system";
    toast.success("显示模式将跟随系统设置");
  }
};

const navigate = async (url: string) => {
  open.value = false;
  await navigateTo(url);
};

const logout = async () => {
  open.value = false;
  global.value.userinfo = {};
  await navigateTo("/");
};

const guestNickname = computed(() => {
  // 生成随机访客ID，格式为访客+4位随机数字
  const randomId = Math.floor(1000 + Math.random() * 9000);
  return `访客${randomId}`;
});

const guestMood = computed(() => {
  // 随机心情状态数组
  const moods = [
    "探索未知世界",
    "静静欣赏美好",
    "感受生活温度",
    "享受此刻宁静",
    "发现点滴惊喜",
    "保持好奇之心",
    "记录美好瞬间",
    "品味生活诗意",
  ];
  return moods[Math.floor(Math.random() * moods.length)];
});
</script>

<style scoped></style>
