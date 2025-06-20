<script setup lang="ts">
import { object, string, type InferType } from "yup";
import { ref, reactive } from 'vue'
import type { FormSubmitEvent } from "@nuxt/ui";
import { useI18n } from "vue-i18n";

definePageMeta({
  layout: "not-authenticated",
  middleware: "no-auth"
});

const { t, locale } = useI18n();

const schema = object({
  email: string()
    .required(t('forgotPassword.validation.required.email'))
    .email(t('forgotPassword.validation.email')),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  email: undefined,
});

const isLoading = ref(false);
const successMessage = ref('');
const errorMessage = ref('');

type ForgotPasswordResponse = {
  message: string;
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  isLoading.value = true;
  successMessage.value = '';
  errorMessage.value = '';
  
  const payload = {
    email: state.email
  };
  
  const { data, error } = await useFetch<ForgotPasswordResponse>("/auth/forgot-password", {
    method: "POST",
    baseURL: useRuntimeConfig().public.apiBase,
    body: payload
  });

  isLoading.value = false;

  if (error.value) {
    console.error('API Error:', error.value);
    errorMessage.value = 'Şifre sıfırlama isteği gönderilirken bir hata oluştu. Lütfen tekrar deneyin.';
    return;
  }

  if (data.value) {
    successMessage.value = t('forgotPassword.successMessage');
    state.email = undefined;
  }
}
</script>

<template>
  <div class="flex h-screen w-full flex-col items-center justify-center gap-y-16">
    <h1 class="text-6xl font-bold">Lifery</h1>
    <UCard class="flex w-full max-w-sm items-center justify-center">
      <h1 class="mb-8 ml-auto text-4xl">{{ t('forgotPassword.title') }}</h1>
      
      <div v-if="successMessage" class="mb-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
        {{ successMessage }}
      </div>
      
      <div v-if="errorMessage" class="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
        {{ errorMessage }}
      </div>

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup :label="t('forgotPassword.email')" name="email">
          <UInput v-model="state.email" type="email" />
        </UFormGroup>

        <div class="flex w-full justify-center">
          <UButton type="submit" class="mt-4" block :loading="isLoading">
            {{ t('forgotPassword.submit') }}
          </UButton>
        </div>

        <div class="flex w-full justify-center mt-4">
          <UButton :to="{ path: '/login', query: { locale } }" variant="link" class="text-sm">
            {{ t('forgotPassword.backToLogin') }}
          </UButton>
        </div>
      </UForm>
    </UCard>
  </div>
</template> 