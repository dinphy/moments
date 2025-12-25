<template>
  <div
    v-if="$route.path !== '/new' && $route.path.indexOf('/edit/') < 0"
    class="header relative mb-14"
  >
    <div
      v-if="$route.path !== '/'"
      :class="{ 
        'bg-[#F5F5F5] dark:bg-[#202020] backdrop-blur-md text-[#202020] dark:text-[#F5F5F5]': y > 100, 
        'text-[#F5F5F5]': y <= 100,
        'sm:mt-14': y > 100 && showNav,
        'mt-0': true
      }"
      class="flex fixed justify-between items-center p-4 w-full md:w-[567px] top-0 z-10 transition-all duration-300"
    >
      <NuxtLink class="flex items-center" title="返回">
        <UIcon
          @click="goBack"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-4"
        />
        <span v-if="$route.path === '/sys/settings'">系统设置</span>
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
          name="i-weui-more-filled"
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
                      <UIcon class="w-4 h-4 text-green-600 dark:text-green-400" name="i-weui-pencil-outlined" />
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
                      <UIcon class="w-4 h-4 text-red-600 dark:text-red-400" name="i-weui-delete-outlined" />
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

      <span
        v-else-if="$route.path === '/sys/settings' && props.showCleanCache"
        class="flex"
        title="清理缓存"
        @click="emit('clean-cache')"
      >
        <UIcon name="i-weui-delete-outlined" class="w-5 h-5 cursor-pointer" />
      </span>

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
          name="i-weui-add-friends-outlined"
          class="w-6 h-6 cursor-pointer"
          @click="$emit('add-friend')"
        />
      </span>
    </div>
    <div
      v-if="$route.path === '/' && global.userinfo.token"
      :class="{ 
        'bg-[#F9F9F9]/95 dark:bg-[#202020]/95 backdrop-blur-md z-10': y > 100, 
        'text-[#F5F5F5]': y <= 100
      }"
      class="flex md:justify-end justify-between items-center px-4 py-3 w-full md:w-[567px] md:absolute fixed top-0 transition-all duration-300"
      title="发表"
    >
      <div class="flex-1 flex justify-center md:hidden">
        <h1 v-show="y > 100">{{ props.user.nickname }}</h1>
      </div>
      <UIcon 
        :name="y > 100 ? 'i-weui-camera-outlined' : 'i-weui-camera-filled'" 
        class="w-5 h-5 cursor-pointer" 
        @click="openMemoDrawer"
      />
    </div>

    <!-- PC端顶部导航 -->
    <div 
      v-if="showNav"
      class="hidden sm:flex fixed top-0 left-0 right-0 bg-white dark:bg-[#202020] backdrop-blur-md z-30 transition-all duration-300"
      :class="{ 'translate-y-0': y > 100, '-translate-y-full': y <= 100 }"
    >
      <div class="max-w-6xl mx-auto w-full flex items-center justify-between px-6 py-3">
        <div class="flex items-center space-x-4">
          <NuxtLink to="/" class="flex items-center space-x-2 group">
            <img
              :src="props.user.avatarUrl"
              class="avatar w-10 h-10 rounded-lg"
            />
            <span class="text-xl font-bold text-gray-800 dark:text-white">{{ props.user.nickname }}</span>
          </NuxtLink>
        </div>

        <div class="hidden md:flex items-center space-x-1">
          <div
            class="relative px-4 py-2 rounded-lg cursor-pointer transition-all duration-200"
            :class="{ 'bg-[#07C160]/10 dark:bg-[#07C160]/20 text-[#07C600]': $route.path === '/' }"
            @click="handleHomeClick"
          >
            <div class="flex items-center space-x-2">
              <UIcon name="i-carbon-activity" class="w-5 h-5" />
              <span class="font-medium">动态</span>
            </div>
          </div>

          <NuxtLink
            to="/discover"
            class="relative px-4 py-2 rounded-lg transition-all duration-200"
            :class="{ 'bg-[#07C160]/10 dark:bg-[#07C160]/20 text-[#07C600]': $route.path === '/discover' }"
          >
            <div class="flex items-center space-x-2">
              <UIcon :name="$route.path === '/discover' ? 'i-weui-discover-filled' : 'i-weui-discover-outlined'" class="w-5 h-5" />
              <span class="font-medium">发现</span>
            </div>
          </NuxtLink>

          <NuxtLink
            v-if="global.userinfo.token"
            to="/user/settings"
            class="relative px-4 py-2 rounded-lg transition-all duration-200"
            :class="{ 'bg-[#07C160]/10 dark:bg-[#07C160]/20 text-[#07C600]': $route.path === '/user/settings' }"
          >
            <div class="flex items-center space-x-2">
              <UIcon :name="$route.path === '/user/settings' ? 'i-weui-me-filled' : 'i-weui-me-outlined'" class="w-5 h-5" />
              <span class="font-medium">我的</span>
            </div>
          </NuxtLink>
          <div
            v-else
            @click="handleUserClick"
            class="relative px-4 py-2 rounded-lg cursor-pointer transition-all duration-200"
          >
            <div class="flex items-center space-x-2">
              <UIcon :name="$route.path === '/user/settings' ? 'i-weui-me-filled' : 'i-weui-me-outlined'" class="w-5 h-5" />
              <span class="font-medium">我</span>
            </div>
          </div>
        </div>

        <!-- 操作区域 -->
        <div class="flex items-center space-x-3">
          <UIcon @click="openSearchDrawer" name="i-heroicons-magnifying-glass" class="w-5 h-5 cursor-pointer" />
        </div>
      </div>
    </div>
    <img class="header-img w-full" :src="props.user.coverUrl" alt="" />
    <div class="absolute right-2 bottom-[-40px]">
      <div class="userinfo flex flex-col">
        <div class="flex flex-row items-center gap-4 justify-end">
          <NuxtLink to="/" class="username text-lg font-bold text-white hover:text-primary-200 transition-colors">
            {{ props.user.nickname }}
          </NuxtLink>
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

    <!-- 发表动态抽屉 -->
    <USlideover
      v-model="showMemoDrawer"
      :ui="{
        width: 'sm:max-w-md md:max-w-lg lg:max-w-xl w-screen',
        overlay: {
          base: 'fixed inset-0 bg-gray-900/50 backdrop-blur-sm',
        },
        background: 'bg-white dark:bg-gray-900',
        ring: '',
        rounded: '',
        shadow: 'shadow-xl',
        padding: 'p-0',
        margin: '',
        height: 'h-screen',
      }"
    >
      <div class="h-[calc(100vh-12px)] overflow-y-auto">
        <MemoEdit :inDrawer="true" @success="handleMemoSuccess" @close="closeMemoDrawer" />
      </div>
    </USlideover>

    <Search />
  </div>
</template>
<script setup lang="ts">
import { toast } from "vue-sonner";
import type { UserVO, MemoVO } from "~/types";
import { useGlobalState } from "~/store";
import { memoReloadEvent } from "~/event";
import MemoEdit from "~/components/MemoEdit.vue";
import Search from "~/components/Search.vue";

const global = useGlobalState();
const route = useRoute();
const router = useRouter();

const props = defineProps<{
  user: UserVO;
  memoItem?: MemoVO;
  showCleanCache?: boolean;
}>();
const emit = defineEmits(["add-friend", "clean-cache"]);
const mode = useColorMode();
const { y } = useWindowScroll();
const loginReg = useState<boolean>("loginReg", () => false);
const moreToolbar = ref(false);

const showMemoDrawer = ref(false);

const handleUserClick = async () => {
  if (global.value.userinfo.token) {
    await navigateTo("/user/settings");
  } else {
    loginReg.value = true;
  }
};

// 添加计算属性判断是否显示导航栏
const showNav = computed(() => {
  return ['/', '/discover', '/user/settings'].includes(route.path);
});

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

// 发表动态抽屉
const openMemoDrawer = () => {
  if (!global.value.userinfo.token) {
    loginReg.value = true;
    return;
  }
  showMemoDrawer.value = true;
};

const closeMemoDrawer = () => {
  showMemoDrawer.value = false;
};

const handleMemoSuccess = async () => {
  closeMemoDrawer();

  if (route.path === '/') {
    memoReloadEvent.emit('refresh');
  }
};

// 打开搜索抽屉
const openSearchDrawer = () => {
  const showSearchDrawer = useState<boolean>('showSearchDrawer');
  showSearchDrawer.value = true;
};
</script>

<style scoped></style>
