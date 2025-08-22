<template>
  <div
    v-if="$route.path !== '/new' && $route.path.indexOf('/edit/') < 0"
    class="header relative mb-14"
  >
    <div
      v-if="$route.path !== '/'"
      :class="{ 'bg-[#4c4c4c]/80': y > 100 }"
      class="flex fixed justify-between items-center p-4 w-full md:w-[567px] text-white top-0 z-10"
    >
      <NuxtLink class="flex items-center" title="返回">
        <UIcon
          @click="goBack"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-4"
        />
        <span v-if="$route.path === '/user/calendar'">日历检索</span>
        <span v-else-if="$route.path === '/sys/settings'">系统设置</span>
        <span v-else-if="$route.path === '/user/manage'">用户管理</span>
        <span v-else-if="$route.path === '/user/settings'">用户中心</span>
        <span v-else-if="$route.path.indexOf('/tags/') >= 0">
          {{ route.params.tag || "话题专栏" }}
        </span>
        <span v-else-if="$route.path === '/friend'">友情链接</span>
        <span v-else-if="$route.path.indexOf('/memo/') >= 0">详情</span>
        <span v-else>{{ props.user.nickname }} 的空间</span>
      </NuxtLink>

      <UPopover
        v-if="$route.path.indexOf('/memo/') >= 0 && memoItem"
        v-model:open="moreToolbar"
        :popper="{ placement: 'bottom-end', strategy: 'fixed' }"
      >
        <UIcon
          v-if="global.userinfo.id === 1 || (memoItem && global.userinfo.id === memoItem.userId)"
          name="i-solar-menu-dots-bold"
          class="w-5 h-5 cursor-pointer"
        />
        <template #panel>
          <div class="w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
            <div class="px-4 py-3 border-b border-gray-100 dark:border-gray-700">
              <h3 class="text-sm text-center font-medium text-gray-700 dark:text-gray-200">操作选项</h3>
            </div>
            <div class="p-1">
              <div class="space-y-1">
                <template v-if="global.userinfo.id === 1">
                  <div
                    class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors duration-200 group"
                    @click="setPinned(memoItem.id)"
                  >
                    <div class="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center group-hover:scale-110 transition-transform duration-200 flex-shrink-0">
                      <UIcon class="w-4 h-4 text-blue-600 dark:text-blue-400" name="i-carbon-pin" />
                    </div>
                    <div class="flex-1">
                      <div class="text-sm text-gray-700 dark:text-gray-300 font-medium">
                        {{ memoItem.pinned ? "取消" : "" }}置顶
                      </div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">{{ memoItem.pinned ? '取消此置顶' : '置顶到顶部' }}</div>
                    </div>
                    <UIcon class="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-200" name="i-carbon-chevron-right" />
                  </div>
                </template>
                <template v-if="global.userinfo.id === memoItem.userId">
                  <div
                    class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors duration-200 group"
                    @click="go2Edit(memoItem.id)"
                  >
                    <div class="w-8 h-8 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center group-hover:scale-110 transition-transform duration-200 flex-shrink-0">
                      <UIcon class="w-4 h-4 text-green-600 dark:text-green-400" name="i-carbon-edit" />
                    </div>
                    <div class="flex-1">
                      <div class="text-sm text-gray-700 dark:text-gray-300 font-medium">编辑</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">修改此内容</div>
                    </div>
                    <UIcon class="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-200" name="i-carbon-chevron-right" />
                  </div>
                </template>
                <template
                  v-if="
                    global.userinfo.id === 1 ||
                    global.userinfo.id === memoItem.userId
                  "
                >
                  <div
                    class="flex items-center gap-3 p-3 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer transition-colors duration-200 group"
                    @click="confirmDelete"
                  >
                    <div class="w-8 h-8 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center group-hover:scale-110 transition-transform duration-200 flex-shrink-0">
                      <UIcon class="w-4 h-4 text-red-600 dark:text-red-400" name="i-carbon-trash-can" />
                    </div>
                    <div class="flex-1">
                      <div class="text-sm text-gray-700 dark:text-gray-300 font-medium">删除</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">永久删除此内容</div>
                    </div>
                    <UIcon class="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-200" name="i-carbon-chevron-right" />
                  </div>
                </template>
              </div>
            </div>
          </div>
        </template>
      </UPopover>

      <NuxtLink
        v-else-if="$route.path === '/user/settings' && global.userinfo.token"
        class="hidden sm:flex"
        title="登出"
        @click="logout"
      >
        <UIcon name="i-carbon-logout" class="w-5 h-5 cursor-pointer" />
      </NuxtLink>
      <span
        v-else-if="$route.path === '/friend' && global.userinfo.id === 1"
        class="flex"
      >
        <UIcon
          name="i-carbon-add"
          class="w-6 h-6 cursor-pointer"
          @click="$emit('add-friend')"
        />
      </span>
    </div>

    <div
      class="dark:bg-neutral-800 hidden sm:flex sm:absolute sm:-right-10 sm:rounded sm:p-2 sm:flex-col sm:w-fit justify-end shadow w-full flex-row top-0 p-1 gap-2 bg-white"
    >
      <svg
        v-if="mode.value === 'light'"
        class="lucide lucide-moon-star-icon cursor-pointer"
        @click="toggleMode"
        xmlns="http://www.w3.org/2000/svg"
        width="20"
        height="20"
        viewBox="0 0 24 24"
        fill="none"
        stroke="#FDE047"
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
        stroke="#FDE047"
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

      <NuxtLink v-if="global.userinfo.token" to="/new" title="发表">
        <UIcon
          name="i-carbon-camera"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/user/calendar' && global.userinfo.token"
        to="/user/calendar"
        title="日历检索"
      >
        <UIcon
          name="i-jam-search-folder"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink v-if="$route.path === '/'" to="/friend" title="友情链接">
        <UIcon
          name="i-carbon-friendship"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path === '/user/settings' && global.userinfo.id === 1"
        to="/user/manage"
        title="用户管理"
      >
        <UIcon
          name="i-carbon-user-multiple"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/sys/settings' && global.userinfo.id === 1"
        to="/sys/settings"
        title="系统设置"
      >
        <UIcon
          name="i-carbon-settings"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <NuxtLink
        v-if="$route.path !== '/user/settings' && global.userinfo.token"
        to="/user/settings"
        title="用户中心"
      >
        <UIcon
          name="i-carbon-user-avatar"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </NuxtLink>
      <div v-if="!global.userinfo.token" title="登录" @click="loginReg = true">
        <UIcon
          name="i-carbon-login"
          class="text-[#9fc84a] w-5 h-5 cursor-pointer"
        />
      </div>
    </div>
    <img class="header-img w-full" :src="props.user.coverUrl" alt="" />
    <div class="absolute right-2 bottom-[-40px]">
      <div class="userinfo flex flex-col">
        <div class="flex flex-row items-center gap-4 justify-end">
          <div class="username text-lg font-bold text-white">
            {{ props.user.nickname }}
          </div>
          <img
            :src="props.user.avatarUrl"
            class="avatar w-[70px] h-[70px] rounded-xl"
          />
        </div>
        <div class="slogon text-gray truncate w-full text-end text-xs mt-2">
          {{ props.user.slogan }}
        </div>
      </div>
    </div>
    <MessageBox v-if="$route.path === '/' && global.userinfo.token" />
    <LoginReg v-model="loginReg" />

    <UModal
      v-model="showDeleteModal"
      :ui="{
        container: 'flex justify-center items-center backdrop-blur-sm',
      }"
    >
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">确认删除</h3>
        </template>
        <p>确定要删除此内容吗？此操作不可恢复。</p>
        <div class="flex justify-end space-x-2 mt-4">
          <UButton color="gray" @click="showDeleteModal = false">取消</UButton>
          <UButton :loading="deleting" @click="doDelete"
            >确认</UButton
          >
        </div>
      </UCard>
    </UModal>
  </div>
</template>
<script setup lang="ts">
import { toast } from "vue-sonner";
import type { UserVO, MemoVO } from "~/types";
import { useGlobalState } from "~/store";
import { memoReloadEvent } from "~/event";

const global = useGlobalState();
const route = useRoute();
const router = useRouter();

const props = defineProps<{
  user: UserVO;
  memoItem?: MemoVO;
}>();
const mode = useColorMode();
const { y } = useWindowScroll();
const loginReg = useState<boolean>("loginReg", () => false);
const moreToolbar = ref(false);

// 删除相关
const showDeleteModal = ref(false);
const deleting = ref(false);

const logout = async () => {
  global.value.userinfo = {};
  await navigateTo("/");
};

const goBack = () => {
  if (window.history.length > 1) {
    router.back();
  } else {
    navigateTo("/");
  }
};

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

// 详情页面操作方法
const go2Edit = async (id: number) => {
  await navigateTo("/edit/" + id);
};

const removeMemo = async (id: number) => {
  await useMyFetch("/memo/remove?id=" + id);
  toast.success("删除成功!");
  if (route.path.startsWith("/memo/")) {
    await navigateTo("/");
  } else {
    memoReloadEvent.emit();
  }
  moreToolbar.value = false;
};

// 确认删除
const confirmDelete = () => {
  showDeleteModal.value = true;
  moreToolbar.value = false;
};

// 执行删除
const doDelete = async () => {
  if (!props.memoItem) return;
  
  deleting.value = true;
  try {
    await useMyFetch("/memo/remove?id=" + props.memoItem.id);
    toast.success("删除成功!");
    showDeleteModal.value = false;
    if (route.path.startsWith("/memo/")) {
      await navigateTo("/");
    } else {
      memoReloadEvent.emit();
    }
  } catch (error) {
    toast.error("删除失败");
  } finally {
    deleting.value = false;
  }
};

const setPinned = async (id: number) => {
  await useMyFetch("/memo/setPinned?id=" + id);
  toast.success("操作成功!");
  if (route.path.startsWith("/memo/")) {
    await navigateTo("/");
  } else {
    memoReloadEvent.emit();
  }
  moreToolbar.value = false;
};
</script>

<style scoped></style>
