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
        <span v-else-if="$route.path === '/user/settings'">个人资料</span>
        <span v-else-if="$route.path.indexOf('/tags/') >= 0">
          {{ route.params.tag || "话题专栏" }}
        </span>
        <span v-else-if="$route.path === '/friend'">友情链接</span>
        <span v-else-if="$route.path === '/discover'">发现</span>
        <span v-else-if="$route.path.indexOf('/memo/') >= 0">详情</span>
        <span v-else>{{ props.user.nickname }} 的空间</span>
      </NuxtLink>

      <!-- 自定义遮罩层 -->
      <div 
        v-if="moreToolbar" 
        class="fixed inset-0 bg-gray-900/60 backdrop-blur-sm"
        @click="moreToolbar = false"
      ></div>
      
      <UPopover
        v-if="$route.path.indexOf('/memo/') >= 0 && memoItem && (global.userinfo.id === 1 || (memoItem && global.userinfo.id === memoItem.userId))"
        v-model:open="moreToolbar"
        :popper="{ placement: 'bottom-end', strategy: 'fixed' }"
      >
        <UIcon
          name="i-solar-menu-dots-bold"
          class="mt-1 w-5 h-5 cursor-pointer"
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
        class="flex"
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

    <!-- PC端顶部导航 -->
    <div 
      class="hidden sm:flex fixed top-0 left-0 right-0 bg-white/95 dark:bg-gray-900/95 backdrop-blur-md shadow-sm z-30 transition-all duration-300"
      :class="{ 'translate-y-0': y > 100, '-translate-y-full': y <= 100 }"
    >
      <div class="max-w-6xl mx-auto w-full flex items-center justify-between px-6 py-3">
        <!-- Logo区域 -->
        <div class="flex items-center space-x-4">
          <NuxtLink to="/" class="flex items-center space-x-2 group">
            <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center shadow-md group-hover:shadow-lg transition-all duration-200">
              <UIcon name="i-carbon-activity" class="w-6 h-6 text-white" />
            </div>
            <span class="text-xl font-bold text-gray-800 dark:text-white">Moments</span>
          </NuxtLink>
        </div>

        <!-- 导航链接区域 -->
        <div class="hidden md:flex items-center space-x-1">
          <!-- 动态 -->
          <div
            class="relative px-4 py-2 rounded-lg cursor-pointer transition-all duration-200"
            :class="{ 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400': $route.path === '/' }"
            @click="handleHomeClick"
          >
            <div class="flex items-center space-x-2">
              <UIcon name="i-carbon-activity" class="w-5 h-5" />
              <span class="font-medium">动态</span>
            </div>
          </div>

          <!-- 发现 -->
          <NuxtLink
            to="/discover"
            class="relative px-4 py-2 rounded-lg transition-all duration-200"
            :class="{ 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400': $route.path === '/discover' }"
          >
            <div class="flex items-center space-x-2">
              <UIcon name="i-carbon-compass" class="w-5 h-5" />
              <span class="font-medium">发现</span>
            </div>
          </NuxtLink>

          <!-- 我 -->
          <NuxtLink
            v-if="global.userinfo.token"
            to="/user/settings"
            class="relative px-4 py-2 rounded-lg transition-all duration-200"
            :class="{ 'bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400': $route.path === '/user/settings' }"
          >
            <div class="flex items-center space-x-2">
              <UIcon name="i-carbon-user-avatar" class="w-5 h-5" />
              <span class="font-medium">我</span>
            </div>
          </NuxtLink>
          <div
            v-else
            @click="handleUserClick"
            class="relative px-4 py-2 rounded-lg cursor-pointer transition-all duration-200"
          >
            <div class="flex items-center space-x-2">
              <UIcon name="i-carbon-user-avatar" class="w-5 h-5" />
              <span class="font-medium">我</span>
            </div>
          </div>
        </div>

        <!-- 操作区域 -->
        <div class="flex items-center space-x-3">
          <!-- 主题切换 -->
          <div class="relative group" :title="modeText">
            <button
              @click="toggleMode"
              class="p-2 rounded-full transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800"
            >
              <svg
                v-if="mode.preference === 'light'"
                class="w-5 h-5 text-yellow-500"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="5"></circle>
                <line x1="12" y1="1" x2="12" y2="3"></line>
                <line x1="12" y1="21" x2="12" y2="23"></line>
                <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
                <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
                <line x1="1" y1="12" x2="3" y2="12"></line>
                <line x1="21" y1="12" x2="23" y2="12"></line>
                <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
                <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
              </svg>

              <svg
                v-else-if="mode.preference === 'dark'"
                class="w-5 h-5 text-gray-700 dark:text-gray-300"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
              </svg>

              <svg
                v-else
                class="w-5 h-5 text-gray-700 dark:text-gray-300"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="9" y1="9" x2="15" y2="9"></line>
                <line x1="9" y1="15" x2="15" y2="15"></line>
              </svg>
            </button>
          </div>

          <!-- 发表动态按钮 -->
          <NuxtLink
            v-if="global.userinfo.token"
            to="/new"
            class="relative group"
            title="发表动态"
          >
            <button class="px-4 py-2 bg-green-500 hover:bg-green-600 text-white rounded-lg flex items-center space-x-2 shadow-sm transition-all duration-200 hover:shadow-md">
              <UIcon name="i-carbon-add" class="w-5 h-5" />
              <span class="font-medium">发表动态</span>
            </button>
          </NuxtLink>
        </div>
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

// 根据当前主题模式返回对应的文案
const modeText = computed(() => {
  if (mode.preference === 'light') {
    return '切换到暗色模式';
  } else if (mode.preference === 'dark') {
    return '切换到亮色模式';
  } else {
    return '跟随系统主题';
  }
});

const handleUserClick = async () => {
  if (global.value.userinfo.token) {
    await navigateTo("/user/settings");
  } else {
    loginReg.value = true;
  }
};

const handleHomeClick = async () => {
  if (route.path === '/') {
    // 如果当前在主页，触发内容刷新事件
    const { memoReloadEvent } = await import('~/event');
    memoReloadEvent.emit('refresh');
  } else {
    // 否则导航到主页
    await navigateTo('/');
  }
};

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
