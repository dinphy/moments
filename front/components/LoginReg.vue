<template>
  <UModal
    v-model="loginReg"
    :ui="{
      overlay: {
        base: 'backdrop-blur-sm bg-gray-900/50 dark:bg-gray-950/70',
      },
      container: 'flex items-center justify-center min-h-screen',
      width: 'w-full max-w-xs sm:max-w-xs',
      height: 'h-auto max-h-[85dvh] sm:max-h-[80vh]',
    }"
  >
    <div class="relative bg-white dark:bg-gray-900 rounded-xl shadow-xl w-full max-h-[85dvh] sm:max-h-[80vh] flex flex-col">
      <div class="px-4 py-3 sm:px-5 sm:py-3 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between">
          <h3 class="text-lg text-gray-900 dark:text-white">
            {{ isLogin ? "用户登录" : "注册账户" }}
          </h3>
          <UIcon @click="loginReg = false" name="i-heroicons-x-mark" class="text-gray-400 hover:text-gray-500 p-2 cursor-pointer" />
        </div>
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
                placeholder="请输入用户名"
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
            <UFormGroup 
              v-if="isLogin"
              name="captcha"
              :ui="{ 
                container: 'space-y-1'
              }"
            >
              <div class="flex gap-2">
                <UInput 
                  v-model="state.captcha" 
                  placeholder="请输入验证码"
                  size="md"
                  :ui="{ 
                    base: 'flex-1 text-sm',
                    rounded: 'rounded-md',
                    placeholder: 'placeholder-gray-400 dark:placeholder-gray-500'
                  }"
                  class="focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
                >
                  <template #leading>
                    <UIcon name="i-heroicons-shield-check" class="w-4 h-4 text-gray-400" />
                  </template>
                </UInput>
                <div 
                  @click="generateCaptcha(); state.captcha = ''"
                  class="flex items-center justify-center w-16 h-9 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md cursor-pointer hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                  :title="'点击刷新验证码'"
                >
                  <span class="text-sm font-mono font-semibold text-gray-700 dark:text-gray-300 select-none">
                    {{ captchaCode }}
                  </span>
                </div>
              </div>
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
              
              <!-- OIDC登录按钮，仅在OIDC启用时显示 -->
              <UButton
                v-if="isLogin && oidcEnabled"
                @click="loginWithOIDC"
                :disabled="oidcPending || !oidcAuthURL"
                :loading="oidcPending"
                size="md"
                block
                class="rounded-md font-medium"
                :ui="{ 
                  base: 'w-full py-2.5 text-sm',
                  rounded: 'rounded-md'
                }"
              >
                使用第三方账号登录
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
                  {{ isLogin ? "没有账户？请注册" : "已有账户？请登录" }}
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
  captcha: "",
  ...(!isLogin.value && { repeatPassword: "" }),
});

const pending = ref(false);
const captchaCode = ref("");
// OIDC相关状态
const oidcEnabled = ref(false);
const oidcAuthURL = ref("");
const oidcPending = ref(false);
const oidcState = ref("");

// 获取OIDC配置
const fetchOIDCConfig = async () => {
  try {
    oidcPending.value = true;
    const response = await useMyFetch<{ enabled: boolean; authURL?: string; state?: string }>("/oidc/config");
    if (response) {
      oidcEnabled.value = response.enabled;
      oidcAuthURL.value = response.authURL || "";
      oidcState.value = response.state || "";
    }
  } catch (error) {
    console.error("获取OIDC配置失败:", error);
    oidcEnabled.value = false;
  } finally {
    oidcPending.value = false;
  }
};

// OIDC登录
const loginWithOIDC = () => {
  if (oidcAuthURL.value) {
    // 存储state到localStorage以便回调时验证
    localStorage.setItem("oidc_state", oidcState.value);
    // 跳转到OIDC授权页面
    window.location.href = oidcAuthURL.value;
  }
};

const generateCaptcha = () => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 4; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  captchaCode.value = result;
};

const validateCaptcha = () => {
  return state.captcha.toLowerCase() === captchaCode.value.toLowerCase();
};

onMounted(() => {
  generateCaptcha();
  // 仅在登录模式下获取OIDC配置
  if (isLogin.value) {
    fetchOIDCConfig();
  }
});

const doLoginReg = async () => {
  pending.value = true
  try {
    if (isLogin.value) {
      if (!state.username) {
        toast.warning("用户名为空或不正确")
        return
      }
      if (!state.password) {
        toast.warning("密码为空或不正确")
        return
      }
      if (!validateCaptcha()) {
        toast.warning("验证码为空或不正确")
        generateCaptcha()
        state.captcha = ""
        return
      }
      
      global.value.userinfo = await useMyFetch<LoginResp>("/user/login", state);
      toast.success("登录成功，正在跳转...");
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
      generateCaptcha();
      state.captcha = "";
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
    // 切换到登录模式时获取OIDC配置
    fetchOIDCConfig();
  } else {
    state.repeatPassword = "";
    // 切换到注册模式时隐藏OIDC按钮
    oidcEnabled.value = false;
  }
  generateCaptcha();
  state.captcha = "";
});
</script>

<style scoped></style>
