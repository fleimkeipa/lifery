<script setup lang="ts">
import { number, object, string } from "yup";
import { ref, reactive } from 'vue'
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});

const { t, locale } = useI18n();

const schema = object({
  username: string()
    .required(t('register.validation.required.username'))
    .min(3, t('register.validation.min.username')),
  email: string()
    .required(t('register.validation.required.email'))
    .email(t('register.validation.email')),
  password: string()
    .required(t('register.validation.required.password'))
    .min(6, t('register.validation.min.password')),
  confirm_password: string()
    .required(t('register.validation.required.confirmPassword'))
    .test('passwords-match', t('register.validation.passwordsMatch'), function (value) {
      return this.parent.password === value
    }),
  role_id: number()
    .required('Rol seçimi gereklidir')
});

const errorMessage = ref('');

const state = reactive({
  username: '',
  email: '',
  password: '',
  confirm_password: '',
  role_id: 0
});

type RegisterResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

async function onSubmit(event: any) {
  errorMessage.value = '';
  
  const { data, error } = await useFetch<RegisterResponse>("/auth/register", {
    method: "POST",
    baseURL: useRuntimeConfig().public.apiBase,
    body: event.data
  });

  if (error.value) {
    console.error(error.value);
    
    const statusCode = error.value.statusCode;
    if (statusCode === 409) {
      errorMessage.value = t('register.error.userExists');
    } else if (statusCode === 400) {
      errorMessage.value = t('register.error.invalidData');
    } else if (statusCode && statusCode >= 500) {
      errorMessage.value = t('register.error.networkError');
    } else if (statusCode && statusCode >= 400) {
      errorMessage.value = t('register.error.general');
    } else {
      errorMessage.value = t('register.error.networkError');
    }
    return;
  }

  if (data.value) {
    if (data.value.type === 'error') {
      errorMessage.value = data.value.message || t('register.error.general');
      return;
    }
    
    localStorage.setItem('auth_token', data.value.token);
    localStorage.setItem('username', data.value.username);
    await navigateTo('/home');
  }
}
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">{{ t('register.title') }}</h1>
      
      <div v-if="errorMessage" class="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-md">
        <div class="flex items-center">
          <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
              clip-rule="evenodd"></path>
          </svg>
          {{ errorMessage }}
        </div>
      </div>
      
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup :label="t('register.username')" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup :label="t('register.email')" name="email">
          <UInput v-model="state.email" type="email" />
        </UFormGroup>

        <UFormGroup :label="t('register.password')" name="password">
          <UInput v-model="state.password" type="password" />
        </UFormGroup>

        <UFormGroup :label="t('register.confirmPassword')" name="confirm_password">
          <UInput v-model="state.confirm_password" type="password" />
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block>{{ t('register.submit') }}</UButton>
        </div>

        <div class="flex w-full justify-center mt-4">
          <UButton :to="{ path: '/login', query: { locale } }" variant="link" class="text-sm">
            {{ t('register.alreadyHaveAccount') }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>
