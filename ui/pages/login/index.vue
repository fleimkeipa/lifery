<script setup lang="ts">
import { object, string } from "yup";
import { ref, reactive } from 'vue'
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});

const { t, locale } = useI18n();

const schema = object({
  username: string().required(t('login.validation.required.username')),
  password: string().required(t('login.validation.required.password')),
});

const show = ref(false);
const errorMessage = ref('');
const isLoading = ref(false);

const state = reactive({
  username: undefined,
  password: undefined,
});

type LoginResponse = {
  type: string;
  token: string;
  username: string;
  message: string;
}

type GoogleAuthURLResponse = {
  auth_url: string;
  message: string;
}

async function onSubmit(event: any) {
  errorMessage.value = '';
  isLoading.value = true;

  const { data, error } = await useFetch<LoginResponse>("/auth/login", {
    method: "POST",
    baseURL: useRuntimeConfig().public.apiBase,
    body: event.data
  });

  isLoading.value = false;

  if (error.value) {
    console.error(error.value);

    const statusCode = error.value.statusCode;
    if (statusCode === 401) {
      errorMessage.value = t('login.error.invalidCredentials');
    } else if (statusCode && statusCode >= 500) {
      errorMessage.value = t('login.error.networkError');
    } else if (statusCode && statusCode >= 400) {
      errorMessage.value = t('login.error.general');
    } else {
      errorMessage.value = t('login.error.networkError');
    }
    return;
  }

  if (data.value) {
    if (data.value.type === 'error') {
      errorMessage.value = data.value.message || t('login.error.general');
      return;
    }

    localStorage.setItem('auth_token', data.value.token);
    localStorage.setItem('auth_type', data.value.type);
    localStorage.setItem('username', data.value.username);
    await navigateTo('/home');
  }
}

async function handleGoogleLogin() {
  isLoading.value = true;
  errorMessage.value = '';

  try {
    const { data, error } = await useFetch<GoogleAuthURLResponse>("/oauth/google/url", {
      method: "GET",
      baseURL: useRuntimeConfig().public.apiBase,
    });

    if (error.value) {
      console.error(error.value);
      errorMessage.value = t('login.error.networkError');
      return;
    }

    if (data.value && data.value.auth_url) {
      // Redirect to Google OAuth
      window.location.href = data.value.auth_url;
    }
  } catch (err) {
    console.error('Google login error:', err);
    errorMessage.value = t('login.error.networkError');
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">{{ t('login.title') }}</h1>

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

      <!-- Google Login Button -->
      <div class="mb-6">
        <UButton 
          @click="handleGoogleLogin" 
          :loading="isLoading"
          color="white" 
          variant="outline" 
          class="w-full mb-4"
          :disabled="isLoading"
        >
          <template #leading>
            <svg class="w-5 h-5" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
          </template>
          {{ t('login.googleLogin') }}
        </UButton>
      </div>

      <!-- Divider -->
      <div class="relative mb-6">
        <div class="absolute inset-0 flex items-center">
          <div class="w-full border-t border-gray-300"></div>
        </div>
        <div class="relative flex justify-center text-sm">
          <span class="px-2 bg-white text-gray-500">{{ t('login.or') }}</span>
        </div>
      </div>

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup :label="t('login.username')" name="username">
          <UInput v-model="state.username" />
        </UFormGroup>

        <UFormGroup :label="t('login.password')" name="password">
          <UInput v-model="state.password" placeholder="Password" :type="show ? 'text' : 'password'">
            <template #trailing>
              <UButton class="pointer-events-auto" color="gray" variant="link"
                :icon="show ? 'i-lucide-eye-off' : 'i-lucide-eye'"
                :aria-label="show ? 'Hide password' : 'Show password'" :aria-pressed="show" aria-controls="password"
                @click.prevent="show = !show" />
            </template>
          </UInput>
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block :loading="isLoading" :disabled="isLoading">
            {{ t('login.submit') }}
          </UButton>
        </div>

        <div class="flex w-full justify-center">
          <UButton :to="{ path: '/forgot-password', query: { locale } }" variant="link" class="text-sm">
            {{ t('login.forgotPassword') }}
          </UButton>
        </div>

        <div class="flex w-full justify-center mt-4">
          <UButton :to="{ path: '/register', query: { locale } }" variant="link" class="text-sm">
            {{ t('login.noAccount') }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template>
