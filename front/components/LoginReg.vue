<template>
  <UModal
    v-model="loginReg"
    :ui="{
      overlay: {
        base: 'backdrop-blur',
      },
      container: 'flex justify-center items-center max-w-72 mx-auto',
    }"
  >
  <div
        class="py-3 text-center text-lg font-sans border-b-[1px] border-neutral-[100] dark:border-neutral-800"
      >
        {{ isLogin ? "用户登录" : "注册用户" }}
      </div>
    <div class="p-5">
      <UForm
        class="space-y-4"
        size="sm"
        :state="state"
        @keyup.enter="doLoginReg"
      >
        <UFormGroup label="用户名" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup label="密码" name="password">
          <UInput type="password" v-model="state.password" />
        </UFormGroup>

        <UFormGroup v-if="!isLogin" label="重复密码" name="repeatPassword">
          <UInput type="password" v-model="state.repeatPassword" />
        </UFormGroup>

        <UButtonGroup size="sm" class="flex justify-center items-center">
          <UButton
            @click="doLoginReg"
            :disabled="pending"
            :loading="pending"
            class="px-10"
          >
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
