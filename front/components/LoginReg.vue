<template>
  <UModal
    v-model="loginReg"
    :ui="{
      overlay: {
        base: 'backdrop-blur-sm bg-gray-900/50 dark:bg-gray-950/70',
      },
      container: 'flex items-start sm:items-center justify-center min-h-screen',
      width: 'w-full max-w-xs sm:max-w-xs',
      height: 'h-auto max-h-[85dvh] sm:max-h-[80vh]',
    }"
  >
    <div class="relative bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-h-[85dvh] sm:max-h-[80vh] flex flex-col">
      <div class="px-4 py-3 sm:px-5 sm:py-3 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between">
          <h3 class="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white">
            {{ isLogin ? "欢迎回来" : "创建账户" }}
          </h3>
          <UIcon @click="loginReg = false" name="i-heroicons-x-mark" class="text-gray-400 hover:text-gray-500 p-2 cursor-pointer" />
        </div>
        <p class="mt-0.5 text-xs sm:text-sm text-gray-600 dark:text-gray-400">
          {{ isLogin ? "登录账户，继续记录美好时光" : "加入我们，开始记录美好时光" }}
        </p>
      </div>
      <div class="flex-1 overflow-y-auto">
        <div class="p-4 sm:p-5">
          <UForm
            class="space-y-4"
            size="sm"
            :state="state"
            @keyup.enter="doLoginReg"
          >
            <UFormGroup 
              name="username"
              :ui="{ 
                container: 'space-y-1'
              }"
            >
              <UInput 
                v-model="state.username" 
                placeholder="请输入账号"
                size="md"
                :ui="{ 
                  base: 'w-full text-sm',
                  rounded: 'rounded-md',
                  placeholder: 'placeholder-gray-400 dark:placeholder-gray-500'
                }"
                class="focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <template #leading>
                  <UIcon name="i-heroicons-user" class="w-4 h-4 text-gray-400" />
                </template>
              </UInput>
            </UFormGroup>
            <UFormGroup 
              name="password"
              :ui="{ 
                container: 'space-y-1'
              }"
            >
              <UInput 
                type="password" 
                v-model="state.password" 
                placeholder="请输入密码"
                size="md"
                :ui="{ 
                  base: 'w-full text-sm',
                  rounded: 'rounded-md',
                  placeholder: 'placeholder-gray-400 dark:placeholder-gray-500'
                }"
                class="focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <template #leading>
                  <UIcon name="i-heroicons-lock-closed" class="w-4 h-4 text-gray-400" />
                </template>
              </UInput>
            </UFormGroup>
            <UFormGroup 
              v-if="!isLogin" 
              name="repeatPassword"
              :ui="{ 
                container: 'space-y-1'
              }"
            >
              <UInput 
                type="password" 
                v-model="state.repeatPassword" 
                placeholder="请确认密码"
                size="md"
                :ui="{ 
                  base: 'w-full text-sm',
                  rounded: 'rounded-md',
                  placeholder: 'placeholder-gray-400 dark:placeholder-gray-500'
                }"
                class="focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
              >
                <template #leading>
                  <UIcon name="i-heroicons-lock-closed" class="w-4 h-4 text-gray-400" />
                </template>
              </UInput>
            </UFormGroup>
            <div class="space-y-4">
              <UButton
                @click="doLoginReg"
                :disabled="pending"
                :loading="pending"
                size="md"
                block
                class="rounded-md font-medium"
                :ui="{ 
                  base: 'w-full py-2.5 text-sm',
                  rounded: 'rounded-md'
                }"
              >
                {{ isLogin ? "登录" : "注册" }}
              </UButton>

              <div class="text-center">
                <UButton
                  color="gray"
                  variant="link"
                  @click="isLogin = !isLogin"
                  v-if="isLogin ? sysConfig.enableRegister : true"
                  class="text-xs sm:text-sm font-medium"
                  size="sm"
                >
                  {{ isLogin ? "还没有账户？立即注册" : "已有账户？立即登录" }}
                </UButton>
              </div>
            </div>
          </UForm>
        </div>
      </div>
    </div>
  </UModal>
</template>

<script setup lang="ts">
import type { LoginResp, SysConfigVO, UserVO } from "~/types";
import { useGlobalState } from "~/store";
import { toast } from "vue-sonner";

const sysConfig = useState<SysConfigVO>("sysConfig");
const currentUser = useState<UserVO>("userinfo");
const global = useGlobalState();
const isLogin = ref(true);
const loginReg = useState<boolean>("loginReg", () => false);
const state = reactive({
  username: "",
  password: "",
  ...(!isLogin.value && { repeatPassword: "" }),
});

const pending = ref(false);

const doLoginReg = async () => {
  pending.value = true;
  try {
    if (isLogin.value) {
      global.value.userinfo = await useMyFetch<LoginResp>("/user/login", state);
      toast.success("登录成功，跳转到首页...");
      loginReg.value = false;
      location.reload();
    } else {
      if (state.username.length < 3) {
        toast.warning("用户名最少3个字符");
        return;
      }
      if (state.password !== state.repeatPassword) {
        toast.warning("两次密码输入不一致");
        return;
      }
      await useMyFetch("/user/reg", state);
      toast.success("注册成功，请登录");
      isLogin.value = true;
    }
  } catch (warning: any) {
    if (isLogin.value) {
      toast.warning("用户不存在或密码不正确");
      return;
    } else {
      toast.warning("注册失败，用户名已存在");
      return;
    }
  } finally {
    pending.value = false;
  }
};

watch(isLogin, (newVal) => {
  if (newVal) {
    delete state.repeatPassword;
  } else {
    state.repeatPassword = "";
  }
});
</script>

<style scoped></style>
