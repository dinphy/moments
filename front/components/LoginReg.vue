<template>
  <UModal
    v-model="loginReg"
    :ui="{
      container:
        'fixed mx-auto max-w-[300px] top-0 left-0 right-0 bottom-0 flex justify-center items-center backdrop-blur',
    }"
  >
    <div
      class="py-5 text-center text-xl font-sans border-b-[1px] border-neutral-[100] dark:border-neutral-800"
    >
      {{ isLogin ? "用户登录" : "账号注册" }}
    </div>
    <div class="p-5">
      <UForm class="space-y-4" size="sm" :state="state" @keyup.enter="submit">
        <UFormGroup label="账号" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup label="密码" name="password">
          <UInput type="password" v-model="state.password" />
        </UFormGroup>

        <UFormGroup v-if="!isLogin" label="重复密码" name="repeatPassword">
          <UInput type="password" v-model="state.repeatPassword" />
        </UFormGroup>

        <UButtonGroup size="sm" class="flex justify-center items-center">
          <UButton @click="submit" :disabled="pending" :loading="pending">
            {{ isLogin ? "登录" : "注册" }}
          </UButton>
          <UButton
            color="gray"
            variant="solid"
            @click="isLogin = !isLogin"
            v-if="isLogin ? sysConfig.enableRegister : true"
          >
            {{ isLogin ? "去注册" : "去登录" }}
          </UButton>
        </UButtonGroup>
      </UForm>
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

const submit = async () => {
  if (!isLogin.value && state.password !== state.repeatPassword) {
    toast.warning("两次密码输入不一致");
    return;
  }

  pending.value = true;
  try {
    if (isLogin.value) {
      global.value.userinfo = await useMyFetch<LoginResp>("/user/login", state);
      toast.success("登录成功，跳转到首页...");
      loginReg.value = false;
      location.reload();
    } else {
      await useMyFetch("/user/reg", state);
      toast.success("注册成功，请登录");
      isLogin.value = true;
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
