<template>
  <Header :user="currentUser" />
  <div class="flex flex-col p-4 my-4 dark:bg-neutral-800">
    <div class="flex justify-between items-center mb-4">
      <div class="text-lg font-bold">
        <UButtonGroup size="sm">
          <UButton
            :color="viewMode === 'card' ? 'primary' : 'gray'"
            variant="ghost"
            @click="viewMode = 'card'"
            icon="i-weui-photo-wall-outlined"
            :title="'卡片视图'"
          />
          <UButton
            :color="viewMode === 'table' ? 'primary' : 'gray'"
            variant="ghost"
            @click="viewMode = 'table'"
            icon="i-carbon-list"
            :title="'列表视图'"
          />
          <UButton
            :color="sortOrder === 'asc' ? 'primary' : 'gray'"
            variant="ghost"
            @click="toggleSort"
            :icon="
              sortOrder === 'asc'
                ? 'i-heroicons-arrow-up'
                : 'i-heroicons-arrow-down'
            "
            :title="sortOrder === 'asc' ? '升序排序' : '降序排序'"
          />
        </UButtonGroup>
      </div>
      <div class="flex items-center space-x-4">
        <UInput
          v-model="keyword"
          placeholder="用户名或昵称..."
          icon="i-weui-search-outlined"
          size="sm"
          @keyup.enter="handleSearch"
          class="w-48"
          :ui="{ icon: { trailing: { pointer: '', wrapper: 'absolute inset-y-0 right-0 flex items-center' } } }"
        >
          <template #trailing v-if="keyword">
            <UButton
              icon="i-heroicons-x-mark-20-solid"
              size="2xs"
              color="gray"
              variant="link"
              :padded="false"
              @click="keyword = ''; handleSearch()"
            />
          </template>
        </UInput>
      </div>
    </div>

    <div v-if="loading && users.length === 0" class="space-y-4">
      <div
        v-if="viewMode === 'card'"
        class="grid grid-cols-2 sm:grid-cols-3 gap-4"
      >
        <USkeleton v-for="i in 6" :key="i" class="h-64" />
      </div>
      <div v-else class="space-y-2">
        <USkeleton class="h-12" />
        <USkeleton v-for="i in 5" :key="i" class="h-16" />
      </div>
    </div>
    <div v-if="!loading && users.length === 0" class="flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400">
      <UIcon name="i-weui-search-outlined" class="text-4xl mb-2" />
      <p>未找到匹配的用户</p>
      <p class="text-sm mt-1">请尝试其他关键词</p>
    </div>

    <div
      v-if="viewMode === 'card'"
      class="grid grid-cols-2 sm:grid-cols-3 gap-4"
    >
      <UCard
        v-for="user in users"
        :key="user.id"
        class="hover:shadow-lg transition-shadow duration-200 relative group"
      >
        <div class="absolute top-2 right-2 z-10">
          <UTooltip text="编辑用户">
            <UButton
              color="primary"
              variant="ghost"
              size="xs"
              icon="i-weui-pencil-outlined"
              @click="openUserSettings(user)"
              class="opacity-100 group-hover:opacity-100 transition-opacity"
            />
          </UTooltip>
        </div>

        <div class="flex flex-col items-center space-y-3">
          <NuxtLink :to="'/user/' + user.id">
            <UAvatar
              :src="user.avatarUrl"
              size="xl"
              class="ring-2 ring-gray-200 dark:ring-gray-700"
            />
          </NuxtLink>
          <div class="text-center">
            <h3 class="font-semibold text-lg">
              {{ user.nickname || user.username }}
            </h3>
            <p class="text-sm text-gray-500">@{{ user.username }}</p>
          </div>
          <div class="text-xs text-gray-500">
            注册于 {{ $dayjs(user.createdAt).format("YYYY-MM-DD") }}
          </div>
          <div class="flex items-center space-x-1">
            <div
              class="w-2 h-2 rounded-full"
              :class="user.id === 1 ? 'bg-green-500' : 'bg-blue-500'"
            ></div>
            <span
              class="text-xs font-medium"
              :class="
                user.id === 1
                  ? 'text-green-700 dark:text-green-400'
                  : 'text-blue-700 dark:text-blue-400'
              "
            >
              {{ user.id === 1 ? "管理员" : "普通用户" }}
            </span>
          </div>
          <div
            class="absolute bottom-0 left-0 right-0 p-2 opacity-0 group-hover:opacity-100 transition-all duration-200 pointer-events-auto"
          >
            <UButton
              v-if="user.id !== 1"
              color="primary"
              variant="soft"
              size="sm"
              icon="i-weui-delete-outlined"
              @click="confirmDelete(user)"
              title="删除用户"
              class="w-full justify-center"
              block
            >
              删除
            </UButton>
          </div>
        </div>
      </UCard>
    </div>

    <div
      v-else-if="viewMode === 'table'"
      class="border rounded-lg overflow-hidden"
    >
      <div class="sm:hidden">
        <div
          v-for="user in users"
          :key="user.id"
          class="p-4 border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 last:border-b-0"
        >
          <div class="flex items-center space-x-3">
            <NuxtLink :to="'/user/' + user.id">
              <UAvatar
                :src="user.avatarUrl"
                size="md"
                class="ring-2 ring-gray-200 dark:ring-gray-700 flex-shrink-0"
              />
            </NuxtLink>
            <div class="flex-1 min-w-0">
              <div class="flex items-center justify-between">
                <div>
                  <div
                    class="text-sm font-medium text-gray-900 dark:text-gray-100"
                  >
                    {{ user.nickname || user.username }}
                  </div>
                  <div class="text-sm text-gray-500 dark:text-gray-400">
                    @{{ user.username }}
                  </div>
                </div>
                <span
                  :class="[
                    'inline-flex items-center px-2 py-1 rounded text-xs font-medium',
                    user.id === 1
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                      : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
                  ]"
                >
                  {{ user.id === 1 ? "管理员" : "普通用户" }}
                </span>
              </div>
              <div class="mt-2 flex items-center justify-between">
                <div class="text-xs text-gray-500">
                  注册于 {{ $dayjs(user.createdAt).format("YYYY-MM-DD") }}
                </div>
                <div class="flex space-x-1">
                  <UTooltip v-if="user.id !== 1" text="删除用户">
                    <UButton
                      color="red"
                      variant="ghost"
                      size="xs"
                      icon="i-weui-delete-outlined"
                      @click="confirmDelete(user)"
                      class="hover:bg-red-50 dark:hover:bg-red-900/20"
                    />
                  </UTooltip>
                  <UTooltip text="编辑用户">
                    <UButton
                      color="primary"
                      variant="ghost"
                      size="xs"
                      icon="i-weui-pencil-outlined"
                      @click="openUserSettings(user)"
                      class="hover:bg-primary-50 dark:hover:bg-primary-900/20"
                    />
                  </UTooltip>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="hidden sm:block">
        <UTable
          :columns="columns"
          :rows="users"
          :loading="loading"
          class="w-full"
          :ui="tableUi"
        >
          <template #id-data="{ row }">
            <span
              class="text-sm font-mono font-medium text-gray-900 dark:text-gray-100"
              >{{ row.id }}</span
            >
          </template>

          <template #username-data="{ row }">
            <div class="flex items-center">
              <NuxtLink :to="'/user/' + row.id">
                <UAvatar
                  :src="row.avatarUrl"
                  size="sm"
                  class="ring-2 ring-gray-200 dark:ring-gray-700"
                />
              </NuxtLink>
              <div class="ml-3">
                <div
                  class="text-sm font-medium text-gray-900 dark:text-gray-100"
                >
                  {{ row.nickname || row.username }}
                </div>
                <div class="text-sm text-gray-500 dark:text-gray-400">
                  <span class="flex items-center">
                    @{{ row.username || "-" }}
                    <span
                      :class="[
                        'inline-flex items-center px-2 py-0.5 rounded text-xs font-medium',
                        row.id === 1
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                          : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
                      ]"
                    >
                      {{ row.id === 1 ? "管理员" : "普通用户" }}
                    </span>
                  </span>
                </div>
              </div>
            </div>
          </template>

          <template #createdAt-data="{ row }">
            <div class="text-sm text-gray-900 dark:text-gray-100">
              <div>{{ $dayjs(row.createdAt).format("YYYY-MM-DD") }}</div>
              <div class="text-xs text-gray-500">
                {{ $dayjs(row.createdAt).format("HH:mm") }}
              </div>
            </div>
          </template>

          <template #actions-data="{ row }">
            <div class="flex items-center justify-end space-x-1">
              <UTooltip v-if="row.id !== 1" text="删除用户">
                <UButton
                  color="red"
                  variant="ghost"
                  size="xs"
                  icon="i-weui-delete-outlined"
                  @click="confirmDelete(row)"
                  class="hover:bg-red-50 dark:hover:bg-red-900/20"
                />
              </UTooltip>
              <UTooltip text="编辑用户">
                <UButton
                  color="primary"
                  variant="ghost"
                  size="xs"
                  icon="i-weui-pencil-outlined"
                  @click="openUserSettings(row)"
                  class="hover:bg-primary-50 dark:hover:bg-primary-900/20"
                />
              </UTooltip>
            </div>
          </template>
        </UTable>
      </div>
    </div>
    <div
      v-if="hasNext"
      ref="loadMoreEle"
      class="text-xs text-center text-gray-500 py-4 cursor-pointer hover:text-gray-700 dark:hover:text-gray-300 transition-colors"
      @click="loadMore"
    >
      点击加载更多
    </div>
    <div 
      v-else-if="users.length > 0" 
      class="text-xs text-center text-gray-400 dark:text-gray-500 py-4"
    >
      已经到底啦
    </div>

    <USlideover
      v-model="showUserSettingsModal"
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
      <div class="sticky top-0 z-10 flex items-center p-4 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900">
        <UIcon
          @click="showUserSettingsModal = false"
          name="i-carbon-chevron-left"
          class="w-5 h-5 cursor-pointer mr-4"
        />
        <h2 class="text-lg">编辑 - {{ settingsTargetUser?.username }}</h2>
      </div>

      <div class="h-[calc(100vh-72px)] overflow-y-auto">
        <UserSettings
          v-if="settingsTargetUser"
          :target-user="settingsTargetUser"
          :is-admin-mode="true"
          :on-save="handleUserSettingsSave"
        />
      </div>
    </USlideover>

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

        <p>
          确定要删除用户 "{{ deleteTarget?.username }}" 吗？此操作不可恢复。
        </p>

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
import type { UserVO, SysConfigVO } from "~/types";
import { toast } from "vue-sonner";
import UserSettings from "~/pages/user/settings.vue";
import { useElementVisibility } from "@vueuse/core";

const currentUser = useState<UserVO>("userinfo");
const sysConfig = useState<SysConfigVO>("sysConfig");

const columns = [
  { key: "id", label: "ID", class: "w-12 md:w-16 hidden sm:table-cell" },
  { key: "username", label: "用户信息", class: "min-w-[200px] flex-1" },
  { key: "createdAt", label: "创建时间", class: "w-28 hidden sm:table-cell" },
  { key: "actions", label: "操作", class: "w-24 text-right" },
];

const tableUi = {
  base: "divide-y divide-gray-200 dark:divide-gray-700",
  thead: "bg-gray-50 dark:bg-gray-800/50",
  tbody: "divide-y divide-gray-200 dark:divide-gray-700",
  tr: {
    base: "hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors duration-150",
    selected: "bg-gray-100 dark:bg-gray-700",
  },
  th: {
    base: "px-3 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider",
  },
  td: {
    base: "px-3 py-4 whitespace-nowrap",
  },
};

const users = ref<UserVO[]>([]);
const loading = ref(false);
const hasNext = ref(false);
const viewMode = ref<"table" | "card">("card");

const state = reactive({
  page: 1,
  size: 12,
  sort: "desc",
  keyword: "",
});

const keyword = ref("");
const sortOrder = ref<"asc" | "desc">("desc");

const loadMoreEle = ref(null);
const targetIsVisible = useElementVisibility(loadMoreEle);

interface UserListResponse {
  list: UserVO[];
  hasNext: boolean;
  keyword?: string;
}

const showUserSettingsModal = ref(false);
const settingsTargetUser = ref<UserVO | null>(null);

const showDeleteModal = ref(false);
const deleteTarget = ref<UserVO | null>(null);
const deleting = ref(false);

const loadUsers = async () => {
  state.page = 1;
  state.sort = sortOrder.value;
  state.keyword = keyword.value;
  const res = await useMyFetch<UserListResponse>("/user/list", state);
  if (res) {
    users.value = res.list || [];
    hasNext.value = res.hasNext || false;
  }
};

const handleSearch = () => {
  loadUsers();
};

const loadMore = async () => {
  state.page = state.page + 1;
  state.sort = sortOrder.value;
  state.keyword = keyword.value;
  const res = await useMyFetch<UserListResponse>("/user/list", state);
  if (res) {
    users.value = [...users.value, ...res.list];
    hasNext.value = res.hasNext || false;
  }
};

watch(targetIsVisible, async (visible) => {
  if (visible && sysConfig.value.enableAutoLoadNextPage) {
    await loadMore();
  }
});

const confirmDelete = (user: UserVO) => {
  deleteTarget.value = user;
  showDeleteModal.value = true;
};

const doDelete = async () => {
  if (!deleteTarget.value) return;

  if (deleteTarget.value.id === 1) {
    toast.warning("不能删除管理员账号");
    showDeleteModal.value = false;
    return;
  }

  deleting.value = true;
  try {
    await useMyFetch(`/user/${deleteTarget.value.id}`, null);
    toast.success("删除成功");
    showDeleteModal.value = false;
    await loadUsers();
  } catch (error) {
    toast.error("删除失败");
  } finally {
    deleting.value = false;
  }
};

const openUserSettings = (user: UserVO) => {
  settingsTargetUser.value = user;
  showUserSettingsModal.value = true;
};

const handleUserSettingsSave = () => {
  showUserSettingsModal.value = false;
  state.page = 1;
  loadUsers();
};

const toggleSort = () => {
  sortOrder.value = sortOrder.value === "asc" ? "desc" : "asc";
  loadUsers();
};

onMounted(async () => {
  if (!currentUser.value || currentUser.value.id !== 1) {
    toast.warning("无权限访问用户管理页面");
    navigateTo("/");
    return;
  }
  loadUsers();
});
</script>
