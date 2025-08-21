<template>
  <Header v-bind:user="currentUser" @add-friend="openAddModal" />
  <div class="bg-white dark:bg-neutral-800">
    <div class="grid sm:grid-cols-2 grid-cols gap-4 p-4">
      <div
        v-for="friend in friendList"
        :key="friend.id"
        class="bg-neutral-100 dark:bg-neutral-700 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 duration-300 relative group"
      >
        <a :href="friend.url" target="_blank" class="block p-4">
          <div class="flex items-center gap-2 mb-2">
            <img
              :src="friend.icon"
              alt="Friend Avatar"
              class="w-8 h-8 rounded-full"
            />
            <span class="text font-semibold">{{ friend.name }}</span>
          </div>
          <p class="text-gray-600 dark:text-gray-300 text-sm">
            {{ friend.desc || "暂无描述" }}
          </p>
        </a>
        <div
          v-if="globalState.userinfo.id === 1"
          class="absolute top-0 right-0 flex gap-1 backdrop-blur-sm m-2 rounded p-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200 z-10"
        >
          <UTooltip text="编辑">
            <UButton
              color="primary"
              variant="ghost"
              size="xs"
              icon="i-heroicons-pencil-square"
              @click="openEditModal(friend)"
              class="hover:bg-primary-50 dark:hover:bg-primary-900/20"
            />
          </UTooltip>
          <UTooltip text="删除">
            <UButton
              color="red"
              variant="ghost"
              size="xs"
              icon="i-carbon-trash-can"
              @click="showConfirmModal(friend.id)"
              class="hover:bg-red-50 dark:hover:bg-red-900/20"
            />
          </UTooltip>
        </div>
      </div>
    </div>
    <div
      class="flex justify-center items-center text-sm text-gray-400 pt-4 pb-10"
    >
      <span v-if="friendList && friendList.length">
        共 {{ friendList.length }} 个朋友
      </span>
      <span v-else class="text-gray-600 dark:text-gray-300 font-semibold">
        空空如也{{ globalState.userinfo.id === 1 ? "，请点击右上角添加" : "" }}
      </span>
    </div>
  </div>

  <UModal
    v-model="showModal"
    :ui="{
      container:
        'flex justify-center items-center backdrop-blur-sm',
    }"
  >
    <div class="p-4 sm:p-6">
      <p class="text-center text-lg font-bold mb-2">
        {{ isEditMode ? "编辑友情链接" : "添加友情链接" }}
      </p>
      <UForm class="space-y-4" size="sm" :state="friendForm">
        <UFormGroup
          label="名称"
          name="name"
          required
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput v-model="friendForm.name" class="mb-2" required />
        </UFormGroup>
        <UFormGroup
          label="图标"
          name="icon"
          required
          :ui="{ label: { base: 'font-bold' } }"
        >
          <div class="flex gap-4">
            <div class="flex-1 space-y-3">
              <UInput
                v-model="friendForm.icon"
                placeholder="输入地址或上传"
                size="sm"
              />

              <div class="flex justify-between gap-3">
                <label class="cursor-pointer group">
                  <UInput
                    type="file"
                    @change="uploadAvatar"
                    accept="image/*"
                    class="hidden"
                  />
                  <div
                    class="flex-col w-24 h-24 border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg flex items-center justify-center text-gray-400 dark:text-gray-500 group-hover:border-gray-400 dark:group-hover:border-gray-500 group-hover:text-gray-500 dark:group-hover:text-gray-400 transition-colors"
                  >
                    <svg
                      class="w-8 h-8"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M12 4v16m8-8H4"
                      ></path>
                    </svg>
                    <span class="text-xs">上传图标</span>
                  </div>
                </label>
                <UAvatar :src="friendForm.icon" size="lg" />
              </div>
            </div>
          </div>
        </UFormGroup>
        <UFormGroup
          label="网址"
          name="url"
          required
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput
            v-model="friendForm.url"
            class="mb-2"
            placeholder="须以 http(s):// 开头"
          />
        </UFormGroup>
        <UFormGroup
          label="描述"
          name="desc"
          :ui="{ label: { base: 'font-bold' } }"
        >
          <UInput v-model="friendForm.desc" class="mb-2" />
        </UFormGroup>
        <div class="flex justify-end gap-2 mt-4">
          <UButton color="white" @click="closeModal">取消</UButton>
          <UButton @click="handleSubmit">
            {{ isEditMode ? "确认更新" : "确认添加" }}
          </UButton>
        </div>
      </UForm>
    </div>
  </UModal>

  <UModal
    v-model="showDeleteModal"
    :ui="{
      container:
        'fixed top-0 left-0 right-0 bottom-0 flex justify-center items-center backdrop-blur-sm',
    }"
  >
    <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">确认删除</h3>
        </template>

        <p>
          确定要删除这个友链吗？此操作不可恢复。
        </p>

        <div class="flex justify-end space-x-2 mt-4">
          <UButton color="white" @click="cancelDelete">取消</UButton>
          <UButton @click="deleteFriend(friendIdToDelete)">确认</UButton>
        </div>
      </UCard>
  </UModal>
</template>

<script setup lang="ts">
import type { Friend, UserVO } from "~/types";
import { toast } from "vue-sonner";
import { useGlobalState } from "~/store";

const DEFAULT_FRIEND = {
  name: "",
  icon: "",
  url: "",
  desc: "",
};

const globalState = useGlobalState();
const currentUser = useState<UserVO>("userinfo");

const friendList = ref<Friend[]>([]);

const showModal = ref(false);
const isEditMode = ref(false);
const friendForm = ref({ ...DEFAULT_FRIEND });
const editingFriendId = ref<number | null>(null);

const showDeleteModal = ref(false);
const friendIdToDelete = ref<number>(0);

const openAddModal = () => {
  isEditMode.value = false;
  friendForm.value = { ...DEFAULT_FRIEND };
  showModal.value = true;
};

const openEditModal = (friend: Friend) => {
  isEditMode.value = true;
  friendForm.value = {
    name: friend.name,
    icon: friend.icon,
    url: friend.url,
    desc: friend.desc || "",
  };
  editingFriendId.value = friend.id;
  showModal.value = true;
};

const closeModal = () => {
  showModal.value = false;
  friendForm.value = { ...DEFAULT_FRIEND };
  editingFriendId.value = null;
};

const handleSubmit = async () => {
  if (!friendForm.value.name) {
    toast.warning("名称不能为空");
    return;
  }

  if (!friendForm.value.icon) {
    toast.warning("图标地址不能为空");
    return;
  }

  if (!friendForm.value.url) {
    toast.warning("网址不能为空");
    return;
  }

  if (!/^https?:\/\//.test(friendForm.value.url)) {
    toast.warning("网址必须以 http 或 https 开头");
    return;
  }

  try {
    if (isEditMode.value && editingFriendId.value) {
      const response = await useMyFetch("/friend/update", {
        id: editingFriendId.value,
        ...friendForm.value,
      });
      toast.success("友情链接更新成功");
    } else {
      const response = await useMyFetch("/friend/add", friendForm.value);
      toast.success("友情链接添加成功");
    }

    await getFriendList();
    closeModal();
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : isEditMode.value
        ? "更新友情链接失败"
        : "添加友情链接失败";
    toast.error(message);
  }
};

const getFriendList = async () => {
  try {
    const response = await useMyFetch("/friend/list");
    const typedResponse = response as { list: Friend[] };
    friendList.value = typedResponse.list;
  } catch (error) {
    friendList.value = [];
  }
};

const showConfirmModal = (id: number) => {
  friendIdToDelete.value = id;
  showDeleteModal.value = true;
};

const cancelDelete = () => {
  showDeleteModal.value = false;
};

const deleteFriend = async (id: number) => {
  try {
    await useMyFetch(`/friend/delete?id=${id}`);
    toast.success("友情链接删除成功");

    await getFriendList();
    showDeleteModal.value = false;
  } catch (error) {
    toast.error("删除友情链接失败");
  }
};

const uploadAvatar = async (files: FileList) => {
  for (let i = 0; i < files.length; i++) {
    if (files[i].type.indexOf("image") < 0) {
      toast.error("只能上传图片");
      return;
    }
  }
  const result = await useUpload(files);
  if (result.length) {
    toast.success("上传成功");
    friendForm.value.icon = result[0];
  }
};

onMounted(() => {
  getFriendList();
});
</script>

<style scoped></style>
